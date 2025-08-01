package service

import (
	"context"
	"errors"
	"time"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/internal/middleware"
	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/repository"
	"github.com/fastenmind/fastener-api/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService interface
type AuthService interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error)
	Register(ctx context.Context, req *model.RegisterRequest) (*model.Account, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*middleware.JWTClaims, error)
}

type authService struct {
	accountRepo repository.AccountRepository
	cfg         *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(accountRepo repository.AccountRepository, cfg *config.Config) AuthService {
	return &authService{
		accountRepo: accountRepo,
		cfg:         cfg,
	}
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	// Get account by username
	account, err := s.accountRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		// Try email as fallback
		account, err = s.accountRepo.GetByEmail(ctx, req.Username)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}

	// Check if account is active
	if !account.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(account)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(account)
	if err != nil {
		return nil, err
	}

	// Update last login
	_ = s.accountRepo.UpdateLastLogin(ctx, account.ID)

	// Parse access token expiration
	accessDuration, _ := time.ParseDuration(s.cfg.JWT.AccessTokenExpiration)
	expiresIn := int64(accessDuration.Seconds())

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         account,
	}, nil
}

func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.Account, error) {
	// Check if username already exists
	if _, err := s.accountRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.accountRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create account
	account := &model.Account{
		CompanyID:       req.CompanyID,
		Username:        req.Username,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		FullName:        req.FullName,
		PhoneNumber:     req.PhoneNumber,
		Role:            req.Role,
		IsActive:        true,
		IsEmailVerified: false,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	// Clear password hash before returning
	account.PasswordHash = ""

	return account, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get account
	accountID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Check if account is active
	if !account.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Generate new tokens
	newAccessToken, err := s.generateAccessToken(account)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(account)
	if err != nil {
		return nil, err
	}

	// Parse access token expiration
	accessDuration, _ := time.ParseDuration(s.cfg.JWT.AccessTokenExpiration)
	expiresIn := int64(accessDuration.Seconds())

	return &model.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		User:         account,
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*middleware.JWTClaims, error) {
	return jwt.ValidateToken(token, s.cfg.JWT.SecretKey)
}

func (s *authService) generateAccessToken(account *model.Account) (string, error) {
	claims := &middleware.JWTClaims{
		UserID:    account.ID.String(),
		CompanyID: account.CompanyID.String(),
		Email:     account.Email,
		Role:      account.Role,
	}

	duration, err := time.ParseDuration(s.cfg.JWT.AccessTokenExpiration)
	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(claims, s.cfg.JWT.SecretKey, duration)
}

func (s *authService) generateRefreshToken(account *model.Account) (string, error) {
	claims := &middleware.JWTClaims{
		UserID:    account.ID.String(),
		CompanyID: account.CompanyID.String(),
		Email:     account.Email,
		Role:      account.Role,
	}

	duration, err := time.ParseDuration(s.cfg.JWT.RefreshTokenExpiration)
	if err != nil {
		return "", err
	}

	return jwt.GenerateToken(claims, s.cfg.JWT.SecretKey, duration)
}
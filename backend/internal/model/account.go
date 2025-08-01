package model

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a user account
type Account struct {
	Base
	CompanyID       uuid.UUID  `json:"company_id" db:"company_id"`
	Username        string     `json:"username" db:"username"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FullName        string     `json:"full_name" db:"full_name"`
	PhoneNumber     *string    `json:"phone_number,omitempty" db:"phone_number"`
	Role            string     `json:"role" db:"role"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsEmailVerified bool       `json:"is_email_verified" db:"is_email_verified"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	
	// Relations
	Company *Company `json:"company,omitempty"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	CompanyID   uuid.UUID `json:"company_id" validate:"required"`
	Username    string    `json:"username" validate:"required,min=3,max=50"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,min=8"`
	FullName    string    `json:"full_name" validate:"required"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Role        string    `json:"role" validate:"required,oneof=admin manager engineer sales viewer"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         *Account `json:"user"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
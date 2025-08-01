package repository

import (
	"context"
	"fmt"

	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/pkg/database"
	"github.com/google/uuid"
)

// AccountRepository interface
type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	GetByUsername(ctx context.Context, username string) (*model.Account, error)
	GetByEmail(ctx context.Context, email string) (*model.Account, error)
	List(ctx context.Context, companyID uuid.UUID, pagination *model.Pagination) ([]*model.Account, error)
	Update(ctx context.Context, account *model.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

type accountRepository struct {
	db *database.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *database.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	account.BeforeCreate()
	
	query := `
		INSERT INTO accounts (
			id, company_id, username, email, password_hash, 
			full_name, phone_number, role, is_active, is_email_verified,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`
	
	_, err := r.db.Exec(ctx, query,
		account.ID, account.CompanyID, account.Username, account.Email, account.PasswordHash,
		account.FullName, account.PhoneNumber, account.Role, account.IsActive, account.IsEmailVerified,
		account.CreatedAt, account.UpdatedAt, account.CreatedBy,
	)
	
	return err
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	account := &model.Account{}
	
	query := `
		SELECT 
			id, company_id, username, email, password_hash,
			full_name, phone_number, role, is_active, is_email_verified,
			last_login_at, created_at, updated_at, deleted_at
		FROM accounts
		WHERE id = $1 AND deleted_at IS NULL`
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID, &account.CompanyID, &account.Username, &account.Email, &account.PasswordHash,
		&account.FullName, &account.PhoneNumber, &account.Role, &account.IsActive, &account.IsEmailVerified,
		&account.LastLoginAt, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return account, nil
}

func (r *accountRepository) GetByUsername(ctx context.Context, username string) (*model.Account, error) {
	account := &model.Account{}
	
	query := `
		SELECT 
			id, company_id, username, email, password_hash,
			full_name, phone_number, role, is_active, is_email_verified,
			last_login_at, created_at, updated_at, deleted_at
		FROM accounts
		WHERE username = $1 AND deleted_at IS NULL`
	
	err := r.db.QueryRow(ctx, query, username).Scan(
		&account.ID, &account.CompanyID, &account.Username, &account.Email, &account.PasswordHash,
		&account.FullName, &account.PhoneNumber, &account.Role, &account.IsActive, &account.IsEmailVerified,
		&account.LastLoginAt, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return account, nil
}

func (r *accountRepository) GetByEmail(ctx context.Context, email string) (*model.Account, error) {
	account := &model.Account{}
	
	query := `
		SELECT 
			id, company_id, username, email, password_hash,
			full_name, phone_number, role, is_active, is_email_verified,
			last_login_at, created_at, updated_at, deleted_at
		FROM accounts
		WHERE email = $1 AND deleted_at IS NULL`
	
	err := r.db.QueryRow(ctx, query, email).Scan(
		&account.ID, &account.CompanyID, &account.Username, &account.Email, &account.PasswordHash,
		&account.FullName, &account.PhoneNumber, &account.Role, &account.IsActive, &account.IsEmailVerified,
		&account.LastLoginAt, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return account, nil
}

func (r *accountRepository) List(ctx context.Context, companyID uuid.UUID, pagination *model.Pagination) ([]*model.Account, error) {
	accounts := []*model.Account{}
	
	// Count total
	countQuery := `SELECT COUNT(*) FROM accounts WHERE company_id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRow(ctx, countQuery, companyID).Scan(&pagination.Total)
	if err != nil {
		return nil, err
	}
	
	// Get list
	query := `
		SELECT 
			id, company_id, username, email, password_hash,
			full_name, phone_number, role, is_active, is_email_verified,
			last_login_at, created_at, updated_at
		FROM accounts
		WHERE company_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(ctx, query, companyID, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		account := &model.Account{}
		err := rows.Scan(
			&account.ID, &account.CompanyID, &account.Username, &account.Email, &account.PasswordHash,
			&account.FullName, &account.PhoneNumber, &account.Role, &account.IsActive, &account.IsEmailVerified,
			&account.LastLoginAt, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	
	return accounts, nil
}

func (r *accountRepository) Update(ctx context.Context, account *model.Account) error {
	account.BeforeUpdate()
	
	query := `
		UPDATE accounts SET
			username = $2, email = $3, full_name = $4, phone_number = $5,
			role = $6, is_active = $7, is_email_verified = $8,
			updated_at = $9, updated_by = $10
		WHERE id = $1 AND deleted_at IS NULL`
	
	result, err := r.db.Exec(ctx, query,
		account.ID, account.Username, account.Email, account.FullName, account.PhoneNumber,
		account.Role, account.IsActive, account.IsEmailVerified,
		account.UpdatedAt, account.UpdatedBy,
	)
	
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found")
	}
	
	return nil
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE accounts SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found")
	}
	
	return nil
}

func (r *accountRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE accounts SET last_login_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
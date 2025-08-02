package repository

import (
	"context"
	"fmt"

	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// accountRepositoryGorm implements AccountRepository using GORM
type accountRepositoryGorm struct {
	db *gorm.DB
}

// NewAccountRepositoryGorm creates a new account repository with GORM
func NewAccountRepositoryGorm(db *gorm.DB) AccountRepository {
	return &accountRepositoryGorm{db: db}
}

func (r *accountRepositoryGorm) Create(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepositoryGorm) GetByUsername(ctx context.Context, username string) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepositoryGorm) GetByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepositoryGorm) List(ctx context.Context, companyID uuid.UUID, pagination *model.Pagination) ([]*model.Account, error) {
	var accounts []*model.Account
	query := r.db.WithContext(ctx).Where("company_id = ?", companyID)
	
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.PageSize
		query = query.Offset(offset).Limit(pagination.PageSize)
	}
	
	err := query.Find(&accounts).Error
	return accounts, err
}

func (r *accountRepositoryGorm) Update(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepositoryGorm) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Account{}, "id = ?", id).Error
}

func (r *accountRepositoryGorm) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Account{}).Where("id = ?", id).Update("last_login", gorm.Expr("NOW()")).Error
}
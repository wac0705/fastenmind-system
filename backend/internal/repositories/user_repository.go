package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/fastenmind/fastener-api/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*model.Account, error) {
	var user model.Account
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUser(id uuid.UUID) (*model.Account, error) {
	return r.GetUserByID(id)
}

func (r *UserRepository) GetUserByEmail(email string) (*model.Account, error) {
	var user model.Account
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *model.Account) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) UpdateUser(user *model.Account) error {
	return r.db.Save(user).Error
}
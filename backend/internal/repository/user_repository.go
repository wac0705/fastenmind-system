package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/fastenmind/fastener-api/internal/model"
)

// UserRepository handles user-related data operations
type UserRepository interface {
	GetUserByID(id uuid.UUID) (*model.Account, error)
	GetUserByEmail(email string) (*model.Account, error)
	CreateUser(user *model.Account) error
	UpdateUser(user *model.Account) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db interface{}) UserRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return &userRepository{}
	}
	return &userRepository{db: gormDB}
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*model.Account, error) {
	var user model.Account
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.Account, error) {
	var user model.Account
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.Account) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user *model.Account) error {
	return r.db.Save(user).Error
}
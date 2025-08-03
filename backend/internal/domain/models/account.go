package models

import (
	"time"

	"github.com/google/uuid"
)

// AccountRole represents user role in the system
type AccountRole string

const (
	RoleAdmin    AccountRole = "admin"
	RoleManager  AccountRole = "manager"
	RoleEngineer AccountRole = "engineer"
	RoleSales    AccountRole = "sales"
	RoleViewer   AccountRole = "viewer"
)

// Account represents a user account
type Account struct {
	BaseModel
	CompanyID        uuid.UUID   `gorm:"type:uuid;not null" json:"company_id"`
	Company          Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Username         string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email            string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash     string      `gorm:"type:varchar(255);not null" json:"-"`
	FullName         string      `gorm:"type:varchar(100);not null" json:"full_name"`
	PhoneNumber      string      `gorm:"type:varchar(50)" json:"phone_number"`
	Role             AccountRole `gorm:"type:varchar(20);not null" json:"role"`
	IsActive         bool        `gorm:"default:true" json:"is_active"`
	IsEmailVerified  bool        `gorm:"default:false" json:"is_email_verified"`
	LastLoginAt      *time.Time  `json:"last_login_at"`
	PasswordResetToken string    `gorm:"type:varchar(255)" json:"-"`
	PasswordResetExpiry *time.Time `json:"-"`
}

func (Account) TableName() string {
	return "accounts"
}

// HasRole checks if account has specific role
func (a *Account) HasRole(role AccountRole) bool {
	return a.Role == role
}

// IsAdmin checks if account is admin
func (a *Account) IsAdmin() bool {
	return a.Role == RoleAdmin
}

// CanManage checks if account can manage resources
func (a *Account) CanManage() bool {
	return a.Role == RoleAdmin || a.Role == RoleManager
}
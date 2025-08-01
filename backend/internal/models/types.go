package models

import (
	"time"
	
	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/model"
)

// Type alias for backward compatibility
type User = model.Account
type Account = model.Account
type Company = model.Company
type Customer = model.Customer
type Inquiry = model.Inquiry
type Product = model.Product

// Additional types that might be in models package
type QuoteActivity struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	QuoteID   uuid.UUID `json:"quote_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Action    string    `json:"action" gorm:"not null"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
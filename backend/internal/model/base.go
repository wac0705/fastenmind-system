package model

import (
	"time"

	"github.com/google/uuid"
)

// Base contains common fields for all models
type Base struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" db:"updated_by"`
}

// BeforeCreate sets default values before creating
func (b *Base) BeforeCreate() {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
}

// BeforeUpdate sets updated_at timestamp
func (b *Base) BeforeUpdate() {
	b.UpdatedAt = time.Now()
}

// Pagination represents pagination parameters
type Pagination struct {
	Page       int `json:"page" query:"page"`
	PageSize   int `json:"page_size" query:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// GetOffset calculates offset for database query
func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the page size
func (p *Pagination) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// ListResponse represents a paginated list response
type ListResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}
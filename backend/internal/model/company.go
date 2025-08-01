package model

import "github.com/google/uuid"

// Company represents a company in the system
type Company struct {
	Base
	Code            string  `json:"code" db:"code"`
	Name            string  `json:"name" db:"name"`
	NameEn          *string `json:"name_en,omitempty" db:"name_en"`
	ShortName       *string `json:"short_name,omitempty" db:"short_name"`
	TaxID           *string `json:"tax_id,omitempty" db:"tax_id"`
	Country         string  `json:"country" db:"country"`
	Address         *string `json:"address,omitempty" db:"address"`
	Phone           *string `json:"phone,omitempty" db:"phone"`
	Fax             *string `json:"fax,omitempty" db:"fax"`
	Email           *string `json:"email,omitempty" db:"email"`
	Website         *string `json:"website,omitempty" db:"website"`
	Type            string  `json:"type" db:"type"` // headquarters, subsidiary, factory
	ParentCompanyID *uuid.UUID `json:"parent_company_id,omitempty" db:"parent_company_id"`
	IsActive        bool    `json:"is_active" db:"is_active"`
	
	// Relations
	ParentCompany *Company   `json:"parent_company,omitempty"`
	Subsidiaries  []*Company `json:"subsidiaries,omitempty"`
}

// CreateCompanyRequest represents company creation request
type CreateCompanyRequest struct {
	Code            string     `json:"code" validate:"required,max=20"`
	Name            string     `json:"name" validate:"required,max=100"`
	NameEn          *string    `json:"name_en,omitempty"`
	ShortName       *string    `json:"short_name,omitempty"`
	TaxID           *string    `json:"tax_id,omitempty"`
	Country         string     `json:"country" validate:"required,len=2"`
	Address         *string    `json:"address,omitempty"`
	Phone           *string    `json:"phone,omitempty"`
	Fax             *string    `json:"fax,omitempty"`
	Email           *string    `json:"email,omitempty" validate:"omitempty,email"`
	Website         *string    `json:"website,omitempty" validate:"omitempty,url"`
	Type            string     `json:"type" validate:"required,oneof=headquarters subsidiary factory"`
	ParentCompanyID *uuid.UUID `json:"parent_company_id,omitempty"`
}

// UpdateCompanyRequest represents company update request
type UpdateCompanyRequest struct {
	Name      *string `json:"name,omitempty"`
	NameEn    *string `json:"name_en,omitempty"`
	ShortName *string `json:"short_name,omitempty"`
	TaxID     *string `json:"tax_id,omitempty"`
	Address   *string `json:"address,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Fax       *string `json:"fax,omitempty"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Website   *string `json:"website,omitempty" validate:"omitempty,url"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
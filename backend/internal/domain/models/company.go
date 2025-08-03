package models

import (
	"github.com/google/uuid"
)

// CompanyType represents the type of company
type CompanyType string

const (
	CompanyTypeHeadquarters CompanyType = "headquarters"
	CompanyTypeSubsidiary   CompanyType = "subsidiary"
	CompanyTypeFactory      CompanyType = "factory"
)

// Company represents a company entity
type Company struct {
	BaseModel
	Code            string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Name            string     `gorm:"type:varchar(100);not null" json:"name"`
	NameEn          string     `gorm:"type:varchar(100)" json:"name_en"`
	ShortName       string     `gorm:"type:varchar(50)" json:"short_name"`
	TaxID           string     `gorm:"type:varchar(50)" json:"tax_id"`
	Country         string     `gorm:"type:varchar(2);not null" json:"country"`
	Address         string     `gorm:"type:text" json:"address"`
	Phone           string     `gorm:"type:varchar(50)" json:"phone"`
	Fax             string     `gorm:"type:varchar(50)" json:"fax"`
	Email           string     `gorm:"type:varchar(100)" json:"email"`
	Website         string     `gorm:"type:varchar(200)" json:"website"`
	Type            CompanyType `gorm:"type:varchar(20);not null" json:"type"`
	ParentCompanyID *uuid.UUID `gorm:"type:uuid" json:"parent_company_id"`
	ParentCompany   *Company   `gorm:"foreignKey:ParentCompanyID" json:"parent_company,omitempty"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	
	// Relationships
	Accounts  []Account  `gorm:"foreignKey:CompanyID" json:"accounts,omitempty"`
	Customers []Customer `gorm:"foreignKey:CompanyID" json:"customers,omitempty"`
}

func (Company) TableName() string {
	return "companies"
}
package repository

import "gorm.io/gorm"

// IntegrationRepository handles integration-related data operations
type IntegrationRepository interface {
	// Add methods as needed
}

type integrationRepository struct {
	db *gorm.DB
}

// NewIntegrationRepository creates a new integration repository
func NewIntegrationRepository(db interface{}) IntegrationRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return &integrationRepository{}
	}
	return &integrationRepository{db: gormDB}
}
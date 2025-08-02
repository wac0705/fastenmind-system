package database

import (
	"github.com/fastenmind/fastener-api/internal/config"
	"gorm.io/gorm"
)

// DBWrapper provides both pgx and GORM interfaces
type DBWrapper struct {
	*DB        // pgx pool
	GormDB *gorm.DB
}

// NewWrapper creates a new database wrapper with both pgx and GORM
func NewWrapper(cfg config.DatabaseConfig) (*DBWrapper, error) {
	// Initialize pgx pool
	pgxDB, err := New(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize GORM
	gormDB, err := NewGorm(cfg)
	if err != nil {
		pgxDB.Close()
		return nil, err
	}

	return &DBWrapper{
		DB:     pgxDB,
		GormDB: gormDB,
	}, nil
}

// Close closes both database connections
func (w *DBWrapper) Close() {
	if w.DB != nil {
		w.DB.Close()
	}
	if w.GormDB != nil {
		if sqlDB, err := w.GormDB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}
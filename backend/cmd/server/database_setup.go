package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/internal/infrastructure/database"
	"github.com/fastenmind/fastener-api/internal/repository"
)

// DatabaseSetup holds all database connections and repositories
type DatabaseSetup struct {
	ReadWriteDB *database.ReadWriteDB
	Repositories *Repositories
}

// Repositories holds all repository instances
type Repositories struct {
	Inquiry *repository.InquiryRepositoryRW
	Quote   *repository.QuoteRepositoryRW
	Order   *repository.OrderRepositoryRW
}

// SetupDatabase initializes database connections with read/write separation  
func SetupDatabase(ctx context.Context) (*DatabaseSetup, error) {
	// Load database configuration
	cfg := config.New()
	dbConfig := cfg.Database
	
	// Create read/write database setup
	rwDB, err := database.NewReadWriteDB(dbConfig.Primary, dbConfig.Replicas)
	if err != nil {
		return nil, fmt.Errorf("failed to setup read/write database: %w", err)
	}
	
	log.Printf("✅ Database setup complete with read/write separation: %v", dbConfig.ReadWriteSeparation)
	
	// Create repositories
	repos := &Repositories{
		Inquiry: repository.NewInquiryRepositoryRW(rwDB),
		Quote:   repository.NewQuoteRepositoryRW(rwDB),
		Order:   repository.NewOrderRepositoryRW(rwDB),
	}
	
	return &DatabaseSetup{
		ReadWriteDB:  rwDB,
		Repositories: repos,
	}, nil
}

// Close closes all database connections
func (ds *DatabaseSetup) Close() error {
	if ds.ReadWriteDB != nil {
		return ds.ReadWriteDB.Close()
	}
	return nil
}
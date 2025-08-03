package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/internal/infrastructure/database"
	"github.com/fastenmind/fastener-api/internal/repository"
	"gorm.io/gorm"
)

// DatabaseSetup holds all database connections and repositories
type DatabaseSetup struct {
	ReadWriteDB *database.ReadWriteDB
	Repositories *Repositories
}

// Repositories holds all repository instances
type Repositories struct {
	Inquiry repository.InquiryRepository
	Quote   repository.QuoteRepository
	Order   repository.OrderRepository
}

// SetupDatabase initializes database connections with read/write separation
func SetupDatabase(ctx context.Context) (*DatabaseSetup, error) {
	// Load database configuration
	dbConfig := config.LoadDatabaseConfig()
	
	// Validate configuration
	if err := dbConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}
	
	// Create database configurations
	primaryConfig := database.Config{
		DSN:             dbConfig.Primary.DSN(),
		MaxOpenConns:    dbConfig.MaxOpenConns,
		MaxIdleConns:    dbConfig.MaxIdleConns,
		ConnMaxLifetime: dbConfig.ConnMaxLifetime,
		SlowQueryThreshold: 200, // milliseconds
	}
	
	// Setup based on configuration
	var rwDB *database.ReadWriteDB
	var err error
	
	if dbConfig.ReadWriteSeparation && len(dbConfig.Replicas) > 0 {
		// Setup with read/write separation
		log.Println("Setting up database with read/write separation...")
		
		// Create replica configurations
		replicaConfigs := make([]database.Config, 0, len(dbConfig.Replicas))
		for _, replica := range dbConfig.Replicas {
			replicaConfig := database.Config{
				DSN:             replica.DSN(),
				MaxOpenConns:    dbConfig.MaxOpenConns,
				MaxIdleConns:    dbConfig.MaxIdleConns,
				ConnMaxLifetime: dbConfig.ConnMaxLifetime,
				SlowQueryThreshold: 200,
			}
			replicaConfigs = append(replicaConfigs, replicaConfig)
		}
		
		// Create read/write database
		rwDB, err = database.NewReadWriteDB(primaryConfig, replicaConfigs)
		if err != nil {
			return nil, fmt.Errorf("failed to setup read/write database: %w", err)
		}
		
		log.Printf("Connected to primary database and %d read replicas", len(dbConfig.Replicas))
	} else {
		// Setup without read/write separation
		log.Println("Setting up database without read/write separation...")
		
		// Connect to primary only
		primaryDB, err := database.Connect(primaryConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		
		// Create read/write DB with primary only
		rwDB = &database.ReadWriteDB{}
		rwDB, err = database.NewReadWriteDB(primaryConfig, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to setup database: %w", err)
		}
		
		log.Println("Connected to primary database (no read replicas)")
	}
	
	// Run migrations on primary
	if err := runMigrations(rwDB.Write()); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	
	// Initialize repositories with read/write separation
	repos := &Repositories{
		Inquiry: repository.NewInquiryRepositoryRW(rwDB),
		Quote:   repository.NewQuoteRepositoryRW(rwDB),
		Order:   repository.NewOrderRepositoryRW(rwDB),
	}
	
	// Test connections
	if err := testConnections(ctx, rwDB); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}
	
	return &DatabaseSetup{
		ReadWriteDB:  rwDB,
		Repositories: repos,
	}, nil
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	// Add your migration logic here
	// For example:
	// return db.AutoMigrate(&models.Inquiry{}, &models.Quote{}, &models.Order{})
	
	log.Println("Running database migrations...")
	// Placeholder - implement actual migrations
	return nil
}

// testConnections tests all database connections
func testConnections(ctx context.Context, rwDB *database.ReadWriteDB) error {
	// Test write connection
	var writeResult int
	if err := rwDB.Write().WithContext(ctx).Raw("SELECT 1").Scan(&writeResult).Error; err != nil {
		return fmt.Errorf("write database test failed: %w", err)
	}
	log.Println("Write database connection: OK")
	
	// Test read connection
	var readResult int
	if err := rwDB.Read().WithContext(ctx).Raw("SELECT 1").Scan(&readResult).Error; err != nil {
		return fmt.Errorf("read database test failed: %w", err)
	}
	log.Println("Read database connection: OK")
	
	return nil
}

// Close closes all database connections
func (ds *DatabaseSetup) Close() error {
	if ds.ReadWriteDB != nil {
		return ds.ReadWriteDB.Close()
	}
	return nil
}

// HealthCheck performs health check on database connections
func (ds *DatabaseSetup) HealthCheck(ctx context.Context) error {
	// Check write DB
	if err := ds.ReadWriteDB.Write().WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return fmt.Errorf("write database unhealthy: %w", err)
	}
	
	// Check read DB
	if err := ds.ReadWriteDB.Read().WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return fmt.Errorf("read database unhealthy: %w", err)
	}
	
	return nil
}

// Example usage in main.go:
/*
func main() {
	ctx := context.Background()
	
	// Setup database
	dbSetup, err := SetupDatabase(ctx)
	if err != nil {
		log.Fatal("Failed to setup database:", err)
	}
	defer dbSetup.Close()
	
	// Initialize services with repositories
	inquiryService := service.NewInquiryService(dbSetup.Repositories.Inquiry, dbSetup.ReadWriteDB.Write())
	quoteService := service.NewQuoteService(dbSetup.Repositories.Quote, dbSetup.ReadWriteDB.Write())
	orderService := service.NewOrderService(dbSetup.Repositories.Order, dbSetup.ReadWriteDB.Write())
	
	// Initialize handlers
	inquiryHandler := rest.NewInquiryHandler(inquiryService)
	quoteHandler := rest.NewQuoteHandler(quoteService)
	orderHandler := rest.NewOrderHandler(orderService)
	
	// Setup routes and start server...
}
*/
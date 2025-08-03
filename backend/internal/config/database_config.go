package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DatabaseConfig holds database configuration with read/write separation
type DatabaseConfig struct {
	// Primary (write) database
	Primary DBConnectionConfig `json:"primary"`
	
	// Read replicas
	Replicas []DBConnectionConfig `json:"replicas"`
	
	// Connection pool settings
	MaxOpenConns    int `json:"max_open_conns"`
	MaxIdleConns    int `json:"max_idle_conns"`
	ConnMaxLifetime int `json:"conn_max_lifetime_minutes"`
	
	// Read/Write separation settings
	ReadWriteSeparation bool `json:"read_write_separation"`
	PreferPrimaryReads  bool `json:"prefer_primary_reads"`
}

// DBConnectionConfig holds configuration for a single database connection
type DBConnectionConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
	
	// Weight for load balancing (higher = more traffic)
	Weight int `json:"weight"`
}

// DSN returns the database connection string
func (c DBConnectionConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	config := &DatabaseConfig{
		Primary: DBConnectionConfig{
			Host:     getEnv("DB_PRIMARY_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PRIMARY_PORT", 5432),
			Database: getEnv("DB_PRIMARY_NAME", "fastenmind"),
			Username: getEnv("DB_PRIMARY_USER", "postgres"),
			Password: getEnv("DB_PRIMARY_PASSWORD", ""),
			SSLMode:  getEnv("DB_PRIMARY_SSL_MODE", "disable"),
			Weight:   1,
		},
		MaxOpenConns:        getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		MaxIdleConns:        getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime:     getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 60),
		ReadWriteSeparation: getEnvAsBool("DB_READ_WRITE_SEPARATION", false),
		PreferPrimaryReads:  getEnvAsBool("DB_PREFER_PRIMARY_READS", false),
	}
	
	// Load read replicas if enabled
	if config.ReadWriteSeparation {
		replicaHosts := getEnv("DB_REPLICA_HOSTS", "")
		if replicaHosts != "" {
			hosts := strings.Split(replicaHosts, ",")
			for i, host := range hosts {
				replica := DBConnectionConfig{
					Host:     strings.TrimSpace(host),
					Port:     getEnvAsInt(fmt.Sprintf("DB_REPLICA_%d_PORT", i+1), config.Primary.Port),
					Database: getEnv(fmt.Sprintf("DB_REPLICA_%d_NAME", i+1), config.Primary.Database),
					Username: getEnv(fmt.Sprintf("DB_REPLICA_%d_USER", i+1), config.Primary.Username),
					Password: getEnv(fmt.Sprintf("DB_REPLICA_%d_PASSWORD", i+1), config.Primary.Password),
					SSLMode:  getEnv(fmt.Sprintf("DB_REPLICA_%d_SSL_MODE", i+1), config.Primary.SSLMode),
					Weight:   getEnvAsInt(fmt.Sprintf("DB_REPLICA_%d_WEIGHT", i+1), 1),
				}
				config.Replicas = append(config.Replicas, replica)
			}
		}
	}
	
	return config
}

// Validate validates the database configuration
func (c *DatabaseConfig) Validate() error {
	// Validate primary connection
	if c.Primary.Host == "" {
		return fmt.Errorf("primary database host is required")
	}
	if c.Primary.Database == "" {
		return fmt.Errorf("primary database name is required")
	}
	if c.Primary.Username == "" {
		return fmt.Errorf("primary database username is required")
	}
	
	// Validate replicas if read/write separation is enabled
	if c.ReadWriteSeparation && len(c.Replicas) == 0 {
		return fmt.Errorf("read/write separation is enabled but no replicas configured")
	}
	
	for i, replica := range c.Replicas {
		if replica.Host == "" {
			return fmt.Errorf("replica %d: host is required", i+1)
		}
		if replica.Database == "" {
			return fmt.Errorf("replica %d: database name is required", i+1)
		}
		if replica.Username == "" {
			return fmt.Errorf("replica %d: username is required", i+1)
		}
	}
	
	// Validate connection pool settings
	if c.MaxOpenConns < 1 {
		return fmt.Errorf("max open connections must be at least 1")
	}
	if c.MaxIdleConns < 0 {
		return fmt.Errorf("max idle connections cannot be negative")
	}
	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("max idle connections cannot exceed max open connections")
	}
	
	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// ExampleEnvFile returns example environment variables for database configuration
func ExampleEnvFile() string {
	return `# Primary (Write) Database Configuration
DB_PRIMARY_HOST=localhost
DB_PRIMARY_PORT=5432
DB_PRIMARY_NAME=fastenmind
DB_PRIMARY_USER=postgres
DB_PRIMARY_PASSWORD=your_password
DB_PRIMARY_SSL_MODE=disable

# Database Connection Pool
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME_MINUTES=60

# Read/Write Separation
DB_READ_WRITE_SEPARATION=true
DB_PREFER_PRIMARY_READS=false

# Read Replicas (comma-separated hosts)
DB_REPLICA_HOSTS=replica1.example.com,replica2.example.com

# Replica 1 Configuration
DB_REPLICA_1_PORT=5432
DB_REPLICA_1_NAME=fastenmind
DB_REPLICA_1_USER=postgres
DB_REPLICA_1_PASSWORD=your_password
DB_REPLICA_1_SSL_MODE=disable
DB_REPLICA_1_WEIGHT=2

# Replica 2 Configuration
DB_REPLICA_2_PORT=5432
DB_REPLICA_2_NAME=fastenmind
DB_REPLICA_2_USER=postgres
DB_REPLICA_2_PASSWORD=your_password
DB_REPLICA_2_SSL_MODE=disable
DB_REPLICA_2_WEIGHT=1`
}
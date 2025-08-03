package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/fastenmind/fastener-api/pkg/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"service":   "FastenMind API",
			"version":   "1.0.0",
			"timestamp": time.Now(),
		})
	})

	// API routes
	api := e.Group("/api/v1")
	
	// Basic API info
	api.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":        "FastenMind API",
			"description": "FastenMind Fastener Manufacturing ERP System",
			"version":     "1.0.0",
			"endpoints": map[string]string{
				"health": "/health",
				"info":   "/api/v1/info",
			},
		})
	})

	// Test database connection if available
	if dbConfig := cfg.Database; dbConfig.Primary.Host != "" {
		go func() {
			if db, err := database.NewGorm(dbConfig); err == nil {
				if sqlDB, err := db.DB(); err == nil {
					if err := sqlDB.Ping(); err == nil {
						log.Println("✅ Database connection successful")
					} else {
						log.Printf("⚠️  Database ping failed: %v", err)
					}
					sqlDB.Close()
				}
			} else {
				log.Printf("⚠️  Database connection failed: %v", err)
			}
		}()
	}

	// Start server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		log.Printf("🚀 Starting FastenMind API server on %s", addr)
		log.Printf("📖 API Documentation: http://localhost%s/api/v1/info", addr)
		log.Printf("❤️  Health Check: http://localhost%s/health", addr)
		
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("❌ Server shutdown error: %v", err)
	} else {
		log.Println("✅ Server stopped gracefully")
	}
}
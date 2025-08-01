package middleware

import (
	"github.com/fastenmind/fastener-api/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORS returns a configured CORS middleware
func CORS(cfg config.CORSConfig) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		MaxAge:           86400, // 24 hours
	})
}
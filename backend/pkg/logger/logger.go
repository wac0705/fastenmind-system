package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// New creates a new Echo logger instance
func New(env string) echo.Logger {
	l := log.New("fastenmind")
	
	switch env {
	case "production":
		l.SetLevel(log.INFO)
		l.EnableColor()
	case "development":
		l.SetLevel(log.DEBUG)
		l.EnableColor()
	default:
		l.SetLevel(log.INFO)
		l.EnableColor()
	}

	return l
}
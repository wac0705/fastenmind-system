package repository

import "gorm.io/gorm"

// TradeRepository handles trade-related data operations
type TradeRepository interface {
	// Add methods as needed
}

type tradeRepository struct {
	db *gorm.DB
}

// NewTradeRepository creates a new trade repository
func NewTradeRepository(db interface{}) TradeRepository {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return &tradeRepository{}
	}
	return &tradeRepository{db: gormDB}
}
package database

import (
	"context"
	"database/sql"
	
	"gorm.io/gorm"
)

// UnitOfWork 工作單元接口
type UnitOfWork interface {
	// Begin 開始事務
	Begin(ctx context.Context) (Transaction, error)
	
	// BeginWithOptions 使用選項開始事務
	BeginWithOptions(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
	
	// Execute 在事務中執行操作
	Execute(ctx context.Context, fn func(tx Transaction) error) error
}

// Transaction 事務接口
type Transaction interface {
	// Commit 提交事務
	Commit() error
	
	// Rollback 回滾事務
	Rollback() error
	
	// GetDB 獲取事務中的數據庫連接
	GetDB() interface{}
}

// GormUnitOfWork GORM工作單元實現
type GormUnitOfWork struct {
	db *gorm.DB
}

// NewGormUnitOfWork 創建GORM工作單元
func NewGormUnitOfWork(db *gorm.DB) UnitOfWork {
	return &GormUnitOfWork{db: db}
}

// Begin 開始事務
func (uow *GormUnitOfWork) Begin(ctx context.Context) (Transaction, error) {
	tx := uow.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	
	return &GormTransaction{tx: tx}, nil
}

// BeginWithOptions 使用選項開始事務
func (uow *GormUnitOfWork) BeginWithOptions(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	tx := uow.db.WithContext(ctx).Begin(opts)
	if tx.Error != nil {
		return nil, tx.Error
	}
	
	return &GormTransaction{tx: tx}, nil
}

// Execute 在事務中執行操作
func (uow *GormUnitOfWork) Execute(ctx context.Context, fn func(tx Transaction) error) error {
	tx, err := uow.Begin(ctx)
	if err != nil {
		return err
	}
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit()
}

// GormTransaction GORM事務實現
type GormTransaction struct {
	tx *gorm.DB
}

// Commit 提交事務
func (t *GormTransaction) Commit() error {
	return t.tx.Commit().Error
}

// Rollback 回滾事務
func (t *GormTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

// GetDB 獲取事務中的數據庫連接
func (t *GormTransaction) GetDB() interface{} {
	return t.tx
}
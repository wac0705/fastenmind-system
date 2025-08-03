package repository

import (
	"fmt"
	"strings"
	
	"github.com/fastenmind/fastener-api/pkg/database"
	"gorm.io/gorm"
)

// SafeLikeQuery 安全地構建 LIKE 查詢
func SafeLikeQuery(db *gorm.DB, column string, pattern string) *gorm.DB {
	// 驗證列名
	safeColumn, err := database.SafeColumnName(column)
	if err != nil {
		// 如果列名無效，返回一個不會匹配任何結果的查詢
		return db.Where("1 = 0")
	}
	
	// 轉義 LIKE 模式中的特殊字符
	safePattern := database.EscapeLike(pattern)
	
	// 構建查詢
	return db.Where(fmt.Sprintf("%s LIKE ?", safeColumn), "%"+safePattern+"%")
}

// SafeILikeQuery 安全地構建 ILIKE 查詢（不區分大小寫）
func SafeILikeQuery(db *gorm.DB, column string, pattern string) *gorm.DB {
	// 驗證列名
	safeColumn, err := database.SafeColumnName(column)
	if err != nil {
		return db.Where("1 = 0")
	}
	
	// 轉義 LIKE 模式中的特殊字符
	safePattern := database.EscapeLike(pattern)
	
	// 構建查詢
	return db.Where(fmt.Sprintf("%s ILIKE ?", safeColumn), "%"+safePattern+"%")
}

// SafeMultiLikeQuery 安全地構建多列 LIKE 查詢
func SafeMultiLikeQuery(db *gorm.DB, columns []string, pattern string) *gorm.DB {
	if len(columns) == 0 || pattern == "" {
		return db
	}
	
	// 轉義 LIKE 模式
	safePattern := database.EscapeLike(pattern)
	likePattern := "%" + safePattern + "%"
	
	var conditions []string
	var args []interface{}
	
	for _, column := range columns {
		safeColumn, err := database.SafeColumnName(column)
		if err != nil {
			continue
		}
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", safeColumn))
		args = append(args, likePattern)
	}
	
	if len(conditions) == 0 {
		return db
	}
	
	// 使用 OR 連接多個條件
	whereClause := strings.Join(conditions, " OR ")
	return db.Where(whereClause, args...)
}

// SafeOrderBy 安全地應用排序
func SafeOrderBy(db *gorm.DB, orderBy string) *gorm.DB {
	safeOrderBy := database.SafeOrderBy(orderBy)
	
	if safeOrderBy == "" {
		return db
	}
	
	return db.Order(safeOrderBy)
}

// SafePagination 安全地應用分頁
func SafePagination(db *gorm.DB, page, pageSize int) *gorm.DB {
	// 驗證頁碼
	if page < 1 {
		page = 1
	}
	
	// 驗證並限制頁面大小
	pageSize = database.ValidateLimit(pageSize, 100)
	
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}

// BuildSafeQuery 構建安全的查詢條件
func BuildSafeQuery(db *gorm.DB, conditions map[string]interface{}) *gorm.DB {
	for column, value := range conditions {
		// 驗證列名
		safeColumn, err := database.SafeColumnName(column)
		if err != nil {
			continue
		}
		
		// 根據值的類型構建不同的查詢
		switch v := value.(type) {
		case []interface{}:
			// IN 查詢
			if len(v) > 0 {
				db = db.Where(fmt.Sprintf("%s IN ?", safeColumn), v)
			}
		case nil:
			// IS NULL 查詢
			db = db.Where(fmt.Sprintf("%s IS NULL", safeColumn))
		default:
			// 普通相等查詢
			db = db.Where(fmt.Sprintf("%s = ?", safeColumn), value)
		}
	}
	
	return db
}
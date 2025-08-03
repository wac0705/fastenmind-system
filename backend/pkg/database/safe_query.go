package database

import (
	"fmt"
	"regexp"
	"strings"
)

// SafeTableName 驗證並清理表名
func SafeTableName(table string) (string, error) {
	// 只允許字母、數字和下劃線
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(table) {
		return "", fmt.Errorf("invalid table name: %s", table)
	}
	return table, nil
}

// SafeColumnName 驗證並清理列名
func SafeColumnName(column string) (string, error) {
	// 只允許字母、數字、下劃線和點（用於表.列格式）
	if !regexp.MustCompile(`^[a-zA-Z0-9_\.]+$`).MatchString(column) {
		return "", fmt.Errorf("invalid column name: %s", column)
	}
	return column, nil
}


// BuildWhereClause 安全地構建 WHERE 子句
func BuildWhereClause(conditions map[string]interface{}) (string, []interface{}) {
	if len(conditions) == 0 {
		return "", nil
	}
	
	var clauses []string
	var args []interface{}
	
	for column, value := range conditions {
		// 驗證列名
		safeCol, err := SafeColumnName(column)
		if err != nil {
			continue // 跳過無效列名
		}
		
		clauses = append(clauses, fmt.Sprintf("%s = ?", safeCol))
		args = append(args, value)
	}
	
	if len(clauses) == 0 {
		return "", nil
	}
	
	return "WHERE " + strings.Join(clauses, " AND "), args
}

// BuildInClause 安全地構建 IN 子句
func BuildInClause(column string, values []interface{}) (string, []interface{}) {
	if len(values) == 0 {
		return "", nil
	}
	
	// 驗證列名
	safeCol, err := SafeColumnName(column)
	if err != nil {
		return "", nil
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	
	clause := fmt.Sprintf("%s IN (%s)", safeCol, strings.Join(placeholders, ", "))
	return clause, values
}

// EscapeLike 轉義 LIKE 查詢中的特殊字符
func EscapeLike(pattern string) string {
	// 轉義 SQL LIKE 中的特殊字符
	pattern = strings.ReplaceAll(pattern, "\\", "\\\\")
	pattern = strings.ReplaceAll(pattern, "%", "\\%")
	pattern = strings.ReplaceAll(pattern, "_", "\\_")
	pattern = strings.ReplaceAll(pattern, "[", "\\[")
	return pattern
}

// ValidateLimit 驗證並限制查詢結果數量
func ValidateLimit(limit int, maxLimit int) int {
	if limit <= 0 {
		return 20 // 預設值
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

// ValidateOffset 驗證並限制偏移量
func ValidateOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

// IsValidTableName 檢查表名是否有效
func IsValidTableName(name string) bool {
	if name == "" {
		return false
	}
	// 只允許字母、數字和下劃線
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, name)
	return matched
}

// IsValidColumnName 檢查列名是否有效
func IsValidColumnName(name string) bool {
	if name == "" {
		return false
	}
	// 只允許字母、數字和下劃線
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, name)
	return matched
}

// SafeOrderBy 創建安全的 ORDER BY 子句（簡化版本）
func SafeOrderBy(field string) string {
	// 只允許字母、數字和下劃線
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, field); matched {
		return field + " DESC"
	}
	return "created_at DESC" // 默認排序
}
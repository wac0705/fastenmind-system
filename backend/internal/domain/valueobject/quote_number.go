package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// QuoteNumber 報價單號值對象
type QuoteNumber struct {
	value string
}

// NewQuoteNumber 創建報價單號
func NewQuoteNumber(value string) (QuoteNumber, error) {
	if value == "" {
		return QuoteNumber{}, errors.New("quote number cannot be empty")
	}
	
	// 驗證格式: Q-YYYYMMDD-XXXX
	pattern := `^Q-\d{8}-[A-Z0-9]{4}$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return QuoteNumber{}, err
	}
	
	if !matched {
		return QuoteNumber{}, errors.New("invalid quote number format")
	}
	
	return QuoteNumber{value: value}, nil
}

// GenerateQuoteNumber 生成報價單號
func GenerateQuoteNumber(companyID uuid.UUID, timestamp time.Time) QuoteNumber {
	// 格式: Q-YYYYMMDD-XXXX (XXXX 是基於公司ID和時間的唯一碼)
	dateStr := timestamp.Format("20060102")
	uniqueCode := generateUniqueCode(companyID, timestamp)
	
	return QuoteNumber{
		value: fmt.Sprintf("Q-%s-%s", dateStr, uniqueCode),
	}
}

// String 返回報價單號字符串
func (q QuoteNumber) String() string {
	return q.value
}

// Equals 比較兩個報價單號是否相等
func (q QuoteNumber) Equals(other QuoteNumber) bool {
	return q.value == other.value
}

// GetDate 獲取報價單號中的日期
func (q QuoteNumber) GetDate() (time.Time, error) {
	if len(q.value) < 11 {
		return time.Time{}, errors.New("invalid quote number format")
	}
	
	dateStr := q.value[2:10]
	return time.Parse("20060102", dateStr)
}

func generateUniqueCode(companyID uuid.UUID, timestamp time.Time) string {
	// 使用公司ID和時間戳生成4位唯一碼
	hash := companyID[0] + companyID[1] + byte(timestamp.Unix()&0xFF) + byte((timestamp.Unix()>>8)&0xFF)
	return fmt.Sprintf("%04X", hash)
}
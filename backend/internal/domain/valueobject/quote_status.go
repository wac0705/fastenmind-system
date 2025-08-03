package valueobject

import "errors"

// QuoteStatus 報價狀態
type QuoteStatus string

const (
	QuoteStatusDraft     QuoteStatus = "DRAFT"      // 草稿
	QuoteStatusPending   QuoteStatus = "PENDING"    // 待審核
	QuoteStatusApproved  QuoteStatus = "APPROVED"   // 已批准
	QuoteStatusRejected  QuoteStatus = "REJECTED"   // 已拒絕
	QuoteStatusExpired   QuoteStatus = "EXPIRED"    // 已過期
	QuoteStatusCancelled QuoteStatus = "CANCELLED"  // 已取消
)

// Validate 驗證報價狀態
func (s QuoteStatus) Validate() error {
	switch s {
	case QuoteStatusDraft, QuoteStatusPending, QuoteStatusApproved, 
	     QuoteStatusRejected, QuoteStatusExpired, QuoteStatusCancelled:
		return nil
	default:
		return errors.New("invalid quote status")
	}
}

// CanTransitionTo 檢查是否可以轉換到目標狀態
func (s QuoteStatus) CanTransitionTo(target QuoteStatus) bool {
	transitions := map[QuoteStatus][]QuoteStatus{
		QuoteStatusDraft: {
			QuoteStatusPending,
			QuoteStatusCancelled,
		},
		QuoteStatusPending: {
			QuoteStatusApproved,
			QuoteStatusRejected,
			QuoteStatusExpired,
			QuoteStatusCancelled,
		},
		QuoteStatusApproved: {
			QuoteStatusExpired,
			QuoteStatusCancelled,
		},
		QuoteStatusRejected: {
			// 拒絕後不能轉換到其他狀態
		},
		QuoteStatusExpired: {
			// 過期後不能轉換到其他狀態
		},
		QuoteStatusCancelled: {
			// 取消後不能轉換到其他狀態
		},
	}
	
	allowedTransitions, exists := transitions[s]
	if !exists {
		return false
	}
	
	for _, allowed := range allowedTransitions {
		if allowed == target {
			return true
		}
	}
	
	return false
}

// IsFinal 檢查是否為最終狀態
func (s QuoteStatus) IsFinal() bool {
	return s == QuoteStatusRejected || s == QuoteStatusExpired || s == QuoteStatusCancelled
}

// IsActive 檢查是否為活躍狀態
func (s QuoteStatus) IsActive() bool {
	return s == QuoteStatusPending || s == QuoteStatusApproved
}

// String 返回狀態字符串
func (s QuoteStatus) String() string {
	return string(s)
}
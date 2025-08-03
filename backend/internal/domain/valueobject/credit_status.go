package valueobject

import "errors"

// CreditStatus 信用狀態
type CreditStatus string

const (
	CreditStatusGood    CreditStatus = "GOOD"     // 良好
	CreditStatusWarning CreditStatus = "WARNING"  // 警告
	CreditStatusBlocked CreditStatus = "BLOCKED"  // 封鎖
	CreditStatusNew     CreditStatus = "NEW"      // 新客戶
)

// Validate 驗證信用狀態
func (s CreditStatus) Validate() error {
	switch s {
	case CreditStatusGood, CreditStatusWarning, CreditStatusBlocked, CreditStatusNew:
		return nil
	default:
		return errors.New("invalid credit status")
	}
}

// CanTransitionTo 檢查是否可以轉換到目標狀態
func (s CreditStatus) CanTransitionTo(target CreditStatus) bool {
	// 定義允許的狀態轉換
	transitions := map[CreditStatus][]CreditStatus{
		CreditStatusNew: {
			CreditStatusGood,
			CreditStatusWarning,
			CreditStatusBlocked,
		},
		CreditStatusGood: {
			CreditStatusWarning,
			CreditStatusBlocked,
		},
		CreditStatusWarning: {
			CreditStatusGood,
			CreditStatusBlocked,
		},
		CreditStatusBlocked: {
			CreditStatusWarning,
			CreditStatusGood,
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

// GetCreditLimit 獲取建議的信用額度係數
func (s CreditStatus) GetCreditLimit() float64 {
	switch s {
	case CreditStatusGood:
		return 1.0 // 100% 信用額度
	case CreditStatusWarning:
		return 0.5 // 50% 信用額度
	case CreditStatusBlocked:
		return 0.0 // 無信用額度
	case CreditStatusNew:
		return 0.3 // 30% 信用額度
	default:
		return 0.0
	}
}

// RequiresApproval 是否需要審核
func (s CreditStatus) RequiresApproval() bool {
	return s == CreditStatusWarning || s == CreditStatusBlocked || s == CreditStatusNew
}

// String 返回狀態字符串
func (s CreditStatus) String() string {
	return string(s)
}
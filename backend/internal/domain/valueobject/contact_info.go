package valueobject

import (
	"errors"
	"regexp"
)

// ContactInfo 聯絡信息值對象
type ContactInfo struct {
	ContactPerson string
	Email         Email
	Phone         Phone
	Mobile        Phone
	Fax           string
	Department    string
	Position      string
}

// Validate 驗證聯絡信息
func (c ContactInfo) Validate() error {
	if c.ContactPerson == "" {
		return errors.New("contact person is required")
	}
	
	if err := c.Email.Validate(); err != nil {
		return err
	}
	
	if err := c.Phone.Validate(); err != nil {
		return err
	}
	
	return nil
}

// Email 電子郵件值對象
type Email struct {
	value string
}

// NewEmail 創建電子郵件
func NewEmail(value string) (Email, error) {
	email := Email{value: value}
	if err := email.Validate(); err != nil {
		return Email{}, err
	}
	return email, nil
}

// String 返回電子郵件字符串
func (e Email) String() string {
	return e.value
}

// Validate 驗證電子郵件
func (e Email) Validate() error {
	if e.value == "" {
		return errors.New("email is required")
	}
	
	// 簡單的電子郵件正則表達式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, e.value)
	if err != nil {
		return err
	}
	
	if !matched {
		return errors.New("invalid email format")
	}
	
	return nil
}

// Phone 電話號碼值對象
type Phone struct {
	CountryCode string
	Number      string
	Extension   string
}

// NewPhone 創建電話號碼
func NewPhone(countryCode, number string) (Phone, error) {
	phone := Phone{
		CountryCode: countryCode,
		Number:      number,
	}
	
	if err := phone.Validate(); err != nil {
		return Phone{}, err
	}
	
	return phone, nil
}

// String 返回電話號碼字符串
func (p Phone) String() string {
	result := ""
	if p.CountryCode != "" {
		result = "+" + p.CountryCode + " "
	}
	result += p.Number
	if p.Extension != "" {
		result += " ext. " + p.Extension
	}
	return result
}

// Validate 驗證電話號碼
func (p Phone) Validate() error {
	if p.Number == "" {
		return errors.New("phone number is required")
	}
	
	// 移除所有非數字字符進行驗證
	digitsOnly := regexp.MustCompile(`\D`).ReplaceAllString(p.Number, "")
	if len(digitsOnly) < 7 || len(digitsOnly) > 15 {
		return errors.New("phone number must be between 7 and 15 digits")
	}
	
	if p.CountryCode != "" {
		// 驗證國家代碼（1-3位數字）
		matched, _ := regexp.MatchString(`^\d{1,3}$`, p.CountryCode)
		if !matched {
			return errors.New("invalid country code")
		}
	}
	
	return nil
}
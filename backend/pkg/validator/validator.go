package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Validator 封裝 validator 實例
type Validator struct {
	validator *validator.Validate
}

// New 創建新的驗證器
func New() *Validator {
	v := validator.New()
	
	// 註冊自定義驗證規則
	registerCustomValidations(v)
	
	return &Validator{
		validator: v,
	}
}

// Validate 驗證結構體
func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// ValidateVar 驗證單個變量
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// registerCustomValidations 註冊自定義驗證規則
func registerCustomValidations(v *validator.Validate) {
	// 電話號碼驗證
	v.RegisterValidation("phone", validatePhone)
	
	// 貨幣代碼驗證
	v.RegisterValidation("currency", validateCurrency)
	
	// 國家代碼驗證
	v.RegisterValidation("country", validateCountry)
	
	// 公司代碼驗證
	v.RegisterValidation("company_code", validateCompanyCode)
	
	// 產品代碼驗證
	v.RegisterValidation("product_code", validateProductCode)
	
	// 安全的檔案名驗證
	v.RegisterValidation("safe_filename", validateSafeFilename)
	
	// UUID 驗證
	v.RegisterValidation("uuid", validateUUID)
	
	// 未來日期驗證
	v.RegisterValidation("future_date", validateFutureDate)
	
	// 價格驗證
	v.RegisterValidation("price", validatePrice)
	
	// 百分比驗證
	v.RegisterValidation("percentage", validatePercentage)
	
	// 國家代碼驗證 (2字母)
	v.RegisterValidation("country_code", validateCountryCode)
	
	// HS Code 驗證
	v.RegisterValidation("hs_code", validateHSCode)
	
	// 正數驗證
	v.RegisterValidation("positive", validatePositive)
}

// validatePhone 驗證電話號碼
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 支援國際電話格式 - more flexible regex
	// Allows formats like: +1-555-123-4567, (555) 123-4567, +886-2-2345-6789, etc.
	// Minimum 7 digits for a valid phone number
	phoneRegex := regexp.MustCompile(`^[\+]?[(]?[0-9]{1,4}[)]?([-\s\.]?)[(]?[0-9]{1,4}[)]?([-\s\.]?)[0-9]{1,4}([-\s\.]?)[0-9]{1,4}$`)
	return phoneRegex.MatchString(phone) && len(regexp.MustCompile(`[0-9]`).FindAllString(phone, -1)) >= 7
}

// validateCurrency 驗證貨幣代碼
func validateCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	validCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CNY", "TWD", "HKD", "SGD", "AUD", "CAD"}
	for _, valid := range validCurrencies {
		if currency == valid {
			return true
		}
	}
	return false
}

// validateCountry 驗證國家代碼
func validateCountry(fl validator.FieldLevel) bool {
	country := fl.Field().String()
	// ISO 3166-1 alpha-2 country codes
	countryRegex := regexp.MustCompile(`^[A-Z]{2}$`)
	return countryRegex.MatchString(country)
}

// validateCompanyCode 驗證公司代碼
func validateCompanyCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	// 公司代碼：3-10個字母數字，可包含連字符
	codeRegex := regexp.MustCompile(`^[A-Z0-9\-]{3,10}$`)
	return codeRegex.MatchString(code)
}

// validateProductCode 驗證產品代碼
func validateProductCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	// 產品代碼：字母開頭，可包含字母、數字、連字符
	codeRegex := regexp.MustCompile(`^[A-Z][A-Z0-9\-]{2,20}$`)
	return codeRegex.MatchString(code)
}

// validateSafeFilename 驗證安全的檔案名
func validateSafeFilename(fl validator.FieldLevel) bool {
	filename := fl.Field().String()
	// 只允許字母、數字、連字符、底線和點
	filenameRegex := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	return filenameRegex.MatchString(filename) && !strings.Contains(filename, "..")
}

// validateUUID 驗證 UUID
func validateUUID(fl validator.FieldLevel) bool {
	uuid := fl.Field().String()
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

// validateFutureDate 驗證未來日期
func validateFutureDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.After(time.Now())
}

// validatePrice 驗證價格（必須為正數，最多兩位小數）
func validatePrice(fl validator.FieldLevel) bool {
	price := fl.Field().Float()
	return price >= 0 && price == float64(int(price*100))/100
}

// validatePercentage 驗證百分比（0-100）
func validatePercentage(fl validator.FieldLevel) bool {
	percentage := fl.Field().Float()
	return percentage >= 0 && percentage <= 100
}

// validateCountryCode 驗證國家代碼 (2字母)
func validateCountryCode(fl validator.FieldLevel) bool {
	country := fl.Field().String()
	// ISO 3166-1 alpha-2 country codes
	countryRegex := regexp.MustCompile(`^[A-Z]{2}$`)
	result := countryRegex.MatchString(country)
	// Debug
	// fmt.Printf("validateCountryCode: '%s' -> %v\n", country, result)
	return result
}

// validateHSCode 驗證 HS Code
func validateHSCode(fl validator.FieldLevel) bool {
	hsCode := fl.Field().String()
	// HS codes should be 4, 6, 8, or 10 digits
	hsCodeRegex := regexp.MustCompile(`^\d{4}$|^\d{6}$|^\d{8}$|^\d{10}$`)
	return hsCodeRegex.MatchString(hsCode)
}

// validatePositive 驗證正數
func validatePositive(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Field().Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Field().Uint() > 0
	case reflect.Float32, reflect.Float64:
		return fl.Field().Float() > 0
	default:
		return false
	}
}

// formatValidationError 格式化驗證錯誤
func formatValidationError(err error) error {
	var errors []string
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, formatFieldError(e))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}
	
	return err
}

// formatFieldError 格式化欄位錯誤
func formatFieldError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "currency":
		return fmt.Sprintf("%s must be a valid currency code", field)
	case "country":
		return fmt.Sprintf("%s must be a valid country code", field)
	case "country_code":
		return fmt.Sprintf("%s must be a valid 2-letter country code", field)
	case "hs_code":
		return fmt.Sprintf("%s must be a valid HS code (4, 6, 8, or 10 digits)", field)
	case "positive":
		return fmt.Sprintf("%s must be a positive number", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "future_date":
		return fmt.Sprintf("%s must be a future date", field)
	case "price":
		return fmt.Sprintf("%s must be a valid price", field)
	case "percentage":
		return fmt.Sprintf("%s must be between 0 and 100", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// ValidationError 自定義驗證錯誤
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors 驗證錯誤集合
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, e := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(messages, "; ")
}
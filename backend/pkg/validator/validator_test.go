package validator

import (
	"testing"
)

func TestValidatePhone(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{
			name:    "Valid US phone",
			phone:   "+1-555-123-4567",
			wantErr: false,
		},
		{
			name:    "Valid international phone",
			phone:   "+886-2-2345-6789",
			wantErr: false,
		},
		{
			name:    "Valid phone without country code",
			phone:   "555-123-4567",
			wantErr: false,
		},
		{
			name:    "Valid phone with parentheses",
			phone:   "(555) 123-4567",
			wantErr: false,
		},
		{
			name:    "Invalid phone - too short",
			phone:   "123",
			wantErr: true,
		},
		{
			name:    "Invalid phone - letters",
			phone:   "555-CALL-NOW",
			wantErr: true,
		},
		{
			name:    "Empty phone",
			phone:   "",
			wantErr: true,
		},
	}

	type TestStruct struct {
		Phone string `validate:"required,phone"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Phone: tt.phone}
			err := v.Validate(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCurrency(t *testing.T) {
	v := New()

	tests := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{
			name:     "Valid USD",
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "Valid EUR",
			currency: "EUR",
			wantErr:  false,
		},
		{
			name:     "Valid TWD",
			currency: "TWD",
			wantErr:  false,
		},
		{
			name:     "Invalid currency - lowercase",
			currency: "usd",
			wantErr:  true,
		},
		{
			name:     "Invalid currency - too long",
			currency: "USDD",
			wantErr:  true,
		},
		{
			name:     "Invalid currency - too short",
			currency: "US",
			wantErr:  true,
		},
		{
			name:     "Invalid currency - numbers",
			currency: "123",
			wantErr:  true,
		},
	}

	type TestStruct struct {
		Currency string `validate:"required,currency"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Currency: tt.currency}
			err := v.Validate(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCountryCode(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		country string
		wantErr bool
	}{
		{
			name:    "Valid US",
			country: "US",
			wantErr: false,
		},
		{
			name:    "Valid TW",
			country: "TW",
			wantErr: false,
		},
		{
			name:    "Valid GB",
			country: "GB",
			wantErr: false,
		},
		{
			name:    "Invalid - lowercase",
			country: "us",
			wantErr: true,
		},
		{
			name:    "Invalid - too long",
			country: "USA",
			wantErr: true,
		},
		{
			name:    "Invalid - single letter",
			country: "U",
			wantErr: true,
		},
		{
			name:    "Invalid - numbers",
			country: "12",
			wantErr: true,
		},
	}

	type TestStruct struct {
		Country string `validate:"required,country_code"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Country: tt.country}
			err := v.Validate(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCountryCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHSCode(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		hsCode  string
		wantErr bool
	}{
		{
			name:    "Valid 4-digit HS code",
			hsCode:  "8501",
			wantErr: false,
		},
		{
			name:    "Valid 6-digit HS code",
			hsCode:  "850110",
			wantErr: false,
		},
		{
			name:    "Valid 8-digit HS code",
			hsCode:  "85011010",
			wantErr: false,
		},
		{
			name:    "Valid 10-digit HS code",
			hsCode:  "8501101000",
			wantErr: false,
		},
		{
			name:    "Invalid - letters",
			hsCode:  "850A",
			wantErr: true,
		},
		{
			name:    "Invalid - too short",
			hsCode:  "850",
			wantErr: true,
		},
		{
			name:    "Invalid - too long",
			hsCode:  "85011010001",
			wantErr: true,
		},
		{
			name:    "Invalid - odd length",
			hsCode:  "85011",
			wantErr: true,
		},
	}

	type TestStruct struct {
		HSCode string `validate:"required,hs_code"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{HSCode: tt.hsCode}
			err := v.Validate(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHSCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePositive(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{
			name:    "Positive integer",
			value:   10,
			wantErr: false,
		},
		{
			name:    "Positive decimal",
			value:   10.5,
			wantErr: false,
		},
		{
			name:    "Zero",
			value:   0,
			wantErr: true,
		},
		{
			name:    "Negative",
			value:   -5,
			wantErr: true,
		},
	}

	type TestStruct struct {
		Value float64 `validate:"positive"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Value: tt.value}
			err := v.Validate(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePositive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePercentage(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{
			name:    "Valid 0%",
			value:   0,
			wantErr: false,
		},
		{
			name:    "Valid 50%",
			value:   50,
			wantErr: false,
		},
		{
			name:    "Valid 100%",
			value:   100,
			wantErr: false,
		},
		{
			name:    "Valid decimal",
			value:   25.5,
			wantErr: false,
		},
		{
			name:    "Invalid - negative",
			value:   -5,
			wantErr: true,
		},
		{
			name:    "Invalid - over 100",
			value:   101,
			wantErr: true,
		},
	}

	type TestStruct struct {
		Value float64 `validate:"percentage"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Value: tt.value}
			err := v.Validate(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePercentage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComplexValidation(t *testing.T) {
	v := New()

	type Address struct {
		Street  string `validate:"required,min=5,max=100"`
		City    string `validate:"required,min=2,max=50"`
		Country string `validate:"required,country_code"`
		Phone   string `validate:"omitempty,phone"`
	}

	type Customer struct {
		Name     string   `validate:"required,min=2,max=100"`
		Email    string   `validate:"required,email"`
		Phone    string   `validate:"required,phone"`
		Currency string   `validate:"required,currency"`
		Address  *Address `validate:"required"`
	}

	tests := []struct {
		name     string
		customer Customer
		wantErr  bool
	}{
		{
			name: "Valid customer",
			customer: Customer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "+1-555-123-4567",
				Currency: "USD",
				Address: &Address{
					Street:  "123 Main Street",
					City:    "New York",
					Country: "US",
					Phone:   "+1-555-987-6543",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid - missing address",
			customer: Customer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "+1-555-123-4567",
				Currency: "USD",
				Address:  nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid - bad email",
			customer: Customer{
				Name:     "John Doe",
				Email:    "not-an-email",
				Phone:    "+1-555-123-4567",
				Currency: "USD",
				Address: &Address{
					Street:  "123 Main Street",
					City:    "New York",
					Country: "US",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid - short street",
			customer: Customer{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "+1-555-123-4567",
				Currency: "USD",
				Address: &Address{
					Street:  "123",
					City:    "New York",
					Country: "US",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.customer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComplexValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
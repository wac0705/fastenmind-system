package valueobject

import (
	"errors"
	"strings"
)

// Address 地址值對象
type Address struct {
	Street1    string
	Street2    string
	City       string
	State      string
	PostalCode string
	Country    Country
	Type       AddressType
}

// NewAddress 創建地址
func NewAddress(street1, city, state, postalCode string, country Country) (Address, error) {
	address := Address{
		Street1:    street1,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
		Type:       AddressTypeBusiness,
	}
	
	if err := address.Validate(); err != nil {
		return Address{}, err
	}
	
	return address, nil
}

// Validate 驗證地址
func (a Address) Validate() error {
	if strings.TrimSpace(a.Street1) == "" {
		return errors.New("street address is required")
	}
	
	if strings.TrimSpace(a.City) == "" {
		return errors.New("city is required")
	}
	
	if strings.TrimSpace(a.PostalCode) == "" {
		return errors.New("postal code is required")
	}
	
	if err := a.Country.Validate(); err != nil {
		return err
	}
	
	if err := a.Type.Validate(); err != nil {
		return err
	}
	
	return nil
}

// String 返回格式化的地址字符串
func (a Address) String() string {
	parts := []string{}
	
	if a.Street1 != "" {
		parts = append(parts, a.Street1)
	}
	if a.Street2 != "" {
		parts = append(parts, a.Street2)
	}
	
	cityStateZip := ""
	if a.City != "" {
		cityStateZip = a.City
	}
	if a.State != "" {
		if cityStateZip != "" {
			cityStateZip += ", "
		}
		cityStateZip += a.State
	}
	if a.PostalCode != "" {
		if cityStateZip != "" {
			cityStateZip += " "
		}
		cityStateZip += a.PostalCode
	}
	
	if cityStateZip != "" {
		parts = append(parts, cityStateZip)
	}
	
	if a.Country != "" {
		parts = append(parts, a.Country.String())
	}
	
	return strings.Join(parts, "\n")
}

// AddressType 地址類型
type AddressType string

const (
	AddressTypeBusiness   AddressType = "BUSINESS"
	AddressTypeShipping   AddressType = "SHIPPING"
	AddressTypeBilling    AddressType = "BILLING"
	AddressTypeRegistered AddressType = "REGISTERED"
)

// Validate 驗證地址類型
func (t AddressType) Validate() error {
	switch t {
	case AddressTypeBusiness, AddressTypeShipping, AddressTypeBilling, AddressTypeRegistered:
		return nil
	default:
		return errors.New("invalid address type")
	}
}

// Country 國家
type Country string

const (
	CountryUS Country = "US"
	CountryCA Country = "CA"
	CountryMX Country = "MX"
	CountryGB Country = "GB"
	CountryDE Country = "DE"
	CountryFR Country = "FR"
	CountryIT Country = "IT"
	CountryES Country = "ES"
	CountryJP Country = "JP"
	CountryKR Country = "KR"
	CountryCN Country = "CN"
	CountryTW Country = "TW"
	CountryHK Country = "HK"
	CountrySG Country = "SG"
	CountryIN Country = "IN"
	CountryAU Country = "AU"
	CountryNZ Country = "NZ"
)

// Validate 驗證國家代碼
func (c Country) Validate() error {
	// ISO 3166-1 alpha-2 country codes
	validCountries := map[Country]bool{
		CountryUS: true, CountryCA: true, CountryMX: true,
		CountryGB: true, CountryDE: true, CountryFR: true,
		CountryIT: true, CountryES: true, CountryJP: true,
		CountryKR: true, CountryCN: true, CountryTW: true,
		CountryHK: true, CountrySG: true, CountryIN: true,
		CountryAU: true, CountryNZ: true,
	}
	
	if _, ok := validCountries[c]; !ok {
		return errors.New("invalid country code")
	}
	
	return nil
}

// String 返回國家代碼字符串
func (c Country) String() string {
	return string(c)
}

// GetCountryName 獲取國家名稱
func (c Country) GetCountryName() string {
	names := map[Country]string{
		CountryUS: "United States",
		CountryCA: "Canada",
		CountryMX: "Mexico",
		CountryGB: "United Kingdom",
		CountryDE: "Germany",
		CountryFR: "France",
		CountryIT: "Italy",
		CountryES: "Spain",
		CountryJP: "Japan",
		CountryKR: "South Korea",
		CountryCN: "China",
		CountryTW: "Taiwan",
		CountryHK: "Hong Kong",
		CountrySG: "Singapore",
		CountryIN: "India",
		CountryAU: "Australia",
		CountryNZ: "New Zealand",
	}
	
	if name, ok := names[c]; ok {
		return name
	}
	
	return string(c)
}
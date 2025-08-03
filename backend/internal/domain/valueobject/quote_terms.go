package valueobject

import (
	"errors"
	"time"
)

// QuoteTerms 報價條款值對象
type QuoteTerms struct {
	PaymentTerms       PaymentTerms
	DeliveryTerms      DeliveryTerms
	WarrantyTerms      WarrantyTerms
	Currency           Currency
	DiscountPercentage float64
	Notes              string
}

// Validate 驗證報價條款
func (t QuoteTerms) Validate() error {
	if err := t.PaymentTerms.Validate(); err != nil {
		return err
	}
	
	if err := t.DeliveryTerms.Validate(); err != nil {
		return err
	}
	
	if err := t.WarrantyTerms.Validate(); err != nil {
		return err
	}
	
	if err := t.Currency.Validate(); err != nil {
		return err
	}
	
	if t.DiscountPercentage < 0 || t.DiscountPercentage > 100 {
		return errors.New("discount percentage must be between 0 and 100")
	}
	
	return nil
}

// PaymentTerms 付款條款
type PaymentTerms struct {
	Type           PaymentTermType
	NetDays        int
	DepositPercent float64
	Description    string
}

// PaymentTermType 付款條款類型
type PaymentTermType string

const (
	PaymentTermTypeNet30      PaymentTermType = "NET30"
	PaymentTermTypeNet60      PaymentTermType = "NET60"
	PaymentTermTypeNet90      PaymentTermType = "NET90"
	PaymentTermTypeCOD        PaymentTermType = "COD"
	PaymentTermTypePrepayment PaymentTermType = "PREPAYMENT"
	PaymentTermTypeCustom     PaymentTermType = "CUSTOM"
)

// Validate 驗證付款條款
func (p PaymentTerms) Validate() error {
	switch p.Type {
	case PaymentTermTypeNet30, PaymentTermTypeNet60, PaymentTermTypeNet90,
		 PaymentTermTypeCOD, PaymentTermTypePrepayment, PaymentTermTypeCustom:
		// Valid
	default:
		return errors.New("invalid payment term type")
	}
	
	if p.NetDays < 0 {
		return errors.New("net days cannot be negative")
	}
	
	if p.DepositPercent < 0 || p.DepositPercent > 100 {
		return errors.New("deposit percentage must be between 0 and 100")
	}
	
	return nil
}

// DeliveryTerms 交貨條款
type DeliveryTerms struct {
	Incoterm      Incoterm
	LeadTimeDays  int
	Location      string
	Description   string
}

// Incoterm 國際貿易術語
type Incoterm string

const (
	IncotermEXW Incoterm = "EXW" // Ex Works
	IncotermFOB Incoterm = "FOB" // Free On Board
	IncotermCIF Incoterm = "CIF" // Cost, Insurance and Freight
	IncotermDDP Incoterm = "DDP" // Delivered Duty Paid
	IncotermFCA Incoterm = "FCA" // Free Carrier
	IncotermCPT Incoterm = "CPT" // Carriage Paid To
	IncotermCIP Incoterm = "CIP" // Carriage and Insurance Paid To
	IncotermDAP Incoterm = "DAP" // Delivered At Place
	IncotermDPU Incoterm = "DPU" // Delivered at Place Unloaded
)

// Validate 驗證交貨條款
func (d DeliveryTerms) Validate() error {
	if err := d.Incoterm.Validate(); err != nil {
		return err
	}
	
	if d.LeadTimeDays < 0 {
		return errors.New("lead time cannot be negative")
	}
	
	if d.Location == "" {
		return errors.New("delivery location is required")
	}
	
	return nil
}

// Validate 驗證國際貿易術語
func (i Incoterm) Validate() error {
	switch i {
	case IncotermEXW, IncotermFOB, IncotermCIF, IncotermDDP,
		 IncotermFCA, IncotermCPT, IncotermCIP, IncotermDAP, IncotermDPU:
		return nil
	default:
		return errors.New("invalid incoterm")
	}
}

// WarrantyTerms 保固條款
type WarrantyTerms struct {
	Duration    time.Duration
	Type        WarrantyType
	Coverage    string
	Exclusions  []string
	Description string
}

// WarrantyType 保固類型
type WarrantyType string

const (
	WarrantyTypeStandard     WarrantyType = "STANDARD"
	WarrantyTypeExtended     WarrantyType = "EXTENDED"
	WarrantyTypeLimited      WarrantyType = "LIMITED"
	WarrantyTypeComprehensive WarrantyType = "COMPREHENSIVE"
)

// Validate 驗證保固條款
func (w WarrantyTerms) Validate() error {
	if w.Duration < 0 {
		return errors.New("warranty duration cannot be negative")
	}
	
	switch w.Type {
	case WarrantyTypeStandard, WarrantyTypeExtended, 
	     WarrantyTypeLimited, WarrantyTypeComprehensive:
		// Valid
	default:
		return errors.New("invalid warranty type")
	}
	
	if w.Coverage == "" {
		return errors.New("warranty coverage is required")
	}
	
	return nil
}

// Currency 貨幣
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyJPY Currency = "JPY"
	CurrencyCNY Currency = "CNY"
	CurrencyTWD Currency = "TWD"
	CurrencyHKD Currency = "HKD"
	CurrencySGD Currency = "SGD"
)

// Validate 驗證貨幣
func (c Currency) Validate() error {
	switch c {
	case CurrencyUSD, CurrencyEUR, CurrencyGBP, CurrencyJPY,
		 CurrencyCNY, CurrencyTWD, CurrencyHKD, CurrencySGD:
		return nil
	default:
		return errors.New("unsupported currency")
	}
}

// String 返回貨幣字符串
func (c Currency) String() string {
	return string(c)
}
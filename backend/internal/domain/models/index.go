package models

// This file exports all domain models from a single location
// to prevent circular dependencies and ensure consistency

// Re-export all models for easy import
type (
	// Base
	BaseModel_ = BaseModel
	
	// Company & Account
	Company_                  = Company
	CompanyType_              = CompanyType
	Account_                  = Account
	AccountRole_              = AccountRole
	
	// Customer
	Customer_                 = Customer
	CustomerTransactionTerms_ = CustomerTransactionTerms
	
	// Product & Process
	Product_                  = Product
	ProductStatus_            = ProductStatus
	ProductProcess_           = ProductProcess
	Process_                  = Process
	Equipment_                = Equipment
	
	// Inquiry
	Inquiry_                  = Inquiry
	InquiryStatus_            = InquiryStatus
	InquiryItem_              = InquiryItem
	
	// Quote
	Quote_                    = Quote
	QuoteStatus_              = QuoteStatus
	QuoteItem_                = QuoteItem
	QuoteRevision_            = QuoteRevision
	
	// Order
	Order_                    = Order
	OrderStatus_              = OrderStatus
	OrderItem_                = OrderItem
	
	// Shipment
	Shipment_                 = Shipment
	ShipmentStatus_           = ShipmentStatus
	ShipmentItem_             = ShipmentItem
	
	// Invoice & Payment
	Invoice_                  = Invoice
	InvoiceStatus_            = InvoiceStatus
	InvoiceItem_              = InvoiceItem
	Payment_                  = Payment
)

// Constants re-export
const (
	// Company Types
	CompanyTypeHeadquarters_ = CompanyTypeHeadquarters
	CompanyTypeSubsidiary_   = CompanyTypeSubsidiary
	CompanyTypeFactory_      = CompanyTypeFactory
	
	// Account Roles
	RoleAdmin_               = RoleAdmin
	RoleManager_             = RoleManager
	RoleEngineer_            = RoleEngineer
	RoleSales_               = RoleSales
	RoleViewer_              = RoleViewer
	
	// Product Status
	ProductStatusActive_     = ProductStatusActive
	ProductStatusDiscontinued_ = ProductStatusDiscontinued
	ProductStatusDraft_      = ProductStatusDraft
	
	// Inquiry Status
	InquiryStatusPending_    = InquiryStatusPending
	InquiryStatusAssigned_   = InquiryStatusAssigned
	InquiryStatusQuoted_     = InquiryStatusQuoted
	InquiryStatusRejected_   = InquiryStatusRejected
	InquiryStatusCancelled_  = InquiryStatusCancelled
	
	// Quote Status
	QuoteStatusDraft_        = QuoteStatusDraft
	QuoteStatusSubmitted_    = QuoteStatusSubmitted
	QuoteStatusApproved_     = QuoteStatusApproved
	QuoteStatusRejected_     = QuoteStatusRejected
	QuoteStatusExpired_      = QuoteStatusExpired
	QuoteStatusOrdered_      = QuoteStatusOrdered
	
	// Order Status
	OrderStatusDraft_        = OrderStatusDraft
	OrderStatusConfirmed_    = OrderStatusConfirmed
	OrderStatusProcessing_   = OrderStatusProcessing
	OrderStatusShipped_      = OrderStatusShipped
	OrderStatusDelivered_    = OrderStatusDelivered
	OrderStatusCancelled_    = OrderStatusCancelled
	
	// Shipment Status
	ShipmentStatusPending_   = ShipmentStatusPending
	ShipmentStatusInTransit_ = ShipmentStatusInTransit
	ShipmentStatusDelivered_ = ShipmentStatusDelivered
	ShipmentStatusCancelled_ = ShipmentStatusCancelled
	
	// Invoice Status
	InvoiceStatusDraft_      = InvoiceStatusDraft
	InvoiceStatusIssued_     = InvoiceStatusIssued
	InvoiceStatusPaid_       = InvoiceStatusPaid
	InvoiceStatusOverdue_    = InvoiceStatusOverdue
	InvoiceStatusCancelled_  = InvoiceStatusCancelled
)
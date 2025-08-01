package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator interface {
	GenerateQuotePDF(quote *models.Quote) ([]byte, error)
}

type pdfGenerator struct{}

func NewPDFGenerator() PDFGenerator {
	return &pdfGenerator{}
}

func (g *pdfGenerator) GenerateQuotePDF(quote *models.Quote) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Set font
	pdf.SetFont("Arial", "", 12)
	
	// Company Header
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "FastenMind Corporation")
	pdf.Ln(8)
	
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, "123 Industrial Avenue, Tech Park")
	pdf.Ln(5)
	pdf.Cell(0, 5, "City, State 12345")
	pdf.Ln(5)
	pdf.Cell(0, 5, "Tel: +1-234-567-8900 | Email: sales@fastenmind.com")
	pdf.Ln(15)
	
	// Quote Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("QUOTATION - %s", quote.QuoteNo))
	pdf.Ln(10)
	
	// Quote Info Box
	pdf.SetFillColor(240, 240, 240)
	pdf.Rect(10, pdf.GetY(), 190, 30, "F")
	
	pdf.SetFont("Arial", "", 10)
	y := pdf.GetY() + 5
	
	// Left column
	pdf.SetXY(15, y)
	pdf.Cell(40, 5, "Date:")
	pdf.Cell(50, 5, time.Now().Format("2006-01-02"))
	
	pdf.SetXY(15, y+6)
	pdf.Cell(40, 5, "Valid Until:")
	pdf.Cell(50, 5, quote.ValidUntil.Format("2006-01-02"))
	
	pdf.SetXY(15, y+12)
	pdf.Cell(40, 5, "Payment Terms:")
	pdf.Cell(50, 5, quote.PaymentTerms)
	
	pdf.SetXY(15, y+18)
	pdf.Cell(40, 5, "Delivery:")
	pdf.Cell(50, 5, fmt.Sprintf("%d days", quote.DeliveryDays))
	
	// Right column
	pdf.SetXY(110, y)
	pdf.Cell(40, 5, "Quote No:")
	pdf.Cell(50, 5, quote.QuoteNo)
	
	if quote.Inquiry != nil {
		pdf.SetXY(110, y+6)
		pdf.Cell(40, 5, "Inquiry No:")
		pdf.Cell(50, 5, quote.Inquiry.InquiryNo)
	}
	
	pdf.SetY(y + 35)
	
	// Customer Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "BILL TO:")
	pdf.Ln(6)
	
	pdf.SetFont("Arial", "", 10)
	if quote.Customer != nil {
		pdf.Cell(0, 5, quote.Customer.Name)
		pdf.Ln(5)
		if quote.Customer.Address != nil && *quote.Customer.Address != "" {
			pdf.Cell(0, 5, *quote.Customer.Address)
			pdf.Ln(5)
		}
		if quote.Customer.ContactPerson != nil && *quote.Customer.ContactPerson != "" {
			pdf.Cell(0, 5, fmt.Sprintf("Attn: %s", *quote.Customer.ContactPerson))
			pdf.Ln(5)
		}
		if quote.Customer.ContactEmail != nil && *quote.Customer.ContactEmail != "" {
			pdf.Cell(0, 5, fmt.Sprintf("Email: %s", *quote.Customer.ContactEmail))
			pdf.Ln(5)
		}
	}
	pdf.Ln(10)
	
	// Product Details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "PRODUCT DETAILS:")
	pdf.Ln(8)
	
	// Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(20, 8, "Item", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 8, "Part Number", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, "Unit Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, "Total Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Currency", "1", 0, "C", true, 0, "")
	pdf.Ln(8)
	
	// Table Content
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	
	if quote.Inquiry != nil {
		pdf.CellFormat(20, 8, "1", "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 8, quote.Inquiry.PartNo, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", quote.Inquiry.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.4f", quote.UnitPrice), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("%.2f", quote.TotalCost), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, quote.Currency, "1", 0, "C", false, 0, "")
		pdf.Ln(8)
	}
	
	pdf.Ln(5)
	
	// Cost Breakdown
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "COST BREAKDOWN:")
	pdf.Ln(8)
	
	pdf.SetFont("Arial", "", 10)
	x := pdf.GetX()
	y = pdf.GetY()
	
	// Cost items
	costItems := []struct {
		label string
		value float64
	}{
		{"Material Cost:", quote.MaterialCost},
		{"Process Cost:", quote.ProcessCost},
		{"Surface Treatment:", quote.SurfaceCost},
		{"Heat Treatment:", quote.HeatTreatCost},
		{"Packaging Cost:", quote.PackagingCost},
		{"Shipping Cost:", quote.ShippingCost},
		{"Tariff Cost:", quote.TariffCost},
	}
	
	for i, item := range costItems {
		pdf.SetXY(x+20, y+float64(i*6))
		pdf.Cell(60, 5, item.label)
		pdf.Cell(30, 5, fmt.Sprintf("$%.2f", item.value))
	}
	
	// Subtotal
	subtotal := quote.MaterialCost + quote.ProcessCost + quote.SurfaceCost + 
		quote.HeatTreatCost + quote.PackagingCost + quote.ShippingCost + quote.TariffCost
	
	pdf.SetXY(x+20, y+float64(len(costItems)*6+6))
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(60, 5, "Production Subtotal:")
	pdf.Cell(30, 5, fmt.Sprintf("$%.2f", subtotal))
	
	// Overhead and Profit
	overheadAmount := subtotal * (quote.OverheadRate / 100)
	subtotalWithOverhead := subtotal + overheadAmount
	profitAmount := subtotalWithOverhead * (quote.ProfitRate / 100)
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(x+20, y+float64(len(costItems)*6+12))
	pdf.Cell(60, 5, fmt.Sprintf("Overhead (%.1f%%):", quote.OverheadRate))
	pdf.Cell(30, 5, fmt.Sprintf("$%.2f", overheadAmount))
	
	pdf.SetXY(x+20, y+float64(len(costItems)*6+18))
	pdf.Cell(60, 5, fmt.Sprintf("Profit (%.1f%%):", quote.ProfitRate))
	pdf.Cell(30, 5, fmt.Sprintf("$%.2f", profitAmount))
	
	// Total
	pdf.SetXY(x+20, y+float64(len(costItems)*6+26))
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 6, "TOTAL COST:")
	pdf.Cell(30, 6, fmt.Sprintf("$%.2f", quote.TotalCost))
	
	pdf.Ln(40)
	
	// Terms and Notes
	if quote.Notes != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, "NOTES:")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 9)
		pdf.MultiCell(0, 5, quote.Notes, "", "", false)
		pdf.Ln(5)
	}
	
	// Footer
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 5, "This quotation is valid for the period stated above.")
	pdf.Ln(4)
	pdf.Cell(0, 5, "Prices are subject to change without notice after the validity period.")
	pdf.Ln(4)
	pdf.Cell(0, 5, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02 15:04:05")))
	
	// Generate PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}
package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/application/query"
	"github.com/fastenmind/fastener-api/internal/domain/aggregate"
	"github.com/fastenmind/fastener-api/internal/domain/repository"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// QuoteQueryHandler 報價查詢處理器
type QuoteQueryHandler struct {
	db        *gorm.DB
	quoteRepo repository.QuoteRepository
	logger    *zap.Logger
}

// NewQuoteQueryHandler 創建報價查詢處理器
func NewQuoteQueryHandler(
	db *gorm.DB,
	quoteRepo repository.QuoteRepository,
	logger *zap.Logger,
) *QuoteQueryHandler {
	return &QuoteQueryHandler{
		db:        db,
		quoteRepo: quoteRepo,
		logger:    logger,
	}
}

// HandleGetQuoteByID 處理根據ID獲取報價查詢
func (h *QuoteQueryHandler) HandleGetQuoteByID(ctx context.Context, q query.GetQuoteByIDQuery) (*query.QuoteDTO, error) {
	// 從倉儲獲取聚合
	quote, err := h.quoteRepo.FindByID(ctx, q.QuoteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("quote not found: %s", q.QuoteID)
		}
		return nil, err
	}
	
	// 轉換為DTO
	dto := h.mapToQuoteDTO(quote)
	
	// 加載關聯數據
	if err := h.loadRelatedData(ctx, dto); err != nil {
		h.logger.Error("Failed to load related data", zap.Error(err))
	}
	
	return dto, nil
}

// HandleGetQuoteByNumber 處理根據編號獲取報價查詢
func (h *QuoteQueryHandler) HandleGetQuoteByNumber(ctx context.Context, q query.GetQuoteByNumberQuery) (*query.QuoteDTO, error) {
	// 創建報價編號值對象
	quoteNumber, err := valueobject.NewQuoteNumber(q.QuoteNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid quote number: %w", err)
	}
	
	// 從倉儲獲取聚合
	quote, err := h.quoteRepo.FindByNumber(ctx, quoteNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("quote not found: %s", q.QuoteNumber)
		}
		return nil, err
	}
	
	// 轉換為DTO
	dto := h.mapToQuoteDTO(quote)
	
	// 加載關聯數據
	if err := h.loadRelatedData(ctx, dto); err != nil {
		h.logger.Error("Failed to load related data", zap.Error(err))
	}
	
	return dto, nil
}

// HandleListQuotes 處理列出報價查詢
func (h *QuoteQueryHandler) HandleListQuotes(ctx context.Context, q query.ListQuotesQuery) (*query.PageResult[query.QuoteDTO], error) {
	// 構建查詢規格
	spec := repository.QuerySpecification{
		Limit:     q.GetLimit(),
		Offset:    q.GetOffset(),
		OrderBy:   q.Sort,
		OrderDesc: q.Order == "desc",
		Filters:   make(map[string]interface{}),
		DateFrom:  q.DateFrom,
		DateTo:    q.DateTo,
	}
	
	// 構建查詢
	var quotes []*aggregate.QuoteAggregate
	var err error
	
	if q.CustomerID != nil {
		quotes, err = h.quoteRepo.FindByCustomer(ctx, *q.CustomerID, spec)
	} else if q.CompanyID != nil {
		quotes, err = h.quoteRepo.FindByCompany(ctx, *q.CompanyID, spec)
	} else if q.Status != nil {
		quotes, err = h.quoteRepo.FindByStatus(ctx, *q.Status, spec)
	} else {
		// 使用原生SQL查詢
		return h.listQuotesWithSQL(ctx, q)
	}
	
	if err != nil {
		return nil, err
	}
	
	// 計算總數
	count, err := h.quoteRepo.Count(ctx, spec)
	if err != nil {
		return nil, err
	}
	
	// 轉換為DTO
	dtos := make([]query.QuoteDTO, len(quotes))
	for i, quote := range quotes {
		dto := h.mapToQuoteDTO(quote)
		if err := h.loadRelatedData(ctx, dto); err != nil {
			h.logger.Error("Failed to load related data", zap.Error(err))
		}
		dtos[i] = *dto
	}
	
	return &query.PageResult[query.QuoteDTO]{
		Items:      dtos,
		TotalItems: count,
		Page:       q.Page,
		PageSize:   q.PageSize,
	}, nil
}

// HandleSearchQuotes 處理搜索報價查詢
func (h *QuoteQueryHandler) HandleSearchQuotes(ctx context.Context, q query.SearchQuotesQuery) (*query.PageResult[query.QuoteDTO], error) {
	// 使用原生SQL進行全文搜索
	var results []struct {
		ID          string
		QuoteNumber string
		CustomerID  string
		CompanyID   string
		Status      string
		ValidUntil  time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Version     int
	}
	
	query := h.db.Table("quotes").
		Select("DISTINCT quotes.*").
		Joins("LEFT JOIN quote_items ON quotes.id = quote_items.quote_id").
		Joins("LEFT JOIN customers ON quotes.customer_id = customers.id").
		Joins("LEFT JOIN products ON quote_items.product_id = products.id")
	
	// 添加搜索條件
	if q.Keyword != "" {
		query = query.Where(
			"quotes.quote_number LIKE ? OR customers.name LIKE ? OR products.name LIKE ?",
			"%"+q.Keyword+"%", "%"+q.Keyword+"%", "%"+q.Keyword+"%",
		)
	}
	
	if q.CustomerName != "" {
		query = query.Where("customers.name LIKE ?", "%"+q.CustomerName+"%")
	}
	
	if q.ProductName != "" {
		query = query.Where("products.name LIKE ?", "%"+q.ProductName+"%")
	}
	
	if q.MinAmount != nil {
		query = query.Where("quotes.total_amount >= ?", *q.MinAmount)
	}
	
	if q.MaxAmount != nil {
		query = query.Where("quotes.total_amount <= ?", *q.MaxAmount)
	}
	
	if q.CompanyID != nil {
		query = query.Where("quotes.company_id = ?", *q.CompanyID)
	}
	
	// 計算總數
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}
	
	// 分頁
	if err := query.
		Limit(q.GetLimit()).
		Offset(q.GetOffset()).
		Order("quotes.created_at DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	
	// 轉換為DTO
	dtos := make([]query.QuoteDTO, len(results))
	for i, result := range results {
		quoteID, _ := uuid.Parse(result.ID)
		quote, err := h.quoteRepo.FindByID(ctx, quoteID)
		if err != nil {
			continue
		}
		
		dto := h.mapToQuoteDTO(quote)
		if err := h.loadRelatedData(ctx, dto); err != nil {
			h.logger.Error("Failed to load related data", zap.Error(err))
		}
		dtos[i] = *dto
	}
	
	return &query.PageResult[query.QuoteDTO]{
		Items:      dtos,
		TotalItems: count,
		Page:       q.Page,
		PageSize:   q.PageSize,
	}, nil
}

// HandleGetExpiringQuotes 處理獲取即將過期的報價查詢
func (h *QuoteQueryHandler) HandleGetExpiringQuotes(ctx context.Context, q query.GetExpiringQuotesQuery) (*query.PageResult[query.QuoteDTO], error) {
	quotes, err := h.quoteRepo.FindExpiring(ctx, q.WithinDays)
	if err != nil {
		return nil, err
	}
	
	// 如果指定了公司ID，進行過濾
	if q.CompanyID != nil {
		filtered := make([]*aggregate.QuoteAggregate, 0)
		for _, quote := range quotes {
			if quote.CompanyID == *q.CompanyID {
				filtered = append(filtered, quote)
			}
		}
		quotes = filtered
	}
	
	// 分頁處理
	total := len(quotes)
	start := q.GetOffset()
	end := start + q.GetLimit()
	
	if start >= total {
		return &query.PageResult[query.QuoteDTO]{
			Items:      []query.QuoteDTO{},
			TotalItems: int64(total),
			Page:       q.Page,
			PageSize:   q.PageSize,
		}, nil
	}
	
	if end > total {
		end = total
	}
	
	pagedQuotes := quotes[start:end]
	
	// 轉換為DTO
	dtos := make([]query.QuoteDTO, len(pagedQuotes))
	for i, quote := range pagedQuotes {
		dto := h.mapToQuoteDTO(quote)
		if err := h.loadRelatedData(ctx, dto); err != nil {
			h.logger.Error("Failed to load related data", zap.Error(err))
		}
		dtos[i] = *dto
	}
	
	return &query.PageResult[query.QuoteDTO]{
		Items:      dtos,
		TotalItems: int64(total),
		Page:       q.Page,
		PageSize:   q.PageSize,
	}, nil
}

// HandleGetQuoteStatistics 處理獲取報價統計查詢
func (h *QuoteQueryHandler) HandleGetQuoteStatistics(ctx context.Context, q query.GetQuoteStatisticsQuery) (*query.QuoteStatisticsDTO, error) {
	stats := &query.QuoteStatisticsDTO{
		QuotesByStatus: make(map[string]int),
		TopCustomers:   make([]query.CustomerStatDTO, 0),
		TopProducts:    make([]query.ProductStatDTO, 0),
		MonthlyTrend:   make([]query.MonthlyTrendDTO, 0),
	}
	
	// 獲取總報價數和狀態分布
	var statusStats []struct {
		Status string
		Count  int
	}
	
	if err := h.db.Table("quotes").
		Select("status, COUNT(*) as count").
		Where("company_id = ? AND created_at BETWEEN ? AND ?", q.CompanyID, q.DateFrom, q.DateTo).
		Group("status").
		Scan(&statusStats).Error; err != nil {
		return nil, err
	}
	
	for _, stat := range statusStats {
		stats.QuotesByStatus[stat.Status] = stat.Count
		stats.TotalQuotes += stat.Count
	}
	
	// 獲取總價值和平均價值
	var valueStats struct {
		TotalValue   float64
		AverageValue float64
	}
	
	if err := h.db.Table("quotes").
		Select("SUM(total_amount) as total_value, AVG(total_amount) as average_value").
		Where("company_id = ? AND created_at BETWEEN ? AND ?", q.CompanyID, q.DateFrom, q.DateTo).
		Scan(&valueStats).Error; err != nil {
		return nil, err
	}
	
	stats.TotalValue = valueStats.TotalValue
	stats.AverageValue = valueStats.AverageValue
	
	// 計算轉換率
	approvedCount := stats.QuotesByStatus["APPROVED"]
	if stats.TotalQuotes > 0 {
		stats.ConversionRate = float64(approvedCount) / float64(stats.TotalQuotes) * 100
	}
	
	// 獲取前10大客戶
	if err := h.db.Table("quotes").
		Select("customer_id, customers.name as customer_name, COUNT(*) as quote_count, SUM(total_amount) as total_value").
		Joins("JOIN customers ON quotes.customer_id = customers.id").
		Where("quotes.company_id = ? AND quotes.created_at BETWEEN ? AND ?", q.CompanyID, q.DateFrom, q.DateTo).
		Group("customer_id, customers.name").
		Order("total_value DESC").
		Limit(10).
		Scan(&stats.TopCustomers).Error; err != nil {
		h.logger.Error("Failed to get top customers", zap.Error(err))
	}
	
	// 獲取前10大產品
	if err := h.db.Table("quote_items").
		Select("product_id, products.name as product_name, products.sku as product_sku, COUNT(DISTINCT quote_id) as quote_count, SUM(quantity) as quantity, SUM(total_price) as total_value").
		Joins("JOIN quotes ON quote_items.quote_id = quotes.id").
		Joins("JOIN products ON quote_items.product_id = products.id").
		Where("quotes.company_id = ? AND quotes.created_at BETWEEN ? AND ?", q.CompanyID, q.DateFrom, q.DateTo).
		Group("product_id, products.name, products.sku").
		Order("total_value DESC").
		Limit(10).
		Scan(&stats.TopProducts).Error; err != nil {
		h.logger.Error("Failed to get top products", zap.Error(err))
	}
	
	// 獲取月度趨勢
	if err := h.db.Table("quotes").
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as quote_count, SUM(total_amount) as total_value").
		Where("company_id = ? AND created_at BETWEEN ? AND ?", q.CompanyID, q.DateFrom, q.DateTo).
		Group("month").
		Order("month").
		Scan(&stats.MonthlyTrend).Error; err != nil {
		h.logger.Error("Failed to get monthly trend", zap.Error(err))
	}
	
	return stats, nil
}

// 私有方法

func (h *QuoteQueryHandler) mapToQuoteDTO(quote *aggregate.QuoteAggregate) *query.QuoteDTO {
	dto := &query.QuoteDTO{
		ID:          quote.ID,
		QuoteNumber: quote.QuoteNumber.String(),
		CustomerID:  quote.CustomerID,
		CompanyID:   quote.CompanyID,
		Status:      quote.Status,
		ValidUntil:  quote.ValidUntil,
		CreatedAt:   quote.CreatedAt,
		UpdatedAt:   quote.UpdatedAt,
		Version:     quote.Version,
		Items:       make([]query.QuoteItemDTO, len(quote.Items)),
		Terms:       h.mapToQuoteTermsDTO(quote.Terms),
		PricingSummary: query.PricingSummaryDTO{
			Subtotal:      quote.PricingSummary.Subtotal,
			TotalTax:      quote.PricingSummary.TotalTax,
			TotalDiscount: quote.PricingSummary.TotalDiscount,
			Total:         quote.PricingSummary.Total,
			Currency:      quote.PricingSummary.Currency.String(),
		},
	}
	
	// 映射項目
	for i, item := range quote.Items {
		dto.Items[i] = query.QuoteItemDTO{
			ID:             item.ID,
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			Specification:  item.Specification,
			Material:       h.mapToMaterialDTO(item.Material),
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			TaxRate:        item.TaxRate,
			DiscountRate:   item.DiscountRate,
			TotalPrice:     item.CalculateTotal(),
			TaxAmount:      item.CalculateTax(),
			DiscountAmount: item.CalculateDiscount(),
			FinalPrice:     item.GetFinalPrice(),
			LeadTime:       item.LeadTime.String(),
			Notes:          item.Notes,
		}
	}
	
	return dto
}

func (h *QuoteQueryHandler) mapToQuoteTermsDTO(terms valueobject.QuoteTerms) query.QuoteTermsDTO {
	return query.QuoteTermsDTO{
		PaymentTerms: query.PaymentTermsDTO{
			Type:           string(terms.PaymentTerms.Type),
			NetDays:        terms.PaymentTerms.NetDays,
			DepositPercent: terms.PaymentTerms.DepositPercent,
			Description:    terms.PaymentTerms.Description,
		},
		DeliveryTerms: query.DeliveryTermsDTO{
			Incoterm:     string(terms.DeliveryTerms.Incoterm),
			LeadTimeDays: terms.DeliveryTerms.LeadTimeDays,
			Location:     terms.DeliveryTerms.Location,
			Description:  terms.DeliveryTerms.Description,
		},
		WarrantyTerms: query.WarrantyTermsDTO{
			Duration:    terms.WarrantyTerms.Duration.String(),
			Type:        string(terms.WarrantyTerms.Type),
			Coverage:    terms.WarrantyTerms.Coverage,
			Exclusions:  terms.WarrantyTerms.Exclusions,
			Description: terms.WarrantyTerms.Description,
		},
		Currency:           terms.Currency.String(),
		DiscountPercentage: terms.DiscountPercentage,
		Notes:              terms.Notes,
	}
}

func (h *QuoteQueryHandler) mapToMaterialDTO(material entity.Material) query.MaterialDTO {
	return query.MaterialDTO{
		Type:        string(material.Type),
		Grade:       material.Grade,
		Standard:    material.Standard,
		Finish:      material.Finish,
		Description: material.Description,
	}
}

func (h *QuoteQueryHandler) loadRelatedData(ctx context.Context, dto *query.QuoteDTO) error {
	// 加載客戶名稱
	var customer struct {
		Name string
	}
	if err := h.db.Table("customers").
		Select("name").
		Where("id = ?", dto.CustomerID).
		First(&customer).Error; err == nil {
		dto.CustomerName = customer.Name
	}
	
	// 加載公司名稱
	var company struct {
		Name string
	}
	if err := h.db.Table("companies").
		Select("name").
		Where("id = ?", dto.CompanyID).
		First(&company).Error; err == nil {
		dto.CompanyName = company.Name
	}
	
	// 加載產品SKU
	for i := range dto.Items {
		var product struct {
			SKU string
		}
		if err := h.db.Table("products").
			Select("sku").
			Where("id = ?", dto.Items[i].ProductID).
			First(&product).Error; err == nil {
			dto.Items[i].ProductSKU = product.SKU
		}
	}
	
	return nil
}

func (h *QuoteQueryHandler) listQuotesWithSQL(ctx context.Context, q query.ListQuotesQuery) (*query.PageResult[query.QuoteDTO], error) {
	// 使用原生SQL查詢以提高性能
	query := h.db.Table("quotes")
	
	// 添加過濾條件
	if q.CustomerID != nil {
		query = query.Where("customer_id = ?", *q.CustomerID)
	}
	
	if q.CompanyID != nil {
		query = query.Where("company_id = ?", *q.CompanyID)
	}
	
	if q.Status != nil {
		query = query.Where("status = ?", *q.Status)
	}
	
	if q.DateFrom != nil {
		query = query.Where("created_at >= ?", *q.DateFrom)
	}
	
	if q.DateTo != nil {
		query = query.Where("created_at <= ?", *q.DateTo)
	}
	
	// 計算總數
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}
	
	// 查詢數據
	var results []struct {
		ID          string
		QuoteNumber string
		CustomerID  string
		CompanyID   string
		Status      string
		ValidUntil  time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Version     int
	}
	
	// 排序和分頁
	orderBy := "created_at"
	if q.Sort != "" {
		orderBy = q.Sort
	}
	
	orderDirection := "DESC"
	if q.Order == "asc" {
		orderDirection = "ASC"
	}
	
	if err := query.
		Order(fmt.Sprintf("%s %s", orderBy, orderDirection)).
		Limit(q.GetLimit()).
		Offset(q.GetOffset()).
		Scan(&results).Error; err != nil {
		return nil, err
	}
	
	// 轉換為DTO
	dtos := make([]query.QuoteDTO, len(results))
	for i, result := range results {
		quoteID, _ := uuid.Parse(result.ID)
		quote, err := h.quoteRepo.FindByID(ctx, quoteID)
		if err != nil {
			continue
		}
		
		dto := h.mapToQuoteDTO(quote)
		if err := h.loadRelatedData(ctx, dto); err != nil {
			h.logger.Error("Failed to load related data", zap.Error(err))
		}
		dtos[i] = *dto
	}
	
	return &query.PageResult[query.QuoteDTO]{
		Items:      dtos,
		TotalItems: count,
		Page:       q.Page,
		PageSize:   q.PageSize,
	}, nil
}
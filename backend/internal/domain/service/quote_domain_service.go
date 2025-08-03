package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/domain/aggregate"
	"github.com/fastenmind/fastener-api/internal/domain/entity"
	"github.com/fastenmind/fastener-api/internal/domain/repository"
	"github.com/fastenmind/fastener-api/internal/domain/valueobject"
)

// QuoteDomainService 報價領域服務
type QuoteDomainService struct {
	quoteRepo      repository.QuoteRepository
	customerRepo   repository.CustomerRepository
	productRepo    repository.ProductRepository
	pricingService PricingService
	eventPublisher EventPublisher
}

// NewQuoteDomainService 創建報價領域服務
func NewQuoteDomainService(
	quoteRepo repository.QuoteRepository,
	customerRepo repository.CustomerRepository,
	productRepo repository.ProductRepository,
	pricingService PricingService,
	eventPublisher EventPublisher,
) *QuoteDomainService {
	return &QuoteDomainService{
		quoteRepo:      quoteRepo,
		customerRepo:   customerRepo,
		productRepo:    productRepo,
		pricingService: pricingService,
		eventPublisher: eventPublisher,
	}
}

// CreateQuote 創建報價
func (s *QuoteDomainService) CreateQuote(ctx context.Context, request CreateQuoteRequest) (*aggregate.QuoteAggregate, error) {
	// 驗證客戶存在
	customer, err := s.customerRepo.FindByID(ctx, request.CustomerID)
	if err != nil {
		return nil, errors.New("customer not found")
	}
	
	// 檢查客戶信用狀態
	if !customer.CanCreateQuote() {
		return nil, errors.New("customer credit status does not allow new quotes")
	}
	
	// 創建報價聚合
	quote, err := aggregate.NewQuoteAggregate(request.CustomerID, request.CompanyID)
	if err != nil {
		return nil, err
	}
	
	// 設置報價條款
	if request.Terms != nil {
		if err := quote.UpdateTerms(*request.Terms); err != nil {
			return nil, err
		}
	}
	
	// 添加報價項目
	for _, itemRequest := range request.Items {
		// 驗證產品存在
		product, err := s.productRepo.FindByID(ctx, itemRequest.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + itemRequest.ProductID.String())
		}
		
		// 創建報價項目
		item, err := entity.NewQuoteItem(
			product.ID,
			product.Name,
			itemRequest.Quantity,
			0, // 初始單價為0，稍後計算
		)
		if err != nil {
			return nil, err
		}
		
		// 設置規格和材料
		item.Specification = itemRequest.Specification
		item.Material = itemRequest.Material
		
		// 計算定價
		pricing, err := s.pricingService.CalculatePrice(ctx, PricingRequest{
			ProductID:    product.ID,
			Quantity:     itemRequest.Quantity,
			CustomerID:   request.CustomerID,
			Material:     itemRequest.Material,
			Specification: itemRequest.Specification,
		})
		if err != nil {
			return nil, err
		}
		
		// 更新項目定價
		item.UpdateUnitPrice(pricing.UnitPrice)
		item.SetTaxRate(pricing.TaxRate)
		item.SetDiscountRate(pricing.DiscountRate)
		item.SetLeadTime(pricing.LeadTime)
		
		// 添加到報價
		if err := quote.AddItem(*item); err != nil {
			return nil, err
		}
	}
	
	// 保存報價
	if err := s.quoteRepo.Save(ctx, quote); err != nil {
		return nil, err
	}
	
	// 發布領域事件
	s.publishEvents(quote)
	
	return quote, nil
}

// UpdateQuoteItems 更新報價項目
func (s *QuoteDomainService) UpdateQuoteItems(ctx context.Context, quoteID uuid.UUID, updates []ItemUpdate) error {
	// 獲取報價
	quote, err := s.quoteRepo.FindByID(ctx, quoteID)
	if err != nil {
		return err
	}
	
	// 處理每個更新
	for _, update := range updates {
		switch update.Action {
		case "add":
			// 驗證產品
			product, err := s.productRepo.FindByID(ctx, update.ProductID)
			if err != nil {
				return errors.New("product not found")
			}
			
			// 創建新項目
			item, err := entity.NewQuoteItem(
				product.ID,
				product.Name,
				update.Quantity,
				0,
			)
			if err != nil {
				return err
			}
			
			// 計算定價
			pricing, err := s.pricingService.CalculatePrice(ctx, PricingRequest{
				ProductID:  product.ID,
				Quantity:   update.Quantity,
				CustomerID: quote.CustomerID,
			})
			if err != nil {
				return err
			}
			
			item.UpdateUnitPrice(pricing.UnitPrice)
			item.SetTaxRate(pricing.TaxRate)
			
			if err := quote.AddItem(*item); err != nil {
				return err
			}
			
		case "update":
			err := quote.UpdateItem(update.ItemID, func(item *entity.QuoteItem) error {
				if update.Quantity > 0 {
					item.UpdateQuantity(update.Quantity)
				}
				if update.UnitPrice > 0 {
					item.UpdateUnitPrice(update.UnitPrice)
				}
				return nil
			})
			if err != nil {
				return err
			}
			
		case "remove":
			if err := quote.RemoveItem(update.ItemID); err != nil {
				return err
			}
		}
	}
	
	// 保存更新
	if err := s.quoteRepo.Update(ctx, quote); err != nil {
		return err
	}
	
	// 發布事件
	s.publishEvents(quote)
	
	return nil
}

// SubmitQuote 提交報價
func (s *QuoteDomainService) SubmitQuote(ctx context.Context, quoteID uuid.UUID) error {
	quote, err := s.quoteRepo.FindByID(ctx, quoteID)
	if err != nil {
		return err
	}
	
	// 重新計算最新價格
	if err := s.recalculatePrices(ctx, quote); err != nil {
		return err
	}
	
	// 提交報價
	if err := quote.Submit(); err != nil {
		return err
	}
	
	// 保存
	if err := s.quoteRepo.Update(ctx, quote); err != nil {
		return err
	}
	
	// 發布事件
	s.publishEvents(quote)
	
	return nil
}

// ApproveQuote 批准報價
func (s *QuoteDomainService) ApproveQuote(ctx context.Context, quoteID, approverID uuid.UUID) error {
	quote, err := s.quoteRepo.FindByID(ctx, quoteID)
	if err != nil {
		return err
	}
	
	// 批准報價
	if err := quote.Approve(approverID); err != nil {
		return err
	}
	
	// 保存
	if err := s.quoteRepo.Update(ctx, quote); err != nil {
		return err
	}
	
	// 發布事件
	s.publishEvents(quote)
	
	return nil
}

// RejectQuote 拒絕報價
func (s *QuoteDomainService) RejectQuote(ctx context.Context, quoteID uuid.UUID, reason string) error {
	quote, err := s.quoteRepo.FindByID(ctx, quoteID)
	if err != nil {
		return err
	}
	
	// 拒絕報價
	if err := quote.Reject(reason); err != nil {
		return err
	}
	
	// 保存
	if err := s.quoteRepo.Update(ctx, quote); err != nil {
		return err
	}
	
	// 發布事件
	s.publishEvents(quote)
	
	return nil
}

// ExtendQuoteValidity 延長報價有效期
func (s *QuoteDomainService) ExtendQuoteValidity(ctx context.Context, quoteID uuid.UUID, newValidUntil time.Time) error {
	quote, err := s.quoteRepo.FindByID(ctx, quoteID)
	if err != nil {
		return err
	}
	
	// 延長有效期
	if err := quote.ExtendValidity(newValidUntil); err != nil {
		return err
	}
	
	// 保存
	if err := s.quoteRepo.Update(ctx, quote); err != nil {
		return err
	}
	
	// 發布事件
	s.publishEvents(quote)
	
	return nil
}

// CloneQuote 複製報價
func (s *QuoteDomainService) CloneQuote(ctx context.Context, quoteID uuid.UUID) (*aggregate.QuoteAggregate, error) {
	originalQuote, err := s.quoteRepo.FindByID(ctx, quoteID)
	if err != nil {
		return nil, err
	}
	
	// 複製報價
	newQuote, err := originalQuote.Clone()
	if err != nil {
		return nil, err
	}
	
	// 重新計算價格（可能價格規則已更新）
	if err := s.recalculatePrices(ctx, newQuote); err != nil {
		return nil, err
	}
	
	// 保存新報價
	if err := s.quoteRepo.Save(ctx, newQuote); err != nil {
		return nil, err
	}
	
	// 發布事件
	s.publishEvents(newQuote)
	
	return newQuote, nil
}

// ProcessExpiredQuotes 處理過期報價
func (s *QuoteDomainService) ProcessExpiredQuotes(ctx context.Context) error {
	// 查找即將過期的報價
	quotes, err := s.quoteRepo.FindExpiring(ctx, 0)
	if err != nil {
		return err
	}
	
	for _, quote := range quotes {
		if err := quote.Expire(); err != nil {
			continue // 記錄錯誤但繼續處理其他報價
		}
		
		if err := s.quoteRepo.Update(ctx, quote); err != nil {
			continue
		}
		
		s.publishEvents(quote)
	}
	
	return nil
}

// 私有方法

func (s *QuoteDomainService) recalculatePrices(ctx context.Context, quote *aggregate.QuoteAggregate) error {
	for i := range quote.Items {
		item := &quote.Items[i]
		
		pricing, err := s.pricingService.CalculatePrice(ctx, PricingRequest{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			CustomerID:    quote.CustomerID,
			Material:      item.Material,
			Specification: item.Specification,
		})
		if err != nil {
			return err
		}
		
		item.UpdateUnitPrice(pricing.UnitPrice)
		item.SetTaxRate(pricing.TaxRate)
		item.SetDiscountRate(pricing.DiscountRate)
	}
	
	return nil
}

func (s *QuoteDomainService) publishEvents(quote *aggregate.QuoteAggregate) {
	events := quote.GetDomainEvents()
	for _, event := range events {
		s.eventPublisher.Publish(event)
	}
	quote.ClearDomainEvents()
}

// 請求和響應結構

// CreateQuoteRequest 創建報價請求
type CreateQuoteRequest struct {
	CustomerID uuid.UUID
	CompanyID  uuid.UUID
	Terms      *valueobject.QuoteTerms
	Items      []CreateQuoteItemRequest
}

// CreateQuoteItemRequest 創建報價項目請求
type CreateQuoteItemRequest struct {
	ProductID     uuid.UUID
	Quantity      int
	Specification string
	Material      entity.Material
}

// ItemUpdate 項目更新
type ItemUpdate struct {
	Action    string    // add, update, remove
	ItemID    uuid.UUID // 用於 update 和 remove
	ProductID uuid.UUID // 用於 add
	Quantity  int
	UnitPrice float64
}

// PricingRequest 定價請求
type PricingRequest struct {
	ProductID     uuid.UUID
	Quantity      int
	CustomerID    uuid.UUID
	Material      entity.Material
	Specification string
}

// PricingResponse 定價響應
type PricingResponse struct {
	UnitPrice    float64
	TaxRate      float64
	DiscountRate float64
	LeadTime     time.Duration
}
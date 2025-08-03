package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	quotev1 "github.com/fastenmind/quote-service/api/v1"
	"github.com/fastenmind/fastener-api/internal/application/command"
	"github.com/fastenmind/fastener-api/internal/application/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"go.uber.org/zap"
)

// QuoteService gRPC 報價服務實現
type QuoteService struct {
	quotev1.UnimplementedQuoteServiceServer
	commandBus command.Bus
	queryBus   query.Bus
	logger     *zap.Logger
}

// NewQuoteService 創建報價服務
func NewQuoteService(commandBus command.Bus, queryBus query.Bus, logger *zap.Logger) *QuoteService {
	return &QuoteService{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

// CreateQuote 創建報價
func (s *QuoteService) CreateQuote(ctx context.Context, req *quotev1.CreateQuoteRequest) (*quotev1.CreateQuoteResponse, error) {
	// 驗證請求
	if req.CustomerId == "" {
		return nil, status.Error(codes.InvalidArgument, "customer_id is required")
	}
	
	if req.CompanyId == "" {
		return nil, status.Error(codes.InvalidArgument, "company_id is required")
	}
	
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one item is required")
	}
	
	// 構建命令
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer_id")
	}
	
	companyID, err := uuid.Parse(req.CompanyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid company_id")
	}
	
	cmd := command.NewCreateQuoteCommand(customerID, companyID)
	
	// 轉換項目
	for _, item := range req.Items {
		productID, err := uuid.Parse(item.ProductId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid product_id")
		}
		
		cmd.Items = append(cmd.Items, command.CreateQuoteItemCommand{
			ProductID:     productID,
			Quantity:      int(item.Quantity),
			Specification: item.Specification,
		})
	}
	
	// 發送命令
	if err := s.commandBus.Send(ctx, cmd); err != nil {
		s.logger.Error("Failed to create quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create quote")
	}
	
	// 返回響應
	return &quotev1.CreateQuoteResponse{
		Id:          cmd.GetCommandID().String(),
		QuoteNumber: "Q-" + cmd.GetCommandID().String()[:8],
	}, nil
}

// GetQuote 獲取報價
func (s *QuoteService) GetQuote(ctx context.Context, req *quotev1.GetQuoteRequest) (*quotev1.Quote, error) {
	// 驗證請求
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	
	quoteID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	
	// 構建查詢
	q := query.NewGetQuoteByIDQuery(quoteID)
	
	// 發送查詢
	result, err := s.queryBus.Send(ctx, q)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "quote not found")
		}
		s.logger.Error("Failed to get quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get quote")
	}
	
	// 轉換響應
	quoteDTO := result.(*query.QuoteDTO)
	return s.mapToProtoQuote(quoteDTO), nil
}

// ListQuotes 列出報價
func (s *QuoteService) ListQuotes(ctx context.Context, req *quotev1.ListQuotesRequest) (*quotev1.ListQuotesResponse, error) {
	// 構建查詢
	q := query.NewListQuotesQuery()
	
	if req.CustomerId != "" {
		customerID, err := uuid.Parse(req.CustomerId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid customer_id")
		}
		q.CustomerID = &customerID
	}
	
	if req.CompanyId != "" {
		companyID, err := uuid.Parse(req.CompanyId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid company_id")
		}
		q.CompanyID = &companyID
	}
	
	if req.Page > 0 {
		q.Page = int(req.Page)
	}
	
	if req.PageSize > 0 {
		q.PageSize = int(req.PageSize)
	}
	
	if req.SortBy != "" {
		q.Sort = req.SortBy
	}
	
	if req.SortDesc {
		q.Order = "desc"
	}
	
	// 發送查詢
	result, err := s.queryBus.Send(ctx, q)
	if err != nil {
		s.logger.Error("Failed to list quotes", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list quotes")
	}
	
	// 轉換響應
	pageResult := result.(*query.PageResult[query.QuoteDTO])
	
	quotes := make([]*quotev1.Quote, len(pageResult.Items))
	for i, dto := range pageResult.Items {
		quotes[i] = s.mapToProtoQuote(&dto)
	}
	
	return &quotev1.ListQuotesResponse{
		Quotes:     quotes,
		TotalCount: pageResult.TotalItems,
		Page:       int32(pageResult.Page),
		PageSize:   int32(pageResult.PageSize),
	}, nil
}

// UpdateQuote 更新報價
func (s *QuoteService) UpdateQuote(ctx context.Context, req *quotev1.UpdateQuoteRequest) (*quotev1.Quote, error) {
	// 驗證請求
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	
	quoteID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	
	// 構建命令
	cmd := command.NewUpdateQuoteItemsCommand(quoteID)
	
	// TODO: 實現更新邏輯
	
	// 發送命令
	if err := s.commandBus.Send(ctx, cmd); err != nil {
		s.logger.Error("Failed to update quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update quote")
	}
	
	// 獲取更新後的報價
	return s.GetQuote(ctx, &quotev1.GetQuoteRequest{Id: req.Id})
}

// SubmitQuote 提交報價
func (s *QuoteService) SubmitQuote(ctx context.Context, req *quotev1.SubmitQuoteRequest) (*emptypb.Empty, error) {
	// 驗證請求
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	
	quoteID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	
	// 構建命令
	cmd := command.NewSubmitQuoteCommand(quoteID)
	
	// 發送命令
	if err := s.commandBus.Send(ctx, cmd); err != nil {
		s.logger.Error("Failed to submit quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to submit quote")
	}
	
	return &emptypb.Empty{}, nil
}

// ApproveQuote 批准報價
func (s *QuoteService) ApproveQuote(ctx context.Context, req *quotev1.ApproveQuoteRequest) (*emptypb.Empty, error) {
	// 驗證請求
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	
	if req.ApproverId == "" {
		return nil, status.Error(codes.InvalidArgument, "approver_id is required")
	}
	
	quoteID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	
	approverID, err := uuid.Parse(req.ApproverId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid approver_id")
	}
	
	// 構建命令
	cmd := command.NewApproveQuoteCommand(quoteID, approverID)
	
	// 發送命令
	if err := s.commandBus.Send(ctx, cmd); err != nil {
		s.logger.Error("Failed to approve quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to approve quote")
	}
	
	return &emptypb.Empty{}, nil
}

// RejectQuote 拒絕報價
func (s *QuoteService) RejectQuote(ctx context.Context, req *quotev1.RejectQuoteRequest) (*emptypb.Empty, error) {
	// 驗證請求
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	
	if req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "reason is required")
	}
	
	quoteID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	
	// 構建命令
	cmd := command.NewRejectQuoteCommand(quoteID, req.Reason)
	
	// 發送命令
	if err := s.commandBus.Send(ctx, cmd); err != nil {
		s.logger.Error("Failed to reject quote", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to reject quote")
	}
	
	return &emptypb.Empty{}, nil
}

// GetQuoteStatistics 獲取報價統計
func (s *QuoteService) GetQuoteStatistics(ctx context.Context, req *quotev1.GetQuoteStatisticsRequest) (*quotev1.QuoteStatistics, error) {
	// 驗證請求
	if req.CompanyId == "" {
		return nil, status.Error(codes.InvalidArgument, "company_id is required")
	}
	
	companyID, err := uuid.Parse(req.CompanyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid company_id")
	}
	
	// 構建查詢
	q := query.NewGetQuoteStatisticsQuery(
		companyID,
		req.StartDate.AsTime(),
		req.EndDate.AsTime(),
	)
	
	// 發送查詢
	result, err := s.queryBus.Send(ctx, q)
	if err != nil {
		s.logger.Error("Failed to get quote statistics", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get quote statistics")
	}
	
	// 轉換響應
	stats := result.(*query.QuoteStatisticsDTO)
	
	return &quotev1.QuoteStatistics{
		TotalQuotes:     int32(stats.TotalQuotes),
		QuotesByStatus:  s.mapToInt32Map(stats.QuotesByStatus),
		TotalValue:      stats.TotalValue,
		AverageValue:    stats.AverageValue,
		ConversionRate:  stats.ConversionRate,
		TopCustomers:    s.mapToProtoCustomerStats(stats.TopCustomers),
		TopProducts:     s.mapToProtoProductStats(stats.TopProducts),
	}, nil
}

// 私有方法

func (s *QuoteService) mapToProtoQuote(dto *query.QuoteDTO) *quotev1.Quote {
	return &quotev1.Quote{
		Id:           dto.ID.String(),
		QuoteNumber:  dto.QuoteNumber,
		CustomerId:   dto.CustomerID.String(),
		CompanyId:    dto.CompanyID.String(),
		Status:       s.mapToProtoStatus(dto.Status),
		ValidUntil:   timestamppb.New(dto.ValidUntil),
		Items:        s.mapToProtoItems(dto.Items),
		Terms:        s.mapToProtoTerms(dto.Terms),
		PricingSummary: &quotev1.PricingSummary{
			Subtotal:      dto.PricingSummary.Subtotal,
			TotalTax:      dto.PricingSummary.TotalTax,
			TotalDiscount: dto.PricingSummary.TotalDiscount,
			Total:         dto.PricingSummary.Total,
			Currency:      dto.PricingSummary.Currency,
		},
		CreatedAt: timestamppb.New(dto.CreatedAt),
		UpdatedAt: timestamppb.New(dto.UpdatedAt),
		Version:   int32(dto.Version),
	}
}

func (s *QuoteService) mapToProtoStatus(status valueobject.QuoteStatus) quotev1.QuoteStatus {
	switch status {
	case valueobject.QuoteStatusDraft:
		return quotev1.QuoteStatus_QUOTE_STATUS_DRAFT
	case valueobject.QuoteStatusPending:
		return quotev1.QuoteStatus_QUOTE_STATUS_PENDING
	case valueobject.QuoteStatusApproved:
		return quotev1.QuoteStatus_QUOTE_STATUS_APPROVED
	case valueobject.QuoteStatusRejected:
		return quotev1.QuoteStatus_QUOTE_STATUS_REJECTED
	case valueobject.QuoteStatusExpired:
		return quotev1.QuoteStatus_QUOTE_STATUS_EXPIRED
	case valueobject.QuoteStatusCancelled:
		return quotev1.QuoteStatus_QUOTE_STATUS_CANCELLED
	default:
		return quotev1.QuoteStatus_QUOTE_STATUS_UNSPECIFIED
	}
}

func (s *QuoteService) mapToProtoItems(items []query.QuoteItemDTO) []*quotev1.QuoteItem {
	result := make([]*quotev1.QuoteItem, len(items))
	for i, item := range items {
		result[i] = &quotev1.QuoteItem{
			Id:           item.ID.String(),
			ProductId:    item.ProductID.String(),
			ProductName:  item.ProductName,
			Specification: item.Specification,
			Material: &quotev1.Material{
				Type:        item.Material.Type,
				Grade:       item.Material.Grade,
				Standard:    item.Material.Standard,
				Finish:      item.Material.Finish,
				Description: item.Material.Description,
			},
			Quantity:     int32(item.Quantity),
			UnitPrice:    item.UnitPrice,
			TaxRate:      item.TaxRate,
			DiscountRate: item.DiscountRate,
			TotalPrice:   item.TotalPrice,
			LeadTime:     item.LeadTime,
			Notes:        item.Notes,
		}
	}
	return result
}

func (s *QuoteService) mapToProtoTerms(terms query.QuoteTermsDTO) *quotev1.QuoteTerms {
	return &quotev1.QuoteTerms{
		PaymentTerms: &quotev1.PaymentTerms{
			Type:           terms.PaymentTerms.Type,
			NetDays:        int32(terms.PaymentTerms.NetDays),
			DepositPercent: terms.PaymentTerms.DepositPercent,
			Description:    terms.PaymentTerms.Description,
		},
		DeliveryTerms: &quotev1.DeliveryTerms{
			Incoterm:     terms.DeliveryTerms.Incoterm,
			LeadTimeDays: int32(terms.DeliveryTerms.LeadTimeDays),
			Location:     terms.DeliveryTerms.Location,
			Description:  terms.DeliveryTerms.Description,
		},
		WarrantyTerms: &quotev1.WarrantyTerms{
			Duration:    terms.WarrantyTerms.Duration,
			Type:        terms.WarrantyTerms.Type,
			Coverage:    terms.WarrantyTerms.Coverage,
			Exclusions:  terms.WarrantyTerms.Exclusions,
			Description: terms.WarrantyTerms.Description,
		},
		Currency:           terms.Currency,
		DiscountPercentage: terms.DiscountPercentage,
		Notes:              terms.Notes,
	}
}

func (s *QuoteService) mapToInt32Map(m map[string]int) map[string]int32 {
	result := make(map[string]int32)
	for k, v := range m {
		result[k] = int32(v)
	}
	return result
}

func (s *QuoteService) mapToProtoCustomerStats(stats []query.CustomerStatDTO) []*quotev1.CustomerStat {
	result := make([]*quotev1.CustomerStat, len(stats))
	for i, stat := range stats {
		result[i] = &quotev1.CustomerStat{
			CustomerId:   stat.CustomerID.String(),
			CustomerName: stat.CustomerName,
			QuoteCount:   int32(stat.QuoteCount),
			TotalValue:   stat.TotalValue,
		}
	}
	return result
}

func (s *QuoteService) mapToProtoProductStats(stats []query.ProductStatDTO) []*quotev1.ProductStat {
	result := make([]*quotev1.ProductStat, len(stats))
	for i, stat := range stats {
		result[i] = &quotev1.ProductStat{
			ProductId:     stat.ProductID.String(),
			ProductName:   stat.ProductName,
			ProductSku:    stat.ProductSKU,
			QuoteCount:    int32(stat.QuoteCount),
			TotalQuantity: int32(stat.Quantity),
			TotalValue:    stat.TotalValue,
		}
	}
	return result
}
package handler

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/fastenmind/fastener-api/internal/application/command"
	"github.com/fastenmind/fastener-api/internal/domain/aggregate"
	"github.com/fastenmind/fastener-api/internal/domain/repository"
	"github.com/fastenmind/fastener-api/internal/domain/service"
	"github.com/fastenmind/fastener-api/pkg/database"
	"go.uber.org/zap"
)

// QuoteCommandHandler 報價命令處理器
type QuoteCommandHandler struct {
	quoteDomainService *service.QuoteDomainService
	quoteRepo         repository.QuoteRepository
	unitOfWork        database.UnitOfWork
	logger            *zap.Logger
}

// NewQuoteCommandHandler 創建報價命令處理器
func NewQuoteCommandHandler(
	quoteDomainService *service.QuoteDomainService,
	quoteRepo repository.QuoteRepository,
	unitOfWork database.UnitOfWork,
	logger *zap.Logger,
) *QuoteCommandHandler {
	return &QuoteCommandHandler{
		quoteDomainService: quoteDomainService,
		quoteRepo:         quoteRepo,
		unitOfWork:        unitOfWork,
		logger:            logger,
	}
}

// HandleCreateQuote 處理創建報價命令
func (h *QuoteCommandHandler) HandleCreateQuote(ctx context.Context, cmd command.CreateQuoteCommand) (uuid.UUID, error) {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return uuid.Nil, err
	}
	
	// 開始事務
	tx, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()
	
	// 構建領域服務請求
	request := service.CreateQuoteRequest{
		CustomerID: cmd.CustomerID,
		CompanyID:  cmd.CompanyID,
		Terms:      cmd.Terms,
		Items:      make([]service.CreateQuoteItemRequest, len(cmd.Items)),
	}
	
	for i, item := range cmd.Items {
		request.Items[i] = service.CreateQuoteItemRequest{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			Specification: item.Specification,
			Material:      item.Material,
		}
	}
	
	// 創建報價
	quote, err := h.quoteDomainService.CreateQuote(ctx, request)
	if err != nil {
		h.logger.Error("Failed to create quote", 
			zap.String("command_id", cmd.GetCommandID().String()),
			zap.Error(err))
		return uuid.Nil, err
	}
	
	// 提交事務
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	
	h.logger.Info("Quote created successfully",
		zap.String("quote_id", quote.ID.String()),
		zap.String("quote_number", quote.QuoteNumber.String()))
	
	return quote.ID, nil
}

// HandleUpdateQuoteItems 處理更新報價項目命令
func (h *QuoteCommandHandler) HandleUpdateQuoteItems(ctx context.Context, cmd command.UpdateQuoteItemsCommand) error {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return err
	}
	
	// 開始事務
	tx, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// 構建領域服務請求
	updates := make([]service.ItemUpdate, len(cmd.Updates))
	for i, update := range cmd.Updates {
		updates[i] = service.ItemUpdate{
			Action:    update.Action,
			ItemID:    update.ItemID,
			ProductID: update.ProductID,
			Quantity:  update.Quantity,
			UnitPrice: update.UnitPrice,
		}
	}
	
	// 更新報價項目
	if err := h.quoteDomainService.UpdateQuoteItems(ctx, cmd.QuoteID, updates); err != nil {
		h.logger.Error("Failed to update quote items",
			zap.String("quote_id", cmd.QuoteID.String()),
			zap.Error(err))
		return err
	}
	
	// 提交事務
	if err := tx.Commit(); err != nil {
		return err
	}
	
	h.logger.Info("Quote items updated successfully",
		zap.String("quote_id", cmd.QuoteID.String()),
		zap.Int("updates_count", len(updates)))
	
	return nil
}

// HandleSubmitQuote 處理提交報價命令
func (h *QuoteCommandHandler) HandleSubmitQuote(ctx context.Context, cmd command.SubmitQuoteCommand) error {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return err
	}
	
	// 開始事務
	tx, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// 提交報價
	if err := h.quoteDomainService.SubmitQuote(ctx, cmd.QuoteID); err != nil {
		h.logger.Error("Failed to submit quote",
			zap.String("quote_id", cmd.QuoteID.String()),
			zap.Error(err))
		return err
	}
	
	// 提交事務
	if err := tx.Commit(); err != nil {
		return err
	}
	
	h.logger.Info("Quote submitted successfully",
		zap.String("quote_id", cmd.QuoteID.String()))
	
	return nil
}

// HandleApproveQuote 處理批准報價命令
func (h *QuoteCommandHandler) HandleApproveQuote(ctx context.Context, cmd command.ApproveQuoteCommand) error {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return err
	}
	
	// 開始事務
	tx, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// 批准報價
	if err := h.quoteDomainService.ApproveQuote(ctx, cmd.QuoteID, cmd.ApproverID); err != nil {
		h.logger.Error("Failed to approve quote",
			zap.String("quote_id", cmd.QuoteID.String()),
			zap.String("approver_id", cmd.ApproverID.String()),
			zap.Error(err))
		return err
	}
	
	// 提交事務
	if err := tx.Commit(); err != nil {
		return err
	}
	
	h.logger.Info("Quote approved successfully",
		zap.String("quote_id", cmd.QuoteID.String()),
		zap.String("approver_id", cmd.ApproverID.String()))
	
	return nil
}

// HandleRejectQuote 處理拒絕報價命令
func (h *QuoteCommandHandler) HandleRejectQuote(ctx context.Context, cmd command.RejectQuoteCommand) error {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return err
	}
	
	// 開始事務
	tx, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// 拒絕報價
	if err := h.quoteDomainService.RejectQuote(ctx, cmd.QuoteID, cmd.Reason); err != nil {
		h.logger.Error("Failed to reject quote",
			zap.String("quote_id", cmd.QuoteID.String()),
			zap.String("reason", cmd.Reason),
			zap.Error(err))
		return err
	}
	
	// 提交事務
	if err := tx.Commit(); err != nil {
		return err
	}
	
	h.logger.Info("Quote rejected successfully",
		zap.String("quote_id", cmd.QuoteID.String()),
		zap.String("reason", cmd.Reason))
	
	return nil
}

// HandleExtendQuoteValidity 處理延長報價有效期命令
func (h *QuoteCommandHandler) HandleExtendQuoteValidity(ctx context.Context, cmd command.ExtendQuoteValidityCommand) error {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return err
	}
	
	// 開始事務
	tx, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// 延長有效期
	if err := h.quoteDomainService.ExtendQuoteValidity(ctx, cmd.QuoteID, cmd.NewValidUntil); err != nil {
		h.logger.Error("Failed to extend quote validity",
			zap.String("quote_id", cmd.QuoteID.String()),
			zap.Time("new_valid_until", cmd.NewValidUntil),
			zap.Error(err))
		return err
	}
	
	// 提交事務
	if err := tx.Commit(); err != nil {
		return err
	}
	
	h.logger.Info("Quote validity extended successfully",
		zap.String("quote_id", cmd.QuoteID.String()),
		zap.Time("new_valid_until", cmd.NewValidUntil))
	
	return nil
}

// HandleCloneQuote 處理複製報價命令
func (h *QuoteCommandHandler) HandleCloneQuote(ctx context.Context, cmd command.CloneQuoteCommand) (uuid.UUID, error) {
	// 驗證命令
	if err := cmd.Validate(); err != nil {
		return uuid.Nil, err
	}
	
	// 開始事務
	tx, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()
	
	// 複製報價
	newQuote, err := h.quoteDomainService.CloneQuote(ctx, cmd.QuoteID)
	if err != nil {
		h.logger.Error("Failed to clone quote",
			zap.String("original_quote_id", cmd.QuoteID.String()),
			zap.Error(err))
		return uuid.Nil, err
	}
	
	// 提交事務
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	
	h.logger.Info("Quote cloned successfully",
		zap.String("original_quote_id", cmd.QuoteID.String()),
		zap.String("new_quote_id", newQuote.ID.String()))
	
	return newQuote.ID, nil
}
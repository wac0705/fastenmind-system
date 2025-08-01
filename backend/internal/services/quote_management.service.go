package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuoteManagementService struct {
	db             *gorm.DB
	quoteRepo      *repositories.QuoteManagementRepository
	pdfService     *PDFGeneratorService
	emailService   *EmailService
	webhookService *WebhookService
}

func NewQuoteManagementService(db *gorm.DB, webhookService *WebhookService) *QuoteManagementService {
	return &QuoteManagementService{
		db:             db,
		quoteRepo:      repositories.NewQuoteManagementRepository(db),
		pdfService:     NewPDFGeneratorService(),
		emailService:   NewEmailServiceDefault(),
		webhookService: webhookService,
	}
}

// CreateQuote 創建報價單
func (s *QuoteManagementService) CreateQuote(req models.CreateQuoteRequest, createdBy uuid.UUID) (*models.Quote, error) {
	// 開始事務
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 生成報價單號
	quoteNo := s.generateQuoteNo()

	// 創建報價單主檔
	quote := &models.Quote{
		QuoteNo:       quoteNo,
		InquiryID:     req.InquiryID,
		CustomerID:    req.CustomerID,
		Status:        "draft",
		ValidityDays:  req.ValidityDays,
		PaymentTerms:  req.PaymentTerms,
		DeliveryTerms: req.DeliveryTerms,
		Remarks:       req.Remarks,
		CreatedBy:     createdBy,
		TemplateID:    req.TemplateID,
	}

	if req.ValidityDays > 0 {
		validUntil := time.Now().AddDate(0, 0, req.ValidityDays)
		quote.ValidUntil = validUntil
	}

	if err := tx.Create(quote).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	// 創建第一個版本
	version := &models.QuoteVersion{
		QuoteID:       quote.ID,
		VersionNumber: 1,
		IsCurrent:     true,
		CreatedBy:     createdBy,
	}

	if err := tx.Create(version).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create quote version: %w", err)
	}

	// 更新報價單的當前版本
	quote.CurrentVersionID = &version.ID
	if err := tx.Save(quote).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update quote current version: %w", err)
	}

	// 創建報價項目
	totalAmount := 0.0
	for i, itemReq := range req.Items {
		item := models.QuoteItem{
			QuoteVersionID:    version.ID,
			ItemNo:            i + 1,
			ProductName:       itemReq.ProductName,
			ProductSpecs:      itemReq.ProductSpecs,
			Quantity:          itemReq.Quantity,
			Unit:              itemReq.Unit,
			UnitPrice:         itemReq.UnitPrice,
			TotalPrice:        float64(itemReq.Quantity) * itemReq.UnitPrice,
			CostCalculationID: itemReq.CostCalculationID,
			Notes:             itemReq.Notes,
		}
		totalAmount += item.TotalPrice

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create quote item: %w", err)
		}
	}

	// 更新報價單總金額
	quote.TotalAmount = totalAmount
	if err := tx.Save(quote).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update quote total amount: %w", err)
	}

	// 創建條款
	if req.UseTemplate && req.TemplateID == nil {
		// 使用預設模板條款
		templates, err := s.quoteRepo.GetDefaultTermsTemplates()
		if err == nil {
			for i, tmpl := range templates {
				term := models.QuoteTerm{
					QuoteVersionID: version.ID,
					TermType:       tmpl.TemplateType,
					TermContent:    tmpl.Content,
					SortOrder:      i + 1,
				}
				if err := tx.Create(&term).Error; err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("failed to create quote term: %w", err)
				}
			}
		}
	} else if len(req.Terms) > 0 {
		// 使用自定義條款
		for _, termReq := range req.Terms {
			term := models.QuoteTerm{
				QuoteVersionID: version.ID,
				TermType:       termReq.TermType,
				TermContent:    termReq.TermContent,
				SortOrder:      termReq.SortOrder,
			}
			if err := tx.Create(&term).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create quote term: %w", err)
			}
		}
	}

	// 記錄活動日誌
	s.logActivity(tx, quote.ID, &version.ID, "created", "Quote created", createdBy)

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 重新載入完整資料
	fullQuote, err := s.quoteRepo.GetQuoteByID(quote.ID)
	if err != nil {
		return nil, err
	}

	// Trigger webhook for quote created
	if s.webhookService != nil {
		go s.webhookService.TriggerQuoteCreated(fullQuote, fullQuote.CompanyID, createdBy)
	}

	return fullQuote, nil
}

// UpdateQuote 更新報價單
func (s *QuoteManagementService) UpdateQuote(quoteID uuid.UUID, req models.UpdateQuoteRequest, updatedBy uuid.UUID) (*models.Quote, error) {
	// 獲取現有報價單
	quote, err := s.quoteRepo.GetQuoteByID(quoteID)
	if err != nil {
		return nil, err
	}

	// 檢查狀態
	if quote.Status != "draft" && quote.Status != "rejected" {
		return nil, errors.New("only draft or rejected quotes can be updated")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var version *models.QuoteVersion

	if req.CreateNewVersion {
		// 創建新版本
		currentVersion, err := s.quoteRepo.GetCurrentVersion(quoteID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 將當前版本設為非當前
		currentVersion.IsCurrent = false
		if err := tx.Save(currentVersion).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 創建新版本
		version = &models.QuoteVersion{
			QuoteID:       quoteID,
			VersionNumber: currentVersion.VersionNumber + 1,
			VersionNotes:  req.VersionNotes,
			IsCurrent:     true,
			CreatedBy:     updatedBy,
		}

		if err := tx.Create(version).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		quote.CurrentVersionID = &version.ID
	} else {
		// 使用當前版本
		var err error
		version, err = s.quoteRepo.GetCurrentVersion(quoteID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 更新報價單基本資訊
	quote.ValidityDays = req.ValidityDays
	quote.PaymentTerms = req.PaymentTerms
	quote.DeliveryTerms = req.DeliveryTerms
	quote.Remarks = req.Remarks
	quote.UpdatedBy = &updatedBy

	if req.ValidityDays > 0 {
		validUntil := time.Now().AddDate(0, 0, req.ValidityDays)
		quote.ValidUntil = validUntil
	}

	// 如果有新項目，刪除舊項目並創建新項目
	if len(req.Items) > 0 {
		// 刪除舊項目
		if err := tx.Where("quote_version_id = ?", version.ID).Delete(&models.QuoteItem{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 創建新項目
		totalAmount := 0.0
		for i, itemReq := range req.Items {
			item := models.QuoteItem{
				QuoteVersionID:    version.ID,
				ItemNo:            i + 1,
				ProductName:       itemReq.ProductName,
				ProductSpecs:      itemReq.ProductSpecs,
				Quantity:          itemReq.Quantity,
				Unit:              itemReq.Unit,
				UnitPrice:         itemReq.UnitPrice,
				TotalPrice:        float64(itemReq.Quantity) * itemReq.UnitPrice,
				CostCalculationID: itemReq.CostCalculationID,
				Notes:             itemReq.Notes,
			}
			totalAmount += item.TotalPrice

			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		quote.TotalAmount = totalAmount
	}

	// 更新條款
	if len(req.Terms) > 0 {
		// 刪除舊條款
		if err := tx.Where("quote_version_id = ?", version.ID).Delete(&models.QuoteTerm{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 創建新條款
		for _, termReq := range req.Terms {
			term := models.QuoteTerm{
				QuoteVersionID: version.ID,
				TermType:       termReq.TermType,
				TermContent:    termReq.TermContent,
				SortOrder:      termReq.SortOrder,
			}
			if err := tx.Create(&term).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 保存報價單
	if err := tx.Save(quote).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 記錄活動日誌
	s.logActivity(tx, quoteID, &version.ID, "updated", "Quote updated", updatedBy)

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.quoteRepo.GetQuoteByID(quoteID)
}

// SubmitForApproval 提交審核
func (s *QuoteManagementService) SubmitForApproval(quoteID uuid.UUID, req models.SubmitApprovalRequest, submittedBy uuid.UUID) error {
	quote, err := s.quoteRepo.GetQuoteByID(quoteID)
	if err != nil {
		return err
	}

	if quote.Status != "draft" && quote.Status != "rejected" {
		return errors.New("only draft or rejected quotes can be submitted for approval")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新狀態
	quote.Status = "pending_approval"
	quote.UpdatedBy = &submittedBy
	if err := tx.Save(quote).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 創建審核記錄
	approvals := s.determineApprovalLevels(quote.TotalAmount)
	for _, approval := range approvals {
		approval.QuoteID = quoteID
		approval.QuoteVersionID = *quote.CurrentVersionID
		if err := tx.Create(&approval).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 記錄活動日誌
	s.logActivity(tx, quoteID, quote.CurrentVersionID, "submitted", req.Notes, submittedBy)

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Trigger webhook for quote submitted for approval
	if s.webhookService != nil {
		go s.webhookService.TriggerQuoteSubmittedForApproval(quoteID, quote.CompanyID, submittedBy)
	}

	return nil
}

// ApproveQuote 審核報價單
func (s *QuoteManagementService) ApproveQuote(quoteID uuid.UUID, req models.ApproveQuoteRequest, approverID uuid.UUID) error {
	quote, err := s.quoteRepo.GetQuoteByID(quoteID)
	if err != nil {
		return err
	}

	if quote.Status != "pending_approval" {
		return errors.New("quote is not pending approval")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 獲取當前待審核記錄
	approval, err := s.quoteRepo.GetPendingApproval(quoteID, approverID)
	if err != nil {
		tx.Rollback()
		return errors.New("no pending approval found for this user")
	}

	// 更新審核記錄
	now := time.Now()
	approval.ActualApproverID = &approverID
	approval.ApprovalNotes = req.Notes
	approval.ApprovedAt = &now

	if req.Approved {
		approval.ApprovalStatus = "approved"
	} else {
		approval.ApprovalStatus = "rejected"
	}

	if err := tx.Save(approval).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 檢查是否所有審核都完成
	allApproved, anyRejected, err := s.checkApprovalStatus(tx, quoteID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新報價單狀態
	if anyRejected {
		quote.Status = "rejected"
		quote.ApprovalStatus = "rejected"
	} else if allApproved {
		quote.Status = "approved"
		quote.ApprovalStatus = "approved"
		quote.ApprovedAt = &now
		quote.ApprovedBy = &approverID
		quote.ApprovedAmount = quote.TotalAmount
	}

	quote.UpdatedBy = &approverID
	if err := tx.Save(quote).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 記錄活動日誌
	activityType := "approved"
	if !req.Approved {
		activityType = "rejected"
	}
	s.logActivity(tx, quoteID, quote.CurrentVersionID, activityType, req.Notes, approverID)

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Trigger webhook based on approval result
	if s.webhookService != nil {
		if req.Approved {
			go s.webhookService.TriggerQuoteApproved(quoteID, approverID, quote.CompanyID, approverID)
		} else {
			go s.webhookService.TriggerQuoteRejected(quoteID, req.Notes, quote.CompanyID, approverID)
		}
	}

	return nil
}

// SendQuote 發送報價單
func (s *QuoteManagementService) SendQuote(quoteID uuid.UUID, req models.SendQuoteRequest, sentBy uuid.UUID) error {
	quote, err := s.quoteRepo.GetQuoteByID(quoteID)
	if err != nil {
		return err
	}

	if quote.Status != "approved" {
		return errors.New("only approved quotes can be sent")
	}

	// 生成 PDF
	var pdfPath string
	if req.AttachPDF {
		pdfPath, err = s.pdfService.GenerateQuotePDF(quote)
		if err != nil {
			return fmt.Errorf("failed to generate PDF: %w", err)
		}
	}

	// 創建發送記錄
	sendLog := &models.QuoteSendLog{
		QuoteID:        quoteID,
		QuoteVersionID: *quote.CurrentVersionID,
		SendMethod:     "email",
		RecipientEmail: req.RecipientEmail,
		RecipientName:  req.RecipientName,
		Subject:        req.Subject,
		Message:        req.Message,
		SendStatus:     "pending",
		CreatedBy:      sentBy,
	}

	if err := s.db.Create(sendLog).Error; err != nil {
		return err
	}

	// 發送郵件
	err = s.emailService.SendQuote(req, pdfPath, quote)
	now := time.Now()
	if err != nil {
		sendLog.SendStatus = "failed"
		sendLog.ErrorMessage = err.Error()
	} else {
		sendLog.SendStatus = "sent"
		sendLog.SentAt = &now
	}

	// 更新發送記錄
	if err := s.db.Save(sendLog).Error; err != nil {
		return err
	}

	if sendLog.SendStatus == "sent" {
		// 記錄活動日誌
		s.logActivity(s.db, quoteID, quote.CurrentVersionID, "sent", 
			fmt.Sprintf("Sent to %s", req.RecipientEmail), sentBy)
		
		// Trigger webhook for quote sent
		if s.webhookService != nil {
			go s.webhookService.TriggerQuoteSent(quoteID, req.RecipientEmail, quote.CompanyID, sentBy)
		}
	}

	return nil
}

// GetQuoteVersions 獲取報價單版本歷史
func (s *QuoteManagementService) GetQuoteVersions(quoteID uuid.UUID) ([]models.QuoteVersion, error) {
	return s.quoteRepo.GetQuoteVersions(quoteID)
}

// GetQuoteActivityLogs 獲取報價單活動日誌
func (s *QuoteManagementService) GetQuoteActivityLogs(quoteID uuid.UUID) ([]models.QuoteActivityLog, error) {
	return s.quoteRepo.GetActivityLogs(quoteID)
}

// Private methods

func (s *QuoteManagementService) generateQuoteNo() string {
	// 格式: Q-YYYYMMDD-XXXX
	date := time.Now().Format("20060102")
	
	var count int64
	s.db.Model(&models.Quote{}).
		Where("quote_no LIKE ?", fmt.Sprintf("Q-%s-%%", date)).
		Count(&count)
	
	return fmt.Sprintf("Q-%s-%04d", date, count+1)
}

func (s *QuoteManagementService) determineApprovalLevels(totalAmount float64) []models.QuoteApproval {
	var approvals []models.QuoteApproval

	// 根據金額決定審核層級
	if totalAmount < 10000 {
		// 只需工程主管審核
		approvals = append(approvals, models.QuoteApproval{
			ApprovalLevel:  1,
			ApproverRole:   "engineer_lead",
			ApprovalStatus: "pending",
		})
	} else if totalAmount < 50000 {
		// 需要工程主管和業務經理審核
		approvals = append(approvals, 
			models.QuoteApproval{
				ApprovalLevel:  1,
				ApproverRole:   "engineer_lead",
				ApprovalStatus: "pending",
			},
			models.QuoteApproval{
				ApprovalLevel:  2,
				ApproverRole:   "sales_manager",
				ApprovalStatus: "pending",
			},
		)
	} else {
		// 需要三級審核
		approvals = append(approvals,
			models.QuoteApproval{
				ApprovalLevel:  1,
				ApproverRole:   "engineer_lead",
				ApprovalStatus: "pending",
			},
			models.QuoteApproval{
				ApprovalLevel:  2,
				ApproverRole:   "sales_manager",
				ApprovalStatus: "pending",
			},
			models.QuoteApproval{
				ApprovalLevel:  3,
				ApproverRole:   "general_manager",
				ApprovalStatus: "pending",
			},
		)
	}

	return approvals
}

func (s *QuoteManagementService) checkApprovalStatus(tx *gorm.DB, quoteID uuid.UUID) (allApproved, anyRejected bool, err error) {
	var approvals []models.QuoteApproval
	if err = tx.Where("quote_id = ?", quoteID).Find(&approvals).Error; err != nil {
		return false, false, err
	}

	allApproved = true
	for _, approval := range approvals {
		if approval.ApprovalStatus == "rejected" {
			anyRejected = true
			allApproved = false
			break
		}
		if approval.ApprovalStatus == "pending" {
			allApproved = false
		}
	}

	return allApproved, anyRejected, nil
}

func (s *QuoteManagementService) logActivity(tx *gorm.DB, quoteID uuid.UUID, versionID *uuid.UUID, activityType, description string, performedBy uuid.UUID) {
	log := &models.QuoteActivityLog{
		QuoteID:             quoteID,
		QuoteVersionID:      versionID,
		ActivityType:        activityType,
		ActivityDescription: description,
		PerformedBy:         performedBy,
	}
	tx.Create(log)
}
package repositories

import (
	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuoteManagementRepository struct {
	db *gorm.DB
}

func NewQuoteManagementRepository(db *gorm.DB) *QuoteManagementRepository {
	return &QuoteManagementRepository{db: db}
}

// GetQuoteByID 根據ID獲取報價單（包含關聯資料）
func (r *QuoteManagementRepository) GetQuoteByID(id uuid.UUID) (*models.Quote, error) {
	var quote models.Quote
	err := r.db.Preload("Customer").
		Preload("Inquiry").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ApprovedByUser").
		Preload("Template").
		First(&quote, "id = ?", id).Error
	return &quote, err
}

// GetQuotes 獲取報價單列表（分頁）
func (r *QuoteManagementRepository) GetQuotes(companyID uuid.UUID, page, pageSize int, status string) ([]models.Quote, int64, error) {
	var quotes []models.Quote
	var total int64

	query := r.db.Model(&models.Quote{}).Where("company_id = ?", companyID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 計算總數
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 獲取分頁資料
	offset := (page - 1) * pageSize
	err := query.Preload("Customer").
		Preload("Inquiry").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&quotes).Error

	return quotes, total, err
}

// GetCurrentVersion 獲取報價單當前版本
func (r *QuoteManagementRepository) GetCurrentVersion(quoteID uuid.UUID) (*models.QuoteVersion, error) {
	var version models.QuoteVersion
	err := r.db.Preload("Items").
		Preload("Items.CostCalculation").
		Preload("Terms").
		Where("quote_id = ? AND is_current = ?", quoteID, true).
		First(&version).Error
	return &version, err
}

// GetQuoteVersions 獲取報價單所有版本
func (r *QuoteManagementRepository) GetQuoteVersions(quoteID uuid.UUID) ([]models.QuoteVersion, error) {
	var versions []models.QuoteVersion
	err := r.db.Where("quote_id = ?", quoteID).
		Preload("Creator").
		Order("version_number DESC").
		Find(&versions).Error
	return versions, err
}

// GetQuoteVersion 根據ID獲取特定版本
func (r *QuoteManagementRepository) GetQuoteVersion(versionID uuid.UUID) (*models.QuoteVersion, error) {
	var version models.QuoteVersion
	err := r.db.Preload("Items").
		Preload("Items.CostCalculation").
		Preload("Terms").
		Preload("Creator").
		First(&version, "id = ?", versionID).Error
	return &version, err
}

// GetPendingApproval 獲取待審核記錄
func (r *QuoteManagementRepository) GetPendingApproval(quoteID uuid.UUID, approverID uuid.UUID) (*models.QuoteApproval, error) {
	var approval models.QuoteApproval
	
	// 根據用戶角色查找待審核記錄
	err := r.db.Joins("JOIN accounts ON accounts.id = ?", approverID).
		Where("quote_approvals.quote_id = ? AND quote_approvals.approval_status = 'pending'", quoteID).
		Where("(quote_approvals.required_approver_id = ? OR " +
			  "(quote_approvals.required_approver_id IS NULL AND " +
			  "quote_approvals.approver_role = accounts.role))", approverID).
		First(&approval).Error
		
	return &approval, err
}

// GetQuoteApprovals 獲取報價單審核記錄
func (r *QuoteManagementRepository) GetQuoteApprovals(quoteID uuid.UUID) ([]models.QuoteApproval, error) {
	var approvals []models.QuoteApproval
	err := r.db.Where("quote_id = ?", quoteID).
		Preload("RequiredApprover").
		Preload("ActualApprover").
		Order("approval_level ASC").
		Find(&approvals).Error
	return approvals, err
}

// GetActivityLogs 獲取活動日誌
func (r *QuoteManagementRepository) GetActivityLogs(quoteID uuid.UUID) ([]models.QuoteActivityLog, error) {
	var logs []models.QuoteActivityLog
	err := r.db.Where("quote_id = ?", quoteID).
		Preload("Performer").
		Order("performed_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetSendLogs 獲取發送記錄
func (r *QuoteManagementRepository) GetSendLogs(quoteID uuid.UUID) ([]models.QuoteSendLog, error) {
	var logs []models.QuoteSendLog
	err := r.db.Where("quote_id = ?", quoteID).
		Preload("Creator").
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetDefaultTermsTemplates 獲取預設條款模板
func (r *QuoteManagementRepository) GetDefaultTermsTemplates() ([]models.QuoteTermsTemplate, error) {
	var templates []models.QuoteTermsTemplate
	err := r.db.Where("is_default = ? AND is_active = ?", true, true).
		Order("template_type").
		Find(&templates).Error
	return templates, err
}

// GetTermsTemplates 獲取所有條款模板
func (r *QuoteManagementRepository) GetTermsTemplates(templateType string) ([]models.QuoteTermsTemplate, error) {
	var templates []models.QuoteTermsTemplate
	query := r.db.Where("is_active = ?", true)
	
	if templateType != "" {
		query = query.Where("template_type = ?", templateType)
	}
	
	err := query.Order("template_type, template_name").Find(&templates).Error
	return templates, err
}

// GetQuoteTemplates 獲取報價單模板
func (r *QuoteManagementRepository) GetQuoteTemplates() ([]models.QuoteTemplate, error) {
	var templates []models.QuoteTemplate
	err := r.db.Where("is_active = ?", true).
		Order("is_default DESC, template_name").
		Find(&templates).Error
	return templates, err
}

// GetDefaultQuoteTemplate 獲取預設報價單模板
func (r *QuoteManagementRepository) GetDefaultQuoteTemplate() (*models.QuoteTemplate, error) {
	var template models.QuoteTemplate
	err := r.db.Where("is_default = ? AND is_active = ?", true, true).
		First(&template).Error
	return &template, err
}

// CreateTermsTemplate 創建條款模板
func (r *QuoteManagementRepository) CreateTermsTemplate(template *models.QuoteTermsTemplate) error {
	return r.db.Create(template).Error
}

// UpdateTermsTemplate 更新條款模板
func (r *QuoteManagementRepository) UpdateTermsTemplate(id uuid.UUID, template models.QuoteTermsTemplate) error {
	return r.db.Model(&models.QuoteTermsTemplate{}).
		Where("id = ?", id).
		Updates(template).Error
}

// CreateQuoteTemplate 創建報價單模板
func (r *QuoteManagementRepository) CreateQuoteTemplate(template *models.QuoteTemplate) error {
	// 如果設為預設，先將其他模板設為非預設
	if template.IsDefault {
		r.db.Model(&models.QuoteTemplate{}).
			Where("is_default = ?", true).
			Update("is_default", false)
	}
	return r.db.Create(template).Error
}

// UpdateQuoteTemplate 更新報價單模板
func (r *QuoteManagementRepository) UpdateQuoteTemplate(id uuid.UUID, template models.QuoteTemplate) error {
	// 如果設為預設，先將其他模板設為非預設
	if template.IsDefault {
		r.db.Model(&models.QuoteTemplate{}).
			Where("is_default = ? AND id != ?", true, id).
			Update("is_default", false)
	}
	return r.db.Model(&models.QuoteTemplate{}).
		Where("id = ?", id).
		Updates(template).Error
}

// GetQuotesByCustomer 根據客戶獲取報價單
func (r *QuoteManagementRepository) GetQuotesByCustomer(customerID uuid.UUID, limit int) ([]models.Quote, error) {
	var quotes []models.Quote
	query := r.db.Where("customer_id = ?", customerID).
		Order("created_at DESC")
		
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&quotes).Error
	return quotes, err
}

// GetQuotesByInquiry 根據詢價單獲取報價單
func (r *QuoteManagementRepository) GetQuotesByInquiry(inquiryID uuid.UUID) ([]models.Quote, error) {
	var quotes []models.Quote
	err := r.db.Where("inquiry_id = ?", inquiryID).
		Preload("CreatedByUser").
		Order("created_at DESC").
		Find(&quotes).Error
	return quotes, err
}

// SearchQuotes 搜尋報價單
func (r *QuoteManagementRepository) SearchQuotes(companyID uuid.UUID, keyword string) ([]models.Quote, error) {
	var quotes []models.Quote
	err := r.db.Where("company_id = ?", companyID).
		Where("quote_no LIKE ? OR remarks LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Preload("Customer").
		Limit(20).
		Find(&quotes).Error
	return quotes, err
}
package handler

import (
	"net/http"
	"strconv"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuoteManagementHandler struct {
	quoteService *services.QuoteManagementService
}

func NewQuoteManagementHandler(quoteService *services.QuoteManagementService) *QuoteManagementHandler {
	return &QuoteManagementHandler{
		quoteService: quoteService,
	}
}

// CreateQuote 創建報價單
// @Summary 創建報價單
// @Description 創建新的報價單
// @Tags Quote Management
// @Accept json
// @Produce json
// @Param request body models.CreateQuoteRequest true "創建報價單請求"
// @Success 201 {object} models.Quote
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes [post]
func (h *QuoteManagementHandler) CreateQuote(c echo.Context) error {
	var req models.CreateQuoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID := c.Get("user_id").(uuid.UUID)

	quote, err := h.quoteService.CreateQuote(req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, quote)
}

// UpdateQuote 更新報價單
// @Summary 更新報價單
// @Description 更新現有報價單
// @Tags Quote Management
// @Accept json
// @Produce json
// @Param id path string true "報價單ID"
// @Param request body models.UpdateQuoteRequest true "更新報價單請求"
// @Success 200 {object} models.Quote
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id} [put]
func (h *QuoteManagementHandler) UpdateQuote(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quote ID"})
	}

	var req models.UpdateQuoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID := c.Get("user_id").(uuid.UUID)

	quote, err := h.quoteService.UpdateQuote(id, req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, quote)
}

// GetQuote 獲取報價單詳情
// @Summary 獲取報價單詳情
// @Description 獲取指定報價單的詳細資訊
// @Tags Quote Management
// @Produce json
// @Param id path string true "報價單ID"
// @Success 200 {object} models.Quote
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/quotes/{id} [get]
func (h *QuoteManagementHandler) GetQuote(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quote ID"})
	}

	quote, err := h.quoteService.GetQuote(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Quote not found"})
	}

	return c.JSON(http.StatusOK, quote)
}

// GetQuotes 獲取報價單列表
// @Summary 獲取報價單列表
// @Description 獲取報價單列表（分頁）
// @Tags Quote Management
// @Produce json
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(20)
// @Param status query string false "狀態篩選"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/quotes [get]
func (h *QuoteManagementHandler) GetQuotes(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}

	status := c.QueryParam("status")
	companyID := c.Get("company_id").(uuid.UUID)

	quotes, total, err := h.quoteService.GetQuotes(companyID, page, pageSize, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": quotes,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// SubmitForApproval 提交審核
// @Summary 提交報價單審核
// @Description 將報價單提交審核流程
// @Tags Quote Management
// @Accept json
// @Produce json
// @Param id path string true "報價單ID"
// @Param request body models.SubmitApprovalRequest true "提交審核請求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id}/submit [post]
func (h *QuoteManagementHandler) SubmitForApproval(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quote ID"})
	}

	var req models.SubmitApprovalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID := c.Get("user_id").(uuid.UUID)

	if err := h.quoteService.SubmitForApproval(id, req, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Quote submitted for approval"})
}

// ApproveQuote 審核報價單
// @Summary 審核報價單
// @Description 審核通過或拒絕報價單
// @Tags Quote Management
// @Accept json
// @Produce json
// @Param id path string true "報價單ID"
// @Param request body models.ApproveQuoteRequest true "審核請求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id}/approve [post]
func (h *QuoteManagementHandler) ApproveQuote(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quote ID"})
	}

	var req models.ApproveQuoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID := c.Get("user_id").(uuid.UUID)

	if err := h.quoteService.ApproveQuote(id, req, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	message := "Quote approved"
	if !req.Approved {
		message = "Quote rejected"
	}

	return c.JSON(http.StatusOK, map[string]string{"message": message})
}

// SendQuote 發送報價單
// @Summary 發送報價單
// @Description 通過郵件發送報價單
// @Tags Quote Management
// @Accept json
// @Produce json
// @Param id path string true "報價單ID"
// @Param request body models.SendQuoteRequest true "發送請求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id}/send [post]
func (h *QuoteManagementHandler) SendQuote(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quote ID"})
	}

	var req models.SendQuoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID := c.Get("user_id").(uuid.UUID)

	if err := h.quoteService.SendQuote(id, req, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Quote sent successfully"})
}

// GetQuoteVersion 獲取報價單版本
// @Summary 獲取報價單版本詳情
// @Description 獲取指定版本的報價單詳情
// @Tags Quote Management
// @Produce json
// @Param id path string true "報價單ID"
// @Param version_id path string true "版本ID"
// @Success 200 {object} models.QuoteVersion
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/quotes/{id}/versions/{version_id} [get]
func (h *QuoteManagementHandler) GetQuoteVersion(c echo.Context) error {
	versionIDStr := c.Param("version_id")
	versionID, err := uuid.Parse(versionIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid version ID"})
	}

	version, err := h.quoteService.GetQuoteVersion(versionID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Version not found"})
	}

	return c.JSON(http.StatusOK, version)
}

// GetQuoteVersions 獲取報價單版本歷史
// @Summary 獲取報價單版本歷史
// @Description 獲取報價單的所有版本列表
// @Tags Quote Management
// @Produce json
// @Param id path string true "報價單ID"
// @Success 200 {array} models.QuoteVersion
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id}/versions [get]
func (h *QuoteManagementHandler) GetQuoteVersions(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quote ID"})
	}

	versions, err := h.quoteService.GetQuoteVersions(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, versions)
}

// GetQuoteActivityLogs 獲取活動日誌
// @Summary 獲取報價單活動日誌
// @Description 獲取報價單的所有活動記錄
// @Tags Quote Management
// @Produce json
// @Param id path string true "報價單ID"
// @Success 200 {array} models.QuoteActivityLog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/quotes/{id}/activities [get]
func (h *QuoteManagementHandler) GetQuoteActivityLogs(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid quote ID"})
	}

	logs, err := h.quoteService.GetQuoteActivityLogs(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, logs)
}

// GetTermsTemplates 獲取條款模板
// @Summary 獲取條款模板列表
// @Description 獲取可用的條款模板
// @Tags Quote Management
// @Produce json
// @Param type query string false "模板類型"
// @Success 200 {array} models.QuoteTermsTemplate
// @Failure 500 {object} map[string]string
// @Router /api/quotes/terms-templates [get]
func (h *QuoteManagementHandler) GetTermsTemplates(c echo.Context) error {
	templateType := c.QueryParam("type")
	
	templates, err := h.quoteService.GetTermsTemplates(templateType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, templates)
}

// GetQuoteTemplates 獲取報價單模板
// @Summary 獲取報價單模板列表
// @Description 獲取可用的報價單模板
// @Tags Quote Management
// @Produce json
// @Success 200 {array} models.QuoteTemplate
// @Failure 500 {object} map[string]string
// @Router /api/quotes/templates [get]
func (h *QuoteManagementHandler) GetQuoteTemplates(c echo.Context) error {
	templates, err := h.quoteService.GetQuoteTemplates()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, templates)
}

// RegisterRoutes 註冊路由
func (h *QuoteManagementHandler) RegisterRoutes(e *echo.Group) {
	quotes := e.Group("/quotes")
	
	// 基本 CRUD
	quotes.POST("", h.CreateQuote)
	quotes.GET("", h.GetQuotes)
	quotes.GET("/:id", h.GetQuote)
	quotes.PUT("/:id", h.UpdateQuote)
	
	// 審核流程
	quotes.POST("/:id/submit", h.SubmitForApproval)
	quotes.POST("/:id/approve", h.ApproveQuote)
	
	// 發送
	quotes.POST("/:id/send", h.SendQuote)
	
	// 版本管理
	quotes.GET("/:id/versions", h.GetQuoteVersions)
	quotes.GET("/:id/versions/:version_id", h.GetQuoteVersion)
	
	// 活動日誌
	quotes.GET("/:id/activities", h.GetQuoteActivityLogs)
	
	// 模板
	quotes.GET("/terms-templates", h.GetTermsTemplates)
	quotes.GET("/templates", h.GetQuoteTemplates)
}
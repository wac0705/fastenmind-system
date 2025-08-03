package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/fastenmind/fastener-api/internal/model"
	"github.com/fastenmind/fastener-api/internal/service"
	"gorm.io/gorm"
)

// CustomerHandler handles customer-related HTTP requests
type CustomerHandler struct {
	BaseHandler
	customerService service.CustomerService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		BaseHandler:     *NewBaseHandler(),
		customerService: customerService,
	}
}

// List handles GET /api/v1/customers
func (h *CustomerHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}
	
	filter := model.CustomerFilter{
		CompanyID: c.Get("company_id").(uuid.UUID),
		Search:    c.QueryParam("search"),
		Country:   c.QueryParam("country"),
		IsActive:  c.QueryParam("is_active") == "true",
		Page:      page,
		PageSize:  pageSize,
	}
	
	customers, pagination, err := h.customerService.List(ctx, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data":       customers,
		"pagination": pagination,
	})
}

// Create handles POST /api/v1/customers
func (h *CustomerHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req model.CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	req.CompanyID = c.Get("company_id").(uuid.UUID)
	req.CreatedBy = c.Get("user_id").(uuid.UUID)
	
	customer, err := h.customerService.Create(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusCreated, customer)
}

// Get handles GET /api/v1/customers/:id
func (h *CustomerHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer ID")
	}
	
	companyID := c.Get("company_id").(uuid.UUID)
	
	customer, err := h.customerService.GetByID(ctx, customerID, companyID)
	if err != nil {
		if err == service.ErrCustomerNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, customer)
}

// Update handles PUT /api/v1/customers/:id
func (h *CustomerHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer ID")
	}
	
	var req model.UpdateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	req.ID = customerID
	req.CompanyID = c.Get("company_id").(uuid.UUID)
	req.UpdatedBy = c.Get("user_id").(uuid.UUID)
	
	customer, err := h.customerService.Update(ctx, req)
	if err != nil {
		if err == service.ErrCustomerNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, customer)
}

// Delete handles DELETE /api/v1/customers/:id
func (h *CustomerHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer ID")
	}
	
	companyID := c.Get("company_id").(uuid.UUID)
	
	err = h.customerService.Delete(ctx, customerID, companyID)
	if err != nil {
		if err == service.ErrCustomerNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Customer not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.NoContent(http.StatusNoContent)
}

// GetStatistics handles GET /api/v1/customers/:id/statistics
func (h *CustomerHandler) GetStatistics(c echo.Context) error {
	ctx := c.Request().Context()
	
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer ID")
	}
	
	companyID := c.Get("company_id").(uuid.UUID)
	
	stats, err := h.customerService.GetStatistics(ctx, customerID, companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, stats)
}

// Export handles GET /api/v1/customers/export
func (h *CustomerHandler) Export(c echo.Context) error {
	// 使用帶超時的上下文
	ctx, cancel := WithContext(c, 60*time.Second)
	defer cancel()
	
	filter := model.CustomerFilter{
		CompanyID: c.Get("company_id").(uuid.UUID),
		Search:    c.QueryParam("search"),
		Country:   c.QueryParam("country"),
	}
	
	format := c.QueryParam("format")
	if format == "" {
		format = "csv"
	}
	
	// 使用臨時文件處理大數據導出
	return h.WithTempFile(c, "customer-export-*."+format, func(tempFile *os.File) error {
		// 在事務中執行導出
		return h.WithTransaction(c, func(tx *gorm.DB) error {
			data, contentType, err := h.customerService.Export(ctx, filter, format)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			
			// 寫入臨時文件
			if _, err := tempFile.Write(data); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to write export data")
			}
			
			// 重置文件指針
			if _, err := tempFile.Seek(0, 0); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to seek file")
			}
			
			// 設置響應頭
			c.Response().Header().Set("Content-Type", contentType)
			c.Response().Header().Set("Content-Disposition", "attachment; filename=customers."+format)
			c.Response().Header().Set("Content-Length", strconv.Itoa(len(data)))
			
			// 流式傳輸文件內容
			return c.Stream(http.StatusOK, contentType, tempFile)
		})
	})
}

// GetCreditHistory handles GET /api/v1/customers/:id/credit-history
func (h *CustomerHandler) GetCreditHistory(c echo.Context) error {
	ctx := c.Request().Context()
	
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid customer ID")
	}
	
	companyID := c.Get("company_id").(uuid.UUID)
	
	history, err := h.customerService.GetCreditHistory(ctx, customerID, companyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": history,
	})
}
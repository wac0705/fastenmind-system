package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fastenmind/fastener-api/internal/api/rest"
	"github.com/fastenmind/fastener-api/internal/domain/models"
	"github.com/fastenmind/fastener-api/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInquiryService is a mock implementation of InquiryService
type MockInquiryService struct {
	mock.Mock
}

func (m *MockInquiryService) Create(ctx context.Context, req service.CreateInquiryRequest) (*models.Inquiry, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inquiry), args.Error(1)
}

func (m *MockInquiryService) GetByID(ctx context.Context, id uuid.UUID) (*models.Inquiry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inquiry), args.Error(1)
}

func (m *MockInquiryService) List(ctx context.Context, params service.ListInquiriesParams) (*service.PaginatedResult[*models.Inquiry], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.PaginatedResult[*models.Inquiry]), args.Error(1)
}

func (m *MockInquiryService) AssignEngineer(ctx context.Context, inquiryID, engineerID, assignedBy uuid.UUID, notes string) (*models.Inquiry, error) {
	args := m.Called(ctx, inquiryID, engineerID, assignedBy, notes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inquiry), args.Error(1)
}

func (m *MockInquiryService) UpdateStatus(ctx context.Context, inquiryID uuid.UUID, status models.InquiryStatus, updatedBy uuid.UUID) (*models.Inquiry, error) {
	args := m.Called(ctx, inquiryID, status, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inquiry), args.Error(1)
}

// Test suite setup
func setupInquiryHandlerTest() (*echo.Echo, *rest.InquiryHandler, *MockInquiryService) {
	e := echo.New()
	mockService := new(MockInquiryService)
	handler := rest.NewInquiryHandler(mockService)
	
	return e, handler, mockService
}

func TestInquiryHandler_Create(t *testing.T) {
	e, handler, mockService := setupInquiryHandlerTest()
	
	t.Run("successful inquiry creation", func(t *testing.T) {
		// Arrange
		reqBody := rest.CreateInquiryRequest{
			CompanyID:       uuid.New().String(),
			CustomerID:      uuid.New().String(),
			SalesID:         uuid.New().String(),
			ProductCategory: "Fasteners",
			ProductName:     "Hex Bolt",
			Quantity:        1000,
			Unit:            "PCS",
			RequiredDate:    time.Now().Add(30 * 24 * time.Hour),
			Incoterm:        "FOB",
			PaymentTerms:    "Net 30",
		}
		
		inquiry := &models.Inquiry{
			BaseModel:       models.BaseModel{ID: uuid.New()},
			InquiryNo:       "INQ-2024-001",
			CompanyID:       uuid.MustParse(reqBody.CompanyID),
			CustomerID:      uuid.MustParse(reqBody.CustomerID),
			SalesID:         uuid.MustParse(reqBody.SalesID),
			ProductCategory: reqBody.ProductCategory,
			ProductName:     reqBody.ProductName,
			Quantity:        reqBody.Quantity,
			Unit:            reqBody.Unit,
			Status:          models.InquiryStatusPending,
		}
		
		mockService.On("Create", mock.Anything, mock.AnythingOfType("service.CreateInquiryRequest")).Return(inquiry, nil)
		
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Act
		err := handler.Create(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		var response rest.InquiryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, inquiry.InquiryNo, response.InquiryNo)
		assert.Equal(t, string(inquiry.Status), response.Status)
		
		mockService.AssertExpectations(t)
	})
	
	t.Run("invalid request body", func(t *testing.T) {
		// Arrange
		reqBody := rest.CreateInquiryRequest{
			// Missing required fields
			ProductName: "Hex Bolt",
		}
		
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Act
		err := handler.Create(c)
		
		// Assert
		assert.NoError(t, err) // Handler returns HTTP error, not Go error
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		
		mockService.AssertNotCalled(t, "Create")
	})
}

func TestInquiryHandler_GetByID(t *testing.T) {
	e, handler, mockService := setupInquiryHandlerTest()
	
	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		inquiry := &models.Inquiry{
			BaseModel:   models.BaseModel{ID: inquiryID},
			InquiryNo:   "INQ-2024-001",
			ProductName: "Hex Bolt",
			Quantity:    1000,
			Status:      models.InquiryStatusPending,
		}
		
		mockService.On("GetByID", mock.Anything, inquiryID).Return(inquiry, nil)
		
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inquiries/"+inquiryID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(inquiryID.String())
		
		// Act
		err := handler.GetByID(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		var response rest.InquiryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, inquiry.InquiryNo, response.InquiryNo)
		
		mockService.AssertExpectations(t)
	})
	
	t.Run("inquiry not found", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		mockService.On("GetByID", mock.Anything, inquiryID).Return(nil, service.ErrNotFound)
		
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inquiries/"+inquiryID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(inquiryID.String())
		
		// Act
		err := handler.GetByID(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		
		mockService.AssertExpectations(t)
	})
	
	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inquiries/invalid-uuid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid-uuid")
		
		// Act
		err := handler.GetByID(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		
		mockService.AssertNotCalled(t, "GetByID")
	})
}

func TestInquiryHandler_List(t *testing.T) {
	e, handler, mockService := setupInquiryHandlerTest()
	
	t.Run("successful list with filters", func(t *testing.T) {
		// Arrange
		companyID := uuid.New()
		inquiries := []*models.Inquiry{
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				InquiryNo:   "INQ-2024-001",
				ProductName: "Hex Bolt",
				Status:      models.InquiryStatusPending,
			},
			{
				BaseModel:   models.BaseModel{ID: uuid.New()},
				InquiryNo:   "INQ-2024-002",
				ProductName: "Flat Washer",
				Status:      models.InquiryStatusPending,
			},
		}
		
		result := &service.PaginatedResult[*models.Inquiry]{
			Items:      inquiries,
			Total:      2,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}
		
		expectedParams := service.ListInquiriesParams{
			CompanyID: companyID,
			Status:    "pending",
			Page:      1,
			PageSize:  10,
		}
		
		mockService.On("List", mock.Anything, mock.MatchedBy(func(params service.ListInquiriesParams) bool {
			return params.CompanyID == expectedParams.CompanyID &&
				params.Status == expectedParams.Status &&
				params.Page == expectedParams.Page &&
				params.PageSize == expectedParams.PageSize
		})).Return(result, nil)
		
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inquiries?company_id="+companyID.String()+"&status=pending&page=1&page_size=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Act
		err := handler.List(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		var response rest.PaginatedInquiryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 2)
		assert.Equal(t, int64(2), response.Total)
		
		mockService.AssertExpectations(t)
	})
}

func TestInquiryHandler_AssignEngineer(t *testing.T) {
	e, handler, mockService := setupInquiryHandlerTest()
	
	t.Run("successful engineer assignment", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		engineerID := uuid.New()
		assignedBy := uuid.New()
		
		reqBody := rest.AssignEngineerRequest{
			EngineerID: engineerID.String(),
			AssignedBy: assignedBy.String(),
			Notes:      "High priority customer",
		}
		
		inquiry := &models.Inquiry{
			BaseModel:          models.BaseModel{ID: inquiryID},
			InquiryNo:          "INQ-2024-001",
			Status:             models.InquiryStatusAssigned,
			AssignedEngineerID: &engineerID,
		}
		
		mockService.On("AssignEngineer", mock.Anything, inquiryID, engineerID, assignedBy, reqBody.Notes).Return(inquiry, nil)
		
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inquiries/"+inquiryID.String()+"/assign", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(inquiryID.String())
		
		// Act
		err := handler.AssignEngineer(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		var response rest.InquiryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, string(models.InquiryStatusAssigned), response.Status)
		assert.NotNil(t, response.AssignedEngineerID)
		
		mockService.AssertExpectations(t)
	})
}

func TestInquiryHandler_UpdateStatus(t *testing.T) {
	e, handler, mockService := setupInquiryHandlerTest()
	
	t.Run("successful status update", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		updatedBy := uuid.New()
		
		reqBody := rest.UpdateStatusRequest{
			Status:    string(models.InquiryStatusQuoted),
			UpdatedBy: updatedBy.String(),
		}
		
		inquiry := &models.Inquiry{
			BaseModel: models.BaseModel{ID: inquiryID},
			InquiryNo: "INQ-2024-001",
			Status:    models.InquiryStatusQuoted,
		}
		
		mockService.On("UpdateStatus", mock.Anything, inquiryID, models.InquiryStatusQuoted, updatedBy).Return(inquiry, nil)
		
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/inquiries/"+inquiryID.String()+"/status", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(inquiryID.String())
		
		// Act
		err := handler.UpdateStatus(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		var response rest.InquiryResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, string(models.InquiryStatusQuoted), response.Status)
		
		mockService.AssertExpectations(t)
	})
	
	t.Run("invalid status value", func(t *testing.T) {
		// Arrange
		inquiryID := uuid.New()
		reqBody := rest.UpdateStatusRequest{
			Status:    "invalid_status",
			UpdatedBy: uuid.New().String(),
		}
		
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/inquiries/"+inquiryID.String()+"/status", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(inquiryID.String())
		
		// Act
		err := handler.UpdateStatus(c)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		
		mockService.AssertNotCalled(t, "UpdateStatus")
	})
}
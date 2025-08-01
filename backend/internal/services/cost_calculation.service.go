package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/fastenmind/fastener-api/internal/models"
	"github.com/fastenmind/fastener-api/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CostCalculationService struct {
	db       *gorm.DB
	costRepo *repositories.CostCalculationRepository
}

func NewCostCalculationService(db *gorm.DB) *CostCalculationService {
	return &CostCalculationService{
		db:       db,
		costRepo: repositories.NewCostCalculationRepository(db),
	}
}

// CalculateCost 計算產品成本
func (s *CostCalculationService) CalculateCost(req models.CostCalculationRequest) (*models.CostCalculation, error) {
	// 生成計算編號
	calcNo := s.generateCalculationNo()

	// 獲取或創建製程路線
	route, err := s.getOrCreateRoute(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get process route: %w", err)
	}

	// 獲取當前成本參數
	params, err := s.costRepo.GetCurrentCostParameters()
	if err != nil {
		return nil, fmt.Errorf("failed to get cost parameters: %w", err)
	}

	// 開始事務
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 創建成本計算主檔
	calculation := &models.CostCalculation{
		InquiryID:        req.InquiryID,
		CalculationNo:    calcNo,
		ProductName:      req.ProductName,
		Quantity:         req.Quantity,
		MaterialCost:     req.MaterialCost,
		RouteID:          &route.ID,
		Status:           "draft",
		MarginPercentage: req.MarginPercentage,
	}

	if err := tx.Create(calculation).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create calculation: %w", err)
	}

	// 計算各製程步驟成本
	var details []models.CostCalculationDetail
	var totalProcessCost float64

	for i, routeDetail := range route.RouteDetails {
		detail, err := s.calculateProcessStepCost(routeDetail, req.Quantity, params, i+1)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to calculate process step cost: %w", err)
		}
		detail.CalculationID = calculation.ID
		details = append(details, *detail)
		totalProcessCost += detail.SubtotalCost + detail.YieldLossCost
	}

	// 批量創建明細
	if len(details) > 0 {
		if err := tx.Create(&details).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create calculation details: %w", err)
		}
	}

	// 計算總成本
	overheadRate := s.getParameterValue(params, "overhead_rate", 150.0) / 100.0
	calculation.ProcessCost = totalProcessCost
	calculation.OverheadCost = totalProcessCost * overheadRate
	calculation.TotalCost = calculation.MaterialCost + calculation.ProcessCost + calculation.OverheadCost
	calculation.UnitCost = calculation.TotalCost / float64(calculation.Quantity)

	// 計算建議售價
	if calculation.MarginPercentage == 0 {
		calculation.MarginPercentage = 30.0 // 預設毛利率
	}
	calculation.SellingPrice = calculation.TotalCost / (1 - calculation.MarginPercentage/100)

	// 更新計算結果
	if err := tx.Save(calculation).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update calculation: %w", err)
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 重新載入完整資料
	return s.costRepo.GetCalculationByID(calculation.ID)
}

// GetCostSummary 獲取成本摘要
func (s *CostCalculationService) GetCostSummary(calculationID uuid.UUID) (*models.CostSummary, error) {
	calc, err := s.costRepo.GetCalculationByID(calculationID)
	if err != nil {
		return nil, err
	}

	summary := &models.CostSummary{
		MaterialCost:     calc.MaterialCost,
		ProcessCost:      calc.ProcessCost,
		OverheadCost:     calc.OverheadCost,
		TotalCost:        calc.TotalCost,
		UnitCost:         calc.UnitCost,
		SuggestedPrice:   calc.SellingPrice,
		MarginPercentage: calc.MarginPercentage,
	}

	// 整理製程成本明細
	for _, detail := range calc.Details {
		breakdown := models.ProcessCostBreakdown{
			ProcessName:     detail.ProcessStep.Name,
			TotalTimeHours:  detail.TotalTimeHours,
			LaborCost:       detail.LaborCost,
			EquipmentCost:   detail.EquipmentCost,
			ElectricityCost: detail.ElectricityCost,
			TotalCost:       detail.SubtotalCost + detail.YieldLossCost,
		}
		if detail.Equipment != nil {
			breakdown.EquipmentName = detail.Equipment.Name
		}
		summary.ProcessBreakdown = append(summary.ProcessBreakdown, breakdown)
	}

	return summary, nil
}

// GetProcessRoutes 獲取製程路線列表
func (s *CostCalculationService) GetProcessRoutes(productCategory string) ([]models.ProductProcessRoute, error) {
	return s.costRepo.GetProcessRoutes(productCategory)
}

// GetProcessSteps 獲取製程步驟列表
func (s *CostCalculationService) GetProcessSteps() ([]models.ProcessStep, error) {
	return s.costRepo.GetAllProcessSteps()
}

// GetEquipmentList 獲取設備列表
func (s *CostCalculationService) GetEquipmentList(categoryID *uuid.UUID) ([]models.Equipment, error) {
	return s.costRepo.GetEquipmentList(categoryID)
}

// ApproveCalculation 審核成本計算
func (s *CostCalculationService) ApproveCalculation(calculationID uuid.UUID, approverID uuid.UUID) error {
	calc, err := s.costRepo.GetCalculationByID(calculationID)
	if err != nil {
		return err
	}

	if calc.Status != "submitted" {
		return errors.New("only submitted calculations can be approved")
	}

	now := time.Now()
	calc.Status = "approved"
	calc.ApprovedBy = &approverID
	calc.ApprovedAt = &now

	return s.db.Save(calc).Error
}

// Private methods

func (s *CostCalculationService) getOrCreateRoute(req models.CostCalculationRequest) (*models.ProductProcessRoute, error) {
	// 如果指定了路線ID，直接使用
	if req.RouteID != nil {
		return s.costRepo.GetProcessRouteByID(*req.RouteID)
	}

	// 如果提供了自定義路線，創建臨時路線
	if len(req.CustomRoute) > 0 {
		route := &models.ProductProcessRoute{
			ProductCategory: req.ProductCategory,
			MaterialType:    req.MaterialType,
			SizeRange:       req.SizeRange,
			RouteName:       fmt.Sprintf("Custom Route - %s", req.ProductName),
			IsDefault:       false,
			IsActive:        true,
		}

		if err := s.db.Create(route).Error; err != nil {
			return nil, err
		}

		// 創建路線明細
		for i, step := range req.CustomRoute {
			detail := models.ProcessRouteDetail{
				RouteID:       route.ID,
				Sequence:      i + 1,
				ProcessStepID: step.ProcessStepID,
				EquipmentID:   &step.EquipmentID,
				YieldRate:     98.0, // 預設良率
			}
			if step.SetupTime > 0 {
				detail.SetupTimeOverride = &step.SetupTime
			}
			if step.CycleTime > 0 {
				detail.CycleTimeOverride = &step.CycleTime
			}
			route.RouteDetails = append(route.RouteDetails, detail)
		}

		if err := s.db.Create(&route.RouteDetails).Error; err != nil {
			return nil, err
		}

		return route, nil
	}

	// 查找預設路線
	routes, err := s.costRepo.GetProcessRoutes(req.ProductCategory)
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.IsDefault {
			fullRoute, err := s.costRepo.GetProcessRouteByID(route.ID)
			if err != nil {
				return nil, err
			}
			return fullRoute, nil
		}
	}

	if len(routes) > 0 {
		// 沒有預設路線，使用第一個
		fullRoute, err := s.costRepo.GetProcessRouteByID(routes[0].ID)
		if err != nil {
			return nil, err
		}
		return fullRoute, nil
	}

	return nil, errors.New("no process route found for product category")
}

func (s *CostCalculationService) calculateProcessStepCost(
	routeDetail models.ProcessRouteDetail,
	quantity int,
	params map[string]models.CostParameter,
	sequence int,
) (*models.CostCalculationDetail, error) {
	// 獲取製程步驟和設備資訊
	step, err := s.costRepo.GetProcessStepByID(routeDetail.ProcessStepID)
	if err != nil {
		return nil, err
	}

	var equipment *models.Equipment
	equipmentID := routeDetail.EquipmentID
	if equipmentID == nil && step.DefaultEquipmentID != nil {
		equipmentID = step.DefaultEquipmentID
	}
	if equipmentID != nil {
		equipment, err = s.costRepo.GetEquipmentByID(*equipmentID)
		if err != nil {
			return nil, err
		}
	}

	// 使用覆蓋值或預設值
	setupTime := step.SetupTimeMinutes
	if routeDetail.SetupTimeOverride != nil {
		setupTime = *routeDetail.SetupTimeOverride
	}

	cycleTime := step.CycleTimeSeconds
	if routeDetail.CycleTimeOverride != nil {
		cycleTime = *routeDetail.CycleTimeOverride
	}

	// 計算總工時（小時）
	totalTimeHours := setupTime/60 + (cycleTime*float64(quantity))/3600

	// 獲取成本參數
	laborCostPerHour := s.getParameterValue(params, "labor_cost", 15.0)
	electricityCostPerKWH := s.getParameterValue(params, "electricity_cost", 0.12)

	// 計算人工成本
	laborCost := totalTimeHours * laborCostPerHour * float64(step.LaborRequired)

	// 計算設備成本（折舊 + 維護）
	var equipmentCost, electricityCost float64
	if equipment != nil {
		// 年工作小時數（假設 250 天 * 8 小時）
		annualHours := 2000.0
		hourlyDepreciation := equipment.PurchaseCost / (float64(equipment.DepreciationYears) * annualHours)
		hourlyMaintenance := equipment.MaintenanceCostPerYear / annualHours
		equipmentCost = totalTimeHours * (hourlyDepreciation + hourlyMaintenance)

		// 計算電力成本
		electricityCost = totalTimeHours * equipment.PowerConsumption * electricityCostPerKWH
	}

	// 小計成本
	subtotalCost := laborCost + equipmentCost + electricityCost

	// 計算良率損失成本
	yieldLoss := (100 - routeDetail.YieldRate) / 100
	yieldLossCost := subtotalCost * yieldLoss

	detail := &models.CostCalculationDetail{
		Sequence:        sequence,
		ProcessStepID:   step.ID,
		ProcessStep:     step,
		EquipmentID:     equipmentID,
		Equipment:       equipment,
		SetupTime:       setupTime,
		CycleTime:       cycleTime,
		TotalTimeHours:  totalTimeHours,
		LaborCost:       math.Round(laborCost*100) / 100,
		EquipmentCost:   math.Round(equipmentCost*100) / 100,
		ElectricityCost: math.Round(electricityCost*100) / 100,
		SubtotalCost:    math.Round(subtotalCost*100) / 100,
		YieldLossCost:   math.Round(yieldLossCost*100) / 100,
	}

	return detail, nil
}

func (s *CostCalculationService) getParameterValue(params map[string]models.CostParameter, paramType string, defaultValue float64) float64 {
	if param, exists := params[paramType]; exists {
		return param.Value
	}
	return defaultValue
}

func (s *CostCalculationService) generateCalculationNo() string {
	// 格式: CALC-YYYYMMDD-XXXX
	date := time.Now().Format("20060102")
	
	// 獲取今日序號
	var count int64
	s.db.Model(&models.CostCalculation{}).
		Where("calculation_no LIKE ?", fmt.Sprintf("CALC-%s-%%", date)).
		Count(&count)
	
	return fmt.Sprintf("CALC-%s-%04d", date, count+1)
}
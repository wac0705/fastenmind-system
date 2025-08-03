package repository

import (
	"strings"
	"time"
	
	"github.com/fastenmind/fastener-api/internal/models"
	"gorm.io/gorm"
)

type MaterialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

// GetMaterials 獲取材料列表
func (r *MaterialRepository) GetMaterials(companyID, materialType string, offset, limit int) ([]models.MaterialCostNew, int64, error) {
	var materials []models.MaterialCostNew
	var total int64
	
	query := r.db.Model(&models.MaterialCostNew{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID)
	
	if materialType != "" {
		query = query.Where("type = ?", materialType)
	}
	
	// 計算總數
	query.Count(&total)
	
	// 獲取分頁數據
	err := query.
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Find(&materials).Error
	
	return materials, total, err
}

// GetByID 根據ID獲取材料
func (r *MaterialRepository) GetByID(id, companyID string) (*models.MaterialCostNew, error) {
	var material models.MaterialCostNew
	err := r.db.Where("id = ? AND company_id = ? AND deleted_at IS NULL", id, companyID).
		First(&material).Error
	return &material, err
}

// Create 創建材料
func (r *MaterialRepository) Create(material *models.MaterialCostNew) error {
	return r.db.Create(material).Error
}

// Update 更新材料
func (r *MaterialRepository) Update(material *models.MaterialCostNew) error {
	return r.db.Save(material).Error
}

// Delete 刪除材料
func (r *MaterialRepository) Delete(id, companyID string) error {
	return r.db.Model(&models.MaterialCostNew{}).
		Where("id = ? AND company_id = ?", id, companyID).
		Update("deleted_at", gorm.DeletedAt{}).Error
}

// SavePriceHistory 保存價格歷史
func (r *MaterialRepository) SavePriceHistory(history *models.MaterialPriceHistory) error {
	return r.db.Create(history).Error
}

// GetPriceHistory 獲取價格歷史
func (r *MaterialRepository) GetPriceHistory(materialID string, limit int) ([]models.MaterialPriceHistory, error) {
	var histories []models.MaterialPriceHistory
	
	err := r.db.Where("material_id = ?", materialID).
		Order("changed_at DESC").
		Limit(limit).
		Find(&histories).Error
	
	return histories, err
}

// SearchMaterials 搜尋材料
func (r *MaterialRepository) SearchMaterials(companyID, keyword string) ([]models.MaterialCostNew, error) {
	var materials []models.MaterialCostNew
	
	// 轉義 LIKE 查詢中的特殊字符
	escapedKeyword := strings.ReplaceAll(keyword, "%", "\\%")
	escapedKeyword = strings.ReplaceAll(escapedKeyword, "_", "\\_")
	likePattern := "%" + escapedKeyword + "%"
	
	err := r.db.Where("company_id = ? AND deleted_at IS NULL", companyID).
		Where("name LIKE ? OR specification LIKE ?", likePattern, likePattern).
		Limit(20).
		Find(&materials).Error
	
	return materials, err
}

// GetMaterialsBySupplier 根據供應商獲取材料
func (r *MaterialRepository) GetMaterialsBySupplier(companyID, supplierID string) ([]models.MaterialCostNew, error) {
	var materials []models.MaterialCostNew
	
	err := r.db.Where("company_id = ? AND supplier_id = ? AND deleted_at IS NULL", 
		companyID, supplierID).
		Find(&materials).Error
	
	return materials, err
}

// UpdateBulkPrices 批量更新價格
func (r *MaterialRepository) UpdateBulkPrices(updates []models.MaterialPriceUpdate, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			// 獲取現有材料
			var material models.MaterialCostNew
			if err := tx.Where("id = ?", update.MaterialID).First(&material).Error; err != nil {
				return err
			}
			
			// 保存價格歷史
			history := &models.MaterialPriceHistory{
				MaterialID: update.MaterialID,
				OldPrice:   material.UnitPrice,
				NewPrice:   update.NewPrice,
				ChangedBy:  userID,
				Reason:     update.Reason,
			}
			if err := tx.Create(history).Error; err != nil {
				return err
			}
			
			// 更新價格
			if err := tx.Model(&models.MaterialCostNew{}).
				Where("id = ?", update.MaterialID).
				Updates(map[string]interface{}{
					"unit_price": update.NewPrice,
					"updated_by": userID,
				}).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// GetMaterialStatistics 獲取材料統計
func (r *MaterialRepository) GetMaterialStatistics(companyID string) (*models.MaterialStatistics, error) {
	stats := &models.MaterialStatistics{}
	
	// 總材料數
	r.db.Model(&models.MaterialCostNew{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Count(&stats.TotalMaterials)
	
	// 各類型材料數
	var typeStats []struct {
		Type  string
		Count int64
	}
	
	r.db.Model(&models.MaterialCostNew{}).
		Select("type, COUNT(*) as count").
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Group("type").
		Scan(&typeStats)
	
	stats.ByType = make(map[string]int64)
	for _, ts := range typeStats {
		stats.ByType[ts.Type] = ts.Count
	}
	
	// 最近更新
	r.db.Model(&models.MaterialCostNew{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Order("updated_at DESC").
		Limit(1).
		Pluck("updated_at", &stats.LastUpdated)
	
	// 價格趨勢
	// 這裡簡化實現，實際應該計算價格變化趨勢
	stats.PriceTrend = []models.PriceTrendData{
		{
			Date:      time.Now(),
			AvgPrice:  0,
			MinPrice:  0,
			MaxPrice:  0,
			ItemCount: 0,
		},
	}
	
	return stats, nil
}
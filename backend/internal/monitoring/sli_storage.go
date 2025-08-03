package monitoring

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormSLIStorage implements SLIStorage using GORM
type GormSLIStorage struct {
	db *gorm.DB
}

// SLIDataPointModel represents SLI data in database
type SLIDataPointModel struct {
	ID        uint      `gorm:"primaryKey"`
	SLIID     string    `gorm:"index;not null"`
	Value     float64   `gorm:"not null"`
	Timestamp time.Time `gorm:"index;not null"`
	Tags      string    `gorm:"type:text"`
	Metadata  string    `gorm:"type:text"`
	CreatedAt time.Time
}

func (SLIDataPointModel) TableName() string {
	return "sli_data_points"
}

// NewGormSLIStorage creates a new GORM-based SLI storage
func NewGormSLIStorage(db *gorm.DB) *GormSLIStorage {
	return &GormSLIStorage{db: db}
}

// Store stores an SLI data point
func (s *GormSLIStorage) Store(ctx context.Context, data *SLIDataPoint) error {
	tagsJSON, err := json.Marshal(data.Tags)
	if err != nil {
		return err
	}

	metadataJSON, err := json.Marshal(data.Metadata)
	if err != nil {
		return err
	}

	model := &SLIDataPointModel{
		SLIID:     data.SLIID.String(),
		Value:     data.Value,
		Timestamp: data.Timestamp,
		Tags:      string(tagsJSON),
		Metadata:  string(metadataJSON),
		CreatedAt: time.Now(),
	}

	return s.db.WithContext(ctx).Create(model).Error
}

// Query queries SLI data points
func (s *GormSLIStorage) Query(ctx context.Context, filter SLIQueryFilter) ([]*SLIDataPoint, error) {
	var models []SLIDataPointModel
	
	query := s.db.WithContext(ctx).Model(&SLIDataPointModel{})

	// Apply filters
	if len(filter.SLIIDs) > 0 {
		sliIDStrings := make([]string, len(filter.SLIIDs))
		for i, id := range filter.SLIIDs {
			sliIDStrings[i] = id.String()
		}
		query = query.Where("sli_id IN ?", sliIDStrings)
	}

	if !filter.StartTime.IsZero() {
		query = query.Where("timestamp >= ?", filter.StartTime)
	}

	if !filter.EndTime.IsZero() {
		query = query.Where("timestamp <= ?", filter.EndTime)
	}

	query = query.Order("timestamp ASC")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	// Convert models to data points
	dataPoints := make([]*SLIDataPoint, len(models))
	for i, model := range models {
		dp, err := s.modelToDataPoint(&model)
		if err != nil {
			return nil, err
		}
		dataPoints[i] = dp
	}

	return dataPoints, nil
}

// Aggregate aggregates SLI data
func (s *GormSLIStorage) Aggregate(ctx context.Context, filter SLIQueryFilter, aggregation AggregationType) (*SLIAggregation, error) {
	dataPoints, err := s.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(dataPoints) == 0 {
		return &SLIAggregation{
			Value:     0,
			Count:     0,
			StartTime: filter.StartTime,
			EndTime:   filter.EndTime,
		}, nil
	}

	values := make([]float64, len(dataPoints))
	for i, dp := range dataPoints {
		values[i] = dp.Value
	}

	var aggregatedValue float64
	switch aggregation {
	case AggregationAverage:
		aggregatedValue = average(values)
	case AggregationSum:
		aggregatedValue = sum(values)
	case AggregationMin:
		aggregatedValue = min(values)
	case AggregationMax:
		aggregatedValue = max(values)
	case AggregationP50:
		aggregatedValue = percentile(values, 0.5)
	case AggregationP90:
		aggregatedValue = percentile(values, 0.9)
	case AggregationP95:
		aggregatedValue = percentile(values, 0.95)
	case AggregationP99:
		aggregatedValue = percentile(values, 0.99)
	case AggregationCount:
		aggregatedValue = float64(len(values))
	default:
		aggregatedValue = average(values)
	}

	startTime := filter.StartTime
	endTime := filter.EndTime
	if len(dataPoints) > 0 {
		if startTime.IsZero() {
			startTime = dataPoints[0].Timestamp
		}
		if endTime.IsZero() {
			endTime = dataPoints[len(dataPoints)-1].Timestamp
		}
	}

	return &SLIAggregation{
		Value:     aggregatedValue,
		Count:     int64(len(dataPoints)),
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// modelToDataPoint converts database model to data point
func (s *GormSLIStorage) modelToDataPoint(model *SLIDataPointModel) (*SLIDataPoint, error) {
	var tags map[string]string
	if err := json.Unmarshal([]byte(model.Tags), &tags); err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(model.Metadata), &metadata); err != nil {
		return nil, err
	}

	sliID, err := uuid.Parse(model.SLIID)
	if err != nil {
		return nil, err
	}

	return &SLIDataPoint{
		SLIID:     sliID,
		Value:     model.Value,
		Timestamp: model.Timestamp,
		Tags:      tags,
		Metadata:  metadata,
	}, nil
}

// Aggregation helper functions

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return sum(values) / float64(len(values))
}

func sum(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy and sort
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	index := p * float64(len(sorted)-1)
	if index == float64(int(index)) {
		return sorted[int(index)]
	}

	lower := int(index)
	upper := lower + 1
	weight := index - float64(lower)

	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// InMemorySLIStorage implements SLIStorage using in-memory storage
type InMemorySLIStorage struct {
	dataPoints []*SLIDataPoint
	mu         sync.RWMutex
}

// NewInMemorySLIStorage creates a new in-memory SLI storage
func NewInMemorySLIStorage() *InMemorySLIStorage {
	return &InMemorySLIStorage{
		dataPoints: make([]*SLIDataPoint, 0),
	}
}

// Store stores an SLI data point in memory
func (s *InMemorySLIStorage) Store(ctx context.Context, data *SLIDataPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy to avoid mutation
	dp := &SLIDataPoint{
		SLIID:     data.SLIID,
		Value:     data.Value,
		Timestamp: data.Timestamp,
		Tags:      make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}

	// Copy tags and metadata
	for k, v := range data.Tags {
		dp.Tags[k] = v
	}
	for k, v := range data.Metadata {
		dp.Metadata[k] = v
	}

	s.dataPoints = append(s.dataPoints, dp)
	return nil
}

// Query queries SLI data points from memory
func (s *InMemorySLIStorage) Query(ctx context.Context, filter SLIQueryFilter) ([]*SLIDataPoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*SLIDataPoint

	for _, dp := range s.dataPoints {
		// Apply filters
		if len(filter.SLIIDs) > 0 {
			found := false
			for _, id := range filter.SLIIDs {
				if dp.SLIID == id {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if !filter.StartTime.IsZero() && dp.Timestamp.Before(filter.StartTime) {
			continue
		}

		if !filter.EndTime.IsZero() && dp.Timestamp.After(filter.EndTime) {
			continue
		}

		// Tag filtering
		if len(filter.Tags) > 0 {
			matches := true
			for key, value := range filter.Tags {
				if dp.Tags[key] != value {
					matches = false
					break
				}
			}
			if !matches {
				continue
			}
		}

		result = append(result, dp)
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	// Apply limit
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result, nil
}

// Aggregate aggregates SLI data from memory
func (s *InMemorySLIStorage) Aggregate(ctx context.Context, filter SLIQueryFilter, aggregation AggregationType) (*SLIAggregation, error) {
	dataPoints, err := s.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(dataPoints) == 0 {
		return &SLIAggregation{
			Value:     0,
			Count:     0,
			StartTime: filter.StartTime,
			EndTime:   filter.EndTime,
		}, nil
	}

	values := make([]float64, len(dataPoints))
	for i, dp := range dataPoints {
		values[i] = dp.Value
	}

	var aggregatedValue float64
	switch aggregation {
	case AggregationAverage:
		aggregatedValue = average(values)
	case AggregationSum:
		aggregatedValue = sum(values)
	case AggregationMin:
		aggregatedValue = min(values)
	case AggregationMax:
		aggregatedValue = max(values)
	case AggregationP50:
		aggregatedValue = percentile(values, 0.5)
	case AggregationP90:
		aggregatedValue = percentile(values, 0.9)
	case AggregationP95:
		aggregatedValue = percentile(values, 0.95)
	case AggregationP99:
		aggregatedValue = percentile(values, 0.99)
	case AggregationCount:
		aggregatedValue = float64(len(values))
	default:
		aggregatedValue = average(values)
	}

	startTime := filter.StartTime
	endTime := filter.EndTime
	if len(dataPoints) > 0 {
		if startTime.IsZero() {
			startTime = dataPoints[0].Timestamp
		}
		if endTime.IsZero() {
			endTime = dataPoints[len(dataPoints)-1].Timestamp
		}
	}

	return &SLIAggregation{
		Value:     aggregatedValue,
		Count:     int64(len(dataPoints)),
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}
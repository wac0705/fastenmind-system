package archival

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ArchivalService manages data archival and retrieval
type ArchivalService struct {
	db              *gorm.DB
	config          *ArchivalConfig
	storageBackends map[string]StorageBackend
	policies        map[string]*ArchivalPolicy
	mu              sync.RWMutex
}

// ArchivalConfig holds archival configuration
type ArchivalConfig struct {
	// Archive settings
	DefaultRetentionDays int
	CompressionEnabled   bool
	EncryptionEnabled    bool
	EncryptionKey        []byte
	
	// Storage settings
	LocalPath      string
	S3Bucket       string
	GlacierVault   string
	AzureContainer string
	
	// Performance settings
	BatchSize        int
	WorkerCount      int
	ArchiveSchedule  string // cron expression
	
	// Notification settings
	NotificationEmail string
	WebhookURL        string
}

// ArchivalPolicy defines archival rules for specific data types
type ArchivalPolicy struct {
	Name             string
	TableName        string
	RetentionDays    int
	ArchiveCondition string // SQL condition
	StorageClass     string // "hot", "warm", "cold"
	CompressData     bool
	EncryptData      bool
	IndexFields      []string
}

// ArchivalRecord tracks archived data
type ArchivalRecord struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	TableName       string
	RecordID        string
	ArchiveDate     time.Time
	StorageBackend  string
	StorageLocation string
	CompressedSize  int64
	OriginalSize    int64
	Checksum        string
	Metadata        json.RawMessage `gorm:"type:jsonb"`
	ExpiryDate      *time.Time
	IsDeleted       bool
	CreatedAt       time.Time
	RestoredAt      *time.Time
}

// StorageBackend interface for different storage systems
type StorageBackend interface {
	Store(ctx context.Context, key string, data []byte, metadata map[string]string) error
	Retrieve(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)
	GetMetadata(ctx context.Context, key string) (map[string]string, error)
}

// NewArchivalService creates a new archival service
func NewArchivalService(db *gorm.DB, config *ArchivalConfig) (*ArchivalService, error) {
	service := &ArchivalService{
		db:              db,
		config:          config,
		storageBackends: make(map[string]StorageBackend),
		policies:        make(map[string]*ArchivalPolicy),
	}
	
	// Initialize storage backends
	if err := service.initializeBackends(); err != nil {
		return nil, err
	}
	
	// Run migrations
	if err := db.AutoMigrate(&ArchivalRecord{}); err != nil {
		return nil, err
	}
	
	return service, nil
}

// initializeBackends sets up storage backends
func (s *ArchivalService) initializeBackends() error {
	// Local filesystem
	if s.config.LocalPath != "" {
		s.storageBackends["local"] = NewLocalStorageBackend(s.config.LocalPath)
	}
	
	// AWS S3
	if s.config.S3Bucket != "" {
		sess, err := session.NewSession()
		if err != nil {
			return err
		}
		s.storageBackends["s3"] = NewS3StorageBackend(sess, s.config.S3Bucket)
	}
	
	// AWS Glacier
	if s.config.GlacierVault != "" {
		sess, err := session.NewSession()
		if err != nil {
			return err
		}
		s.storageBackends["glacier"] = NewGlacierStorageBackend(sess, s.config.GlacierVault)
	}
	
	return nil
}

// RegisterPolicy registers an archival policy
func (s *ArchivalService) RegisterPolicy(policy *ArchivalPolicy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[policy.Name] = policy
}

// ArchiveTable archives data from a table based on policy
func (s *ArchivalService) ArchiveTable(ctx context.Context, tableName string) error {
	policy, exists := s.policies[tableName]
	if !exists {
		return fmt.Errorf("no archival policy found for table %s", tableName)
	}
	
	// Calculate cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -policy.RetentionDays)
	
	// Start archival process
	return s.archiveData(ctx, policy, cutoffDate)
}

// archiveData performs the actual archival
func (s *ArchivalService) archiveData(ctx context.Context, policy *ArchivalPolicy, cutoffDate time.Time) error {
	// Create workers
	workerCount := s.config.WorkerCount
	if workerCount == 0 {
		workerCount = 4
	}
	
	jobs := make(chan map[string]interface{}, s.config.BatchSize)
	errors := make(chan error, workerCount)
	
	var wg sync.WaitGroup
	
	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range jobs {
				if err := s.archiveRecord(ctx, policy, record); err != nil {
					errors <- err
				}
			}
		}()
	}
	
	// Query data to archive
	rows := make([]map[string]interface{}, 0, s.config.BatchSize)
	
	query := s.db.WithContext(ctx).Table(policy.TableName)
	if policy.ArchiveCondition != "" {
		query = query.Where(policy.ArchiveCondition)
	}
	query = query.Where("created_at < ?", cutoffDate)
	
	// Process in batches
	offset := 0
	for {
		batch := make([]map[string]interface{}, 0, s.config.BatchSize)
		err := query.Offset(offset).Limit(s.config.BatchSize).Find(&batch).Error
		if err != nil {
			close(jobs)
			return err
		}
		
		if len(batch) == 0 {
			break
		}
		
		// Send to workers
		for _, record := range batch {
			jobs <- record
		}
		
		offset += len(batch)
		
		// Check for cancellation
		select {
		case <-ctx.Done():
			close(jobs)
			return ctx.Err()
		default:
		}
	}
	
	close(jobs)
	wg.Wait()
	
	// Check for errors
	close(errors)
	var firstError error
	errorCount := 0
	for err := range errors {
		if firstError == nil {
			firstError = err
		}
		errorCount++
	}
	
	if errorCount > 0 {
		return fmt.Errorf("archival failed with %d errors, first error: %w", errorCount, firstError)
	}
	
	return nil
}

// archiveRecord archives a single record
func (s *ArchivalService) archiveRecord(ctx context.Context, policy *ArchivalPolicy, record map[string]interface{}) error {
	// Extract record ID
	recordID := fmt.Sprintf("%v", record["id"])
	
	// Serialize record
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	
	originalSize := int64(len(data))
	
	// Compress if enabled
	if policy.CompressData || s.config.CompressionEnabled {
		data, err = s.compress(data)
		if err != nil {
			return err
		}
	}
	
	// Encrypt if enabled
	if policy.EncryptData || s.config.EncryptionEnabled {
		data, err = s.encrypt(data)
		if err != nil {
			return err
		}
	}
	
	// Generate storage key
	storageKey := s.generateStorageKey(policy.TableName, recordID)
	
	// Determine storage backend based on policy
	backendName := s.getStorageBackend(policy.StorageClass)
	backend, exists := s.storageBackends[backendName]
	if !exists {
		return fmt.Errorf("storage backend %s not found", backendName)
	}
	
	// Prepare metadata
	metadata := map[string]string{
		"table":        policy.TableName,
		"record_id":    recordID,
		"archive_date": time.Now().Format(time.RFC3339),
		"compressed":   fmt.Sprintf("%v", policy.CompressData),
		"encrypted":    fmt.Sprintf("%v", policy.EncryptData),
	}
	
	// Add index fields to metadata
	for _, field := range policy.IndexFields {
		if value, exists := record[field]; exists {
			metadata[field] = fmt.Sprintf("%v", value)
		}
	}
	
	// Store in backend
	if err := backend.Store(ctx, storageKey, data, metadata); err != nil {
		return err
	}
	
	// Record archival
	archivalRecord := &ArchivalRecord{
		ID:              uuid.New(),
		TableName:       policy.TableName,
		RecordID:        recordID,
		ArchiveDate:     time.Now(),
		StorageBackend:  backendName,
		StorageLocation: storageKey,
		CompressedSize:  int64(len(data)),
		OriginalSize:    originalSize,
		Checksum:        s.calculateChecksum(data),
		Metadata:        s.extractMetadata(record, policy.IndexFields),
		CreatedAt:       time.Now(),
	}
	
	if policy.RetentionDays > 0 {
		expiryDate := time.Now().AddDate(0, 0, policy.RetentionDays)
		archivalRecord.ExpiryDate = &expiryDate
	}
	
	if err := s.db.Create(archivalRecord).Error; err != nil {
		// Rollback storage
		backend.Delete(ctx, storageKey)
		return err
	}
	
	// Delete from source table if configured
	if err := s.db.Table(policy.TableName).Where("id = ?", recordID).Delete(nil).Error; err != nil {
		// Log error but don't fail archival
		fmt.Printf("Failed to delete archived record %s: %v\n", recordID, err)
	}
	
	return nil
}

// RestoreRecord restores an archived record
func (s *ArchivalService) RestoreRecord(ctx context.Context, archivalID uuid.UUID) (map[string]interface{}, error) {
	// Get archival record
	var archivalRecord ArchivalRecord
	if err := s.db.Where("id = ?", archivalID).First(&archivalRecord).Error; err != nil {
		return nil, err
	}
	
	// Get storage backend
	backend, exists := s.storageBackends[archivalRecord.StorageBackend]
	if !exists {
		return nil, fmt.Errorf("storage backend %s not found", archivalRecord.StorageBackend)
	}
	
	// Retrieve data
	data, err := backend.Retrieve(ctx, archivalRecord.StorageLocation)
	if err != nil {
		return nil, err
	}
	
	// Verify checksum
	if checksum := s.calculateChecksum(data); checksum != archivalRecord.Checksum {
		return nil, fmt.Errorf("checksum mismatch: expected %s, got %s", archivalRecord.Checksum, checksum)
	}
	
	// Decrypt if needed
	policy := s.policies[archivalRecord.TableName]
	if policy != nil && (policy.EncryptData || s.config.EncryptionEnabled) {
		data, err = s.decrypt(data)
		if err != nil {
			return nil, err
		}
	}
	
	// Decompress if needed
	if policy != nil && (policy.CompressData || s.config.CompressionEnabled) {
		data, err = s.decompress(data)
		if err != nil {
			return nil, err
		}
	}
	
	// Deserialize
	var record map[string]interface{}
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	
	// Update restoration timestamp
	now := time.Now()
	archivalRecord.RestoredAt = &now
	s.db.Save(&archivalRecord)
	
	return record, nil
}

// SearchArchive searches archived records
func (s *ArchivalService) SearchArchive(ctx context.Context, criteria SearchCriteria) ([]*ArchivalRecord, error) {
	query := s.db.WithContext(ctx).Model(&ArchivalRecord{})
	
	if criteria.TableName != "" {
		query = query.Where("table_name = ?", criteria.TableName)
	}
	
	if criteria.StartDate != nil {
		query = query.Where("archive_date >= ?", *criteria.StartDate)
	}
	
	if criteria.EndDate != nil {
		query = query.Where("archive_date <= ?", *criteria.EndDate)
	}
	
	if len(criteria.Metadata) > 0 {
		for key, value := range criteria.Metadata {
			query = query.Where("metadata->? = ?", key, value)
		}
	}
	
	var records []*ArchivalRecord
	err := query.Order("archive_date DESC").Limit(criteria.Limit).Find(&records).Error
	
	return records, err
}

// SearchCriteria defines search parameters
type SearchCriteria struct {
	TableName string
	StartDate *time.Time
	EndDate   *time.Time
	Metadata  map[string]string
	Limit     int
}

// Cleanup removes expired archives
func (s *ArchivalService) Cleanup(ctx context.Context) error {
	// Find expired records
	var expiredRecords []*ArchivalRecord
	err := s.db.Where("expiry_date < ? AND is_deleted = ?", time.Now(), false).Find(&expiredRecords).Error
	if err != nil {
		return err
	}
	
	// Delete from storage
	for _, record := range expiredRecords {
		backend, exists := s.storageBackends[record.StorageBackend]
		if !exists {
			continue
		}
		
		if err := backend.Delete(ctx, record.StorageLocation); err != nil {
			// Log error but continue
			fmt.Printf("Failed to delete archive %s: %v\n", record.ID, err)
			continue
		}
		
		// Mark as deleted
		record.IsDeleted = true
		s.db.Save(record)
	}
	
	return nil
}

// Helper methods

func (s *ArchivalService) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	
	if err := gz.Close(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

func (s *ArchivalService) decompress(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	
	return io.ReadAll(gz)
}

func (s *ArchivalService) encrypt(data []byte) ([]byte, error) {
	// Use encryption service
	// Placeholder implementation
	return data, nil
}

func (s *ArchivalService) decrypt(data []byte) ([]byte, error) {
	// Use encryption service
	// Placeholder implementation
	return data, nil
}

func (s *ArchivalService) calculateChecksum(data []byte) string {
	// Calculate SHA256 checksum
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func (s *ArchivalService) generateStorageKey(tableName, recordID string) string {
	date := time.Now().Format("2006/01/02")
	return fmt.Sprintf("%s/%s/%s/%s", tableName, date, uuid.New().String(), recordID)
}

func (s *ArchivalService) getStorageBackend(storageClass string) string {
	switch storageClass {
	case "hot":
		return "s3"
	case "warm":
		return "s3" // Could use S3 IA
	case "cold":
		return "glacier"
	default:
		return "local"
	}
}

func (s *ArchivalService) extractMetadata(record map[string]interface{}, indexFields []string) json.RawMessage {
	metadata := make(map[string]interface{})
	
	for _, field := range indexFields {
		if value, exists := record[field]; exists {
			metadata[field] = value
		}
	}
	
	data, _ := json.Marshal(metadata)
	return data
}

// LocalStorageBackend implements local filesystem storage
type LocalStorageBackend struct {
	basePath string
}

func NewLocalStorageBackend(basePath string) *LocalStorageBackend {
	return &LocalStorageBackend{basePath: basePath}
}

func (l *LocalStorageBackend) Store(ctx context.Context, key string, data []byte, metadata map[string]string) error {
	fullPath := filepath.Join(l.basePath, key)
	
	// Create directory
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	
	// Write data
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return err
	}
	
	// Write metadata
	metadataPath := fullPath + ".metadata"
	metadataData, _ := json.Marshal(metadata)
	return os.WriteFile(metadataPath, metadataData, 0644)
}

func (l *LocalStorageBackend) Retrieve(ctx context.Context, key string) ([]byte, error) {
	fullPath := filepath.Join(l.basePath, key)
	return os.ReadFile(fullPath)
}

func (l *LocalStorageBackend) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(l.basePath, key)
	os.Remove(fullPath + ".metadata")
	return os.Remove(fullPath)
}

func (l *LocalStorageBackend) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	
	searchPath := filepath.Join(l.basePath, prefix)
	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && !strings.HasSuffix(path, ".metadata") {
			relPath, _ := filepath.Rel(l.basePath, path)
			keys = append(keys, relPath)
		}
		
		return nil
	})
	
	return keys, err
}

func (l *LocalStorageBackend) GetMetadata(ctx context.Context, key string) (map[string]string, error) {
	metadataPath := filepath.Join(l.basePath, key) + ".metadata"
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}
	
	var metadata map[string]string
	err = json.Unmarshal(data, &metadata)
	return metadata, err
}

// S3StorageBackend implements AWS S3 storage
type S3StorageBackend struct {
	client *s3.S3
	bucket string
}

func NewS3StorageBackend(sess *session.Session, bucket string) *S3StorageBackend {
	return &S3StorageBackend{
		client: s3.New(sess),
		bucket: bucket,
	}
}

func (s *S3StorageBackend) Store(ctx context.Context, key string, data []byte, metadata map[string]string) error {
	input := &s3.PutObjectInput{
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		Body:     bytes.NewReader(data),
		Metadata: aws.StringMap(metadata),
	}
	
	_, err := s.client.PutObjectWithContext(ctx, input)
	return err
}

func (s *S3StorageBackend) Retrieve(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	
	result, err := s.client.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	
	return io.ReadAll(result.Body)
}

func (s *S3StorageBackend) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	
	_, err := s.client.DeleteObjectWithContext(ctx, input)
	return err
}

func (s *S3StorageBackend) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}
	
	err := s.client.ListObjectsV2PagesWithContext(ctx, input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
		return true
	})
	
	return keys, err
}

func (s *S3StorageBackend) GetMetadata(ctx context.Context, key string) (map[string]string, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	
	result, err := s.client.HeadObjectWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	
	metadata := make(map[string]string)
	for k, v := range result.Metadata {
		metadata[k] = *v
	}
	
	return metadata, nil
}

// GlacierStorageBackend implements AWS Glacier storage
type GlacierStorageBackend struct {
	client *glacier.Glacier
	vault  string
}

func NewGlacierStorageBackend(sess *session.Session, vault string) *GlacierStorageBackend {
	return &GlacierStorageBackend{
		client: glacier.New(sess),
		vault:  vault,
	}
}

func (g *GlacierStorageBackend) Store(ctx context.Context, key string, data []byte, metadata map[string]string) error {
	// Convert metadata to description
	metadataJSON, _ := json.Marshal(metadata)
	
	input := &glacier.UploadArchiveInput{
		VaultName:          aws.String(g.vault),
		Body:               bytes.NewReader(data),
		ArchiveDescription: aws.String(string(metadataJSON)),
	}
	
	_, err := g.client.UploadArchiveWithContext(ctx, input)
	return err
}

func (g *GlacierStorageBackend) Retrieve(ctx context.Context, key string) ([]byte, error) {
	// Glacier retrieval is async and complex
	// This is a simplified placeholder
	return nil, fmt.Errorf("glacier retrieval not implemented")
}

func (g *GlacierStorageBackend) Delete(ctx context.Context, key string) error {
	input := &glacier.DeleteArchiveInput{
		VaultName: aws.String(g.vault),
		ArchiveId: aws.String(key),
	}
	
	_, err := g.client.DeleteArchiveWithContext(ctx, input)
	return err
}

func (g *GlacierStorageBackend) List(ctx context.Context, prefix string) ([]string, error) {
	// Glacier inventory is async
	// This is a simplified placeholder
	return nil, fmt.Errorf("glacier listing not implemented")
}

func (g *GlacierStorageBackend) GetMetadata(ctx context.Context, key string) (map[string]string, error) {
	// Would need to retrieve from inventory
	return nil, fmt.Errorf("glacier metadata not implemented")
}
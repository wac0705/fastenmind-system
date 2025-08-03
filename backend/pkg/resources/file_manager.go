package resources

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileManager 文件資源管理器
type FileManager struct {
	openFiles map[string]*ManagedFile
	tempFiles map[string]string // 臨時文件路徑
	mu        sync.RWMutex
}

// ManagedFile 被管理的文件
type ManagedFile struct {
	file      *os.File
	path      string
	openedAt  time.Time
	lastUsed  time.Time
	readOnly  bool
	mu        sync.Mutex
}

// NewFileManager 創建新的文件管理器
func NewFileManager() *FileManager {
	return &FileManager{
		openFiles: make(map[string]*ManagedFile),
		tempFiles: make(map[string]string),
	}
}

// OpenFile 打開文件（自動管理關閉）
func (m *FileManager) OpenFile(path string, flag int, perm os.FileMode) (*ManagedFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 檢查文件是否已打開
	if mf, exists := m.openFiles[path]; exists {
		mf.mu.Lock()
		mf.lastUsed = time.Now()
		mf.mu.Unlock()
		return mf, nil
	}
	
	// 打開新文件
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	
	mf := &ManagedFile{
		file:     file,
		path:     path,
		openedAt: time.Now(),
		lastUsed: time.Now(),
		readOnly: flag&os.O_WRONLY == 0 && flag&os.O_RDWR == 0,
	}
	
	m.openFiles[path] = mf
	return mf, nil
}

// CreateTempFile 創建臨時文件
func (m *FileManager) CreateTempFile(dir, pattern string) (*ManagedFile, error) {
	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	
	path := file.Name()
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	mf := &ManagedFile{
		file:     file,
		path:     path,
		openedAt: time.Now(),
		lastUsed: time.Now(),
		readOnly: false,
	}
	
	m.openFiles[path] = mf
	m.tempFiles[path] = path
	
	return mf, nil
}

// Read 讀取文件
func (f *ManagedFile) Read(p []byte) (n int, err error) {
	f.mu.Lock()
	f.lastUsed = time.Now()
	f.mu.Unlock()
	
	return f.file.Read(p)
}

// Write 寫入文件
func (f *ManagedFile) Write(p []byte) (n int, err error) {
	if f.readOnly {
		return 0, fmt.Errorf("file is read-only")
	}
	
	f.mu.Lock()
	f.lastUsed = time.Now()
	f.mu.Unlock()
	
	return f.file.Write(p)
}

// Path 獲取文件路徑
func (f *ManagedFile) Path() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.path
}

// Close 關閉文件
func (f *ManagedFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if f.file == nil {
		return nil
	}
	
	err := f.file.Close()
	f.file = nil
	return err
}

// CloseFile 關閉指定文件
func (m *FileManager) CloseFile(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	mf, exists := m.openFiles[path]
	if !exists {
		return nil
	}
	
	err := mf.Close()
	delete(m.openFiles, path)
	
	// 如果是臨時文件，刪除它
	if _, isTemp := m.tempFiles[path]; isTemp {
		os.Remove(path)
		delete(m.tempFiles, path)
	}
	
	return err
}

// CloseAllFiles 關閉所有打開的文件
func (m *FileManager) CloseAllFiles() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var firstErr error
	
	// 關閉所有文件
	for path, mf := range m.openFiles {
		if err := mf.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(m.openFiles, path)
	}
	
	// 刪除所有臨時文件
	for path := range m.tempFiles {
		os.Remove(path)
		delete(m.tempFiles, path)
	}
	
	return firstErr
}

// CleanupIdleFiles 清理空閒文件
func (m *FileManager) CleanupIdleFiles(maxIdleTime time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	now := time.Now()
	closed := 0
	
	for path, mf := range m.openFiles {
		mf.mu.Lock()
		if now.Sub(mf.lastUsed) > maxIdleTime {
			mf.file.Close()
			delete(m.openFiles, path)
			
			// 如果是臨時文件，刪除它
			if _, isTemp := m.tempFiles[path]; isTemp {
				os.Remove(path)
				delete(m.tempFiles, path)
			}
			
			closed++
		}
		mf.mu.Unlock()
	}
	
	return closed
}

// Stats 獲取統計信息
func (m *FileManager) Stats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := make(map[string]interface{})
	stats["open_files"] = len(m.openFiles)
	stats["temp_files"] = len(m.tempFiles)
	
	files := make([]map[string]interface{}, 0, len(m.openFiles))
	for path, mf := range m.openFiles {
		mf.mu.Lock()
		fileInfo := map[string]interface{}{
			"path":      path,
			"opened_at": mf.openedAt,
			"last_used": mf.lastUsed,
			"read_only": mf.readOnly,
			"is_temp":   m.tempFiles[path] != "",
		}
		mf.mu.Unlock()
		files = append(files, fileInfo)
	}
	stats["files"] = files
	
	return stats
}

// FileReader 文件讀取器（自動關閉）
type FileReader struct {
	file   *os.File
	closed bool
	mu     sync.Mutex
}

// NewFileReader 創建文件讀取器
func NewFileReader(path string) (*FileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	
	return &FileReader{
		file:   file,
		closed: false,
	}, nil
}

// Read 讀取數據
func (r *FileReader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.closed {
		return 0, io.ErrClosedPipe
	}
	
	return r.file.Read(p)
}

// Close 關閉文件
func (r *FileReader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.closed {
		return nil
	}
	
	r.closed = true
	return r.file.Close()
}

// ReadFileWithCleanup 讀取文件並自動清理
func ReadFileWithCleanup(ctx context.Context, path string, handler func(io.Reader) error) error {
	reader, err := NewFileReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	
	// 創建一個通道來處理上下文取消
	done := make(chan error, 1)
	
	go func() {
		done <- handler(reader)
	}()
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// WriteFileWithCleanup 寫入文件並自動清理
func WriteFileWithCleanup(ctx context.Context, path string, handler func(io.Writer) error) error {
	// 創建臨時文件
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	temp, err := os.CreateTemp(dir, base+".tmp*")
	if err != nil {
		return err
	}
	
	tempPath := temp.Name()
	
	// 確保清理臨時文件
	defer func() {
		temp.Close()
		os.Remove(tempPath)
	}()
	
	// 寫入數據
	if err := handler(temp); err != nil {
		return err
	}
	
	// 關閉臨時文件
	if err := temp.Close(); err != nil {
		return err
	}
	
	// 原子性地替換原文件
	return os.Rename(tempPath, path)
}
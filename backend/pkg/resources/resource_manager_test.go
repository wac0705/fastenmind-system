package resources

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientManager(t *testing.T) {
	manager := NewHTTPClientManager()

	// Test getting HTTP client
	client1 := manager.GetClient("test1", 10*time.Second)
	assert.NotNil(t, client1)

	// Test getting same client returns cached instance
	client2 := manager.GetClient("test1", 10*time.Second)
	assert.Equal(t, client1, client2)

	// Test different name returns different client
	client3 := manager.GetClient("test2", 5*time.Second)
	assert.NotNil(t, client3)
	assert.NotEqual(t, client1, client3)

	// Test cleanup idle clients
	cleaned := manager.CleanupIdleClients(1 * time.Nanosecond)
	assert.GreaterOrEqual(t, cleaned, 0)
}

func TestHTTPClientManagerRequests(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	manager := NewHTTPClientManager()

	client := manager.GetClient("test", 5*time.Second)
	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFileManager(t *testing.T) {
	manager := NewFileManager()

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	require.NoError(t, err)
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Write test data
	testData := []byte("test data")
	err = os.WriteFile(tmpPath, testData, 0644)
	require.NoError(t, err)

	// Test opening file
	file1, err := manager.OpenFile(tmpPath, os.O_RDONLY, 0644)
	require.NoError(t, err)
	defer manager.CloseFile(tmpPath)

	// Test file is cached
	file2, err := manager.OpenFile(tmpPath, os.O_RDONLY, 0644)
	require.NoError(t, err)
	assert.Equal(t, file1, file2)

	// Test reading from file
	buf := make([]byte, len(testData))
	n, err := file1.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Equal(t, testData, buf)

	// Test closing file
	err = manager.CloseFile(tmpPath)
	assert.NoError(t, err)

	// Test cleanup idle files
	cleaned := manager.CleanupIdleFiles(1 * time.Nanosecond)
	assert.GreaterOrEqual(t, cleaned, 0)
}

func TestFileManagerMultipleFiles(t *testing.T) {
	manager := NewFileManager()

	// Create multiple temporary files
	files := make([]string, 3)
	for i := range files {
		tmpFile, err := os.CreateTemp("", "test_file_*.txt")
		require.NoError(t, err)
		files[i] = tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(files[i])
	}

	// Open all files
	for _, path := range files {
		_, err := manager.OpenFile(path, os.O_RDONLY, 0644)
		assert.NoError(t, err)
	}

	// Close all files
	for _, path := range files {
		err := manager.CloseFile(path)
		assert.NoError(t, err)
	}
}

func TestResourceManager(t *testing.T) {
	manager := NewResourceManager()

	// Test HTTP client access
	httpManager := manager.HTTPClients()
	assert.NotNil(t, httpManager)
	client := httpManager.GetClient("test", 10*time.Second)
	assert.NotNil(t, client)

	// Test file manager access
	fileManager := manager.Files()
	assert.NotNil(t, fileManager)

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	require.NoError(t, err)
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Test file operations
	file, err := fileManager.OpenFile(tmpPath, os.O_RDONLY, 0644)
	require.NoError(t, err)
	defer fileManager.CloseFile(tmpPath)

	assert.NotNil(t, file)

	// Test starting cleanup goroutine
	ctx := context.Background()
	manager.StartCleanup(ctx)
}

func TestResourceManagerConcurrent(t *testing.T) {
	manager := NewResourceManager()

	done := make(chan bool, 10)

	// Concurrent HTTP client access
	httpManager := manager.HTTPClients()
	for i := 0; i < 5; i++ {
		go func(id int) {
			client := httpManager.GetClient("test", 10*time.Second)
			assert.NotNil(t, client)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}
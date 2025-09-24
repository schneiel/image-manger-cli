package log_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

// MockLocalizer is a simple mock for testing.
type MockLocalizer struct {
	TranslateFunc func(messageID string, templateData ...map[string]interface{}) string
}

func (m *MockLocalizer) Translate(messageID string, templateData ...map[string]interface{}) string {
	if m.TranslateFunc != nil {
		return m.TranslateFunc(messageID, templateData...)
	}
	return messageID
}

func (m *MockLocalizer) GetCurrentLanguage() string { return "en" }
func (m *MockLocalizer) SetLanguage(_ string) error { return nil }
func (m *MockLocalizer) IsInitialized() bool        { return true }

// MockFileSystem is a simple mock for testing.
type MockFileSystem struct {
	OpenFileFunc func(name string, flag int, perm os.FileMode) (filesystem.File, error)
	StatFunc     func(name string) (os.FileInfo, error)
}

func (m *MockFileSystem) OpenFile(name string, flag int, perm os.FileMode) (filesystem.File, error) {
	if m.OpenFileFunc != nil {
		return m.OpenFileFunc(name, flag, perm)
	}
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	return nil, errors.New("not implemented")
}

// Implement other required methods with no-op implementations.
func (m *MockFileSystem) Create(_ string) (filesystem.File, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) Open(_ string) (filesystem.File, error) {
	return nil, errors.New("not implemented")
}
func (m *MockFileSystem) Remove(_ string) error    { return errors.New("not implemented") }
func (m *MockFileSystem) RemoveAll(_ string) error { return errors.New("not implemented") }
func (m *MockFileSystem) Rename(_, _ string) error { return errors.New("not implemented") }
func (m *MockFileSystem) Mkdir(_ string, _ os.FileMode) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) MkdirAll(_ string, _ os.FileMode) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) ReadDir(_ string) ([]os.DirEntry, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) Lstat(_ string) (os.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) ReadFile(_ string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) WriteFile(_ string, _ []byte, _ os.FileMode) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) CreateTemp(_, _ string) (filesystem.File, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) MkdirTemp(_, _ string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *MockFileSystem) Chmod(_ string, _ os.FileMode) error {
	return errors.New("not implemented")
}
func (m *MockFileSystem) Chown(_ string, _, _ int) error { return errors.New("not implemented") }
func (m *MockFileSystem) Chtimes(_ string, _, _ time.Time) error {
	return errors.New("not implemented")
}
func (m *MockFileSystem) Getwd() (string, error)    { return "", errors.New("not implemented") }
func (m *MockFileSystem) Chdir(_ string) error      { return errors.New("not implemented") }
func (m *MockFileSystem) Symlink(_, _ string) error { return errors.New("not implemented") }
func (m *MockFileSystem) Link(_, _ string) error    { return errors.New("not implemented") }
func (m *MockFileSystem) Readlink(_ string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *MockFileSystem) WalkDir(_ string, _ fs.WalkDirFunc) error {
	return errors.New("not implemented")
}
func (m *MockFileSystem) IsNotExist(err error) bool { return os.IsNotExist(err) }

// MockFile is a simple mock for testing.
type MockFile struct {
	*os.File
}

func (m *MockFile) Close() error {
	if m.File != nil {
		return m.File.Close()
	}
	return nil
}

func TestNewDefaultFileLogger_Success(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	logger, err := log.NewDefaultFileLogger(logPath, log.INFO)

	require.NoError(t, err)
	assert.NotNil(t, logger)

	// Clean up
	err = logger.Close()
	assert.NoError(t, err)
}

func TestNewDefaultFileLogger_FileCreationError(t *testing.T) {
	t.Parallel()

	// Use an invalid path that should fail
	invalidPath := "/invalid/path/that/does/not/exist/test.log"

	logger, err := log.NewDefaultFileLogger(invalidPath, log.INFO)

	require.Error(t, err)
	assert.Nil(t, logger)
	assert.Contains(t, err.Error(), "ErrorOpeningLogFile")
}

func TestNewDefaultFileLoggerWithLocalizer_Success(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")
	mockLocalizer := &testutils.MockLocalizer{}

	logger, err := log.NewDefaultFileLoggerWithLocalizer(logPath, log.DEBUG, mockLocalizer)

	require.NoError(t, err)
	assert.NotNil(t, logger)

	// Clean up
	err = logger.Close()
	assert.NoError(t, err)
}

func TestNewDefaultFileLoggerWithLocalizer_FileCreationErrorWithLocalizer(t *testing.T) {
	t.Parallel()

	invalidPath := "/invalid/path/that/does/not/exist/test.log"
	mockLocalizer := &MockLocalizer{
		TranslateFunc: func(messageID string, _ ...map[string]interface{}) string {
			if messageID == "ErrorOpeningLogFile" {
				return "Failed to open log file at /invalid/path/that/does/not/exist/test.log: permission denied"
			}
			return messageID
		},
	}

	logger, err := log.NewDefaultFileLoggerWithLocalizer(invalidPath, log.INFO, mockLocalizer)

	require.Error(t, err)
	assert.Nil(t, logger)
	assert.Contains(t, err.Error(), "Failed to open log file")
}

func TestNewDefaultFileLoggerWithDependencies_Success(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	mockFS := &MockFileSystem{
		OpenFileFunc: func(_ string, flag int, perm os.FileMode) (filesystem.File, error) {
			// Create a real temp file for testing
			return os.OpenFile(logPath, flag, perm) //nolint:gosec // G304: Test needs to open temp file with variable path
		},
	}
	mockLocalizer := &MockLocalizer{}

	logger, err := log.NewDefaultFileLoggerWithDependencies(logPath, log.WARN, mockLocalizer, mockFS)

	require.NoError(t, err)
	assert.NotNil(t, logger)

	// Clean up
	err = logger.Close()
	assert.NoError(t, err)
}

func TestNewDefaultFileLoggerWithDependencies_NilFileSystemPanic(t *testing.T) {
	t.Parallel()

	mockLocalizer := &testutils.MockLocalizer{}

	_, err := log.NewDefaultFileLoggerWithDependencies("/test/path.log", log.INFO, mockLocalizer, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fileSystem cannot be nil")
}

func TestNewDefaultFileLoggerWithDependencies_FileOpenError(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("permission denied")
	mockFS := &MockFileSystem{
		OpenFileFunc: func(_ string, _ int, _ os.FileMode) (filesystem.File, error) {
			return nil, expectedError
		},
	}
	mockLocalizer := &MockLocalizer{
		TranslateFunc: func(messageID string, _ ...map[string]interface{}) string {
			if messageID == "ErrorOpeningLogFile" {
				return "Failed to open log file at /test/path.log: permission denied"
			}
			return messageID
		},
	}

	logger, err := log.NewDefaultFileLoggerWithDependencies("/test/path.log", log.INFO, mockLocalizer, mockFS)

	require.Error(t, err)
	assert.Nil(t, logger)
	assert.Contains(t, err.Error(), "Failed to open log file")
}

func TestNewDefaultFileLoggerWithDependencies_FileOpenErrorNoLocalizer(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("permission denied")
	mockFS := &MockFileSystem{
		OpenFileFunc: func(_ string, _ int, _ os.FileMode) (filesystem.File, error) {
			return nil, expectedError
		},
	}

	logger, err := log.NewDefaultFileLoggerWithDependencies("/test/path.log", log.INFO, nil, mockFS)

	require.Error(t, err)
	assert.Nil(t, logger)
	assert.Contains(t, err.Error(), "ErrorOpeningLogFile")
}

func TestDefaultFileLogger_Close_Success(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	logger, err := log.NewDefaultFileLogger(logPath, log.INFO)
	require.NoError(t, err)

	err = logger.Close()
	require.NoError(t, err)

	// Second close should not error
	err = logger.Close()
	assert.NoError(t, err)
}

func TestDefaultFileLogger_SetLevel(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	logger, err := log.NewDefaultFileLogger(logPath, log.INFO)
	require.NoError(t, err)
	defer func() { _ = logger.Close() }()

	// Test setting different levels
	logger.SetLevel(log.DEBUG)
	logger.SetLevel(log.ERROR)
	// No direct way to verify level was set, but method should not panic
}

func TestDefaultFileLogger_LogMethods(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	logger, err := log.NewDefaultFileLogger(logPath, log.DEBUG)
	require.NoError(t, err)
	defer func() { _ = logger.Close() }()

	// Test all log methods
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Test formatted log methods
	logger.Debugf("debug %s", "formatted")
	logger.Infof("info %d", 42)
	logger.Warnf("warn %v", true)
	logger.Errorf("error %s %d", "formatted", 123)

	// Verify log file was created and contains content
	content, err := os.ReadFile(logPath) //nolint:gosec // G304: Test needs to read temp file with variable path
	require.NoError(t, err)

	logContent := string(content)
	assert.Contains(t, logContent, "[DEBUG] debug message")
	assert.Contains(t, logContent, "[INFO] info message")
	assert.Contains(t, logContent, "[WARN] warn message")
	assert.Contains(t, logContent, "[ERROR] error message")
	assert.Contains(t, logContent, "[DEBUG] debug formatted")
	assert.Contains(t, logContent, "[INFO] info 42")
	assert.Contains(t, logContent, "[WARN] warn true")
	assert.Contains(t, logContent, "[ERROR] error formatted 123")
}

func TestDefaultFileLogger_LogLevelFiltering(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Create logger with WARN level (should filter out DEBUG and INFO)
	logger, err := log.NewDefaultFileLogger(logPath, log.WARN)
	require.NoError(t, err)
	defer func() { _ = logger.Close() }()

	// Log messages at different levels
	logger.Debug("debug message - should be filtered")
	logger.Info("info message - should be filtered")
	logger.Warn("warn message - should appear")
	logger.Error("error message - should appear")

	// Verify only WARN and ERROR messages appear
	content, err := os.ReadFile(logPath) //nolint:gosec // G304: Test needs to read temp file with variable path
	require.NoError(t, err)

	logContent := string(content)
	assert.NotContains(t, logContent, "[DEBUG] debug message")
	assert.NotContains(t, logContent, "[INFO] info message")
	assert.Contains(t, logContent, "[WARN] warn message")
	assert.Contains(t, logContent, "[ERROR] error message")
}

func TestDefaultFileLogger_LogLevelFilteringAfterSetLevel(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	// Create logger with DEBUG level initially
	logger, err := log.NewDefaultFileLogger(logPath, log.DEBUG)
	require.NoError(t, err)
	defer func() { _ = logger.Close() }()

	// Log a debug message (should appear)
	logger.Debug("debug message 1")

	// Change level to ERROR
	logger.SetLevel(log.ERROR)

	// Log messages at different levels
	logger.Debug("debug message 2 - should be filtered")
	logger.Info("info message - should be filtered")
	logger.Warn("warn message - should be filtered")
	logger.Error("error message - should appear")

	// Verify filtering behavior changed
	content, err := os.ReadFile(logPath) //nolint:gosec // G304: Test needs to read temp file with variable path
	require.NoError(t, err)

	logContent := string(content)
	assert.Contains(t, logContent, "[DEBUG] debug message 1")    // Before level change
	assert.NotContains(t, logContent, "[DEBUG] debug message 2") // After level change
	assert.NotContains(t, logContent, "[INFO] info message")
	assert.NotContains(t, logContent, "[WARN] warn message")
	assert.Contains(t, logContent, "[ERROR] error message")
}

func TestDefaultFileLogger_ConcurrentLogging(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "concurrent.log")

	logger, err := log.NewDefaultFileLogger(logPath, log.DEBUG)
	require.NoError(t, err)
	defer func() { _ = logger.Close() }()

	// Test concurrent logging
	const numGoroutines = 10
	const messagesPerGoroutine = 5

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer func() { done <- true }()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Infof("Worker %d message %d", workerID, j)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all messages were logged
	content, err := os.ReadFile(logPath) //nolint:gosec // G304: Test needs to read temp file with variable path
	require.NoError(t, err)

	logContent := string(content)
	messageCount := strings.Count(logContent, "[INFO] Worker")
	assert.Equal(t, numGoroutines*messagesPerGoroutine, messageCount)
}

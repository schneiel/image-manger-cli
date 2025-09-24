package shared

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// MockFileSystem implements filesystem.FileSystem for testing.
type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) Create(name string) (filesystem.File, error) {
	args := m.Called(name)
	return args.Get(0).(filesystem.File), args.Error(1)
}

// Implement other required methods with minimal functionality.
func (m *MockFileSystem) Open(_ string) (filesystem.File, error) { return nil, nil }

func (m *MockFileSystem) OpenFile(_ string, _ int, _ os.FileMode) (filesystem.File, error) {
	return nil, nil
}
func (m *MockFileSystem) Remove(_ string) error                             { return nil }
func (m *MockFileSystem) RemoveAll(_ string) error                          { return nil }
func (m *MockFileSystem) Rename(_, _ string) error                          { return nil }
func (m *MockFileSystem) Mkdir(_ string, _ os.FileMode) error               { return nil }
func (m *MockFileSystem) MkdirAll(_ string, _ os.FileMode) error            { return nil }
func (m *MockFileSystem) ReadDir(_ string) ([]os.DirEntry, error)           { return nil, nil }
func (m *MockFileSystem) Stat(_ string) (os.FileInfo, error)                { return nil, nil }
func (m *MockFileSystem) Lstat(_ string) (os.FileInfo, error)               { return nil, nil }
func (m *MockFileSystem) ReadFile(_ string) ([]byte, error)                 { return nil, nil }
func (m *MockFileSystem) WriteFile(_ string, _ []byte, _ os.FileMode) error { return nil }
func (m *MockFileSystem) CreateTemp(_, _ string) (filesystem.File, error)   { return nil, nil }
func (m *MockFileSystem) MkdirTemp(_, _ string) (string, error)             { return "", nil }
func (m *MockFileSystem) Chmod(_ string, _ os.FileMode) error               { return nil }
func (m *MockFileSystem) Chown(_ string, _, _ int) error                    { return nil }
func (m *MockFileSystem) Chtimes(_ string, _, _ time.Time) error            { return nil }
func (m *MockFileSystem) Getwd() (string, error)                            { return "/test", nil }
func (m *MockFileSystem) Chdir(_ string) error                              { return nil }
func (m *MockFileSystem) Symlink(_, _ string) error                         { return nil }
func (m *MockFileSystem) Link(_, _ string) error                            { return nil }
func (m *MockFileSystem) Readlink(_ string) (string, error)                 { return "", nil }
func (m *MockFileSystem) WalkDir(_ string, _ fs.WalkDirFunc) error          { return nil }
func (m *MockFileSystem) IsNotExist(_ error) bool                           { return false }

// MockFile implements filesystem.File for testing.
type MockFile struct {
	mock.Mock
	name string
}

func (m *MockFile) Read(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockFile) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockFile) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Implement other required methods with minimal functionality.
func (m *MockFile) Seek(_ int64, _ int) (int64, error)   { return 0, nil }
func (m *MockFile) Stat() (os.FileInfo, error)           { return nil, nil }
func (m *MockFile) Sync() error                          { return nil }
func (m *MockFile) Truncate(_ int64) error               { return nil }
func (m *MockFile) Name() string                         { return m.name }
func (m *MockFile) Readdir(_ int) ([]os.FileInfo, error) { return nil, nil }
func (m *MockFile) Readdirnames(_ int) ([]string, error) { return nil, nil }

// MockLogger implements log.Logger for testing.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) SetLevel(level log.Level) {
	m.Called(level)
}

func (m *MockLogger) Debug(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Info(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Warn(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Error(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Debugf(format string, v ...interface{}) {
	m.Called(format, v)
}

func (m *MockLogger) Infof(format string, v ...interface{}) {
	m.Called(format, v)
}

func (m *MockLogger) Warnf(format string, v ...interface{}) {
	m.Called(format, v)
}

func (m *MockLogger) Errorf(format string, v ...interface{}) {
	m.Called(format, v)
}

// MockLocalizer implements i18n.Localizer for testing.
type MockLocalizer struct {
	mock.Mock
}

// Ensure MockLocalizer implements i18n.Localizer.
var _ i18n.Localizer = (*MockLocalizer)(nil)

func (m *MockLocalizer) Translate(messageID string, templateData ...map[string]interface{}) string {
	args := m.Called(messageID, templateData)
	return args.String(0)
}

func (m *MockLocalizer) GetCurrentLanguage() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLocalizer) SetLanguage(lang string) error {
	args := m.Called(lang)
	return args.Error(0)
}

func (m *MockLocalizer) IsInitialized() bool {
	args := m.Called()
	return args.Bool(0)
}

// TestNewCSVResource tests the creation of a CSV resource.
func TestNewCSVResource(t *testing.T) {
	mockLogger := &MockLogger{}
	mockLocalizer := &MockLocalizer{}

	resource := NewCSVResource("test.csv", []string{"header1", "header2"}, mockLogger, mockLocalizer)

	assert.NotNil(t, resource)
}

// TestNewCSVResourceWithFilesystem tests creation with custom filesystem.
func TestNewCSVResourceWithFilesystem(t *testing.T) {
	mockLogger := &MockLogger{}
	mockLocalizer := &MockLocalizer{}
	mockFS := &MockFileSystem{}

	resource := NewCSVResourceWithFilesystem(
		"test.csv",
		[]string{"header1", "header2"},
		mockLogger,
		mockLocalizer,
		mockFS,
	)

	assert.NotNil(t, resource)
}

// TestNewCSVResourceWithFilesystem_NilFilesystem_Panics tests panic on nil filesystem.
func TestNewCSVResourceWithFilesystem_NilFilesystem_Panics(t *testing.T) {
	mockLogger := &MockLogger{}
	mockLocalizer := &MockLocalizer{}

	assert.Panics(t, func() {
		NewCSVResourceWithFilesystem("test.csv", []string{"header1", "header2"}, mockLogger, mockLocalizer, nil)
	})
}

// TestDefaultCSVResource_Setup_Success tests successful setup.
func TestDefaultCSVResource_Setup_Success(t *testing.T) {
	mockLogger := &MockLogger{}
	mockFS := &MockFileSystem{}
	mockFile := &MockFile{name: "test.csv"}

	mockFS.On("Create", "test.csv").Return(mockFile, nil)
	mockFile.On("Write", mock.AnythingOfType("[]uint8")).Return(16, nil)

	resource := NewCSVResourceWithFilesystem(
		"test.csv",
		[]string{"header1", "header2"},
		mockLogger,
		&MockLocalizer{},
		mockFS,
	)

	err := resource.Setup()

	require.NoError(t, err)
	mockFS.AssertExpectations(t)
	mockFile.AssertExpectations(t)

	// Adding test for nil logger
	resource = NewCSVResourceWithFilesystem("test.csv", []string{"header1", "header2"}, nil, &MockLocalizer{}, mockFS)
	err = resource.Setup()
	require.NoError(t, err)
}

// TestDefaultCSVResource_Setup_Error tests setup failure.
func TestDefaultCSVResource_Setup_Error(t *testing.T) {
	mockLogger := &MockLogger{}
	mockFS := &MockFileSystem{}

	mockFS.On("Create", "test.csv").Return((*MockFile)(nil), errors.New("create failed"))

	resource := NewCSVResourceWithFilesystem(
		"test.csv",
		[]string{"header1", "header2"},
		mockLogger,
		&MockLocalizer{},
		mockFS,
	)

	err := resource.Setup()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "create failed")
	mockFS.AssertExpectations(t)
}

// TestDefaultCSVResource_WriteRow_NotInitialized tests writing without setup.
func TestDefaultCSVResource_WriteRow_NotInitialized(t *testing.T) {
	mockLogger := &MockLogger{}
	mockFS := &MockFileSystem{}

	resource := NewCSVResourceWithFilesystem(
		"test.csv",
		[]string{"header1", "header2"},
		mockLogger,
		&MockLocalizer{},
		mockFS,
	)

	// Try to write without setup
	row := []string{"value1", "value2"}
	err := resource.WriteRow(row)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

// TestDefaultCSVResource_Teardown_NoSetup tests teardown without setup.
func TestDefaultCSVResource_Teardown_NoSetup(t *testing.T) {
	mockLogger := &MockLogger{}
	mockFS := &MockFileSystem{}

	resource := NewCSVResourceWithFilesystem(
		"test.csv",
		[]string{"header1", "header2"},
		mockLogger,
		&MockLocalizer{},
		mockFS,
	)

	// Teardown without setup should not error
	err := resource.Teardown()

	require.NoError(t, err)
}

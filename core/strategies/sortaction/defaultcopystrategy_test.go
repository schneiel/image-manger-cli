package sortaction

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultCopyStrategy(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	strategy, err := NewDefaultCopyStrategy(mockLogger, mockLocalizer)
	require.NoError(t, err)

	if strategy == nil {
		t.Fatal("Expected strategy to be created, got nil")
	}

	if strategy.logger != mockLogger {
		t.Error("Expected logger to be injected")
	}
	if strategy.localizer != mockLocalizer {
		t.Error("Expected localizer to be injected")
	}
	if strategy.fileSystem == nil {
		t.Error("Expected fileSystem to be created")
	}
}

func TestNewDefaultCopyStrategyWithFilesystem(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileSystem := &testutils.MockFileSystem{}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)

	if strategy == nil {
		t.Fatal("Expected strategy to be created, got nil")
	}

	if strategy.logger != mockLogger {
		t.Error("Expected logger to be injected")
	}
	if strategy.localizer != mockLocalizer {
		t.Error("Expected localizer to be injected")
	}
	if strategy.fileSystem != mockFileSystem {
		t.Error("Expected fileSystem to be injected")
	}
}

func TestNewDefaultCopyStrategyWithFilesystem_NilFilesystem(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	_, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fileSystem cannot be nil")
}

func TestDefaultCopyStrategy_Execute_Success(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "FileCopied"
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}

	// Mock directory creation
	mockFileSystem.MkdirAllFunc = func(_ string, _ os.FileMode) error {
		return nil
	}

	// Mock file operations
	sourceContent := []byte("test file content")
	mockSourceFile := &testutils.MockFile{
		ReadFunc: func(p []byte) (int, error) {
			copy(p, sourceContent)
			return len(sourceContent), io.EOF
		},
		CloseFunc: func() error {
			return nil
		},
	}
	mockDestFile := &testutils.MockFile{
		WriteFunc: func(p []byte) (int, error) {
			return len(p), nil
		},
		CloseFunc: func() error {
			return nil
		},
	}

	mockFileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockSourceFile, nil
	}
	mockFileSystem.CreateFunc = func(_ string) (filesystem.File, error) {
		return mockDestFile, nil
	}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)
	source := "/path/to/source/image.jpg"
	destination := "/path/to/destination/image.jpg"

	err = strategy.Execute(source, destination)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDefaultCopyStrategy_Execute_MkdirAllError(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "FileCopied"
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}

	// Mock directory creation error
	mkdirError := errors.New("mkdir error")
	mockFileSystem.MkdirAllFunc = func(_ string, _ os.FileMode) error {
		return mkdirError
	}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)
	source := "/path/to/source/image.jpg"
	destination := "/path/to/destination/image.jpg"

	err = strategy.Execute(source, destination)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "failed to create destination directory: mkdir error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDefaultCopyStrategy_Execute_OpenSourceError(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "FileCopied"
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}

	// Mock directory creation success
	mockFileSystem.MkdirAllFunc = func(_ string, _ os.FileMode) error {
		return nil
	}

	// Mock source file open error
	openError := errors.New("open error")
	mockFileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return nil, openError
	}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)
	source := "/path/to/source/image.jpg"
	destination := "/path/to/destination/image.jpg"

	err = strategy.Execute(source, destination)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "failed to open source file: open error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDefaultCopyStrategy_Execute_CreateDestError(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "FileCopied"
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}

	// Mock directory creation success
	mockFileSystem.MkdirAllFunc = func(_ string, _ os.FileMode) error {
		return nil
	}

	// Mock source file open success
	mockSourceFile := &testutils.MockFile{
		CloseFunc: func() error {
			return nil
		},
	}
	mockFileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockSourceFile, nil
	}

	// Mock destination file create error
	createError := errors.New("create error")
	mockFileSystem.CreateFunc = func(_ string) (filesystem.File, error) {
		return nil, createError
	}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)
	source := "/path/to/source/image.jpg"
	destination := "/path/to/destination/image.jpg"

	err = strategy.Execute(source, destination)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "failed to create destination file: create error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDefaultCopyStrategy_Execute_CopyError(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "FileCopied"
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}

	// Mock directory creation success
	mockFileSystem.MkdirAllFunc = func(_ string, _ os.FileMode) error {
		return nil
	}

	// Mock source file open success
	mockSourceFile := &testutils.MockFile{
		ReadFunc: func(_ []byte) (int, error) {
			return 0, errors.New("read error")
		},
		CloseFunc: func() error {
			return nil
		},
	}
	mockFileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockSourceFile, nil
	}

	// Mock destination file create success
	mockDestFile := &testutils.MockFile{
		CloseFunc: func() error {
			return nil
		},
	}
	mockFileSystem.CreateFunc = func(_ string) (filesystem.File, error) {
		return mockDestFile, nil
	}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)
	source := "/path/to/source/image.jpg"
	destination := "/path/to/destination/image.jpg"

	err = strategy.Execute(source, destination)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "failed to copy file: read error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDefaultCopyStrategy_Execute_ComplexPaths(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "FileCopied"
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}

	// Mock all operations to succeed
	mockFileSystem.MkdirAllFunc = func(_ string, _ os.FileMode) error {
		return nil
	}
	mockSourceFile := &testutils.MockFile{
		ReadFunc: func(_ []byte) (int, error) {
			return 0, io.EOF
		},
		CloseFunc: func() error {
			return nil
		},
	}
	mockDestFile := &testutils.MockFile{
		WriteFunc: func(p []byte) (int, error) {
			return len(p), nil
		},
		CloseFunc: func() error {
			return nil
		},
	}
	mockFileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockSourceFile, nil
	}
	mockFileSystem.CreateFunc = func(_ string) (filesystem.File, error) {
		return mockDestFile, nil
	}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)

	testCases := []struct {
		source      string
		destination string
	}{
		{"/path/with spaces/file name.jpg", "/destination/with spaces/file name.jpg"},
		{"/path/with/unicode/测试图片.jpeg", "/destination/with/unicode/测试图片.jpeg"},
		{"/path/with/dots/file.name.with.dots.tiff", "/destination/with/dots/file.name.with.dots.tiff"},
		{"/path/with/numbers/IMG_2023_001.jpg", "/destination/with/numbers/IMG_2023_001.jpg"},
		{"/path/with/special/chars/file@#$%^&*().png", "/destination/with/special/chars/file@#$%^&*().png"},
	}

	for _, tc := range testCases {
		t.Run(tc.source, func(_ *testing.T) {
			err := strategy.Execute(tc.source, tc.destination)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}

func TestDefaultCopyStrategy_GetResources(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileSystem := &testutils.MockFileSystem{}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)

	resources := strategy.GetResources()

	if resources != nil {
		t.Errorf("Expected nil resources, got %+v", resources)
	}
}

func TestDefaultCopyStrategy_Execute_LoggingVerification(t *testing.T) {
	t.Parallel()
	loggedMessages := []string{}
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(format string, _ ...interface{}) {
			loggedMessages = append(loggedMessages, format)
		},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "FileCopied"
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}

	// Mock all operations to succeed
	mockFileSystem.MkdirAllFunc = func(_ string, _ os.FileMode) error {
		return nil
	}
	mockSourceFile := &testutils.MockFile{
		ReadFunc: func(_ []byte) (int, error) {
			return 0, io.EOF
		},
		CloseFunc: func() error {
			return nil
		},
	}
	mockDestFile := &testutils.MockFile{
		WriteFunc: func(p []byte) (int, error) {
			return len(p), nil
		},
		CloseFunc: func() error {
			return nil
		},
	}
	mockFileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockSourceFile, nil
	}
	mockFileSystem.CreateFunc = func(_ string) (filesystem.File, error) {
		return mockDestFile, nil
	}

	strategy, err := NewDefaultCopyStrategyWithFilesystem(mockLogger, mockLocalizer, mockFileSystem)
	require.NoError(t, err)
	source := "/path/to/source/image.jpg"
	destination := "/path/to/destination/image.jpg"

	err = strategy.Execute(source, destination)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(loggedMessages) != 1 {
		t.Errorf("Expected 1 logged message, got %d", len(loggedMessages))
	}

	if loggedMessages[0] != "FileCopied" {
		t.Errorf("Expected logged message 'FileCopied', got '%s'", loggedMessages[0])
	}
}

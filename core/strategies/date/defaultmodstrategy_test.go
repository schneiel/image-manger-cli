package date

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultModTimeStrategy(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultModTimeStrategy()

	if strategy == nil {
		t.Fatal("Expected strategy to be created, got nil")
	}

	if strategy.fileSystem == nil {
		t.Error("Expected filesystem to be initialized")
	}
}

func TestNewDefaultModTimeStrategyWithFilesystem(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}

	strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)

	if strategy == nil {
		t.Fatal("Expected strategy to be created, got nil")
	}

	if strategy.fileSystem != mockFileSystem {
		t.Error("Expected injected filesystem to be used")
	}
}

func TestNewDefaultModTimeStrategyWithFilesystem_NilFilesystem(t *testing.T) {
	t.Parallel()
	// This should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when filesystem is nil")
		}
	}()

	NewDefaultModTimeStrategyWithFilesystem(nil)
}

func TestDefaultModTimeStrategy_Extract_Success(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestDefaultModTimeStrategy_Extract_StatError(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}

	expectedError := errors.New("file not found")
	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return nil, expectedError
	}

	strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error '%s', got '%s'", expectedError.Error(), err.Error())
	}

	// Should return zero time when stat fails
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultModTimeStrategy_Extract_NilFields(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(nil, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestDefaultModTimeStrategy_Extract_EmptyFields(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
	fields := map[string]interface{}{}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestDefaultModTimeStrategy_Extract_ComplexFilePath(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}

	testCases := []string{
		"/path/to/test/image.jpg",
		"/path/with spaces/file name.png",
		"/path/with/unicode/测试图片.jpeg",
		"/path/with/dots/file.name.with.dots.tiff",
		"/path/with/numbers/IMG_2023_001.jpg",
		"/path/with/special/chars/file@#$%^&*().png",
	}

	for _, filePath := range testCases {
		t.Run(filePath, func(_ *testing.T) {
			result, err := strategy.Extract(fields, filePath)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !result.Equal(expectedTime) {
				t.Errorf("Expected %v, got %v", expectedTime, result)
			}
		})
	}
}

func TestDefaultModTimeStrategy_Extract_VariousTimes(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}

	testCases := []struct {
		name     string
		modTime  time.Time
		expected time.Time
	}{
		{
			name:     "Recent time",
			modTime:  time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:     "Old time",
			modTime:  time.Date(2020, 12, 25, 23, 59, 59, 0, time.UTC),
			expected: time.Date(2020, 12, 25, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "Future time",
			modTime:  time.Date(2030, 6, 15, 14, 20, 0, 0, time.UTC),
			expected: time.Date(2030, 6, 15, 14, 20, 0, 0, time.UTC),
		},
		{
			name:     "Zero time",
			modTime:  time.Time{},
			expected: time.Time{},
		},
		{
			name:     "Unix epoch",
			modTime:  time.Unix(0, 0),
			expected: time.Unix(0, 0),
		},
		{
			name:     "With nanoseconds",
			modTime:  time.Date(2023, 1, 15, 12, 30, 0, 123456789, time.UTC),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 123456789, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			mockFileInfo := &testutils.MockFileInfo{}
			mockFileInfo.ModTimeFunc = func() time.Time {
				return tc.modTime
			}

			mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
				return mockFileInfo, nil
			}

			strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
			fields := map[string]interface{}{
				"DateTimeOriginal": "2023:01:15 12:30:00",
			}
			filePath := "/path/to/test/image.jpg"

			result, err := strategy.Extract(fields, filePath)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultModTimeStrategy_Extract_TimezoneHandling(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockFileInfo := &testutils.MockFileInfo{}

	testCases := []struct {
		name     string
		modTime  time.Time
		expected time.Time
	}{
		{
			name:     "UTC time",
			modTime:  time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:     "EST time",
			modTime:  time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("EST", -5*3600)),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("EST", -5*3600)),
		},
		{
			name:     "PST time",
			modTime:  time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("PST", -8*3600)),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("PST", -8*3600)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			mockFileInfo.ModTimeFunc = func() time.Time {
				return tc.modTime
			}

			mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
				return mockFileInfo, nil
			}

			strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
			fields := map[string]interface{}{
				"DateTimeOriginal": "2023:01:15 12:30:00",
			}
			filePath := "/path/to/test/image.jpg"

			result, err := strategy.Extract(fields, filePath)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultModTimeStrategy_Extract_FileSystemErrors(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}

	testCases := []struct {
		name          string
		statError     error
		expectedError string
	}{
		{
			name:          "File not found",
			statError:     errors.New("file not found"),
			expectedError: "file not found",
		},
		{
			name:          "Permission denied",
			statError:     errors.New("permission denied"),
			expectedError: "permission denied",
		},
		{
			name:          "Path too long",
			statError:     errors.New("path too long"),
			expectedError: "path too long",
		},
		{
			name:          "Network error",
			statError:     errors.New("network unreachable"),
			expectedError: "network unreachable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
				return nil, tc.statError
			}

			strategy := NewDefaultModTimeStrategyWithFilesystem(mockFileSystem)
			fields := map[string]interface{}{
				"DateTimeOriginal": "2023:01:15 12:30:00",
			}
			filePath := "/path/to/test/image.jpg"

			result, err := strategy.Extract(fields, filePath)

			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if err.Error() != tc.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tc.expectedError, err.Error())
			}

			// Should return zero time when stat fails
			zeroTime := time.Time{}
			if !result.Equal(zeroTime) {
				t.Errorf("Expected zero time, got %v", result)
			}
		})
	}
}

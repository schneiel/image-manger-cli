package date

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultModificationTimeStrategy(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

	extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)

	if extractor == nil {
		t.Fatal("Expected extractor to be created, got nil")
	}

	if extractor.fileSystem != mockFileSystem {
		t.Error("Expected filesystem to be injected")
	}

	if extractor.localizer != mockLocalizer {
		t.Error("Expected localizer to be injected")
	}
}

func TestNewDefaultModificationTimeStrategy_NilFilesystem(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}

	_, err := NewDefaultModificationTimeStrategy(nil, mockLocalizer)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fileSystem cannot be nil")
}

func TestNewDefaultModificationTimeStrategy_NilLocalizer(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}

	_, err := NewDefaultModificationTimeStrategy(mockFileSystem, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "localizer cannot be nil")
}

func TestDefaultModificationTimeStrategy_Extract_Success(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
	filePath := "/path/to/test/image.jpg"
	exifData := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}

	result, success := extractor.Extract(filePath, exifData)

	if !success {
		t.Fatal("Expected success, got failure")
	}

	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestDefaultModificationTimeStrategy_Extract_StatError(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

	expectedError := errors.New("file not found")
	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return nil, expectedError
	}

	extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
	filePath := "/path/to/test/image.jpg"
	exifData := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}

	result, success := extractor.Extract(filePath, exifData)

	if success {
		t.Fatal("Expected failure, got success")
	}

	// Should return zero time when stat fails
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultModificationTimeStrategy_Extract_NilExifData(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
	filePath := "/path/to/test/image.jpg"

	result, success := extractor.Extract(filePath, nil)

	if !success {
		t.Fatal("Expected success, got failure")
	}

	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestDefaultModificationTimeStrategy_Extract_EmptyExifData(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
	filePath := "/path/to/test/image.jpg"
	exifData := map[string]interface{}{}

	result, success := extractor.Extract(filePath, exifData)

	if !success {
		t.Fatal("Expected success, got failure")
	}

	if !result.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, result)
	}
}

func TestDefaultModificationTimeStrategy_Extract_ComplexFilePath(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileInfo := &testutils.MockFileInfo{}

	expectedTime := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockFileInfo.ModTimeFunc = func() time.Time {
		return expectedTime
	}

	mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
		return mockFileInfo, nil
	}

	extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
	exifData := map[string]interface{}{
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
			result, success := extractor.Extract(filePath, exifData)

			if !success {
				t.Fatal("Expected success, got failure")
			}

			if !result.Equal(expectedTime) {
				t.Errorf("Expected %v, got %v", expectedTime, result)
			}
		})
	}
}

func TestDefaultModificationTimeStrategy_Extract_VariousTimes(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

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

			extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
			filePath := "/path/to/test/image.jpg"
			exifData := map[string]interface{}{
				"DateTimeOriginal": "2023:01:15 12:30:00",
			}

			result, success := extractor.Extract(filePath, exifData)

			if !success {
				t.Fatal("Expected success, got failure")
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultModificationTimeStrategy_Extract_TimezoneHandling(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

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
			mockFileInfo := &testutils.MockFileInfo{}
			mockFileInfo.ModTimeFunc = func() time.Time {
				return tc.modTime
			}

			mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
				return mockFileInfo, nil
			}

			extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
			filePath := "/path/to/test/image.jpg"
			exifData := map[string]interface{}{
				"DateTimeOriginal": "2023:01:15 12:30:00",
			}

			result, success := extractor.Extract(filePath, exifData)

			if !success {
				t.Fatal("Expected success, got failure")
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultModificationTimeStrategy_Extract_FileSystemErrors(t *testing.T) {
	t.Parallel()
	mockFileSystem := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

	testCases := []struct {
		name      string
		statError error
	}{
		{
			name:      "File not found",
			statError: errors.New("file not found"),
		},
		{
			name:      "Permission denied",
			statError: errors.New("permission denied"),
		},
		{
			name:      "Path too long",
			statError: errors.New("path too long"),
		},
		{
			name:      "Network error",
			statError: errors.New("network unreachable"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			mockFileSystem.StatFunc = func(_ string) (os.FileInfo, error) {
				return nil, tc.statError
			}

			extractor, err := NewDefaultModificationTimeStrategy(mockFileSystem, mockLocalizer)
	require.NoError(t, err)
			filePath := "/path/to/test/image.jpg"
			exifData := map[string]interface{}{
				"DateTimeOriginal": "2023:01:15 12:30:00",
			}

			result, success := extractor.Extract(filePath, exifData)

			if success {
				t.Fatal("Expected failure, got success")
			}

			// Should return zero time when stat fails
			zeroTime := time.Time{}
			if !result.Equal(zeroTime) {
				t.Errorf("Expected zero time, got %v", result)
			}
		})
	}
}

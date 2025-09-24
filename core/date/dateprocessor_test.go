package date

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewStrategyProcessor(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	processor := NewStrategyProcessor(mockStrategy)

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.strategy != mockStrategy {
		t.Error("Expected strategy to be injected")
	}
}

func TestStrategyProcessor_GetBestAvailableDate_Success(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}

	processor := NewStrategyProcessor(mockStrategy)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := processor.GetBestAvailableDate(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return local time
	expectedLocal := expectedDate.Local()
	if !result.Equal(expectedLocal) {
		t.Errorf("Expected %v, got %v", expectedLocal, result)
	}
}

func TestStrategyProcessor_GetBestAvailableDate_StrategyError(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	expectedError := errors.New("strategy error")
	mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, expectedError
	}

	processor := NewStrategyProcessor(mockStrategy)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := processor.GetBestAvailableDate(fields, filePath)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error contains the expected components
	if !strings.Contains(err.Error(), "failed to extract date from") {
		t.Errorf("Expected error to contain 'failed to extract date from', got '%s'", err.Error())
	}
	if !strings.Contains(err.Error(), filePath) {
		t.Errorf("Expected error to contain file path '%s', got '%s'", filePath, err.Error())
	}
	if !strings.Contains(err.Error(), expectedError.Error()) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError.Error(), err.Error())
	}

	// Should return zero time
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestStrategyProcessor_GetBestAvailableDate_ZeroDate(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	// Strategy returns zero time (no error)
	mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, nil
	}

	processor := NewStrategyProcessor(mockStrategy)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := processor.GetBestAvailableDate(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return zero time
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestStrategyProcessor_GetBestAvailableDate_NilFields(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}

	processor := NewStrategyProcessor(mockStrategy)
	filePath := "/path/to/test/image.jpg"

	result, err := processor.GetBestAvailableDate(nil, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return local time
	expectedLocal := expectedDate.Local()
	if !result.Equal(expectedLocal) {
		t.Errorf("Expected %v, got %v", expectedLocal, result)
	}
}

func TestStrategyProcessor_GetBestAvailableDate_EmptyFields(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}

	processor := NewStrategyProcessor(mockStrategy)
	fields := map[string]interface{}{}
	filePath := "/path/to/test/image.jpg"

	result, err := processor.GetBestAvailableDate(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return local time
	expectedLocal := expectedDate.Local()
	if !result.Equal(expectedLocal) {
		t.Errorf("Expected %v, got %v", expectedLocal, result)
	}
}

func TestStrategyProcessor_GetBestAvailableDate_ComplexFields(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}

	processor := NewStrategyProcessor(mockStrategy)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
		"DateTime":         "2023:01:15 12:30:00",
		"CreateDate":       "2023:01:15 12:30:00",
		"ModifyDate":       "2023:01:15 12:30:00",
		"GPSDateStamp":     "2023:01:15",
		"GPSTimeStamp":     "12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := processor.GetBestAvailableDate(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return local time
	expectedLocal := expectedDate.Local()
	if !result.Equal(expectedLocal) {
		t.Errorf("Expected %v, got %v", expectedLocal, result)
	}
}

func TestStrategyProcessor_GetBestAvailableDate_ComplexFilePath(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}

	processor := NewStrategyProcessor(mockStrategy)
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
		t.Run(filePath, func(subT *testing.T) {
			subT.Parallel()
			result, err := processor.GetBestAvailableDate(fields, filePath)
			if err != nil {
				subT.Fatalf("Expected no error, got %v", err)
			}

			// Should return local time
			expectedLocal := expectedDate.Local()
			if !result.Equal(expectedLocal) {
				subT.Errorf("Expected %v, got %v", expectedLocal, result)
			}
		})
	}
}

func TestStrategyProcessor_GetBestAvailableDate_TimezoneHandling(t *testing.T) {
	t.Parallel()
	mockStrategy := &testutils.MockDateStrategy{}

	// Test with different timezones
	testCases := []struct {
		name     string
		utcTime  time.Time
		expected time.Time
	}{
		{
			name:     "UTC time",
			utcTime:  time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC).Local(),
		},
		{
			name:     "EST time",
			utcTime:  time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("EST", -5*3600)),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("EST", -5*3600)).Local(),
		},
		{
			name:     "PST time",
			utcTime:  time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("PST", -8*3600)),
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.FixedZone("PST", -8*3600)).Local(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(subT *testing.T) {
			mockStrategy.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
				return tc.utcTime, nil
			}

			processor := NewStrategyProcessor(mockStrategy)
			fields := map[string]interface{}{
				"DateTimeOriginal": "2023:01:15 12:30:00",
			}
			filePath := "/path/to/test/image.jpg"

			result, err := processor.GetBestAvailableDate(fields, filePath)
			if err != nil {
				subT.Fatalf("Expected no error, got %v", err)
			}

			// Should return local time
			if !result.Equal(tc.expected) {
				subT.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

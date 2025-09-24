package date

import (
	"testing"
	"time"
)

func TestNewExtractorAdapter(t *testing.T) {
	t.Parallel()
	extractor := func(_ string, _ map[string]interface{}) (time.Time, bool) {
		return time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC), true
	}

	adapter := NewExtractorAdapter(extractor)

	if adapter == nil {
		t.Fatal("Expected adapter to be created, got nil")
	}

	// Test that it implements DateStrategy interface
	_ = adapter
}

func TestDateExtractorAdapter_Extract_Success(t *testing.T) {
	t.Parallel()
	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	extractor := func(_ string, _ map[string]interface{}) (time.Time, bool) {
		return expectedDate, true
	}

	adapter := NewExtractorAdapter(extractor)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := adapter.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDateExtractorAdapter_Extract_ExtractorReturnsFalse(t *testing.T) {
	t.Parallel()
	extractor := func(_ string, _ map[string]interface{}) (time.Time, bool) {
		return time.Time{}, false
	}

	adapter := NewExtractorAdapter(extractor)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := adapter.Extract(fields, filePath)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "date not found by extractor"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}

	// Should return zero time
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDateExtractorAdapter_Extract_NilFields(t *testing.T) {
	t.Parallel()
	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	extractor := func(_ string, _ map[string]interface{}) (time.Time, bool) {
		return expectedDate, true
	}

	adapter := NewExtractorAdapter(extractor)
	filePath := "/path/to/test/image.jpg"

	result, err := adapter.Extract(nil, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDateExtractorAdapter_Extract_EmptyFields(t *testing.T) {
	t.Parallel()
	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	extractor := func(_ string, _ map[string]interface{}) (time.Time, bool) {
		return expectedDate, true
	}

	adapter := NewExtractorAdapter(extractor)
	fields := map[string]interface{}{}
	filePath := "/path/to/test/image.jpg"

	result, err := adapter.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDateExtractorAdapter_Extract_ComplexFilePath(t *testing.T) {
	t.Parallel()
	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	extractor := func(_ string, _ map[string]interface{}) (time.Time, bool) {
		return expectedDate, true
	}

	adapter := NewExtractorAdapter(extractor)
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
			result, err := adapter.Extract(fields, filePath)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !result.Equal(expectedDate) {
				t.Errorf("Expected %v, got %v", expectedDate, result)
			}
		})
	}
}

func TestDateExtractorAdapter_Extract_ExtractorUsesParameters(t *testing.T) {
	t.Parallel()
	extractor := func(filePath string, exifData map[string]interface{}) (time.Time, bool) {
		// Check that the extractor receives the correct parameters
		if filePath != "/path/to/test/image.jpg" {
			return time.Time{}, false
		}
		if exifData == nil {
			return time.Time{}, false
		}
		if exifData["DateTimeOriginal"] != "2023:01:15 12:30:00" {
			return time.Time{}, false
		}
		return time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC), true
	}

	adapter := NewExtractorAdapter(extractor)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := adapter.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDateExtractorAdapter_Extract_ExtractorReturnsZeroTime(t *testing.T) {
	t.Parallel()
	extractor := func(_ string, _ map[string]interface{}) (time.Time, bool) {
		return time.Time{}, true // Returns zero time but success
	}

	adapter := NewExtractorAdapter(extractor)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := adapter.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return zero time
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

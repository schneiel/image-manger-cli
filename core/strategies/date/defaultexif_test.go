package date

import (
	"testing"
	"time"
)

func TestNewDefaultExifDateStrategy(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	if strategy == nil {
		t.Fatal("Expected strategy to be created, got nil")
	}
}

func TestDefaultExifDateStrategy_Extract_Success(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDefaultExifDateStrategy_Extract_NilFields(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(nil, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return zero time when fields is nil
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultExifDateStrategy_Extract_EmptyFields(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	fields := map[string]interface{}{}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return zero time when DateTimeOriginal is not present
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultExifDateStrategy_Extract_DateTimeOriginalNotString(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	fields := map[string]interface{}{
		"DateTimeOriginal": 12345, // Not a string
	}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return zero time when DateTimeOriginal is not a string
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultExifDateStrategy_Extract_InvalidDateFormat(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	fields := map[string]interface{}{
		"DateTimeOriginal": "invalid-date-format",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)

	if err == nil {
		t.Fatal("Expected error for invalid date format, got nil")
	}

	// Should return zero time when date parsing fails
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultExifDateStrategy_Extract_VariousDateFormats(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	testCases := []struct {
		name        string
		dateStr     string
		expected    time.Time
		shouldError bool
	}{
		{
			name:        "Standard format",
			dateStr:     "2023:01:15 12:30:00",
			expected:    time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
			shouldError: false,
		},
		{
			name:        "Different year",
			dateStr:     "2020:12:25 23:59:59",
			expected:    time.Date(2020, 12, 25, 23, 59, 59, 0, time.UTC),
			shouldError: false,
		},
		{
			name:        "Leap year",
			dateStr:     "2024:02:29 14:30:00",
			expected:    time.Date(2024, 2, 29, 14, 30, 0, 0, time.UTC),
			shouldError: false,
		},
		{
			name:        "Midnight",
			dateStr:     "2023:06:15 00:00:00",
			expected:    time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			shouldError: false,
		},
		{
			name:        "End of day",
			dateStr:     "2023:12:31 23:59:59",
			expected:    time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			shouldError: false,
		},
		{
			name:        "Invalid format - wrong separator",
			dateStr:     "2023-01-15 12:30:00",
			expected:    time.Time{},
			shouldError: true,
		},
		{
			name:        "Invalid format - wrong time separator",
			dateStr:     "2023:01:15T12:30:00",
			expected:    time.Time{},
			shouldError: true,
		},
		{
			name:        "Invalid format - missing time",
			dateStr:     "2023:01:15",
			expected:    time.Time{},
			shouldError: true,
		},
		{
			name:        "Invalid format - empty string",
			dateStr:     "",
			expected:    time.Time{},
			shouldError: true,
		},
		{
			name:        "Invalid format - non-numeric",
			dateStr:     "abc:def:ghi jkl:mno:pqr",
			expected:    time.Time{},
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			fields := map[string]interface{}{
				"DateTimeOriginal": tc.dateStr,
			}
			filePath := "/path/to/test/image.jpg"

			result, err := strategy.Extract(fields, filePath)

			if tc.shouldError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultExifDateStrategy_Extract_ComplexFilePath(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
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

			if !result.Equal(expectedDate) {
				t.Errorf("Expected %v, got %v", expectedDate, result)
			}
		})
	}
}

func TestDefaultExifDateStrategy_Extract_OtherExifFields(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	// Test that only DateTimeOriginal is used, other fields are ignored
	fields := map[string]interface{}{
		"DateTime":         "2023:01:16 14:20:00", // Should be ignored
		"CreateDate":       "2023:01:17 16:10:00", // Should be ignored
		"ModifyDate":       "2023:01:18 18:00:00", // Should be ignored
		"DateTimeOriginal": "2023:01:15 12:30:00", // Should be used
	}
	filePath := "/path/to/test/image.jpg"

	result, err := strategy.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDefaultExifDateStrategy_Extract_EdgeCases(t *testing.T) {
	t.Parallel()
	strategy := NewDefaultExifDateStrategy()

	testCases := []struct {
		name     string
		fields   map[string]interface{}
		expected time.Time
	}{
		{
			name:     "DateTimeOriginal as interface{}",
			fields:   map[string]interface{}{"DateTimeOriginal": interface{}("2023:01:15 12:30:00")},
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:     "DateTimeOriginal as any",
			fields:   map[string]interface{}{"DateTimeOriginal": any("2023:01:15 12:30:00")},
			expected: time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:     "DateTimeOriginal as []byte",
			fields:   map[string]interface{}{"DateTimeOriginal": []byte("2023:01:15 12:30:00")},
			expected: time.Time{},
		},
		{
			name:     "DateTimeOriginal as int",
			fields:   map[string]interface{}{"DateTimeOriginal": 12345},
			expected: time.Time{},
		},
		{
			name:     "DateTimeOriginal as float",
			fields:   map[string]interface{}{"DateTimeOriginal": 123.45},
			expected: time.Time{},
		},
		{
			name:     "DateTimeOriginal as bool",
			fields:   map[string]interface{}{"DateTimeOriginal": true},
			expected: time.Time{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			filePath := "/path/to/test/image.jpg"

			result, err := strategy.Extract(tc.fields, filePath)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

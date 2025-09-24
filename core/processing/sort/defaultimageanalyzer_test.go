package sort

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/exif"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultImageAnalyzer(t *testing.T) {
	t.Parallel()

	mockDateProcessor := &testutils.MockDateProcessor{}
	mockExifReader := &exif.DefaultReader{}
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	analyzer, err := NewDefaultImageAnalyzer(mockDateProcessor, mockExifReader, mockLogger, mockLocalizer)
	require.NoError(t, err)

	if analyzer == nil {
		t.Fatal("Expected analyzer to be created, got nil")
	}

	defaultAnalyzer, ok := analyzer.(*DefaultImageAnalyzer)
	if !ok {
		t.Fatal("Expected DefaultImageAnalyzer type")
	}

	if defaultAnalyzer.dateProcessor != mockDateProcessor {
		t.Error("Expected dateProcessor to be injected")
	}
	if defaultAnalyzer.exifReader != mockExifReader {
		t.Error("Expected exifReader to be injected")
	}
	if defaultAnalyzer.logger != mockLogger {
		t.Error("Expected logger to be injected")
	}
	if defaultAnalyzer.localizer != mockLocalizer {
		t.Error("Expected localizer to be injected")
	}
}

func TestDefaultImageAnalyzer_Analyze(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		filePath      string
		exifData      map[string]interface{}
		exifError     error
		dateError     error
		expectedDate  time.Time
		expectedError bool
	}{
		{
			name:         "Successful analysis",
			filePath:     "/path/to/image.jpg",
			exifData:     map[string]interface{}{"DateTimeOriginal": "2023:12:25 14:30:00"},
			expectedDate: time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC),
		},
		{
			name:         "EXIF decode error",
			filePath:     "/path/to/image.jpg",
			exifError:    errors.New("exif decode error"),
			expectedDate: time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC),
		},
		{
			name:          "Date determination error",
			filePath:      "/path/to/image.jpg",
			exifData:      map[string]interface{}{"DateTimeOriginal": "2023:12:25 14:30:00"},
			dateError:     errors.New("date determination error"),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDateProcessor := &testutils.MockDateProcessor{
				GetBestAvailableDateFunc: func(_ map[string]interface{}, _ string) (time.Time, error) {
					return tc.expectedDate, tc.dateError
				},
			}
			mockExifReader := &testutils.MockExifReader{
				ReadExifFunc: func(_ string) (map[string]interface{}, error) {
					return tc.exifData, tc.exifError
				},
			}
			mockLogger := &testutils.MockLogger{}
			mockLocalizer := &testutils.MockLocalizer{
				TranslateFunc: func(key string, _ ...map[string]interface{}) string {
					return key
				},
			}

			analyzer, err := NewDefaultImageAnalyzer(mockDateProcessor, mockExifReader, mockLogger, mockLocalizer)
			require.NoError(t, err)

			image, err := analyzer.Analyze(tc.filePath)

			if tc.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if !image.Date.Equal(tc.expectedDate) {
					t.Errorf("Expected date %v, got %v", tc.expectedDate, image.Date)
				}
			}
		})
	}
}

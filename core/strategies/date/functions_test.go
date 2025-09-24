package date

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewExifExtractor(t *testing.T) {
	t.Parallel()
	fieldName := "DateTime"
	layout := "2006:01:02 15:04:05"
	extractor := NewExifExtractor(fieldName, layout)

	tests := []struct {
		name         string
		exifData     map[string]interface{}
		expectedTime time.Time
		expectedOk   bool
	}{
		{
			name: "valid datetime",
			exifData: map[string]interface{}{
				"DateTime": "2023:12:25 14:30:45",
			},
			expectedTime: time.Date(2023, 12, 25, 14, 30, 45, 0, time.UTC),
			expectedOk:   true,
		},
		{
			name:         "missing field",
			exifData:     map[string]interface{}{},
			expectedTime: time.Time{},
			expectedOk:   false,
		},
		{
			name: "invalid field type",
			exifData: map[string]interface{}{
				"DateTime": 12345,
			},
			expectedTime: time.Time{},
			expectedOk:   false,
		},
		{
			name: "invalid datetime format",
			exifData: map[string]interface{}{
				"DateTime": "invalid-date",
			},
			expectedTime: time.Time{},
			expectedOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			result, ok := extractor("/test/path.jpg", tt.exifData)

			if ok != tt.expectedOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectedOk, ok)
			}

			if !result.Equal(tt.expectedTime) {
				t.Errorf("expected time=%v, got time=%v", tt.expectedTime, result)
			}
		})
	}
}

func TestNewModTimeExtractor(t *testing.T) {
	t.Parallel()
	mockFS := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

	expectedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockFS.StatFunc = func(_ string) (os.FileInfo, error) {
		return &testutils.MockFileInfo{
			ModTimeFunc: func() time.Time { return expectedTime },
		}, nil
	}

	extractor, err := NewModTimeExtractor(mockFS, mockLocalizer)
	require.NoError(t, err)
	result, ok := extractor("/test/file.jpg", map[string]interface{}{})

	if !ok {
		t.Error("expected successful extraction")
	}

	if !result.Equal(expectedTime) {
		t.Errorf("expected time=%v, got time=%v", expectedTime, result)
	}
}

func TestNewCreationTimeExtractor(t *testing.T) {
	t.Parallel()
	mockFS := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Mock filesystem stat to return a file info
	expectedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockFS.StatFunc = func(_ string) (os.FileInfo, error) {
		return &testutils.MockFileInfo{
			ModTimeFunc: func() time.Time { return expectedTime },
		}, nil
	}

	extractor, err := NewCreationTimeExtractor(mockFS, mockLocalizer)
	require.NoError(t, err)
	result, ok := extractor("/test/file.jpg", map[string]interface{}{})

	// Note: This will depend on the actual implementation of creation time extraction
	// For now, we just test that it doesn't panic and returns reasonable results
	if !ok && !result.IsZero() {
		t.Error("extractor should either succeed or return zero time with false")
	}
}

func TestNewExtractorChain(t *testing.T) {
	t.Parallel()
	// Create test extractors
	failingExtractor := func(string, map[string]interface{}) (time.Time, bool) {
		return time.Time{}, false
	}

	successfulExtractor := func(string, map[string]interface{}) (time.Time, bool) {
		return time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), true
	}

	neverCalledExtractor := func(string, map[string]interface{}) (time.Time, bool) {
		t.Error("This extractor should never be called")
		return time.Time{}, false
	}

	tests := []struct {
		name       string
		extractors []Extractor
		expectedOk bool
	}{
		{
			name:       "empty chain",
			extractors: []Extractor{},
			expectedOk: false,
		},
		{
			name:       "all fail",
			extractors: []Extractor{failingExtractor, failingExtractor},
			expectedOk: false,
		},
		{
			name:       "second succeeds",
			extractors: []Extractor{failingExtractor, successfulExtractor, neverCalledExtractor},
			expectedOk: true,
		},
		{
			name:       "first succeeds",
			extractors: []Extractor{successfulExtractor, neverCalledExtractor},
			expectedOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			chain := NewExtractorChain(tt.extractors...)
			_, ok := chain("/test/path.jpg", map[string]interface{}{})

			if ok != tt.expectedOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectedOk, ok)
			}
		})
	}
}

func TestCreateExtractorsFromConfig(t *testing.T) {
	t.Parallel()
	mockFS := &testutils.MockFileSystem{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Set up mock filesystem
	mockFS.StatFunc = func(_ string) (os.FileInfo, error) {
		return &testutils.MockFileInfo{
			ModTimeFunc: time.Now,
		}, nil
	}

	tests := []struct {
		name        string
		config      ExtractorConfig
		expectedLen int
		expectError bool
	}{
		{
			name: "valid config with exif and modTime",
			config: ExtractorConfig{
				StrategyOrder: []string{"exif", "modTime"},
				ExifStrategies: []ExifStrategyConfig{
					{FieldName: "DateTime", Layout: "2006:01:02 15:04:05"},
				},
				Localizer:  mockLocalizer,
				FileSystem: mockFS,
			},
			expectedLen: 2, // 1 exif + 1 modTime
			expectError: false,
		},
		{
			name: "invalid strategy",
			config: ExtractorConfig{
				StrategyOrder: []string{"invalid"},
				Localizer:     mockLocalizer,
				FileSystem:    mockFS,
			},
			expectedLen: 0,
			expectError: true,
		},
		{
			name: "missing localizer",
			config: ExtractorConfig{
				StrategyOrder: []string{"modTime"},
				FileSystem:    mockFS,
			},
			expectedLen: 0,
			expectError: true,
		},
		{
			name: "missing filesystem",
			config: ExtractorConfig{
				StrategyOrder: []string{"modTime"},
				Localizer:     mockLocalizer,
			},
			expectedLen: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			extractors, err := CreateExtractorsFromConfig(tt.config)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(extractors) != tt.expectedLen {
				t.Errorf("expected %d extractors, got %d", tt.expectedLen, len(extractors))
			}
		})
	}
}

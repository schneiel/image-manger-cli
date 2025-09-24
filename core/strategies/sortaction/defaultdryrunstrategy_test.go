package sortaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultDryRunStrategy(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	strategy, err := NewDefaultDryRunStrategy(mockLogger, mockLocalizer)
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
}

func TestDefaultDryRunStrategy_Execute_Success(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "DryRunWouldMoveFile"
		},
	}

	strategy, err := NewDefaultDryRunStrategy(mockLogger, mockLocalizer)
	require.NoError(t, err)
	source := "/path/to/source/image.jpg"
	destination := "/path/to/destination/image.jpg"

	err = strategy.Execute(source, destination)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDefaultDryRunStrategy_Execute_EmptyPaths(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "DryRunWouldMoveFile"
		},
	}

	strategy, err := NewDefaultDryRunStrategy(mockLogger, mockLocalizer)
	require.NoError(t, err)

	testCases := []struct {
		source      string
		destination string
	}{
		{"", "/path/to/destination"},
		{"/path/to/source", ""},
		{"", ""},
		{"/path/to/source", "/path/to/destination"},
	}

	for _, tc := range testCases {
		t.Run(tc.source+"_to_"+tc.destination, func(_ *testing.T) {
			err := strategy.Execute(tc.source, tc.destination)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}

func TestDefaultDryRunStrategy_Execute_ComplexPaths(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "DryRunWouldMoveFile"
		},
	}

	strategy, err := NewDefaultDryRunStrategy(mockLogger, mockLocalizer)
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

// Note: The Execute method calls localizer.Translate() without nil checks.
// This would cause a panic with nil localizer, so we don't test that scenario.
// The code should be refactored to handle nil localizer gracefully.

// Note: The Execute method calls logger.Infof() and localizer.Translate() without nil checks.
// These would cause panics with nil dependencies, so we don't test those scenarios.
// The code should be refactored to handle nil dependencies gracefully.

func TestDefaultDryRunStrategy_GetResources(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	strategy, err := NewDefaultDryRunStrategy(mockLogger, mockLocalizer)
	require.NoError(t, err)

	resources := strategy.GetResources()

	if resources != nil {
		t.Errorf("Expected nil resources, got %+v", resources)
	}
}

func TestDefaultDryRunStrategy_Execute_LoggingVerification(t *testing.T) {
	t.Parallel()
	loggedMessages := []string{}
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(format string, _ ...interface{}) {
			loggedMessages = append(loggedMessages, format)
		},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "DryRunWouldMoveFile"
		},
	}

	strategy, err := NewDefaultDryRunStrategy(mockLogger, mockLocalizer)
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

	if loggedMessages[0] != "DryRunWouldMoveFile" {
		t.Errorf("Expected logged message 'DryRunWouldMoveFile', got '%s'", loggedMessages[0])
	}
}

func TestDefaultDryRunStrategy_Execute_MultipleCalls(t *testing.T) {
	t.Parallel()
	callCount := 0
	mockLogger := &testutils.MockLogger{
		InfofFunc: func(_ string, _ ...interface{}) {
			callCount++
		},
	}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(_ string, _ ...map[string]interface{}) string {
			return "DryRunWouldMoveFile"
		},
	}

	strategy, err := NewDefaultDryRunStrategy(mockLogger, mockLocalizer)
	require.NoError(t, err)

	// Execute multiple times
	for i := 0; i < 5; i++ {
		err := strategy.Execute("/source"+string(rune(i+'0'))+".jpg", "/dest"+string(rune(i+'0'))+".jpg")
		if err != nil {
			t.Fatalf("Expected no error on call %d, got %v", i+1, err)
		}
	}

	if callCount != 5 {
		t.Errorf("Expected 5 logged messages, got %d", callCount)
	}
}

package services_test

import (
	"testing"

	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	"github.com/schneiel/ImageManagerGo/internal/services"
)

const (
	sortSrcDestRequiredKey = "SortSrcDestRequired"
	sortSrcDestRequiredMsg = "Sort source and destination are required"
)

func TestNewSortConfigValidator(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}

	validator := services.NewSortConfigValidator(mockLocalizer)

	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}
}

func TestNewSortConfigValidator_NilLocalizer(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil localizer, but none occurred")
		}
	}()

	services.NewSortConfigValidator(nil)
}

func TestSortConfigValidator_Validate_Success(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		Source:         "/test/source",
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSortConfigValidator_Validate_MissingSource(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		// Source is empty
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := sortSrcDestRequiredMsg
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSortConfigValidator_Validate_MissingDestination(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		Source: "/test/source",
		// Destination is empty
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := sortSrcDestRequiredMsg
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSortConfigValidator_Validate_MissingBoth(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		// Both Source and Destination are empty
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := sortSrcDestRequiredMsg
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSortConfigValidator_Validate_EmptySource(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		Source:         "", // Empty source
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := sortSrcDestRequiredMsg
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSortConfigValidator_Validate_EmptyDestination(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		Source:         "/test/source",
		Destination:    "", // Empty destination
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := sortSrcDestRequiredMsg
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSortConfigValidator_Validate_WhitespaceSource(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		Source:         "   ", // Whitespace-only source (current implementation treats this as valid)
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)
	// Current implementation doesn't check for whitespace-only strings
	if err != nil {
		t.Fatalf("Expected no error for whitespace source, got %v", err)
	}
}

func TestSortConfigValidator_Validate_WhitespaceDestination(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		Source:         "/test/source",
		Destination:    "\t", // Whitespace-only destination (current implementation treats this as valid)
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)
	// Current implementation doesn't check for whitespace-only strings
	if err != nil {
		t.Fatalf("Expected no error for whitespace destination, got %v", err)
	}
}

func TestSortConfigValidator_Validate_MinimalConfig(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewSortConfigValidator(mockLocalizer)

	// Minimal valid config with only required fields
	cfg := &cmdconfig.SortConfig{
		Source:      "/test/source",
		Destination: "/test/destination",
		// ActionStrategy is empty
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSortConfigValidator_Validate_NilConfig(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewSortConfigValidator(mockLocalizer)

	// Should panic when config is nil (current implementation doesn't handle nil)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil config, but none occurred")
		}
	}()

	err := validator.Validate(nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSortConfigValidator_Validate_LocalizerCalls(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	translateCalls := 0
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		translateCalls++
		if key == sortSrcDestRequiredKey {
			return sortSrcDestRequiredMsg
		}

		return key
	}

	validator := services.NewSortConfigValidator(mockLocalizer)

	cfg := &cmdconfig.SortConfig{
		// Missing source to trigger error
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	err := validator.Validate(cfg)
	if err == nil {
		t.Fatal("Expected error for missing source, got nil")
	}

	// Verify localizer was called for error message
	if translateCalls != 1 {
		t.Errorf("Expected 1 translate call, got %d", translateCalls)
	}

	// Verify the error message is correct
	expectedError := sortSrcDestRequiredMsg
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSortConfigValidator_Validate_ValidPaths(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewSortConfigValidator(mockLocalizer)

	testCases := []struct {
		name        string
		source      string
		destination string
		shouldPass  bool
	}{
		{
			name:        "Absolute paths",
			source:      "/absolute/source/path",
			destination: "/absolute/destination/path",
			shouldPass:  true,
		},
		{
			name:        "Relative paths",
			source:      "./relative/source",
			destination: "../relative/destination",
			shouldPass:  true,
		},
		{
			name:        "Home directory paths",
			source:      "~/source",
			destination: "~/destination",
			shouldPass:  true,
		},
		{
			name:        "Current directory",
			source:      ".",
			destination: "./destination",
			shouldPass:  true,
		},
		{
			name:        "Parent directory",
			source:      "..",
			destination: "../destination",
			shouldPass:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg := &cmdconfig.SortConfig{
				Source:         tc.source,
				Destination:    tc.destination,
				ActionStrategy: "copy",
			}

			err := validator.Validate(cfg)

			if tc.shouldPass && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if !tc.shouldPass && err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

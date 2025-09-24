package services_test

import (
	"testing"

	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	"github.com/schneiel/ImageManagerGo/internal/services"
)

const (
	dedupTargetDirMissingKey = "DedupTargetDirMissing"
	dedupTargetDirMissingMsg = "Deduplication target directory is missing"
)

func TestNewDedupConfigValidator(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}

	validator := services.NewDedupConfigValidator(mockLocalizer)

	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}
}

func TestNewDedupConfigValidator_NilLocalizer(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil localizer, but none occurred")
		}
	}()

	services.NewDedupConfigValidator(nil)
}

func TestDedupConfigValidator_Validate_Success(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		Source:         "/test/source",
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
		TrashPath:      ".trash",
		Workers:        4,
		Threshold:      1,
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDedupConfigValidator_Validate_MissingSource(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == dedupTargetDirMissingKey {
			return dedupTargetDirMissingMsg
		}

		return key
	}

	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		// Source is empty
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
		Workers:        4,
		Threshold:      1,
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "Deduplication target directory is missing"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDedupConfigValidator_Validate_EmptySource(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == dedupTargetDirMissingKey {
			return dedupTargetDirMissingMsg
		}

		return key
	}

	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		Source:         "", // Empty source
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
		Workers:        4,
		Threshold:      1,
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "Deduplication target directory is missing"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDedupConfigValidator_Validate_WhitespaceSource(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == dedupTargetDirMissingKey {
			return dedupTargetDirMissingMsg
		}

		return key
	}

	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		Source:         "   ", // Whitespace-only source (current implementation treats this as valid)
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
		Workers:        4,
		Threshold:      1,
	}

	err := validator.Validate(cfg)
	// Current implementation doesn't check for whitespace-only strings
	if err != nil {
		t.Fatalf("Expected no error for whitespace source, got %v", err)
	}
}

func TestDedupConfigValidator_Validate_ZeroWorkers(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		Source:    "/test/source",
		Workers:   0, // Zero workers should be set to default
		Threshold: 1,
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify workers was set to default value
	defaultCfg := cmdconfig.DefaultDedupConfig()
	if cfg.Workers != defaultCfg.Workers {
		t.Errorf("Expected Workers to be set to default %d, got %d", defaultCfg.Workers, cfg.Workers)
	}
}

func TestDedupConfigValidator_Validate_NegativeWorkers(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		Source:    "/test/source",
		Workers:   -1, // Negative workers should be set to default
		Threshold: 1,
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify workers was set to default value
	defaultCfg := cmdconfig.DefaultDedupConfig()
	if cfg.Workers != defaultCfg.Workers {
		t.Errorf("Expected Workers to be set to default %d, got %d", defaultCfg.Workers, cfg.Workers)
	}
}

func TestDedupConfigValidator_Validate_NegativeThreshold(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		Source:    "/test/source",
		Workers:   4,
		Threshold: -1, // Negative threshold should cause error
	}

	err := validator.Validate(cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "threshold must be non-negative"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDedupConfigValidator_Validate_ZeroThreshold(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		Source:    "/test/source",
		Workers:   4,
		Threshold: 0, // Zero threshold should be valid
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDedupConfigValidator_Validate_MinimalConfig(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewDedupConfigValidator(mockLocalizer)

	// Minimal valid config with only required fields
	cfg := &cmdconfig.DedupConfig{
		Source: "/test/source",
		// All other fields are zero/empty values
	}

	err := validator.Validate(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify workers was set to default value
	defaultCfg := cmdconfig.DefaultDedupConfig()
	if cfg.Workers != defaultCfg.Workers {
		t.Errorf("Expected Workers to be set to default %d, got %d", defaultCfg.Workers, cfg.Workers)
	}
}

func TestDedupConfigValidator_Validate_NilConfig(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	validator := services.NewDedupConfigValidator(mockLocalizer)

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

func TestDedupConfigValidator_Validate_LocalizerCalls(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	translateCalls := 0
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		translateCalls++
		if key == dedupTargetDirMissingKey {
			return dedupTargetDirMissingMsg
		}

		return key
	}

	validator := services.NewDedupConfigValidator(mockLocalizer)

	cfg := &cmdconfig.DedupConfig{
		// Missing source to trigger error
		Workers:   4,
		Threshold: 1,
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
	expectedError := dedupTargetDirMissingMsg
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

package handlers

import (
	"errors"
	"testing"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

const (
	dedupProcessStartingKey  = "DedupProcessStarting"
	dedupProcessCompletedKey = "DedupProcessCompleted"
	dedupProcessStartingMsg  = "Starting deduplication process"
	dedupProcessCompletedMsg = "Deduplication process completed"
	invalidConfigMsg         = "not a dedup config"
)

func TestNewDedupHandler(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.BaseHandler != mockBaseHandler {
		t.Error("Expected BaseHandler to be injected")
	}

	if handler.Executor != mockExecutor {
		t.Error("Expected Executor to be injected")
	}

	if handler.Applier != mockApplier {
		t.Error("Expected Applier to be injected")
	}

	if handler.Validator != mockValidator {
		t.Error("Expected Validator to be injected")
	}

	if handler.Localizer != mockLocalizer {
		t.Error("Expected Localizer to be injected")
	}
}

func TestNewDedupHandler_NilDependencies(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Test nil executor
	t.Run("Nil Executor", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil executor, but none occurred")
			}
		}()
		NewDedupHandler(mockBaseHandler, nil, mockApplier, mockValidator, mockLocalizer)
	})

	// Test nil applier
	t.Run("Nil Applier", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil applier, but none occurred")
			}
		}()
		NewDedupHandler(mockBaseHandler, mockExecutor, nil, mockValidator, mockLocalizer)
	})

	// Test nil validator
	t.Run("Nil Validator", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil validator, but none occurred")
			}
		}()
		NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, nil, mockLocalizer)
	})

	// Test nil localizer
	t.Run("Nil Localizer", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil localizer, but none occurred")
			}
		}()
		NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, nil)
	})
}

func TestDedupHandler_RunE_Success(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Configure infrastructure mocks for handler testing (appropriate pattern)
	mockValidator.ValidateFunc = func(_ *cmdconfig.DedupConfig) error {
		return nil
	}

	mockApplier.ApplyFunc = func(src *cmdconfig.DedupConfig, dest *config.Config) {
		dest.Deduplicator.Source = src.Source
	}

	mockExecutor.ExecuteFunc = func(_ string, _ config.Config) error {
		return nil
	}

	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case dedupProcessStartingKey:
			return dedupProcessStartingMsg
		case dedupProcessCompletedKey:
			return dedupProcessCompletedMsg
		default:
			return key
		}
	}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)
	handler.Config = &cmdconfig.DedupConfig{
		Source:         "/test/source",
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
	}

	err := handler.RunE(nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDedupHandler_RunE_ValidationError(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	expectedError := errors.New("validation failed")
	mockValidator.ValidateFunc = func(_ *cmdconfig.DedupConfig) error {
		return expectedError
	}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)
	handler.Config = &cmdconfig.DedupConfig{
		Source: "/test/source",
	}

	err := handler.RunE(nil, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedWrappedError := "dedup configuration validation failed: " + expectedError.Error()
	if err.Error() != expectedWrappedError {
		t.Errorf("Expected error '%s', got '%s'", expectedWrappedError, err.Error())
	}
}

func TestDedupHandler_RunE_ExecutionError(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Configure infrastructure mocks for handler testing (appropriate pattern)
	mockValidator.ValidateFunc = func(_ *cmdconfig.DedupConfig) error {
		return nil
	}

	mockApplier.ApplyFunc = func(src *cmdconfig.DedupConfig, dest *config.Config) {
		dest.Deduplicator.Source = src.Source
	}

	expectedError := errors.New("execution failed")
	mockExecutor.ExecuteFunc = func(_ string, _ config.Config) error {
		return expectedError
	}

	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case dedupProcessStartingKey:
			return dedupProcessStartingMsg
		case "DedupProcessFailed":
			return "Deduplication process failed"
		default:
			return key
		}
	}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)
	handler.Config = &cmdconfig.DedupConfig{
		Source: "/test/source",
	}

	err := handler.RunE(nil, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedErrorMsg := "Deduplication process failed"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestDedupHandler_Execute_Success(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Configure infrastructure mocks for handler testing (appropriate pattern)
	mockValidator.ValidateFunc = func(_ *cmdconfig.DedupConfig) error {
		return nil
	}

	mockApplier.ApplyFunc = func(src *cmdconfig.DedupConfig, dest *config.Config) {
		dest.Deduplicator.Source = src.Source
	}

	mockExecutor.ExecuteFunc = func(_ string, _ config.Config) error {
		return nil
	}

	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case dedupProcessStartingKey:
			return dedupProcessStartingMsg
		case dedupProcessCompletedKey:
			return dedupProcessCompletedMsg
		default:
			return key
		}
	}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)
	dedupConfig := &cmdconfig.DedupConfig{
		Source:         "/test/source",
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
	}

	err := handler.Execute(dedupConfig)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if handler.Config != dedupConfig {
		t.Error("Expected Config to be set")
	}
}

func TestDedupHandler_Execute_InvalidConfigType(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Configure infrastructure mocks for handler testing (appropriate pattern) to handle nil config gracefully
	mockValidator.ValidateFunc = func(_ *cmdconfig.DedupConfig) error {
		return nil
	}

	mockApplier.ApplyFunc = func(src *cmdconfig.DedupConfig, dest *config.Config) {
		if src != nil {
			dest.Deduplicator.Source = src.Source
		}
	}

	mockExecutor.ExecuteFunc = func(_ string, _ config.Config) error {
		return nil
	}

	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case dedupProcessStartingKey:
			return dedupProcessStartingMsg
		case "DedupProcessCompleted":
			return dedupProcessCompletedMsg
		default:
			return key
		}
	}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)

	// Pass wrong config type - should be ignored and execution should fail due to nil config
	invalidConfig := invalidConfigMsg
	err := handler.Execute(invalidConfig)

	if err == nil {
		t.Fatal("Expected error for invalid config type, got nil")
	}

	// Config should not be set for invalid type
	if handler.Config != nil {
		t.Error("Expected Config to remain nil for invalid config type")
	}
}

func TestDedupHandler_Validate_Success(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	mockValidator.ValidateFunc = func(_ *cmdconfig.DedupConfig) error {
		return nil
	}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)
	dedupConfig := &cmdconfig.DedupConfig{
		Source: "/test/source",
	}

	err := handler.Validate(dedupConfig)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDedupHandler_Validate_Error(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	expectedError := errors.New("validation failed")
	mockValidator.ValidateFunc = func(_ *cmdconfig.DedupConfig) error {
		return expectedError
	}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)
	dedupConfig := &cmdconfig.DedupConfig{
		Source: "/test/source",
	}

	err := handler.Validate(dedupConfig)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedWrappedError := "dedup validation failed: " + expectedError.Error()
	if err.Error() != expectedWrappedError {
		t.Errorf("Expected error '%s', got '%s'", expectedWrappedError, err.Error())
	}
}

func TestDedupHandler_Validate_InvalidConfigType(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.DedupConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.DedupConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	handler := NewDedupHandler(mockBaseHandler, mockExecutor, mockApplier, mockValidator, mockLocalizer)

	// Pass wrong config type
	invalidConfig := invalidConfigMsg
	err := handler.Validate(invalidConfig)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "invalid config type for dedup handler"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

package handlers

import (
	"errors"
	"testing"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

const (
	sortProcessStartingKey  = "SortProcessStarting"
	sortProcessCompletedKey = "SortProcessCompleted"
	sortProcessFailedKey    = "SortProcessFailed"
	sortProcessStartingMsg  = "Starting sort process"
	sortProcessCompletedMsg = "Sort process completed"
	sortProcessFailedMsg    = "Sort process failed"
	invalidSortConfigMsg    = "not a sort config"
)

func TestNewSortHandler(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.SortConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockConfig := &cmdconfig.SortConfig{}

	handler := &SortHandler{
		BaseHandler: mockBaseHandler,
		Executor:    mockExecutor,
		Applier:     mockApplier,
		Validator:   mockValidator,
		Localizer:   mockLocalizer,
		Config:      mockConfig,
	}

	if handler.BaseHandler != mockBaseHandler {
		t.Error("BaseHandler not set correctly")
	}
	if handler.Executor != mockExecutor {
		t.Error("Executor not set correctly")
	}
	if handler.Applier != mockApplier {
		t.Error("Applier not set correctly")
	}
	if handler.Validator != mockValidator {
		t.Error("Validator not set correctly")
	}
	if handler.Localizer != mockLocalizer {
		t.Error("Localizer not set correctly")
	}
	if handler.Config != mockConfig {
		t.Error("Config not set correctly")
	}
}

func TestSortHandler_RunE_Success(t *testing.T) {
	t.Parallel()
	mockBaseHandler := &BaseHandler{
		Logger: &testutils.MockLogger{},
		Config: &config.Config{},
	}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.SortConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Configure infrastructure mocks for handler testing (appropriate pattern)
	mockValidator.ValidateFunc = func(_ *cmdconfig.SortConfig) error {
		return nil
	}

	mockApplier.ApplyFunc = func(src *cmdconfig.SortConfig, dest *config.Config) {
		dest.Sorter.Source = src.Source
	}

	mockExecutor.ExecuteFunc = func(_ string, _ config.Config) error {
		return nil
	}

	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case sortProcessStartingKey:
			return sortProcessStartingMsg
		case sortProcessCompletedKey:
			return sortProcessCompletedMsg
		default:
			return key
		}
	}

	handler := &SortHandler{
		BaseHandler: mockBaseHandler,
		Executor:    mockExecutor,
		Applier:     mockApplier,
		Validator:   mockValidator,
		Localizer:   mockLocalizer,
		Config:      &cmdconfig.SortConfig{},
	}

	err := handler.RunE(nil, nil)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestSortHandler_RunE_ValidationError(t *testing.T) {
	t.Parallel()
	validationError := errors.New("validation failed")
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	mockValidator.ValidateFunc = func(_ *cmdconfig.SortConfig) error {
		return validationError
	}

	handler := &SortHandler{
		BaseHandler: &BaseHandler{
			Logger: &testutils.MockLogger{},
			Config: &config.Config{},
		},
		Validator: mockValidator,
		Config:    &cmdconfig.SortConfig{},
	}

	err := handler.RunE(nil, nil)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
	if !errors.Is(err, validationError) {
		t.Errorf("Expected validation error wrapped, got: %v", err)
	}
}

func TestSortHandler_RunE_ExecutorError(t *testing.T) {
	t.Parallel()
	executorError := errors.New("execution failed")
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.SortConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	mockValidator.ValidateFunc = func(_ *cmdconfig.SortConfig) error {
		return nil
	}

	mockApplier.ApplyFunc = func(src *cmdconfig.SortConfig, dest *config.Config) {
		dest.Sorter.Source = src.Source
	}

	mockExecutor.ExecuteFunc = func(_ string, _ config.Config) error {
		return executorError
	}

	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case sortProcessStartingKey:
			return sortProcessStartingMsg
		case sortProcessFailedKey:
			return sortProcessFailedMsg
		default:
			return key
		}
	}

	handler := &SortHandler{
		BaseHandler: &BaseHandler{
			Logger: &testutils.MockLogger{},
			Config: &config.Config{},
		},
		Executor:  mockExecutor,
		Applier:   mockApplier,
		Validator: mockValidator,
		Localizer: mockLocalizer,
		Config:    &cmdconfig.SortConfig{},
	}

	err := handler.RunE(nil, nil)

	if err == nil {
		t.Error("Expected executor error, got nil")
	}
	if err.Error() != sortProcessFailedMsg {
		t.Errorf("Expected error message %q, got: %v", sortProcessFailedMsg, err)
	}
}

func TestSortHandler_Execute_Success(t *testing.T) {
	t.Parallel()
	mockExecutor := &testutils.MockTaskExecutor{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.SortConfig]{}
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	mockLocalizer := &testutils.MockLocalizer{}

	sortConfig := &cmdconfig.SortConfig{}

	mockValidator.ValidateFunc = func(_ *cmdconfig.SortConfig) error {
		return nil
	}

	mockApplier.ApplyFunc = func(src *cmdconfig.SortConfig, dest *config.Config) {
		dest.Sorter.Source = src.Source
	}

	mockExecutor.ExecuteFunc = func(_ string, _ config.Config) error {
		return nil
	}

	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case sortProcessStartingKey:
			return sortProcessStartingMsg
		case sortProcessCompletedKey:
			return sortProcessCompletedMsg
		default:
			return key
		}
	}

	handler := &SortHandler{
		BaseHandler: &BaseHandler{
			Logger: &testutils.MockLogger{},
			Config: &config.Config{},
		},
		Executor:  mockExecutor,
		Applier:   mockApplier,
		Validator: mockValidator,
		Localizer: mockLocalizer,
	}

	err := handler.Execute(sortConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if handler.Config != sortConfig {
		t.Error("Config not set correctly during execution")
	}
}

func TestSortHandler_Execute_NilValidator(t *testing.T) {
	t.Parallel()
	handler := &SortHandler{
		BaseHandler: &BaseHandler{
			Logger: &testutils.MockLogger{},
			Config: &config.Config{},
		},
		Validator: nil,
	}

	err := handler.Execute(&cmdconfig.SortConfig{})

	if err == nil {
		t.Error("Expected nil validator error, got nil")
	}
	expectedMsg := "sort handler validator is nil (dependency injection failure)"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error %q, got: %v", expectedMsg, err)
	}
}

func TestSortHandler_Execute_InvalidConfigType(t *testing.T) {
	t.Parallel()
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	mockApplier := &testutils.MockConfigApplier[*cmdconfig.SortConfig]{}
	mockExecutor := &testutils.MockTaskExecutor{}
	mockLocalizer := &testutils.MockLocalizer{}

	handler := &SortHandler{
		BaseHandler: &BaseHandler{
			Logger: &testutils.MockLogger{},
			Config: &config.Config{},
		},
		Validator: mockValidator,
		Applier:   mockApplier,
		Executor:  mockExecutor,
		Localizer: mockLocalizer,
	}

	// Pass wrong config type - should trigger nil config and return immediately
	err := handler.Execute(&cmdconfig.DedupConfig{})

	if err != nil {
		t.Errorf("Expected no error for type assertion failure, got: %v", err)
	}
}

func TestSortHandler_Validate_Success(t *testing.T) {
	t.Parallel()
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	sortConfig := &cmdconfig.SortConfig{}

	mockValidator.ValidateFunc = func(_ *cmdconfig.SortConfig) error {
		return nil
	}

	handler := &SortHandler{
		Validator: mockValidator,
	}

	err := handler.Validate(sortConfig)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestSortHandler_Validate_ValidationError(t *testing.T) {
	t.Parallel()
	validationError := errors.New("validation failed")
	mockValidator := &testutils.MockConfigValidator[*cmdconfig.SortConfig]{}
	sortConfig := &cmdconfig.SortConfig{}

	mockValidator.ValidateFunc = func(_ *cmdconfig.SortConfig) error {
		return validationError
	}

	handler := &SortHandler{
		Validator: mockValidator,
	}

	err := handler.Validate(sortConfig)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
	if !errors.Is(err, validationError) {
		t.Errorf("Expected validation error wrapped, got: %v", err)
	}
}

func TestSortHandler_Validate_InvalidConfigType(t *testing.T) {
	t.Parallel()
	handler := &SortHandler{
		Validator: &testutils.MockConfigValidator[*cmdconfig.SortConfig]{},
	}

	err := handler.Validate(&cmdconfig.DedupConfig{})

	if err == nil {
		t.Error("Expected invalid config type error, got nil")
	}
	expectedMsg := "invalid config type for sort handler"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error %q, got: %v", expectedMsg, err)
	}
}

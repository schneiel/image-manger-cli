package handlers

import (
	"testing"

	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

func TestWithHandlerBaseHandler_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	baseHandler := &BaseHandler{}
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerBaseHandler[*cmdconfig.SortConfig](baseHandler)

	// Act
	err := option(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if config.BaseHandler != baseHandler {
		t.Error("Expected BaseHandler to be set correctly")
	}
}

func TestWithHandlerBaseHandler_NilHandler(t *testing.T) {
	t.Parallel()

	// Arrange
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerBaseHandler[*cmdconfig.SortConfig](nil)

	// Act
	err := option(config)

	// Assert
	if err == nil {
		t.Error("Expected error for nil baseHandler")
	}
	if err.Error() != "baseHandler cannot be nil" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestWithHandlerExecutor_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	executor := &testutils.FakeTaskExecutor{}
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerExecutor[*cmdconfig.SortConfig](executor)

	// Act
	err := option(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if config.Executor != executor {
		t.Error("Expected Executor to be set correctly")
	}
}

func TestWithHandlerExecutor_NilExecutor(t *testing.T) {
	t.Parallel()

	// Arrange
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerExecutor[*cmdconfig.SortConfig](nil)

	// Act
	err := option(config)

	// Assert
	if err == nil {
		t.Error("Expected error for nil executor")
	}
	if err.Error() != "executor cannot be nil" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestWithHandlerApplier_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	applier := &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{}
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerApplier[*cmdconfig.SortConfig](applier)

	// Act
	err := option(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if config.Applier != applier {
		t.Error("Expected Applier to be set correctly")
	}
}

func TestWithHandlerApplier_NilApplier(t *testing.T) {
	t.Parallel()

	// Arrange
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerApplier[*cmdconfig.SortConfig](nil)

	// Act
	err := option(config)

	// Assert
	if err == nil {
		t.Error("Expected error for nil applier")
	}
	if err.Error() != "applier cannot be nil" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestWithHandlerValidator_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	validator := &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{}
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerValidator[*cmdconfig.SortConfig](validator)

	// Act
	err := option(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if config.Validator != validator {
		t.Error("Expected Validator to be set correctly")
	}
}

func TestWithHandlerValidator_NilValidator(t *testing.T) {
	t.Parallel()

	// Arrange
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerValidator[*cmdconfig.SortConfig](nil)

	// Act
	err := option(config)

	// Assert
	if err == nil {
		t.Error("Expected error for nil validator")
	}
	if err.Error() != "validator cannot be nil" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestWithHandlerLocalizer_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerLocalizer[*cmdconfig.SortConfig](localizer)

	// Act
	err := option(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if config.Localizer != localizer {
		t.Error("Expected Localizer to be set correctly")
	}
}

func TestWithHandlerLocalizer_NilLocalizer(t *testing.T) {
	t.Parallel()

	// Arrange
	config := &HandlerConfig[*cmdconfig.SortConfig]{}
	option := WithHandlerLocalizer[*cmdconfig.SortConfig](nil)

	// Act
	err := option(config)

	// Assert
	if err == nil {
		t.Error("Expected error for nil localizer")
	}
	if err.Error() != "localizer cannot be nil" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestValidateHandlerConfig_AllFieldsSet(t *testing.T) {
	t.Parallel()

	// Arrange
	config := &HandlerConfig[*cmdconfig.SortConfig]{
		BaseHandler: &BaseHandler{},
		Executor:    &testutils.FakeTaskExecutor{},
		Applier:     &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{},
		Validator:   &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{},
		Localizer:   testutils.NewFakeLocalizer(),
	}

	// Act
	err := ValidateHandlerConfig(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}
}

func TestValidateHandlerConfig_MissingFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *HandlerConfig[*cmdconfig.SortConfig]
		wantErr string
	}{
		{
			name: "missing baseHandler",
			config: &HandlerConfig[*cmdconfig.SortConfig]{
				Executor:  &testutils.FakeTaskExecutor{},
				Applier:   &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{},
				Validator: &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{},
				Localizer: testutils.NewFakeLocalizer(),
			},
			wantErr: "baseHandler is required",
		},
		{
			name: "missing executor",
			config: &HandlerConfig[*cmdconfig.SortConfig]{
				BaseHandler: &BaseHandler{},
				Applier:     &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{},
				Validator:   &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{},
				Localizer:   testutils.NewFakeLocalizer(),
			},
			wantErr: "executor is required",
		},
		{
			name: "missing applier",
			config: &HandlerConfig[*cmdconfig.SortConfig]{
				BaseHandler: &BaseHandler{},
				Executor:    &testutils.FakeTaskExecutor{},
				Validator:   &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{},
				Localizer:   testutils.NewFakeLocalizer(),
			},
			wantErr: "applier is required",
		},
		{
			name: "missing validator",
			config: &HandlerConfig[*cmdconfig.SortConfig]{
				BaseHandler: &BaseHandler{},
				Executor:    &testutils.FakeTaskExecutor{},
				Applier:     &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{},
				Localizer:   testutils.NewFakeLocalizer(),
			},
			wantErr: "validator is required",
		},
		{
			name: "missing localizer",
			config: &HandlerConfig[*cmdconfig.SortConfig]{
				BaseHandler: &BaseHandler{},
				Executor:    &testutils.FakeTaskExecutor{},
				Applier:     &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{},
				Validator:   &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{},
			},
			wantErr: "localizer is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			err := ValidateHandlerConfig(tt.config)

			// Assert
			if err == nil {
				t.Error("Expected error for invalid config")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Expected error '%s', got '%s'", tt.wantErr, err.Error())
			}
		})
	}
}

// Integration test using real sort handler builder.
func TestSortHandlerBuilder_Integration(t *testing.T) {
	t.Parallel()

	// Test that sort handler validation error flows through correctly
	handler, err := NewSortHandlerWithOptions()

	// Assert
	if err == nil {
		t.Error("Expected validation error")
	}
	if handler != nil {
		t.Error("Expected nil handler when validation fails")
	}
	if err.Error() != "baseHandler is required" {
		t.Errorf("Expected validation error message, got: %v", err)
	}
}

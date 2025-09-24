package handlers

import (
	"testing"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

func TestCreateDedupHandler(t *testing.T) {
	t.Parallel()

	// Arrange
	baseHandler := &BaseHandler{}
	executor := &testutils.FakeTaskExecutor{}
	applier := &testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{}
	validator := &testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{}
	localizer := testutils.NewFakeLocalizer()

	// Act
	handler := NewDedupHandler(baseHandler, executor, applier, validator, localizer)

	// Assert
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}
	if handler.BaseHandler != baseHandler {
		t.Error("Expected BaseHandler to be set correctly")
	}
	if handler.Executor != executor {
		t.Error("Expected Executor to be set correctly")
	}
	if handler.Applier != applier {
		t.Error("Expected Applier to be set correctly")
	}
	if handler.Validator != validator {
		t.Error("Expected Validator to be set correctly")
	}
	if handler.Localizer != localizer {
		t.Error("Expected Localizer to be set correctly")
	}
	if handler.Config != nil {
		t.Error("Expected Config to be nil initially")
	}
}

func TestCreateDedupHandler_ValidDependencies(t *testing.T) {
	t.Parallel()

	// Arrange
	baseHandler := &BaseHandler{}
	executor := &testutils.FakeTaskExecutor{}
	applier := &testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{}
	validator := &testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{}
	localizer := testutils.NewFakeLocalizer()

	// Act
	handler := createDedupHandler(baseHandler, executor, applier, validator, localizer)

	// Assert
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}
	if handler.BaseHandler != baseHandler {
		t.Error("Expected BaseHandler to be set correctly")
	}
	if handler.Executor != executor {
		t.Error("Expected Executor to be set correctly")
	}
	if handler.Applier != applier {
		t.Error("Expected Applier to be set correctly")
	}
	if handler.Validator != validator {
		t.Error("Expected Validator to be set correctly")
	}
	if handler.Localizer != localizer {
		t.Error("Expected Localizer to be set correctly")
	}
	if handler.Config != nil {
		t.Error("Expected Config to be nil initially")
	}
}

func TestCreateDedupHandler_NilDependencies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		baseHandler *BaseHandler
		executor    TaskExecutor
		applier     ConfigApplier[*cmdconfig.DedupConfig]
		validator   ConfigValidator[*cmdconfig.DedupConfig]
		localizer   i18n.Localizer
		wantPanic   string
	}{
		{
			name:        "nil baseHandler",
			baseHandler: nil,
			executor:    &testutils.FakeTaskExecutor{},
			applier:     &testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{},
			validator:   &testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{},
			localizer:   testutils.NewFakeLocalizer(),
			wantPanic:   "baseHandler cannot be nil",
		},
		{
			name:        "nil executor",
			baseHandler: &BaseHandler{},
			executor:    nil,
			applier:     &testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{},
			validator:   &testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{},
			localizer:   testutils.NewFakeLocalizer(),
			wantPanic:   "executor cannot be nil",
		},
		{
			name:        "nil applier",
			baseHandler: &BaseHandler{},
			executor:    &testutils.FakeTaskExecutor{},
			applier:     nil,
			validator:   &testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{},
			localizer:   testutils.NewFakeLocalizer(),
			wantPanic:   "applier cannot be nil",
		},
		{
			name:        "nil validator",
			baseHandler: &BaseHandler{},
			executor:    &testutils.FakeTaskExecutor{},
			applier:     &testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{},
			validator:   nil,
			localizer:   testutils.NewFakeLocalizer(),
			wantPanic:   "validator cannot be nil",
		},
		{
			name:        "nil localizer",
			baseHandler: &BaseHandler{},
			executor:    &testutils.FakeTaskExecutor{},
			applier:     &testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{},
			validator:   &testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{},
			localizer:   nil,
			wantPanic:   "localizer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act & Assert
			defer func() {
				if r := recover(); r != nil {
					if r != tt.wantPanic {
						t.Errorf("Expected panic '%s', got '%v'", tt.wantPanic, r)
					}
				} else {
					t.Error("Expected panic, but function completed normally")
				}
			}()

			createDedupHandler(tt.baseHandler, tt.executor, tt.applier, tt.validator, tt.localizer)
		})
	}
}

func TestNewDedupHandlerWithOptions_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	baseHandler := &BaseHandler{}
	executor := &testutils.FakeTaskExecutor{}
	applier := &testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{}
	validator := &testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{}
	localizer := testutils.NewFakeLocalizer()

	options := []DedupHandlerOption{
		WithBaseHandler(baseHandler),
		WithTaskExecutor(executor),
		WithDedupConfigApplier(applier),
		WithDedupConfigValidator(validator),
		WithLocalizer(localizer),
	}

	// Act
	handler, err := NewDedupHandlerWithOptions(options...)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}
	if handler.BaseHandler != baseHandler {
		t.Error("Expected BaseHandler to be set correctly")
	}
	if handler.Executor != executor {
		t.Error("Expected Executor to be set correctly")
	}
	if handler.Applier != applier {
		t.Error("Expected Applier to be set correctly")
	}
	if handler.Validator != validator {
		t.Error("Expected Validator to be set correctly")
	}
	if handler.Localizer != localizer {
		t.Error("Expected Localizer to be set correctly")
	}
}

func TestNewDedupHandlerWithOptions_MissingDependencies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options []DedupHandlerOption
		wantErr string
	}{
		{
			name:    "missing all dependencies",
			options: []DedupHandlerOption{},
			wantErr: "baseHandler is required",
		},
		{
			name: "missing executor",
			options: []DedupHandlerOption{
				WithBaseHandler(&BaseHandler{}),
			},
			wantErr: "executor is required",
		},
		{
			name: "missing applier",
			options: []DedupHandlerOption{
				WithBaseHandler(&BaseHandler{}),
				WithTaskExecutor(&testutils.FakeTaskExecutor{}),
			},
			wantErr: "applier is required",
		},
		{
			name: "missing validator",
			options: []DedupHandlerOption{
				WithBaseHandler(&BaseHandler{}),
				WithTaskExecutor(&testutils.FakeTaskExecutor{}),
				WithDedupConfigApplier(&testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{}),
			},
			wantErr: "validator is required",
		},
		{
			name: "missing localizer",
			options: []DedupHandlerOption{
				WithBaseHandler(&BaseHandler{}),
				WithTaskExecutor(&testutils.FakeTaskExecutor{}),
				WithDedupConfigApplier(&testutils.FakeConfigApplier[*cmdconfig.DedupConfig]{}),
				WithDedupConfigValidator(&testutils.FakeConfigValidator[*cmdconfig.DedupConfig]{}),
			},
			wantErr: "localizer is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			handler, err := NewDedupHandlerWithOptions(tt.options...)

			// Assert
			if err == nil {
				t.Error("Expected error, got nil")
			}
			if handler != nil {
				t.Error("Expected nil handler when error occurs")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Expected error '%s', got '%s'", tt.wantErr, err.Error())
			}
		})
	}
}

func TestNewDedupHandlerWithOptions_NilDependencies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options []DedupHandlerOption
		wantErr string
	}{
		{
			name: "nil base handler",
			options: []DedupHandlerOption{
				WithBaseHandler(nil),
			},
			wantErr: "baseHandler cannot be nil",
		},
		{
			name: "nil executor",
			options: []DedupHandlerOption{
				WithTaskExecutor(nil),
			},
			wantErr: "executor cannot be nil",
		},
		{
			name: "nil applier",
			options: []DedupHandlerOption{
				WithDedupConfigApplier(nil),
			},
			wantErr: "applier cannot be nil",
		},
		{
			name: "nil validator",
			options: []DedupHandlerOption{
				WithDedupConfigValidator(nil),
			},
			wantErr: "validator cannot be nil",
		},
		{
			name: "nil localizer",
			options: []DedupHandlerOption{
				WithLocalizer(nil),
			},
			wantErr: "localizer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			handler, err := NewDedupHandlerWithOptions(tt.options...)

			// Assert
			if err == nil {
				t.Error("Expected error, got nil")
			}
			if handler != nil {
				t.Error("Expected nil handler when error occurs")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Expected error '%s', got '%s'", tt.wantErr, err.Error())
			}
		})
	}
}

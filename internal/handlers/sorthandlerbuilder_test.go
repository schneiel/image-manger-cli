package handlers

import (
	"testing"

	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

func TestCreateSortHandler(t *testing.T) {
	t.Parallel()

	// Arrange
	baseHandler := &BaseHandler{}
	executor := &testutils.FakeTaskExecutor{}
	applier := &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{}
	validator := &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{}
	localizer := testutils.NewFakeLocalizer()

	// Act
	handler := NewSortHandler(baseHandler, executor, applier, validator, localizer)

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
	if handler.Config == nil {
		t.Error("Expected Config to be initialized")
	}
}

func TestNewSortHandlerWithOptions_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	baseHandler := &BaseHandler{}
	executor := &testutils.FakeTaskExecutor{}
	applier := &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{}
	validator := &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{}
	localizer := testutils.NewFakeLocalizer()

	options := []SortHandlerOption{
		WithSortBaseHandler(baseHandler),
		WithSortTaskExecutor(executor),
		WithSortConfigApplier(applier),
		WithSortConfigValidator(validator),
		WithSortLocalizer(localizer),
	}

	// Act
	handler, err := NewSortHandlerWithOptions(options...)

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

func TestNewSortHandlerWithOptions_MissingDependencies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options []SortHandlerOption
		wantErr string
	}{
		{
			name:    "missing all dependencies",
			options: []SortHandlerOption{},
			wantErr: "baseHandler is required",
		},
		{
			name: "missing executor",
			options: []SortHandlerOption{
				WithSortBaseHandler(&BaseHandler{}),
			},
			wantErr: "executor is required",
		},
		{
			name: "missing applier",
			options: []SortHandlerOption{
				WithSortBaseHandler(&BaseHandler{}),
				WithSortTaskExecutor(&testutils.FakeTaskExecutor{}),
			},
			wantErr: "applier is required",
		},
		{
			name: "missing validator",
			options: []SortHandlerOption{
				WithSortBaseHandler(&BaseHandler{}),
				WithSortTaskExecutor(&testutils.FakeTaskExecutor{}),
				WithSortConfigApplier(&testutils.FakeConfigApplier[*cmdconfig.SortConfig]{}),
			},
			wantErr: "validator is required",
		},
		{
			name: "missing localizer",
			options: []SortHandlerOption{
				WithSortBaseHandler(&BaseHandler{}),
				WithSortTaskExecutor(&testutils.FakeTaskExecutor{}),
				WithSortConfigApplier(&testutils.FakeConfigApplier[*cmdconfig.SortConfig]{}),
				WithSortConfigValidator(&testutils.FakeConfigValidator[*cmdconfig.SortConfig]{}),
			},
			wantErr: "localizer is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			handler, err := NewSortHandlerWithOptions(tt.options...)

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

func TestNewSortHandlerWithOptions_NilDependencies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options []SortHandlerOption
		wantErr string
	}{
		{
			name: "nil base handler",
			options: []SortHandlerOption{
				WithSortBaseHandler(nil),
			},
			wantErr: "baseHandler cannot be nil",
		},
		{
			name: "nil executor",
			options: []SortHandlerOption{
				WithSortTaskExecutor(nil),
			},
			wantErr: "executor cannot be nil",
		},
		{
			name: "nil applier",
			options: []SortHandlerOption{
				WithSortConfigApplier(nil),
			},
			wantErr: "applier cannot be nil",
		},
		{
			name: "nil validator",
			options: []SortHandlerOption{
				WithSortConfigValidator(nil),
			},
			wantErr: "validator cannot be nil",
		},
		{
			name: "nil localizer",
			options: []SortHandlerOption{
				WithSortLocalizer(nil),
			},
			wantErr: "localizer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			handler, err := NewSortHandlerWithOptions(tt.options...)

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

func TestCreateSortHandlerFunction(t *testing.T) {
	t.Parallel()

	// Arrange
	baseHandler := &BaseHandler{}
	executor := &testutils.FakeTaskExecutor{}
	applier := &testutils.FakeConfigApplier[*cmdconfig.SortConfig]{}
	validator := &testutils.FakeConfigValidator[*cmdconfig.SortConfig]{}
	localizer := testutils.NewFakeLocalizer()

	// Act
	handler := createSortHandler(baseHandler, executor, applier, validator, localizer)

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
	if handler.Config == nil {
		t.Error("Expected Config to be initialized")
	}
}

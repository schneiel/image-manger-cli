package task

import (
	"testing"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

// TestNewDedupTaskBuilder tests the creation of a new builder.
func TestNewDedupTaskBuilder(t *testing.T) {
	t.Parallel()
	builder := NewDedupTaskBuilder()
	if builder == nil {
		t.Fatal("NewDedupTaskBuilder should not return nil")
	}
	if builder.task == nil {
		t.Fatal("Builder should initialize with a task instance")
	}
	if builder.err != nil {
		t.Fatalf("Builder should initialize without error, got: %v", builder.err)
	}
}

// TestDedupTaskBuilder_Build_MissingConfig tests error handling for missing config.
func TestDedupTaskBuilder_Build_MissingConfig(t *testing.T) {
	t.Parallel()
	builder := NewDedupTaskBuilder()
	task, err := builder.Build()

	if err == nil {
		t.Fatal("Build should fail with missing config")
	}
	if task != nil {
		t.Fatal("Build should return nil task on error")
	}
	if err.Error() != "config is required" {
		t.Errorf("Expected error 'config is required', got %q", err.Error())
	}
}

// TestDedupTaskBuilder_WithNilConfig tests nil config validation.
func TestDedupTaskBuilder_WithNilConfig(t *testing.T) {
	t.Parallel()
	builder := NewDedupTaskBuilder().WithConfig(nil)
	task, err := builder.Build()

	if err == nil {
		t.Fatal("Build should fail with nil config")
	}
	if task != nil {
		t.Fatal("Build should return nil task on error")
	}
	if err.Error() != "config cannot be nil" {
		t.Errorf("Expected error 'config cannot be nil', got %q", err.Error())
	}
}

// TestDedupTaskBuilder_FluentInterface tests basic method chaining.
func TestDedupTaskBuilder_FluentInterface(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Deduplicator: config.DefaultDeduplicatorConfig(),
	}
	logger := &testutils.MockLogger{}

	// Test that chaining returns the builder
	builder := NewDedupTaskBuilder().WithConfig(cfg).WithLogger(logger)
	if builder == nil {
		t.Fatal("Fluent interface should return builder instance")
	}
}

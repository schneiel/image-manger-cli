package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	"github.com/schneiel/ImageManagerGo/internal/services"
)

func TestNewSortConfigApplier(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}

	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	if applier == nil {
		t.Fatal("Expected applier to be created, got nil")
	}
}

func TestNewSortConfigApplier_NilLogger(t *testing.T) {
	t.Parallel()
	_, err := services.NewSortConfigApplier(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "logger cannot be nil")
}

func TestSortConfigApplier_Apply_AllFields(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.SortConfig{
		Source:         "/test/source",
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	dest := &config.Config{}

	applier.Apply(src, dest)

	// Verify all fields were applied
	if dest.Sorter.Source != src.Source {
		t.Errorf("Expected Source %s, got %s", src.Source, dest.Sorter.Source)
	}

	if dest.Sorter.Destination != src.Destination {
		t.Errorf("Expected Destination %s, got %s", src.Destination, dest.Sorter.Destination)
	}

	if dest.Sorter.ActionStrategy != src.ActionStrategy {
		t.Errorf("Expected ActionStrategy %s, got %s", src.ActionStrategy, dest.Sorter.ActionStrategy)
	}
}

func TestSortConfigApplier_Apply_EmptyFields(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	// Create source with empty fields
	src := &cmdconfig.SortConfig{}

	// Create destination with existing values
	dest := &config.Config{
		Sorter: config.SorterConfig{
			Source:         "existing_source",
			Destination:    "existing_destination",
			ActionStrategy: "existing_action",
		},
	}

	applier.Apply(src, dest)

	// Verify existing values were preserved (not overwritten with empty values)
	if dest.Sorter.Source != "existing_source" {
		t.Errorf("Expected Source to remain 'existing_source', got %s", dest.Sorter.Source)
	}

	if dest.Sorter.Destination != "existing_destination" {
		t.Errorf("Expected Destination to remain 'existing_destination', got %s", dest.Sorter.Destination)
	}

	if dest.Sorter.ActionStrategy != "existing_action" {
		t.Errorf("Expected ActionStrategy to remain 'existing_action', got %s", dest.Sorter.ActionStrategy)
	}
}

func TestSortConfigApplier_Apply_PartialFields(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	// Create source with only some fields set
	src := &cmdconfig.SortConfig{
		Source: "/test/source",
		// Destination and ActionStrategy are empty
	}

	dest := &config.Config{}

	applier.Apply(src, dest)

	// Verify only set fields were applied
	if dest.Sorter.Source != src.Source {
		t.Errorf("Expected Source %s, got %s", src.Source, dest.Sorter.Source)
	}

	// Verify unset fields remain at zero values
	if dest.Sorter.Destination != "" {
		t.Errorf("Expected Destination to be empty, got %s", dest.Sorter.Destination)
	}

	if dest.Sorter.ActionStrategy != "" {
		t.Errorf("Expected ActionStrategy to be empty, got %s", dest.Sorter.ActionStrategy)
	}
}

func TestSortConfigApplier_Apply_LoggerCalls(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	debugCalls := 0
	mockLogger.DebugFunc = func(_ string) {
		debugCalls++
	}

	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.SortConfig{
		Source:         "/test/source",
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	dest := &config.Config{}

	applier.Apply(src, dest)

	// Verify logger was called for each field that was applied
	expectedCalls := 3 // Source, Destination, ActionStrategy
	if debugCalls != expectedCalls {
		t.Errorf("Expected %d debug calls, got %d", expectedCalls, debugCalls)
	}
}

func TestSortConfigApplier_Apply_NilSource(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	dest := &config.Config{
		Sorter: config.SorterConfig{
			Source: "existing_source",
		},
	}

	// Should panic when source is nil (current implementation doesn't handle nil)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil source, but none occurred")
		}
	}()

	applier.Apply(nil, dest)
}

func TestSortConfigApplier_Apply_NilDestination(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.SortConfig{
		Source: "/test/source",
	}

	// Should panic when destination is nil (current implementation doesn't handle nil)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil destination, but none occurred")
		}
	}()

	applier.Apply(src, nil)
}

func TestSortConfigApplier_Apply_WhitespaceFields(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewSortConfigApplier(mockLogger)
	require.NoError(t, err)

	// Create source with whitespace-only fields
	src := &cmdconfig.SortConfig{
		Source:         "   ",
		Destination:    "\t",
		ActionStrategy: "\n",
	}

	dest := &config.Config{
		Sorter: config.SorterConfig{
			Source:         "existing_source",
			Destination:    "existing_destination",
			ActionStrategy: "existing_action",
		},
	}

	applier.Apply(src, dest)

	// Verify whitespace-only fields were applied (they are not empty strings)
	if dest.Sorter.Source != "   " {
		t.Errorf("Expected Source to be '   ', got '%s'", dest.Sorter.Source)
	}

	if dest.Sorter.Destination != "\t" {
		t.Errorf("Expected Destination to be '\\t', got '%s'", dest.Sorter.Destination)
	}

	if dest.Sorter.ActionStrategy != "\n" {
		t.Errorf("Expected ActionStrategy to be '\\n', got '%s'", dest.Sorter.ActionStrategy)
	}
}

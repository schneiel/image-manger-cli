package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	"github.com/schneiel/ImageManagerGo/internal/services"
)

func TestNewDedupConfigApplier(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}

	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	if applier == nil {
		t.Fatal("Expected applier to be created, got nil")
	}
}

func TestNewDedupConfigApplier_NilLogger(t *testing.T) {
	t.Parallel()
	_, err := services.NewDedupConfigApplier(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "logger cannot be nil")
}

func TestDedupConfigApplier_Apply_AllFields(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.DedupConfig{
		Source:         "/test/source",
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
		TrashPath:      ".trash",
		Workers:        4,
		Threshold:      1,
	}

	dest := &config.Config{}

	applier.Apply(src, dest)

	// Verify all fields were applied
	if dest.Deduplicator.Source != src.Source {
		t.Errorf("Expected Source %s, got %s", src.Source, dest.Deduplicator.Source)
	}

	if dest.Deduplicator.ActionStrategy != src.ActionStrategy {
		t.Errorf("Expected ActionStrategy %s, got %s", src.ActionStrategy, dest.Deduplicator.ActionStrategy)
	}

	if dest.Deduplicator.KeepStrategy != src.KeepStrategy {
		t.Errorf("Expected KeepStrategy %s, got %s", src.KeepStrategy, dest.Deduplicator.KeepStrategy)
	}

	if dest.Deduplicator.TrashPath != src.TrashPath {
		t.Errorf("Expected TrashPath %s, got %s", src.TrashPath, dest.Deduplicator.TrashPath)
	}

	if dest.Deduplicator.Workers != src.Workers {
		t.Errorf("Expected Workers %d, got %d", src.Workers, dest.Deduplicator.Workers)
	}

	if dest.Deduplicator.Threshold != src.Threshold {
		t.Errorf("Expected Threshold %d, got %d", src.Threshold, dest.Deduplicator.Threshold)
	}
}

func TestDedupConfigApplier_Apply_EmptyFields(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	// Create source with empty fields
	src := &cmdconfig.DedupConfig{}

	// Create destination with existing values
	dest := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source:         "existing_source",
			ActionStrategy: "existing_action",
			KeepStrategy:   "existing_keep",
			TrashPath:      "existing_trash",
			Workers:        2,
			Threshold:      0,
		},
	}

	applier.Apply(src, dest)

	// Verify existing values were preserved (not overwritten with empty values)
	if dest.Deduplicator.Source != "existing_source" {
		t.Errorf("Expected Source to remain 'existing_source', got %s", dest.Deduplicator.Source)
	}

	if dest.Deduplicator.ActionStrategy != "existing_action" {
		t.Errorf("Expected ActionStrategy to remain 'existing_action', got %s", dest.Deduplicator.ActionStrategy)
	}

	if dest.Deduplicator.KeepStrategy != "existing_keep" {
		t.Errorf("Expected KeepStrategy to remain 'existing_keep', got %s", dest.Deduplicator.KeepStrategy)
	}

	if dest.Deduplicator.TrashPath != "existing_trash" {
		t.Errorf("Expected TrashPath to remain 'existing_trash', got %s", dest.Deduplicator.TrashPath)
	}

	if dest.Deduplicator.Workers != 2 {
		t.Errorf("Expected Workers to remain 2, got %d", dest.Deduplicator.Workers)
	}

	if dest.Deduplicator.Threshold != 0 {
		t.Errorf("Expected Threshold to remain 0, got %d", dest.Deduplicator.Threshold)
	}
}

func TestDedupConfigApplier_Apply_PartialFields(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	// Create source with only some fields set
	src := &cmdconfig.DedupConfig{
		Source:         "/test/source",
		ActionStrategy: "dryRun",
		// KeepStrategy, TrashPath, Workers, Threshold are empty/zero
	}

	dest := &config.Config{}

	applier.Apply(src, dest)

	// Verify only set fields were applied
	if dest.Deduplicator.Source != src.Source {
		t.Errorf("Expected Source %s, got %s", src.Source, dest.Deduplicator.Source)
	}

	if dest.Deduplicator.ActionStrategy != src.ActionStrategy {
		t.Errorf("Expected ActionStrategy %s, got %s", src.ActionStrategy, dest.Deduplicator.ActionStrategy)
	}

	// Verify unset fields remain at zero values
	if dest.Deduplicator.KeepStrategy != "" {
		t.Errorf("Expected KeepStrategy to be empty, got %s", dest.Deduplicator.KeepStrategy)
	}

	if dest.Deduplicator.TrashPath != "" {
		t.Errorf("Expected TrashPath to be empty, got %s", dest.Deduplicator.TrashPath)
	}

	if dest.Deduplicator.Workers != 0 {
		t.Errorf("Expected Workers to be 0, got %d", dest.Deduplicator.Workers)
	}

	if dest.Deduplicator.Threshold != 0 {
		t.Errorf("Expected Threshold to be 0, got %d", dest.Deduplicator.Threshold)
	}
}

func TestDedupConfigApplier_Apply_ZeroWorkers(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.DedupConfig{
		Workers: 0, // Zero workers should not be applied
	}

	dest := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Workers: 4, // Existing value
		},
	}

	applier.Apply(src, dest)

	// Verify workers was not changed (zero value should not overwrite)
	if dest.Deduplicator.Workers != 4 {
		t.Errorf("Expected Workers to remain 4, got %d", dest.Deduplicator.Workers)
	}
}

func TestDedupConfigApplier_Apply_ZeroThreshold(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.DedupConfig{
		Threshold: 0, // Zero threshold should be applied (>= 0 condition)
	}

	dest := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Threshold: 1, // Existing value
		},
	}

	applier.Apply(src, dest)

	// Verify threshold was applied (zero is valid for threshold)
	if dest.Deduplicator.Threshold != 0 {
		t.Errorf("Expected Threshold to be 0, got %d", dest.Deduplicator.Threshold)
	}
}

func TestDedupConfigApplier_Apply_NegativeThreshold(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.DedupConfig{
		Threshold: -1, // Negative threshold should not be applied
	}

	dest := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Threshold: 1, // Existing value
		},
	}

	applier.Apply(src, dest)

	// Verify threshold was not changed (negative value should not overwrite)
	if dest.Deduplicator.Threshold != 1 {
		t.Errorf("Expected Threshold to remain 1, got %d", dest.Deduplicator.Threshold)
	}
}

func TestDedupConfigApplier_Apply_LoggerCalls(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	debugCalls := 0
	mockLogger.DebugFunc = func(_ string) {
		debugCalls++
	}

	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.DedupConfig{
		Source:         "/test/source",
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
		TrashPath:      ".trash",
		Workers:        4,
		Threshold:      1,
	}

	dest := &config.Config{}

	applier.Apply(src, dest)

	// Verify logger was called for each field that was applied
	expectedCalls := 6 // Source, ActionStrategy, KeepStrategy, TrashPath, Workers, Threshold
	if debugCalls != expectedCalls {
		t.Errorf("Expected %d debug calls, got %d", expectedCalls, debugCalls)
	}
}

func TestDedupConfigApplier_Apply_NilSource(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	dest := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
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

func TestDedupConfigApplier_Apply_NilDestination(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	applier, err := services.NewDedupConfigApplier(mockLogger)
	require.NoError(t, err)

	src := &cmdconfig.DedupConfig{
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

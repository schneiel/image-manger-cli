package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewBaseHandler(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockConfig := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/destination",
		},
	}

	handler, err := NewBaseHandler(mockLogger, mockConfig)
	require.NoError(t, err)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.Logger != mockLogger {
		t.Error("Expected logger to be injected")
	}

	if handler.Config != mockConfig {
		t.Error("Expected config to be injected")
	}
}

func TestNewBaseHandler_NilLogger(t *testing.T) {
	t.Parallel()
	mockConfig := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/destination",
		},
	}

	handler, err := NewBaseHandler(nil, mockConfig)

	if err == nil {
		t.Fatal("Expected error for nil logger, got nil")
	}

	if handler != nil {
		t.Error("Expected handler to be nil when error occurs")
	}
}

func TestNewBaseHandler_NilConfig(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}

	handler, err := NewBaseHandler(mockLogger, nil)

	if err == nil {
		t.Fatal("Expected error for nil config, got nil")
	}

	if handler != nil {
		t.Error("Expected handler to be nil when error occurs")
	}
}

func TestNewBaseHandler_BothNil(t *testing.T) {
	t.Parallel()
	handler, err := NewBaseHandler(nil, nil)

	if err == nil {
		t.Fatal("Expected error for nil parameters, got nil")
	}

	if handler != nil {
		t.Error("Expected handler to be nil when error occurs")
	}
}

func TestBaseHandler_WithCompleteConfig(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockConfig := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source:         "/test/source",
			ActionStrategy: "dryRun",
			KeepStrategy:   "keepOldest",
		},
		Sorter: config.SorterConfig{
			Source:         "/test/source",
			Destination:    "/test/destination",
			ActionStrategy: "copy",
		},
		Files: config.FilesConfig{
			ApplicationLog: "test.log",
		},
		AllowedImageExtensions: []string{".jpg", ".png", ".jpeg"},
	}

	handler, err := NewBaseHandler(mockLogger, mockConfig)
	require.NoError(t, err)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.Logger != mockLogger {
		t.Error("Expected logger to be injected")
	}

	if handler.Config != mockConfig {
		t.Error("Expected config to be injected")
	}

	// Verify config fields are accessible
	if handler.Config.Deduplicator.Source != "/test/source" {
		t.Errorf("Expected source '/test/source', got '%s'", handler.Config.Deduplicator.Source)
	}

	if handler.Config.Sorter.Destination != "/test/destination" {
		t.Errorf("Expected destination '/test/destination', got '%s'", handler.Config.Sorter.Destination)
	}

	if handler.Config.Files.ApplicationLog != "test.log" {
		t.Errorf("Expected application log 'test.log', got '%s'", handler.Config.Files.ApplicationLog)
	}

	if len(handler.Config.AllowedImageExtensions) != 3 {
		t.Errorf("Expected 3 extensions, got %d", len(handler.Config.AllowedImageExtensions))
	}
}

func TestBaseHandler_WithMinimalConfig(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockConfig := &config.Config{}

	handler, err := NewBaseHandler(mockLogger, mockConfig)
	require.NoError(t, err)
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.Logger != mockLogger {
		t.Error("Expected logger to be injected")
	}

	if handler.Config != mockConfig {
		t.Error("Expected config to be injected")
	}
}

func TestBaseHandler_LoggerInteraction(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockConfig := &config.Config{}

	handler, err := NewBaseHandler(mockLogger, mockConfig)
	require.NoError(t, err)

	// Test that we can use the logger
	handler.Logger.Info("Test message")

	// Verify the logger was called (if we had a way to track calls)
	// For now, just verify the handler has access to the logger
	if handler.Logger == nil {
		t.Error("Expected logger to be available")
	}
}

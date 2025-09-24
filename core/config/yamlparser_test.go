package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAMLParser_ParseWithDefaults(t *testing.T) {
	t.Parallel()
	parser := NewYAMLParser()

	// Test with minimal YAML that doesn't specify log paths
	yamlData := []byte(`
files:
  applicationLog: "/tmp/app.log"
`)

	config, err := parser.Parse(yamlData)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	// Check that specified values are used
	if config.Files.ApplicationLog != "/tmp/app.log" {
		t.Errorf("Expected application log to be '/tmp/app.log', got '%s'", config.Files.ApplicationLog)
	}
}

func TestYAMLParser_Parse(t *testing.T) {
	t.Parallel()
	// Arrange
	parser := NewYAMLParser()
	data := []byte(`
deduplicator:
  actionStrategy: dryRun
  keepStrategy: keepOldest
  log: deduplicator.log
  trashPath: .trash
  workers: 4
  threshold: 1
sorter:
  actionStrategy: dryRun
  log: sorter.log
  date:
    strategyOrder:
      - exif
      - modTime
      - creationTime
    exifStrategies:
      - fieldName: DateTimeOriginal
        layout: "2006:01:02 15:04:05"
      - fieldName: DateTime
        layout: "2006:01:02 15:04:05"
      - fieldName: SubSecDateTimeOriginal
        layout: "2006:01:02 15:04:05.00"
      - fieldName: GPSDateStamp
        layout: "2006:01:02"
files:
  applicationLog: application.log
  dedupDryRunLog: dedup_dry_run_log.csv
  sortDryRunLog: sort_dry_run_log.csv
allowedImageExtensions:
  - .jpg
  - .jpeg
  - .png
  - .gif
`)

	// Act
	cfg, err := parser.Parse(data)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "dryRun", cfg.Deduplicator.ActionStrategy)
	assert.Equal(t, "keepOldest", cfg.Deduplicator.KeepStrategy)
	assert.Equal(t, "deduplicator.log", cfg.Deduplicator.Log)
	assert.Equal(t, ".trash", cfg.Deduplicator.TrashPath)
	assert.Equal(t, 4, cfg.Deduplicator.Workers)
	assert.Equal(t, 1, cfg.Deduplicator.Threshold)
	assert.Equal(t, "dryRun", cfg.Sorter.ActionStrategy)
	assert.Equal(t, "sorter.log", cfg.Sorter.Log)
	assert.Equal(t, []string{"exif", "modTime", "creationTime"}, cfg.Sorter.Date.StrategyOrder)
	assert.Len(t, cfg.Sorter.Date.ExifStrategies, 4)
	assert.Equal(t, "application.log", cfg.Files.ApplicationLog)
	assert.Equal(t, "dedup_dry_run_log.csv", cfg.Files.DedupDryRunLog)
	assert.Equal(t, "sort_dry_run_log.csv", cfg.Files.SortDryRunLog)
	assert.Equal(t, []string{".jpg", ".jpeg", ".png", ".gif"}, cfg.AllowedImageExtensions)
}

func TestYAMLParser_Parse_Error(t *testing.T) {
	t.Parallel()
	// Arrange
	parser := NewYAMLParser()
	data := []byte(`invalid yaml`)

	// Act
	cfg, err := parser.Parse(data)

	// Assert
	require.Error(t, err)
	assert.Nil(t, cfg)
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.NotNil(t, config.Deduplicator)
	assert.NotNil(t, config.Sorter)
	assert.NotNil(t, config.Files)
	assert.NotEmpty(t, config.AllowedImageExtensions)
}

func TestSorterConfig_DefaultValues(t *testing.T) {
	t.Parallel()
	config := DefaultSorterConfig()

	assert.Equal(t, "", config.ActionStrategy) // Empty - set by CLI flags or defaults to dryRun for safety
	assert.Equal(t, "sorter.log", config.Log)
	assert.NotNil(t, config.Date)
}

func TestDeduplicatorConfig_DefaultValues(t *testing.T) {
	t.Parallel()
	config := DefaultDeduplicatorConfig()

	assert.Equal(t, "", config.ActionStrategy) // Empty - set by CLI flags or defaults to dryRun for safety
	assert.Equal(t, "keepOldest", config.KeepStrategy)
	assert.Equal(t, "deduplicator.log", config.Log)
	assert.Equal(t, ".trash", config.TrashPath)
	assert.Positive(t, config.Workers)
	assert.Equal(t, 1, config.Threshold)
}

func TestDateConfig_DefaultValues(t *testing.T) {
	t.Parallel()
	config := DefaultDateConfig()

	assert.NotEmpty(t, config.StrategyOrder)
	assert.NotEmpty(t, config.ExifStrategies)
}

func TestFilesConfig_DefaultValues(t *testing.T) {
	t.Parallel()
	config := DefaultFilesConfig()

	assert.NotEmpty(t, config.DedupDryRunLog)
}

func TestConfig_Structure(t *testing.T) {
	t.Parallel()
	config := &Config{
		AllowedImageExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"},
	}

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.AllowedImageExtensions)
	assert.Contains(t, config.AllowedImageExtensions, ".jpg")
	assert.Contains(t, config.AllowedImageExtensions, ".jpeg")
	assert.Contains(t, config.AllowedImageExtensions, ".png")
}

func TestConfig_AllowedExtensions(t *testing.T) {
	t.Parallel()
	config := &Config{
		AllowedImageExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"},
	}

	assert.Contains(t, config.AllowedImageExtensions, ".jpg")
	assert.Contains(t, config.AllowedImageExtensions, ".jpeg")
	assert.Contains(t, config.AllowedImageExtensions, ".png")
	assert.Contains(t, config.AllowedImageExtensions, ".gif")
	assert.Contains(t, config.AllowedImageExtensions, ".bmp")
	assert.Contains(t, config.AllowedImageExtensions, ".tiff")
}

func TestConfig_InvalidExtensions(t *testing.T) {
	t.Parallel()
	config := &Config{
		AllowedImageExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"},
	}

	assert.NotContains(t, config.AllowedImageExtensions, ".pdf")
	assert.NotContains(t, config.AllowedImageExtensions, ".txt")
	assert.NotContains(t, config.AllowedImageExtensions, ".docx")
}

func TestConfig_EmptyAllowedExtensions(t *testing.T) {
	t.Parallel()
	config := &Config{
		AllowedImageExtensions: []string{},
	}

	assert.Empty(t, config.AllowedImageExtensions)
}

func TestConfig_NilLogger(t *testing.T) {
	t.Parallel()
	config := &Config{
		Logger: nil,
	}

	assert.Nil(t, config.Logger)
}

func TestConfig_EmptyDeduplicatorConfig(t *testing.T) {
	t.Parallel()
	config := &Config{
		Deduplicator: DeduplicatorConfig{},
	}

	assert.NotNil(t, config.Deduplicator)
}

func TestConfig_EmptySorterConfig(t *testing.T) {
	t.Parallel()
	config := &Config{
		Sorter: SorterConfig{},
	}

	assert.NotNil(t, config.Sorter)
}

func TestConfig_EmptyFilesConfig(t *testing.T) {
	t.Parallel()
	config := &Config{
		Files: FilesConfig{},
	}

	assert.NotNil(t, config.Files)
}

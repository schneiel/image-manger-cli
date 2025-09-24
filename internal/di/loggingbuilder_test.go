package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/time"
	"github.com/schneiel/ImageManagerGo/internal/di"
)

func TestNewLoggingBuilder(t *testing.T) {
	t.Parallel()
	// t.Parallel() removed due to race condition in i18n.NewBundleLocalizer().

	builder := di.NewLoggingBuilder()
	require.NotNil(t, builder)
	assert.IsType(t, &di.LoggingBuilder{}, builder)
}

func TestLoggingBuilder_Build_RequiresParameters(t *testing.T) {
	t.Parallel()

	builder := di.NewLoggingBuilder()
	_, err := builder.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Build() requires parameters")
}

func TestLoggingBuilder_BuildLogging_Success(t *testing.T) {
	t.Parallel()
	// t.Parallel() removed due to race condition in i18n.NewBundleLocalizer().

	builder := di.NewLoggingBuilder()

	// Create a test config with logging settings
	cfg := &config.Config{
		Sorter: config.SorterConfig{
			Log: "sorter.log",
		},
		Deduplicator: config.DeduplicatorConfig{
			Log: "dedup.log",
		},
	}

	// Create mock core dependencies
	core := &di.CoreDependencies{
		FileSystem:   filesystem.NewDefaultFileSystem(),
		TimeProvider: time.NewDefaultTimeProvider(),
		Localizer:    &MockLocalizer{},
	}

	logging, err := builder.BuildLogging(core, cfg)
	require.NoError(t, err)
	require.NotNil(t, logging)
	assert.NotNil(t, logging.SortLogger)
	assert.NotNil(t, logging.DedupLogger)
}

func TestLoggingBuilder_BuildLogging_WithDifferentLogFiles(t *testing.T) {
	t.Parallel()
	// t.Parallel() removed due to race condition in i18n.NewBundleLocalizer().

	builder := di.NewLoggingBuilder()

	tests := []struct {
		name       string
		sorterFile string
		dedupFile  string
	}{
		{
			name:       "different files",
			sorterFile: "sorter.log",
			dedupFile:  "dedup.log",
		},
		{
			name:       "same file",
			sorterFile: "app.log",
			dedupFile:  "app.log",
		},
		{
			name:       "console logging",
			sorterFile: "console",
			dedupFile:  "console",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Config{
				Sorter: config.SorterConfig{
					Log: tt.sorterFile,
				},
				Deduplicator: config.DeduplicatorConfig{
					Log: tt.dedupFile,
				},
			}

			core := &di.CoreDependencies{
				FileSystem:   filesystem.NewDefaultFileSystem(),
				TimeProvider: time.NewDefaultTimeProvider(),
				Localizer:    &MockLocalizer{},
			}

			logging, err := builder.BuildLogging(core, cfg)
			require.NoError(t, err)
			require.NotNil(t, logging)
			assert.NotNil(t, logging.SortLogger)
			assert.NotNil(t, logging.DedupLogger)
		})
	}
}

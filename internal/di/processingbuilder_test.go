package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/time"
	"github.com/schneiel/ImageManagerGo/internal/di"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewProcessingBuilder(t *testing.T) {
	t.Parallel()

	builder := di.NewProcessingBuilder()
	require.NotNil(t, builder)
	assert.IsType(t, &di.ProcessingBuilder{}, builder)
}

func TestProcessingBuilder_BuildProcessing_Success(t *testing.T) {
	t.Parallel()

	builder := di.NewProcessingBuilder()

	// Create test config with minimal required settings
	cfg := &config.Config{
		AllowedImageExtensions: []string{".jpg", ".png"},
		Sorter: config.SorterConfig{
			Date: config.DateConfig{
				StrategyOrder: []string{"exif"},
				ExifStrategies: []config.ExifConfig{
					{
						FieldName: "DateTime",
						Layout:    "2006:01:02 15:04:05",
					},
				},
			},
		},
		Deduplicator: config.DeduplicatorConfig{
			Workers:   4,
			Threshold: 1,
		},
	}

	// Create mock dependencies
	core := &di.CoreDependencies{
		FileSystem:   filesystem.NewDefaultFileSystem(),
		TimeProvider: time.NewDefaultTimeProvider(),
		Localizer:    &MockLocalizer{},
	}

	logging := &di.LoggingDependencies{
		SortLogger:  &testutils.MockLogger{},
		DedupLogger: &testutils.MockLogger{},
	}

	processing, err := builder.BuildProcessing(core, cfg, logging)
	require.NoError(t, err)
	require.NotNil(t, processing)

	// Verify that all components were created
	assert.NotNil(t, processing.ExifDecoder)
	assert.NotNil(t, processing.ExifReader)
	assert.NotNil(t, processing.DateProcessor)
	assert.NotNil(t, processing.ImageFinder)
	assert.NotNil(t, processing.ImageAnalyzer)
	assert.NotNil(t, processing.ImageProcessor)
	assert.NotNil(t, processing.SizeScanner)
	assert.NotNil(t, processing.PHasher)
	assert.NotNil(t, processing.DistanceGrouper)
}

func TestProcessingBuilder_BuildSorterActionStrategy_DryRun(t *testing.T) {
	t.Parallel()

	builder := di.NewProcessingBuilder()

	cfg := &config.Config{
		Sorter: config.SorterConfig{
			ActionStrategy: config.ActionStrategyDryRun,
		},
	}

	core := &di.CoreDependencies{
		FileSystem: filesystem.NewDefaultFileSystem(),
		Localizer:  &MockLocalizer{},
	}

	logging := &di.LoggingDependencies{
		SortLogger:  &testutils.MockLogger{},
		DedupLogger: &testutils.MockLogger{},
	}

	strategy, err := builder.BuildSorterActionStrategy(cfg, core, logging)
	require.NoError(t, err)
	assert.NotNil(t, strategy)
}

func TestProcessingBuilder_BuildSorterActionStrategy_Copy(t *testing.T) {
	t.Parallel()

	builder := di.NewProcessingBuilder()

	cfg := &config.Config{
		Sorter: config.SorterConfig{
			ActionStrategy: config.ActionStrategyCopy,
		},
	}

	core := &di.CoreDependencies{
		FileSystem: filesystem.NewDefaultFileSystem(),
		Localizer:  &MockLocalizer{},
	}

	logging := &di.LoggingDependencies{
		SortLogger:  &testutils.MockLogger{},
		DedupLogger: &testutils.MockLogger{},
	}

	strategy, err := builder.BuildSorterActionStrategy(cfg, core, logging)
	require.NoError(t, err)
	assert.NotNil(t, strategy)
}

func TestProcessingBuilder_BuildDedupActionStrategy_DryRun(t *testing.T) {
	t.Parallel()

	builder := di.NewProcessingBuilder()

	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			ActionStrategy: config.ActionStrategyDryRun,
		},
		Files: config.FilesConfig{
			DedupDryRunLog: "dedup.csv",
		},
	}

	core := &di.CoreDependencies{
		FileSystem: filesystem.NewDefaultFileSystem(),
		Localizer:  &MockLocalizer{},
	}

	logging := &di.LoggingDependencies{
		SortLogger:  &testutils.MockLogger{},
		DedupLogger: &testutils.MockLogger{},
	}

	strategy, err := builder.BuildDedupActionStrategy(cfg, core, logging)
	require.NoError(t, err)
	assert.NotNil(t, strategy)
}

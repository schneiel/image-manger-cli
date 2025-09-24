package task

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/processing/sort"
	"github.com/schneiel/ImageManagerGo/core/strategies/shared"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

// Test implementations of required interfaces

type mockImageProcessor struct{}

func (m *mockImageProcessor) Process(_ string) []image.Image {
	return []image.Image{}
}

type mockFileUtils struct{}

func (m *mockFileUtils) CopyFile(_, _ string) error {
	return nil
}

func (m *mockFileUtils) Exists(_ string) bool {
	return true
}

func (m *mockFileUtils) EnsureDir(_ string) error {
	return nil
}

type mockStrategy struct{}

func (m *mockStrategy) Execute(_, _ string) error {
	return nil
}

func (m *mockStrategy) GetResources() shared.ActionResource {
	return nil
}

func TestNewSortTaskBuilder(t *testing.T) {
	t.Parallel()

	builder := NewSortTaskBuilder()

	require.NotNil(t, builder)
	assert.NotNil(t, builder.task)
	assert.NoError(t, builder.err)
}

func TestSortTaskBuilder_WithConfig(t *testing.T) {
	t.Parallel()

	t.Run("Valid config", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		cfg := &config.Config{}

		result := builder.WithConfig(cfg)

		assert.Same(t, builder, result)
		assert.Equal(t, cfg, builder.task.config)
		assert.NoError(t, builder.err)
	})

	t.Run("Nil config", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()

		result := builder.WithConfig(nil)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.config)
		assert.EqualError(t, builder.err, "config cannot be nil")
	})

	t.Run("Builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		builder.err = errors.New("existing error")
		cfg := &config.Config{}

		result := builder.WithConfig(cfg)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.config)
		assert.EqualError(t, builder.err, "existing error")
	})
}

func TestSortTaskBuilder_WithLogger(t *testing.T) {
	t.Parallel()

	t.Run("Valid logger", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		logger := testutils.NewFakeLogger()

		result := builder.WithLogger(logger)

		assert.Same(t, builder, result)
		assert.Equal(t, logger, builder.task.logger)
		assert.NoError(t, builder.err)
	})

	t.Run("Nil logger", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()

		result := builder.WithLogger(nil)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.logger)
		assert.EqualError(t, builder.err, "logger cannot be nil")
	})

	t.Run("Builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		builder.err = errors.New("existing error")
		logger := testutils.NewFakeLogger()

		result := builder.WithLogger(logger)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.logger)
		assert.EqualError(t, builder.err, "existing error")
	})
}

func TestSortTaskBuilder_WithLocalizer(t *testing.T) {
	t.Parallel()

	t.Run("Valid localizer", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		localizer := testutils.NewFakeLocalizer()

		result := builder.WithLocalizer(localizer)

		assert.Same(t, builder, result)
		assert.Equal(t, localizer, builder.task.localizer)
		assert.NoError(t, builder.err)
	})

	t.Run("Nil localizer", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()

		result := builder.WithLocalizer(nil)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.localizer)
		assert.EqualError(t, builder.err, "localizer cannot be nil")
	})

	t.Run("Builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		builder.err = errors.New("existing error")
		localizer := testutils.NewFakeLocalizer()

		result := builder.WithLocalizer(localizer)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.localizer)
		assert.EqualError(t, builder.err, "existing error")
	})
}

func TestSortTaskBuilder_WithImageProcessor(t *testing.T) {
	t.Parallel()

	t.Run("Valid image processor", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		imageProcessor := &mockImageProcessor{}

		result := builder.WithImageProcessor(imageProcessor)

		assert.Same(t, builder, result)
		assert.Equal(t, imageProcessor, builder.task.imageProcessor)
		assert.NoError(t, builder.err)
	})

	t.Run("Nil image processor", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()

		result := builder.WithImageProcessor(nil)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.imageProcessor)
		assert.EqualError(t, builder.err, "imageProcessor cannot be nil")
	})

	t.Run("Builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		builder.err = errors.New("existing error")
		imageProcessor := &mockImageProcessor{}

		result := builder.WithImageProcessor(imageProcessor)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.imageProcessor)
		assert.EqualError(t, builder.err, "existing error")
	})
}

func TestSortTaskBuilder_WithFileUtils(t *testing.T) {
	t.Parallel()

	t.Run("Valid file utils", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		fileUtils := &mockFileUtils{}

		result := builder.WithFileUtils(fileUtils)

		assert.Same(t, builder, result)
		assert.Equal(t, fileUtils, builder.task.fileUtils)
		assert.NoError(t, builder.err)
	})

	t.Run("Nil file utils", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()

		result := builder.WithFileUtils(nil)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.fileUtils)
		assert.EqualError(t, builder.err, "fileUtils cannot be nil")
	})

	t.Run("Builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		builder.err = errors.New("existing error")
		fileUtils := &mockFileUtils{}

		result := builder.WithFileUtils(fileUtils)

		assert.Same(t, builder, result)
		assert.Nil(t, builder.task.fileUtils)
		assert.EqualError(t, builder.err, "existing error")
	})
}

func TestSortTaskBuilder_WithStrategyFactory(t *testing.T) {
	t.Parallel()

	t.Run("Valid action strategy", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		actionStrategy := &mockStrategy{}

		result := builder.WithStrategyFactory(func() sort.Strategy { return actionStrategy })

		assert.Same(t, builder, result)
		// actionStrategy field no longer exists - strategy is now created via factory function
		assert.NoError(t, builder.err)
	})

	t.Run("Nil action strategy", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()

		result := builder.WithStrategyFactory(nil)

		assert.Same(t, builder, result)
		// actionStrategy field no longer exists
		assert.EqualError(t, builder.err, "strategyFactory cannot be nil")
	})

	t.Run("Builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		builder.err = errors.New("existing error")
		actionStrategy := &mockStrategy{}

		result := builder.WithStrategyFactory(func() sort.Strategy { return actionStrategy })

		assert.Same(t, builder, result)
		// actionStrategy field no longer exists
		assert.EqualError(t, builder.err, "existing error")
	})
}

func TestSortTaskBuilder_Build(t *testing.T) {
	t.Parallel()

	t.Run("Complete valid build", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		cfg := &config.Config{}
		logger := testutils.NewFakeLogger()
		localizer := testutils.NewFakeLocalizer()
		imageProcessor := &mockImageProcessor{}
		fileUtils := &mockFileUtils{}
		actionStrategy := &mockStrategy{}

		task, err := builder.
			WithConfig(cfg).
			WithLogger(logger).
			WithLocalizer(localizer).
			WithImageProcessor(imageProcessor).
			WithFileUtils(fileUtils).
			WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
			Build()

		require.NoError(t, err)
		require.NotNil(t, task)
		assert.Equal(t, cfg, task.config)
		assert.Equal(t, logger, task.logger)
		assert.Equal(t, localizer, task.localizer)
		assert.Equal(t, imageProcessor, task.imageProcessor)
		assert.Equal(t, fileUtils, task.fileUtils)
		// actionStrategy field no longer exists - strategy is now created via factory function
	})

	t.Run("Build with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		builder.err = errors.New("existing error")

		task, err := builder.Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "existing error")
	})

	t.Run("Build missing config", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		logger := testutils.NewFakeLogger()
		localizer := testutils.NewFakeLocalizer()
		imageProcessor := &mockImageProcessor{}
		fileUtils := &mockFileUtils{}
		actionStrategy := &mockStrategy{}

		task, err := builder.
			WithLogger(logger).
			WithLocalizer(localizer).
			WithImageProcessor(imageProcessor).
			WithFileUtils(fileUtils).
			WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
			Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "config is required")
	})

	t.Run("Build missing logger", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		cfg := &config.Config{}
		localizer := testutils.NewFakeLocalizer()
		imageProcessor := &mockImageProcessor{}
		fileUtils := &mockFileUtils{}
		actionStrategy := &mockStrategy{}

		task, err := builder.
			WithConfig(cfg).
			WithLocalizer(localizer).
			WithImageProcessor(imageProcessor).
			WithFileUtils(fileUtils).
			WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
			Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "logger is required")
	})

	t.Run("Build missing localizer", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		cfg := &config.Config{}
		logger := testutils.NewFakeLogger()
		imageProcessor := &mockImageProcessor{}
		fileUtils := &mockFileUtils{}
		actionStrategy := &mockStrategy{}

		task, err := builder.
			WithConfig(cfg).
			WithLogger(logger).
			WithImageProcessor(imageProcessor).
			WithFileUtils(fileUtils).
			WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
			Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "localizer is required")
	})

	t.Run("Build missing imageProcessor", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		cfg := &config.Config{}
		logger := testutils.NewFakeLogger()
		localizer := testutils.NewFakeLocalizer()
		fileUtils := &mockFileUtils{}
		actionStrategy := &mockStrategy{}

		task, err := builder.
			WithConfig(cfg).
			WithLogger(logger).
			WithLocalizer(localizer).
			WithFileUtils(fileUtils).
			WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
			Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "imageProcessor is required")
	})

	t.Run("Build missing fileUtils", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		cfg := &config.Config{}
		logger := testutils.NewFakeLogger()
		localizer := testutils.NewFakeLocalizer()
		imageProcessor := &mockImageProcessor{}
		actionStrategy := &mockStrategy{}

		task, err := builder.
			WithConfig(cfg).
			WithLogger(logger).
			WithLocalizer(localizer).
			WithImageProcessor(imageProcessor).
			WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
			Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "fileUtils is required")
	})

	t.Run("Build missing actionStrategy", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		cfg := &config.Config{}
		logger := testutils.NewFakeLogger()
		localizer := testutils.NewFakeLocalizer()
		imageProcessor := &mockImageProcessor{}
		fileUtils := &mockFileUtils{}

		task, err := builder.
			WithConfig(cfg).
			WithLogger(logger).
			WithLocalizer(localizer).
			WithImageProcessor(imageProcessor).
			WithFileUtils(fileUtils).
			Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "strategyFactory is required")
	})
}

func TestSortTaskBuilder_ChainedCalls(t *testing.T) {
	t.Parallel()

	t.Run("Chain calls with error propagation", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()

		task, err := builder.
			WithConfig(nil). // This will set an error
			WithLogger(testutils.NewFakeLogger()).
			WithLocalizer(testutils.NewFakeLocalizer()).
			Build()

		assert.Nil(t, task)
		assert.EqualError(t, err, "config cannot be nil")
	})

	t.Run("Chain calls stopping at first error", func(t *testing.T) {
		t.Parallel()
		builder := NewSortTaskBuilder()
		logger := testutils.NewFakeLogger()

		// After setting error with nil config, logger should not be set
		result := builder.
			WithConfig(nil).   // Sets error
			WithLogger(logger) // Should not execute due to error

		assert.Nil(t, result.task.config)
		assert.Nil(t, result.task.logger) // Should not be set due to error
		assert.EqualError(t, result.err, "config cannot be nil")
	})
}

func TestSortTaskBuilder_BuilderPattern(t *testing.T) {
	t.Parallel()

	// Test that all methods return the same builder instance for chaining
	builder := NewSortTaskBuilder()
	cfg := &config.Config{}
	logger := testutils.NewFakeLogger()
	localizer := testutils.NewFakeLocalizer()
	imageProcessor := &mockImageProcessor{}
	fileUtils := &mockFileUtils{}
	actionStrategy := &mockStrategy{}

	result1 := builder.WithConfig(cfg)
	result2 := result1.WithLogger(logger)
	result3 := result2.WithLocalizer(localizer)
	result4 := result3.WithImageProcessor(imageProcessor)
	result5 := result4.WithFileUtils(fileUtils)
	result6 := result5.WithStrategyFactory(func() sort.Strategy { return actionStrategy })

	// All results should be the same builder instance
	assert.Same(t, builder, result1)
	assert.Same(t, builder, result2)
	assert.Same(t, builder, result3)
	assert.Same(t, builder, result4)
	assert.Same(t, builder, result5)
	assert.Same(t, builder, result6)
}

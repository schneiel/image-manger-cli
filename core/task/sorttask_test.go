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

// Local mock for ActionStrategy.
type mockActionStrategy struct {
	GetResourcesFunc func() shared.ActionResource
	ExecuteFunc      func(source, destination string) error
}

func (m *mockActionStrategy) GetResources() shared.ActionResource {
	if m.GetResourcesFunc != nil {
		return m.GetResourcesFunc()
	}
	return nil
}

func (m *mockActionStrategy) Execute(source, destination string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(source, destination)
	}
	return nil
}

func TestNewSortTask(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	imageProcessor := &testutils.MockImageProcessor{}
	fileUtils := &testutils.MockFileUtils{}
	actionStrategy := &mockActionStrategy{}

	task, err := NewSortTaskBuilder().
		WithConfig(cfg).
		WithLogger(logger).
		WithLocalizer(localizer).
		WithImageProcessor(imageProcessor).
		WithFileUtils(fileUtils).
		WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
		Build()

	require.NoError(t, err)

	assert.NotNil(t, task)
	assert.Equal(t, cfg, task.config)
	assert.Equal(t, logger, task.logger)
	assert.Equal(t, localizer, task.localizer)
	assert.Equal(t, imageProcessor, task.imageProcessor)
	assert.Equal(t, fileUtils, task.fileUtils)
	// actionStrategy field no longer exists - strategy is now created via factory function
}

func TestSortTask_Run_NoImages(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/dest",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	imageProcessor := &testutils.MockImageProcessor{}
	fileUtils := &testutils.MockFileUtils{}
	actionStrategy := &mockActionStrategy{}

	// Mock image processor to return no images
	imageProcessor.ProcessFunc = func(_ string) []image.Image {
		return []image.Image{}
	}

	task, buildErr := NewSortTaskBuilder().
		WithConfig(cfg).
		WithLogger(logger).
		WithLocalizer(localizer).
		WithImageProcessor(imageProcessor).
		WithFileUtils(fileUtils).
		WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
		Build()

	require.NoError(t, buildErr)

	err := task.Run()

	require.NoError(t, err)
}

func TestSortTask_Run_WithImages_NoResources(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/dest",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	imageProcessor := &testutils.MockImageProcessor{}
	fileUtils := &testutils.MockFileUtils{}
	actionStrategy := &mockActionStrategy{}

	// Mock image processor to return some images
	testImages := []image.Image{
		{FilePath: "/test/image1.jpg"},
		{FilePath: "/test/image2.jpg"},
	}
	imageProcessor.ProcessFunc = func(_ string) []image.Image {
		return testImages
	}

	// Mock action strategy with no resources
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return nil
	}

	executeCallCount := 0
	actionStrategy.ExecuteFunc = func(_, _ string) error {
		executeCallCount++
		return nil
	}

	task, buildErr := NewSortTaskBuilder().
		WithConfig(cfg).
		WithLogger(logger).
		WithLocalizer(localizer).
		WithImageProcessor(imageProcessor).
		WithFileUtils(fileUtils).
		WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
		Build()

	require.NoError(t, buildErr)

	err := task.Run()

	require.NoError(t, err)
	assert.Equal(t, 2, executeCallCount) // Should execute for both images
}

func TestSortTask_Run_WithImages_WithResources(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/dest",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	imageProcessor := &testutils.MockImageProcessor{}
	fileUtils := &testutils.MockFileUtils{}
	actionStrategy := &mockActionStrategy{}

	// Mock image processor to return some images
	testImages := []image.Image{
		{FilePath: "/test/image1.jpg"},
	}
	imageProcessor.ProcessFunc = func(_ string) []image.Image {
		return testImages
	}

	// Mock action strategy with resources
	mockResource := &testutils.MockResource{}
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return mockResource
	}

	setupCalled := false
	teardownCalled := false
	mockResource.SetupFunc = func() error {
		setupCalled = true
		return nil
	}
	mockResource.TeardownFunc = func() error {
		teardownCalled = true
		return nil
	}

	actionStrategy.ExecuteFunc = func(_, _ string) error {
		return nil
	}

	task, buildErr := NewSortTaskBuilder().
		WithConfig(cfg).
		WithLogger(logger).
		WithLocalizer(localizer).
		WithImageProcessor(imageProcessor).
		WithFileUtils(fileUtils).
		WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
		Build()

	require.NoError(t, buildErr)

	err := task.Run()

	require.NoError(t, err)
	assert.True(t, setupCalled)
	assert.True(t, teardownCalled)
}

func TestSortTask_Run_ResourceSetupError(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/dest",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	imageProcessor := &testutils.MockImageProcessor{}
	fileUtils := &testutils.MockFileUtils{}
	actionStrategy := &mockActionStrategy{}

	// Mock image processor to return some images
	testImages := []image.Image{
		{FilePath: "/test/image1.jpg"},
	}
	imageProcessor.ProcessFunc = func(_ string) []image.Image {
		return testImages
	}

	// Mock action strategy with resources that fail setup
	mockResource := &testutils.MockResource{}
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return mockResource
	}

	expectedError := errors.New("setup failed")
	mockResource.SetupFunc = func() error {
		return expectedError
	}

	task, buildErr := NewSortTaskBuilder().
		WithConfig(cfg).
		WithLogger(logger).
		WithLocalizer(localizer).
		WithImageProcessor(imageProcessor).
		WithFileUtils(fileUtils).
		WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
		Build()

	require.NoError(t, buildErr)

	err := task.Run()

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestSortTask_InterfaceCompliance(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	imageProcessor := &testutils.MockImageProcessor{}
	fileUtils := &testutils.MockFileUtils{}
	actionStrategy := &mockActionStrategy{}

	task, buildErr := NewSortTaskBuilder().
		WithConfig(cfg).
		WithLogger(logger).
		WithLocalizer(localizer).
		WithImageProcessor(imageProcessor).
		WithFileUtils(fileUtils).
		WithStrategyFactory(func() sort.Strategy { return actionStrategy }).
		Build()

	require.NoError(t, buildErr)

	// Verify that SortTask implements the Task interface
	var _ Task = task
	assert.Implements(t, (*Task)(nil), task)
}

package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
	"github.com/schneiel/ImageManagerGo/internal/di"
)

// Simple mock task for testing.
type mockTask struct {
	runFunc func() error
}

func (m *mockTask) Run() error {
	if m.runFunc != nil {
		return m.runFunc()
	}
	return nil
}

func TestNewHandlersBuilder(t *testing.T) {
	t.Parallel()

	builder := di.NewHandlersBuilder()
	assert.NotNil(t, builder)
}

func TestHandlersBuilder_Build_RequiresParameters(t *testing.T) {
	t.Parallel()

	builder := di.NewHandlersBuilder()
	_, err := builder.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Build() requires parameters")
}

func TestHandlersBuilder_BuildHandlers_Success(t *testing.T) {
	t.Parallel()

	// Setup infrastructure mocks for CLI layer testing (appropriate pattern)
	mockLocalizer := &testutils.MockLocalizer{}
	mockLogger := &testutils.MockLogger{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &mockTask{}
	mockDedupTask := &mockTask{}

	// Create dependencies
	core := &di.CoreDependencies{
		Localizer: mockLocalizer,
		FileUtils: mockFileUtils,
	}

	logging := &di.LoggingDependencies{
		SortLogger:  mockLogger,
		DedupLogger: mockLogger,
	}

	cfg := &config.Config{}
	args := []string{"test", "--language", "en"}

	builder := di.NewHandlersBuilder()
	deps, err := builder.BuildHandlers(args, core, cfg, logging, mockSortTask, mockDedupTask)

	require.NoError(t, err)
	assert.NotNil(t, deps)
	assert.NotNil(t, deps.TaskExecutor)
	assert.NotNil(t, deps.SortConfigApplier)
	assert.NotNil(t, deps.SortConfigValidator)
	assert.NotNil(t, deps.DedupConfigApplier)
	assert.NotNil(t, deps.DedupConfigValidator)
	assert.NotNil(t, deps.CommandLocalizer)
	assert.NotNil(t, deps.SortFlagSetup)
	assert.NotNil(t, deps.DedupFlagSetup)
	assert.NotNil(t, deps.SortHandler)
	assert.NotNil(t, deps.DedupHandler)
}

func TestHandlersBuilder_BuildHandlers_GermanLocalizer(t *testing.T) {
	t.Parallel()

	// Setup infrastructure mocks for CLI layer testing (appropriate pattern)
	mockLocalizer := &testutils.MockLocalizer{}
	mockLogger := &testutils.MockLogger{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &mockTask{}
	mockDedupTask := &mockTask{}

	// Create dependencies
	core := &di.CoreDependencies{
		Localizer: mockLocalizer,
		FileUtils: mockFileUtils,
	}

	logging := &di.LoggingDependencies{
		SortLogger:  mockLogger,
		DedupLogger: mockLogger,
	}

	cfg := &config.Config{}
	args := []string{"test", "--language", "de"}

	builder := di.NewHandlersBuilder()
	deps, err := builder.BuildHandlers(args, core, cfg, logging, mockSortTask, mockDedupTask)

	require.NoError(t, err)
	assert.NotNil(t, deps)
	assert.NotNil(t, deps.CommandLocalizer)
}

func TestHandlersBuilder_BuildHandlers_DefaultToEnglish(t *testing.T) {
	t.Parallel()

	// Setup infrastructure mocks for CLI layer testing (appropriate pattern)
	mockLocalizer := &testutils.MockLocalizer{}
	mockLogger := &testutils.MockLogger{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &mockTask{}
	mockDedupTask := &mockTask{}

	// Create dependencies
	core := &di.CoreDependencies{
		Localizer: mockLocalizer,
		FileUtils: mockFileUtils,
	}

	logging := &di.LoggingDependencies{
		SortLogger:  mockLogger,
		DedupLogger: mockLogger,
	}

	cfg := &config.Config{}
	args := []string{"test", "--language", "fr"} // Unsupported language

	builder := di.NewHandlersBuilder()
	deps, err := builder.BuildHandlers(args, core, cfg, logging, mockSortTask, mockDedupTask)

	require.NoError(t, err)
	assert.NotNil(t, deps)
	assert.NotNil(t, deps.CommandLocalizer)
}

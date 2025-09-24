package task

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/processing/dedup"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupaction"
	"github.com/schneiel/ImageManagerGo/core/strategies/shared"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

// mockDedupActionStrategy implements imagededuplicator.ActionStrategy for testing.
type mockDedupActionStrategy struct {
	GetResourcesFunc func() shared.ActionResource
	ExecuteFunc      func(source, target *image.Image) error
}

func (m *mockDedupActionStrategy) GetResources() shared.ActionResource {
	if m.GetResourcesFunc != nil {
		return m.GetResourcesFunc()
	}
	return nil
}

func (m *mockDedupActionStrategy) Execute(source, target *image.Image) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(source, target)
	}
	return nil
}

// mockScanner implements dedup.Scanner for testing.
type mockScanner struct {
	ScanFunc func(_ string) (dedup.FileGroup, error)
}

func (m *mockScanner) Scan(rootPath string) (dedup.FileGroup, error) {
	if m.ScanFunc != nil {
		return m.ScanFunc(rootPath)
	}
	return nil, nil
}

// mockHasher implements dedup.Hasher for testing.
type mockHasher struct {
	HashFilesFunc func(_ []string) ([]*image.HashInfo, error)
}

func (m *mockHasher) HashFiles(files []string) ([]*image.HashInfo, error) {
	if m.HashFilesFunc != nil {
		return m.HashFilesFunc(files)
	}
	return nil, nil
}

// mockGrouper implements dedup.Grouper for testing.
type mockGrouper struct {
	GroupFunc func(_ []*image.HashInfo) ([]dedup.DuplicateGroup, error)
}

func (m *mockGrouper) Group(hashes []*image.HashInfo) ([]dedup.DuplicateGroup, error) {
	if m.GroupFunc != nil {
		return m.GroupFunc(hashes)
	}
	return nil, nil
}

// mockKeepFunc creates a mock keep function for testing.
func mockKeepFunc(files []string) (string, []string) {
	if len(files) > 0 {
		return files[0], files[1:]
	}
	return "", nil
}

func TestNewDedupTask(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	assert.NotNil(t, task)
	assert.Equal(t, cfg, task.config)
	assert.Equal(t, logger, task.logger)
	assert.Equal(t, localizer, task.localizer)
	assert.Equal(t, filesystem, task.filesystem)
	assert.Equal(t, scanner, task.scanner)
	assert.Equal(t, hasher, task.hasher)
	assert.Equal(t, grouper, task.grouper)
	// actionStrategy field no longer exists - strategy is now created via factory function
	assert.NotNil(t, task.groupFlattener)
}

func TestNewDedupTask_DefaultTrashPath(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
			// TrashPath not set
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	expectedTrashPath := filepath.Join(string(filepath.Separator), "test", "source", ".trash")
	assert.Equal(t, expectedTrashPath, task.config.Deduplicator.TrashPath)
}

func TestNewDedupTask_PanicOnNilDependencies(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	// Test error on nil logger
	_, err := NewDedupTask(cfg, nil, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "logger cannot be nil")

	// Test error on nil localizer
	_, err = NewDedupTask(cfg, logger, nil, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "localizer cannot be nil")

	// Test error on nil filesystem
	_, err = NewDedupTask(cfg, logger, localizer, nil, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "filesystem cannot be nil")

	// Test error on nil scanner
	_, err = NewDedupTask(cfg, logger, localizer, filesystem, nil, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scanner cannot be nil")

	// Test error on nil hasher
	_, err = NewDedupTask(cfg, logger, localizer, filesystem, scanner, nil, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hasher cannot be nil")

	// Test error on nil grouper
	_, err = NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, nil, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grouper cannot be nil")

	// Test error on nil keepFunc
	_, err = NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, nil, func() dedupaction.Strategy { return actionStrategy })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "keepFunc cannot be nil")

	// Test error on nil actionStrategy factory
	_, err = NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "strategyFactory cannot be nil")
}

func TestDedupTask_Run_Success(t *testing.T) {
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	// Mock scanner to return potential groups
	potentialGroups := dedup.FileGroup{
		{"/test/file1.jpg", "/test/file2.jpg"},
		{"/test/file3.jpg"},
	}
	scanner.ScanFunc = func(_ string) (dedup.FileGroup, error) {
		return potentialGroups, nil
	}

	// Mock hasher to return hash results
	hashResults := []*image.HashInfo{
		{FilePath: "/test/file1.jpg", Hash: nil},
		{FilePath: "/test/file2.jpg", Hash: nil},
		{FilePath: "/test/file3.jpg", Hash: nil},
	}
	hasher.HashFilesFunc = func(_ []string) ([]*image.HashInfo, error) {
		return hashResults, nil
	}

	// Mock grouper to return duplicate groups
	duplicateGroups := []dedup.DuplicateGroup{
		dedup.DuplicateGroup([]string{"/test/file1.jpg", "/test/file2.jpg"}),
	}
	grouper.GroupFunc = func(_ []*image.HashInfo) ([]dedup.DuplicateGroup, error) {
		return duplicateGroups, nil
	}

	// Custom keep function for this test
	customKeepFunc := func(files []string) (string, []string) {
		return files[0], files[1:]
	}

	// Mock action strategy with no resources
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return nil
	}

	executeCallCount := 0
	actionStrategy.ExecuteFunc = func(_, _ *image.Image) error {
		executeCallCount++
		return nil
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return key
	}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, customKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	err = task.Run()

	require.NoError(t, err)
	assert.Equal(t, 1, executeCallCount) // Should execute once for the duplicate
}

func TestDedupTask_Run_ScanError(t *testing.T) {
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	// Mock scanner to return error
	expectedError := errors.New("scan failed")
	scanner.ScanFunc = func(_ string) (dedup.FileGroup, error) {
		return nil, expectedError
	}

	// Mock action strategy with no resources
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return nil
	}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	err = task.Run()

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestDedupTask_Run_HashError(t *testing.T) {
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	// Mock scanner to return potential groups
	potentialGroups := dedup.FileGroup{
		{"/test/file1.jpg"},
	}
	scanner.ScanFunc = func(_ string) (dedup.FileGroup, error) {
		return potentialGroups, nil
	}

	// Mock hasher to return error
	expectedError := errors.New("hash failed")
	hasher.HashFilesFunc = func(_ []string) ([]*image.HashInfo, error) {
		return nil, expectedError
	}

	// Mock action strategy with no resources
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return nil
	}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	err = task.Run()

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestDedupTask_Run_GroupError(t *testing.T) {
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	// Mock scanner to return potential groups
	potentialGroups := dedup.FileGroup{
		{"/test/file1.jpg"},
	}
	scanner.ScanFunc = func(_ string) (dedup.FileGroup, error) {
		return potentialGroups, nil
	}

	// Mock hasher to return hash results
	hashResults := []*image.HashInfo{
		{FilePath: "/test/file1.jpg", Hash: nil},
	}
	hasher.HashFilesFunc = func(_ []string) ([]*image.HashInfo, error) {
		return hashResults, nil
	}

	// Mock grouper to return error
	expectedError := errors.New("group failed")
	grouper.GroupFunc = func(_ []*image.HashInfo) ([]dedup.DuplicateGroup, error) {
		return nil, expectedError
	}

	// Mock action strategy with no resources
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return nil
	}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	err = task.Run()

	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestDedupTask_Run_ResourceSetupError(t *testing.T) {
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	// Mock action strategy with resources that fail setup
	mockResource := &testutils.MockResource{}
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return mockResource
	}

	expectedError := errors.New("setup failed")
	mockResource.SetupFunc = func() error {
		return expectedError
	}

	// Mock localizer
	localizer.TranslateFunc = func(_ string, _ ...map[string]interface{}) string {
		return "ActionStrategyError"
	}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	err = task.Run()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ActionStrategyError")
	assert.Contains(t, err.Error(), "setup failed")
}

func TestDedupTask_Run_WithResources(t *testing.T) {
	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	// Mock scanner to return potential groups
	potentialGroups := dedup.FileGroup{
		{"/test/file1.jpg"},
	}
	scanner.ScanFunc = func(_ string) (dedup.FileGroup, error) {
		return potentialGroups, nil
	}

	// Mock hasher to return hash results
	hashResults := []*image.HashInfo{
		{FilePath: "/test/file1.jpg", Hash: nil},
	}
	hasher.HashFilesFunc = func(_ []string) ([]*image.HashInfo, error) {
		return hashResults, nil
	}

	// Mock grouper to return no duplicate groups
	grouper.GroupFunc = func(_ []*image.HashInfo) ([]dedup.DuplicateGroup, error) {
		return []dedup.DuplicateGroup{}, nil
	}

	// Mock action strategy with no resources
	actionStrategy.GetResourcesFunc = func() shared.ActionResource {
		return nil
	}

	// Mock localizer
	localizer.TranslateFunc = func(_ string, _ ...map[string]interface{}) string {
		return "SummaryNoDuplicates"
	}

	// Track logger calls
	infoCallCount := 0
	logger.InfoFunc = func(_ string) {
		infoCallCount++
	}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	err = task.Run()

	require.NoError(t, err)
	assert.Equal(t, 1, infoCallCount) // Should log "no duplicates" message
}

func TestDedupTask_InterfaceCompliance(t *testing.T) {
	cfg := &config.Config{}
	logger := &testutils.MockLogger{}
	localizer := &testutils.MockLocalizer{}
	filesystem := &testutils.MockFileSystem{}
	scanner := &mockScanner{}
	hasher := &mockHasher{}
	grouper := &mockGrouper{}
	actionStrategy := &mockDedupActionStrategy{}

	task, err := NewDedupTask(cfg, logger, localizer, filesystem, scanner, hasher, grouper, mockKeepFunc, func() dedupaction.Strategy { return actionStrategy })
	require.NoError(t, err)

	// Verify that DedupTask implements the Task interface
	var _ Task = task
	assert.Implements(t, (*Task)(nil), task)
}

package dedupaction

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/strategies/shared"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultMoveToTrashStrategy(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	trashPath := "/tmp/trash"

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
		TrashPath: trashPath,
	}

	strategy, err := NewDefaultMoveToTrashStrategy(config)
	require.NoError(t, err)

	assert.NotNil(t, strategy)
	assert.Equal(t, mockLogger, strategy.logger)
	assert.Equal(t, mockLocalizer, strategy.localizer)
	assert.NotNil(t, strategy.fileSystem)
	assert.NotNil(t, strategy.trashResource)
}

func TestNewDefaultMoveToTrashStrategyWithFilesystem(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileSystem := &testutils.MockFileSystem{}
	trashPath := "/tmp/trash"

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
		TrashPath: trashPath,
	}

	strategy, err := NewDefaultMoveToTrashStrategyWithFilesystem(config, mockFileSystem)
	require.NoError(t, err)

	assert.NotNil(t, strategy)
	assert.Equal(t, mockLogger, strategy.logger)
	assert.Equal(t, mockLocalizer, strategy.localizer)
	assert.Equal(t, mockFileSystem, strategy.fileSystem)
	assert.NotNil(t, strategy.trashResource)
}

func TestNewDefaultMoveToTrashStrategyWithFilesystem_NilFilesystem(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	trashPath := "/tmp/trash"

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
		TrashPath: trashPath,
	}

	_, err := NewDefaultMoveToTrashStrategyWithFilesystem(config, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fileSystem cannot be nil")
}

func TestDefaultMoveToTrashStrategy_Execute_Success(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileSystem := &testutils.MockFileSystem{}
	trashPath := "/tmp/trash"

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "MovingFile":
			return "Moving file"
		default:
			return key
		}
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	mockFileSystem.RenameFunc = func(_, _ string) error {
		return nil
	}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
		TrashPath: trashPath,
	}

	strategy, err := NewDefaultMoveToTrashStrategyWithFilesystem(config, mockFileSystem)
	require.NoError(t, err)

	// Create test images
	original := &image.Image{
		FilePath:         "/path/to/original.jpg",
		OriginalFileName: "original.jpg",
		Date:             time.Now(),
	}
	duplicate := &image.Image{
		FilePath:         "/path/to/duplicate.jpg",
		OriginalFileName: "duplicate.jpg",
		Date:             time.Now(),
	}

	err = strategy.Execute(duplicate, original)

	require.NoError(t, err)
}

func TestDefaultMoveToTrashStrategy_Execute_NilDuplicate(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{
		TranslateFunc: func(key string, _ ...map[string]interface{}) string {
			return key
		},
	}
	mockFileSystem := &testutils.MockFileSystem{}
	trashPath := "/tmp/trash"

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
		TrashPath: trashPath,
	}

	strategy, err := NewDefaultMoveToTrashStrategyWithFilesystem(config, mockFileSystem)
	require.NoError(t, err)

	// Create test original image
	original := &image.Image{
		FilePath:         "/path/to/original.jpg",
		OriginalFileName: "original.jpg",
		Date:             time.Now(),
	}

	// Should handle nil duplicate gracefully without panicking
	err = strategy.Execute(nil, original)

	// Should return an error when duplicate is nil
	if err == nil {
		t.Error("Expected error when duplicate is nil")
	}
}

func TestDefaultMoveToTrashStrategy_Execute_RenameError(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileSystem := &testutils.MockFileSystem{}
	trashPath := "/tmp/trash"

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "MovingFile":
			return "Moving file"
		case "MovingFileError":
			return "Error moving file"
		default:
			return key
		}
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	mockLogger.ErrorfFunc = func(_ string, _ ...interface{}) {
		// Verify error logging calls
	}

	renameError := errors.New("rename failed")
	mockFileSystem.RenameFunc = func(_, _ string) error {
		return renameError
	}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
		TrashPath: trashPath,
	}

	strategy, err := NewDefaultMoveToTrashStrategyWithFilesystem(config, mockFileSystem)
	require.NoError(t, err)

	// Create test duplicate image
	duplicate := &image.Image{
		FilePath:         "/path/to/duplicate.jpg",
		OriginalFileName: "duplicate.jpg",
		Date:             time.Now(),
	}

	err = strategy.Execute(duplicate, nil)

	require.Error(t, err)
	assert.Equal(t, renameError, err)
}

func TestDefaultMoveToTrashStrategy_GetResources(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileSystem := &testutils.MockFileSystem{}
	trashPath := "/tmp/trash"

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
		TrashPath: trashPath,
	}

	strategy, err := NewDefaultMoveToTrashStrategyWithFilesystem(config, mockFileSystem)
	require.NoError(t, err)

	resources := strategy.GetResources()

	assert.NotNil(t, resources)
	assert.Equal(t, strategy.trashResource, resources)
}

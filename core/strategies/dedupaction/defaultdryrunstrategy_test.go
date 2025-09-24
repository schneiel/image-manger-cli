package dedupaction

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/strategies/shared"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultDryRunStrategy(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
	}

	strategy, err := NewDefaultDryRunStrategy(config)
	require.NoError(t, err)

	assert.NotNil(t, strategy)
	assert.Equal(t, mockLogger, strategy.logger)
	assert.Equal(t, mockLocalizer, strategy.localizer)
}

func TestDefaultDryRunStrategy_Execute_Success(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "DryRunWouldMoveFile":
			return "Dry run would move file"
		default:
			return key
		}
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
	}

	strategy, err := NewDefaultDryRunStrategy(config)
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

	err = strategy.Execute(original, duplicate)

	require.NoError(t, err)
}

func TestDefaultDryRunStrategy_Execute_WithNilLocalizer(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: nil,
	}

	_, err := NewDefaultDryRunStrategy(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "localizer cannot be nil")
}

func TestDefaultDryRunStrategy_Execute_WithNilImages(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "DryRunWouldMoveFile":
			return "Dry run would move file"
		default:
			return key
		}
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
	}

	strategy, err := NewDefaultDryRunStrategy(config)
	require.NoError(t, err)

	// Test with nil images
	err = strategy.Execute(nil, nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "original and duplicate images cannot be nil")
}

func TestDefaultDryRunStrategy_Execute_WithNilOriginal(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "DryRunWouldMoveFile":
			return "Dry run would move file"
		default:
			return key
		}
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
	}

	strategy, err := NewDefaultDryRunStrategy(config)
	require.NoError(t, err)

	// Create test duplicate image
	duplicate := &image.Image{
		FilePath:         "/path/to/duplicate.jpg",
		OriginalFileName: "duplicate.jpg",
		Date:             time.Now(),
	}

	err = strategy.Execute(nil, duplicate)

	require.Error(t, err)
	require.Contains(t, err.Error(), "original and duplicate images cannot be nil")
}

func TestDefaultDryRunStrategy_Execute_WithNilDuplicate(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "DryRunWouldMoveFile":
			return "Dry run would move file"
		default:
			return key
		}
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
	}

	strategy, err := NewDefaultDryRunStrategy(config)
	require.NoError(t, err)

	// Create test original image
	original := &image.Image{
		FilePath:         "/path/to/original.jpg",
		OriginalFileName: "original.jpg",
		Date:             time.Now(),
	}

	err = strategy.Execute(original, nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "original and duplicate images cannot be nil")
}

func TestDefaultDryRunStrategy_Execute_WithSpecialCharacters(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "DryRunWouldMoveFile":
			return "Dry run would move file"
		default:
			return key
		}
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
	}

	strategy, err := NewDefaultDryRunStrategy(config)
	require.NoError(t, err)

	// Create test images with special characters in paths
	original := &image.Image{
		FilePath:         "/path/to/my photos/original (1).jpg",
		OriginalFileName: "original (1).jpg",
		Date:             time.Now(),
	}
	duplicate := &image.Image{
		FilePath:         "/path/with spaces/duplicate file.jpg",
		OriginalFileName: "duplicate file.jpg",
		Date:             time.Now(),
	}

	err = strategy.Execute(original, duplicate)

	require.NoError(t, err)
}

func TestDefaultDryRunStrategy_GetResources(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	config := shared.ActionConfig{
		Logger:    mockLogger,
		Localizer: mockLocalizer,
	}

	strategy, err := NewDefaultDryRunStrategy(config)
	require.NoError(t, err)

	resources := strategy.GetResources()

	// Dry run strategy should return nil resources
	assert.Nil(t, resources)
}

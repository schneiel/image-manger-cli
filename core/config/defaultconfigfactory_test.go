package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfigFactory(t *testing.T) {
	t.Parallel()

	// Act
	config := DefaultConfig()

	// Assert
	assert.NotNil(t, config)
	assert.NotNil(t, config.Deduplicator)
	assert.NotNil(t, config.Sorter)
	assert.NotNil(t, config.Files)
	assert.NotEmpty(t, config.AllowedImageExtensions)
	assert.Equal(t, []string{".jpg", ".jpeg", ".png", ".gif"}, config.AllowedImageExtensions)
}

func TestDefaultConfig_Consistency(t *testing.T) {
	t.Parallel()

	// Act - Multiple calls should return equivalent configurations
	config1 := DefaultConfig()
	config2 := DefaultConfig()

	// Assert - Both calls should return equivalent configurations
	assert.Equal(t, config1.AllowedImageExtensions, config2.AllowedImageExtensions)
	assert.Equal(t, config1.Deduplicator, config2.Deduplicator)
	assert.Equal(t, config1.Sorter, config2.Sorter)
	assert.Equal(t, config1.Files, config2.Files)
}

func TestLoadConfigFromFile_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	filename := "test-config.yaml"
	expectedConfig := &Config{
		AllowedImageExtensions: []string{".jpg", ".png"},
	}
	mockLoader := &MockConfigLoader{
		loadConfig: expectedConfig,
		loadErr:    nil,
	}
	mockLoader.On("Load", filename).Return(expectedConfig, nil)

	// Act
	config, err := LoadConfigFromFile(filename, mockLoader)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, expectedConfig, config)
	mockLoader.AssertExpectations(t)
}

func TestLoadConfigFromFile_LoadError(t *testing.T) {
	t.Parallel()

	// Arrange
	filename := "nonexistent-config.yaml"
	expectedErr := errors.New("file not found")
	mockLoader := &MockConfigLoader{
		loadConfig: nil,
		loadErr:    expectedErr,
	}
	mockLoader.On("Load", filename).Return(nil, expectedErr)

	// Act
	config, err := LoadConfigFromFile(filename, mockLoader)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
	assert.Contains(t, err.Error(), filename)
	assert.Nil(t, config)
	mockLoader.AssertExpectations(t)
}

func TestLoadConfigFromFile_EmptyFilename(t *testing.T) {
	t.Parallel()

	// Arrange
	filename := ""
	defaultConfig := DefaultConfig()
	mockLoader := &MockConfigLoader{
		loadConfig: defaultConfig,
		loadErr:    nil,
	}
	mockLoader.On("Load", filename).Return(defaultConfig, nil)

	// Act
	config, err := LoadConfigFromFile(filename, mockLoader)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, config)
	mockLoader.AssertExpectations(t)
}

func TestLoadConfigFromFile_NilLoader(t *testing.T) {
	t.Parallel()

	// Arrange
	filename := "test-config.yaml"

	// Act & Assert - This should panic due to nil pointer dereference
	assert.Panics(t, func() {
		_, _ = LoadConfigFromFile(filename, nil) // Ignore return values since we expect panic
	})
}

func TestNewConfigLoader_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	mockFileReader := &MockFileReader{}
	mockParser := &MockConfigParser{}
	mockLocalizer := &MockLocalizer{}

	// Act
	loader := NewConfigLoader(mockFileReader, mockParser, mockLocalizer)

	// Assert
	assert.NotNil(t, loader)
	assert.IsType(t, &DefaultConfigLoader{}, loader)
}

func TestNewConfigLoader_WithNilDependencies(t *testing.T) {
	t.Parallel()

	// Act
	loader := NewConfigLoader(nil, nil, nil)

	// Assert
	assert.NotNil(t, loader)
	assert.IsType(t, &DefaultConfigLoader{}, loader)
}

func TestNewConfigLoader_ReturnsDefaultConfigLoader(t *testing.T) {
	t.Parallel()

	// Arrange
	mockFileReader := &MockFileReader{}
	mockParser := &MockConfigParser{}
	mockLocalizer := &MockLocalizer{}

	// Act
	loader := NewConfigLoader(mockFileReader, mockParser, mockLocalizer)
	expectedLoader := NewDefaultConfigLoader(mockFileReader, mockParser, mockLocalizer)

	// Assert
	assert.IsType(t, expectedLoader, loader)
}

func TestNewConfigLoader_Integration(t *testing.T) {
	t.Parallel()

	// Arrange
	testData := []byte("test: value")
	expectedConfig := &Config{AllowedImageExtensions: []string{".test"}}

	mockFileReader := &MockFileReader{
		readFileData: testData,
		readFileErr:  nil,
	}
	mockParser := &MockConfigParser{
		parseConfig: expectedConfig,
		parseErr:    nil,
	}
	mockLocalizer := &MockLocalizer{}

	// Act
	loader := NewConfigLoader(mockFileReader, mockParser, mockLocalizer)
	config, err := loader.Load("test.yaml")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, []string{".test"}, config.AllowedImageExtensions)
}

// MockConfigLoader is a mock implementation for testing LoadConfigFromFile.
type MockConfigLoader struct {
	mock.Mock
	loadConfig *Config
	loadErr    error
}

func (m *MockConfigLoader) Load(filename string) (*Config, error) {
	m.Called(filename)
	return m.loadConfig, m.loadErr
}

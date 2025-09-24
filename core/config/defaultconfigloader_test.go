package config

import (
	"errors"
	"testing"
)

func TestDefaultConfigLoader_Load(t *testing.T) {
	t.Parallel()
	// Arrange
	mockFileReader := &MockFileReader{
		readFileData: []byte("test content"),
		readFileErr:  nil,
	}
	mockParser := &MockConfigParser{
		parseConfig: &Config{},
		parseErr:    nil,
	}
	mockLocalizer := &MockLocalizer{}
	loader := NewDefaultConfigLoader(mockFileReader, mockParser, mockLocalizer)

	// Act
	cfg, err := loader.Load("test.yaml")
	// Assert
	if err != nil {
		t.Errorf("Load() unexpected error: %v", err)
	}
	if cfg == nil {
		t.Errorf("Load() returned nil config")
	}

	// Adding test for nil file reader
	loader = NewDefaultConfigLoader(nil, mockParser, mockLocalizer)
	cfg, err = loader.Load("test.yaml")
	if err == nil {
		t.Errorf("Load() expected error for nil file reader")
	}
	if cfg != nil {
		t.Errorf("Load() expected nil config for nil file reader")
	}
}

func TestDefaultConfigLoader_Load_FileNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedErr := errors.New("file not found")
	mockFileReader := &MockFileReader{
		readFileData: nil,
		readFileErr:  expectedErr,
	}
	mockParser := &MockConfigParser{}
	mockLocalizer := &MockLocalizer{}
	loader := NewDefaultConfigLoader(mockFileReader, mockParser, mockLocalizer)

	// Act
	cfg, err := loader.Load("nonexistent.yaml")

	// Assert
	if err == nil {
		t.Errorf("Load() error = nil, want %v", expectedErr)
	}
	if cfg != nil {
		t.Errorf("Load() returned non-nil config")
	}
}

func TestDefaultConfigLoader_Load_ParseError(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedErr := errors.New("parse error")
	mockFileReader := &MockFileReader{
		readFileData: []byte("invalid content"),
		readFileErr:  nil,
	}
	mockParser := &MockConfigParser{
		parseConfig: nil,
		parseErr:    expectedErr,
	}
	mockLocalizer := &MockLocalizer{}
	loader := NewDefaultConfigLoader(mockFileReader, mockParser, mockLocalizer)

	// Act
	cfg, err := loader.Load("invalid.yaml")

	// Assert
	if err == nil {
		t.Errorf("Load() error = nil, want %v", expectedErr)
	}
	if cfg != nil {
		t.Errorf("Load() returned non-nil config")
	}
}

// MockFileReader is a mock implementation for testing.
type MockFileReader struct {
	readFileData []byte
	readFileErr  error
}

func (m *MockFileReader) ReadFile(_ string) ([]byte, error) {
	return m.readFileData, m.readFileErr
}

// MockConfigParser is a mock implementation for testing.
type MockConfigParser struct {
	parseConfig *Config
	parseErr    error
}

func (m *MockConfigParser) Parse(_ []byte) (*Config, error) {
	return m.parseConfig, m.parseErr
}

// MockLocalizer is a mock implementation for testing.
type MockLocalizer struct{}

func (m *MockLocalizer) Translate(_ string, _ ...map[string]interface{}) string {
	return "translated message"
}

func (m *MockLocalizer) GetCurrentLanguage() string {
	return "en"
}

func (m *MockLocalizer) SetLanguage(_ string) error {
	return nil
}

func (m *MockLocalizer) IsInitialized() bool {
	return true
}

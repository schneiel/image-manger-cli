package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultLevelParser_Parse(t *testing.T) {
	t.Parallel()
	// Arrange
	mockLocalizer := &MockLocalizer{}
	parser, err := NewDefaultLevelParser(mockLocalizer)
	require.NoError(t, err)

	tests := []struct {
		name        string
		input       string
		expected    Level
		expectError bool
	}{
		{
			name:        "should parse DEBUG",
			input:       "DEBUG",
			expected:    DEBUG,
			expectError: false,
		},
		{
			name:        "should parse debug (case insensitive)",
			input:       "debug",
			expected:    DEBUG,
			expectError: false,
		},
		{
			name:        "should parse INFO",
			input:       "INFO",
			expected:    INFO,
			expectError: false,
		},
		{
			name:        "should parse WARN",
			input:       "WARN",
			expected:    WARN,
			expectError: false,
		},
		{
			name:        "should parse ERROR",
			input:       "ERROR",
			expected:    ERROR,
			expectError: false,
		},
		{
			name:        "should return error for invalid level",
			input:       "INVALID",
			expected:    DEBUG,
			expectError: true,
		},
		{
			name:        "should return error for empty string",
			input:       "",
			expected:    DEBUG,
			expectError: true,
		},
		{
			name:        "should return error for whitespace string",
			input:       "   ",
			expected:    DEBUG,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			// Act
			result, err := parser.Parse(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("Parse() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Parse() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Parse() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

// MockLocalizer is a simple mock for testing.
type MockLocalizer struct{}

func (m *MockLocalizer) Translate(key string, _ ...map[string]interface{}) string {
	return key
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

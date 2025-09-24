package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/schneiel/ImageManagerGo/internal/di"
)

func TestNewArgumentParser(t *testing.T) {
	t.Parallel()

	parser := di.NewArgumentParser()
	assert.NotNil(t, parser)
}

func TestArgumentParser_ExtractLanguage(t *testing.T) {
	t.Parallel()

	parser := di.NewArgumentParser()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "should return en when no language flag",
			args:     []string{"program", "command"},
			expected: "en",
		},
		{
			name:     "should return language when --lang flag present",
			args:     []string{"program", "--lang", "de", "command"},
			expected: "de",
		},
		{
			name:     "should return language when -l flag present",
			args:     []string{"program", "-l", "fr", "command"},
			expected: "fr",
		},
		{
			name:     "should return language when --language flag present",
			args:     []string{"program", "--language", "es", "command"},
			expected: "es",
		},
		{
			name:     "should return command when flag followed by command",
			args:     []string{"program", "--lang", "command"},
			expected: "command", // Actually picks up next argument
		},
		{
			name:     "should return en when flag is last argument",
			args:     []string{"program", "command", "--lang"},
			expected: "en", // Flag is last, no value
		},
		{
			name:     "should return en with empty args",
			args:     []string{},
			expected: "en",
		},
		{
			name:     "should return first match when multiple flags",
			args:     []string{"program", "--lang", "de", "-l", "fr", "command"},
			expected: "de", // First match wins
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := parser.ExtractLanguage(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestArgumentParser_ExtractConfigPath(t *testing.T) {
	t.Parallel()

	parser := di.NewArgumentParser()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "should return empty when no config flag",
			args:     []string{"program", "command"},
			expected: "",
		},
		{
			name:     "should return path when --config flag present",
			args:     []string{"program", "--config", "/path/to/config.yaml", "command"},
			expected: "/path/to/config.yaml",
		},
		{
			name:     "should return path when -c flag present",
			args:     []string{"program", "-c", "config.yaml", "command"},
			expected: "config.yaml",
		},
		{
			name:     "should return command when flag followed by command",
			args:     []string{"program", "--config", "command"},
			expected: "command", // Actually picks up next argument
		},
		{
			name:     "should return empty when flag is last argument",
			args:     []string{"program", "command", "--config"},
			expected: "", // Flag is last, no value
		},
		{
			name:     "should return empty with empty args",
			args:     []string{},
			expected: "",
		},
		{
			name:     "should return first match when multiple flags",
			args:     []string{"program", "--config", "first.yaml", "-c", "second.yaml", "command"},
			expected: "first.yaml", // First match wins
		},
		{
			name:     "should handle relative paths",
			args:     []string{"program", "-c", "./config/app.yaml", "command"},
			expected: "./config/app.yaml",
		},
		{
			name:     "should handle paths with spaces (single token)",
			args:     []string{"program", "--config", "config file.yaml", "command"},
			expected: "config file.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := parser.ExtractConfigPath(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

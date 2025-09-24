package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/internal/config"
)

func TestDefaultSortConfig(t *testing.T) {
	t.Parallel()
	cfg := config.DefaultSortConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "dryRun", cfg.ActionStrategy)
	assert.Empty(t, cfg.Source)      // Should be empty by default
	assert.Empty(t, cfg.Destination) // Should be empty by default
}

func TestSortConfig_Fields(t *testing.T) {
	t.Parallel()
	cfg := &config.SortConfig{
		Source:         "/path/to/source",
		Destination:    "/path/to/destination",
		ActionStrategy: "copy",
	}

	assert.Equal(t, "/path/to/source", cfg.Source)
	assert.Equal(t, "/path/to/destination", cfg.Destination)
	assert.Equal(t, "copy", cfg.ActionStrategy)
}

func TestSortConfig_ZeroValues(t *testing.T) {
	t.Parallel()
	cfg := &config.SortConfig{}

	assert.Empty(t, cfg.Source)
	assert.Empty(t, cfg.Destination)
	assert.Empty(t, cfg.ActionStrategy)
}

func TestSortConfig_Validate_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.SortConfig
	}{
		{
			name: "valid config with copy strategy",
			config: &config.SortConfig{
				Source:         "/path/to/source",
				Destination:    "/path/to/destination",
				ActionStrategy: "copy",
			},
		},
		{
			name: "valid config with move strategy",
			config: &config.SortConfig{
				Source:         "/home/user/photos",
				Destination:    "/home/user/sorted",
				ActionStrategy: "move",
			},
		},
		{
			name: "valid config with dryRun strategy",
			config: &config.SortConfig{
				Source:         "/tmp/images",
				Destination:    "/tmp/sorted",
				ActionStrategy: "dryRun",
			},
		},
		{
			name: "valid config with relative paths",
			config: &config.SortConfig{
				Source:         "./source",
				Destination:    "./destination",
				ActionStrategy: "copy",
			},
		},
		{
			name: "valid config with same source and destination",
			config: &config.SortConfig{
				Source:         "/path/to/images",
				Destination:    "/path/to/images",
				ActionStrategy: "dryRun",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.Validate()

			assert.NoError(t, err)
		})
	}
}

func TestSortConfig_Validate_MissingSource(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.SortConfig
	}{
		{
			name: "empty source",
			config: &config.SortConfig{
				Source:         "",
				Destination:    "/path/to/destination",
				ActionStrategy: "copy",
			},
		},
		{
			name: "nil source (zero value)",
			config: &config.SortConfig{
				Destination:    "/path/to/destination",
				ActionStrategy: "copy",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.Validate()

			require.Error(t, err)
			assert.Contains(t, err.Error(), "source directory is required")
		})
	}
}

func TestSortConfig_Validate_MissingDestination(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.SortConfig
	}{
		{
			name: "empty destination",
			config: &config.SortConfig{
				Source:         "/path/to/source",
				Destination:    "",
				ActionStrategy: "copy",
			},
		},
		{
			name: "nil destination (zero value)",
			config: &config.SortConfig{
				Source:         "/path/to/source",
				ActionStrategy: "copy",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.Validate()

			require.Error(t, err)
			assert.Contains(t, err.Error(), "destination directory is required")
		})
	}
}

func TestSortConfig_Validate_MissingActionStrategy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.SortConfig
	}{
		{
			name: "empty action strategy",
			config: &config.SortConfig{
				Source:         "/path/to/source",
				Destination:    "/path/to/destination",
				ActionStrategy: "",
			},
		},
		{
			name: "nil action strategy (zero value)",
			config: &config.SortConfig{
				Source:      "/path/to/source",
				Destination: "/path/to/destination",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.Validate()

			require.Error(t, err)
			assert.Contains(t, err.Error(), "action strategy is required")
		})
	}
}

func TestSortConfig_Validate_MultipleErrors(t *testing.T) {
	t.Parallel()

	// Test that validation returns the first error encountered.
	cfg := &config.SortConfig{
		Source:         "",
		Destination:    "",
		ActionStrategy: "",
	}

	err := cfg.Validate()

	require.Error(t, err)
	// Should return the first error (source directory is required)
	assert.Contains(t, err.Error(), "source directory is required")
}

func TestSortConfig_Validate_EdgeCases(
	t *testing.T,
) { //nolint:funlen // comprehensive edge case testing with multiple scenarios
	t.Parallel()

	testCases := []struct {
		name        string
		config      *config.SortConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "whitespace-only source",
			config: &config.SortConfig{
				Source:         "   ",
				Destination:    "/path/to/destination",
				ActionStrategy: "copy",
			},
			expectError: false, // Whitespace is considered valid
		},
		{
			name: "whitespace-only destination",
			config: &config.SortConfig{
				Source:         "/path/to/source",
				Destination:    "   ",
				ActionStrategy: "copy",
			},
			expectError: false, // Whitespace is considered valid
		},
		{
			name: "whitespace-only action strategy",
			config: &config.SortConfig{
				Source:         "/path/to/source",
				Destination:    "/path/to/destination",
				ActionStrategy: "   ",
			},
			expectError: false, // Whitespace is considered valid
		},
		{
			name: "very long paths",
			config: &config.SortConfig{
				Source:         "/very/long/path/that/might/exceed/some/filesystem/limits/but/should/still/be/valid/for/configuration/purposes/source",
				Destination:    "/very/long/path/that/might/exceed/some/filesystem/limits/but/should/still/be/valid/for/configuration/purposes/destination",
				ActionStrategy: "copy",
			},
			expectError: false,
		},
		{
			name: "paths with special characters",
			config: &config.SortConfig{
				Source:         "/path/with spaces & special chars!@#$%/source",
				Destination:    "/path/with spaces & special chars!@#$%/destination",
				ActionStrategy: "move",
			},
			expectError: false,
		},
		{
			name: "unicode paths",
			config: &config.SortConfig{
				Source:         "/пуѕть/к/источнику",
				Destination:    "/путь/к/назначению",
				ActionStrategy: "dryRun",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.Validate()

			if tc.expectError {
				require.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/internal/config"
)

func TestDefaultDedupConfig(t *testing.T) {
	t.Parallel()
	cfg := config.DefaultDedupConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "dryRun", cfg.ActionStrategy)
	assert.Equal(t, "keepOldest", cfg.KeepStrategy)
	assert.Equal(t, 1, cfg.Threshold)
	assert.Empty(t, cfg.Source) // Should be empty by default
}

func TestDedupConfig_Fields(t *testing.T) {
	t.Parallel()
	cfg := &config.DedupConfig{
		Source:         "/path/to/source",
		ActionStrategy: "moveToTrash",
		KeepStrategy:   "keepShortest",
		Threshold:      10,
	}

	assert.Equal(t, "/path/to/source", cfg.Source)
	assert.Equal(t, "moveToTrash", cfg.ActionStrategy)
	assert.Equal(t, "keepShortest", cfg.KeepStrategy)
	assert.Equal(t, 10, cfg.Threshold)
}

func TestDedupConfig_ZeroValues(t *testing.T) {
	t.Parallel()
	cfg := &config.DedupConfig{}

	assert.Empty(t, cfg.Source)
	assert.Empty(t, cfg.ActionStrategy)
	assert.Empty(t, cfg.KeepStrategy)
	assert.Equal(t, 0, cfg.Threshold)
}

func TestDedupConfig_Validate_Success(
	t *testing.T,
) { //nolint:funlen // table-driven test with multiple validation cases
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.DedupConfig
	}{
		{
			name: "valid config with moveToTrash strategy",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepOldest",
				Workers:        4,
				Threshold:      5,
			},
		},
		{
			name: "valid config with dryRun strategy",
			config: &config.DedupConfig{
				Source:         "/home/user/photos",
				ActionStrategy: "dryRun",
				KeepStrategy:   "keepShortest",
				Workers:        2,
				Threshold:      10,
			},
		},
		{
			name: "valid config with minimum threshold",
			config: &config.DedupConfig{
				Source:         "/tmp/images",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepOldest",
				Workers:        1,
				Threshold:      1,
			},
		},
		{
			name: "valid config with high threshold",
			config: &config.DedupConfig{
				Source:         "/path/to/images",
				ActionStrategy: "dryRun",
				KeepStrategy:   "keepShortest",
				Workers:        8,
				Threshold:      100,
			},
		},
		{
			name: "valid config with relative path",
			config: &config.DedupConfig{
				Source:         "./source",
				ActionStrategy: "dryRun",
				KeepStrategy:   "keepOldest",
				Workers:        2,
				Threshold:      5,
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

func TestDedupConfig_Validate_MissingSource(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.DedupConfig
	}{
		{
			name: "empty source",
			config: &config.DedupConfig{
				Source:         "",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepOldest",
				Workers:        4,
				Threshold:      5,
			},
		},
		{
			name: "nil source (zero value)",
			config: &config.DedupConfig{
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepOldest",
				Workers:        4,
				Threshold:      5,
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

func TestDedupConfig_Validate_MissingActionStrategy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.DedupConfig
	}{
		{
			name: "empty action strategy",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "",
				KeepStrategy:   "keepOldest",
				Workers:        4,
				Threshold:      5,
			},
		},
		{
			name: "nil action strategy (zero value)",
			config: &config.DedupConfig{
				Source:       "/path/to/source",
				KeepStrategy: "keepOldest",
				Workers:      4,
				Threshold:    5,
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

func TestDedupConfig_Validate_MissingKeepStrategy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *config.DedupConfig
	}{
		{
			name: "empty keep strategy",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "",
				Workers:        4,
				Threshold:      5,
			},
		},
		{
			name: "nil keep strategy (zero value)",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "moveToTrash",
				Workers:        4,
				Threshold:      5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.config.Validate()

			require.Error(t, err)
			assert.Contains(t, err.Error(), "keep strategy is required")
		})
	}
}

func TestDedupConfig_Validate_InvalidThreshold(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		threshold int
		errorMsg  string
	}{
		{
			name:      "negative threshold",
			threshold: -1,
			errorMsg:  "threshold must be non-negative",
		},
		{
			name:      "very negative threshold",
			threshold: -100,
			errorMsg:  "threshold must be non-negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			config := &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepOldest",
				Threshold:      tc.threshold,
			}

			err := config.Validate()

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

func TestDedupConfig_Validate_MultipleErrors(t *testing.T) {
	t.Parallel()

	// Test that validation returns the first error encountered
	cfg := &config.DedupConfig{
		Source:         "",
		ActionStrategy: "",
		KeepStrategy:   "",
		Threshold:      0,
	}

	err := cfg.Validate()

	require.Error(t, err)
	// Should return the first error (source directory is required)
	assert.Contains(t, err.Error(), "source directory is required")
}

func TestDedupConfig_Validate_EdgeCases(
	t *testing.T,
) { //nolint:funlen // comprehensive edge case testing with multiple scenarios
	t.Parallel()

	testCases := []struct {
		name        string
		config      *config.DedupConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "whitespace-only source",
			config: &config.DedupConfig{
				Source:         "   ",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepOldest",
				Workers:        4,
				Threshold:      5,
			},
			expectError: false, // Whitespace-only source is valid (not empty string)
		},
		{
			name: "whitespace-only action strategy",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "   ",
				KeepStrategy:   "keepOldest",
				Workers:        4,
				Threshold:      5,
			},
			expectError: true, // Whitespace-only action strategy should fail
		},
		{
			name: "whitespace-only keep strategy",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "   ",
				Workers:        4,
				Threshold:      5,
			},
			expectError: true, // Whitespace-only keep strategy should fail
		},
		{
			name: "very long path",
			config: &config.DedupConfig{
				Source:         "/very/long/path/that/might/exceed/some/filesystem/limits/but/should/still/be/valid/for/configuration/purposes/source",
				ActionStrategy: "dryRun",
				KeepStrategy:   "keepOldest",
				Workers:        6,
				Threshold:      5,
			},
			expectError: false,
		},
		{
			name: "path with special characters",
			config: &config.DedupConfig{
				Source:         "/path/with spaces & special chars!@#$%/source",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepShortest",
				Workers:        3,
				Threshold:      10,
			},
			expectError: false,
		},
		{
			name: "unicode path",
			config: &config.DedupConfig{
				Source:         "/пуѕть/к/источнику",
				ActionStrategy: "dryRun",
				KeepStrategy:   "keepOldest",
				Workers:        4,
				Threshold:      3,
			},
			expectError: false,
		},
		{
			name: "maximum reasonable threshold",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "dryRun",
				KeepStrategy:   "keepOldest",
				Workers:        8,
				Threshold:      1000,
			},
			expectError: false,
		},
		{
			name: "threshold of 1",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "moveToTrash",
				KeepStrategy:   "keepShortest",
				Workers:        2,
				Threshold:      1,
			},
			expectError: false,
		},
		{
			name: "threshold of 0 (exact matches only)",
			config: &config.DedupConfig{
				Source:         "/path/to/source",
				ActionStrategy: "dryRun",
				KeepStrategy:   "keepOldest",
				Workers:        1,
				Threshold:      0,
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

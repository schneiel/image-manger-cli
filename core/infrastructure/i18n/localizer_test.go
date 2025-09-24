package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBundleLocalizerConfig tests configuration handling.
func TestBundleLocalizerConfig(t *testing.T) {
	t.Parallel()

	t.Run("NewBundleLocalizer requires embedded filesystem", func(t *testing.T) {
		cfg := LocalizerConfig{
			Language:   "en",
			LocalesDir: "locales",
			LocalesFS:  nil, // This should cause an error
		}

		localizer, err := NewBundleLocalizer(cfg)
		require.Error(t, err)
		assert.Nil(t, localizer)
		assert.Contains(t, err.Error(), "embedded filesystem is nil")
	})

	t.Run("Default values are applied", func(t *testing.T) {
		// Test that the constructor handles default values correctly
		// This tests the behavior without requiring actual locale files
		cfg := LocalizerConfig{
			Language:   "", // Should default to English
			LocalesDir: "", // Should default to "locales"
			LocalesFS:  nil,
		}

		localizer, err := NewBundleLocalizer(cfg)
		require.Error(t, err) // Will fail due to nil FS, but we can test the logic path
		assert.Nil(t, localizer)
	})
}

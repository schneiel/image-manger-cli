package i18n

import (
	"embed"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

//go:embed testdata
var testLocalesFS embed.FS

// TestNewBundleLocalizer_Success tests successful creation of a BundleLocalizer.
func TestNewBundleLocalizer_Success(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)

	require.NoError(t, err)
	require.NotNil(t, localizer)
	assert.Equal(t, "en", localizer.GetCurrentLanguage())
	assert.True(t, localizer.IsInitialized())
}

// TestNewBundleLocalizer_NilFS tests error handling when LocalesFS is nil.
func TestNewBundleLocalizer_NilFS(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  nil,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)

	require.Error(t, err)
	assert.Nil(t, localizer)
	assert.Contains(t, err.Error(), "embedded filesystem is nil")
}

// TestNewBundleLocalizer_DefaultLanguage tests creation with empty language (should use default).
func TestNewBundleLocalizer_DefaultLanguage(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "", // Empty language should use default
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)

	require.NoError(t, err)
	require.NotNil(t, localizer)
	assert.Equal(t, DefaultLanguage.String(), localizer.GetCurrentLanguage())
}

// TestNewBundleLocalizer_ReuseBundle tests that bundle is reused on subsequent calls.
func TestNewBundleLocalizer_ReuseBundle(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	// First creation
	localizer1, err1 := NewBundleLocalizer(cfg)
	require.NoError(t, err1)
	require.NotNil(t, localizer1)

	// Second creation should reuse bundle
	localizer2, err2 := NewBundleLocalizer(cfg)
	require.NoError(t, err2)
	require.NotNil(t, localizer2)

	// Both should work
	assert.True(t, localizer1.IsInitialized())
	assert.True(t, localizer2.IsInitialized())
}

// TestBundleLocalizer_Translate_Success tests successful translation.
func TestBundleLocalizer_Translate_Success(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)
	require.NoError(t, err)

	// Test basic translation
	result := localizer.Translate("test_message")
	assert.Equal(t, "Hello, World!", result)
}

// TestBundleLocalizer_Translate_WithTemplateData tests translation with template data.
func TestBundleLocalizer_Translate_WithTemplateData(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)
	require.NoError(t, err)

	// Test translation with template data
	templateData := map[string]interface{}{
		"Name": "John",
	}
	result := localizer.Translate("test_template", templateData)
	assert.Equal(t, "Hello, John!", result)
}

// TestBundleLocalizer_Translate_MessageNotFound tests fallback for missing messages.
func TestBundleLocalizer_Translate_MessageNotFound(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)
	require.NoError(t, err)

	// Test translation for non-existent message
	result := localizer.Translate("non_existent_message")
	assert.Contains(t, result, "Translation for 'non_existent_message' not found.")
}

// TestBundleLocalizer_Translate_NotInitialized tests behavior when localizer is not initialized.
func TestBundleLocalizer_Translate_NotInitialized(t *testing.T) {
	// Create uninitialized localizer
	localizer := &BundleLocalizer{}

	result := localizer.Translate("test_message")
	assert.Contains(t, result, "Localizer not initialized. MessageID: test_message")
}

// TestBundleLocalizer_SetLanguage_Success tests successful language change.
func TestBundleLocalizer_SetLanguage_Success(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)
	require.NoError(t, err)

	// Change to German
	err = localizer.SetLanguage("de")
	require.NoError(t, err)
	assert.Equal(t, "de", localizer.GetCurrentLanguage())

	// Test translation in German
	result := localizer.Translate("test_message")
	assert.Equal(t, "Hallo, Welt!", result)
}

// TestBundleLocalizer_SetLanguage_UnsupportedLanguage tests error for unsupported language.
func TestBundleLocalizer_SetLanguage_UnsupportedLanguage(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	localizer, err := NewBundleLocalizer(cfg)
	require.NoError(t, err)

	// Try to set unsupported language
	err = localizer.SetLanguage("fr")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "language 'fr' is not supported")
	// Language should remain unchanged
	assert.Equal(t, "en", localizer.GetCurrentLanguage())
}

// TestBundleLocalizer_IsInitialized tests initialization status.
func TestBundleLocalizer_IsInitialized(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		localizer *BundleLocalizer
		expected  bool
	}{
		{
			name:      "nil localizer",
			localizer: nil,
			expected:  false,
		},
		{
			name:      "empty localizer",
			localizer: &BundleLocalizer{},
			expected:  false,
		},
		{
			name: "localizer with nil internal localizer",
			localizer: &BundleLocalizer{
				currentLang: "en",
				localizer:   nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.localizer.IsInitialized()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLoadLocaleFiles_NilFS tests error handling when filesystem is nil.
func TestLoadLocaleFiles_NilFS(t *testing.T) {
	t.Parallel()

	bundle := i18n.NewBundle(DefaultLanguage)
	err := loadLocaleFiles(bundle, nil, "locales")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "embedded filesystem is nil")
}

// TestLoadLocaleFiles_NilBundle tests error handling when bundle is nil.
func TestLoadLocaleFiles_NilBundle(t *testing.T) {
	t.Parallel()

	err := loadLocaleFiles(nil, &testLocalesFS, "locales")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bundle is nil")
}

// TestLoadLocaleFiles_InvalidDirectory tests error handling for invalid directory.
func TestLoadLocaleFiles_InvalidDirectory(t *testing.T) {
	// Initialize bundle first
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}
	_, err := NewBundleLocalizer(cfg)
	require.NoError(t, err)

	// Now test with invalid directory
	bundle := i18n.NewBundle(DefaultLanguage)
	err = loadLocaleFiles(bundle, &testLocalesFS, "invalid_directory")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading locales directory")
}

// TestDefaultLanguage tests the default language constant.
func TestDefaultLanguage(t *testing.T) {
	t.Parallel()

	assert.Equal(t, language.English, DefaultLanguage)
	assert.Equal(t, "en", DefaultLanguage.String())
}

// TestConcurrentAccess tests thread safety of bundle creation.
func TestConcurrentAccess(t *testing.T) {
	cfg := LocalizerConfig{
		Language:   "en",
		LocalesFS:  &testLocalesFS,
		LocalesDir: "testdata",
	}

	// Create multiple localizers concurrently
	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := NewBundleLocalizer(cfg)
			results <- err
		}()
	}

	// Check all results
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		require.NoError(t, err)
	}
}

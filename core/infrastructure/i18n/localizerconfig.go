package i18n

import "embed"

// LocalizerConfig holds configuration for creating a localizer.
type LocalizerConfig struct {
	// Language is the target language code (e.g., "en", "de")
	Language string

	// LocalesFS is the embedded filesystem containing locale files
	LocalesFS *embed.FS

	// LocalesDir is the directory path within the filesystem containing locales
	LocalesDir string

	// FallbackLanguage is used when the target language is not available
	FallbackLanguage string
}

// Package i18n provides internationalization support for the application.
package i18n

// Localizer defines the interface for internationalization services.
type Localizer interface {
	// Translate translates a message by its ID with optional template data
	Translate(messageID string, templateData ...map[string]interface{}) string

	// GetCurrentLanguage returns the currently configured language
	GetCurrentLanguage() string

	// SetLanguage changes the current language
	SetLanguage(lang string) error

	// IsInitialized returns true if the localizer has been properly initialized
	IsInitialized() bool
}

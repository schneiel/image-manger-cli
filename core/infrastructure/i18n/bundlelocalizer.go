package i18n

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// DefaultLanguage is the fallback language if no other is specified.
var DefaultLanguage = language.English

// BundleLocalizer implements the Localizer interface using go-i18n bundles.
type BundleLocalizer struct {
	currentLang string
	localizer   *i18n.Localizer
	bundle      *i18n.Bundle
}

// NewBundleLocalizer creates a new localizer for the given language.
func NewBundleLocalizer(cfg LocalizerConfig) (Localizer, error) {
	bundle := i18n.NewBundle(DefaultLanguage)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	localesDir := cfg.LocalesDir
	if localesDir == "" {
		localesDir = "locales"
	}
	err := loadLocaleFiles(bundle, cfg.LocalesFS, localesDir)
	if err != nil {
		return nil, err
	}

	targetLang := cfg.Language
	if targetLang == "" {
		targetLang = DefaultLanguage.String()
	}

	localizer := i18n.NewLocalizer(bundle, targetLang)
	return &BundleLocalizer{
		currentLang: targetLang,
		localizer:   localizer,
		bundle:      bundle,
	}, nil
}

func loadLocaleFiles(bundle *i18n.Bundle, fs *embed.FS, dir string) error {
	if fs == nil {
		return errors.New("embedded filesystem is nil")
	}

	if bundle == nil {
		return errors.New("bundle is nil")
	}

	files, err := fs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading locales directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			path := fmt.Sprintf("%s/%s", dir, file.Name())
			if _, err := bundle.LoadMessageFileFS(fs, path); err != nil {
				return fmt.Errorf("failed to load message file %s: %w", path, err)
			}
		}
	}
	return nil
}

// Translate translates a message by its ID with optional template data.
func (l *BundleLocalizer) Translate(messageID string, templateData ...map[string]interface{}) string {
	if !l.IsInitialized() {
		return "Localizer not initialized. MessageID: " + messageID
	}
	var data interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}
	localized, err := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		// Fallback for untranslated messages
		return fmt.Sprintf("Translation for '%s' not found.", messageID)
	}
	return localized
}

// GetCurrentLanguage returns the currently configured language.
func (l *BundleLocalizer) GetCurrentLanguage() string {
	return l.currentLang
}

// SetLanguage changes the current language.
func (l *BundleLocalizer) SetLanguage(lang string) error {
	if l.bundle == nil {
		return errors.New("bundle is nil")
	}

	// Check if the language is supported
	isSupported := false
	for _, tag := range l.bundle.LanguageTags() {
		if tag.String() == lang {
			isSupported = true
			break
		}
	}
	if !isSupported {
		return fmt.Errorf("language '%s' is not supported", lang)
	}
	l.currentLang = lang
	l.localizer = i18n.NewLocalizer(l.bundle, lang)
	return nil
}

// IsInitialized returns true if the localizer has been properly initialized.
func (l *BundleLocalizer) IsInitialized() bool {
	return l != nil && l.localizer != nil
}

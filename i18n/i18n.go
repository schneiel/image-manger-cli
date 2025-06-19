// Package i18n handles internationalization for the application.
package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var bundle *i18n.Bundle
var localizer *i18n.Localizer

// Init initializes the i18n bundle and localizer for a given language.
// It loads all YAML translation files from the specified directory.
func Init(lang, localeDir string) error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	localeFiles, err := os.ReadDir(localeDir)
	if err != nil {
		return fmt.Errorf("could not read locale directory: %w", err)
	}

	for _, file := range localeFiles {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml") {
			path := filepath.Join(localeDir, file.Name())
			bundle.MustLoadMessageFile(path)
		}
	}

	localizer = i18n.NewLocalizer(bundle, lang)
	return nil
}

// T translates a message by its ID.
// It supports template data for dynamic values and falls back to the message ID
// if a translation is not found.
func T(messageID string, templateData ...map[string]interface{}) string {
	if localizer == nil {
		return messageID
	}

	var data interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}

	localized, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	if err != nil {
		// Log the error to help find missing translations during development.
		fmt.Fprintf(os.Stderr, "Translation failed for ID '%s': %v\n", messageID, err)
		return messageID
	}

	return localized
}

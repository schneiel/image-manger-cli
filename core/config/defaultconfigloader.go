package config

import (
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultConfigLoader implements Loader with dependency injection.
type DefaultConfigLoader struct {
	fileReader FileReader
	parser     Parser
	localizer  i18n.Localizer
}

// NewDefaultConfigLoader creates a new DefaultConfigLoader.
func NewDefaultConfigLoader(fileReader FileReader, parser Parser, localizer i18n.Localizer) *DefaultConfigLoader {
	return &DefaultConfigLoader{
		fileReader: fileReader,
		parser:     parser,
		localizer:  localizer,
	}
}

// Load loads configuration from the specified file path.
func (l *DefaultConfigLoader) Load(path string) (*Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}

	if l.fileReader == nil {
		return nil, errors.New(
			l.localizer.Translate("ConfigLoaderError", map[string]interface{}{"Error": "file reader is nil"}),
		)
	}

	data, err := l.fileReader.ReadFile(path)
	if err != nil {
		// Check if this is the default config file and it doesn't exist
		if path == "config.yaml" {
			// Return default config if default config file doesn't exist
			return DefaultConfig(), nil
		}
		// For explicitly specified config files, return error
		return nil, fmt.Errorf("config read error: %w", err)
	}

	cfg, err := l.parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	return cfg, nil
}

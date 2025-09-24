package config

import (
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// LoadConfigFromFile creates a configuration from a file using the provided loader.
// This replaces the over-engineered factory wrapper with a simple function.
func LoadConfigFromFile(filename string, loader Loader) (*Config, error) {
	config, err := loader.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", filename, err)
	}
	return config, nil
}

// NewConfigLoader creates a new Loader with the specified dependencies.
// This replaces the over-engineered factory wrapper with a simple constructor.
func NewConfigLoader(fileReader FileReader, parser Parser, localizer i18n.Localizer) Loader {
	return NewDefaultConfigLoader(fileReader, parser, localizer)
}

// Package config provides configuration management for the ImageManager application.
package config

import (
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// Config represents the application configuration structure.
type Config struct {
	Deduplicator DeduplicatorConfig `yaml:"deduplicator"`
	Sorter       SorterConfig       `yaml:"sorter"`
	Files        FilesConfig        `yaml:"files"`

	AllowedImageExtensions []string   `yaml:"allowedImageExtensions"`
	Logger                 log.Logger `yaml:"-"`
}

// DefaultConfig returns a Config instance with default values.
func DefaultConfig() *Config {
	return &Config{
		Deduplicator:           DefaultDeduplicatorConfig(),
		Sorter:                 DefaultSorterConfig(),
		Files:                  DefaultFilesConfig(),
		AllowedImageExtensions: []string{".jpg", ".jpeg", ".png", ".gif"},
	}
}

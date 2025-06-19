package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application's main configuration, mapping directly
// to the structure of the config.yaml file. It holds settings for all major
// components of the application.
type Config struct {
	// Source is the default input directory for operations like sorting or
	// deduplication. This value can be overridden by the --source command-line flag.
	Source string `yaml:"source,omitempty"`
	// Destination is the default output directory for the sort operation.
	// This value can be overridden by the --destination command-line flag.
	Destination string `yaml:"destination,omitempty"`

	DryRun bool `yaml:"dry_run"`
	// Deduplicator holds the configuration specific to the deduplication process.
	Deduplicator DeduplicatorConfig `yaml:"deduplicator"`
	// Date holds the configuration for the date extraction strategies.
	Date DateConfig `yaml:"date"`

	AllowedImageExtensions []string
}

// DeduplicatorConfig contains settings for how the deduplication process
// should behave.
type DeduplicatorConfig struct {
	// ActionStrategy defines what to do with identified duplicate files.
	// Examples: "dryRun", "moveToTrash".
	ActionStrategy string `yaml:"actionStrategy"`
	// KeepStrategy defines the criteria for which file to keep when a
	// set of duplicates is found. Examples: "keepOldest", "keepShortestPath".
	KeepStrategy string `yaml:"keepStrategy"`
}

// DateConfig specifies the strategies for extracting a timestamp from an image file.
type DateConfig struct {
	// Strategies is a list of methods to try for date extraction, in order of
	// priority. The first strategy that successfully finds a date will be used.
	// Examples: "exif", "creationTime", "modTime".
	Strategies []string `yaml:"strategies"`
}

// DefaultConfig returns a pointer to a Config struct with sensible default values.
// This is used as a fallback when no configuration file is provided or if the
// file is missing certain fields.
func DefaultConfig() *Config {
	return &Config{
		Deduplicator: DeduplicatorConfig{
			KeepStrategy: "keepOldest",
		},
		Date: DateConfig{
			Strategies: []string{"exif", "creationTime", "modTime"},
		},
		AllowedImageExtensions: []string{
			".jpg",
			".jpeg",
			".png",
			".raw",
		},
	}
}

// LoadConfig reads a YAML configuration file from the given path, unmarshals it
// into a Config struct, and returns it. If the path is empty, it returns the
// default configuration.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return cfg, nil
}

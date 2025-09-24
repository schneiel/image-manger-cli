package config

import (
	"errors"

	coreconfig "github.com/schneiel/ImageManagerGo/core/config"
)

// SortConfig holds the configuration for the sorting process.
type SortConfig struct {
	Source         string
	Destination    string
	ActionStrategy string
}

// DefaultSortConfig returns a SortConfig with sensible defaults.
func DefaultSortConfig() *SortConfig {
	return &SortConfig{
		ActionStrategy: coreconfig.ActionStrategyDryRun, // Default for CLI - will be shown in help
		Source:         "",
		Destination:    "",
	}
}

// Validate checks the SortConfig for any inconsistencies or missing values.
func (c *SortConfig) Validate() error {
	if c.Source == "" {
		return errors.New("source directory is required")
	}
	if c.Destination == "" {
		return errors.New("destination directory is required")
	}
	if c.ActionStrategy == "" {
		return errors.New("action strategy is required")
	}
	return nil
}

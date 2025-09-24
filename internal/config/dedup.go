// Package config contains configuration types for CLI commands.
package config

import (
	"errors"
	"runtime"
	"strings"

	coreconfig "github.com/schneiel/ImageManagerGo/core/config"
)

// DedupConfig holds configuration specific to deduplication.
type DedupConfig struct {
	Source         string
	ActionStrategy string
	KeepStrategy   string
	TrashPath      string
	Workers        int
	Threshold      int
}

// DefaultDedupConfig returns a DedupConfig with sensible defaults.
func DefaultDedupConfig() *DedupConfig {
	return &DedupConfig{
		ActionStrategy: coreconfig.ActionStrategyDryRun,
		KeepStrategy:   coreconfig.KeepStrategyOldest,
		TrashPath:      ".trash",
		Workers:        runtime.NumCPU(),
		Threshold:      1,
	}
}

// Validate checks the DedupConfig for any inconsistencies or missing values.
// It applies defaults for unset numeric values before validation.
func (c *DedupConfig) Validate() error {
	if c.Source == "" {
		return errors.New("source directory is required")
	}
	if strings.TrimSpace(c.ActionStrategy) == "" {
		return errors.New("action strategy is required")
	}
	if strings.TrimSpace(c.KeepStrategy) == "" {
		return errors.New("keep strategy is required")
	}

	// Apply defaults for numeric values only
	if c.Workers <= 0 {
		c.Workers = runtime.NumCPU()
	}
	if c.Threshold < 0 {
		return errors.New("threshold must be non-negative")
	}
	return nil
}

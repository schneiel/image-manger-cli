package config

import "runtime"

// DeduplicatorConfig contains settings for how the deduplication process
// should behave.
type DeduplicatorConfig struct {
	Source string `yaml:"source"`
	// ActionStrategy defines what to do with identified duplicate files.
	// Examples: "dryRun", "moveToTrash".
	ActionStrategy string `yaml:"actionStrategy"`
	// KeepStrategy defines the criteria for which file to keep when a
	// set of duplicates is found. Examples: "keepOldest", "keepShortestPath".
	KeepStrategy string `yaml:"keepStrategy"`

	Log string `yaml:"log"`

	TrashPath string `yaml:"trashPath"`
	Workers   int    `yaml:"workers"`
	Threshold int    `yaml:"threshold"`
}

// DefaultDeduplicatorConfig returns a DeduplicatorConfig instance with default values.
func DefaultDeduplicatorConfig() DeduplicatorConfig {
	return DeduplicatorConfig{
		ActionStrategy: "", // Empty string - will be set by CLI flags or default to dryRun for safety
		KeepStrategy:   KeepStrategyOldest,
		Log:            "deduplicator.log",
		TrashPath:      ".trash",
		Workers:        runtime.NumCPU(),
		Threshold:      1,
	}
}

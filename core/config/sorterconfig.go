package config

// SorterConfig defines configuration for the image sorting functionality.
type SorterConfig struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`

	// Examples: "dryRun", "copy".
	ActionStrategy string `yaml:"actionStrategy"`

	Log string `yaml:"log"`

	Date DateConfig `yaml:"github.com/schneiel/ImageManagerGo/core/date"`
}

// DefaultSorterConfig returns a SorterConfig instance with default values.
func DefaultSorterConfig() SorterConfig {
	return SorterConfig{
		ActionStrategy: "", // Empty string - will be set by CLI flags or default to dryRun for safety
		Log:            "sorter.log",
		Date:           DefaultDateConfig(),
	}
}

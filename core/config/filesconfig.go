package config

// FilesConfig defines file paths used by the application.
type FilesConfig struct {
	ApplicationLog string `yaml:"applicationLog"`
	DedupDryRunLog string `yaml:"dedupDryRunLog"`
	SortDryRunLog  string `yaml:"sortDryRunLog"`
}

// DefaultFilesConfig returns a FilesConfig instance with default values.
func DefaultFilesConfig() FilesConfig {
	return FilesConfig{
		ApplicationLog: "application.log",
		DedupDryRunLog: "dedup_dry_run_log.csv",
		SortDryRunLog:  "sort_dry_run_log.csv",
	}
}

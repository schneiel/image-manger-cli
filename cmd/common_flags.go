// Package cmd defines the interfaces and structures for the application's command-line commands.
package cmd

import (
	"flag"
)

// CommonFlags holds the flags that are shared across multiple commands.
type CommonFlags struct {
	ConfigPath string
	DryRun     bool
}

// AddCommonFlags adds the common flags (-config, -dry-run) to the given FlagSet.
// It links the flags to the fields in the provided CommonFlags struct.
func AddCommonFlags(fs *flag.FlagSet, flags *CommonFlags) {
	fs.StringVar(&flags.ConfigPath, "config", "", "Path to an optional config.yaml file.")
	fs.BoolVar(&flags.DryRun, "dry-run", false, "Simulates execution without making changes.")
}

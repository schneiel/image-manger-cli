// Package cmd defines the interfaces and structures for the application's command-line commands.
package cmd

import "flag"

// Command is the interface implemented by all commands (e.g., sort, dedup).
// It abstracts the initialization, execution, and provision of usage information.
type Command interface {
	// Name returns the name of the command.
	Name() string
	// Init initializes the command with the given command-line arguments.
	Init(args []string)
	// Run executes the command's logic.
	Run() error
	// Usage returns the command's flag set for displaying help.
	Usage() *flag.FlagSet
}

// Package imagesorter contains the logic for processing and sorting images.
package action_strategy

// SortActionStrategy defines the interface for actions to be taken on files during sorting.
type SortActionStrategy interface {
	// Setup performs any necessary initialization before execution.
	Setup() error
	// Execute performs the action on the identified file.
	Execute(sourcePath, destinationPath string) error
	// Teardown cleans up resources after execution.
	Teardown()
}

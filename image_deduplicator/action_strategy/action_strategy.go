// Package imagededuplicator contains the logic for finding and handling duplicate images.
package action_strategy

// ActionStrategy defines the interface for actions to be taken on duplicate files.
type ActionStrategy interface {
	// Setup performs any necessary initialization before execution.
	Setup() error
	// Execute performs the action on the identified duplicate files.
	Execute(toKeep string, toRemove []string)
	// Teardown cleans up resources after execution.
	Teardown()
}

// Package sortaction provides action strategies for image sorting operations.
package sortaction

import "github.com/schneiel/ImageManagerGo/core/strategies/shared"

// Strategy defines the interface for executing sort operations in an OOP manner.
// Implementing classes should provide specific logic for different action types.
type Strategy interface {
	// Execute performs the action on the source and destination paths
	Execute(sourcePath, destinationPath string) error
	// GetResources returns the resource manager for this strategy
	GetResources() shared.ActionResource
}

// Package shared provides common action strategy patterns for image operations.
package shared

// ActionResource manages resources for action execution (setup/teardown).
// Strategies that don't need resources can return nil from GetResources().
type ActionResource interface {
	Setup() error
	Teardown() error
}

// CSVResource manages CSV file resources with write capability.
type CSVResource interface {
	ActionResource
	WriteRow(row []string) error
}

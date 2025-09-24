// Package task defines a common interface for executable operations.
package task

// Task defines the interface for executable tasks.
type Task interface {
	Run() error
}

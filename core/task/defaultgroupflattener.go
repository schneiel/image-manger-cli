// Package task provides task implementations for image processing operations.
package task

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// DefaultGroupFlattener provides functionality for flattening groups of strings.
type DefaultGroupFlattener struct {
	logger log.Logger
}

// NewDefaultGroupFlattener creates a new DefaultGroupFlattener with injected dependencies.
func NewDefaultGroupFlattener(logger log.Logger) (*DefaultGroupFlattener, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &DefaultGroupFlattener{logger: logger}, nil
}

// Flatten converts a slice of string slices into a single flat slice of strings.
// Logs operation details for debugging and monitoring.
func (f *DefaultGroupFlattener) Flatten(groups [][]string) []string {
	if len(groups) == 0 {
		f.logger.Debug("No groups to flatten")
		return []string{}
	}

	// Calculate total size for preallocation
	totalSize := 0
	for _, group := range groups {
		totalSize += len(group)
	}

	flat := make([]string, 0, totalSize)
	for i, group := range groups {
		if len(group) == 0 {
			f.logger.Debugf("Skipping empty group at index %d", i)
			continue
		}
		flat = append(flat, group...)
	}

	f.logger.Debugf("Flattened %d groups into %d files", len(groups), len(flat))
	return flat
}

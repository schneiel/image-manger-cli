package services

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// SortConfigApplier implements ConfigApplier for SortConfig.
type SortConfigApplier struct {
	logger log.Logger
}

// NewSortConfigApplier creates a new SortConfigApplier.
func NewSortConfigApplier(logger log.Logger) (*SortConfigApplier, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &SortConfigApplier{logger: logger}, nil
}

// Apply implements ConfigApplier interface.
func (a *SortConfigApplier) Apply(src *cmdconfig.SortConfig, dest *config.Config) {
	if src.Source != "" {
		dest.Sorter.Source = src.Source
		a.logger.Debug("SourceDirSet")
	}

	if src.Destination != "" {
		dest.Sorter.Destination = src.Destination
		a.logger.Debug("DestDirSet")
	}

	if src.ActionStrategy != "" {
		dest.Sorter.ActionStrategy = src.ActionStrategy
		a.logger.Debug("ActionStrategySet")
	}
}

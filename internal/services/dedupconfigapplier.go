// Package services provides configuration application services for CLI commands
package services

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// DedupConfigApplier implements ConfigApplier for DedupConfig.
type DedupConfigApplier struct {
	logger log.Logger
}

// NewDedupConfigApplier creates a new DedupConfigApplier.
func NewDedupConfigApplier(logger log.Logger) (*DedupConfigApplier, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &DedupConfigApplier{logger: logger}, nil
}

// Apply implements ConfigApplier interface.
func (a *DedupConfigApplier) Apply(src *cmdconfig.DedupConfig, dest *config.Config) {
	if src.Source != "" {
		dest.Deduplicator.Source = src.Source
		a.logger.Debug("SourceDirSet")
	}

	if src.ActionStrategy != "" {
		dest.Deduplicator.ActionStrategy = src.ActionStrategy
		a.logger.Debug("ActionStrategySet")
	}

	if src.KeepStrategy != "" {
		dest.Deduplicator.KeepStrategy = src.KeepStrategy
		a.logger.Debug("KeepStrategySet")
	}

	if src.TrashPath != "" {
		dest.Deduplicator.TrashPath = src.TrashPath
		a.logger.Debug("TrashPathSet")
	}

	if src.Workers > 0 {
		dest.Deduplicator.Workers = src.Workers
		a.logger.Debug("WorkersSet")
	}

	if src.Threshold >= 0 {
		dest.Deduplicator.Threshold = src.Threshold
		a.logger.Debug("ThresholdSet")
	}
}

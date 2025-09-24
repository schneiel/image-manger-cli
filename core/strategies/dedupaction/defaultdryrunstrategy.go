package dedupaction

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	shared "github.com/schneiel/ImageManagerGo/core/strategies/shared"
)

// DefaultDryRunStrategy implements Strategy for dry-run mode.
type DefaultDryRunStrategy struct {
	logger    log.Logger
	localizer i18n.Localizer
}

// NewDefaultDryRunStrategy creates a new dry-run strategy.
func NewDefaultDryRunStrategy(config shared.ActionConfig) (*DefaultDryRunStrategy, error) {
	if config.Logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if config.Localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	return &DefaultDryRunStrategy{
		logger:    config.Logger,
		localizer: config.Localizer,
	}, nil
}

// Execute simulates the action without making changes.
func (d *DefaultDryRunStrategy) Execute(original, duplicate *image.Image) error {
	if original == nil || duplicate == nil {
		if d.logger != nil {
			d.logger.Errorf("Cannot execute dry run strategy: original or duplicate image is nil")
		}
		return errors.New("original and duplicate images cannot be nil")
	}

	if d.logger != nil && d.localizer != nil {
		d.logger.Infof(
			d.localizer.Translate(
				"DryRunWouldMoveFile",
				map[string]interface{}{"Source": duplicate.FilePath, "Destination": original.FilePath},
			),
		)
	}
	return nil
}

// GetResources returns nil for dry-run strategy.
func (d *DefaultDryRunStrategy) GetResources() shared.ActionResource {
	return nil
}

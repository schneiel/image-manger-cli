package sortaction

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	shared "github.com/schneiel/ImageManagerGo/core/strategies/shared"
)

// DefaultDryRunStrategy implements the dry-run action strategy for imagesorter.
type DefaultDryRunStrategy struct {
	logger    log.Logger
	localizer i18n.Localizer
}

// NewDefaultDryRunStrategy creates a new DefaultDryRunStrategy.
func NewDefaultDryRunStrategy(logger log.Logger, localizer i18n.Localizer) (*DefaultDryRunStrategy, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	return &DefaultDryRunStrategy{
		logger:    logger,
		localizer: localizer,
	}, nil
}

// Execute simulates the copy action without making changes.
func (s *DefaultDryRunStrategy) Execute(source, destination string) error {
	s.logger.Infof(
		s.localizer.Translate(
			"DryRunWouldMoveFile",
			map[string]interface{}{"Source": source, "Destination": destination},
		),
	)
	return nil
}

// GetResources returns the resource manager for this strategy.
// For dry-run strategy, no resources are needed, so it returns nil.
func (s *DefaultDryRunStrategy) GetResources() shared.ActionResource {
	return nil
}

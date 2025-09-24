package services

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// SortConfigValidator implements ConfigValidator for SortConfig.
type SortConfigValidator struct {
	localizer i18n.Localizer
}

// NewSortConfigValidator creates a new SortConfigValidator.
func NewSortConfigValidator(localizer i18n.Localizer) *SortConfigValidator {
	if localizer == nil {
		panic("localizer cannot be nil")
	}

	return &SortConfigValidator{localizer: localizer}
}

// Validate implements ConfigValidator interface.
func (v *SortConfigValidator) Validate(cfg *cmdconfig.SortConfig) error {
	// Validate required fields
	if cfg.Source == "" || cfg.Destination == "" {
		return errors.New(v.localizer.Translate("SortSrcDestRequired"))
	}

	return nil
}

package services

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// DedupConfigValidator implements ConfigValidator for DedupConfig.
type DedupConfigValidator struct {
	localizer i18n.Localizer
}

// NewDedupConfigValidator creates a new DedupConfigValidator.
func NewDedupConfigValidator(localizer i18n.Localizer) *DedupConfigValidator {
	if localizer == nil {
		panic("localizer cannot be nil")
	}

	return &DedupConfigValidator{localizer: localizer}
}

// Validate implements ConfigValidator interface.
func (v *DedupConfigValidator) Validate(cfg *cmdconfig.DedupConfig) error {
	// Validate required fields
	if cfg.Source == "" {
		return errors.New(v.localizer.Translate("DedupTargetDirMissing"))
	}

	// Apply defaults if values are not set
	if cfg.Workers <= 0 {
		defaultCfg := cmdconfig.DefaultDedupConfig()
		cfg.Workers = defaultCfg.Workers
	}

	if cfg.Threshold < 0 {
		return errors.New("threshold must be non-negative")
	}

	return nil
}

package handlers

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// DedupHandler handles the deduplication command execution with improved architecture.
type DedupHandler struct {
	*BaseHandler
	Executor  TaskExecutor
	Applier   ConfigApplier[*cmdconfig.DedupConfig]
	Validator ConfigValidator[*cmdconfig.DedupConfig]
	Localizer i18n.Localizer
	Config    *cmdconfig.DedupConfig
}

// RunE is the entry point for the command.
func (h *DedupHandler) RunE(_ *cobra.Command, _ []string) error {
	err := h.Validator.Validate(h.Config)
	if err != nil {
		return fmt.Errorf("dedup configuration validation failed: %w", err)
	}

	// Apply configuration to global config using injected applier
	h.Applier.Apply(h.Config, h.BaseHandler.Config)

	// Execute the task using injected executor
	h.Logger.Info(h.Localizer.Translate("DedupProcessStarting"))

	err = h.Executor.Execute("dedup", *h.BaseHandler.Config)
	if err != nil {
		errorMsg := h.Localizer.Translate("DedupProcessFailed", map[string]interface{}{"Error": err})
		h.Logger.Errorf("%s", errorMsg)

		return errors.New(errorMsg)
	}

	h.Logger.Info(h.Localizer.Translate("DedupProcessCompleted"))

	return nil
}

// Execute implements the builders.CommandHandler interface.
func (h *DedupHandler) Execute(config interface{}) error {
	if h.Validator == nil {
		return errors.New("dedup handler validator is nil (dependency injection failure)")
	}
	if dedupConfig, ok := config.(*cmdconfig.DedupConfig); ok {
		h.Config = dedupConfig
		return h.RunE(nil, nil)
	}

	return errors.New("invalid config type for dedup handler")
}

// Validate implements the builders.CommandHandler interface.
func (h *DedupHandler) Validate(config interface{}) error {
	if dedupConfig, ok := config.(*cmdconfig.DedupConfig); ok {
		err := h.Validator.Validate(dedupConfig)
		if err != nil {
			return fmt.Errorf("dedup validation failed: %w", err)
		}

		return nil
	}

	return errors.New("invalid config type for dedup handler")
}

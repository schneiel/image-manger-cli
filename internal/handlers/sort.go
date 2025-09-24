package handlers

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// SortHandler handles the sort command execution with improved architecture.
type SortHandler struct {
	*BaseHandler
	Executor  TaskExecutor
	Applier   ConfigApplier[*cmdconfig.SortConfig]
	Validator ConfigValidator[*cmdconfig.SortConfig]
	Localizer i18n.Localizer
	Config    *cmdconfig.SortConfig
}

// RunE is the entry point for the command.
func (h *SortHandler) RunE(_ *cobra.Command, _ []string) error {
	err := h.Validator.Validate(h.Config)
	if err != nil {
		return fmt.Errorf("sort configuration validation failed: %w", err)
	}

	// Apply configuration to global config using injected applier
	h.Applier.Apply(h.Config, h.BaseHandler.Config)

	// Execute the task using injected executor
	h.Logger.Info(h.Localizer.Translate("SortProcessStarting"))

	err = h.Executor.Execute("sort", *h.BaseHandler.Config)
	if err != nil {
		errorMsg := h.Localizer.Translate("SortProcessFailed")
		h.Logger.Errorf("%s", errorMsg)

		return errors.New(errorMsg)
	}

	h.Logger.Info(h.Localizer.Translate("SortProcessCompleted"))

	return nil
}

// Execute implements the builders.CommandHandler interface.
func (h *SortHandler) Execute(config interface{}) error {
	if h.Validator == nil {
		return errors.New("sort handler validator is nil (dependency injection failure)")
	}
	if sortConfig, ok := config.(*cmdconfig.SortConfig); ok {
		h.Config = sortConfig
	}

	return h.RunE(nil, nil)
}

// Validate implements the builders.CommandHandler interface.
func (h *SortHandler) Validate(config interface{}) error {
	if sortConfig, ok := config.(*cmdconfig.SortConfig); ok {
		err := h.Validator.Validate(sortConfig)
		if err != nil {
			return fmt.Errorf("sort validation failed: %w", err)
		}

		return nil
	}

	return errors.New("invalid config type for sort handler")
}

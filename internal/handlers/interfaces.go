package handlers

import (
	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/config"
)

// TaskExecutor interface for executing tasks.
type TaskExecutor interface {
	Execute(taskName string, cfg config.Config) error
}

// ConfigApplier interface for applying configuration changes.
type ConfigApplier[T any] interface {
	Apply(src T, dest *config.Config)
}

// ConfigValidator interface for configuration validation.
type ConfigValidator[T any] interface {
	Validate(cfg T) error
}

// CommandHandler is an interface for command handlers that can be executed by cobra.
type CommandHandler interface {
	RunE(cmd *cobra.Command, args []string) error
	Execute(config interface{}) error
	Validate(config interface{}) error
}

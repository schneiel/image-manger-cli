// Package builders provides command building utilities for the CLI
package builders

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/internal/handlers"
)

// CommandConfig holds configuration for building commands.
type CommandConfig struct {
	use           string
	shortDesc     string
	longDesc      string
	handler       handlers.CommandHandler
	configBuilder func() any
	flagSetup     func(*cobra.Command, any)
}

// CommandOption defines functional options for command configuration.
type CommandOption func(*CommandConfig)

// WithDescription sets the short and long description for the command.
func WithDescription(short, long string) CommandOption {
	return func(c *CommandConfig) {
		c.shortDesc = short
		c.longDesc = long
	}
}

// WithHandler sets the command handler.
func WithHandler(handler handlers.CommandHandler) CommandOption {
	return func(c *CommandConfig) {
		c.handler = handler
	}
}

// WithConfig sets the configuration builder function.
func WithConfig(configBuilder func() any) CommandOption {
	return func(c *CommandConfig) {
		c.configBuilder = configBuilder
	}
}

// WithFlags sets the flag setup function.
func WithFlags(flagSetup func(*cobra.Command, any)) CommandOption {
	return func(c *CommandConfig) {
		c.flagSetup = flagSetup
	}
}

// NewCommand creates a Cobra command using functional options.
func NewCommand(use string, options ...CommandOption) *cobra.Command {
	config := &CommandConfig{
		use: use,
	}

	for _, option := range options {
		option(config)
	}

	var cmdConfig any
	if config.configBuilder != nil {
		cmdConfig = config.configBuilder()
	}

	cmd := &cobra.Command{
		Use:   config.use,
		Short: config.shortDesc,
		Long:  config.longDesc,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := config.handler.Validate(cmdConfig)
			if err != nil {
				return fmt.Errorf("command validation failed: %w", err)
			}

			err = config.handler.Execute(cmdConfig)
			if err != nil {
				return fmt.Errorf("command execution failed: %w", err)
			}
			return nil
		},
		DisableFlagParsing: false,
	}

	cmd.InitDefaultHelpFlag()

	if config.flagSetup != nil && cmdConfig != nil {
		config.flagSetup(cmd, cmdConfig)
	}

	return cmd
}

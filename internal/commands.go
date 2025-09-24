// Package cli provides command constructor functions for creating CLI components
package cli

import (
	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/internal/builders"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	"github.com/schneiel/ImageManagerGo/internal/handlers"
)

// createCommand is a generic helper function to create a cobra command.
func createCommand(
	use, shortDesc, longDesc string,
	handler handlers.CommandHandler,
	configProvider func() any,
	flagSetup func(*cobra.Command, any),
) *cobra.Command {
	cmd := builders.NewCommand(use,
		builders.WithDescription(shortDesc, longDesc),
		builders.WithHandler(handler),
		builders.WithConfig(configProvider),
		builders.WithFlags(flagSetup),
	)
	// Don't override cmd.RunE - let the builder handle it properly
	return cmd
}

// NewDedupCommand creates a dedup command following best practices.
func NewDedupCommand(
	cfg *config.Config,
	_ i18n.Localizer,
	dedupHandler handlers.CommandHandler,
	dedupFlagSetup handlers.FlagSetup,
) *cobra.Command {
	return createCommand(
		"dedup", "DedupCommandDesc", "DedupCommandLongDesc", dedupHandler,
		// Always return a new config object per invocation
		func() any {
			global := cfg
			cfg := cmdconfig.DefaultDedupConfig()
			cfg.Source = global.Deduplicator.Source
			cfg.ActionStrategy = global.Deduplicator.ActionStrategy
			cfg.KeepStrategy = global.Deduplicator.KeepStrategy
			cfg.TrashPath = global.Deduplicator.TrashPath
			cfg.Workers = global.Deduplicator.Workers
			cfg.Threshold = global.Deduplicator.Threshold

			return cfg
		},
		func(cmd *cobra.Command, cfg any) {
			dedupFlagSetup.SetupFlags(cmd, cfg)
		},
	)
}

// NewSortCommand creates a sort command following best practices.
func NewSortCommand(
	cfg *config.Config,
	_ i18n.Localizer,
	sortHandler handlers.CommandHandler,
	sortFlagSetup handlers.FlagSetup,
) *cobra.Command {
	return createCommand(
		"sort", "SortCommandDesc", "SortCommandLongDesc", sortHandler,
		// Always return a new config object per invocation
		func() any {
			global := cfg
			cfg := cmdconfig.DefaultSortConfig()
			cfg.Source = global.Sorter.Source
			cfg.Destination = global.Sorter.Destination
			cfg.ActionStrategy = global.Sorter.ActionStrategy

			return cfg
		},
		func(cmd *cobra.Command, cfg any) {
			sortFlagSetup.SetupFlags(cmd, cfg)
		},
	)
}

// NewAllCommands creates all commands and returns them as a slice.
func NewAllCommands(
	cfg *config.Config,
	localizer i18n.Localizer,
	dedupHandler handlers.CommandHandler,
	sortHandler handlers.CommandHandler,
	dedupFlagSetup handlers.FlagSetup,
	sortFlagSetup handlers.FlagSetup,
) []*cobra.Command {
	return []*cobra.Command{
		NewDedupCommand(cfg, localizer, dedupHandler, dedupFlagSetup),
		NewSortCommand(cfg, localizer, sortHandler, sortFlagSetup),
	}
}

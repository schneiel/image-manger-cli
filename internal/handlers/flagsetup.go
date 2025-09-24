package handlers

import (
	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// FlagSetup provides functionality for setting up command flags.
type FlagSetup interface {
	SetupFlags(cmd *cobra.Command, cfg interface{})
}

// SortFlagSetup provides functionality for setting up sort command flags.
type SortFlagSetup struct {
	localizer i18n.Localizer
}

// NewSortFlagSetup creates a new SortFlagSetup with injected dependencies.
func NewSortFlagSetup(localizer i18n.Localizer) *SortFlagSetup {
	if localizer == nil {
		panic("localizer cannot be nil")
	}

	return &SortFlagSetup{localizer: localizer}
}

// SetupFlags sets up the flags for the sort command.
func (s *SortFlagSetup) SetupFlags(cmd *cobra.Command, cfg interface{}) {
	sortCfg, ok := cfg.(*cmdconfig.SortConfig)
	if !ok {
		return
	}

	cmd.Flags().StringVarP(&sortCfg.Source, "source", "s", sortCfg.Source, "")
	cmd.Flags().StringVarP(&sortCfg.Destination, "destination", "d", sortCfg.Destination, "")
	cmd.Flags().StringVarP(&sortCfg.ActionStrategy, "actionStrategy", "", sortCfg.ActionStrategy, "")
}

// DedupFlagSetup provides functionality for setting up dedup command flags.
type DedupFlagSetup struct {
	localizer i18n.Localizer
}

// NewDedupFlagSetup creates a new DedupFlagSetup with injected dependencies.
func NewDedupFlagSetup(localizer i18n.Localizer) *DedupFlagSetup {
	if localizer == nil {
		panic("localizer cannot be nil")
	}

	return &DedupFlagSetup{localizer: localizer}
}

// SetupFlags sets up the flags for the dedup command.
func (d *DedupFlagSetup) SetupFlags(cmd *cobra.Command, cfg interface{}) {
	dedupCfg, ok := cfg.(*cmdconfig.DedupConfig)
	if !ok {
		return
	}

	cmd.Flags().StringVarP(&dedupCfg.Source, "source", "s", dedupCfg.Source, "")
	cmd.Flags().StringVarP(&dedupCfg.ActionStrategy, "actionStrategy", "", dedupCfg.ActionStrategy, "")
	cmd.Flags().StringVarP(&dedupCfg.KeepStrategy, "keepStrategy", "", dedupCfg.KeepStrategy, "")
	cmd.Flags().StringVarP(&dedupCfg.TrashPath, "trashPath", "", dedupCfg.TrashPath, "")
	cmd.Flags().IntVarP(&dedupCfg.Workers, "workers", "w", dedupCfg.Workers, "")
	cmd.Flags().IntVarP(&dedupCfg.Threshold, "threshold", "t", dedupCfg.Threshold, "")
}

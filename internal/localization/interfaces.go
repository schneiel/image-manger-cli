package localization

import (
	"github.com/spf13/cobra"
)

// BaseLocalizer defines the interface for command localization.
type BaseLocalizer interface {
	LocalizeCommand(cmd *cobra.Command, shortKey, longKey string, flagUsages map[string]string)
	LocalizeAllCommands(cmd *cobra.Command)
	SetCobraTemplates(cmd *cobra.Command)
	LocalizeBuiltinCommands(cmd *cobra.Command)
}

// Package localization provides internationalization support for CLI commands
package localization

import (
	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

const (
	// Command names.
	sortCommand  = "sort"
	dedupCommand = "dedup"
)

// CommandLocalizer provides base localization functionality.
type CommandLocalizer struct {
	language  string
	localizer i18n.Localizer
}

// NewCommandLocalizer creates a new CommandLocalizer.
func NewCommandLocalizer(language string, localizer i18n.Localizer) *CommandLocalizer {
	return &CommandLocalizer{
		language:  language,
		localizer: localizer,
	}
}

// LocalizeCommand localizes a specific command.
func (cl *CommandLocalizer) LocalizeCommand(
	cmd *cobra.Command,
	shortKey, longKey string,
	flagUsages map[string]string,
) {
	cmd.Short = cl.localizer.Translate(shortKey)
	cmd.Long = cl.localizer.Translate(longKey)
	for flagName, usageKey := range flagUsages {
		if flag := cmd.Flag(flagName); flag != nil {
			flag.Usage = cl.localizer.Translate(usageKey)
		}
	}
}

// LocalizeAllCommands localizes all commands in the command tree.
func (cl *CommandLocalizer) LocalizeAllCommands(cmd *cobra.Command) {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Short != "" {
			subCmd.Short = cl.localizer.Translate(subCmd.Short)
		}
		if subCmd.Long != "" {
			subCmd.Long = cl.localizer.Translate(subCmd.Long)
		}

		cl.localizeCommandFlags(subCmd)
	}

	cl.LocalizeBuiltinCommands(cmd)
}

// localizeCommandFlags localizes flags for specific commands.
func (cl *CommandLocalizer) localizeCommandFlags(cmd *cobra.Command) {
	switch cmd.Use {
	case sortCommand:
		cl.localizeSortFlags(cmd)
	case dedupCommand:
		cl.localizeDedupFlags(cmd)
	}
}

// localizeSortFlags localizes flags for the sort command.
func (cl *CommandLocalizer) localizeSortFlags(cmd *cobra.Command) {
	flagMappings := map[string]string{
		"source":         "SortSourceFlagDesc",
		"destination":    "SortDestinationFlagDesc",
		"actionStrategy": "SortActionStrategyFlagDesc",
		"help":           "help_for_sort",
	}
	cl.applyFlagTranslations(cmd, flagMappings)
}

// localizeDedupFlags localizes flags for the dedup command.
func (cl *CommandLocalizer) localizeDedupFlags(cmd *cobra.Command) {
	flagMappings := map[string]string{
		"source":         "DedupSourceFlagDesc",
		"actionStrategy": "DedupActionStrategyFlagDesc",
		"keepStrategy":   "DedupKeepStrategyFlagDesc",
		"trashPath":      "DedupTrashPathFlagDesc",
		"workers":        "DedupWorkersFlagDesc",
		"threshold":      "DedupThresholdFlagDesc",
		"help":           "help_for_dedup",
	}
	cl.applyFlagTranslations(cmd, flagMappings)
}

// applyFlagTranslations applies translations to flags based on provided mappings.
func (cl *CommandLocalizer) applyFlagTranslations(cmd *cobra.Command, flagMappings map[string]string) {
	for flagName, translationKey := range flagMappings {
		if flag := cmd.Flag(flagName); flag != nil {
			flag.Usage = cl.localizer.Translate(translationKey)
		}
	}
}

// SetCobraTemplates sets default Cobra templates (no-op for base implementation).
func (cl *CommandLocalizer) SetCobraTemplates(_ *cobra.Command) {
	// Default implementation does nothing
}

// LocalizeBuiltinCommands localizes built-in commands (no-op for base implementation).
func (cl *CommandLocalizer) LocalizeBuiltinCommands(_ *cobra.Command) {
	// Default implementation does nothing
}

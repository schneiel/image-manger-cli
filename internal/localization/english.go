package localization

import (
	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// EnglishLocalizer handles English-specific localization.
type EnglishLocalizer struct {
	*CommandLocalizer
}

// NewEnglishLocalizer creates a new EnglishLocalizer instance.
func NewEnglishLocalizer(localizer i18n.Localizer) *EnglishLocalizer {
	return &EnglishLocalizer{
		CommandLocalizer: NewCommandLocalizer("en", localizer),
	}
}

// SetCobraTemplates applies English Cobra templates (uses defaults).
func (el *EnglishLocalizer) SetCobraTemplates(_ *cobra.Command) {
	// English uses default Cobra templates - no special formatting needed
}

// LocalizeBuiltinCommands localizes built-in Cobra commands for English.
func (el *EnglishLocalizer) LocalizeBuiltinCommands(cmd *cobra.Command) {
	for _, command := range cmd.Commands() {
		switch command.Use {
		case "help":
			command.Short = el.localizer.Translate("help_command_short")
		case "completion":
			command.Short = el.localizer.Translate("completion_command_short")
		case "sort":
			// Override the automatic "help for sort" text
			command.Short = el.localizer.Translate("help_for_sort")
		case "dedup":
			// Override the automatic "help for dedup" text
			command.Short = el.localizer.Translate("help_for_dedup")
		}
	}
}

// LocalizeAllCommands applies English localization to all registered commands.
func (el *EnglishLocalizer) LocalizeAllCommands(cmd *cobra.Command) {
	// Call parent implementation first
	el.CommandLocalizer.LocalizeAllCommands(cmd)

	// Apply English-specific localization if needed
	// Currently, English uses the default behavior
}

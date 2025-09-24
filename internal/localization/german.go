package localization

import (
	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// GermanLocalizer handles German-specific localization.
type GermanLocalizer struct {
	*CommandLocalizer
}

// NewGermanLocalizer creates a new GermanLocalizer instance.
func NewGermanLocalizer(localizer i18n.Localizer) *GermanLocalizer {
	return &GermanLocalizer{
		CommandLocalizer: NewCommandLocalizer("de", localizer),
	}
}

// SetCobraTemplates applies German Cobra templates.
func (gl *GermanLocalizer) SetCobraTemplates(cmd *cobra.Command) {
	cmd.Root().SetUsageTemplate(gl.getGermanUsageTemplate())
	cmd.Root().SetHelpTemplate(gl.getGermanHelpTemplate())
}

// LocalizeBuiltinCommands localizes built-in Cobra commands for German.
func (gl *GermanLocalizer) LocalizeBuiltinCommands(cmd *cobra.Command) {
	for _, command := range cmd.Commands() {
		switch command.Use {
		case "help":
			command.Short = gl.localizer.Translate("help_command_short")
		case "completion":
			command.Short = gl.localizer.Translate("completion_command_short")
		case "sort":
			command.Short = gl.localizer.Translate("SortCommandDesc")
		case "dedup":
			command.Short = gl.localizer.Translate("DedupCommandDesc")
		}

		// Special handling for help command if it has a different pattern
		if command.Name() == "help" {
			command.Short = gl.localizer.Translate("help_command_short")
		}
	}
}

// LocalizeAllCommands applies German localization to all registered commands.
func (gl *GermanLocalizer) LocalizeAllCommands(cmd *cobra.Command) {
	gl.CommandLocalizer.LocalizeAllCommands(cmd)
}

// getGermanUsageTemplate returns the German usage template.
func (gl *GermanLocalizer) getGermanUsageTemplate() string {
	return `Verwendung:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [Befehl]{{end}}{{if gt (len .Aliases) 0}}

Aliase:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Beispiele:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Verfuegbare Befehle:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Globale Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Weitere Hilfe-Befehle:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Verwenden Sie "{{.CommandPath}} [Befehl] --help" fuer mehr Informationen ueber einen Befehl.{{end}}
`
}

// getGermanHelpTemplate returns the German help template.
func (gl *GermanLocalizer) getGermanHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

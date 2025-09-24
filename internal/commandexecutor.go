// Package cli provides the main CLI command structure and execution logic
package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/internal/handlers"
	"github.com/schneiel/ImageManagerGo/internal/localization"
)

// DefaultCommandExecutor implements CommandExecutor with proper dependency injection.
type DefaultCommandExecutor struct {
	args             []string
	localizer        i18n.Localizer
	fileReader       config.FileReader
	parser           config.Parser
	config           *config.Config
	sortHandler      handlers.CommandHandler
	dedupHandler     handlers.CommandHandler
	sortFlagSetup    handlers.FlagSetup
	dedupFlagSetup   handlers.FlagSetup
	commandLocalizer localization.BaseLocalizer
}

// Execute runs the CLI application with the provided localizer.
func (e *DefaultCommandExecutor) Execute() error {
	lang := e.parseLanguageFlag()

	// Update the localizer to use the parsed language
	err := e.localizer.SetLanguage(lang)
	if err != nil {
		return fmt.Errorf("failed to set language: %w", err)
	}

	rootCmd := e.createRootCommand()

	// Use injected command factory
	subcommands := NewAllCommands(
		e.config,
		e.localizer,
		e.dedupHandler,
		e.sortHandler,
		e.dedupFlagSetup,
		e.sortFlagSetup,
	)

	e.setupLocalization(rootCmd, lang, e.localizer)
	e.setupRootCommandText(rootCmd, e.localizer)

	rootCmd.AddCommand(subcommands...)
	e.commandLocalizer.LocalizeAllCommands(rootCmd)

	// Force Cobra to initialize built-in commands by calling InitDefaultHelpCmd
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultCompletionCmd()

	// Localize built-in commands after they're initialized
	e.commandLocalizer.LocalizeBuiltinCommands(rootCmd)

	// Localize help flag after commands are initialized
	if flag := rootCmd.Flag("help"); flag != nil {
		flag.Usage = e.localizer.Translate("help_flag_desc")
	}

	err = rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

// parseLanguageFlag extracts the language flag from command line arguments.
func (e *DefaultCommandExecutor) parseLanguageFlag() string {
	if len(e.args) <= 1 {
		return "en" // default language
	}
	args := e.args[1:]
	language := "en" // default language
	for index, arg := range args {
		if arg == "-l" || arg == "--language" {
			if index+1 >= len(args) {
				continue
			}

			language = args[index+1]
		}
		if strings.HasPrefix(arg, "--language=") {
			language = strings.TrimPrefix(arg, "--language=")
		}
	}

	return language
}

// createRootCommand creates the root cobra command.
func (e *DefaultCommandExecutor) createRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "image-manager",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	cmd.PersistentFlags().StringP("language", "l", "en", "Language for localization (e.g., 'en', 'de')")
	cmd.PersistentFlags().StringP("github.com/schneiel/ImageManagerGo/core/config", "c", "config.yaml", "Custom config file path")

	return cmd
}

// setupLocalization configures localization for the root command.
func (e *DefaultCommandExecutor) setupLocalization(rootCmd *cobra.Command, _ string, _ i18n.Localizer) {
	e.commandLocalizer.SetCobraTemplates(rootCmd)
}

// setupRootCommandText sets up the text for the root command.
func (e *DefaultCommandExecutor) setupRootCommandText(rootCmd *cobra.Command, localizer i18n.Localizer) {
	rootCmd.Short = localizer.Translate("root_command_short_description")
	rootCmd.Long = localizer.Translate("root_command_long_description")

	if flag := rootCmd.Flag("github.com/schneiel/ImageManagerGo/core/config"); flag != nil {
		flag.Usage = localizer.Translate("CustomConfigFlagDesc")
	}
	if flag := rootCmd.Flag("language"); flag != nil {
		flag.Usage = localizer.Translate("LanguageFlagDesc")
	}
	if flag := rootCmd.Flag("help"); flag != nil {
		flag.Usage = localizer.Translate("help_flag_desc")
	}
}

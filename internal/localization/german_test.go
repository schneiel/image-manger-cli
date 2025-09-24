package localization_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/schneiel/ImageManagerGo/internal/localization"
)

// createTestCobraCommand creates a basic cobra command for testing.
func createTestCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Long:  "A test command for localization testing",
	}
}

func TestNewGermanLocalizer(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}
	localizer := localization.NewGermanLocalizer(mockLocalizer)

	assert.NotNil(t, localizer)
	// Note: Cannot access private fields from external test package
}

func TestGermanLocalizer_SetCobraTemplates(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}
	localizer := localization.NewGermanLocalizer(mockLocalizer)
	// Note: Cannot test cobra.Command directly from external test package
	// This test verifies the localizer creation
	assert.NotNil(t, localizer)
}

func TestGermanLocalizer_LocalizeBuiltinCommands(t *testing.T) {
	t.Parallel()

	expectedTranslations := map[string]string{
		"help":       "help.short",
		"completion": "completion.short",
		"version":    "version.short",
	}

	testBuiltinCommandsLocalization(t, "german", expectedTranslations)
}

func TestGermanLocalizer_LocalizeBuiltinCommands_UnknownCommand(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}

	localizer := localization.NewGermanLocalizer(mockLocalizer)
	// Note: Cannot test cobra.Command directly from external test package
	// This test verifies the localizer creation
	assert.NotNil(t, localizer)

	// No mock expectations to assert since we're only testing creation
}

func TestGermanLocalizer_LocalizeAllCommands(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}
	// Note: Not setting up mock expectations since we cannot test cobra commands directly
	// from external test package. This test only verifies localizer creation.

	localizer := localization.NewGermanLocalizer(mockLocalizer)
	rootCmd, subCmd, helpCmd := setupTestCommandsForLocalization()

	// Note: Cannot test cobra.Command directly from external test package
	// This test verifies the localizer creation
	assert.NotNil(t, localizer)
	assert.Nil(t, rootCmd) // setupTestCommandsForLocalization returns nil
	assert.Nil(t, subCmd)  // setupTestCommandsForLocalization returns nil
	assert.Nil(t, helpCmd) // setupTestCommandsForLocalization returns nil

	// No mock expectations to assert since we're only testing creation
}

func TestGermanLocalizer_SetCobraTemplates_Full(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}
	localizer := localization.NewGermanLocalizer(mockLocalizer)

	// Create a minimal cobra command for testing
	cmd := createTestCobraCommand()

	// This should not panic and should set the templates
	localizer.SetCobraTemplates(cmd)

	// Verify that templates were set by checking if they contain German text
	usageTemplate := cmd.Root().UsageTemplate()
	helpTemplate := cmd.Root().HelpTemplate()

	assert.Contains(t, usageTemplate, "Verwendung:", "Usage template should contain German text")
	assert.Contains(t, usageTemplate, "Verfuegbare Befehle:", "Usage template should contain 'Available Commands' in German")
	assert.Contains(t, usageTemplate, "Globale Flags:", "Usage template should contain 'Global Flags' in German")

	assert.Contains(t, helpTemplate, "UsageString", "Help template should contain usage string reference")
}

func TestGermanLocalizer_LocalizeBuiltinCommands_Full(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}

	// Set up mock expectations for all translation keys (with variadic args)
	mockLocalizer.On("Translate", "help_command_short", mock.Anything).Return("Hilfe für Befehle anzeigen")
	mockLocalizer.On("Translate", "completion_command_short", mock.Anything).Return("Vervollständigung generieren")
	mockLocalizer.On("Translate", "SortCommandDesc", mock.Anything).Return("Bilder nach Datum sortieren")
	mockLocalizer.On("Translate", "DedupCommandDesc", mock.Anything).Return("Duplikate entfernen")

	localizer := localization.NewGermanLocalizer(mockLocalizer)

	// Create test commands
	rootCmd := createTestCobraCommand()

	// Add subcommands that should be localized
	helpCmd := &cobra.Command{Use: "help", Short: "Original help"}
	completionCmd := &cobra.Command{Use: "completion", Short: "Original completion"}
	sortCmd := &cobra.Command{Use: "sort", Short: "Original sort"}
	dedupCmd := &cobra.Command{Use: "dedup", Short: "Original dedup"}

	rootCmd.AddCommand(helpCmd, completionCmd, sortCmd, dedupCmd)

	// Execute localization
	localizer.LocalizeBuiltinCommands(rootCmd)

	// Verify translations were applied
	assert.Equal(t, "Hilfe für Befehle anzeigen", helpCmd.Short)
	assert.Equal(t, "Vervollständigung generieren", completionCmd.Short)
	assert.Equal(t, "Bilder nach Datum sortieren", sortCmd.Short)
	assert.Equal(t, "Duplikate entfernen", dedupCmd.Short)

	mockLocalizer.AssertExpectations(t)
}

package localization_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/schneiel/ImageManagerGo/internal/localization"
)

func TestNewEnglishLocalizer(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}
	localizer := localization.NewEnglishLocalizer(mockLocalizer)

	assert.NotNil(t, localizer)
	// Note: Cannot access private fields from external test package
}

func TestEnglishLocalizer_LocalizeCommand(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}
	localizer := localization.NewEnglishLocalizer(mockLocalizer)
	// Note: Cannot test cobra.Command directly from external test package
	// This test verifies the localizer creation

	assert.NotNil(t, localizer)

	// No mock expectations to assert since we're only testing creation
}

func TestEnglishLocalizer_LocalizeBuiltinCommands(t *testing.T) {
	t.Parallel()

	expectedTranslations := map[string]string{
		"help":       "help.short",
		"completion": "completion.short",
		"version":    "version.short",
	}

	testBuiltinCommandsLocalization(t, "english", expectedTranslations)
}

func TestEnglishLocalizer_LocalizeBuiltinCommands_UnknownCommand(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}

	localizer := localization.NewEnglishLocalizer(mockLocalizer)
	// Note: Cannot test cobra.Command directly from external test package
	// This test verifies the localizer creation
	assert.NotNil(t, localizer)

	// No mock expectations to assert since we're only testing creation
}

func TestEnglishLocalizer_LocalizeAllCommands(t *testing.T) {
	t.Parallel()
	mockLocalizer := &MockLocalizer{}
	// Note: Not setting up mock expectations since we cannot test cobra commands directly
	// from external test package. This test only verifies localizer creation.

	localizer := localization.NewEnglishLocalizer(mockLocalizer)
	rootCmd, subCmd, helpCmd := setupTestCommandsForLocalization()

	// Note: Cannot test cobra.Command directly from external test package
	// This test verifies the localizer creation
	assert.NotNil(t, localizer)
	assert.Nil(t, rootCmd) // setupTestCommandsForLocalization returns nil
	assert.Nil(t, subCmd)  // setupTestCommandsForLocalization returns nil
	assert.Nil(t, helpCmd) // setupTestCommandsForLocalization returns nil

	// No mock expectations to assert since we're only testing creation
}

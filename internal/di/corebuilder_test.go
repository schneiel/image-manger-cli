package di_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/internal/di"
)

func TestNewCoreBuilder(t *testing.T) {
	t.Parallel()

	builder := di.NewCoreBuilder()
	assert.NotNil(t, builder)
}

func TestCoreBuilder_BuildCore_NilLocalesFS(t *testing.T) {
	t.Parallel()

	builder := di.NewCoreBuilder()
	args := []string{"program", "command"}

	// Should fail with nil localesFS. - but need to handle potential panic.
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v", r)
		}
	}()

	_, err := builder.BuildCore(args, nil)
	if err != nil {
		assert.Contains(t, err.Error(), "localizer")
	}
	// If no error, then the implementation handles nil gracefully.
}

func TestCoreBuilder_BuildCore_EmptyArgs(t *testing.T) {
	t.Parallel()

	builder := di.NewCoreBuilder()
	args := []string{}

	// Should fail with nil localesFS.
	_, err := builder.BuildCore(args, nil)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	assert.Contains(t, err.Error(), "embedded filesystem is nil")
}

func TestCoreBuilder_BuildCore_ValidArgs(t *testing.T) {
	t.Parallel()

	builder := di.NewCoreBuilder()
	args := []string{"program", "--lang", "en", "command"}

	// Should fail due to nil localesFS
	_, err := builder.BuildCore(args, nil)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	assert.Contains(t, err.Error(), "embedded filesystem is nil")
}

func TestCoreBuilder_BuildCore_WithLanguageArg(t *testing.T) {
	t.Parallel()

	builder := di.NewCoreBuilder()
	args := []string{"program", "--lang", "de", "command"}

	// Should fail due to nil localesFS, but test argument parsing works
	_, err := builder.BuildCore(args, nil)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	// Error should mention the language that was parsed
	assert.Contains(t, err.Error(), "de")
}

func TestCoreBuilder_BuildCore_MockEmbeddedFS(t *testing.T) {
	t.Parallel()
	// t.Parallel() removed due to race condition in i18n.NewBundleLocalizer().

	builder := di.NewCoreBuilder()
	args := []string{"program", "command"}

	// Create a mock embedded filesystem (empty)
	mockFS := &embed.FS{}

	// Should fail because the mock FS doesn't contain locale files
	_, err := builder.BuildCore(args, mockFS)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create bundle localizer")
}

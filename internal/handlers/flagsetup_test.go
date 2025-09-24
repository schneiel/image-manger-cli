package handlers

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

func TestNewSortFlagSetup(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}

	flagSetup := NewSortFlagSetup(mockLocalizer)

	if flagSetup == nil {
		t.Fatal("Expected flag setup to be created, got nil")
	}

	// Note: Cannot test unexported field localizer from external test package
}

func TestNewSortFlagSetup_NilLocalizer(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, but none occurred")
		}
	}()

	NewSortFlagSetup(nil)
}

func TestSortFlagSetup_SetupFlags_Success(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	flagSetup := NewSortFlagSetup(mockLocalizer)

	cmd := &cobra.Command{}
	sortConfig := &cmdconfig.SortConfig{
		Source:         "/test/source",
		Destination:    "/test/destination",
		ActionStrategy: "copy",
	}

	flagSetup.SetupFlags(cmd, sortConfig)

	// Verify flags were set up
	if cmd.Flags().Lookup("source") == nil {
		t.Error("Expected source flag to be created")
	}

	if cmd.Flags().Lookup("destination") == nil {
		t.Error("Expected destination flag to be created")
	}

	if cmd.Flags().Lookup("actionStrategy") == nil {
		t.Error("Expected actionStrategy flag to be created")
	}
}

func TestSortFlagSetup_SetupFlags_InvalidConfigType(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	flagSetup := NewSortFlagSetup(mockLocalizer)

	cmd := &cobra.Command{}
	invalidConfig := "not a sort config"

	// Should not panic or error
	flagSetup.SetupFlags(cmd, invalidConfig)

	// Verify no flags were created
	if cmd.Flags().Lookup("source") != nil {
		t.Error("Expected no source flag to be created for invalid config")
	}
}

func TestNewDedupFlagSetup(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}

	flagSetup := NewDedupFlagSetup(mockLocalizer)

	if flagSetup == nil {
		t.Fatal("Expected flag setup to be created, got nil")
	}

	// Note: Cannot test unexported field localizer from external test package
}

func TestNewDedupFlagSetup_NilLocalizer(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, but none occurred")
		}
	}()

	NewDedupFlagSetup(nil)
}

func TestDedupFlagSetup_SetupFlags_Success(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	flagSetup := NewDedupFlagSetup(mockLocalizer)

	cmd := &cobra.Command{}
	dedupConfig := &cmdconfig.DedupConfig{
		Source:         "/test/source",
		ActionStrategy: "dryRun",
		KeepStrategy:   "keepOldest",
		TrashPath:      ".trash",
		Workers:        4,
		Threshold:      1,
	}

	flagSetup.SetupFlags(cmd, dedupConfig)

	// Verify flags were set up
	if cmd.Flags().Lookup("source") == nil {
		t.Error("Expected source flag to be created")
	}

	if cmd.Flags().Lookup("actionStrategy") == nil {
		t.Error("Expected actionStrategy flag to be created")
	}

	if cmd.Flags().Lookup("keepStrategy") == nil {
		t.Error("Expected keepStrategy flag to be created")
	}

	if cmd.Flags().Lookup("trashPath") == nil {
		t.Error("Expected trashPath flag to be created")
	}

	if cmd.Flags().Lookup("workers") == nil {
		t.Error("Expected workers flag to be created")
	}

	if cmd.Flags().Lookup("threshold") == nil {
		t.Error("Expected threshold flag to be created")
	}
}

func TestDedupFlagSetup_SetupFlags_InvalidConfigType(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	flagSetup := NewDedupFlagSetup(mockLocalizer)

	cmd := &cobra.Command{}
	invalidConfig := "not a dedup config"

	// Should not panic or error
	flagSetup.SetupFlags(cmd, invalidConfig)

	// Verify no flags were created
	if cmd.Flags().Lookup("source") != nil {
		t.Error("Expected no source flag to be created for invalid config")
	}
}

func TestFlagSetup_InterfaceCompliance(t *testing.T) {
	t.Parallel()
	// Test that both flag setups implement the FlagSetup interface
	var _ FlagSetup = (*SortFlagSetup)(nil)
	var _ FlagSetup = (*DedupFlagSetup)(nil)
}

func TestSortFlagSetup_SetupFlags_EmptyConfig(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	flagSetup := NewSortFlagSetup(mockLocalizer)

	cmd := &cobra.Command{}
	sortConfig := &cmdconfig.SortConfig{}

	flagSetup.SetupFlags(cmd, sortConfig)

	// Verify flags were set up even with empty config
	if cmd.Flags().Lookup("source") == nil {
		t.Error("Expected source flag to be created")
	}

	if cmd.Flags().Lookup("destination") == nil {
		t.Error("Expected destination flag to be created")
	}

	if cmd.Flags().Lookup("actionStrategy") == nil {
		t.Error("Expected actionStrategy flag to be created")
	}
}

func TestDedupFlagSetup_SetupFlags_EmptyConfig(t *testing.T) {
	t.Parallel()
	mockLocalizer := &testutils.MockLocalizer{}
	flagSetup := NewDedupFlagSetup(mockLocalizer)

	cmd := &cobra.Command{}
	dedupConfig := &cmdconfig.DedupConfig{}

	flagSetup.SetupFlags(cmd, dedupConfig)

	// Verify flags were set up even with empty config
	if cmd.Flags().Lookup("source") == nil {
		t.Error("Expected source flag to be created")
	}

	if cmd.Flags().Lookup("actionStrategy") == nil {
		t.Error("Expected actionStrategy flag to be created")
	}

	if cmd.Flags().Lookup("keepStrategy") == nil {
		t.Error("Expected keepStrategy flag to be created")
	}

	if cmd.Flags().Lookup("trashPath") == nil {
		t.Error("Expected trashPath flag to be created")
	}

	if cmd.Flags().Lookup("workers") == nil {
		t.Error("Expected workers flag to be created")
	}

	if cmd.Flags().Lookup("threshold") == nil {
		t.Error("Expected threshold flag to be created")
	}
}

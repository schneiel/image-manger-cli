package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	coretestutils "github.com/schneiel/ImageManagerGo/core/testutils"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	clitestutils "github.com/schneiel/ImageManagerGo/internal/testutils"
)

func TestCreateCommand(t *testing.T) {
	t.Parallel()

	handler := &clitestutils.MockSortHandler{}
	flagSetup := &clitestutils.MockFlagSetup{}

	cmd := createCommand(
		"test",
		"short desc",
		"long desc",
		handler,
		func() any { return &cmdconfig.SortConfig{} },
		func(cmd *cobra.Command, cfg any) {
			flagSetup.SetupFlags(cmd, cfg)
		},
	)

	assert.NotNil(t, cmd)
	assert.Equal(t, "test", cmd.Use)
	assert.Equal(t, "short desc", cmd.Short)
	assert.Equal(t, "long desc", cmd.Long)
}

func TestNewDedupCommand(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source:         "/test/source",
			ActionStrategy: "dryRun",
			KeepStrategy:   "oldest",
			TrashPath:      "/test/trash",
			Workers:        4,
			Threshold:      95,
		},
	}

	localizer := coretestutils.NewFakeLocalizer()
	handler := &clitestutils.MockDedupHandler{}
	flagSetup := &clitestutils.MockFlagSetup{}

	cmd := NewDedupCommand(cfg, localizer, handler, flagSetup)

	require.NotNil(t, cmd)
	assert.Equal(t, "dedup", cmd.Use)
	assert.Equal(t, "DedupCommandDesc", cmd.Short)
	assert.Equal(t, "DedupCommandLongDesc", cmd.Long)
}

func TestNewDedupCommand_ConfigInitialization(t *testing.T) {
	t.Parallel()

	globalCfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source:         "/global/source",
			ActionStrategy: "moveToTrash",
			KeepStrategy:   "shortestPath",
			TrashPath:      "/global/trash",
			Workers:        8,
			Threshold:      90,
		},
	}

	localizer := coretestutils.NewFakeLocalizer()
	handler := &clitestutils.MockDedupHandler{}
	flagSetup := &clitestutils.MockFlagSetup{
		SetupFlagsFunc: func(_ *cobra.Command, cfg interface{}) {
			// Verify the config has the global values
			dedupCfg, ok := cfg.(*cmdconfig.DedupConfig)
			require.True(t, ok)
			assert.Equal(t, "/global/source", dedupCfg.Source)
			assert.Equal(t, "moveToTrash", dedupCfg.ActionStrategy)
			assert.Equal(t, "shortestPath", dedupCfg.KeepStrategy)
			assert.Equal(t, "/global/trash", dedupCfg.TrashPath)
			assert.Equal(t, 8, dedupCfg.Workers)
			assert.Equal(t, 90, dedupCfg.Threshold)
		},
	}

	cmd := NewDedupCommand(globalCfg, localizer, handler, flagSetup)

	// Trigger flag setup by accessing the command's flag functionality
	_ = cmd.Flag("help") // This should trigger the flag setup function

	require.NotNil(t, cmd)
}

func TestNewSortCommand(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		Sorter: config.SorterConfig{
			Source:         "/test/source",
			Destination:    "/test/dest",
			ActionStrategy: "copy",
		},
	}

	localizer := coretestutils.NewFakeLocalizer()
	handler := &clitestutils.MockSortHandler{}
	flagSetup := &clitestutils.MockFlagSetup{}

	cmd := NewSortCommand(cfg, localizer, handler, flagSetup)

	require.NotNil(t, cmd)
	assert.Equal(t, "sort", cmd.Use)
	assert.Equal(t, "SortCommandDesc", cmd.Short)
	assert.Equal(t, "SortCommandLongDesc", cmd.Long)
}

func TestNewSortCommand_ConfigInitialization(t *testing.T) {
	t.Parallel()

	globalCfg := &config.Config{
		Sorter: config.SorterConfig{
			Source:         "/global/source",
			Destination:    "/global/dest",
			ActionStrategy: "dryRun",
		},
	}

	localizer := coretestutils.NewFakeLocalizer()
	handler := &clitestutils.MockSortHandler{}
	flagSetup := &clitestutils.MockFlagSetup{
		SetupFlagsFunc: func(_ *cobra.Command, cfg interface{}) {
			// Verify the config has the global values
			sortCfg, ok := cfg.(*cmdconfig.SortConfig)
			require.True(t, ok)
			assert.Equal(t, "/global/source", sortCfg.Source)
			assert.Equal(t, "/global/dest", sortCfg.Destination)
			assert.Equal(t, "dryRun", sortCfg.ActionStrategy)
		},
	}

	cmd := NewSortCommand(globalCfg, localizer, handler, flagSetup)

	// Trigger flag setup by accessing the command's flag functionality
	_ = cmd.Flag("help") // This should trigger the flag setup function

	require.NotNil(t, cmd)
}

func TestNewAllCommands(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source:         "/test/dedup",
			ActionStrategy: "dryRun",
			KeepStrategy:   "oldest",
			TrashPath:      "/test/trash",
			Workers:        4,
			Threshold:      95,
		},
		Sorter: config.SorterConfig{
			Source:         "/test/sort",
			Destination:    "/test/dest",
			ActionStrategy: "copy",
		},
	}

	localizer := coretestutils.NewFakeLocalizer()
	dedupHandler := &clitestutils.MockDedupHandler{}
	sortHandler := &clitestutils.MockSortHandler{}
	dedupFlagSetup := &clitestutils.MockFlagSetup{}
	sortFlagSetup := &clitestutils.MockFlagSetup{}

	commands := NewAllCommands(
		cfg,
		localizer,
		dedupHandler,
		sortHandler,
		dedupFlagSetup,
		sortFlagSetup,
	)

	require.Len(t, commands, 2)

	// Check dedup command
	dedupCmd := commands[0]
	assert.Equal(t, "dedup", dedupCmd.Use)
	assert.Equal(t, "DedupCommandDesc", dedupCmd.Short)

	// Check sort command
	sortCmd := commands[1]
	assert.Equal(t, "sort", sortCmd.Use)
	assert.Equal(t, "SortCommandDesc", sortCmd.Short)
}

func TestNewAllCommands_EmptyConfig(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{} // Empty config

	localizer := coretestutils.NewFakeLocalizer()
	dedupHandler := &clitestutils.MockDedupHandler{}
	sortHandler := &clitestutils.MockSortHandler{}
	dedupFlagSetup := &clitestutils.MockFlagSetup{}
	sortFlagSetup := &clitestutils.MockFlagSetup{}

	commands := NewAllCommands(
		cfg,
		localizer,
		dedupHandler,
		sortHandler,
		dedupFlagSetup,
		sortFlagSetup,
	)

	require.Len(t, commands, 2)
	assert.Equal(t, "dedup", commands[0].Use)
	assert.Equal(t, "sort", commands[1].Use)
}

func TestCommandConfigProviders_ReturnNewInstances(t *testing.T) {
	t.Parallel()

	t.Run("DedupCommand config validation", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Deduplicator: config.DeduplicatorConfig{
				Source: "/test/source",
			},
		}

		localizer := coretestutils.NewFakeLocalizer()
		handler := &clitestutils.MockDedupHandler{}

		var capturedConfig *cmdconfig.DedupConfig
		flagSetup := &clitestutils.MockFlagSetup{
			SetupFlagsFunc: func(_ *cobra.Command, cfg interface{}) {
				dedupCfg, ok := cfg.(*cmdconfig.DedupConfig)
				require.True(t, ok)
				capturedConfig = dedupCfg
			},
		}

		cmd := NewDedupCommand(cfg, localizer, handler, flagSetup)

		// Accessing flags should trigger config provider
		_ = cmd.Flag("help")

		// Verify config was provided with correct values
		require.NotNil(t, capturedConfig)
		assert.Equal(t, "/test/source", capturedConfig.Source)
	})

	t.Run("SortCommand config validation", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Sorter: config.SorterConfig{
				Source: "/test/source",
			},
		}

		localizer := coretestutils.NewFakeLocalizer()
		handler := &clitestutils.MockSortHandler{}

		var capturedConfig *cmdconfig.SortConfig
		flagSetup := &clitestutils.MockFlagSetup{
			SetupFlagsFunc: func(_ *cobra.Command, cfg interface{}) {
				sortCfg, ok := cfg.(*cmdconfig.SortConfig)
				require.True(t, ok)
				capturedConfig = sortCfg
			},
		}

		cmd := NewSortCommand(cfg, localizer, handler, flagSetup)

		// Accessing flags should trigger config provider
		_ = cmd.Flag("help")

		// Verify config was provided with correct values
		require.NotNil(t, capturedConfig)
		assert.Equal(t, "/test/source", capturedConfig.Source)
	})
}

// Test that the created commands have the expected structure.
func TestCommandStructure(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{}
	localizer := coretestutils.NewFakeLocalizer()

	t.Run("DedupCommand structure", func(t *testing.T) {
		t.Parallel()
		handler := &clitestutils.MockDedupHandler{}
		flagSetup := &clitestutils.MockFlagSetup{}

		cmd := NewDedupCommand(cfg, localizer, handler, flagSetup)

		assert.Equal(t, "dedup", cmd.Use)
		assert.NotEmpty(t, cmd.Short)
		assert.NotEmpty(t, cmd.Long)
		assert.NotNil(t, cmd.RunE)
	})

	t.Run("SortCommand structure", func(t *testing.T) {
		t.Parallel()
		handler := &clitestutils.MockSortHandler{}
		flagSetup := &clitestutils.MockFlagSetup{}

		cmd := NewSortCommand(cfg, localizer, handler, flagSetup)

		assert.Equal(t, "sort", cmd.Use)
		assert.NotEmpty(t, cmd.Short)
		assert.NotEmpty(t, cmd.Long)
		assert.NotNil(t, cmd.RunE)
	})
}

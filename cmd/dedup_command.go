package cmd

import (
	"ImageManager/config"
	"ImageManager/i18n"
	"ImageManager/image_deduplicator"
	"ImageManager/image_deduplicator/action_strategy"
	"ImageManager/image_deduplicator/keep_strategy"
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
)

// DedupCommand implements the Command interface for the 'dedup' command.
type DedupCommand struct {
	fs           *flag.FlagSet
	commonFlags  CommonFlags
	targetDir    string
	trashPath    string
	workers      int
	keepStrategy string
	threshold    int
}

// NewDedupCommand creates a new DedupCommand.
func NewDedupCommand() *DedupCommand {
	dc := &DedupCommand{
		fs: flag.NewFlagSet("dedup", flag.ExitOnError),
	}
	// Add command-specific flags
	dc.fs.StringVar(&dc.targetDir, "source", "", "Directory to scan for duplicates. (required)")
	dc.fs.StringVar(&dc.trashPath, "trash-dir", "", i18n.T("DedupTrashDirUsage"))
	dc.fs.IntVar(&dc.workers, "workers", runtime.NumCPU(), i18n.T("DedupWorkersUsage"))
	dc.fs.StringVar(&dc.keepStrategy, "keep", "", i18n.T("DedupKeepUsage"))
	dc.fs.IntVar(&dc.threshold, "threshold", 5, i18n.T("DedupThresholdUsage"))

	// Add the common, polymorphic flags with "
	AddCommonFlags(dc.fs, &dc.commonFlags)

	return dc
}

// Name returns the command's name.
func (c *DedupCommand) Name() string { return c.fs.Name() }

// Init initializes the command with arguments.
func (c *DedupCommand) Init(args []string) { c.fs.Parse(args) }

// Usage returns the command's flag set.
func (c *DedupCommand) Usage() *flag.FlagSet { return c.fs }

// Run executes the dedup command's logic.
func (c *DedupCommand) Run() error {
	// Load configuration from file using the common config flag
	appConfig, err := config.LoadConfig(c.commonFlags.ConfigPath)
	if err != nil {
		return fmt.Errorf("error loading config file: %w", err)
	}

	// Override with config file value if flag is not provided
	if c.targetDir == "" {
		c.targetDir = appConfig.Source
	}

	// Override with config file value if flag is not provided
	if !c.commonFlags.DryRun {
		c.commonFlags.DryRun = appConfig.DryRun
	}

	if c.targetDir == "" {
		return fmt.Errorf(i18n.T("DedupTargetDirMissing"))
	}

	if c.trashPath == "" {
		c.trashPath = filepath.Join(c.targetDir, ".trash")
	}

	// --- Determine Strategy ---
	// Priority: Flag > Config File > Default
	finalKeepStrategy := appConfig.Deduplicator.KeepStrategy
	if c.keepStrategy != "" { // Flag overrides config
		finalKeepStrategy = c.keepStrategy
	}

	var keepStrategy keep_strategy.KeepStrategy
	switch finalKeepStrategy {
	case "short-path", "keepShortestPath":
		keepStrategy = &keep_strategy.ShortestPathStrategy{}
	case "oldest", "keepOldest":
		keepStrategy = &keep_strategy.OldestFileStrategy{}
	default:
		return fmt.Errorf(i18n.T("UnknownKeepStrategy", map[string]interface{}{"Strategy": finalKeepStrategy}))
	}

	// Determine action strategy using the common dry-run flag
	var actionStrategy action_strategy.ActionStrategy
	if c.commonFlags.DryRun {
		actionStrategy = &action_strategy.DryRunStrategy{
			CsvPath: "dedup_dry_run_log.csv",
			Header:  []string{i18n.T("CsvHeaderOriginal"), i18n.T("CsvHeaderDuplicate")},
		}
	} else {
		actionStrategy = &keep_strategy.MoveToTrashStrategy{TrashDir: c.trashPath}
	}

	dedupConfig := imagededuplicator.Config{
		TargetPath:   c.targetDir,
		NumWorkers:   c.workers,
		Threshold:    c.threshold,
		KeepStrategy: keepStrategy,
		Action:       actionStrategy,
	}

	deduplicator := imagededuplicator.NewDeduplicator(dedupConfig)
	return deduplicator.FindAndRemoveDuplicates()
}

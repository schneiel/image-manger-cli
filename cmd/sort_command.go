package cmd

import (
	"ImageManager/config"
	"ImageManager/i18n"
	"ImageManager/image_sorter"
	"ImageManager/image_sorter/action_strategy"
	"flag"
	"fmt"
	"os"
)

// SortCommand implements the Command interface for the 'sort' command.
type SortCommand struct {
	fs          *flag.FlagSet
	commonFlags CommonFlags
	sourceDir   string
	destDir     string
}

// NewSortCommand creates a new SortCommand.
func NewSortCommand() *SortCommand {
	sc := &SortCommand{
		fs: flag.NewFlagSet("sort", flag.ExitOnError),
	}
	// Add command-specific flags
	sc.fs.StringVar(&sc.sourceDir, "source", "", i18n.T("SortSrcUsage"))
	sc.fs.StringVar(&sc.destDir, "destination", "", i18n.T("SortDestUsage"))

	// Add the common, polymorphic flags
	AddCommonFlags(sc.fs, &sc.commonFlags)
	return sc
}

// Name returns the command's name.
func (c *SortCommand) Name() string { return c.fs.Name() }

// Init initializes the command with arguments.
func (c *SortCommand) Init(args []string) { c.fs.Parse(args) }

// Usage returns the command's flag set.
func (c *SortCommand) Usage() *flag.FlagSet { return c.fs }

// Run executes the sort command's logic.
func (c *SortCommand) Run() error {
	// Load configuration using the common config flag
	app_config, err := config.LoadConfig(c.commonFlags.ConfigPath)
	if err != nil {
		return fmt.Errorf("error loading config file: %w", err)
	}

	// Override with config file values if flags are not provided
	if c.sourceDir == "" {
		c.sourceDir = app_config.Source
	}

	if c.destDir == "" {
		c.destDir = app_config.Destination
	}

	// Override with config file value if flag is not provided
	if !c.commonFlags.DryRun {
		c.commonFlags.DryRun = app_config.DryRun
	}

	if c.sourceDir == "" || c.destDir == "" {
		return fmt.Errorf(i18n.T("SortSrcDestRequired"))
	}

	if _, err := os.Stat(c.destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.destDir, 0755); err != nil {
			return fmt.Errorf(i18n.T("SortCreatingDestFailed", map[string]interface{}{"Path": c.destDir, "Error": err}))
		}
	}

	fmt.Println(i18n.T("SortStarting", map[string]interface{}{"Src": c.sourceDir, "Dest": c.destDir}))

	imageProcessor, err := imagesorter.NewImageProcessor(*app_config)
	if err != nil {
		return fmt.Errorf("error initializing ImageProcessor: %w", err)
	}

	fmt.Println(i18n.T("SortScanningSource", map[string]interface{}{"Path": c.sourceDir}))
	images := imageProcessor.HandleImages(c.sourceDir)
	fmt.Println(i18n.T("SortImagesFound", map[string]interface{}{"Count": len(images)}))

	if len(images) > 0 {
		// Determine which action strategy to use based on the common dry-run flag
		var strategy action_strategy.SortActionStrategy
		if c.commonFlags.DryRun {
			strategy = action_strategy.NewDryRunStrategy("sort_dry_run_log.csv")
		} else {
			strategy = action_strategy.NewCopyStrategy()
		}

		sorter := imagesorter.NewSorter(c.destDir, strategy)

		fmt.Println(i18n.T("SortSortingAndCopying"))
		if err := sorter.SortImages(images); err != nil {
			return err
		}
		fmt.Println(i18n.T("SortFinished"))
	} else {
		fmt.Println(i18n.T("SortNoImages"))
	}
	return nil
}

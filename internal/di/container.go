// Package di provides dependency injection container and builders for the CLI application.
package di

import (
	"embed"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/date"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/exif"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	coretimepkg "github.com/schneiel/ImageManagerGo/core/infrastructure/time"
	"github.com/schneiel/ImageManagerGo/core/processing/dedup"
	"github.com/schneiel/ImageManagerGo/core/processing/sort"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupaction"
	"github.com/schneiel/ImageManagerGo/core/strategies/shared"
	"github.com/schneiel/ImageManagerGo/core/strategies/sortaction"
	"github.com/schneiel/ImageManagerGo/core/task"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	"github.com/schneiel/ImageManagerGo/internal/handlers"
	"github.com/schneiel/ImageManagerGo/internal/localization"
)

// Container holds all the dependencies for the application.
type Container struct {
	// Core dependencies.
	Args         []string
	Localizer    i18n.Localizer
	FileUtils    filesystem.FileUtils
	FileReader   config.FileReader
	Parser       config.Parser
	FileSystem   filesystem.FileSystem
	TimeProvider coretimepkg.TimeProvider
	LocalesFS    *embed.FS

	// Configuration.
	Config *config.Config

	// Logging.
	SortLogger  log.Logger
	DedupLogger log.Logger

	// EXIF processing.
	ExifDecoder exif.Decoder
	ExifReader  *exif.DefaultReader

	// Date processing.
	DateProcessor date.DateProcessor

	// Image sorting.
	ImageFinder    sort.ImageFinder
	ImageAnalyzer  sort.ImageAnalyzer
	ImageProcessor sort.ImageProcessor

	// Image deduplication.
	SizeScanner     dedup.Scanner
	PHasher         dedup.Hasher
	DistanceGrouper dedup.Grouper

	// Tasks.
	SortTask  task.Task
	DedupTask task.Task

	// Services.
	TaskExecutor         handlers.TaskExecutor
	SortConfigApplier    handlers.ConfigApplier[*cmdconfig.SortConfig]
	SortConfigValidator  handlers.ConfigValidator[*cmdconfig.SortConfig]
	DedupConfigApplier   handlers.ConfigApplier[*cmdconfig.DedupConfig]
	DedupConfigValidator handlers.ConfigValidator[*cmdconfig.DedupConfig]

	// Handlers.
	SortHandler  handlers.CommandHandler
	DedupHandler handlers.CommandHandler

	// Flag Setup.
	SortFlagSetup  handlers.FlagSetup
	DedupFlagSetup handlers.FlagSetup

	// Localization.
	CommandLocalizer localization.BaseLocalizer
}

// ContainerBuilder provides a fluent interface for building the container.
type ContainerBuilder struct {
	container *Container
	err       error

	// Builders for different concerns.
	coreBuilder       *CoreBuilder
	configBuilder     *ConfigBuilder
	loggingBuilder    *LoggingBuilder
	processingBuilder *ProcessingBuilder
	handlersBuilder   *HandlersBuilder
}

// NewContainerBuilder creates a new container builder.
func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		container:         &Container{},
		coreBuilder:       NewCoreBuilder(),
		configBuilder:     NewConfigBuilder(),
		loggingBuilder:    NewLoggingBuilder(),
		processingBuilder: NewProcessingBuilder(),
		handlersBuilder:   NewHandlersBuilder(),
	}
}

// WithArgs sets the command line arguments.
func (b *ContainerBuilder) WithArgs(args []string) *ContainerBuilder {
	if b.err != nil {
		return b
	}
	b.container.Args = args

	return b
}

// WithLocalesFS sets the embedded locales filesystem.
func (b *ContainerBuilder) WithLocalesFS(localesFS *embed.FS) *ContainerBuilder {
	if b.err != nil {
		return b
	}
	b.container.LocalesFS = localesFS

	return b
}

// Build creates and returns the container with all dependencies initialized.
func (b *ContainerBuilder) Build() (*Container, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Build core dependencies.
	core, err := b.coreBuilder.BuildCore(b.container.Args, b.container.LocalesFS)
	if err != nil {
		return nil, err
	}
	b.copyCoreDependencies(core)

	// Build configuration.
	cfg, err := b.configBuilder.BuildConfig(b.container.Args, core)
	if err != nil {
		return nil, err
	}
	b.container.Config = cfg

	// Build logging.
	logging, err := b.loggingBuilder.BuildLogging(core, cfg)
	if err != nil {
		return nil, err
	}
	b.copyLoggingDependencies(logging)

	// Build processing.
	processing, err := b.processingBuilder.BuildProcessing(core, cfg, logging)
	if err != nil {
		return nil, err
	}
	b.copyProcessingDependencies(processing)

	// Build tasks.
	err = b.buildTasks(core, cfg, logging, processing)
	if err != nil {
		return nil, err
	}

	// Build handlers.
	handles, err := b.handlersBuilder.BuildHandlers(
		b.container.Args, core, cfg, logging,
		b.container.SortTask, b.container.DedupTask,
	)
	if err != nil {
		return nil, err
	}
	b.copyHandlersDependencies(handles)

	return b.container, nil
}

// copyCoreDependencies copies core dependencies to the container.
func (b *ContainerBuilder) copyCoreDependencies(core *CoreDependencies) {
	b.container.Parser = core.Parser
	b.container.FileReader = core.FileReader
	b.container.FileSystem = core.FileSystem
	b.container.TimeProvider = core.TimeProvider
	b.container.Localizer = core.Localizer
	b.container.FileUtils = core.FileUtils
}

// copyLoggingDependencies copies logging dependencies to the container.
func (b *ContainerBuilder) copyLoggingDependencies(logging *LoggingDependencies) {
	b.container.SortLogger = logging.SortLogger
	b.container.DedupLogger = logging.DedupLogger
}

// copyProcessingDependencies copies processing dependencies to the container.
func (b *ContainerBuilder) copyProcessingDependencies(processing *ProcessingDependencies) {
	b.container.ExifDecoder = processing.ExifDecoder
	b.container.ExifReader = processing.ExifReader
	b.container.DateProcessor = processing.DateProcessor
	b.container.ImageFinder = processing.ImageFinder
	b.container.ImageAnalyzer = processing.ImageAnalyzer
	b.container.ImageProcessor = processing.ImageProcessor
	b.container.SizeScanner = processing.SizeScanner
	b.container.PHasher = processing.PHasher
	b.container.DistanceGrouper = processing.DistanceGrouper
}

// copyHandlersDependencies copies handlers dependencies to the container.
func (b *ContainerBuilder) copyHandlersDependencies(deps *HandlersDependencies) {
	b.container.TaskExecutor = deps.TaskExecutor
	b.container.SortConfigApplier = deps.SortConfigApplier
	b.container.SortConfigValidator = deps.SortConfigValidator
	b.container.DedupConfigApplier = deps.DedupConfigApplier
	b.container.DedupConfigValidator = deps.DedupConfigValidator
	b.container.CommandLocalizer = deps.CommandLocalizer
	b.container.SortFlagSetup = deps.SortFlagSetup
	b.container.DedupFlagSetup = deps.DedupFlagSetup
	b.container.SortHandler = deps.SortHandler
	b.container.DedupHandler = deps.DedupHandler
}

// buildTasks creates the task instances.
func (b *ContainerBuilder) buildTasks(
	core *CoreDependencies,
	cfg *config.Config,
	logging *LoggingDependencies,
	processing *ProcessingDependencies,
) error {
	err := b.buildSortTask(core, cfg, logging, processing)
	if err != nil {
		return fmt.Errorf("failed to build sort task: %w", err)
	}

	err = b.buildDedupTask(core, cfg, logging, processing)
	if err != nil {
		return fmt.Errorf("failed to build dedup task: %w", err)
	}

	return nil
}

func (b *ContainerBuilder) buildSortTask(
	core *CoreDependencies,
	cfg *config.Config,
	logging *LoggingDependencies,
	processing *ProcessingDependencies,
) error {
	sorterStrategyFactory := func() sort.Strategy {
		strategy, err := b.processingBuilder.BuildSorterActionStrategy(cfg, core, logging)
		if err != nil {
			fallbackStrategy, fallbackErr := sortaction.NewDefaultDryRunStrategy(logging.SortLogger, core.Localizer)
			if fallbackErr != nil {
				panic("failed to create fallback strategy: " + fallbackErr.Error())
			}
			return fallbackStrategy
		}
		return strategy
	}

	sortTask, err := task.NewSortTaskBuilder().
		WithConfig(cfg).
		WithLogger(logging.SortLogger).
		WithLocalizer(core.Localizer).
		WithImageProcessor(processing.ImageProcessor).
		WithFileUtils(core.FileUtils).
		WithStrategyFactory(sorterStrategyFactory).
		Build()
	if err != nil {
		return err
	}
	b.container.SortTask = sortTask
	return nil
}

func (b *ContainerBuilder) buildDedupTask(
	core *CoreDependencies,
	cfg *config.Config,
	logging *LoggingDependencies,
	processing *ProcessingDependencies,
) error {
	dedupStrategyFactory := func() dedupaction.Strategy {
		strategy, err := b.processingBuilder.BuildDedupActionStrategy(cfg, core, logging)
		if err != nil {
			fallbackStrategy, fallbackErr := dedupaction.NewDefaultDryRunStrategy(shared.ActionConfig{
				Logger:    logging.DedupLogger,
				Localizer: core.Localizer,
				FileUtils: core.FileUtils,
				TrashPath: cfg.Deduplicator.TrashPath,
				CsvPath:   cfg.Files.DedupDryRunLog,
			})
			if fallbackErr != nil {
				panic("failed to create fallback dedup strategy: " + fallbackErr.Error())
			}
			return fallbackStrategy
		}
		return strategy
	}

	dedupTask, err := task.NewDedupTaskBuilder().
		WithConfig(cfg).
		WithLogger(logging.DedupLogger).
		WithLocalizer(core.Localizer).
		WithFilesystem(core.FileSystem).
		WithScanner(processing.SizeScanner).
		WithHasher(processing.PHasher).
		WithGrouper(processing.DistanceGrouper).
		WithKeepFunc(b.processingBuilder.BuildKeepFunc(cfg, core)).
		WithStrategyFactory(dedupStrategyFactory).
		Build()
	if err != nil {
		return err
	}
	b.container.DedupTask = dedupTask
	return nil
}

// NewContainer creates a new container using the builder pattern.
func NewContainer(args []string, localesFS *embed.FS) (*Container, error) {
	return NewContainerBuilder().
		WithArgs(args).
		WithLocalesFS(localesFS).
		Build()
}

package di

import (
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/task"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
	"github.com/schneiel/ImageManagerGo/internal/handlers"
	"github.com/schneiel/ImageManagerGo/internal/localization"
	"github.com/schneiel/ImageManagerGo/internal/services"
)

// HandlersBuilder handles CLI handlers and services setup.
type HandlersBuilder struct {
	argParser *ArgumentParser
}

// NewHandlersBuilder creates a new handlers builder.
func NewHandlersBuilder() *HandlersBuilder {
	return &HandlersBuilder{
		argParser: NewArgumentParser(),
	}
}

// BuildHandlers sets up CLI handlers and services.
func (hb *HandlersBuilder) BuildHandlers( //nolint:funlen // dependency injection setup function
	args []string,
	core *CoreDependencies,
	cfg *config.Config,
	logging *LoggingDependencies,
	sortTask task.Task,
	dedupTask task.Task,
) (*HandlersDependencies, error) {
	deps := &HandlersDependencies{}

	// Task executor
	deps.TaskExecutor = handlers.NewDefaultTaskExecutor(
		logging.SortLogger,
		core.Localizer,
		core.FileUtils,
		sortTask,
		dedupTask,
	)

	// Config services
	sortConfigApplier, err := services.NewSortConfigApplier(logging.SortLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create sort config applier: %w", err)
	}
	deps.SortConfigApplier = sortConfigApplier
	
	deps.SortConfigValidator = services.NewSortConfigValidator(core.Localizer)
	
	dedupConfigApplier, err := services.NewDedupConfigApplier(logging.DedupLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dedup config applier: %w", err)
	}
	deps.DedupConfigApplier = dedupConfigApplier
	
	deps.DedupConfigValidator = services.NewDedupConfigValidator(core.Localizer)

	// Command localizer
	language := hb.argParser.ExtractLanguage(args)
	switch language {
	case "de":
		deps.CommandLocalizer = localization.NewGermanLocalizer(core.Localizer)
	case "en":
		deps.CommandLocalizer = localization.NewEnglishLocalizer(core.Localizer)
	default:
		deps.CommandLocalizer = localization.NewEnglishLocalizer(core.Localizer)
	}

	// Flag setup
	deps.SortFlagSetup = handlers.NewSortFlagSetup(core.Localizer)
	deps.DedupFlagSetup = handlers.NewDedupFlagSetup(core.Localizer)

	// Command handlers using builder pattern
	baseSortHandler, err := handlers.NewBaseHandler(logging.SortLogger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base sort handler: %w", err)
	}
	sortHandler, err := handlers.NewSortHandlerWithOptions(
		handlers.WithSortBaseHandler(baseSortHandler),
		handlers.WithSortTaskExecutor(deps.TaskExecutor),
		handlers.WithSortConfigApplier(deps.SortConfigApplier),
		handlers.WithSortConfigValidator(deps.SortConfigValidator),
		handlers.WithSortLocalizer(core.Localizer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sort handler: %w", err)
	}
	deps.SortHandler = sortHandler

	baseDedupHandler, err := handlers.NewBaseHandler(logging.DedupLogger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base dedup handler: %w", err)
	}
	dedupHandler, err := handlers.NewDedupHandlerWithOptions(
		handlers.WithBaseHandler(baseDedupHandler),
		handlers.WithTaskExecutor(deps.TaskExecutor),
		handlers.WithDedupConfigApplier(deps.DedupConfigApplier),
		handlers.WithDedupConfigValidator(deps.DedupConfigValidator),
		handlers.WithLocalizer(core.Localizer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dedup handler: %w", err)
	}
	deps.DedupHandler = dedupHandler

	return deps, nil
}

// Build is an alias for BuildHandlers to satisfy the standard builder pattern
// This method provides a simplified interface but requires parameters.
func (hb *HandlersBuilder) Build() (*HandlersDependencies, error) {
	// For the standard Build() method, we need to handle the lack of parameters
	return nil, errors.New(
		"Build() requires parameters - use BuildHandlers(args, core, cfg, logging, " +
			"sortTask, dedupTask) instead",
	)
}

// HandlersDependencies holds CLI handlers and services.
type HandlersDependencies struct {
	TaskExecutor         handlers.TaskExecutor
	SortConfigApplier    handlers.ConfigApplier[*cmdconfig.SortConfig]
	SortConfigValidator  handlers.ConfigValidator[*cmdconfig.SortConfig]
	DedupConfigApplier   handlers.ConfigApplier[*cmdconfig.DedupConfig]
	DedupConfigValidator handlers.ConfigValidator[*cmdconfig.DedupConfig]
	CommandLocalizer     localization.BaseLocalizer
	SortFlagSetup        handlers.FlagSetup
	DedupFlagSetup       handlers.FlagSetup
	SortHandler          handlers.CommandHandler
	DedupHandler         handlers.CommandHandler
}

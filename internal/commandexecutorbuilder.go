package cli

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/internal/handlers"
	"github.com/schneiel/ImageManagerGo/internal/localization"
)

// CommandExecutorBuilder provides a fluent interface for building DefaultCommandExecutor instances
// with complex dependency injection requirements.
type CommandExecutorBuilder struct {
	executor *DefaultCommandExecutor
	err      error
}

// NewCommandExecutorBuilder creates a new builder instance for constructing DefaultCommandExecutor.
func NewCommandExecutorBuilder() *CommandExecutorBuilder {
	return &CommandExecutorBuilder{
		executor: &DefaultCommandExecutor{},
	}
}

// WithArgs sets the command line arguments for the executor.
func (b *CommandExecutorBuilder) WithArgs(args []string) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	b.executor.args = args

	return b
}

// WithLocalizer sets the localizer for the executor.
func (b *CommandExecutorBuilder) WithLocalizer(localizer i18n.Localizer) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if localizer == nil {
		b.err = errors.New("localizer cannot be nil")

		return b
	}
	b.executor.localizer = localizer

	return b
}

// WithFileReader sets the file reader for the executor.
func (b *CommandExecutorBuilder) WithFileReader(fileReader config.FileReader) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if fileReader == nil {
		b.err = errors.New("fileReader cannot be nil")

		return b
	}
	b.executor.fileReader = fileReader

	return b
}

// WithParser sets the config parser for the executor.
func (b *CommandExecutorBuilder) WithParser(parser config.Parser) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if parser == nil {
		b.err = errors.New("parser cannot be nil")

		return b
	}
	b.executor.parser = parser

	return b
}

// WithConfig sets the configuration for the executor.
func (b *CommandExecutorBuilder) WithConfig(cfg *config.Config) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if cfg == nil {
		b.err = errors.New("config cannot be nil")

		return b
	}
	b.executor.config = cfg

	return b
}

// WithHandlers sets both sort and dedup handlers for the executor.
func (b *CommandExecutorBuilder) WithHandlers(
	sortHandler, dedupHandler handlers.CommandHandler,
) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if sortHandler == nil {
		b.err = errors.New("sortHandler cannot be nil")

		return b
	}
	if dedupHandler == nil {
		b.err = errors.New("dedupHandler cannot be nil")

		return b
	}
	b.executor.sortHandler = sortHandler
	b.executor.dedupHandler = dedupHandler

	return b
}

// WithSortHandler sets the sort handler for the executor.
func (b *CommandExecutorBuilder) WithSortHandler(sortHandler handlers.CommandHandler) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if sortHandler == nil {
		b.err = errors.New("sortHandler cannot be nil")

		return b
	}
	b.executor.sortHandler = sortHandler

	return b
}

// WithDedupHandler sets the dedup handler for the executor.
func (b *CommandExecutorBuilder) WithDedupHandler(dedupHandler handlers.CommandHandler) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if dedupHandler == nil {
		b.err = errors.New("dedupHandler cannot be nil")

		return b
	}
	b.executor.dedupHandler = dedupHandler

	return b
}

// WithFlagSetups sets both sort and dedup flag setups for the executor.
func (b *CommandExecutorBuilder) WithFlagSetups(
	sortFlagSetup, dedupFlagSetup handlers.FlagSetup,
) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if sortFlagSetup == nil {
		b.err = errors.New("sortFlagSetup cannot be nil")

		return b
	}
	if dedupFlagSetup == nil {
		b.err = errors.New("dedupFlagSetup cannot be nil")

		return b
	}
	b.executor.sortFlagSetup = sortFlagSetup
	b.executor.dedupFlagSetup = dedupFlagSetup

	return b
}

// WithSortFlagSetup sets the sort flag setup for the executor.
func (b *CommandExecutorBuilder) WithSortFlagSetup(sortFlagSetup handlers.FlagSetup) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if sortFlagSetup == nil {
		b.err = errors.New("sortFlagSetup cannot be nil")

		return b
	}
	b.executor.sortFlagSetup = sortFlagSetup

	return b
}

// WithDedupFlagSetup sets the dedup flag setup for the executor.
func (b *CommandExecutorBuilder) WithDedupFlagSetup(dedupFlagSetup handlers.FlagSetup) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if dedupFlagSetup == nil {
		b.err = errors.New("dedupFlagSetup cannot be nil")

		return b
	}
	b.executor.dedupFlagSetup = dedupFlagSetup

	return b
}

// WithCommandLocalizer sets the command localizer for the executor.
func (b *CommandExecutorBuilder) WithCommandLocalizer(
	commandLocalizer localization.BaseLocalizer,
) *CommandExecutorBuilder {
	if b.err != nil {
		return b
	}
	if commandLocalizer == nil {
		b.err = errors.New("commandLocalizer cannot be nil")

		return b
	}
	b.executor.commandLocalizer = commandLocalizer

	return b
}

// Build constructs the final DefaultCommandExecutor with validation.
func (b *CommandExecutorBuilder) Build() (*DefaultCommandExecutor, error) {
	if b.err != nil {
		return nil, b.err
	}

	err := b.validateRequiredDependencies()
	if err != nil {
		return nil, err
	}

	// Args can be empty for testing, but should be set
	if b.executor.args == nil {
		b.executor.args = []string{}
	}

	return b.executor, nil
}

// validateRequiredDependencies validates all required dependencies are set.
func (b *CommandExecutorBuilder) validateRequiredDependencies() error {
	requiredFields := map[string]interface{}{
		"localizer":  b.executor.localizer,
		"fileReader": b.executor.fileReader,
		"parser":     b.executor.parser,
		"github.com/schneiel/ImageManagerGo/core/config": b.executor.config,
		"sortHandler":      b.executor.sortHandler,
		"dedupHandler":     b.executor.dedupHandler,
		"sortFlagSetup":    b.executor.sortFlagSetup,
		"dedupFlagSetup":   b.executor.dedupFlagSetup,
		"commandLocalizer": b.executor.commandLocalizer,
	}

	for fieldName, fieldValue := range requiredFields {
		if fieldValue == nil {
			return errors.New(fieldName + " is required")
		}
	}

	return nil
}

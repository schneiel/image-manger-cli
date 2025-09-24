//nolint:dupl // These files have similar patterns due to type aliases and builders
package handlers

import (
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// DedupHandlerConfig is a type alias for the generic handler configuration.
type DedupHandlerConfig = HandlerConfig[*cmdconfig.DedupConfig]

// DedupHandlerOption is a type alias for the generic handler option.
type DedupHandlerOption = HandlerOption[*cmdconfig.DedupConfig]

// WithBaseHandler sets the base handler for the dedup handler.
var WithBaseHandler = WithHandlerBaseHandler[*cmdconfig.DedupConfig]

// WithTaskExecutor sets the task executor for the dedup handler.
var WithTaskExecutor = WithHandlerExecutor[*cmdconfig.DedupConfig]

// WithDedupConfigApplier sets the config applier for the dedup handler.
var WithDedupConfigApplier = WithHandlerApplier[*cmdconfig.DedupConfig]

// WithDedupConfigValidator sets the config validator for the dedup handler.
var WithDedupConfigValidator = WithHandlerValidator[*cmdconfig.DedupConfig]

// WithLocalizer sets the localizer for the dedup handler.
var WithLocalizer = WithHandlerLocalizer[*cmdconfig.DedupConfig]

// createDedupHandler is the creator function for DedupHandler instances.
func createDedupHandler(
	baseHandler *BaseHandler,
	executor TaskExecutor,
	applier ConfigApplier[*cmdconfig.DedupConfig],
	validator ConfigValidator[*cmdconfig.DedupConfig],
	localizer i18n.Localizer,
) *DedupHandler {
	// Validate required dependencies
	if baseHandler == nil {
		panic("baseHandler cannot be nil")
	}
	if executor == nil {
		panic("executor cannot be nil")
	}
	if applier == nil {
		panic("applier cannot be nil")
	}
	if validator == nil {
		panic("validator cannot be nil")
	}
	if localizer == nil {
		panic("localizer cannot be nil")
	}

	return &DedupHandler{
		BaseHandler: baseHandler,
		Executor:    executor,
		Applier:     applier,
		Validator:   validator,
		Localizer:   localizer,
		Config:      nil, // Initialize as nil
	}
}

// NewDedupHandler creates a new dedup handler with all required dependencies.
func NewDedupHandler(
	baseHandler *BaseHandler,
	executor TaskExecutor,
	applier ConfigApplier[*cmdconfig.DedupConfig],
	validator ConfigValidator[*cmdconfig.DedupConfig],
	localizer i18n.Localizer,
) *DedupHandler {
	return createDedupHandler(baseHandler, executor, applier, validator, localizer)
}

// NewDedupHandlerWithOptions creates a new dedup handler using functional options.
func NewDedupHandlerWithOptions(options ...DedupHandlerOption) (*DedupHandler, error) {
	return NewHandlerWithOptions(createDedupHandler, options...)
}

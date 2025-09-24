//nolint:dupl // These files have similar patterns due to type aliases and builders
package handlers

import (
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	cmdconfig "github.com/schneiel/ImageManagerGo/internal/config"
)

// SortHandlerConfig is a type alias for the generic handler configuration.
type SortHandlerConfig = HandlerConfig[*cmdconfig.SortConfig]

// SortHandlerOption is a type alias for the generic handler option.
type SortHandlerOption = HandlerOption[*cmdconfig.SortConfig]

// WithSortBaseHandler sets the base handler for the sort handler.
var WithSortBaseHandler = WithHandlerBaseHandler[*cmdconfig.SortConfig]

// WithSortTaskExecutor sets the task executor for the sort handler.
var WithSortTaskExecutor = WithHandlerExecutor[*cmdconfig.SortConfig]

// WithSortConfigApplier sets the config applier for the sort handler.
var WithSortConfigApplier = WithHandlerApplier[*cmdconfig.SortConfig]

// WithSortConfigValidator sets the config validator for the sort handler.
var WithSortConfigValidator = WithHandlerValidator[*cmdconfig.SortConfig]

// WithSortLocalizer sets the localizer for the sort handler.
var WithSortLocalizer = WithHandlerLocalizer[*cmdconfig.SortConfig]

// createSortHandler is the creator function for SortHandler instances.
func createSortHandler(
	baseHandler *BaseHandler,
	executor TaskExecutor,
	applier ConfigApplier[*cmdconfig.SortConfig],
	validator ConfigValidator[*cmdconfig.SortConfig],
	localizer i18n.Localizer,
) *SortHandler {
	return &SortHandler{
		BaseHandler: baseHandler,
		Executor:    executor,
		Applier:     applier,
		Validator:   validator,
		Localizer:   localizer,
		Config:      &cmdconfig.SortConfig{}, // Initialize empty config
	}
}

// NewSortHandler creates a new sort handler with all required dependencies.
func NewSortHandler(
	baseHandler *BaseHandler,
	executor TaskExecutor,
	applier ConfigApplier[*cmdconfig.SortConfig],
	validator ConfigValidator[*cmdconfig.SortConfig],
	localizer i18n.Localizer,
) *SortHandler {
	return createSortHandler(baseHandler, executor, applier, validator, localizer)
}

// NewSortHandlerWithOptions creates a new sort handler using functional options.
func NewSortHandlerWithOptions(options ...SortHandlerOption) (*SortHandler, error) {
	return NewHandlerWithOptions(createSortHandler, options...)
}

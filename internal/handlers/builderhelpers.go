package handlers

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// HandlerConfig represents a generic handler configuration.
type HandlerConfig[T any] struct {
	BaseHandler *BaseHandler
	Executor    TaskExecutor
	Applier     ConfigApplier[T]
	Validator   ConfigValidator[T]
	Localizer   i18n.Localizer
}

// HandlerOption is a functional option for configuring handlers.
type HandlerOption[T any] func(*HandlerConfig[T]) error

// WithHandlerBaseHandler sets the base handler.
func WithHandlerBaseHandler[T any](baseHandler *BaseHandler) HandlerOption[T] {
	return func(config *HandlerConfig[T]) error {
		if baseHandler == nil {
			return errors.New("baseHandler cannot be nil")
		}
		config.BaseHandler = baseHandler

		return nil
	}
}

// WithHandlerExecutor sets the task executor.
func WithHandlerExecutor[T any](executor TaskExecutor) HandlerOption[T] {
	return func(config *HandlerConfig[T]) error {
		if executor == nil {
			return errors.New("executor cannot be nil")
		}
		config.Executor = executor

		return nil
	}
}

// WithHandlerApplier sets the config applier.
func WithHandlerApplier[T any](applier ConfigApplier[T]) HandlerOption[T] {
	return func(config *HandlerConfig[T]) error {
		if applier == nil {
			return errors.New("applier cannot be nil")
		}
		config.Applier = applier

		return nil
	}
}

// WithHandlerValidator sets the config validator.
func WithHandlerValidator[T any](validator ConfigValidator[T]) HandlerOption[T] {
	return func(config *HandlerConfig[T]) error {
		if validator == nil {
			return errors.New("validator cannot be nil")
		}
		config.Validator = validator

		return nil
	}
}

// WithHandlerLocalizer sets the localizer.
func WithHandlerLocalizer[T any](localizer i18n.Localizer) HandlerOption[T] {
	return func(config *HandlerConfig[T]) error {
		if localizer == nil {
			return errors.New("localizer cannot be nil")
		}
		config.Localizer = localizer

		return nil
	}
}

// ValidateHandlerConfig validates that all required fields are set.
func ValidateHandlerConfig[T any](config *HandlerConfig[T]) error {
	if config.BaseHandler == nil {
		return errors.New("baseHandler is required")
	}
	if config.Executor == nil {
		return errors.New("executor is required")
	}
	if config.Applier == nil {
		return errors.New("applier is required")
	}
	if config.Validator == nil {
		return errors.New("validator is required")
	}
	if config.Localizer == nil {
		return errors.New("localizer is required")
	}

	return nil
}

// HandlerCreator is a function type for creating specific handler instances.
type HandlerCreator[T any, H any] func(
	baseHandler *BaseHandler,
	executor TaskExecutor,
	applier ConfigApplier[T],
	validator ConfigValidator[T],
	localizer i18n.Localizer,
) H

// NewHandlerFromConfig creates a handler from a validated config using the provided creator function.
func NewHandlerFromConfig[T any, H any](
	config *HandlerConfig[T],
	creator HandlerCreator[T, H],
) H { //nolint:ireturn // generic factory returns interface by design
	return creator(
		config.BaseHandler,
		config.Executor,
		config.Applier,
		config.Validator,
		config.Localizer,
	)
}

// NewHandlerWithOptions creates a handler using functional options and a creator function.
func NewHandlerWithOptions[T any, H any](
	creator HandlerCreator[T, H],
	options ...HandlerOption[T],
) (H, error) { //nolint:ireturn // generic factory returns interface by design
	config := &HandlerConfig[T]{}

	// Apply all options
	for _, option := range options {
		err := option(config)
		if err != nil {
			var zero H

			return zero, err
		}
	}

	// Validate all required dependencies are provided
	err := ValidateHandlerConfig(config)
	if err != nil {
		var zero H

		return zero, err
	}

	return NewHandlerFromConfig(config, creator), nil
}

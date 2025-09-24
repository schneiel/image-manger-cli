// Package handlers provides command handlers for CLI operations
package handlers

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// BaseHandler provides a base implementation for command handlers.
type BaseHandler struct {
	Logger log.Logger
	Config *config.Config
}

// NewBaseHandler creates a new BaseHandler.
func NewBaseHandler(logger log.Logger, cfg *config.Config) (*BaseHandler, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if cfg == nil {
		return nil, errors.New("cfg cannot be nil")
	}

	return &BaseHandler{
		Logger: logger,
		Config: cfg,
	}, nil
}

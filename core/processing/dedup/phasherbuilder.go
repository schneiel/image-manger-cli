package dedup

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// PHasherConfig holds the configuration for building a DefaultPHasher.
type PHasherConfig struct {
	numWorkers int
	logger     log.Logger
	filesystem filesystem.FileSystem
}

// PHasherOption is a functional option for configuring DefaultPHasher.
type PHasherOption func(*PHasherConfig) error

// WithWorkers sets the number of workers for the hasher.
func WithWorkers(count int) PHasherOption {
	return func(config *PHasherConfig) error {
		if count <= 0 {
			return errors.New("worker count must be greater than 0")
		}
		config.numWorkers = count
		return nil
	}
}

// WithLogger sets the logger for the hasher.
func WithLogger(logger log.Logger) PHasherOption {
	return func(config *PHasherConfig) error {
		if logger == nil {
			return errors.New("logger cannot be nil")
		}
		config.logger = logger
		return nil
	}
}

// WithFilesystem sets the filesystem for the hasher.
func WithFilesystem(fs filesystem.FileSystem) PHasherOption {
	return func(config *PHasherConfig) error {
		if fs == nil {
			return errors.New("filesystem cannot be nil")
		}
		config.filesystem = fs
		return nil
	}
}

// NewDefaultPHasherWithOptions creates a new DefaultPHasher using functional options.
func NewDefaultPHasherWithOptions(localizer i18n.Localizer, options ...PHasherOption) (Hasher, error) {
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}

	// Set default configuration
	config := &PHasherConfig{
		numWorkers: 4, // Default worker count
	}

	// Apply all options
	for _, option := range options {
		err := option(config)
		if err != nil {
			return nil, err
		}
	}

	// Validate required dependencies are provided
	if config.logger == nil {
		return nil, errors.New("logger is required")
	}
	if config.filesystem == nil {
		return nil, errors.New("filesystem is required")
	}

	return &DefaultPHasher{
		numWorkers: config.numWorkers,
		logger:     config.logger,
		fs:         config.filesystem,
		localizer:  localizer,
	}, nil
}

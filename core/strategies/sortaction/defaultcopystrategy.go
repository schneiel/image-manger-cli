// Package sortaction provides action strategies for image sorting operations.
package sortaction

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	shared "github.com/schneiel/ImageManagerGo/core/strategies/shared"
)

// DefaultCopyStrategy implements the action strategy for copying files.
type DefaultCopyStrategy struct {
	logger     log.Logger
	localizer  i18n.Localizer
	fileSystem filesystem.FileSystem
}

// NewDefaultCopyStrategy creates a new DefaultCopyStrategy with injected dependencies.
func NewDefaultCopyStrategy(logger log.Logger, localizer i18n.Localizer) (*DefaultCopyStrategy, error) {
	return NewDefaultCopyStrategyWithFilesystem(logger, localizer, &filesystem.DefaultFileSystem{})
}

// NewDefaultCopyStrategyWithFilesystem creates a new DefaultCopyStrategy with injected filesystem.
func NewDefaultCopyStrategyWithFilesystem(
	logger log.Logger,
	localizer i18n.Localizer,
	fileSystem filesystem.FileSystem,
) (*DefaultCopyStrategy, error) {
	if fileSystem == nil {
		return nil, errors.New("fileSystem cannot be nil")
	}
	return &DefaultCopyStrategy{
		logger:     logger,
		localizer:  localizer,
		fileSystem: fileSystem,
	}, nil
}

// Execute copies the file from source to destination.
func (s *DefaultCopyStrategy) Execute(source, destination string) error {
	s.logger.Infof(
		s.localizer.Translate("FileCopied", map[string]interface{}{"Source": source, "Destination": destination}),
	)
	if err := s.fileSystem.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	srcFile, err := s.fileSystem.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := s.fileSystem.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() { _ = dstFile.Close() }()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

// GetResources returns the resource manager for this strategy.
// For copy strategy, no resources are needed, so it returns nil.
func (s *DefaultCopyStrategy) GetResources() shared.ActionResource {
	return nil
}

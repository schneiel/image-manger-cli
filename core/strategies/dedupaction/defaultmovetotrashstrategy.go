package dedupaction

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	shared "github.com/schneiel/ImageManagerGo/core/strategies/shared"
)

// DefaultMoveToTrashStrategy implements Strategy for moving duplicate files to trash.
type DefaultMoveToTrashStrategy struct {
	trashResource *TrashDirectoryResource
	logger        log.Logger
	localizer     i18n.Localizer
	fileSystem    filesystem.FileSystem
}

// NewDefaultMoveToTrashStrategy creates a new DefaultMoveToTrashStrategy instance.
func NewDefaultMoveToTrashStrategy(config shared.ActionConfig) (*DefaultMoveToTrashStrategy, error) {
	return NewDefaultMoveToTrashStrategyWithFilesystem(config, &filesystem.DefaultFileSystem{})
}

// NewDefaultMoveToTrashStrategyWithFilesystem creates a new DefaultMoveToTrashStrategy with injected filesystem.
func NewDefaultMoveToTrashStrategyWithFilesystem(
	config shared.ActionConfig,
	fileSystem filesystem.FileSystem,
) (*DefaultMoveToTrashStrategy, error) {
	if fileSystem == nil {
		return nil, errors.New("fileSystem cannot be nil")
	}

	trashResource, err := NewTrashDirectoryResourceWithFilesystem(
		config.TrashPath,
		config.Logger,
		config.Localizer,
		fileSystem,
	)
	if err != nil {
		return nil, err
	}
	return &DefaultMoveToTrashStrategy{
		trashResource: trashResource,
		logger:        config.Logger,
		localizer:     config.Localizer,
		fileSystem:    fileSystem,
	}, nil
}

// Execute performs the move to trash action on duplicate files.
func (m *DefaultMoveToTrashStrategy) Execute(duplicate *image.Image, _ *image.Image) error {
	if duplicate == nil {
		if m.logger != nil {
			m.logger.Errorf("Cannot execute move to trash strategy: duplicate image is nil")
		}
		return errors.New("duplicate image cannot be nil")
	}

	path := duplicate.FilePath
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(path))
	newPath := filepath.Join(m.trashResource.trashDir, uniqueName)

	if m.logger != nil && m.localizer != nil {
		m.logger.Infof(m.localizer.Translate("MovingFile", map[string]interface{}{"From": path, "To": newPath}))
	}

	err := m.fileSystem.Rename(path, newPath)
	if err != nil {
		if m.logger != nil && m.localizer != nil {
			m.logger.Errorf(
				m.localizer.Translate("MovingFileError", map[string]interface{}{"FilePath": path, "Error": err}),
			)
		}
		return err
	}

	return nil
}

// GetResources returns the trash directory resource.
func (m *DefaultMoveToTrashStrategy) GetResources() shared.ActionResource {
	return m.trashResource
}

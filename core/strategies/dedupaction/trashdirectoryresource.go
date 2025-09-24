package dedupaction

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// TrashDirectoryResource manages the trash directory.
type TrashDirectoryResource struct {
	trashDir   string
	logger     log.Logger
	localizer  i18n.Localizer
	fileSystem filesystem.FileSystem
}

// NewTrashDirectoryResource creates a new TrashDirectoryResource.
func NewTrashDirectoryResource(trashDir string, logger log.Logger, localizer i18n.Localizer) (*TrashDirectoryResource, error) {
	return NewTrashDirectoryResourceWithFilesystem(trashDir, logger, localizer, &filesystem.DefaultFileSystem{})
}

// NewTrashDirectoryResourceWithFilesystem creates a new TrashDirectoryResource with injected filesystem.
func NewTrashDirectoryResourceWithFilesystem(
	trashDir string,
	logger log.Logger,
	localizer i18n.Localizer,
	fileSystem filesystem.FileSystem,
) (*TrashDirectoryResource, error) {
	if fileSystem == nil {
		return nil, errors.New("fileSystem cannot be nil")
	}
	return &TrashDirectoryResource{
		trashDir:   trashDir,
		logger:     logger,
		localizer:  localizer,
		fileSystem: fileSystem,
	}, nil
}

// Setup creates the trash directory if it doesn't exist.
func (r *TrashDirectoryResource) Setup() error {
	r.logger.Infof(r.localizer.Translate("MoveToTrashSetup", map[string]interface{}{"Dir": r.trashDir}))
	// #nosec G301 -- 0o755 is appropriate for trash directory, needs to be accessible
	return r.fileSystem.MkdirAll(r.trashDir, 0o755)
}

// Teardown performs cleanup after execution.
func (r *TrashDirectoryResource) Teardown() error {
	return nil
}

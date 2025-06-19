package action_strategy

import (
	"ImageManager/i18n"
	"ImageManager/log"
	"ImageManager/util"
	"os"
)

// CopyStrategy defines a strategy for moving/copying files to the destination.
type CopyStrategy struct{}

// NewMoveStrategy creates a new CopyStrategy.
func NewCopyStrategy() *CopyStrategy {
	return &CopyStrategy{}
}

// Setup performs no special setup.
func (s *CopyStrategy) Setup() error {
	return nil // Nothing to set up
}

// Execute copies the file from source to destination.
// It logs a warning if the destination file already exists.
func (s *CopyStrategy) Execute(sourcePath, destinationPath string) error {
	// Skip if the file already exists.
	if _, err := os.Stat(destinationPath); err == nil {
		log.LogWarn(i18n.T("SortSkippingExisting", map[string]interface{}{"Path": destinationPath}))
		return nil
	}

	// Copy the file.
	if err := util.CopyFile(sourcePath, destinationPath); err != nil {
		log.LogError(i18n.T("ErrorCopyingFile", map[string]interface{}{"Source": sourcePath, "Destination": destinationPath, "Error": err}))
		return err
	}
	return nil
}

// Teardown performs no special cleanup.
func (s *CopyStrategy) Teardown() {
	// Nothing to tear down
}

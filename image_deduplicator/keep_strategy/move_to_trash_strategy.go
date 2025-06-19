package keep_strategy

import (
	"ImageManager/i18n"
	"ImageManager/log"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MoveToTrashStrategy defines a strategy for moving duplicate files to a specified trash directory.
type MoveToTrashStrategy struct {
	TrashDir string
}

// Setup ensures the trash directory exists.
func (a *MoveToTrashStrategy) Setup() error {
	log.LogInfo(i18n.T("MoveToTrashSetup", map[string]interface{}{"Dir": a.TrashDir}))
	return os.MkdirAll(a.TrashDir, 0755)
}

// Execute moves the files marked for removal to the trash directory.
// A timestamp is prepended to the filename to avoid collisions.
func (a *MoveToTrashStrategy) Execute(toKeep string, toRemove []string) {
	for _, path := range toRemove {
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(path))
		newPath := filepath.Join(a.TrashDir, uniqueName)
		log.LogInfo(i18n.T("MovingFile", map[string]interface{}{"From": path, "To": newPath}))
		if err := os.Rename(path, newPath); err != nil {
			log.LogError(i18n.T("MovingFileError", map[string]interface{}{"Path": path, "Error": err}))
		}
	}
}

// Teardown performs no action as file handles are closed within Execute.
func (a *MoveToTrashStrategy) Teardown() {}

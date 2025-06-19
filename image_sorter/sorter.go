package imagesorter

import (
	"ImageManager/i18n"
	"ImageManager/image"
	"ImageManager/image_sorter/action_strategy"
	"ImageManager/log"
	"fmt"
	"os"
	"path/filepath"
)

// Sorter is responsible for sorting images into the destination directory structure.
type Sorter struct {
	root     string
	strategy action_strategy.SortActionStrategy
}

// NewSorter creates a new Sorter with a given destination root and an action strategy.
func NewSorter(root string, strategy action_strategy.SortActionStrategy) *Sorter {
	return &Sorter{
		root:     root,
		strategy: strategy,
	}
}

// SortImages organizes images into a folder structure based on their date (YYYY/MM/DD)
// and uses the configured strategy to either copy them or log the action.
func (s *Sorter) SortImages(images []image.Image) error {
	if err := s.strategy.Setup(); err != nil {
		return fmt.Errorf("error setting up sort strategy: %w", err)
	}
	defer s.strategy.Teardown()

	dirMap := make(map[string][]image.Image)

	for _, img := range images {
		if img.Date.IsZero() {
			log.LogError(fmt.Sprintf("Image %s has an invalid (zero) date and will be skipped.", img.FilePath))
			continue
		}

		// Create the destination path based on the date.
		year := fmt.Sprintf("%04d", img.Date.Year())
		month := fmt.Sprintf("%02d", img.Date.Month())
		day := fmt.Sprintf("%02d", img.Date.Day())
		dirPath := filepath.Join(s.root, year, month, day)

		dirMap[dirPath] = append(dirMap[dirPath], img)
	}

	for dirPath, imgsInDir := range dirMap {
		// Create the destination directory if it doesn't exist and we're not in dry-run mode.
		// For dry-run, we don't create directories.
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if _, isDryRun := s.strategy.(*action_strategy.DryRunStrategy); !isDryRun {
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					log.LogError(i18n.T("ErrorCreatingDir", map[string]interface{}{"Path": dirPath, "Error": err}))
					continue
				}
			}
		}

		for _, img := range imgsInDir {
			newPath := filepath.Join(dirPath, img.OriginalFileName)
			if err := s.strategy.Execute(img.FilePath, newPath); err != nil {
				// The strategy itself will log detailed errors, here we can just note that it failed.
				log.LogError(fmt.Sprintf("Failed to process %s", img.FilePath))
			}
		}
	}
	return nil
}

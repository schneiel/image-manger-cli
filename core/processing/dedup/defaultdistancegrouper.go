// Package dedup provides functionality for detecting and handling duplicate images.
package dedup

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// DefaultDistanceGrouper implements the Grouper interface by comparing hash distances.
type DefaultDistanceGrouper struct {
	threshold int
	logger    log.Logger
	localizer i18n.Localizer
}

// NewDistanceGrouper creates a grouper that finds duplicates within a given distance threshold.
func NewDistanceGrouper(threshold int, logger log.Logger, localizer i18n.Localizer) (Grouper, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}

	return &DefaultDistanceGrouper{
		threshold: threshold,
		logger:    logger,
		localizer: localizer,
	}, nil
}

// Group organizes image hashes into duplicate groups based on hash distance comparison.
// It compares each hash with all others and groups files whose hash distance is within the threshold.
func (g *DefaultDistanceGrouper) Group(hashes []*image.HashInfo) ([]DuplicateGroup, error) {
	g.logger.Info(g.localizer.Translate("GroupingDuplicatesStarted"))
	var duplicateGroups []DuplicateGroup
	processedFiles := make(map[string]bool)

	for i := 0; i < len(hashes); i++ {
		if processedFiles[hashes[i].FilePath] {
			continue
		}

		currentGroup := DuplicateGroup{hashes[i].FilePath}
		for j := i + 1; j < len(hashes); j++ {
			if processedFiles[hashes[j].FilePath] {
				continue
			}

			distance, err := hashes[i].Hash.Distance(hashes[j].Hash)
			if err != nil {
				// Log error but continue, as one failed comparison shouldn't stop the whole process.
				g.logger.Warnf(
					g.localizer.Translate(
						"HashCompareError",
						map[string]interface{}{"File1": hashes[i].FilePath, "File2": hashes[j].FilePath, "Error": err},
					),
				)
				continue
			}

			if distance <= g.threshold {
				currentGroup = append(currentGroup, hashes[j].FilePath)
			}
		}

		if len(currentGroup) > 1 {
			duplicateGroups = append(duplicateGroups, currentGroup)
			// Mark all members of the found group as processed
			for _, path := range currentGroup {
				processedFiles[path] = true
			}
		}
	}

	g.logger.Infof(
		g.localizer.Translate("GroupingDuplicatesFinished", map[string]interface{}{"Count": len(duplicateGroups)}),
	)
	return duplicateGroups, nil
}

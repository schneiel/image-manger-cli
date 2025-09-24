package dedup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// DefaultSizeScanner implements the Scanner interface by grouping files by size.
type DefaultSizeScanner struct {
	AllowedExtensions []string
	Logger            log.Logger
	localizer         i18n.Localizer
}

// NewDefaultSizeScanner creates a new scanner that filters by extension and groups by size.
func NewDefaultSizeScanner(allowedExtensions []string, logger log.Logger, localizer i18n.Localizer) (Scanner, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}

	return &DefaultSizeScanner{
		AllowedExtensions: allowedExtensions,
		Logger:            logger,
		localizer:         localizer,
	}, nil
}

// Scan finds potential duplicate files by grouping them by file size.
func (s *DefaultSizeScanner) Scan(rootPath string) (FileGroup, error) {
	s.Logger.Info(s.localizer.Translate("ScanningForFiles"))
	sizes := make(map[int64][]string)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			s.Logger.Warnf(
				s.localizer.Translate("ErrorAccessingPath", map[string]interface{}{"FilePath": path, "Error": err}),
			)
			return nil
		}

		if info.IsDir() || info.Size() < 1 {
			return nil
		}

		if !slices.Contains(s.AllowedExtensions, strings.ToLower(filepath.Ext(path))) {
			return nil
		}

		sizes[info.Size()] = append(sizes[info.Size()], path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", rootPath, err)
	}

	// Filter out sizes with only one file
	var groups FileGroup
	for _, paths := range sizes {
		if len(paths) >= 2 {
			groups = append(groups, paths)
		}
	}
	s.Logger.Infof(s.localizer.Translate("PotentialDuplicateGroupsFound", map[string]interface{}{"Count": len(groups)}))
	return groups, nil
}

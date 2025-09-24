// Package sort provides functionality for analyzing and processing images.
package sort

import (
	"errors"
	"path/filepath"

	"github.com/schneiel/ImageManagerGo/core/date"
	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/exif"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// DefaultImageAnalyzer is responsible for analyzing a single image file to extract its metadata,
// primarily its creation date.
type DefaultImageAnalyzer struct {
	dateProcessor date.DateProcessor
	exifReader    exif.Reader
	logger        log.Logger
	localizer     i18n.Localizer
}

// NewDefaultImageAnalyzer creates a new image analyzer with injected dependencies for date processing and EXIF reading.
func NewDefaultImageAnalyzer(
	dateProcessor date.DateProcessor,
	exifReader exif.Reader,
	logger log.Logger,
	localizer i18n.Localizer,
) (ImageAnalyzer, error) {
	if dateProcessor == nil {
		return nil, errors.New("dateProcessor cannot be nil")
	}
	if exifReader == nil {
		return nil, errors.New("exifReader cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}

	return &DefaultImageAnalyzer{
		dateProcessor: dateProcessor,
		exifReader:    exifReader,
		logger:        logger,
		localizer:     localizer,
	}, nil
}

// Analyze takes a file path, reads its EXIF data, determines its date,
// and returns a populated image.Image struct.
func (a *DefaultImageAnalyzer) Analyze(filePath string) (image.Image, error) {
	exifFieldsMap, err := a.exifReader.ReadExif(filePath)
	if err != nil {
		a.logger.Warnf(
			a.localizer.Translate("ExifDecodeWarning", map[string]interface{}{"FilePath": filePath, "Error": err}),
		)
	}

	imgDate, determinationErr := a.dateProcessor.GetBestAvailableDate(exifFieldsMap, filePath)
	if determinationErr != nil {
		return image.Image{}, errors.New(
			a.localizer.Translate(
				"ErrorGettingDate",
				map[string]interface{}{"FilePath": filePath, "Error": determinationErr},
			),
		)
	}

	return image.Image{
		FilePath:         filePath,
		OriginalFileName: filepath.Base(filePath),
		Date:             imgDate,
	}, nil
}

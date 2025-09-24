// Package exif provides functionality to read EXIF metadata from image files.
package exif

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultReader is responsible for extracting EXIF metadata from a file.
type DefaultReader struct {
	fileSystem  filesystem.FileSystem
	exifDecoder Decoder
	localizer   i18n.Localizer
}

// NewDefaultReader creates a new EXIF reader with the given dependencies.
func NewDefaultReader(localizer i18n.Localizer, fileSystem filesystem.FileSystem, exifDecoder Decoder) (*DefaultReader, error) {
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	if fileSystem == nil {
		return nil, errors.New("fileSystem cannot be nil")
	}
	if exifDecoder == nil {
		return nil, errors.New("exifDecoder cannot be nil")
	}

	return &DefaultReader{
		fileSystem:  fileSystem,
		exifDecoder: exifDecoder,
		localizer:   localizer,
	}, nil
}

// validatePath validates a file path to prevent path traversal attacks.
func (r *DefaultReader) validatePath(path string) error {
	// Clean the path to resolve any '..' or '.' components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return errors.New("path traversal detected: path contains '..'")
	}

	return nil
}

// ReadExif extracts relevant date-related EXIF fields from the given file.
func (r *DefaultReader) ReadExif(filePath string) (map[string]interface{}, error) {
	if err := r.validatePath(filePath); err != nil {
		return nil, err
	}
	// #nosec G304 -- filePath is validated by validatePath to prevent traversal attacks
	file, err := r.fileSystem.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(r.localizer.Translate("ErrorOpeningFile", map[string]interface{}{"Error": err}), err)
	}
	defer func() {
		_ = file.Close()
	}()

	x, err := r.exifDecoder.Decode(file)
	if err != nil {
		// This is not a fatal error for the map, just means it will be empty.
		// The caller can decide how to handle a decoding error.
		return nil, errors.New(r.localizer.Translate("ExifDecodeError", map[string]interface{}{
			"FilePath": filePath,
			"Error":    err,
		}))
	}

	return r.extractDateFields(x), nil
}

// extractDateFields pulls specific date/time tags from the decoded EXIF data.
func (r *DefaultReader) extractDateFields(x *exif.Exif) map[string]interface{} {
	exifFieldsMap := make(map[string]interface{})

	// Handle nil input
	if x == nil {
		return exifFieldsMap
	}

	if tag, err := x.Get(exif.DateTimeOriginal); err == nil {
		strVal, _ := tag.StringVal()
		exifFieldsMap["DateTimeOriginal"] = strVal
	}

	if tag, err := x.Get(exif.DateTimeDigitized); err == nil {
		strVal, _ := tag.StringVal()
		exifFieldsMap["DateTimeDigitized"] = strVal
	}

	if tag, err := x.Get(exif.DateTime); err == nil {
		strVal, _ := tag.StringVal()
		exifFieldsMap["DateTime"] = strVal
	}

	// Handle sub-second data
	dtOrigTag, _ := x.Get(exif.DateTimeOriginal)
	subSecTag, _ := x.Get(exif.SubSecTimeOriginal)
	if dtOrigTag != nil && subSecTag != nil {
		dtOrigStr, _ := dtOrigTag.StringVal()
		subSecStr, _ := subSecTag.StringVal()
		if dtOrigStr != "" && subSecStr != "" {
			// Combine date/time with sub-seconds
			combined := fmt.Sprintf("%s.%s", dtOrigStr, strings.TrimSpace(subSecStr))
			exifFieldsMap["SubSecDateTimeOriginal"] = combined
		}
	}

	return exifFieldsMap
}

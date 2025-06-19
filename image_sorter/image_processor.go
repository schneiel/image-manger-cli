// Package imagesorter contains the logic for processing and sorting images.
package imagesorter

import (
	"ImageManager/config"
	"ImageManager/date"
	"ImageManager/i18n"
	"ImageManager/image"
	"ImageManager/log"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
)

// ImageProcessor is responsible for reading and processing image files.
type ImageProcessor struct {
	datePriorityProcessor *date.DatePriorityProcessor
	AllowedExtensions     map[string]struct{}
	fileCache             sync.Map
}

// NewImageProcessor creates a new ImageProcessor with the given configuration.
func NewImageProcessor(cfg config.Config) (*ImageProcessor, error) {
	extMap := make(map[string]struct{})
	for _, ext := range cfg.AllowedImageExtensions {
		extMap[ext] = struct{}{}
	}

	// Use the strategies from the config to create the date processor
	dateProcessor, err := date.NewDatePriorityProcessorFromStrategies(cfg.Date.Strategies)
	if err != nil {
		return nil, err
	}

	return &ImageProcessor{
		datePriorityProcessor: dateProcessor,
		AllowedExtensions:     extMap,
	}, nil
}

// HandleImages walks a directory, finds images, and extracts their metadata concurrently.
func (ip *ImageProcessor) HandleImages(dirPath string) []image.Image {
	var wg sync.WaitGroup
	processedImagesChan := make(chan image.Image, 100)

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.LogError(i18n.T("ErrorAccessingPath", map[string]interface{}{"Path": path, "Error": err}))
			return err
		}

		if !d.IsDir() && ip.isImage(path) {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				img, procErr := ip.handleSingleImageFile(filePath)
				if procErr != nil {
					log.LogError(i18n.T("ErrorProcessingImage", map[string]interface{}{"Path": filePath, "Error": procErr}))
					return
				}
				processedImagesChan <- img
			}(path)
		}
		return nil
	}

	if err := filepath.WalkDir(dirPath, walkFunc); err != nil {
		log.LogError(i18n.T("ErrorWalkingDir", map[string]interface{}{"Path": dirPath, "Error": err}))
	}

	go func() {
		wg.Wait()
		close(processedImagesChan)
	}()

	var result []image.Image
	for img := range processedImagesChan {
		result = append(result, img)
	}
	return result
}

// handleSingleImageFile processes a single image file to extract its date.
func (ip *ImageProcessor) handleSingleImageFile(filePath string) (image.Image, error) {
	exifFieldsMap := make(map[string]interface{})

	f, err := os.Open(filePath)
	if err != nil {
		return image.Image{}, err
	}
	defer f.Close()

	x, exifErr := exif.Decode(f)
	if exifErr != nil {
		log.LogWarn(fmt.Sprintf("Error decoding EXIF for %s: %v. Trying fallback date.", filePath, exifErr))
	}

	if x != nil {
		if tag, errGet := x.Get(exif.DateTimeOriginal); errGet == nil {
			strVal, _ := tag.StringVal()
			exifFieldsMap[string(date.DateTimeOriginal)] = strVal
		}
		if tag, errGet := x.Get(exif.DateTimeDigitized); errGet == nil {
			strVal, _ := tag.StringVal()
			exifFieldsMap[string(date.DateTimeDigitized)] = strVal
		}
		dtOrigTag, _ := x.Get(exif.DateTimeOriginal)
		subSecTag, _ := x.Get(exif.SubSecTimeOriginal)
		if dtOrigTag != nil && subSecTag != nil {
			dtOrigStr, _ := dtOrigTag.StringVal()
			subSecStr, _ := subSecTag.StringVal()
			if dtOrigStr != "" && subSecStr != "" {
				exifFieldsMap[string(date.SubSecDateTimeOriginal)] = fmt.Sprintf("%s.%s", dtOrigStr, strings.TrimSpace(subSecStr))
			}
		}
		if tag, errGet := x.Get(exif.DateTime); errGet == nil {
			strVal, _ := tag.StringVal()
			exifFieldsMap[string(date.DateTime)] = strVal
		}
	}

	// Use the DatePriorityProcessor to determine the best date.
	imgDate, determinationErr := ip.datePriorityProcessor.GetBestAvailableDate(exifFieldsMap, filePath)
	if determinationErr != nil {
		return image.Image{}, fmt.Errorf(i18n.T("ErrorGettingDate", map[string]interface{}{"Path": filePath, "Error": determinationErr}))
	}

	return image.Image{
		FilePath:         filePath,
		OriginalFileName: filepath.Base(filePath),
		Date:             imgDate,
	}, nil
}

// isImage checks if a file has a supported image extension.
func (ip *ImageProcessor) isImage(path string) bool {
	if val, ok := ip.fileCache.Load(path); ok {
		return val.(bool)
	}
	ext := strings.ToLower(filepath.Ext(path))
	_, exists := ip.AllowedExtensions[ext]
	ip.fileCache.Store(path, exists)
	return exists
}

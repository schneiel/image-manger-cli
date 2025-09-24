// Package shared provides common action strategy patterns for image operations.
package shared

import (
	"encoding/csv"
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// DefaultCSVResource implements CSV resource management.
type DefaultCSVResource struct {
	filePath   string
	header     []string
	logger     log.Logger
	localizer  i18n.Localizer
	fileSystem filesystem.FileSystem
	file       filesystem.File
	writer     *csv.Writer
}

// NewCSVResource creates a CSV resource manager.
func NewCSVResource(csvFilePath string, header []string, logger log.Logger, localizer i18n.Localizer) CSVResource {
	return NewCSVResourceWithFilesystem(csvFilePath, header, logger, localizer, &filesystem.DefaultFileSystem{})
}

// NewCSVResourceWithFilesystem creates a CSV resource manager with injected filesystem.
func NewCSVResourceWithFilesystem(
	csvFilePath string,
	header []string,
	logger log.Logger,
	localizer i18n.Localizer,
	fileSystem filesystem.FileSystem,
) CSVResource {
	if fileSystem == nil {
		panic("fileSystem cannot be nil")
	}
	return &DefaultCSVResource{
		filePath:   csvFilePath,
		header:     header,
		logger:     logger,
		localizer:  localizer,
		fileSystem: fileSystem,
	}
}

// Setup initializes the CSV file and writer.
func (r *DefaultCSVResource) Setup() error {
	file, err := r.fileSystem.Create(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file %s: %w", r.filePath, err)
	}

	r.file = file
	r.writer = csv.NewWriter(file)

	// Write header if provided
	if len(r.header) > 0 {
		err := r.writer.Write(r.header)
		if err != nil {
			_ = r.file.Close()
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
		r.writer.Flush()
	}

	return nil
}

// WriteRow writes a row to the CSV file.
func (r *DefaultCSVResource) WriteRow(row []string) error {
	if r.writer == nil {
		return errors.New("CSV writer not initialized, call Setup() first")
	}

	err := r.writer.Write(row)
	if err != nil {
		return fmt.Errorf("failed to write CSV row: %w", err)
	}

	r.writer.Flush()
	return r.writer.Error()
}

// Teardown closes the CSV file and writer.
func (r *DefaultCSVResource) Teardown() error {
	if r.writer != nil {
		r.writer.Flush()
	}
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

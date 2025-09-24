package config

import (
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
)

// DefaultFileReader implements FileReader using injected filesystem.
type DefaultFileReader struct {
	fileSystem filesystem.FileSystem
}

// NewDefaultFileReader creates a new DefaultFileReader.
func NewDefaultFileReader() FileReader {
	return NewDefaultFileReaderWithFilesystem(&filesystem.DefaultFileSystem{})
}

// NewDefaultFileReaderWithFilesystem creates a new DefaultFileReader with injected filesystem.
func NewDefaultFileReaderWithFilesystem(fileSystem filesystem.FileSystem) FileReader {
	if fileSystem == nil {
		panic("fileSystem cannot be nil")
	}
	return &DefaultFileReader{fileSystem: fileSystem}
}

// ReadFile reads a file and returns its contents.
func (r *DefaultFileReader) ReadFile(filename string) ([]byte, error) {
	data, err := r.fileSystem.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return data, nil
}

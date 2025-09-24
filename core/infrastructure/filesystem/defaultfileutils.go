package filesystem

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultFileUtils provides file system operations with dependency injection.
type DefaultFileUtils struct {
	fs        FileSystem
	localizer i18n.Localizer
}

// NewFileUtils creates a new FileUtils instance with the given file system and localizer.
func NewFileUtils(fs FileSystem, localizer i18n.Localizer) (FileUtils, error) {
	if fs == nil {
		return nil, errors.New("FileUtils requires non-nil filesystem.FileSystem")
	}
	if localizer == nil {
		return nil, errors.New("FileUtils requires non-nil i18n.Localizer")
	}
	return &DefaultFileUtils{fs: fs, localizer: localizer}, nil
}

// GetFileSystem returns the underlying FileSystem instance.
func (fu *DefaultFileUtils) GetFileSystem() FileSystem {
	return fu.fs
}

// validatePath validates a file path to prevent path traversal attacks.
func (fu *DefaultFileUtils) validatePath(path string) error {
	// Clean the path to resolve any '..' or '.' components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return errors.New("path traversal detected: path contains '..'")
	}

	return nil
}

// CopyFile copies a file from sourcePath to destinationPath.
// It first attempts to create a hard link to speed up the process.
// If that fails (e.g., across different partitions), it falls back to a standard copy.
func (fu *DefaultFileUtils) CopyFile(sourcePath, destinationPath string) error {
	// Validate paths to prevent traversal attacks
	if err := fu.validatePath(sourcePath); err != nil {
		return err
	}
	if err := fu.validatePath(destinationPath); err != nil {
		return err
	}

	// no-op when copying file to itself
	if sourcePath == destinationPath {
		return nil
	}
	err := fu.fs.Link(sourcePath, destinationPath)
	if err == nil {
		return nil
	}
	// #nosec G304 -- sourcePath is validated by validatePath to prevent traversal attacks
	src, err := fu.fs.Open(sourcePath)
	if err != nil {
		return errors.New(
			fu.localizer.Translate("SourceFileOpenError", map[string]interface{}{"FilePath": sourcePath, "Error": err}),
		)
	}
	defer func() {
		_ = src.Close()
	}()

	// #nosec G304 -- destinationPath is validated by validatePath to prevent traversal attacks
	dst, err := fu.fs.Create(destinationPath)
	if err != nil {
		return errors.New(
			fu.localizer.Translate(
				"DestFileCreateError",
				map[string]interface{}{"FilePath": destinationPath, "Error": err},
			),
		)
	}
	defer func() {
		_ = dst.Close()
	}()

	if _, err = io.Copy(dst, src); err != nil {
		return errors.New(
			fu.localizer.Translate(
				"FileCopyContentError",
				map[string]interface{}{"Source": sourcePath, "Destination": destinationPath, "Error": err},
			),
		)
	}

	si, err := fu.fs.Stat(sourcePath)
	if err != nil {
		return errors.New(
			fu.localizer.Translate("FileInfoGetError", map[string]interface{}{"FilePath": sourcePath, "Error": err}),
		)
	}

	err = fu.fs.Chmod(destinationPath, si.Mode())
	if err != nil {
		return errors.New(
			fu.localizer.Translate(
				"FilePermissionsSetError",
				map[string]interface{}{"FilePath": destinationPath, "Error": err},
			),
		)
	}

	return nil
}

// Exists returns true if the file or directory at the given path exists.
func (fu *DefaultFileUtils) Exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := fu.fs.Stat(path)
	return err == nil
}

// EnsureDir creates a directory (and parents) if it does not exist.
// Returns an error if the path exists and is not a directory, or if creation fails.
func (fu *DefaultFileUtils) EnsureDir(path string) error {
	if path == "" {
		return errors.New("path is empty")
	}

	info, err := fu.fs.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists and is not a directory: %s", path)
		}
		return nil
	}

	if !fu.fs.IsNotExist(err) {
		return fmt.Errorf("failed to stat directory %s: %w", path, err)
	}

	if err := fu.fs.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

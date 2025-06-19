// Package util provides helper functions for file system operations.
package util

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies a file from sourcePath to destinationPath.
// It first attempts to create a hard link to speed up the process.
// If that fails (e.g., across different partitions), it falls back to a standard copy.
func CopyFile(sourcePath, destinationPath string) error {
	// Attempt to create a hard link. This is much faster than a copy.
	err := os.Link(sourcePath, destinationPath)
	if err == nil {
		return nil // Success
	}

	// Fallback: Manual copy
	src, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("could not open source file %s: %w", sourcePath, err)
	}
	defer src.Close()

	dst, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("could not create destination file %s: %w", destinationPath, err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("could not copy file content from %s to %s: %w", sourcePath, destinationPath, err)
	}

	// Copy file permissions.
	si, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("could not get file info for %s: %w", sourcePath, err)
	}
	err = os.Chmod(destinationPath, si.Mode())
	if err != nil {
		return fmt.Errorf("could not set permissions for %s: %w", destinationPath, err)
	}
	return nil
}

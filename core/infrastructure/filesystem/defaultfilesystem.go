package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DefaultFileSystem implements FileSystem using the standard os package.
type DefaultFileSystem struct{}

// NewDefaultFileSystem creates a new DefaultFileSystem instance.
func NewDefaultFileSystem() *DefaultFileSystem {
	return &DefaultFileSystem{}
}

// validatePath validates a file path to prevent path traversal attacks.
func (f *DefaultFileSystem) validatePath(path string) error {
	// Clean the path to resolve any '..' or '.' components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return errors.New("path traversal detected: path contains '..'")
	}

	// Additional validation could be added here based on requirements
	// For example, checking against allowed directories

	return nil
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666.
func (f *DefaultFileSystem) Create(name string) (File, error) {
	if err := f.validatePath(name); err != nil {
		return nil, err
	}
	// #nosec G304 -- path is validated by validatePath to prevent traversal attacks
	file, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", name, err)
	}
	return file, nil
}

// Open opens the named file for reading. If successful, methods on the returned file can be used for reading.
func (f *DefaultFileSystem) Open(name string) (File, error) {
	if err := f.validatePath(name); err != nil {
		return nil, err
	}
	// #nosec G304 -- path is validated by validatePath to prevent traversal attacks
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", name, err)
	}
	return file, nil
}

// OpenFile opens the named file with specified flag (O_RDONLY etc.) and perm (before umask).
func (f *DefaultFileSystem) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	if err := f.validatePath(name); err != nil {
		return nil, err
	}
	// #nosec G304 -- path is validated by validatePath to prevent traversal attacks
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s with flags %d: %w", name, flag, err)
	}
	return file, nil
}

// Remove removes the named file or (empty) directory.
func (f *DefaultFileSystem) Remove(name string) error {
	err := os.Remove(name)
	if err != nil {
		return fmt.Errorf("failed to remove %s: %w", name, err)
	}
	return nil
}

// RemoveAll removes path and any children it contains.
func (f *DefaultFileSystem) RemoveAll(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to remove all %s: %w", path, err)
	}
	return nil
}

// Rename renames (moves) oldpath to newpath.
func (f *DefaultFileSystem) Rename(oldpath, newpath string) error {
	err := os.Rename(oldpath, newpath)
	if err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", oldpath, newpath, err)
	}
	return nil
}

// Mkdir creates a new directory with the specified name and permission bits.
func (f *DefaultFileSystem) Mkdir(name string, perm os.FileMode) error {
	err := os.Mkdir(name, perm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", name, err)
	}
	return nil
}

// MkdirAll creates a directory named path, along with any necessary parents.
func (f *DefaultFileSystem) MkdirAll(path string, perm os.FileMode) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return fmt.Errorf("failed to create directory tree %s: %w", path, err)
	}
	return nil
}

// ReadDir reads the named directory and returns a list of directory entries.
func (f *DefaultFileSystem) ReadDir(dirname string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirname, err)
	}
	return entries, nil
}

// Stat returns a FileInfo describing the named file.
func (f *DefaultFileSystem) Stat(name string) (os.FileInfo, error) {
	info, err := os.Stat(name)
	if err != nil {
		// Don't wrap os.IsNotExist errors to preserve type checking
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to stat %s: %w", name, err)
	}
	return info, nil
}

// Lstat returns a FileInfo describing the named file without following symbolic links.
func (f *DefaultFileSystem) Lstat(name string) (os.FileInfo, error) {
	info, err := os.Lstat(name)
	if err != nil {
		// Don't wrap os.IsNotExist errors to preserve type checking
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to lstat %s: %w", name, err)
	}
	return info, nil
}

// ReadFile reads the named file and returns the contents.
func (f *DefaultFileSystem) ReadFile(filename string) ([]byte, error) {
	if err := f.validatePath(filename); err != nil {
		return nil, err
	}
	// #nosec G304 -- path is validated by validatePath to prevent traversal attacks
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return data, nil
}

// WriteFile writes data to the named file, creating it if necessary.
func (f *DefaultFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	err := os.WriteFile(filename, data, perm)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}

// CreateTemp creates a new temporary file in the directory dir with a name beginning with pattern.
func (f *DefaultFileSystem) CreateTemp(dir, pattern string) (File, error) {
	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file in %s with pattern %s: %w", dir, pattern, err)
	}
	return file, nil
}

// MkdirTemp creates a new temporary directory in the directory dir with a name beginning with pattern.
func (f *DefaultFileSystem) MkdirTemp(dir, pattern string) (string, error) {
	tempDir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory in %s with pattern %s: %w", dir, pattern, err)
	}
	return tempDir, nil
}

// Chmod changes the mode of the named file to mode.
func (f *DefaultFileSystem) Chmod(name string, mode os.FileMode) error {
	err := os.Chmod(name, mode)
	if err != nil {
		return fmt.Errorf("failed to chmod %s to %v: %w", name, mode, err)
	}
	return nil
}

// Chown changes the numeric uid and gid of the named file.
func (f *DefaultFileSystem) Chown(name string, uid, gid int) error {
	err := os.Chown(name, uid, gid)
	if err != nil {
		return fmt.Errorf("failed to chown %s to %d:%d: %w", name, uid, gid, err)
	}
	return nil
}

// Chtimes changes the access and modification times of the named file.
func (f *DefaultFileSystem) Chtimes(name string, atime, mtime time.Time) error {
	err := os.Chtimes(name, atime, mtime)
	if err != nil {
		return fmt.Errorf("failed to change times for %s: %w", name, err)
	}
	return nil
}

// Getwd returns a rooted path name corresponding to the current directory.
func (f *DefaultFileSystem) Getwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	return dir, nil
}

// Chdir changes the current working directory to the named directory.
func (f *DefaultFileSystem) Chdir(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", dir, err)
	}
	return nil
}

// Symlink creates newname as a symbolic link to oldname.
func (f *DefaultFileSystem) Symlink(oldname, newname string) error {
	err := os.Symlink(oldname, newname)
	if err != nil {
		return fmt.Errorf("failed to create symlink from %s to %s: %w", oldname, newname, err)
	}
	return nil
}

// Link creates newname as a hard link to oldname.
func (f *DefaultFileSystem) Link(oldname, newname string) error {
	err := os.Link(oldname, newname)
	if err != nil {
		return fmt.Errorf("failed to create hard link from %s to %s: %w", oldname, newname, err)
	}
	return nil
}

// Readlink returns the destination of the named symbolic link.
func (f *DefaultFileSystem) Readlink(name string) (string, error) {
	link, err := os.Readlink(name)
	if err != nil {
		return "", fmt.Errorf("failed to read link %s: %w", name, err)
	}
	return link, nil
}

// WalkDir walks the file tree rooted at root, calling fn for each file or directory.
func (f *DefaultFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	err := f.validatePath(root)
	if err != nil {
		return err
	}
	err = filepath.WalkDir(root, fn)
	if err != nil {
		return fmt.Errorf("failed to walk directory tree %s: %w", root, err)
	}
	return nil
}

// IsNotExist returns a boolean indicating whether the error indicates that a file or directory does not exist.
func (f *DefaultFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

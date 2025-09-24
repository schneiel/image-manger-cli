// Package filesystem provides abstraction over file system operations for better testability.
package filesystem

import (
	"io/fs"
	"os"
	"time"
)

// FileReader provides read-only file operations.
type FileReader interface {
	Open(name string) (File, error)
	ReadFile(filename string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
	Lstat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
}

// FileWriter provides write operations for files.
type FileWriter interface {
	Create(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Remove(name string) error
	Rename(oldpath, newpath string) error
}

// DirectoryManager provides directory-specific operations.
type DirectoryManager interface {
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	ReadDir(dirname string) ([]os.DirEntry, error)
	RemoveAll(path string) error
	WalkDir(root string, fn fs.WalkDirFunc) error
}

// TempManager provides temporary file and directory operations.
type TempManager interface {
	CreateTemp(dir, pattern string) (File, error)
	MkdirTemp(dir, pattern string) (string, error)
}

// PermissionManager provides permission and ownership operations.
type PermissionManager interface {
	Chmod(name string, mode os.FileMode) error
	Chown(name string, uid, gid int) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
}

// WorkingDirManager provides working directory operations.
type WorkingDirManager interface {
	Getwd() (string, error)
	Chdir(dir string) error
}

// LinkManager provides symbolic and hard link operations.
type LinkManager interface {
	Symlink(oldname, newname string) error
	Link(oldname, newname string) error
	Readlink(name string) (string, error)
}

// FileSystem provides comprehensive filesystem operations by embedding focused interfaces.
// Components should prefer using specific interfaces (FileReader, FileWriter, etc.) over this composite.
type FileSystem interface {
	FileReader
	FileWriter
	DirectoryManager
	TempManager
	PermissionManager
	WorkingDirManager
	LinkManager
}

// FileUtils provides utility functions for file operations.
type FileUtils interface {
	CopyFile(sourcePath, destinationPath string) error
	Exists(path string) bool
	EnsureDir(path string) error
}

// File provides an interface for file operations.
type File interface {
	Close() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Seek(int64, int) (int64, error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	Name() string
	Readdir(count int) ([]os.FileInfo, error)
	Readdirnames(n int) ([]string, error)
}

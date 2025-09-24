package dedupkeep

import (
	"os"
	"time"
)

// DummyFileInfo implements os.FileInfo for error cases when file stats cannot be retrieved.
type DummyFileInfo struct {
	name string
}

// Name returns the base name of the file.
func (d *DummyFileInfo) Name() string { return d.name }

// Size returns the length in bytes for regular files; system-dependent for others.
func (d *DummyFileInfo) Size() int64 { return 0 }

// Mode returns the file mode bits.
func (d *DummyFileInfo) Mode() os.FileMode { return 0 }

// ModTime returns the modification time.
func (d *DummyFileInfo) ModTime() time.Time { return time.Time{} }

// IsDir returns true if file is a directory.
func (d *DummyFileInfo) IsDir() bool { return false }

// Sys returns underlying data source (can return nil).
func (d *DummyFileInfo) Sys() any { return nil }

package config

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
)

func TestDefaultFileReader_ReadFile(t *testing.T) {
	t.Parallel()
	// Arrange
	mockFileSystem := &MockFileSystem{
		readFileData: []byte("test content"),
		readFileErr:  nil,
	}
	reader := NewDefaultFileReaderWithFilesystem(mockFileSystem)

	// Act
	data, err := reader.ReadFile("test.txt")
	// Assert
	if err != nil {
		t.Errorf("ReadFile() unexpected error: %v", err)
	}
	if string(data) != "test content" {
		t.Errorf("ReadFile() = %s, want %s", string(data), "test content")
	}
}

func TestDefaultFileReader_ReadFile_Error(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedErr := errors.New("file not found")
	mockFileSystem := &MockFileSystem{
		readFileData: nil,
		readFileErr:  expectedErr,
	}
	reader := NewDefaultFileReaderWithFilesystem(mockFileSystem)

	// Act
	data, err := reader.ReadFile("nonexistent.txt")

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Errorf("ReadFile() error = %v, want %v", err, expectedErr)
	}
	if data != nil {
		t.Errorf("ReadFile() data = %v, want nil", data)
	}
}

func TestDefaultFileReader_Constructor_NilFilesystem(t *testing.T) {
	t.Parallel()
	// Arrange & Act & Assert
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Constructor should panic with nil filesystem")
		}
	}()

	NewDefaultFileReaderWithFilesystem(nil)
}

// MockFileSystem is a mock implementation for testing.
type MockFileSystem struct {
	readFileData []byte
	readFileErr  error
}

func (m *MockFileSystem) Create(_ string) (filesystem.File, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) Open(_ string) (filesystem.File, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) OpenFile(_ string, _ int, _ os.FileMode) (filesystem.File, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) Remove(_ string) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) RemoveAll(_ string) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Rename(_, _ string) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Mkdir(_ string, _ os.FileMode) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) MkdirAll(_ string, _ os.FileMode) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) ReadDir(_ string) ([]os.DirEntry, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) Stat(_ string) (os.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) Lstat(_ string) (os.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) ReadFile(_ string) ([]byte, error) {
	return m.readFileData, m.readFileErr
}

func (m *MockFileSystem) WriteFile(_ string, _ []byte, _ os.FileMode) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) CreateTemp(_, _ string) (filesystem.File, error) {
	return nil, errors.New("not implemented")
}

func (m *MockFileSystem) MkdirTemp(_, _ string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *MockFileSystem) Chmod(_ string, _ os.FileMode) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Chown(_ string, _, _ int) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Chtimes(_ string, _, _ time.Time) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Getwd() (string, error) {
	return "", errors.New("not implemented")
}

func (m *MockFileSystem) Chdir(_ string) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Symlink(_, _ string) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Link(_, _ string) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) Readlink(_ string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *MockFileSystem) WalkDir(_ string, _ fs.WalkDirFunc) error {
	return errors.New("not implemented")
}

func (m *MockFileSystem) IsNotExist(err error) bool {
	return errors.Is(err, errors.New("not implemented"))
}

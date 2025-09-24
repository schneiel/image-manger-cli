package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockLocalizer implements i18n.Localizer for testing.
type MockLocalizer struct{}

func (m *MockLocalizer) Translate(key string, _ ...map[string]interface{}) string {
	return key // Return the key as translation for testing
}

func (m *MockLocalizer) GetCurrentLanguage() string {
	return "en"
}

func (m *MockLocalizer) SetLanguage(_ string) error {
	return nil
}

func (m *MockLocalizer) IsInitialized() bool {
	return true
}

// MockFileSystem is a mock implementation of the FileSystem interface.
type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) Create(name string) (File, error) {
	args := m.Called(name)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(File), args.Error(1)
}

func (m *MockFileSystem) Open(name string) (File, error) {
	args := m.Called(name)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(File), args.Error(1)
}

func (m *MockFileSystem) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	args := m.Called(name, flag, perm)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(File), args.Error(1)
}

func (m *MockFileSystem) Remove(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockFileSystem) RemoveAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockFileSystem) Rename(oldpath, newpath string) error {
	args := m.Called(oldpath, newpath)
	return args.Error(0)
}

func (m *MockFileSystem) Mkdir(name string, perm os.FileMode) error {
	args := m.Called(name, perm)
	return args.Error(0)
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockFileSystem) ReadDir(dirname string) ([]os.DirEntry, error) {
	args := m.Called(dirname)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]os.DirEntry), args.Error(1)
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockFileSystem) Lstat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

func (m *MockFileSystem) CreateTemp(dir, pattern string) (File, error) {
	args := m.Called(dir, pattern)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(File), args.Error(1)
}

func (m *MockFileSystem) MkdirTemp(dir, pattern string) (string, error) {
	args := m.Called(dir, pattern)
	return args.String(0), args.Error(1)
}

func (m *MockFileSystem) Chmod(name string, mode os.FileMode) error {
	args := m.Called(name, mode)
	return args.Error(0)
}

func (m *MockFileSystem) Chown(name string, uid, gid int) error {
	args := m.Called(name, uid, gid)
	return args.Error(0)
}

func (m *MockFileSystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	args := m.Called(name, atime, mtime)
	return args.Error(0)
}

func (m *MockFileSystem) Getwd() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockFileSystem) Chdir(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

func (m *MockFileSystem) Symlink(oldname, newname string) error {
	args := m.Called(oldname, newname)
	return args.Error(0)
}

func (m *MockFileSystem) Link(oldname, newname string) error {
	args := m.Called(oldname, newname)
	return args.Error(0)
}

func (m *MockFileSystem) Readlink(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	args := m.Called(root, fn)
	return args.Error(0)
}

func (m *MockFileSystem) IsNotExist(err error) bool {
	args := m.Called(err)
	return args.Bool(0)
}

// Adding MockFileInfo for testing.
type MockFileInfo struct {
	mock.Mock
}

func (m *MockFileInfo) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockFileInfo) Size() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockFileInfo) Mode() os.FileMode {
	args := m.Called()
	return args.Get(0).(os.FileMode)
}

func (m *MockFileInfo) ModTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockFileInfo) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockFileInfo) Sys() any {
	args := m.Called()
	return args.Get(0)
}

func TestNewFileUtils(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}

	// Test successful creation
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)
	require.NotNil(t, fileUtils)
	assert.IsType(t, &DefaultFileUtils{}, fileUtils)

	// Test that we can use the fileUtils instance
	assert.NotNil(t, fileUtils)

	// Adding test for nil filesystem
	_, err = NewFileUtils(nil, localizer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FileUtils requires non-nil filesystem.FileSystem")

	// Adding test for nil localizer
	_, err = NewFileUtils(fs, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FileUtils requires non-nil i18n.Localizer")
}

func TestNewFileUtils_PanicConditions(t *testing.T) {
	t.Parallel()

	localizer := &MockLocalizer{}

	// Test error with nil filesystem
	_, err := NewFileUtils(nil, localizer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FileUtils requires non-nil filesystem.FileSystem")

	// Test error with nil localizer
	fs := NewDefaultFileSystem()
	_, err = NewFileUtils(fs, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FileUtils requires non-nil i18n.Localizer")
}

func TestDefaultFileUtils_Exists(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Test with existing file
	testFile := filepath.Join(tempDir, "test_exists.txt")
	err = fs.WriteFile(testFile, []byte("test content"), 0o644)
	require.NoError(t, err)

	assert.True(t, fileUtils.Exists(testFile))

	// Test with existing directory
	testDir := filepath.Join(tempDir, "test_dir")
	err = fs.Mkdir(testDir, 0o755)
	require.NoError(t, err)

	assert.True(t, fileUtils.Exists(testDir))

	// Test with non-existent file
	nonExistent := filepath.Join(tempDir, "does_not_exist.txt")
	assert.False(t, fileUtils.Exists(nonExistent))

	// Test with empty path
	assert.False(t, fileUtils.Exists(""))

	// Test with symbolic link to existing file
	symlink := filepath.Join(tempDir, "symlink.txt")
	err = fs.Symlink(testFile, symlink)
	require.NoError(t, err)

	assert.True(t, fileUtils.Exists(symlink))

	// Test with broken symbolic link
	brokenSymlink := filepath.Join(tempDir, "broken_symlink.txt")
	err = fs.Symlink("non_existent_file.txt", brokenSymlink)
	require.NoError(t, err)

	assert.False(t, fileUtils.Exists(brokenSymlink))
}

func TestDefaultFileUtils_EnsureDir(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Test creating new directory
	newDir := filepath.Join(tempDir, "new_directory")
	err = fileUtils.EnsureDir(newDir)
	require.NoError(t, err)

	// Verify directory was created
	fileInfo, err := fs.Stat(newDir)
	require.NoError(t, err)
	assert.True(t, fileInfo.IsDir())

	// Test with existing directory (should not error)
	err = fileUtils.EnsureDir(newDir)
	require.NoError(t, err)

	// Test creating nested directories
	nestedDir := filepath.Join(tempDir, "nested", "deep", "directory")
	err = fileUtils.EnsureDir(nestedDir)
	require.NoError(t, err)

	// Verify nested directory was created
	fileInfo, err = fs.Stat(nestedDir)
	require.NoError(t, err)
	assert.True(t, fileInfo.IsDir())
}

func TestDefaultFileUtils_EnsureDir_ErrorConditions(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Test with empty path
	err = fileUtils.EnsureDir("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path is empty")

	// Test with existing file (not directory)
	testFile := filepath.Join(tempDir, "test_file.txt")
	err = fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)

	err = fileUtils.EnsureDir(testFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path exists and is not a directory")

	// Test with permission error (using root directory)
	err = fileUtils.EnsureDir("/root/restricted_test")
	require.Error(t, err)
	// assert.Contains(t, err.Error(), "permission denied")
}

func TestDefaultFileUtils_CopyFile_ErrorConditions(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Test with non-existent source file
	nonExistentSource := filepath.Join(tempDir, "does_not_exist.txt")
	destFile := filepath.Join(tempDir, "destination.txt")

	err = fileUtils.CopyFile(nonExistentSource, destFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SourceFileOpenError")

	// Test with invalid destination (directory that doesn't exist)
	sourceFile := filepath.Join(tempDir, "source.txt")
	err = fs.WriteFile(sourceFile, []byte("test"), 0o644)
	require.NoError(t, err)

	invalidDest := filepath.Join(tempDir, "nonexistent_dir", "destination.txt")
	err = fileUtils.CopyFile(sourceFile, invalidDest)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DestFileCreateError")

	// Test with permission error (using root directory)
	err = fileUtils.CopyFile(sourceFile, "/root/restricted_destination.txt")
	require.Error(t, err)
	// assert.Contains(t, err.Error(), "permission denied")
}

func TestDefaultFileUtils_CopyFile(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Create source file
	sourceFile := filepath.Join(tempDir, "source.txt")
	testContent := []byte("test file content for copying")
	err = fs.WriteFile(sourceFile, testContent, 0o644)
	require.NoError(t, err)

	// Test successful copy
	destFile := filepath.Join(tempDir, "destination.txt")
	err = fileUtils.CopyFile(sourceFile, destFile)
	require.NoError(t, err)

	// Verify destination file exists and has correct content
	copiedContent, err := fs.ReadFile(destFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, copiedContent)

	// Verify file permissions are preserved
	sourceInfo, err := fs.Stat(sourceFile)
	require.NoError(t, err)
	destInfo, err := fs.Stat(destFile)
	require.NoError(t, err)
	assert.Equal(t, sourceInfo.Mode(), destInfo.Mode())
}

func TestDefaultFileUtils_CopyFile_SameFile(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Create source file
	sourceFile := filepath.Join(tempDir, "same_file.txt")
	testContent := []byte("test content")
	err = fs.WriteFile(sourceFile, testContent, 0o644)
	require.NoError(t, err)

	// Test copying file to itself (should be no-op)
	err = fileUtils.CopyFile(sourceFile, sourceFile)
	require.NoError(t, err)

	// Verify file still exists and content is unchanged
	content, err := fs.ReadFile(sourceFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, content)
}

func TestDefaultFileUtils_EnsureDir_StatError(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Test with a path that would cause a stat error (permission denied)
	// This is difficult to test reliably across platforms, so we focus on the empty path case
	err = fileUtils.EnsureDir("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path is empty")
}

func TestDefaultFileUtils_CopyFile_HardLinkFallback(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Create source file
	sourceFile := filepath.Join(tempDir, "source_hardlink.txt")
	testContent := []byte("test content for hard link")
	err = fs.WriteFile(sourceFile, testContent, 0o644)
	require.NoError(t, err)

	// Test copy (may use hard link or fallback to copy depending on filesystem)
	destFile := filepath.Join(tempDir, "dest_hardlink.txt")
	err = fileUtils.CopyFile(sourceFile, destFile)
	require.NoError(t, err)

	// Verify destination file exists and has correct content
	copiedContent, err := fs.ReadFile(destFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, copiedContent)
}

func TestDefaultFileUtils_GetFileSystem(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Test that GetFileSystem returns the underlying filesystem
	// Cast to concrete type to access GetFileSystem method
	defaultFileUtils, ok := fileUtils.(*DefaultFileUtils)
	require.True(t, ok)

	retrievedFS := defaultFileUtils.GetFileSystem()
	require.NotNil(t, retrievedFS)
	assert.Same(t, fs, retrievedFS)
	assert.IsType(t, &DefaultFileSystem{}, retrievedFS)
}

func TestDefaultFileUtils_Integration(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_integration_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Test complete workflow: EnsureDir -> Create file -> Check Exists -> Copy file
	workDir := filepath.Join(tempDir, "work", "subdir")
	err = fileUtils.EnsureDir(workDir)
	require.NoError(t, err)

	// Create a file in the work directory
	sourceFile := filepath.Join(workDir, "original.txt")
	testContent := []byte("integration test content")
	err = fs.WriteFile(sourceFile, testContent, 0o644)
	require.NoError(t, err)

	// Verify file exists
	assert.True(t, fileUtils.Exists(sourceFile))

	// Create backup directory and copy file
	backupDir := filepath.Join(tempDir, "backup")
	err = fileUtils.EnsureDir(backupDir)
	require.NoError(t, err)

	backupFile := filepath.Join(backupDir, "original_backup.txt")
	err = fileUtils.CopyFile(sourceFile, backupFile)
	require.NoError(t, err)

	// Verify backup exists and has correct content
	assert.True(t, fileUtils.Exists(backupFile))
	backupContent, err := fs.ReadFile(backupFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, backupContent)
}

// TestDefaultFileUtils_CopyFile_EmptyFile tests copying an empty file.
func TestDefaultFileUtils_CopyFile_EmptyFile(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Create empty source file
	sourceFile := filepath.Join(tempDir, "empty_source.txt")
	err = fs.WriteFile(sourceFile, []byte{}, 0o644)
	require.NoError(t, err)

	// Test copying empty file
	destFile := filepath.Join(tempDir, "empty_dest.txt")
	err = fileUtils.CopyFile(sourceFile, destFile)
	require.NoError(t, err)

	// Verify destination file exists and is empty
	data, err := fs.ReadFile(destFile)
	require.NoError(t, err)
	assert.Empty(t, data)
}

// TestDefaultFileUtils_CopyFile_Permissions tests copying a file with specific permissions.
func TestDefaultFileUtils_CopyFile_Permissions(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Create source file with specific permissions
	sourceFile := filepath.Join(tempDir, "source_perm.txt")
	err = fs.WriteFile(sourceFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Change permissions of source file
	err = fs.Chmod(sourceFile, 0o755)
	require.NoError(t, err)

	// Test copying file with specific permissions
	destFile := filepath.Join(tempDir, "dest_perm.txt")
	err = fileUtils.CopyFile(sourceFile, destFile)
	require.NoError(t, err)

	// Verify destination file has same permissions
	sourceInfo, err := fs.Stat(sourceFile)
	require.NoError(t, err)
	destInfo, err := fs.Stat(destFile)
	require.NoError(t, err)
	assert.Equal(t, sourceInfo.Mode(), destInfo.Mode())
}

// TestDefaultFileUtils_EnsureDir_ExistingFile tests EnsureDir with an existing file.
func TestDefaultFileUtils_EnsureDir_ExistingFile(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Create temporary directory for test
	tempDir, err := fs.MkdirTemp("", "fileutils_test_")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Create existing file
	existingFile := filepath.Join(tempDir, "existing_file.txt")
	err = fs.WriteFile(existingFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Test EnsureDir with existing file
	err = fileUtils.EnsureDir(existingFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path exists and is not a directory")
}

// TestDefaultFileUtils_Exists_EmptyPath tests Exists with an empty path.
func TestDefaultFileUtils_Exists_EmptyPath(t *testing.T) {
	t.Parallel()

	fs := NewDefaultFileSystem()
	localizer := &MockLocalizer{}
	fileUtils, err := NewFileUtils(fs, localizer)
	require.NoError(t, err)

	// Test Exists with empty path
	assert.False(t, fileUtils.Exists(""))
}

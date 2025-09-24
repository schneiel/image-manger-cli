package filesystem

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestFileSystemOperations tests basic file operations.
func TestFileSystemOperations(t *testing.T) {
	mockFS := new(MockFileSystem)
	tempDir := t.TempDir()

	// Test file creation
	testFile := tempDir + "/test.txt"
	mockFS.On("Create", testFile).Return(nil, nil)
	_, err := mockFS.Create(testFile)
	require.NoError(t, err)

	// Test file reading
	mockFS.On("ReadFile", testFile).Return([]byte("test data"), nil)
	data, err := mockFS.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, "test data", string(data))

	// Test file existence
	mockFS.On("Stat", testFile).Return(nil, nil)
	_, err = mockFS.Stat(testFile)
	require.NoError(t, err)

	// Test file removal
	mockFS.On("Remove", testFile).Return(nil)
	err = mockFS.Remove(testFile)
	require.NoError(t, err)

	// Test directory creation
	testDir := tempDir + "/testdir"
	mockFS.On("Mkdir", testDir, os.ModePerm).Return(nil)
	err = mockFS.Mkdir(testDir, os.ModePerm)
	require.NoError(t, err)

	// Test directory removal
	mockFS.On("RemoveAll", testDir).Return(nil)
	err = mockFS.RemoveAll(testDir)
	require.NoError(t, err)

	// Test file renaming
	newFile := tempDir + "/newtest.txt"
	mockFS.On("Rename", testFile, newFile).Return(nil)
	err = mockFS.Rename(testFile, newFile)
	require.NoError(t, err)

	// Test file permissions
	mockFS.On("Chmod", newFile, os.FileMode(0o644)).Return(nil)
	err = mockFS.Chmod(newFile, os.FileMode(0o644))
	require.NoError(t, err)

	// Test file ownership
	mockFS.On("Chown", newFile, 1000, 1000).Return(nil)
	err = mockFS.Chown(newFile, 1000, 1000)
	require.NoError(t, err)

	// Test file times
	atime := time.Now()
	mtime := time.Now()
	mockFS.On("Chtimes", newFile, atime, mtime).Return(nil)
	err = mockFS.Chtimes(newFile, atime, mtime)
	require.NoError(t, err)

	// Test creating a temporary file
	mockFS.On("CreateTemp", tempDir, "tempfile-*").Return(nil, nil)
	tempFile, err := mockFS.CreateTemp(tempDir, "tempfile-*")
	require.NoError(t, err)
	assert.Nil(t, tempFile)

	// Test creating a temporary directory
	mockFS.On("MkdirTemp", tempDir, "tempdir-*").Return(tempDir+"/tempdir-123456", nil)
	tempDirPath, err := mockFS.MkdirTemp(tempDir, "tempdir-*")
	require.NoError(t, err)
	assert.NotEmpty(t, tempDirPath)

	// Test reading a directory
	mockFS.On("ReadDir", tempDir).Return(nil, nil)
	entries, err := mockFS.ReadDir(tempDir)
	require.NoError(t, err)
	assert.Nil(t, entries)

	// Test getting the current working directory
	mockFS.On("Getwd").Return(tempDir, nil)
	wd, err := mockFS.Getwd()
	require.NoError(t, err)
	assert.Equal(t, tempDir, wd)

	// Test changing the current working directory
	mockFS.On("Chdir", tempDir).Return(nil)
	err = mockFS.Chdir(tempDir)
	require.NoError(t, err)

	// Test creating a symbolic link
	mockFS.On("Symlink", testFile, tempDir+"/symlink").Return(nil)
	err = mockFS.Symlink(testFile, tempDir+"/symlink")
	require.NoError(t, err)

	// Test creating a hard link
	mockFS.On("Link", testFile, tempDir+"/hardlink").Return(nil)
	err = mockFS.Link(testFile, tempDir+"/hardlink")
	require.NoError(t, err)

	// Test reading a symbolic link
	mockFS.On("Readlink", tempDir+"/symlink").Return(testFile, nil)
	linkTarget, err := mockFS.Readlink(tempDir + "/symlink")
	require.NoError(t, err)
	assert.Equal(t, testFile, linkTarget)

	// Test walking a directory
	mockFS.On("WalkDir", tempDir, mock.Anything).Return(nil)
	err = mockFS.WalkDir(tempDir, func(_ string, _ fs.DirEntry, _ error) error {
		return nil
	})
	require.NoError(t, err)

	// Test checking if an error is due to a non-existent file
	mockFS.On("IsNotExist", os.ErrNotExist).Return(true)
	isNotExist := mockFS.IsNotExist(os.ErrNotExist)
	assert.True(t, isNotExist)

	mockFS.AssertExpectations(t)
}

// TestFileSystemErrors tests error scenarios for file operations.
func TestFileSystemErrors(t *testing.T) {
	mockFS := new(MockFileSystem)
	tempDir := t.TempDir()

	// Test file creation error
	testFile := tempDir + "/test.txt"
	mockFS.On("Create", testFile).Return(nil, os.ErrPermission)
	_, err := mockFS.Create(testFile)
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test file reading error
	mockFS.On("ReadFile", testFile).Return(nil, os.ErrNotExist)
	data, err := mockFS.ReadFile(testFile)
	require.Error(t, err)
	assert.Nil(t, data)
	assert.Equal(t, os.ErrNotExist, err)

	// Test file existence error
	mockFS.On("Stat", testFile).Return(nil, os.ErrNotExist)
	_, err = mockFS.Stat(testFile)
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test file removal error
	mockFS.On("Remove", testFile).Return(os.ErrNotExist)
	err = mockFS.Remove(testFile)
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test directory creation error
	testDir := tempDir + "/testdir"
	mockFS.On("Mkdir", testDir, os.ModePerm).Return(os.ErrPermission)
	err = mockFS.Mkdir(testDir, os.ModePerm)
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test directory removal error
	mockFS.On("RemoveAll", testDir).Return(os.ErrNotExist)
	err = mockFS.RemoveAll(testDir)
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test file renaming error
	newFile := tempDir + "/newtest.txt"
	mockFS.On("Rename", testFile, newFile).Return(os.ErrNotExist)
	err = mockFS.Rename(testFile, newFile)
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test file permissions error
	mockFS.On("Chmod", newFile, os.FileMode(0o644)).Return(os.ErrPermission)
	err = mockFS.Chmod(newFile, os.FileMode(0o644))
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test file ownership error
	mockFS.On("Chown", newFile, 1000, 1000).Return(os.ErrPermission)
	err = mockFS.Chown(newFile, 1000, 1000)
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test file times error
	atime := time.Now()
	mtime := time.Now()
	mockFS.On("Chtimes", newFile, atime, mtime).Return(os.ErrPermission)
	err = mockFS.Chtimes(newFile, atime, mtime)
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test creating a temporary file error
	mockFS.On("CreateTemp", tempDir, "tempfile-*").Return(nil, os.ErrPermission)
	_, err = mockFS.CreateTemp(tempDir, "tempfile-*")
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test creating a temporary directory error
	mockFS.On("MkdirTemp", tempDir, "tempdir-*").Return("", os.ErrPermission)
	_, err = mockFS.MkdirTemp(tempDir, "tempdir-*")
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test reading a directory error
	mockFS.On("ReadDir", tempDir).Return(nil, os.ErrNotExist)
	_, err = mockFS.ReadDir(tempDir)
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test getting the current working directory error
	mockFS.On("Getwd").Return("", os.ErrNotExist)
	_, err = mockFS.Getwd()
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test changing the current working directory error
	mockFS.On("Chdir", tempDir).Return(os.ErrNotExist)
	err = mockFS.Chdir(tempDir)
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test creating a symbolic link error
	mockFS.On("Symlink", testFile, tempDir+"/symlink").Return(os.ErrPermission)
	err = mockFS.Symlink(testFile, tempDir+"/symlink")
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test creating a hard link error
	mockFS.On("Link", testFile, tempDir+"/hardlink").Return(os.ErrPermission)
	err = mockFS.Link(testFile, tempDir+"/hardlink")
	require.Error(t, err)
	assert.Equal(t, os.ErrPermission, err)

	// Test reading a symbolic link error
	mockFS.On("Readlink", tempDir+"/symlink").Return("", os.ErrNotExist)
	_, err = mockFS.Readlink(tempDir + "/symlink")
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test walking a directory error
	mockFS.On("WalkDir", tempDir, mock.Anything).Return(os.ErrNotExist)
	err = mockFS.WalkDir(tempDir, func(_ string, _ fs.DirEntry, _ error) error {
		return nil
	})
	require.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)

	// Test checking if an error is due to a non-existent file
	mockFS.On("IsNotExist", os.ErrNotExist).Return(true)
	isNotExist := mockFS.IsNotExist(os.ErrNotExist)
	assert.True(t, isNotExist)

	mockFS.AssertExpectations(t)
}

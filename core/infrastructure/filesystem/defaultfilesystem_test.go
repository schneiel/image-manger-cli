package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultFileSystem_Integration(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}

	// Test basic file operations
	testFile := "test_file.txt"
	testContent := []byte("test content")

	// Write file
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Read file
	data, err := fs.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)

	// Stat file
	fileInfo, err := fs.Stat(testFile)
	require.NoError(t, err)
	assert.NotNil(t, fileInfo)
	assert.Equal(t, testFile, fileInfo.Name())

	// Remove file
	err = fs.Remove(testFile)
	require.NoError(t, err)

	// Verify file is removed
	_, err = fs.Stat(testFile)
	require.Error(t, err)
}

func TestDefaultFileSystem_DirectoryOperations(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "test_directory")

	// Create directory
	err := fs.Mkdir(testDir, 0o755)
	require.NoError(t, err)

	// Check directory exists
	fileInfo, err := fs.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, fileInfo.IsDir())

	// Remove directory
	err = fs.RemoveAll(testDir)
	require.NoError(t, err)

	// Verify directory is removed
	_, err = fs.Stat(testDir)
	require.Error(t, err)

	// Test permission error
	err = fs.Mkdir("/root/restricted_test_directory", 0o755)
	require.Error(t, err)
	// assert.Contains(t, err.Error(), "permission denied")
}

func TestDefaultFileSystem_FileOperations_ErrorConditions(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}

	// Test non-existent file error
	_, err := fs.Open("non_existent_file.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")

	// Test non-existent file error for OpenFile
	_, err = fs.OpenFile("non_existent_file.txt", os.O_RDONLY, 0o644)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")

	// Test non-existent file error for ReadFile
	_, err = fs.ReadFile("non_existent_file.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")

	// Test non-existent file error for Remove
	err = fs.Remove("non_existent_file.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")

	// Test non-existent directory error for ReadDir
	_, err = fs.ReadDir("non_existent_directory")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")

	// Test non-existent file error for Stat
	_, err = fs.Stat("non_existent_file.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")

	// Test non-existent file error for Lstat
	_, err = fs.Lstat("non_existent_file.txt")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")

	// Test permission errors only attempt operations that reliably fail
	// Most systems will deny access to /root/ for non-root users
	_, err = fs.Create("/root/test_file.txt")
	if err == nil {
		// If running as root or on a permissive system, skip permission tests
		t.Skip("Skipping permission tests - running with elevated privileges or permissive filesystem")
	}
	// If we get here, permission restrictions are working as expected
}

func TestDefaultFileSystem_FileInfo(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}

	testFile := "test_info.txt"
	testContent := []byte("test content")

	// Create file
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(testFile) }()

	// Get file info
	fileInfo, err := fs.Stat(testFile)
	require.NoError(t, err)

	// Verify file info
	assert.Equal(t, testFile, fileInfo.Name())
	assert.Equal(t, int64(len(testContent)), fileInfo.Size())
	assert.False(t, fileInfo.IsDir())
	assert.True(t, fileInfo.ModTime().Before(time.Now().Add(time.Second)))
}

func TestDefaultFileSystem_ErrorHandling(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}

	// Test reading non-existent file
	_, err := fs.ReadFile("nonexistent.txt")
	require.Error(t, err)

	// Test statting non-existent file
	_, err = fs.Stat("nonexistent.txt")
	require.Error(t, err)

	// Test removing non-existent file
	err = fs.Remove("nonexistent.txt")
	require.Error(t, err)
}

func TestNewDefaultFileSystem(t *testing.T) {
	t.Parallel()
	fs := NewDefaultFileSystem()
	assert.NotNil(t, fs)
	assert.IsType(t, &DefaultFileSystem{}, fs)
}

func TestDefaultFileSystem_OpenFile(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile.txt")
	testContent := []byte("test content for openfile")

	// Create file first
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Test OpenFile with different modes
	file, err := fs.OpenFile(testFile, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Read from opened file
	data := make([]byte, len(testContent))
	n, err := file.Read(data)
	require.NoError(t, err)
	assert.Equal(t, len(testContent), n)
	assert.Equal(t, testContent, data)

	err = file.Close()
	require.NoError(t, err)
}

// TestDefaultFileSystem_OpenFile_Read tests reading from an opened file.
func TestDefaultFileSystem_OpenFile_Read(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_read.txt")
	testContent := []byte("test content for openfile read")

	// Create file first
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Open file for reading
	file, err := fs.OpenFile(testFile, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Read from opened file
	data := make([]byte, len(testContent))
	n, err := file.Read(data)
	require.NoError(t, err)
	assert.Equal(t, len(testContent), n)
	assert.Equal(t, testContent, data)

	err = file.Close()
	require.NoError(t, err)
}

// TestDefaultFileSystem_OpenFile_Write tests writing to an opened file.
func TestDefaultFileSystem_OpenFile_Write(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_write.txt")

	// Open file for writing
	file, err := fs.OpenFile(testFile, os.O_WRONLY|os.O_CREATE, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Write to opened file
	testContent := []byte("test content for openfile write")
	n, err := file.Write(testContent)
	require.NoError(t, err)
	assert.Equal(t, len(testContent), n)

	err = file.Close()
	require.NoError(t, err)

	// Verify file was written
	data, err := fs.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)
}

// TestDefaultFileSystem_OpenFile_Append tests appending to an opened file.
func TestDefaultFileSystem_OpenFile_Append(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_append.txt")
	initialContent := []byte("initial content")
	appendContent := []byte(" appended content")

	// Create file with initial content
	err := fs.WriteFile(testFile, initialContent, 0o644)
	require.NoError(t, err)

	// Open file for appending
	file, err := fs.OpenFile(testFile, os.O_WRONLY|os.O_APPEND, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Append to opened file
	n, err := file.Write(appendContent)
	require.NoError(t, err)
	assert.Equal(t, len(appendContent), n)

	err = file.Close()
	require.NoError(t, err)

	// Verify file was appended
	data, err := fs.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, append(initialContent, appendContent...), data)
}

// TestDefaultFileSystem_OpenFile_Seek tests seeking within an opened file.
func TestDefaultFileSystem_OpenFile_Seek(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_seek.txt")
	testContent := []byte("test content for openfile seek")

	// Create file first
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Open file for reading
	file, err := fs.OpenFile(testFile, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Seek to middle of file
	offset, err := file.Seek(5, io.SeekStart)
	require.NoError(t, err)
	assert.Equal(t, int64(5), offset)

	// Read from seek position
	data := make([]byte, 7)
	n, err := file.Read(data)
	require.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, testContent[5:12], data)

	err = file.Close()
	require.NoError(t, err)
}

// TestDefaultFileSystem_OpenFile_Truncate tests truncating an opened file.
func TestDefaultFileSystem_OpenFile_Truncate(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_truncate.txt")
	testContent := []byte("test content for openfile truncate")

	// Create file first
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Open file for writing
	file, err := fs.OpenFile(testFile, os.O_WRONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Truncate file
	err = file.Truncate(5)
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	// Verify file was truncated
	data, err := fs.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, testContent[:5], data)
}

// TestDefaultFileSystem_OpenFile_Sync tests syncing an opened file.
func TestDefaultFileSystem_OpenFile_Sync(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_sync.txt")
	testContent := []byte("test content for openfile sync")

	// Open file for writing
	file, err := fs.OpenFile(testFile, os.O_WRONLY|os.O_CREATE, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Write to opened file
	n, err := file.Write(testContent)
	require.NoError(t, err)
	assert.Equal(t, len(testContent), n)

	// Sync file
	err = file.Sync()
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	// Verify file was written
	data, err := fs.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)
}

// TestDefaultFileSystem_OpenFile_Readdir tests reading directory entries from an opened directory.
func TestDefaultFileSystem_OpenFile_Readdir(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_readdir.txt")

	// Create file first
	err := fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Open directory for reading
	file, err := fs.OpenFile(tempDir, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Read directory entries
	entries, err := file.Readdir(1)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "test_openfile_readdir.txt", entries[0].Name())

	err = file.Close()
	require.NoError(t, err)
}

// TestDefaultFileSystem_OpenFile_Readdirnames tests reading directory entry names from an opened directory.
func TestDefaultFileSystem_OpenFile_Readdirnames(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_readdirnames.txt")

	// Create file first
	err := fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Open directory for reading
	file, err := fs.OpenFile(tempDir, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Read directory entry names
	names, err := file.Readdirnames(1)
	require.NoError(t, err)
	assert.Len(t, names, 1)
	assert.Equal(t, "test_openfile_readdirnames.txt", names[0])

	err = file.Close()
	require.NoError(t, err)
}

// TestDefaultFileSystem_OpenFile_Name tests getting the name of an opened file.
func TestDefaultFileSystem_OpenFile_Name(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_name.txt")

	// Create file first
	err := fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Open file for reading
	file, err := fs.OpenFile(testFile, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Get file name
	name := file.Name()
	assert.Equal(t, testFile, name)

	err = file.Close()
	require.NoError(t, err)
}

// TestDefaultFileSystem_OpenFile_Stat tests getting file info from an opened file.
func TestDefaultFileSystem_OpenFile_Stat(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_openfile_stat.txt")
	testContent := []byte("test content for openfile stat")

	// Create file first
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Open file for reading
	file, err := fs.OpenFile(testFile, os.O_RDONLY, 0o644)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Get file info
	info, err := file.Stat()
	require.NoError(t, err)
	assert.Equal(t, "test_openfile_stat.txt", info.Name())
	assert.Equal(t, int64(len(testContent)), info.Size())

	err = file.Close()
	require.NoError(t, err)
}

func TestDefaultFileSystem_Create(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_create.txt")

	// Create file
	file, err := fs.Create(testFile)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Write to created file
	testContent := []byte("created content")
	n, err := file.Write(testContent)
	require.NoError(t, err)
	assert.Equal(t, len(testContent), n)

	err = file.Close()
	require.NoError(t, err)

	// Verify file was created
	data, err := fs.ReadFile(testFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)
}

func TestDefaultFileSystem_Open(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_open.txt")
	testContent := []byte("content to open")

	// Create file first
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Open file
	file, err := fs.Open(testFile)
	require.NoError(t, err)
	require.NotNil(t, file)

	// Read from opened file
	data := make([]byte, len(testContent))
	n, err := file.Read(data)
	require.NoError(t, err)
	assert.Equal(t, len(testContent), n)
	assert.Equal(t, testContent, data)

	err = file.Close()
	require.NoError(t, err)
}

func TestDefaultFileSystem_Rename(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	oldFile := "test_rename_old.txt"
	newFile := "test_rename_new.txt"
	testContent := []byte("content to rename")

	// Create file
	err := fs.WriteFile(oldFile, testContent, 0o644)
	require.NoError(t, err)

	// Rename file
	err = fs.Rename(oldFile, newFile)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(newFile) }()

	// Verify old file doesn't exist
	_, err = fs.Stat(oldFile)
	require.Error(t, err)

	// Verify new file exists with same content
	data, err := fs.ReadFile(newFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)
}

func TestDefaultFileSystem_MkdirAll(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	testDir := "test_dir/nested/deep"

	// Create nested directories
	err := fs.MkdirAll(testDir, 0o755)
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll("test_dir") }()

	// Verify directory exists
	info, err := fs.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestDefaultFileSystem_ReadDir(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "test_readdir")

	// Create directory and files
	err := fs.Mkdir(testDir, 0o755)
	require.NoError(t, err)

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := fs.WriteFile(filePath, []byte("test"), 0o644)
		require.NoError(t, err)
	}

	// Read directory
	entries, err := fs.ReadDir(testDir)
	require.NoError(t, err)
	assert.Len(t, entries, len(testFiles))

	// Verify all files are present
	for _, entry := range entries {
		assert.Contains(t, testFiles, entry.Name())
		assert.False(t, entry.IsDir())
	}
}

func TestDefaultFileSystem_Lstat(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_lstat.txt")
	testContent := []byte("content for lstat")

	// Create file
	err := fs.WriteFile(testFile, testContent, 0o644)
	require.NoError(t, err)

	// Get file info with Lstat
	info, err := fs.Lstat(testFile)
	require.NoError(t, err)
	assert.Equal(t, filepath.Base(testFile), info.Name())
	assert.Equal(t, int64(len(testContent)), info.Size())
	assert.False(t, info.IsDir())
}

func TestDefaultFileSystem_CreateTemp(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}

	// Create temporary file
	file, err := fs.CreateTemp("", "test_temp_*.txt")
	require.NoError(t, err)
	require.NotNil(t, file)

	filename := file.Name()
	defer func() { _ = fs.Remove(filename) }()

	// Write to temp file
	testContent := []byte("temporary content")
	n, err := file.Write(testContent)
	require.NoError(t, err)
	assert.Equal(t, len(testContent), n)

	err = file.Close()
	require.NoError(t, err)

	// Verify temp file exists and has content
	data, err := fs.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)
}

func TestDefaultFileSystem_MkdirTemp(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}

	// Create temporary directory
	tempDir, err := fs.MkdirTemp("", "test_temp_dir_*")
	require.NoError(t, err)
	require.NotEmpty(t, tempDir)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Verify temp directory exists
	info, err := fs.Stat(tempDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestDefaultFileSystem_Chmod(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_chmod.txt")

	// Create file
	err := fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Change permissions
	err = fs.Chmod(testFile, 0o600)
	require.NoError(t, err)

	// Verify permissions changed
	info, err := fs.Stat(testFile)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestDefaultFileSystem_Chtimes(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	testFile := "test_chtimes.txt"

	// Create file
	err := fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(testFile) }()

	// Set specific times
	atime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mtime := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

	err = fs.Chtimes(testFile, atime, mtime)
	require.NoError(t, err)

	// Verify modification time changed
	info, err := fs.Stat(testFile)
	require.NoError(t, err)
	assert.True(t, info.ModTime().Equal(mtime))
}

func TestDefaultFileSystem_WorkingDirectory(t *testing.T) {
	// Don't run in parallel since we're changing working directory
	fs := &DefaultFileSystem{}

	// Get current working directory
	originalWd, err := fs.Getwd()
	require.NoError(t, err)
	require.NotEmpty(t, originalWd)

	// Create temp directory
	tempDir, err := fs.MkdirTemp("", "test_wd_*")
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(tempDir) }()

	// Resolve symlinks for comparison (macOS has /var -> /private/var symlink)
	resolvedTempDir, err := filepath.EvalSymlinks(tempDir)
	require.NoError(t, err)

	// Change to temp directory
	err = fs.Chdir(tempDir)
	require.NoError(t, err)
	defer func() { _ = fs.Chdir(originalWd) }() // Restore original

	// Verify we're in the temp directory
	currentWd, err := fs.Getwd()
	require.NoError(t, err)
	assert.Equal(t, resolvedTempDir, currentWd)
}

func TestDefaultFileSystem_Links(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	originalFile := "test_link_original.txt"
	hardLink := "test_hardlink.txt"
	symLink := "test_symlink.txt"
	testContent := []byte("content for linking")

	// Create original file
	err := fs.WriteFile(originalFile, testContent, 0o644)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(originalFile) }()

	// Create hard link
	err = fs.Link(originalFile, hardLink)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(hardLink) }()

	// Verify hard link has same content
	data, err := fs.ReadFile(hardLink)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)

	// Create symbolic link
	err = fs.Symlink(originalFile, symLink)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(symLink) }()

	// Verify symbolic link points to original
	target, err := fs.Readlink(symLink)
	require.NoError(t, err)
	assert.Equal(t, originalFile, target)

	// Verify symbolic link has same content
	data, err = fs.ReadFile(symLink)
	require.NoError(t, err)
	assert.Equal(t, testContent, data)
}

func TestDefaultFileSystem_WalkDir(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	rootDir := "test_walkdir"

	// Create directory structure
	err := fs.MkdirAll(filepath.Join(rootDir, "subdir1"), 0o755)
	require.NoError(t, err)
	err = fs.MkdirAll(filepath.Join(rootDir, "subdir2"), 0o755)
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(rootDir) }()

	// Create files
	testFiles := []string{
		filepath.Join(rootDir, "file1.txt"),
		filepath.Join(rootDir, "subdir1", "file2.txt"),
		filepath.Join(rootDir, "subdir2", "file3.txt"),
	}
	for _, file := range testFiles {
		err := fs.WriteFile(file, []byte("test"), 0o644)
		require.NoError(t, err)
	}

	// Walk directory tree
	var visitedPaths []string
	err = fs.WalkDir(rootDir, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		visitedPaths = append(visitedPaths, path)
		return nil
	})
	require.NoError(t, err)

	// Verify all paths were visited
	assert.Contains(t, visitedPaths, rootDir)
	assert.Contains(t, visitedPaths, filepath.Join(rootDir, "subdir1"))
	assert.Contains(t, visitedPaths, filepath.Join(rootDir, "subdir2"))
	for _, file := range testFiles {
		assert.Contains(t, visitedPaths, file)
	}
}

func TestDefaultFileSystem_IsNotExist(t *testing.T) {
	// Not running in parallel due to file system interactions
	fs := &DefaultFileSystem{}

	// Test with file that doesn't exist
	_, err := fs.Stat("nonexistent_file.txt")
	require.Error(t, err)
	assert.True(t, fs.IsNotExist(err))

	// Test with file that exists
	testFile := "test_exists.txt"
	err = fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(testFile) }()

	_, err = fs.Stat(testFile)
	require.NoError(t, err)
	assert.False(t, fs.IsNotExist(err))
}

func TestDefaultFileSystem_EdgeCases(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}

	// Test creating file with empty content
	emptyFile := "test_empty.txt"
	err := fs.WriteFile(emptyFile, []byte{}, 0o644)
	require.NoError(t, err)
	defer func() { _ = fs.Remove(emptyFile) }()

	data, err := fs.ReadFile(emptyFile)
	require.NoError(t, err)
	assert.Empty(t, data)

	// Test mkdir with existing directory (should succeed)
	testDir := "test_existing_dir"
	err = fs.Mkdir(testDir, 0o755)
	require.NoError(t, err)
	defer func() { _ = fs.RemoveAll(testDir) }()

	// Second mkdir should fail
	err = fs.Mkdir(testDir, 0o755)
	require.Error(t, err)

	// But MkdirAll should succeed
	err = fs.MkdirAll(testDir, 0o755)
	require.NoError(t, err)
}

// TestDefaultFileSystem_Chown tests the Chown method.
func TestDefaultFileSystem_Chown(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_chown.txt")

	// Create file
	err := fs.WriteFile(testFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Change ownership (may fail on some systems due to permissions)
	err = fs.Chown(testFile, 1000, 1000)
	if err != nil {
		// Skip test if we don't have permission (common on macOS)
		t.Skipf("Chown test skipped due to permission restrictions: %v", err)
		return
	}

	// Verify ownership changed
	info, err := fs.Stat(testFile)
	require.NoError(t, err)
	assert.Equal(t, 1000, info.Sys().(*syscall.Stat_t).Uid)
	assert.Equal(t, 1000, info.Sys().(*syscall.Stat_t).Gid)
}

// TestDefaultFileSystem_Symlink tests the Symlink method.
func TestDefaultFileSystem_Symlink(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	originalFile := filepath.Join(tempDir, "original.txt")
	symlinkFile := filepath.Join(tempDir, "symlink.txt")

	// Create original file
	err := fs.WriteFile(originalFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Create symlink
	err = fs.Symlink(originalFile, symlinkFile)
	require.NoError(t, err)

	// Verify symlink points to original file
	target, err := fs.Readlink(symlinkFile)
	require.NoError(t, err)
	assert.Equal(t, originalFile, target)

	// Verify symlink content
	data, err := fs.ReadFile(symlinkFile)
	require.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
}

// TestDefaultFileSystem_Link tests the Link method.
func TestDefaultFileSystem_Link(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	originalFile := filepath.Join(tempDir, "original.txt")
	hardlinkFile := filepath.Join(tempDir, "hardlink.txt")

	// Create original file
	err := fs.WriteFile(originalFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Create hard link
	err = fs.Link(originalFile, hardlinkFile)
	require.NoError(t, err)

	// Verify hard link content
	data, err := fs.ReadFile(hardlinkFile)
	require.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
}

// TestDefaultFileSystem_Readlink tests the Readlink method.
func TestDefaultFileSystem_Readlink(t *testing.T) {
	t.Parallel()
	fs := &DefaultFileSystem{}
	tempDir := t.TempDir()
	originalFile := filepath.Join(tempDir, "original.txt")
	symlinkFile := filepath.Join(tempDir, "symlink.txt")

	// Create original file
	err := fs.WriteFile(originalFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Create symlink
	err = fs.Symlink(originalFile, symlinkFile)
	require.NoError(t, err)

	// Verify symlink points to original file
	target, err := fs.Readlink(symlinkFile)
	require.NoError(t, err)
	assert.Equal(t, originalFile, target)
}

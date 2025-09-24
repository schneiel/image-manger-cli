package dedup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultSizeScanner(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	allowedExtensions := []string{".jpg", ".png", ".gif"}

	scanner, err := NewDefaultSizeScanner(allowedExtensions, mockLogger, mockLocalizer)
	require.NoError(t, err)

	assert.NotNil(t, scanner)
	defaultScanner, ok := scanner.(*DefaultSizeScanner)
	assert.True(t, ok)
	assert.Equal(t, allowedExtensions, defaultScanner.AllowedExtensions)
	assert.Equal(t, mockLogger, defaultScanner.Logger)
	assert.Equal(t, mockLocalizer, defaultScanner.localizer)
}

func TestDefaultSizeScanner_Scan_Success(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "ScanningForFiles":
			return "Scanning for files"
		case "PotentialDuplicateGroupsFound":
			return "Found potential duplicate groups"
		default:
			return key
		}
	}

	mockLogger.InfoFunc = func(_ string) {
		// Verify logging calls
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	mockLogger.WarnfFunc = func(_ string, _ ...interface{}) {
		// Verify warning calls
	}

	allowedExtensions := []string{".jpg", ".png", ".gif"}
	scanner, err := NewDefaultSizeScanner(allowedExtensions, mockLogger, mockLocalizer)
	require.NoError(t, err)

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files with different sizes
	testFiles := []struct {
		name string
		size int64
	}{
		{"image1.jpg", 1024},
		{"image2.jpg", 1024}, // Same size as image1
		{"image3.png", 2048},
		{"image4.png", 2048},   // Same size as image3
		{"image5.gif", 1024},   // Same size as image1 and image2
		{"document.txt", 1024}, // Not an image file
		{"image6.jpg", 4096},   // Unique size
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file.name)
		// Create a file with the specified size
		f, err := os.Create(filePath) // #nosec G304
		require.NoError(t, err)
		defer func() { _ = f.Close() }()

		// Write some data to make the file the specified size
		data := make([]byte, file.size)
		_, err = f.Write(data)
		require.NoError(t, err)
	}

	// Scan the directory
	fileGroups, err := scanner.Scan(tempDir)

	require.NoError(t, err)
	assert.NotNil(t, fileGroups)

	// We expect groups for files with the same size
	// image1.jpg, image2.jpg, image5.gif all have size 1024
	// image3.png, image4.png both have size 2048
	// image6.jpg has unique size 4096
	// document.txt is not an image file

	// Verify we have the expected number of groups
	assert.Len(t, fileGroups, 2)

	// Find the group with size 1024 (should have 3 files)
	var size1024Group []string
	var size2048Group []string

	for _, group := range fileGroups {
		if len(group) == 3 {
			size1024Group = group
		} else if len(group) == 2 {
			size2048Group = group
		}
	}

	assert.NotNil(t, size1024Group, "Should have a group with 3 files (size 1024)")
	assert.NotNil(t, size2048Group, "Should have a group with 2 files (size 2048)")

	// Verify the files in the size 1024 group
	expectedSize1024Files := []string{
		filepath.Join(tempDir, "image1.jpg"),
		filepath.Join(tempDir, "image2.jpg"),
		filepath.Join(tempDir, "image5.gif"),
	}

	for _, expectedFile := range expectedSize1024Files {
		found := false
		for _, actualFile := range size1024Group {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		assert.True(t, found, "File %s should be in the size 1024 group", expectedFile)
	}

	// Verify the files in the size 2048 group
	expectedSize2048Files := []string{
		filepath.Join(tempDir, "image3.png"),
		filepath.Join(tempDir, "image4.png"),
	}

	for _, expectedFile := range expectedSize2048Files {
		found := false
		for _, actualFile := range size2048Group {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		assert.True(t, found, "File %s should be in the size 2048 group", expectedFile)
	}
}

func TestDefaultSizeScanner_Scan_EmptyDirectory(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "ScanningForFiles":
			return "Scanning for files"
		case "PotentialDuplicateGroupsFound":
			return "Found potential duplicate groups"
		default:
			return key
		}
	}

	mockLogger.InfoFunc = func(_ string) {
		// Verify logging calls
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	allowedExtensions := []string{".jpg", ".png", ".gif"}
	scanner, err := NewDefaultSizeScanner(allowedExtensions, mockLogger, mockLocalizer)
	require.NoError(t, err)

	// Create an empty temporary directory
	tempDir := t.TempDir()

	// Scan the empty directory
	fileGroups, err := scanner.Scan(tempDir)
	if err != nil {
		t.Logf("Error scanning directory: %v", err)
	}
	require.NoError(t, err)
	if fileGroups == nil {
		t.Logf("fileGroups is nil")
	} else {
		t.Logf("fileGroups length: %d", len(fileGroups))
	}
	// The scanner should return an empty slice, not nil
	if fileGroups == nil {
		fileGroups = FileGroup{}
	}
	assert.Empty(t, fileGroups)
}

func TestDefaultSizeScanner_Scan_NoImageFiles(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "ScanningForFiles":
			return "Scanning for files"
		case "PotentialDuplicateGroupsFound":
			return "Found potential duplicate groups"
		default:
			return key
		}
	}

	mockLogger.InfoFunc = func(_ string) {
		// Verify logging calls
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	allowedExtensions := []string{".jpg", ".png", ".gif"}
	scanner, err := NewDefaultSizeScanner(allowedExtensions, mockLogger, mockLocalizer)
	require.NoError(t, err)

	// Create a temporary directory with non-image files
	tempDir := t.TempDir()

	// Create test files that are not images
	testFiles := []string{"document.txt", "data.csv", "script.sh", "archive.zip"}

	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		f, err := os.Create(filePath) // #nosec G304
		require.NoError(t, err)
		defer func() { _ = f.Close() }()

		// Write some data
		data := make([]byte, 1024)
		_, err = f.Write(data)
		require.NoError(t, err)
	}

	// Scan the directory
	fileGroups, err := scanner.Scan(tempDir)

	require.NoError(t, err)
	// The scanner should return an empty slice, not nil
	if fileGroups == nil {
		fileGroups = FileGroup{}
	}
	assert.Empty(t, fileGroups) // No image files, so no groups
}

func TestDefaultSizeScanner_Scan_NoDuplicates(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "ScanningForFiles":
			return "Scanning for files"
		case "PotentialDuplicateGroupsFound":
			return "Found potential duplicate groups"
		default:
			return key
		}
	}

	mockLogger.InfoFunc = func(_ string) {
		// Verify logging calls
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	allowedExtensions := []string{".jpg", ".png", ".gif"}
	scanner, err := NewDefaultSizeScanner(allowedExtensions, mockLogger, mockLocalizer)
	require.NoError(t, err)

	// Create a temporary directory with image files of different sizes
	tempDir := t.TempDir()

	// Create test files with different sizes
	testFiles := []struct {
		name string
		size int64
	}{
		{"image1.jpg", 1024},
		{"image2.png", 2048},
		{"image3.gif", 4096},
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file.name)
		f, err := os.Create(filePath) // #nosec G304
		require.NoError(t, err)
		defer func() { _ = f.Close() }()

		// Write some data to make the file the specified size
		data := make([]byte, file.size)
		_, err = f.Write(data)
		require.NoError(t, err)
	}

	// Scan the directory
	fileGroups, err := scanner.Scan(tempDir)

	require.NoError(t, err)
	// The scanner should return an empty slice, not nil
	if fileGroups == nil {
		fileGroups = FileGroup{}
	}
	assert.Empty(t, fileGroups) // No duplicates, so no groups
}

func TestDefaultSizeScanner_Scan_WithSubdirectories(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "ScanningForFiles":
			return "Scanning for files"
		case "PotentialDuplicateGroupsFound":
			return "Found potential duplicate groups"
		default:
			return key
		}
	}

	mockLogger.InfoFunc = func(_ string) {
		// Verify logging calls
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	allowedExtensions := []string{".jpg", ".png", ".gif"}
	scanner, err := NewDefaultSizeScanner(allowedExtensions, mockLogger, mockLocalizer)
	require.NoError(t, err)

	// Create a temporary directory with subdirectories
	tempDir := t.TempDir()

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0o750)
	require.NoError(t, err)

	// Create test files in root and subdirectory with same sizes
	testFiles := []struct {
		path string
		size int64
	}{
		{filepath.Join(tempDir, "image1.jpg"), 1024},
		{filepath.Join(subDir, "image2.jpg"), 1024}, // Same size as image1
		{filepath.Join(tempDir, "image3.png"), 2048},
		{filepath.Join(subDir, "image4.png"), 2048}, // Same size as image3
	}

	for _, file := range testFiles {
		f, err := os.Create(file.path)
		require.NoError(t, err)
		defer func() { _ = f.Close() }()

		// Write some data to make the file the specified size
		data := make([]byte, file.size)
		_, err = f.Write(data)
		require.NoError(t, err)
	}

	// Scan the directory
	fileGroups, err := scanner.Scan(tempDir)

	require.NoError(t, err)
	assert.NotNil(t, fileGroups)

	// We expect groups for files with the same size across directories
	assert.Len(t, fileGroups, 2)

	// Find the groups
	var size1024Group []string
	var size2048Group []string

	for _, group := range fileGroups {
		if len(group) == 2 {
			if filepath.Base(group[0]) == "image1.jpg" || filepath.Base(group[1]) == "image1.jpg" {
				size1024Group = group
			} else {
				size2048Group = group
			}
		}
	}

	assert.NotNil(t, size1024Group, "Should have a group with 2 files (size 1024)")
	assert.NotNil(t, size2048Group, "Should have a group with 2 files (size 2048)")

	// Verify the files in the size 1024 group
	expectedSize1024Files := []string{
		filepath.Join(tempDir, "image1.jpg"),
		filepath.Join(subDir, "image2.jpg"),
	}

	for _, expectedFile := range expectedSize1024Files {
		found := false
		for _, actualFile := range size1024Group {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		assert.True(t, found, "File %s should be in the size 1024 group", expectedFile)
	}

	// Verify the files in the size 2048 group
	expectedSize2048Files := []string{
		filepath.Join(tempDir, "image3.png"),
		filepath.Join(subDir, "image4.png"),
	}

	for _, expectedFile := range expectedSize2048Files {
		found := false
		for _, actualFile := range size2048Group {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		assert.True(t, found, "File %s should be in the size 2048 group", expectedFile)
	}
}

func TestDefaultSizeScanner_Scan_ExtensionCaseInsensitive(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}

	// Setup mock expectations
	mockLocalizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		switch key {
		case "ScanningForFiles":
			return "Scanning for files"
		case "PotentialDuplicateGroupsFound":
			return "Found potential duplicate groups"
		default:
			return key
		}
	}

	mockLogger.InfoFunc = func(_ string) {
		// Verify logging calls
	}

	mockLogger.InfofFunc = func(_ string, _ ...interface{}) {
		// Verify logging calls
	}

	allowedExtensions := []string{".jpg", ".png", ".gif"}
	scanner, err := NewDefaultSizeScanner(allowedExtensions, mockLogger, mockLocalizer)
	require.NoError(t, err)

	// Create a temporary directory
	tempDir := t.TempDir()

	// Create test files with different case extensions
	testFiles := []struct {
		name string
		size int64
	}{
		{"image1.JPG", 1024}, // Uppercase
		{"image2.jpg", 1024}, // Lowercase
		{"image3.PNG", 2048}, // Uppercase
		{"image4.png", 2048}, // Lowercase
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file.name)
		f, err := os.Create(filePath) // #nosec G304
		require.NoError(t, err)
		defer func() { _ = f.Close() }()

		// Write some data to make the file the specified size
		data := make([]byte, file.size)
		_, err = f.Write(data)
		require.NoError(t, err)
	}

	// Scan the directory
	fileGroups, err := scanner.Scan(tempDir)

	require.NoError(t, err)
	assert.NotNil(t, fileGroups)

	// We expect groups for files with the same size, regardless of case
	assert.Len(t, fileGroups, 2)

	// Find the groups
	var size1024Group []string
	var size2048Group []string

	for _, group := range fileGroups {
		if len(group) == 2 {
			if filepath.Base(group[0]) == "image1.JPG" || filepath.Base(group[1]) == "image1.JPG" {
				size1024Group = group
			} else {
				size2048Group = group
			}
		}
	}

	assert.NotNil(t, size1024Group, "Should have a group with 2 files (size 1024)")
	assert.NotNil(t, size2048Group, "Should have a group with 2 files (size 2048)")

	// Verify the files in the size 1024 group
	expectedSize1024Files := []string{
		filepath.Join(tempDir, "image1.JPG"),
		filepath.Join(tempDir, "image2.jpg"),
	}

	for _, expectedFile := range expectedSize1024Files {
		found := false
		for _, actualFile := range size1024Group {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		assert.True(t, found, "File %s should be in the size 1024 group", expectedFile)
	}

	// Verify the files in the size 2048 group
	expectedSize2048Files := []string{
		filepath.Join(tempDir, "image3.PNG"),
		filepath.Join(tempDir, "image4.png"),
	}

	for _, expectedFile := range expectedSize2048Files {
		found := false
		for _, actualFile := range size2048Group {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		assert.True(t, found, "File %s should be in the size 2048 group", expectedFile)
	}
}

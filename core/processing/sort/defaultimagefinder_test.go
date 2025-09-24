package sort

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultImageFinder(t *testing.T) {
	t.Parallel()
	allowedExtensions := []string{".jpg", ".png", ".gif", ".JPEG", ".TIFF"}

	finder := NewDefaultImageFinder(allowedExtensions)

	require.NotNil(t, finder, "Expected finder to be created")

	defaultFinder, ok := finder.(*DefaultImageFinder)
	require.True(t, ok, "Expected DefaultImageFinder type")

	// Check that extensions are converted to lowercase
	expectedExtensions := []string{".jpg", ".png", ".gif", ".jpeg", ".tiff"}
	assert.Len(t, defaultFinder.allowedExtensions, len(expectedExtensions))

	for i, ext := range expectedExtensions {
		assert.Equal(t, ext, defaultFinder.allowedExtensions[i])
	}
}

func TestDefaultImageFinder_Find(t *testing.T) { //nolint:gocognit // comprehensive test with multiple scenarios
	t.Parallel()

	testCases := []struct {
		name          string
		setup         func(t *testing.T) string
		allowedExts   []string
		expectedCount int
		expectedError bool
	}{
		{
			name: "Successful find",
			setup: func(t *testing.T) string {
				t.Helper()
				tempDir, err := os.MkdirTemp("", "imagefinder_test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

				testFiles := []string{
					"image1.jpg",
					"image2.png",
					"document.txt",
					"image3.gif",
					"image4.JPEG",
					"image5.tiff",
					"subdir/image6.jpg",
					"subdir/document.pdf",
				}

				subDir := filepath.Join(tempDir, "subdir")
				if err := os.Mkdir(subDir, 0o750); err != nil {
					t.Fatalf("Failed to create subdir: %v", err)
				}

				for _, fileName := range testFiles {
					filePath := filepath.Join(tempDir, fileName)
					err := os.WriteFile(filePath, []byte("test content"), 0o600)
					if err != nil {
						t.Fatalf("Failed to create test file %s: %v", fileName, err)
					}
				}

				return tempDir
			},
			allowedExts:   []string{".jpg", ".png", ".gif", ".jpeg", ".tiff"},
			expectedCount: 6,
		},
		{
			name: "Empty directory",
			setup: func(t *testing.T) string {
				t.Helper()
				tempDir, err := os.MkdirTemp("", "imagefinder_empty_test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				t.Cleanup(func() { _ = os.RemoveAll(tempDir) })
				return tempDir
			},
			allowedExts:   []string{".jpg", ".png"},
			expectedCount: 0,
		},
		{
			name: "No image files",
			setup: func(t *testing.T) string {
				t.Helper()
				tempDir, err := os.MkdirTemp("", "imagefinder_noimages_test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

				testFiles := []string{
					"document.txt",
					"data.csv",
					"script.py",
					"config.json",
				}

				for _, fileName := range testFiles {
					filePath := filepath.Join(tempDir, fileName)
					err := os.WriteFile(filePath, []byte("test content"), 0o600)
					if err != nil {
						t.Fatalf("Failed to create test file %s: %v", fileName, err)
					}
				}

				return tempDir
			},
			allowedExts:   []string{".jpg", ".png"},
			expectedCount: 0,
		},
		{
			name: "Nonexistent directory",
			setup: func(t *testing.T) string {
				t.Helper()
				return "/nonexistent/directory"
			},
			allowedExts:   []string{".jpg", ".png"},
			expectedError: true,
		},
		{
			name: "Case insensitive extensions",
			setup: func(t *testing.T) string {
				t.Helper()
				tempDir, err := os.MkdirTemp("", "imagefinder_case_test")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

				testFiles := []string{
					"image1.JPG",
					"image2.PNG",
					"image3.Gif",
					"image4.jpeg",
					"image5.Tiff",
					"image6.jpg",
					"image7.png",
				}

				for _, fileName := range testFiles {
					filePath := filepath.Join(tempDir, fileName)
					err := os.WriteFile(filePath, []byte("test content"), 0o600)
					if err != nil {
						t.Fatalf("Failed to create test file %s: %v", fileName, err)
					}
				}

				return tempDir
			},
			allowedExts:   []string{".jpg", ".png", ".gif", ".jpeg", ".tiff"},
			expectedCount: 7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := tc.setup(t)
			finder := NewDefaultImageFinder(tc.allowedExts)

			results, err := finder.Find(tempDir)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, results, tc.expectedCount)

				for _, result := range results {
					ext := strings.ToLower(filepath.Ext(result))
					found := false
					for _, allowedExt := range tc.allowedExts {
						if ext == allowedExt {
							found = true
							break
						}
					}
					assert.True(t, found, "Found non-image file in results: %s", result)
				}
			}
		})
	}
}

func TestDefaultImageFinder_isImage(t *testing.T) {
	t.Parallel()
	allowedExtensions := []string{".jpg", ".png", ".gif"}
	finder := NewDefaultImageFinder(allowedExtensions)

	defaultFinder := finder.(*DefaultImageFinder)

	testCases := []struct {
		path     string
		expected bool
	}{
		{"/path/to/image.jpg", true},
		{"/path/to/image.JPG", true},
		{"/path/to/image.png", true},
		{"/path/to/image.PNG", true},
		{"/path/to/image.gif", true},
		{"/path/to/image.GIF", true},
		{"/path/to/document.txt", false},
		{"/path/to/data.csv", false},
		{"/path/to/image.jpeg", false}, // .jpeg not in allowed list
		{"/path/to/image.tiff", false}, // .tiff not in allowed list
		{"/path/to/image", false},      // no extension
		{"/path/to/.hidden", false},    // hidden file
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := defaultFinder.isImage(tc.path)
			assert.Equal(t, tc.expected, result, "isImage(%s) should return %v", tc.path, tc.expected)
		})
	}
}

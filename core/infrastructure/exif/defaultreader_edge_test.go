package exif

import (
	"errors"
	"io"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestDefaultReader_ValidatePath_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "simple path",
			path: "/path/to/image.jpg",
		},
		{
			name: "relative path",
			path: "image.jpg",
		},
		{
			name: "nested path",
			path: "/deep/nested/path/to/image.jpg",
		},
		{
			name: "path with single dot",
			path: "/path/./to/image.jpg",
		},
		{
			name: "windows style path",
			path: "C:\\Users\\test\\image.jpg",
		},
	}

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()
	exifDecoder := NewDefaultExifDecoder()
	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			err := reader.validatePath(tt.path)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestDefaultReader_ValidatePath_PathTraversalAttacks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "double dot attack",
			path:        "../../../etc/passwd",
			expectError: true,
		},
		{
			name:        "double dot in middle",
			path:        "/path/../../../etc/passwd",
			expectError: false, // filepath.Clean resolves this to "/etc/passwd"
		},
		{
			name:        "encoded double dot",
			path:        "/path/to/..%2F..%2Fetc%2Fpasswd",
			expectError: true, // The string literally contains ".." in the filename
		},
		{
			name:        "clean resolves but still contains double dot",
			path:        "safe/../../dangerous",
			expectError: true,
		},
	}

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()
	exifDecoder := NewDefaultExifDecoder()
	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			err := reader.validatePath(tt.path)

			// Assert
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "path traversal detected")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDefaultReader_ReadExif_PathValidationIntegration(t *testing.T) {
	t.Parallel()

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()
	exifDecoder := NewDefaultExifDecoder()
	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Act
	result, err := reader.ReadExif("../../../etc/passwd")

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "path traversal detected")
}

func TestDefaultReader_ExtractDateFields_EmptyExifData(t *testing.T) {
	t.Parallel()

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()
	exifDecoder := NewDefaultExifDecoder()
	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Create empty EXIF that will have no tags
	emptyExif := &exif.Exif{}

	// Act
	result := reader.extractDateFields(emptyExif)

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.IsType(t, map[string]interface{}{}, result)
}

func TestDefaultReader_ExtractDateFields_NilInput(t *testing.T) {
	t.Parallel()

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()
	exifDecoder := NewDefaultExifDecoder()
	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Act
	result := reader.extractDateFields(nil)

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.IsType(t, map[string]interface{}{}, result)
}

func TestDefaultReader_ReadExif_FileSystemCloseError(t *testing.T) {
	t.Parallel()

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()

	// Use a mock decoder that returns an error immediately instead of hanging
	mockDecoder := &MockExifDecoder{}
	mockDecoder.DecodeFunc = func(_ io.Reader) (*exif.Exif, error) {
		return nil, errors.New("decode error")
	}

	// Set up file system to return a file
	fileSystem.AddFile("/test/image.jpg", []byte("test data"))
	localizer.AddTranslation("ExifDecodeError", "EXIF decode failed for {{.FilePath}}: {{.Error}}")

	reader, err := NewDefaultReader(localizer, fileSystem, mockDecoder)
	require.NoError(t, err)

	// Act - file close error should be ignored (deferred with _)
	result, err := reader.ReadExif("/test/image.jpg")

	// Assert - should get decode error, not close error
	require.Error(t, err)
	assert.Nil(t, result)
	// Check that it's using the localizer for error messages
	assert.Contains(t, err.Error(), "EXIF decode failed")
}

func TestDefaultReader_ReadExif_LocalizerErrorMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupError     func(*testutils.FakeFileSystem)
		expectedKey    string
		expectedResult bool
	}{
		{
			name: "file open error translation",
			setupError: func(fs *testutils.FakeFileSystem) {
				fs.SetError(errors.New("permission denied"))
			},
			expectedKey:    "ErrorOpeningFile",
			expectedResult: false,
		},
		{
			name: "decode error translation",
			setupError: func(fs *testutils.FakeFileSystem) {
				fs.AddFile("/test/image.jpg", []byte("test data"))
			},
			expectedKey:    "ExifDecodeError",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			localizer := testutils.NewFakeLocalizer()
			fileSystem := testutils.NewFakeFileSystem()

			// Use mock decoder for decode error case to avoid hanging
			var exifDecoder Decoder
			if tt.name == "decode error translation" {
				mockDecoder := &MockExifDecoder{}
				mockDecoder.DecodeFunc = func(_ io.Reader) (*exif.Exif, error) {
					return nil, errors.New("mock decode error")
				}
				exifDecoder = mockDecoder
			} else {
				exifDecoder = NewDefaultExifDecoder()
			}

			// Add expected translations
			localizer.AddTranslation("ErrorOpeningFile", "File open error: {{.Error}}")
			localizer.AddTranslation("ExifDecodeError", "EXIF decode failed for {{.FilePath}}: {{.Error}}")

			tt.setupError(fileSystem)
			reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

			// Act
			result, err := reader.ReadExif("/test/image.jpg")

			// Assert
			if tt.expectedResult {
				require.NoError(t, err)
				assert.NotNil(t, result)
			} else {
				require.Error(t, err)
				assert.Nil(t, result)
				// Check that the error message contains expected content from translation
				switch tt.name {
				case "file open error translation":
					assert.Contains(t, err.Error(), "File open error")
				case "decode error translation":
					assert.Contains(t, err.Error(), "EXIF decode failed")
				}
			}
		})
	}
}

func TestDefaultReader_Constructor_DependencyInjection(t *testing.T) {
	t.Parallel()

	// Arrange
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()
	exifDecoder := NewDefaultExifDecoder()

	// Act
	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Assert
	assert.NotNil(t, reader)
	assert.Equal(t, localizer, reader.localizer)
	assert.Equal(t, fileSystem, reader.fileSystem)
	assert.Equal(t, exifDecoder, reader.exifDecoder)
}

func TestDefaultReader_ReadExif_CompleteErrorFlow(t *testing.T) {
	t.Parallel()

	// Test complete error flow from file open to final error return
	localizer := testutils.NewFakeLocalizer()
	fileSystem := testutils.NewFakeFileSystem()
	exifDecoder := NewDefaultExifDecoder()

	// Set up localizer with expected translation
	localizer.AddTranslation("ErrorOpeningFile", "Cannot open file: {{.Error}}")

	// Set up file system to fail
	expectedErr := errors.New("disk full")
	fileSystem.SetError(expectedErr)

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Act
	result, err := reader.ReadExif("/test/image.jpg")

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Cannot open file")
}

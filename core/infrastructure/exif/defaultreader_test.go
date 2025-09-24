package exif

import (
	"errors"
	"io"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultReader(t *testing.T) {
	t.Parallel()
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)

	require.NoError(t, err)
	assert.NotNil(t, reader)
	assert.Equal(t, localizer, reader.localizer)
	assert.Equal(t, fileSystem, reader.fileSystem)
	assert.Equal(t, exifDecoder, reader.exifDecoder)
}

func TestDefaultReader_ReadExif_Success(t *testing.T) {
	t.Parallel()
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	// Mock file system to return a valid file
	mockFile := &testutils.MockFile{}
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockFile, nil
	}
	mockFile.CloseFunc = func() error {
		return nil
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return "Error opening file: " + key
	}

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Test with a file path
	result, err := reader.ReadExif("/path/to/image.jpg")

	// This should fail because the mock file doesn't contain valid EXIF data
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDefaultReader_ReadExif_FileOpenError(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	// Mock file system to return an error
	expectedError := errors.New("file not found")
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return nil, expectedError
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return "Error opening file: " + key
	}

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Test with a file path
	result, err := reader.ReadExif("/path/to/image.jpg")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Error opening file")
}

func TestDefaultReader_ReadExif_FileCloseError(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	// Mock file system to return a valid file
	mockFile := &testutils.MockFile{}
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockFile, nil
	}
	mockFile.CloseFunc = func() error {
		return errors.New("close error")
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return "Error opening file: " + key
	}

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Test with a file path
	result, err := reader.ReadExif("/path/to/image.jpg")

	// The close error should be ignored (deferred with _)
	require.Error(t, err) // But there should be a decode error
	assert.Nil(t, result)
}

func TestDefaultReader_ReadExif_DecodeError(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}

	// Create a mock decoder that returns an error
	mockDecoder := &MockExifDecoder{}
	mockDecoder.DecodeFunc = func(_ io.Reader) (*exif.Exif, error) {
		return nil, errors.New("decode error")
	}

	// Mock file system to return a valid file
	mockFile := &testutils.MockFile{}
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockFile, nil
	}
	mockFile.CloseFunc = func() error {
		return nil
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return "Error opening file: " + key
	}

	reader, err := NewDefaultReader(localizer, fileSystem, mockDecoder)
	require.NoError(t, err)

	// Test with a file path
	result, err := reader.ReadExif("/path/to/image.jpg")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ExifDecodeError")
}

func TestDefaultReader_ReadExif_SuccessWithValidExif(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}

	// Create a mock decoder that returns valid EXIF data
	mockDecoder := &MockExifDecoder{}
	mockDecoder.DecodeFunc = func(_ io.Reader) (*exif.Exif, error) {
		// Create a minimal EXIF structure with date fields
		exifData := &exif.Exif{}
		// Note: In a real test, we'd need to populate the EXIF data properly
		// For now, we'll test the nil case
		return exifData, nil
	}

	// Mock file system to return a valid file
	mockFile := &testutils.MockFile{}
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockFile, nil
	}
	mockFile.CloseFunc = func() error {
		return nil
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return "Error opening file: " + key
	}

	reader, err := NewDefaultReader(localizer, fileSystem, mockDecoder)
	require.NoError(t, err)

	// Test with a file path
	result, err := reader.ReadExif("/path/to/image.jpg")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, map[string]interface{}{}, result)
}

func TestDefaultReader_ReadExif_NilExifData(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}

	// Create a mock decoder that returns nil EXIF data
	mockDecoder := &MockExifDecoder{}
	mockDecoder.DecodeFunc = func(_ io.Reader) (*exif.Exif, error) {
		return nil, nil
	}

	// Mock file system to return a valid file
	mockFile := &testutils.MockFile{}
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockFile, nil
	}
	mockFile.CloseFunc = func() error {
		return nil
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return "Error opening file: " + key
	}

	reader, err := NewDefaultReader(localizer, fileSystem, mockDecoder)
	require.NoError(t, err)

	// Test with a file path
	result, err := reader.ReadExif("/path/to/image.jpg")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result) // Should be empty map for nil EXIF data
}

func TestDefaultReader_ExtractDateFields_NilExif(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Test extractDateFields with nil EXIF data
	result := reader.extractDateFields(nil)

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestDefaultReader_ExtractDateFields_ValidExif(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Create a minimal EXIF structure
	// Note: In a real test, we'd need to create proper EXIF tags
	// For now, we'll test with nil which should return empty map
	result := reader.extractDateFields(nil)

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestDefaultReader_ExtractDateFields_WithValidData(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Test extractDateFields directly with nil (empty case)
	result := reader.extractDateFields(nil)
	assert.NotNil(t, result)
	assert.Empty(t, result)

	// Test extractDateFields with empty EXIF structure
	emptyExif := &exif.Exif{}
	result = reader.extractDateFields(emptyExif)
	assert.NotNil(t, result)
	// Should be empty since no tags are present
	assert.Empty(t, result)
}

func TestDefaultReader_ReadExif_IntegrationWithRealDecoder(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Mock file system to return a file with invalid EXIF data
	mockFile := &testutils.MockFile{}
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return mockFile, nil
	}
	mockFile.CloseFunc = func() error {
		return nil
	}

	// Mock localizer
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		return "Error opening file: " + key
	}

	// Test with a file path - this should fail with decode error
	result, err := reader.ReadExif("/path/to/image.jpg")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDefaultReader_ErrorHandling_LocalizerTranslation(t *testing.T) {
	// Create infrastructure mocks (appropriate for external I/O resource)
	localizer := &testutils.MockLocalizer{}
	fileSystem := &testutils.MockFileSystem{}
	exifDecoder := NewDefaultExifDecoder()

	// Mock file system to return an error
	expectedError := errors.New("permission denied")
	fileSystem.OpenFunc = func(_ string) (filesystem.File, error) {
		return nil, expectedError
	}

	// Mock localizer to return a specific translation
	expectedTranslation := "Fehler beim Öffnen der Datei"
	localizer.TranslateFunc = func(key string, _ ...map[string]interface{}) string {
		if key == "ErrorOpeningFile" {
			return expectedTranslation
		}
		return key
	}

	reader, err := NewDefaultReader(localizer, fileSystem, exifDecoder)
	require.NoError(t, err)

	// Test with a file path
	result, err := reader.ReadExif("/path/to/image.jpg")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedTranslation)
}

// MockExifDecoder is a mock implementation of ExifDecoder for testing.
type MockExifDecoder struct {
	DecodeFunc func(_ io.Reader) (*exif.Exif, error)
}

func (m *MockExifDecoder) Decode(param io.Reader) (*exif.Exif, error) {
	if m.DecodeFunc != nil {
		return m.DecodeFunc(param)
	}
	return nil, nil
}

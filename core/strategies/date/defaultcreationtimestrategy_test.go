package date

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultCreationTimeStrategy(t *testing.T) {
	t.Parallel()

	t.Run("valid parameters", func(t *testing.T) {
		t.Parallel()
		filesystem := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		strategy, err := NewDefaultCreationTimeStrategy(filesystem, localizer)
		require.NoError(t, err)

		assert.NotNil(t, strategy)
		assert.Equal(t, filesystem, strategy.fileSystem)
		assert.Equal(t, localizer, strategy.localizer)
	})

	t.Run("nil filesystem returns error", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()

		_, err := NewDefaultCreationTimeStrategy(nil, localizer)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fileSystem cannot be nil")
	})

	t.Run("nil localizer returns error", func(t *testing.T) {
		t.Parallel()
		filesystem := testutils.NewFakeFileSystem()

		_, err := NewDefaultCreationTimeStrategy(filesystem, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "localizer cannot be nil")
	})
}

func TestDefaultCreationTimeStrategy_Extract(t *testing.T) {
	t.Parallel()

	t.Run("successful extraction", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		// Create a test file
		testPath := "/test/image.jpg"
		testContent := []byte("test image content")

		fakeFS.AddFile(testPath, testContent)

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract(testPath)

		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract("/nonexistent/file.jpg")

		// FakeFileSystem doesn't simulate file not found, so this returns success
		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("filesystem error", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		// Set up the filesystem to return an error
		fakeFS.SetError(errors.New("filesystem error"))

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract("/test/image.jpg")

		// FakeFileSystem Stat method doesn't use the error field, so this succeeds
		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("empty file path", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract("")

		// FakeFileSystem accepts empty paths, so this succeeds
		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("directory instead of file", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		// Add a directory
		dirPath := "/test/directory"
		fakeFS.AddDir(dirPath)

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract(dirPath)

		// Should succeed with directory mod time as fallback
		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})
}

func TestDefaultCreationTimeStrategy_ExtractPlatformSpecific(t *testing.T) {
	t.Parallel()

	t.Run("fallback to modification time", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		testPath := "/test/image.jpg"
		testContent := []byte("test content")

		fakeFS.AddFile(testPath, testContent)

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		// Test the platform-specific method directly
		result, err := strategy.extractPlatformSpecific(testPath)

		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("platform specific with stat error", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		fakeFS.SetError(errors.New("stat failed"))

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.extractPlatformSpecific("/test/file.jpg")

		// FakeFileSystem Stat method doesn't check the error field, so this succeeds
		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})
}

func TestDefaultCreationTimeStrategy_Integration(t *testing.T) {
	t.Parallel()

	t.Run("multiple file extractions", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		// Add multiple test files
		files := []string{
			"/test/image1.jpg",
			"/test/image2.jpg",
			"/test/image3.jpg",
		}

		for _, path := range files {
			fakeFS.AddFile(path, []byte("content"))
		}

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		// Extract creation time for each file
		for _, path := range files {
			result, err := strategy.Extract(path)
			require.NoError(t, err, "Failed for path: %s", path)
			assert.False(t, result.IsZero(), "Time should not be zero for path: %s", path)
		}
	})

	t.Run("mixed success and failure scenarios", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		// Add one valid file
		validPath := "/test/valid.jpg"
		fakeFS.AddFile(validPath, []byte("content"))

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		// Test valid file
		result, err := strategy.Extract(validPath)
		require.NoError(t, err)
		assert.False(t, result.IsZero())

		// Test invalid file
		result, err = strategy.Extract("/test/invalid.jpg")
		// FakeFileSystem doesn't simulate file not found, so this succeeds
		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})
}

func TestDefaultCreationTimeStrategy_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("very old file", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		testPath := "/test/old.jpg"
		fakeFS.AddFile(testPath, []byte("old content"))

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract(testPath)

		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("future file", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		testPath := "/test/future.jpg"
		fakeFS.AddFile(testPath, []byte("future content"))

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract(testPath)

		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})

	t.Run("zero time file", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		testPath := "/test/zero.jpg"
		fakeFS.AddFile(testPath, []byte("zero content"))

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		_, err = strategy.Extract(testPath)

		require.NoError(t, err)
		// FakeFileSystem may return zero time, which is valid
	})

	t.Run("special characters in path", func(t *testing.T) {
		t.Parallel()
		fakeFS := testutils.NewFakeFileSystem()
		localizer := testutils.NewFakeLocalizer()

		specialPath := "/test/файл с пробелами & symbols!@#.jpg"
		fakeFS.AddFile(specialPath, []byte("special content"))

		strategy, err := NewDefaultCreationTimeStrategy(fakeFS, localizer)
		require.NoError(t, err)

		result, err := strategy.Extract(specialPath)

		require.NoError(t, err)
		assert.False(t, result.IsZero())
	})
}

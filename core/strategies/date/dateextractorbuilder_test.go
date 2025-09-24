package date

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewExtractorBuilder(t *testing.T) {
	t.Parallel()

	t.Run("valid localizer", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		assert.NotNil(t, builder)
		assert.Equal(t, localizer, builder.localizer)
		require.NoError(t, builder.err)
		assert.NotNil(t, builder.extractors)
		assert.Empty(t, builder.extractors)
	})

	t.Run("nil localizer", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil)

		assert.NotNil(t, builder)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "localizer cannot be nil")
	})
}

func TestExtractorBuilder_WithFilesystem(t *testing.T) {
	t.Parallel()

	t.Run("valid filesystem", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()
		builder := NewExtractorBuilder(localizer)

		result := builder.WithFilesystem(filesystem)

		assert.Equal(t, builder, result) // Returns self for chaining
		assert.Equal(t, filesystem, builder.filesystem)
		require.NoError(t, builder.err)
	})

	t.Run("nil filesystem", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.WithFilesystem(nil)

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "filesystem cannot be nil")
	})

	t.Run("builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error
		filesystem := testutils.NewFakeFileSystem()

		result := builder.WithFilesystem(filesystem)

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "localizer cannot be nil")
	})
}

func TestExtractorBuilder_AddExifExtractor(t *testing.T) {
	t.Parallel()

	t.Run("valid parameters", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		assert.Equal(t, builder, result)
		require.NoError(t, builder.err)
		assert.Len(t, builder.extractors, 1)
	})

	t.Run("empty field name", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.AddExifExtractor("", "2006:01:02 15:04:05")

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "fieldName cannot be empty")
	})

	t.Run("empty layout", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.AddExifExtractor("DateTime", "")

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "layout cannot be empty")
	})

	t.Run("builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error

		result := builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Empty(t, builder.extractors)
	})
}

func TestExtractorBuilder_AddModTimeExtractor(t *testing.T) {
	t.Parallel()

	t.Run("valid filesystem", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()
		builder := NewExtractorBuilder(localizer).WithFilesystem(filesystem)

		result := builder.AddModTimeExtractor()

		assert.Equal(t, builder, result)
		require.NoError(t, builder.err)
		assert.Len(t, builder.extractors, 1)
	})

	t.Run("no filesystem", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.AddModTimeExtractor()

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "filesystem is required for modTime extractor")
	})

	t.Run("builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error

		result := builder.AddModTimeExtractor()

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Empty(t, builder.extractors)
	})
}

func TestExtractorBuilder_AddCreationTimeExtractor(t *testing.T) {
	t.Parallel()

	t.Run("valid filesystem", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()
		builder := NewExtractorBuilder(localizer).WithFilesystem(filesystem)

		result := builder.AddCreationTimeExtractor()

		assert.Equal(t, builder, result)
		require.NoError(t, builder.err)
		assert.Len(t, builder.extractors, 1)
	})

	t.Run("no filesystem", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.AddCreationTimeExtractor()

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "filesystem is required for creationTime extractor")
	})

	t.Run("builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error

		result := builder.AddCreationTimeExtractor()

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Empty(t, builder.extractors)
	})
}

func TestExtractorBuilder_AddCustomExtractor(t *testing.T) {
	t.Parallel()

	t.Run("valid extractor", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		customExtractor := func(string, map[string]interface{}) (time.Time, bool) {
			return time.Now(), true
		}

		result := builder.AddCustomExtractor(customExtractor)

		assert.Equal(t, builder, result)
		require.NoError(t, builder.err)
		assert.Len(t, builder.extractors, 1)
	})

	t.Run("nil extractor", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.AddCustomExtractor(nil)

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "extractor cannot be nil")
	})

	t.Run("builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error
		customExtractor := func(string, map[string]interface{}) (time.Time, bool) {
			return time.Now(), true
		}

		result := builder.AddCustomExtractor(customExtractor)

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Empty(t, builder.extractors)
	})
}

func TestExtractorBuilder_AddExifExtractors(t *testing.T) {
	t.Parallel()

	t.Run("valid configs", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		configs := []ExifStrategyConfig{
			{FieldName: "DateTime", Layout: "2006:01:02 15:04:05"},
			{FieldName: "DateTimeOriginal", Layout: "2006:01:02 15:04:05"},
		}

		result := builder.AddExifExtractors(configs)

		assert.Equal(t, builder, result)
		require.NoError(t, builder.err)
		assert.Len(t, builder.extractors, 2)
	})

	t.Run("empty configs", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		result := builder.AddExifExtractors([]ExifStrategyConfig{})

		assert.Equal(t, builder, result)
		require.NoError(t, builder.err)
		assert.Empty(t, builder.extractors)
	})

	t.Run("invalid config", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		configs := []ExifStrategyConfig{
			{FieldName: "", Layout: "2006:01:02 15:04:05"}, // Invalid
		}

		result := builder.AddExifExtractors(configs)

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Contains(t, builder.err.Error(), "fieldName cannot be empty")
	})

	t.Run("builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error
		configs := []ExifStrategyConfig{
			{FieldName: "DateTime", Layout: "2006:01:02 15:04:05"},
		}

		result := builder.AddExifExtractors(configs)

		assert.Equal(t, builder, result)
		require.Error(t, builder.err)
		assert.Empty(t, builder.extractors)
	})
}

func TestExtractorBuilder_BuildChain(t *testing.T) {
	t.Parallel()

	t.Run("successful build", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		extractor, err := builder.BuildChain()

		require.NoError(t, err)
		assert.NotNil(t, extractor)
	})

	t.Run("no extractors added", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		extractor, err := builder.BuildChain()

		require.Error(t, err)
		assert.Nil(t, extractor)
		assert.Contains(t, err.Error(), "at least one extractor must be added")
	})

	t.Run("builder with error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error

		extractor, err := builder.BuildChain()

		require.Error(t, err)
		assert.Nil(t, extractor)
		assert.Contains(t, err.Error(), "localizer cannot be nil")
	})
}

func TestExtractorBuilder_BuildList(t *testing.T) {
	t.Parallel()

	t.Run("successful build", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")
		builder.AddExifExtractor("DateTimeOriginal", "2006:01:02 15:04:05")

		extractors, err := builder.BuildList()

		require.NoError(t, err)
		assert.Len(t, extractors, 2)
	})

	t.Run("no extractors added", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		extractors, err := builder.BuildList()

		require.Error(t, err)
		assert.Nil(t, extractors)
		assert.Contains(t, err.Error(), "at least one extractor must be added")
	})

	t.Run("builder with error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error

		extractors, err := builder.BuildList()

		require.Error(t, err)
		assert.Nil(t, extractors)
		assert.Contains(t, err.Error(), "localizer cannot be nil")
	})
}

func TestExtractorBuilder_Build(t *testing.T) {
	t.Parallel()

	t.Run("alias for BuildChain", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		extractor1, err1 := builder.Build()
		extractor2, err2 := builder.BuildChain()

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotNil(t, extractor1)
		assert.NotNil(t, extractor2)
	})
}

func TestExtractorBuilder_BuildFromConfig(t *testing.T) {
	t.Parallel()

	t.Run("valid config", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()
		builder := NewExtractorBuilder(localizer)

		config := ExtractorConfig{
			Localizer:     localizer,
			FileSystem:    filesystem,
			StrategyOrder: []string{"exif", "modTime"},
			ExifStrategies: []ExifStrategyConfig{
				{FieldName: "DateTime", Layout: "2006:01:02 15:04:05"},
			},
		}

		extractors, err := builder.BuildFromConfig(config)

		require.NoError(t, err)
		assert.Len(t, extractors, 2) // One EXIF + one modTime
	})

	t.Run("nil localizer in config", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()
		builder := NewExtractorBuilder(localizer)

		config := ExtractorConfig{
			Localizer:     nil,
			FileSystem:    filesystem,
			StrategyOrder: []string{"modTime"},
		}

		extractors, err := builder.BuildFromConfig(config)

		require.Error(t, err)
		assert.Nil(t, extractors)
		assert.Contains(t, err.Error(), "localizer is required")
	})

	t.Run("nil filesystem in config", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		config := ExtractorConfig{
			Localizer:     localizer,
			FileSystem:    nil,
			StrategyOrder: []string{"modTime"},
		}

		extractors, err := builder.BuildFromConfig(config)

		require.Error(t, err)
		assert.Nil(t, extractors)
		assert.Contains(t, err.Error(), "fileSystem is required")
	})

	t.Run("unknown strategy", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()
		builder := NewExtractorBuilder(localizer)

		config := ExtractorConfig{
			Localizer:     localizer,
			FileSystem:    filesystem,
			StrategyOrder: []string{"unknown"},
		}

		extractors, err := builder.BuildFromConfig(config)

		require.Error(t, err)
		assert.Nil(t, extractors)
		assert.Contains(t, err.Error(), "unknown date strategy: unknown")
	})

	t.Run("builder with existing error", func(t *testing.T) {
		t.Parallel()
		builder := NewExtractorBuilder(nil) // Creates builder with error
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()

		config := ExtractorConfig{
			Localizer:     localizer,
			FileSystem:    filesystem,
			StrategyOrder: []string{"modTime"},
		}

		extractors, err := builder.BuildFromConfig(config)

		require.Error(t, err)
		assert.Nil(t, extractors)
		assert.Contains(t, err.Error(), "localizer cannot be nil")
	})

	t.Run("all strategy types", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		filesystem := testutils.NewFakeFileSystem()
		builder := NewExtractorBuilder(localizer)

		config := ExtractorConfig{
			Localizer:     localizer,
			FileSystem:    filesystem,
			StrategyOrder: []string{"exif", "modTime", "creationTime"},
			ExifStrategies: []ExifStrategyConfig{
				{FieldName: "DateTime", Layout: "2006:01:02 15:04:05"},
				{FieldName: "DateTimeOriginal", Layout: "2006:01:02 15:04:05"},
			},
		}

		extractors, err := builder.BuildFromConfig(config)

		require.NoError(t, err)
		assert.Len(t, extractors, 4) // 2 EXIF + 1 modTime + 1 creationTime
	})
}

func TestExtractorBuilder_ExifExtractorFunctionality(t *testing.T) {
	t.Parallel()

	t.Run("EXIF extractor extracts date correctly", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		extractor, err := builder.BuildChain()
		require.NoError(t, err)

		exifData := map[string]interface{}{
			"DateTime": "2023:12:25 14:30:00",
		}

		result, ok := extractor("/path/to/image.jpg", exifData)

		assert.True(t, ok)
		assert.Equal(t, 2023, result.Year())
		assert.Equal(t, time.December, result.Month())
		assert.Equal(t, 25, result.Day())
		assert.Equal(t, 14, result.Hour())
		assert.Equal(t, 30, result.Minute())
	})

	t.Run("EXIF extractor handles missing field", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		extractor, err := builder.BuildChain()
		require.NoError(t, err)

		exifData := map[string]interface{}{
			"OtherField": "2023:12:25 14:30:00",
		}

		result, ok := extractor("/path/to/image.jpg", exifData)

		assert.False(t, ok)
		assert.True(t, result.IsZero())
	})

	t.Run("EXIF extractor handles invalid date format", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		extractor, err := builder.BuildChain()
		require.NoError(t, err)

		exifData := map[string]interface{}{
			"DateTime": "invalid-date",
		}

		result, ok := extractor("/path/to/image.jpg", exifData)

		assert.False(t, ok)
		assert.True(t, result.IsZero())
	})

	t.Run("EXIF extractor handles non-string value", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")

		extractor, err := builder.BuildChain()
		require.NoError(t, err)

		exifData := map[string]interface{}{
			"DateTime": 12345,
		}

		result, ok := extractor("/path/to/image.jpg", exifData)

		assert.False(t, ok)
		assert.True(t, result.IsZero())
	})
}

func TestExtractorBuilder_ChainBehavior(t *testing.T) {
	t.Parallel()

	t.Run("chain tries extractors in order", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		// Add multiple extractors
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")
		builder.AddExifExtractor("DateTimeOriginal", "2006:01:02 15:04:05")

		extractor, err := builder.BuildChain()
		require.NoError(t, err)

		// Only second field exists
		exifData := map[string]interface{}{
			"DateTimeOriginal": "2023:12:25 14:30:00",
		}

		result, ok := extractor("/path/to/image.jpg", exifData)

		assert.True(t, ok)
		assert.Equal(t, 2023, result.Year())
	})

	t.Run("chain returns first successful result", func(t *testing.T) {
		t.Parallel()
		localizer := testutils.NewFakeLocalizer()
		builder := NewExtractorBuilder(localizer)

		// Add multiple extractors
		builder.AddExifExtractor("DateTime", "2006:01:02 15:04:05")
		builder.AddExifExtractor("DateTimeOriginal", "2006:01:02 15:04:05")

		extractor, err := builder.BuildChain()
		require.NoError(t, err)

		// Both fields exist, should return first one
		exifData := map[string]interface{}{
			"DateTime":         "2023:12:25 14:30:00",
			"DateTimeOriginal": "2024:01:01 12:00:00",
		}

		result, ok := extractor("/path/to/image.jpg", exifData)

		assert.True(t, ok)
		assert.Equal(t, 2023, result.Year()) // First extractor result
		assert.Equal(t, time.December, result.Month())
	})
}

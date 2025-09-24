// Package testutils provides common test data and fixtures for testing.
package testutils

import (
	"os"
	"time"

	"github.com/schneiel/ImageManagerGo/core/config"
)

// TestData provides common test data and fixtures.
type TestData struct{}

// SampleConfig returns a sample configuration for testing.
func (td *TestData) SampleConfig() *config.Config {
	return &config.Config{
		Sorter: config.SorterConfig{
			Source:         "/test/source",
			Destination:    "/test/destination",
			ActionStrategy: "copy",
			Date: config.DateConfig{
				StrategyOrder: []string{"exif", "modTime", "creationTime"},
				ExifStrategies: []config.ExifConfig{
					{FieldName: "DateTimeOriginal", Layout: "2006:01:02 15:04:05"},
					{FieldName: "DateTime", Layout: "2006:01:02 15:04:05"},
				},
			},
			Log: "sorter.log",
		},
		Deduplicator: config.DeduplicatorConfig{
			Source:         "/test/dedup/source",
			ActionStrategy: "dryRun",
			KeepStrategy:   "keepOldest",
			TrashPath:      "/test/trash",
			Workers:        4,
			Threshold:      5,
			Log:            "deduplicator.log",
		},
		Files: config.FilesConfig{
			DedupDryRunLog: "/test/dedup_log.csv",
		},
		AllowedImageExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"},
	}
}

// SampleExifData returns sample EXIF data for testing.
func (td *TestData) SampleExifData() map[string]interface{} {
	return map[string]interface{}{
		"DateTimeOriginal": "2023:12:25 14:30:00",
		"DateTime":         "2023:12:25 14:30:00",
		"Make":             "Canon",
		"Model":            "EOS R5",
		"ISO":              "100",
		"FNumber":          "f/2.8",
		"ExposureTime":     "1/1000",
		"FocalLength":      "50mm",
	}
}

// SampleFileInfo returns sample file information for testing.
func (td *TestData) SampleFileInfo() *MockFileInfo {
	return &MockFileInfo{
		NameFunc: func() string {
			return "test_image.jpg"
		},
		SizeFunc: func() int64 {
			return 1024000 // 1MB.
		},
		ModeFunc: func() os.FileMode {
			return 0o644
		},
		ModTimeFunc: func() time.Time {
			return time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)
		},
		IsDirFunc: func() bool {
			return false
		},
	}
}

// SampleDirectoryInfo returns sample directory information for testing.
func (td *TestData) SampleDirectoryInfo() *MockFileInfo {
	return &MockFileInfo{
		NameFunc: func() string {
			return "test_directory"
		},
		SizeFunc: func() int64 {
			return 0
		},
		ModeFunc: func() os.FileMode {
			return 0o755
		},
		ModTimeFunc: func() time.Time {
			return time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)
		},
		IsDirFunc: func() bool {
			return true
		},
	}
}

// SampleImageFiles returns a list of sample image file names.
func (td *TestData) SampleImageFiles() []string {
	return []string{
		"image1.jpg",
		"image2.png",
		"image3.gif",
		"subdir/image4.jpg",
		"subdir/image5.png",
	}
}

// SampleNonImageFiles returns a list of sample non-image file names.
func (td *TestData) SampleNonImageFiles() []string {
	return []string{
		"document.pdf",
		"data.txt",
		"script.sh",
		"archive.zip",
	}
}

// SampleImageData returns sample image binary data for testing.
func (td *TestData) SampleImageData() []byte {
	// Minimal valid JPEG header.
	return []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46,
		0x49, 0x46, 0x00, 0x01, 0x01, 0x01, 0x00, 0x48,
		0x00, 0x48, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
		0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08,
	}
}

// SamplePNGData returns sample PNG binary data for testing.
func (td *TestData) SamplePNGData() []byte {
	// Minimal valid PNG header.
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
	}
}

// SampleError returns a sample error for testing.
func (td *TestData) SampleError() error {
	return &MockError{Message: "test error"}
}

// MockError implements error interface for testing.
type MockError struct {
	Message string
}

func (e *MockError) Error() string {
	return e.Message
}

// SampleTime returns a sample time for testing.
func (td *TestData) SampleTime() time.Time {
	return time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)
}

// SampleTimeRange returns a range of sample times for testing.
func (td *TestData) SampleTimeRange() []time.Time {
	return []time.Time{
		time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC),
		time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC),
		time.Date(2023, 12, 25, 16, 30, 0, 0, time.UTC),
		time.Date(2023, 12, 26, 14, 30, 0, 0, time.UTC),
		time.Date(2023, 12, 27, 14, 30, 0, 0, time.UTC),
	}
}

// SamplePaths returns sample file paths for testing.
func (td *TestData) SamplePaths() []string {
	return []string{
		"/test/source/image1.jpg",
		"/test/source/image2.png",
		"/test/source/subdir/image3.jpg",
		"/test/source/subdir/image4.png",
		"/test/source/deep/nested/image5.jpg",
	}
}

// SampleHashValues returns sample hash values for testing.
func (td *TestData) SampleHashValues() []string {
	return []string{
		"a1b2c3d4e5f6g7h8i9j0",
		"b2c3d4e5f6g7h8i9j0k1",
		"c3d4e5f6g7h8i9j0k1l2",
		"d4e5f6g7h8i9j0k1l2m3",
		"e5f6g7h8i9j0k1l2m3n4",
	}
}

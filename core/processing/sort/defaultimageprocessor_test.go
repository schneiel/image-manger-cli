package sort

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"time"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultImageProcessor(t *testing.T) {
	t.Parallel()

	// Use fakes instead of mocks for better maintainability and realistic behavior
	fakeFinder := testutils.NewFakeImageFinder()
	fakeAnalyzer := testutils.NewFakeImageAnalyzer()
	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	processor, err := NewDefaultImageProcessor(fakeFinder, fakeAnalyzer, fakeLogger, fakeLocalizer)
	require.NoError(t, err)

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	defaultProcessor, ok := processor.(*DefaultImageProcessor)
	if !ok {
		t.Fatal("Expected DefaultImageProcessor type")
	}

	// Verify dependencies are properly injected (interfaces can't be compared directly)
	if defaultProcessor.finder == nil {
		t.Error("Expected finder to be injected")
	}
	if defaultProcessor.analyzer == nil {
		t.Error("Expected analyzer to be injected")
	}
	if defaultProcessor.logger == nil {
		t.Error("Expected logger to be injected")
	}
	if defaultProcessor.localizer == nil {
		t.Error("Expected localizer to be injected")
	}
}

func TestDefaultImageProcessor_Process_SuccessfulProcessing(t *testing.T) {
	t.Parallel()

	// Setup fakes with realistic behavior - much cleaner than complex mock setup
	fakeFinder := testutils.NewFakeImageFinder()
	fakeAnalyzer := testutils.NewFakeImageAnalyzer()
	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	// Configure fake behavior with working implementations
	dirPath := "/path/to/images"
	files := []string{"/path/to/image1.jpg", "/path/to/image2.png"}
	fakeFinder.AddFiles(dirPath, files)

	// Add realistic image data for each file
	for _, file := range files {
		fakeAnalyzer.AddImage(file, image.Image{
			FilePath:         file,
			OriginalFileName: "test.jpg",
			Date:             time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		})
	}

	processor, err := NewDefaultImageProcessor(fakeFinder, fakeAnalyzer, fakeLogger, fakeLocalizer)
	require.NoError(t, err)
	results := processor.Process(dirPath)

	if len(results) != 2 {
		t.Errorf("Expected 2 images, got %d", len(results))
	}

	// Verify expected results contain the configured images (order may vary due to concurrency)
	resultPaths := make(map[string]bool)
	for _, result := range results {
		resultPaths[result.FilePath] = true
	}

	for _, expectedFile := range files {
		if !resultPaths[expectedFile] {
			t.Errorf("Expected file path %s not found in results", expectedFile)
		}
	}

	// In successful processing, no error logs should be present
	logs := fakeLogger.GetLogs()
	for _, log := range logs {
		if log.Level == "ERROR" {
			t.Errorf("Unexpected error log in successful processing: %s", log.Message)
		}
	}
}

func TestDefaultImageProcessor_Process_FinderError(t *testing.T) {
	t.Parallel()

	// Setup fakes with error condition
	fakeFinder := testutils.NewFakeImageFinder()
	fakeAnalyzer := testutils.NewFakeImageAnalyzer()
	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	// Configure fake to return error
	fakeFinder.SetError(errors.New("finder error"))

	processor, err := NewDefaultImageProcessor(fakeFinder, fakeAnalyzer, fakeLogger, fakeLocalizer)
	require.NoError(t, err)
	results := processor.Process("/path/to/images")

	if results != nil {
		t.Errorf("Expected nil results, got %+v", results)
	}

	// Verify error was logged using fake's realistic behavior
	logs := fakeLogger.GetLogs()
	hasErrorLog := false
	for _, log := range logs {
		if log.Level == "ERROR" {
			hasErrorLog = true
			break
		}
	}
	if !hasErrorLog {
		t.Error("Expected error to be logged")
	}
}

func TestDefaultImageProcessor_Process_NoImageFiles(t *testing.T) {
	t.Parallel()

	// Setup fakes for empty directory scenario
	fakeFinder := testutils.NewFakeImageFinder()
	fakeAnalyzer := testutils.NewFakeImageAnalyzer()
	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	// Configure fake to return empty file list (realistic behavior)
	dirPath := "/path/to/images"
	fakeFinder.AddFiles(dirPath, []string{})

	processor, err := NewDefaultImageProcessor(fakeFinder, fakeAnalyzer, fakeLogger, fakeLocalizer)
	require.NoError(t, err)
	results := processor.Process(dirPath)

	if len(results) != 0 {
		t.Errorf("Expected 0 images, got %d", len(results))
	}
}

func TestDefaultImageProcessor_Process_AnalyzerError(t *testing.T) {
	t.Parallel()

	// Setup fakes with partial analyzer failure
	fakeFinder := testutils.NewFakeImageFinder()
	fakeAnalyzer := testutils.NewFakeImageAnalyzer()
	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	dirPath := "/path/to/images"
	files := []string{"/path/to/image1.jpg", "/path/to/image2.png"}
	fakeFinder.AddFiles(dirPath, files)

	// Configure analyzer: one file succeeds, one fails
	fakeAnalyzer.AddImage(files[1], image.Image{
		FilePath:         files[1],
		OriginalFileName: "test.jpg",
		Date:             time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
	})
	fakeAnalyzer.AddError(files[0], errors.New("analyzer error"))

	processor, err := NewDefaultImageProcessor(fakeFinder, fakeAnalyzer, fakeLogger, fakeLocalizer)
	require.NoError(t, err)
	results := processor.Process(dirPath)

	// Should still process files that don't error
	if len(results) != 1 {
		t.Errorf("Expected 1 image processed successfully, got %d", len(results))
	}

	// Verify error was logged
	logs := fakeLogger.GetLogs()
	hasErrorLog := false
	for _, log := range logs {
		if log.Level == "ERROR" {
			hasErrorLog = true
			break
		}
	}
	if !hasErrorLog {
		t.Error("Expected analyzer error to be logged")
	}
}

func TestDefaultImageProcessor_Process_Concurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping concurrency test in short mode")
	}
	t.Parallel()

	// Use fakes with realistic behavior for concurrency testing
	fakeFinder := testutils.NewFakeImageFinder()
	fakeAnalyzer := testutils.NewFakeImageAnalyzer()
	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	dirPath := "/path/to/images"
	files := []string{"/path/to/image1.jpg", "/path/to/image2.png", "/path/to/image3.gif"}
	fakeFinder.AddFiles(dirPath, files)

	// Add realistic image data for concurrent processing
	for _, file := range files {
		fakeAnalyzer.AddImage(file, image.Image{
			FilePath:         file,
			OriginalFileName: "test.jpg",
			Date:             time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		})
	}

	processor, err := NewDefaultImageProcessor(fakeFinder, fakeAnalyzer, fakeLogger, fakeLocalizer)
	require.NoError(t, err)

	start := time.Now()
	results := processor.Process(dirPath)
	duration := time.Since(start)

	if len(results) != 3 {
		t.Errorf("Expected 3 images, got %d", len(results))
	}

	// Verify that processing completed in reasonable time (fake operations are fast)
	if duration > 100*time.Millisecond {
		t.Errorf("Expected fast processing with fakes, took %v", duration)
	}

	// In successful concurrent processing, no error logs should be present
	logs := fakeLogger.GetLogs()
	for _, log := range logs {
		if log.Level == "ERROR" {
			t.Errorf("Unexpected error log in concurrent processing: %s", log.Message)
		}
	}
}

// BenchmarkDefaultImageProcessor_SyncPoolOptimization tests the performance impact of sync.Pool.
func BenchmarkDefaultImageProcessor_SyncPoolOptimization(b *testing.B) {
	// Setup test data
	fakeFinder := testutils.NewFakeImageFinder()
	fakeAnalyzer := testutils.NewFakeImageAnalyzer()
	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	dirPath := "/bench/images"
	// Create a larger dataset to see pool benefits
	files := make([]string, 500)
	for i := 0; i < 500; i++ {
		file := fmt.Sprintf("/bench/image%d.jpg", i)
		files[i] = file
		fakeAnalyzer.AddImage(file, image.Image{
			FilePath:         file,
			OriginalFileName: fmt.Sprintf("image%d.jpg", i),
			Date:             time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC),
		})
	}
	fakeFinder.AddFiles(dirPath, files)

	processor, err := NewDefaultImageProcessor(fakeFinder, fakeAnalyzer, fakeLogger, fakeLocalizer)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark the pool-optimized implementation
	b.Run("WithSyncPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			results := processor.Process(dirPath)
			if len(results) != 500 {
				b.Errorf("Expected 500 images, got %d", len(results))
			}
		}
	})
}

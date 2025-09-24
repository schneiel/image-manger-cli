package sort

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

// Contract tests verify that interface implementations maintain consistent behavior
// These tests focus on the contracts defined by interfaces rather than specific implementations

func TestImageFinder_Contract(t *testing.T) {
	t.Parallel()

	implementations := []struct {
		name   string
		finder ImageFinder
	}{
		{
			"DefaultImageFinder",
			NewDefaultImageFinder([]string{".jpg", ".jpeg", ".png", ".gif"}),
		},
	}

	for _, impl := range implementations {
		// capture loop variable
		t.Run(impl.name, func(t *testing.T) {
			t.Parallel()
			testImageFinderContract(t, impl.finder)
		})
	}
}

func testImageFinderContract(t *testing.T, finder ImageFinder) {
	t.Helper()
	t.Run("FindNonExistentDirectory", func(t *testing.T) {
		// Contract: Finding images in non-existent directory should return error
		results, err := finder.Find("/definitely/does/not/exist/path/12345")

		// Contract verification
		require.Error(t, err, "Finding images in non-existent directory should return error")
		assert.Nil(t, results, "Results should be nil when error occurs")
	})

	t.Run("FindConsistentBehavior", func(t *testing.T) {
		// Contract: Multiple calls with same path should be consistent
		testPath := "/definitely/does/not/exist/path/54321"

		result1, err1 := finder.Find(testPath)
		result2, err2 := finder.Find(testPath)

		// Contract verification
		assert.Equal(t, err1 != nil, err2 != nil, "Error behavior should be consistent")
		if err1 == nil && err2 == nil {
			assert.Len(t, result2, len(result1), "Result length should be consistent")
			// Results should contain same files (order may vary)
			assert.ElementsMatch(t, result1, result2, "Results should contain same files")
		}
	})

	t.Run("ExtensionFiltering", func(t *testing.T) {
		// Contract: Only files with allowed extensions should be returned
		// This test verifies the filtering logic by checking extension matching behavior

		// Create a DefaultImageFinder with specific extensions for testing
		testFinder := NewDefaultImageFinder([]string{".jpg", ".png"})

		// Test the internal isImage method if it's accessible, or test the behavior indirectly
		if defaultFinder, ok := testFinder.(*DefaultImageFinder); ok {
			// Test case-insensitive extension matching
			assert.True(t, defaultFinder.isImage("/test/image.jpg"), "Should accept .jpg files")
			assert.True(t, defaultFinder.isImage("/test/image.JPG"), "Should accept .JPG files (case insensitive)")
			assert.True(t, defaultFinder.isImage("/test/image.png"), "Should accept .png files")
			assert.False(t, defaultFinder.isImage("/test/document.txt"), "Should reject .txt files")
			assert.False(t, defaultFinder.isImage("/test/data.json"), "Should reject .json files")
			assert.False(
				t,
				defaultFinder.isImage("/test/image.gif"),
				"Should reject .gif files when not in allowed list",
			)
		}
	})

	t.Run("NilResultContract", func(_ *testing.T) {
		// Contract: Results should never be nil when no error occurs
		// Since we can't easily create a guaranteed empty directory,
		// we test that the contract is maintained by checking result structure

		// This is a documentation test for the contract
		// In a real scenario with a guaranteed empty directory:
		// results, err := finder.Find("/guaranteed/empty/directory")
		// assert.NoError(t, err)
		// assert.NotNil(t, results, "Results should never be nil")
		// assert.Empty(t, results, "Results should be empty for empty directory")

		// Contract documented: Results should never be nil when no error occurs
		// This is a documentation test placeholder
	})
}

func TestImageAnalyzer_Contract(t *testing.T) {
	t.Parallel()

	// Create the analyzer outside the slice to handle the error return
	analyzer, err := NewDefaultImageAnalyzer(
		&testutils.MockDateProcessor{
			GetBestAvailableDateFunc: func(_ map[string]interface{}, _ string) (time.Time, error) {
				return time.Now(), nil
			},
		},
		&testutils.MockExifReader{
			ReadExifFunc: func(_ string) (map[string]interface{}, error) {
				return make(map[string]interface{}), nil
			},
		},
		&testutils.MockLogger{},
		&testutils.MockLocalizer{
			TranslateFunc: func(key string, _ ...map[string]interface{}) string {
				return key
			},
		},
	)
	require.NoError(t, err)

	implementations := []struct {
		name     string
		analyzer ImageAnalyzer
	}{
		{
			"DefaultImageAnalyzer",
			analyzer,
		},
	}

	for _, impl := range implementations {
		// capture loop variable
		t.Run(impl.name, func(t *testing.T) {
			t.Parallel()
			testImageAnalyzerContract(t, impl.analyzer)
		})
	}
}

func testImageAnalyzerContract(t *testing.T, analyzer ImageAnalyzer) {
	t.Helper()
	t.Run("AnalyzeValidImage", func(t *testing.T) {
		// Contract: Analyzing valid image should return Image with populated fields
		testPath := "/test/valid_image.jpg"

		result, err := analyzer.Analyze(testPath)

		// Contract verification
		require.NoError(t, err, "Analyzing valid image should not return error")
		assert.Equal(t, testPath, result.FilePath, "Result should have correct file path")
		assert.NotEmpty(t, result.OriginalFileName, "Result should have original file name")
		assert.False(t, result.Date.IsZero(), "Result should have valid date")
	})

	t.Run("AnalyzeWithDateProcessorError", func(t *testing.T) {
		// Contract: Date processor error should propagate as analysis error
		testPath := "/test/error_image.jpg"

		// Create analyzer with failing date processor
		failingAnalyzer, err := NewDefaultImageAnalyzer(
			&testutils.MockDateProcessor{
				GetBestAvailableDateFunc: func(_ map[string]interface{}, _ string) (time.Time, error) {
					return time.Time{}, errors.New("date processing failed")
				},
			},
			&testutils.MockExifReader{
				ReadExifFunc: func(_ string) (map[string]interface{}, error) {
					return make(map[string]interface{}), nil
				},
			},
			&testutils.MockLogger{},
			&testutils.MockLocalizer{
				TranslateFunc: func(key string, _ ...map[string]interface{}) string {
					return key
				},
			},
		)
		require.NoError(t, err)

		result, err := failingAnalyzer.Analyze(testPath)

		// Contract verification
		require.Error(t, err, "Date processor error should propagate")
		assert.Contains(t, err.Error(), "ErrorGettingDate", "Error should contain localized error key")
		assert.Equal(t, image.Image{}, result, "Result should be zero value on error")
	})

	t.Run("AnalyzeEmptyPath", func(t *testing.T) {
		// Contract: Empty path should be handled gracefully
		result, err := analyzer.Analyze("")

		// Contract verification (implementation may vary, but should be consistent)
		if err != nil {
			// If error is returned, result should be zero value
			assert.Equal(t, image.Image{}, result, "Result should be zero value when error occurs")
		} else {
			// If no error, result should have empty path.
			assert.Empty(t, result.FilePath, "FilePath should match input")
		}
	})

	t.Run("AnalyzeConsistentBehavior", func(t *testing.T) {
		// Contract: Multiple analyses of same file should be consistent
		testPath := "/test/consistent_image.jpg"

		result1, err1 := analyzer.Analyze(testPath)
		result2, err2 := analyzer.Analyze(testPath)

		// Contract verification
		assert.Equal(t, err1 != nil, err2 != nil, "Error behavior should be consistent")
		if err1 == nil && err2 == nil {
			assert.Equal(t, result1.FilePath, result2.FilePath, "FilePath should be consistent")
			assert.Equal(t, result1.OriginalFileName, result2.OriginalFileName, "OriginalFileName should be consistent")
			// Note: Date might vary if based on current time, so we don't assert equality
		}
	})
}

func TestImageProcessor_Contract(t *testing.T) {
	t.Parallel()

	// Create mock dependencies
	mockFinder := &testutils.MockImageFinder{
		FindFunc: func(_ string) ([]string, error) {
			return []string{"/test/image1.jpg", "/test/image2.png"}, nil
		},
	}

	mockAnalyzer := &testutils.MockImageAnalyzer{
		AnalyzeFunc: func(filePath string) (image.Image, error) {
			return image.Image{
				FilePath:         filePath,
				OriginalFileName: "test_image.jpg",
				Date:             time.Now(),
			}, nil
		},
	}

	// Create the processor outside the slice to handle the error return
	processor, err := NewDefaultImageProcessor(
		mockFinder,
		mockAnalyzer,
		&testutils.MockLogger{},
		&testutils.MockLocalizer{
			TranslateFunc: func(key string, _ ...map[string]interface{}) string {
				return key
			},
		},
	)
	require.NoError(t, err)

	implementations := []struct {
		name      string
		processor ImageProcessor
	}{
		{
			"DefaultImageProcessor",
			processor,
		},
	}

	for _, impl := range implementations {
		// capture loop variable
		t.Run(impl.name, func(t *testing.T) {
			t.Parallel()
			testImageProcessorContract(t, impl.processor, mockFinder, mockAnalyzer)
		})
	}
}

func testImageProcessorContract(
	t *testing.T,
	processor ImageProcessor,
	mockFinder *testutils.MockImageFinder,
	mockAnalyzer *testutils.MockImageAnalyzer,
) {
	t.Helper()
	t.Run("ProcessValidDirectory", func(t *testing.T) {
		// Contract: Processing valid directory should return processed images
		results := processor.Process("/test/directory")

		// Contract verification
		assert.NotNil(t, results, "Results should not be nil")
		assert.Len(t, results, 2, "Should process all found images")

		for _, result := range results {
			assert.NotEmpty(t, result.FilePath, "Each result should have file path")
			assert.NotEmpty(t, result.OriginalFileName, "Each result should have original file name")
			assert.False(t, result.Date.IsZero(), "Each result should have valid date")
		}
	})

	t.Run("ProcessEmptyDirectory", func(t *testing.T) {
		// Contract: Processing empty directory should return empty results
		// Temporarily modify mock to return empty results
		originalFindFunc := mockFinder.FindFunc
		mockFinder.FindFunc = func(_ string) ([]string, error) {
			return []string{}, nil
		}
		defer func() {
			mockFinder.FindFunc = originalFindFunc
		}()

		results := processor.Process("/empty/directory")

		// Contract verification
		assert.NotNil(t, results, "Results should not be nil even for empty directory")
		assert.Empty(t, results, "Results should be empty for empty directory")
	})

	t.Run("ProcessWithFinderError", func(t *testing.T) {
		// Contract: Finder error should result in nil results
		originalFindFunc := mockFinder.FindFunc
		mockFinder.FindFunc = func(_ string) ([]string, error) {
			return nil, errors.New("finder error")
		}
		defer func() {
			mockFinder.FindFunc = originalFindFunc
		}()

		results := processor.Process("/error/directory")

		// Contract verification
		assert.Nil(t, results, "Results should be nil when finder returns error")
	})

	t.Run("ProcessWithAnalyzerErrors", func(t *testing.T) {
		// Contract: Analyzer errors should not stop processing other images
		originalAnalyzeFunc := mockAnalyzer.AnalyzeFunc
		mockAnalyzer.AnalyzeFunc = func(filePath string) (image.Image, error) {
			if filePath == "/test/image1.jpg" {
				return image.Image{}, errors.New("analyzer error")
			}
			return image.Image{
				FilePath:         filePath,
				OriginalFileName: "test_image.jpg",
				Date:             time.Now(),
			}, nil
		}
		defer func() {
			mockAnalyzer.AnalyzeFunc = originalAnalyzeFunc
		}()

		results := processor.Process("/test/directory")

		// Contract verification
		assert.NotNil(t, results, "Results should not be nil")
		assert.Len(t, results, 1, "Should process successful images even when some fail")
		assert.Equal(t, "/test/image2.png", results[0].FilePath, "Should contain the successfully processed image")
	})

	t.Run("ProcessConcurrentSafety", func(t *testing.T) {
		// Contract: Concurrent processing should be safe and consistent
		const numGoroutines = 10
		results := make([][]image.Image, numGoroutines)

		// Process same directory concurrently
		done := make(chan int, numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				results[index] = processor.Process("/test/directory")
				done <- index
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Contract verification
		for i, result := range results {
			assert.NotNil(t, result, "Result %d should not be nil", i)
			assert.Len(t, result, 2, "Result %d should have consistent length", i)
		}

		// Verify consistency across all results
		firstResult := results[0]
		for i := 1; i < numGoroutines; i++ {
			assert.Len(t, results[i], len(firstResult), "All results should have same length")
			// Note: Due to concurrent processing, order might vary, so we check content presence
			for _, expectedImg := range firstResult {
				found := false
				for _, actualImg := range results[i] {
					if actualImg.FilePath == expectedImg.FilePath {
						found = true
						break
					}
				}
				assert.True(t, found, "Result %d should contain image with path %s", i, expectedImg.FilePath)
			}
		}
	})
}

func TestActionStrategy_Contract(t *testing.T) {
	t.Parallel()

	// Note: This would test ActionStrategy implementations when they're available
	// For now, we define the contract that implementations should follow

	t.Run("ContractDefinition", func(_ *testing.T) {
		// Document the expected contract for ActionStrategy implementations:

		// 1. Execute should be idempotent when called with same parameters
		// 2. Execute should handle missing source files gracefully
		// 3. Execute should create destination directories as needed
		// 4. Execute should return appropriate errors for invalid operations
		// 5. GetResources should return non-nil resource manager

		// This test serves as documentation until implementations are tested
		// ActionStrategy contract documented
	})
}

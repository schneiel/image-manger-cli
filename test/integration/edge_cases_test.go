package integrationtest

import (
	"path/filepath"
	"testing"
)

func TestGlobalCommand_Help(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run main command with help flag
	result := env.RunCommand("--help")

	// Assert command succeeded
	result.AssertSuccess(t)

	// Assert help text contains expected content
	result.AssertContains(t, "Image Manager")
	result.AssertContains(t, "Available Commands")
	result.AssertContains(t, "sort")
	result.AssertContains(t, "dedup")
	result.AssertContains(t, "completion")
}

func TestGlobalCommand_Version(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run main command without arguments (should show help)
	result := env.RunCommand()

	// Assert command succeeded and shows help
	result.AssertSuccess(t)
	result.AssertContains(t, "Usage:")
}

func TestEdgeCase_LongPaths(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create a deeply nested directory structure
	if err := env.CreateTestImage(filepath.Join("very", "long", "nested", "directory", "structure", "for", "testing", "purposes", "image.png"), TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Run sort command with long paths
	result := env.RunCommand("sort",
		"--source", env.SourceDir(),
		"--destination", env.DestDir(),
		"--actionStrategy", "dryRun")

	// Assert command succeeded
	result.AssertSuccess(t)
}

func TestEdgeCase_SpecialCharacters(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create files with special characters (that are safe for filesystem)
	if err := env.CreateTestImage("image-with-dashes.png", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	if err := env.CreateTestImage("image_with_underscores.jpg", TestImageJPEG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	if err := env.CreateTestImage("image.with.dots.png", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Run sort command
	result := env.RunCommand("sort",
		"--source", env.SourceDir(),
		"--destination", env.DestDir(),
		"--actionStrategy", "dryRun")

	// Assert command succeeded
	result.AssertSuccess(t)
}

func TestEdgeCase_EmptyDirectoryNames(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Test with empty strings as parameters
	result := env.RunCommand("sort",
		"--source", "",
		"--destination", "",
		"--actionStrategy", "dryRun")

	// CLI properly validates required parameters
	result.AssertFailure(t)
	result.AssertStderrContains(t, "required")
}

func TestEdgeCase_NonExistentFlags(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Test with non-existent flag
	result := env.RunCommand("sort",
		"--nonexistentflag", "value",
		"--source", env.SourceDir(),
		"--destination", env.DestDir())

	// Should fail with unknown flag error
	result.AssertFailure(t)
	result.AssertStderrContains(t, "unknown flag")
}

func TestEdgeCase_DuplicateFlags(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test image
	if err := env.CreateTestImage("image.png", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test with duplicate flags (last one should win)
	result := env.RunCommand("sort",
		"--source", "/invalid/path",
		"--source", env.SourceDir(),
		"--destination", env.DestDir(),
		"--actionStrategy", "copy",
		"--actionStrategy", "dryRun")

	// Should succeed with the last values
	result.AssertSuccess(t)
}

func TestEdgeCase_MixedShortLongFlags(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test image
	if err := env.CreateTestImage("image.png", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test mixing short and long flags
	result := env.RunCommand("sort",
		"-s", env.SourceDir(),
		"--destination", env.DestDir(),
		"-l", "de",
		"--actionStrategy", "dryRun")

	// Should succeed
	result.AssertSuccess(t)
}

func TestEdgeCase_CaseSensitivity(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create images with different case extensions
	if err := env.CreateTestImage("image.PNG", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	if err := env.CreateTestImage("image.JPG", TestImageJPEG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	if err := env.CreateTestImage("image.jpeg", TestImageJPEG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Run sort command
	result := env.RunCommand("sort",
		"--source", env.SourceDir(),
		"--destination", env.DestDir(),
		"--actionStrategy", "dryRun")

	// Should handle case variations
	result.AssertSuccess(t)
}

func TestEdgeCase_InvalidLanguage(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test image
	if err := env.CreateTestImage("image.png", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test with invalid language code
	result := env.RunCommand("--language", "invalid",
		"sort",
		"--source", env.SourceDir(),
		"--destination", env.DestDir(),
		"--actionStrategy", "dryRun")

	// CLI properly validates language codes
	result.AssertFailure(t)
	result.AssertStderrContains(t, "not supported")
}

func TestEdgeCase_VeryLargeThreshold(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test image
	if err := env.CreateTestImage("image.png", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test with very large threshold value
	result := env.RunCommand("dedup",
		"--source", env.SourceDir(),
		"--actionStrategy", "dryRun",
		"--threshold", "999999")

	// Should handle large values gracefully
	result.AssertSuccess(t)
}

func TestEdgeCase_NegativeValues(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Create test image
	if err := env.CreateTestImage("image.png", TestImagePNG); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test with negative values
	result := env.RunCommand("dedup",
		"--source", env.SourceDir(),
		"--actionStrategy", "dryRun",
		"--threshold", "-1",
		"--workers", "-5")

	// CLI properly validates that threshold must be non-negative
	result.AssertFailure(t)
	result.AssertStderrContains(t, "non-negative")
}

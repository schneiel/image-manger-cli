// Package integrationtest provides integration testing utilities for the Image Manager CLI.
package integrationtest

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Test image constants - minimal valid image data for testing purposes.
var (
	// TestImagePNG is a minimal valid 1x1 PNG image for testing.
	TestImagePNG = mustDecodeBase64("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==")

	// TestImageJPEG is a minimal valid 1x1 JPEG image for testing.
	TestImageJPEG = mustDecodeBase64("/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwA/AB8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==")
)

// mustDecodeBase64 decodes base64 string and panics on error (for constant initialization).
func mustDecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// TestEnvironment provides a controlled environment for integration testing.
type TestEnvironment struct {
	tempDir     string
	binaryPath  string
	originalDir string
}

// NewTestEnvironment creates a new test environment with temporary directories and builds the CLI binary.
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	// Create temporary directory for this test
	tempDir, err := os.MkdirTemp("", "imagemanager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Get current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Build the CLI binary in temp directory
	binaryPath := filepath.Join(tempDir, "imagemanager")
	if err := buildBinary(originalDir, binaryPath); err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("Failed to build binary: %v", err)
	}

	return &TestEnvironment{
		tempDir:     tempDir,
		binaryPath:  binaryPath,
		originalDir: originalDir,
	}
}

// RunCommand executes the CLI command with the given arguments and returns the result.
func (te *TestEnvironment) RunCommand(args ...string) *CommandResult {
	// Use CommandContext with timeout for test execution
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// #nosec G204 - Test environment with controlled input
	cmd := exec.CommandContext(ctx, te.binaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set working directory to temp dir
	cmd.Dir = te.tempDir

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return &CommandResult{
		ExitCode: exitCode,
		Stdout:   strings.TrimSpace(stdout.String()),
		Stderr:   strings.TrimSpace(stderr.String()),
		Args:     append([]string{"imagemanager"}, args...),
	}
}

// CreateTestFile creates a test file with the given name and content in the test environment.
func (te *TestEnvironment) CreateTestFile(relativePath string, content []byte) error {
	fullPath := filepath.Join(te.tempDir, relativePath)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	return os.WriteFile(fullPath, content, 0600)
}

// CreateTestDir creates a test directory in the test environment.
func (te *TestEnvironment) CreateTestDir(relativePath string) error {
	fullPath := filepath.Join(te.tempDir, relativePath)
	return os.MkdirAll(fullPath, 0750)
}

// TempDir returns the temporary directory path for this test environment.
func (te *TestEnvironment) TempDir() string {
	return te.tempDir
}

// SourceDir returns a source directory path within the temp directory.
func (te *TestEnvironment) SourceDir() string {
	return filepath.Join(te.tempDir, "source")
}

// DestDir returns a destination directory path within the temp directory.
func (te *TestEnvironment) DestDir() string {
	return filepath.Join(te.tempDir, "destination")
}

// ConfigFile returns the path to the test configuration file.
func (te *TestEnvironment) ConfigFile() string {
	return filepath.Join(te.tempDir, "config.yaml")
}

// CreateTestImage creates a test image file with the given path and content.
func (te *TestEnvironment) CreateTestImage(relativePath string, imageData []byte) error {
	return te.CreateTestFile(relativePath, imageData)
}

// CreateRealTestImage creates a test image file with the given path, content, and modification time.
func (te *TestEnvironment) CreateRealTestImage(relativePath string, imageData []byte, modTime time.Time) error {
	if err := te.CreateTestFile(relativePath, imageData); err != nil {
		return err
	}

	// Set the modification time
	fullPath := filepath.Join(te.tempDir, relativePath)
	return os.Chtimes(fullPath, modTime, modTime)
}

// CreateTestConfig creates a test configuration file with the given content.
func (te *TestEnvironment) CreateTestConfig(content string) error {
	return te.CreateTestFile("config.yaml", []byte(content))
}

// Cleanup removes the temporary directory and all its contents.
func (te *TestEnvironment) Cleanup() {
	if te.tempDir != "" {
		_ = os.RemoveAll(te.tempDir)
	}
}

// buildBinary builds the CLI binary from the project root.
func buildBinary(projectRoot, outputPath string) error {
	// Use CommandContext with timeout for build operations
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "build", "-o", outputPath, ".")
	cmd.Dir = projectRoot

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// countFiles counts the number of files recursively in a directory.
func countFiles(t *testing.T, dir string) int {
	t.Helper()
	count := 0
	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		t.Logf("Warning: error counting files in %s: %v", dir, err)
		return 0
	}
	return count
}

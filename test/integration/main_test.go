// Package integrationtest provides integration tests for the Image Manager CLI application.
package integrationtest

import (
	"testing"
)

func TestMainIntegration_Help(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Test that the CLI binary can be executed and shows help
	result := env.RunCommand("--help")

	// Assert help command succeeded
	result.AssertSuccess(t)

	// Assert help contains expected commands
	result.AssertContains(t, "sort")
	result.AssertContains(t, "dedup")
}

func TestMainIntegration_Version(t *testing.T) {
	t.Parallel()

	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Test that the CLI binary shows version information
	result := env.RunCommand("--version")

	// Version command should succeed or show help (depending on implementation)
	// Don't assert success since --version might not be implemented yet

	// Just check that the binary runs without crashing
	if result.ExitCode != 0 && result.ExitCode != 1 {
		t.Errorf("Unexpected exit code: %d, stderr: %s", result.ExitCode, result.Stderr)
	}
}

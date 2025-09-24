// Package integrationtest provides integration testing utilities for the Image Manager CLI.
package integrationtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// CommandResult represents the result of running a CLI command.
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Args     []string
}

// AssertSuccess asserts that the command executed successfully (exit code 0).
func (cr *CommandResult) AssertSuccess(t *testing.T) {
	t.Helper()
	assert.Equal(t, 0, cr.ExitCode,
		"Command failed with exit code %d\nArgs: %v\nStdout: %s\nStderr: %s",
		cr.ExitCode, cr.Args, cr.Stdout, cr.Stderr)
}

// AssertFailure asserts that the command failed (non-zero exit code).
func (cr *CommandResult) AssertFailure(t *testing.T) {
	t.Helper()
	assert.NotEqual(t, 0, cr.ExitCode,
		"Expected command to fail but it succeeded\nArgs: %v\nStdout: %s",
		cr.Args, cr.Stdout)
}

// AssertContains asserts that stdout contains the expected string.
func (cr *CommandResult) AssertContains(t *testing.T, expected string) {
	t.Helper()
	assert.Contains(t, cr.Stdout, expected,
		"Expected stdout to contain '%s'\nActual stdout: %s", expected, cr.Stdout)
}

// AssertStderrContains asserts that stderr contains the expected string.
func (cr *CommandResult) AssertStderrContains(t *testing.T, expected string) {
	t.Helper()
	assert.Contains(t, cr.Stderr, expected,
		"Expected stderr to contain '%s'\nActual stderr: %s", expected, cr.Stderr)
}

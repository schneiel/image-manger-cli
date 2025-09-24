package integrationtest

import (
	"testing"
)

func TestCompletionCommand_Help(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run completion command with help flag
	result := env.RunCommand("completion", "--help")

	// Assert command succeeded
	result.AssertSuccess(t)

	// Assert help text contains expected content
	result.AssertContains(t, "completion")
	result.AssertContains(t, "bash")
	result.AssertContains(t, "fish")
	result.AssertContains(t, "powershell")
	result.AssertContains(t, "zsh")
}

func TestCompletionCommand_Bash(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run bash completion generation
	result := env.RunCommand("completion", "bash")

	// Assert command succeeded
	result.AssertSuccess(t)

	// Assert output contains bash completion script
	result.AssertContains(t, "bash")
}

func TestCompletionCommand_Fish(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run fish completion generation
	result := env.RunCommand("completion", "fish")

	// Assert command succeeded
	result.AssertSuccess(t)

	// Assert output contains fish completion script
	result.AssertContains(t, "fish")
}

func TestCompletionCommand_Zsh(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run zsh completion generation
	result := env.RunCommand("completion", "zsh")

	// Assert command succeeded
	result.AssertSuccess(t)

	// Assert output contains zsh completion script
	result.AssertContains(t, "zsh")
}

func TestCompletionCommand_PowerShell(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run powershell completion generation
	result := env.RunCommand("completion", "powershell")

	// Assert command succeeded
	result.AssertSuccess(t)

	// Assert output contains powershell completion script
	result.AssertContains(t, "powershell")
}

func TestCompletionCommand_InvalidShell(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run completion with invalid shell
	result := env.RunCommand("completion", "invalidshell")

	// CLI shows help for invalid subcommands rather than failing
	result.AssertSuccess(t)
	result.AssertContains(t, "Available Commands")
}

func TestCompletionCommand_NoArguments(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t)
	defer env.Cleanup()

	// Run completion command without arguments
	result := env.RunCommand("completion")

	// Assert command succeeded and shows help
	result.AssertSuccess(t)
	result.AssertContains(t, "Available Commands")
}

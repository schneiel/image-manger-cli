package builders_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/internal/builders"
)

const (
	testConfigValue = "test config"
)

// MockCommandHandler is a mock implementation of handlers.CommandHandler.
type MockCommandHandler struct {
	mock.Mock
}

func (m *MockCommandHandler) RunE(cmd *cobra.Command, args []string) error {
	args2 := m.Called(cmd, args)

	err := args2.Error(0)
	if err != nil {
		return fmt.Errorf("mock RunE error: %w", err)
	}
	return nil
}

func (m *MockCommandHandler) Execute(config any) error {
	args := m.Called(config)

	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("mock Execute error: %w", err)
	}
	return nil
}

func (m *MockCommandHandler) Validate(config any) error {
	args := m.Called(config)

	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("mock Validate error: %w", err)
	}
	return nil
}

func TestWithDescription(t *testing.T) {
	t.Parallel()
	// Test that WithDescription option can be created without error
	option := builders.WithDescription("short", "long")
	assert.NotNil(t, option)

	// Adding test for empty description
	option = builders.WithDescription("", "")
	assert.NotNil(t, option)
}

func TestWithHandler(t *testing.T) {
	t.Parallel()
	// Test that WithHandler option can be created without error
	handler := &MockCommandHandler{}
	option := builders.WithHandler(handler)
	assert.NotNil(t, option)
}

func TestWithConfig(t *testing.T) {
	t.Parallel()
	// Test that WithConfig option can be created without error
	configBuilder := func() any { return testConfigValue }
	option := builders.WithConfig(configBuilder)
	assert.NotNil(t, option)
}

func TestWithFlags(t *testing.T) {
	t.Parallel()
	// Test that WithFlags option can be created without error
	flagSetup := func(cmd *cobra.Command, _ any) {
		cmd.Flags().String("test", "", "test flag")
	}
	option := builders.WithFlags(flagSetup)
	assert.NotNil(t, option)
}

func TestNewCommand_Basic(t *testing.T) {
	t.Parallel()
	cmd := builders.NewCommand("test")

	assert.NotNil(t, cmd)
	assert.Equal(t, "test", cmd.Use)
	assert.Empty(t, cmd.Short)
	assert.Empty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}

func TestNewCommand_WithDescription(t *testing.T) {
	t.Parallel()
	cmd := builders.NewCommand("test",
		builders.WithDescription("short desc", "long desc"),
	)

	assert.NotNil(t, cmd)
	assert.Equal(t, "test", cmd.Use)
	assert.Equal(t, "short desc", cmd.Short)
	assert.Equal(t, "long desc", cmd.Long)
}

func TestNewCommand_WithHandler(t *testing.T) {
	t.Parallel()
	handler := &MockCommandHandler{}
	handler.On("Validate", mock.Anything).Return(nil)
	handler.On("Execute", mock.Anything).Return(nil)

	cmd := builders.NewCommand("test",
		builders.WithHandler(handler),
	)

	assert.NotNil(t, cmd)
	assert.Equal(t, "test", cmd.Use)
	assert.NotNil(t, cmd.RunE)

	// Test RunE function
	err := cmd.RunE(nil, nil)
	require.NoError(t, err)

	handler.AssertExpectations(t)
}

func TestNewCommand_WithConfig(t *testing.T) {
	t.Parallel()
	handler := &MockCommandHandler{}
	handler.On("Validate", testConfigValue).Return(nil)
	handler.On("Execute", testConfigValue).Return(nil)

	configBuilder := func() any { return testConfigValue }

	cmd := builders.NewCommand("test",
		builders.WithHandler(handler),
		builders.WithConfig(configBuilder),
	)

	assert.NotNil(t, cmd)

	// Test RunE function
	err := cmd.RunE(nil, nil)
	require.NoError(t, err)

	handler.AssertExpectations(t)
}

func TestNewCommand_WithFlags(t *testing.T) {
	t.Parallel()
	handler := &MockCommandHandler{}
	handler.On("Validate", mock.Anything).Return(nil)
	handler.On("Execute", mock.Anything).Return(nil)

	flagSetup := func(cmd *cobra.Command, _ any) {
		cmd.Flags().String("test", "", "test flag")
	}

	configBuilder := func() any { return testConfigValue }

	cmd := builders.NewCommand("test",
		builders.WithHandler(handler),
		builders.WithConfig(configBuilder),
		builders.WithFlags(flagSetup),
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, cmd.Flag("test"))

	// Test RunE function
	err := cmd.RunE(nil, nil)
	require.NoError(t, err)

	handler.AssertExpectations(t)
}

func TestNewCommand_ValidationError(t *testing.T) {
	t.Parallel()
	handler := &MockCommandHandler{}
	handler.On("Validate", mock.Anything).Return(assert.AnError)

	cmd := builders.NewCommand("test",
		builders.WithHandler(handler),
	)

	// Test RunE function
	err := cmd.RunE(nil, nil)
	require.Error(t, err)

	handler.AssertExpectations(t)
}

func TestNewCommand_ExecuteError(t *testing.T) {
	t.Parallel()
	handler := &MockCommandHandler{}
	handler.On("Validate", mock.Anything).Return(nil)
	handler.On("Execute", mock.Anything).Return(assert.AnError)

	cmd := builders.NewCommand("test",
		builders.WithHandler(handler),
	)

	// Test RunE function
	err := cmd.RunE(nil, nil)
	require.Error(t, err)

	handler.AssertExpectations(t)
}

func TestNewCommand_NoHandler(t *testing.T) {
	t.Parallel()
	cmd := builders.NewCommand("test")

	// Test RunE function - should panic when no handler is provided
	assert.Panics(t, func() {
		_ = cmd.RunE(nil, nil) // Error is expected and handled by panic
	})
}

func TestNewCommand_HelpFlag(t *testing.T) {
	t.Parallel()
	cmd := builders.NewCommand("test")

	// Should have help flag by default
	helpFlag := cmd.Flags().Lookup("help")
	assert.NotNil(t, helpFlag)
}

func TestNewCommand_DisableFlagParsing(t *testing.T) {
	t.Parallel()
	cmd := builders.NewCommand("test")

	// Should not disable flag parsing by default
	assert.False(t, cmd.DisableFlagParsing)
}

package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/internal/di"
)

func TestNewContainerBuilder(t *testing.T) {
	t.Parallel()
	builder := di.NewContainerBuilder()
	assert.NotNil(t, builder)
}

func TestContainerBuilder_WithArgs(t *testing.T) {
	t.Parallel()
	args := []string{"test", "arg1", "arg2"}
	builder := di.NewContainerBuilder().WithArgs(args)
	assert.NotNil(t, builder)
}

func TestContainerBuilder_WithLocalesFS(t *testing.T) {
	t.Parallel()
	args := []string{"test"}
	builder := di.NewContainerBuilder().WithArgs(args).WithLocalesFS(nil)
	assert.NotNil(t, builder)
}

func TestContainerBuilder_ErrorPropagation(t *testing.T) {
	t.Parallel()

	// Test error propagation in WithArgs
	builder := di.NewContainerBuilder()
	// Simulate error condition by providing invalid input that would cause an error during build
	builder = builder.WithArgs(nil)
	assert.NotNil(t, builder)

	// Test error propagation in WithLocalesFS
	builder = builder.WithLocalesFS(nil)
	assert.NotNil(t, builder)
}

func TestContainerBuilder_BuildInvalidArgs(t *testing.T) {
	t.Parallel()

	// Test with invalid args that would cause build errors
	args := []string{} // Empty args should cause parsing issues
	builder := di.NewContainerBuilder().WithArgs(args).WithLocalesFS(nil)

	_, err := builder.Build()
	// We expect this to fail because of missing embedded filesystem and invalid config
	require.Error(t, err)
}

func TestNewContainer_InvalidInput(t *testing.T) {
	t.Parallel()

	// Test NewContainer with invalid inputs
	_, err := di.NewContainer([]string{}, nil)
	require.Error(t, err, "Should fail with nil localesFS and empty args")
}

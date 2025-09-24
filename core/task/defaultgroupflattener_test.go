package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultGroupFlattener(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	assert.NotNil(t, flattener)
	assert.Equal(t, logger, flattener.logger)
}

func TestNewDefaultGroupFlattener_NilLogger(t *testing.T) {
	t.Parallel()
	_, err := NewDefaultGroupFlattener(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "logger cannot be nil")
}

// Removed TestNewDefaultGroupFlattener_NilLocalizer as localizer is no longer a dependency

func TestDefaultGroupFlattener_Flatten_EmptyInput(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	// Track logger calls
	debugCallCount := 0
	logger.DebugFunc = func(_ string) {
		debugCallCount++
	}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	result := flattener.Flatten([][]string{})

	assert.Empty(t, result)
	assert.Equal(t, 1, debugCallCount) // Should log flattening operation
}

func TestDefaultGroupFlattener_Flatten_SingleGroup(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	// Track logger calls
	debugfCallCount := 0
	logger.DebugfFunc = func(_ string, _ ...interface{}) {
		debugfCallCount++
	}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	input := [][]string{
		{"file1.jpg", "file2.jpg", "file3.jpg"},
	}

	result := flattener.Flatten(input)

	expected := []string{"file1.jpg", "file2.jpg", "file3.jpg"}
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, debugfCallCount) // Should log flattening operation
}

func TestDefaultGroupFlattener_Flatten_MultipleGroups(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	// Track logger calls
	debugfCallCount := 0
	logger.DebugfFunc = func(_ string, _ ...interface{}) {
		debugfCallCount++
	}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	input := [][]string{
		{"file1.jpg", "file2.jpg"},
		{"file3.jpg"},
		{"file4.jpg", "file5.jpg", "file6.jpg"},
	}

	result := flattener.Flatten(input)

	expected := []string{"file1.jpg", "file2.jpg", "file3.jpg", "file4.jpg", "file5.jpg", "file6.jpg"}
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, debugfCallCount) // Should log flattening operation
}

func TestDefaultGroupFlattener_Flatten_EmptyGroups(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	// Track logger calls (both Debug and Debugf will be called)
	debugfCallCount := 0
	logger.DebugfFunc = func(_ string, _ ...interface{}) {
		debugfCallCount++
	}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	input := [][]string{
		{"file1.jpg"},
		{}, // Empty group
		{"file2.jpg", "file3.jpg"},
		{}, // Another empty group
	}

	result := flattener.Flatten(input)

	expected := []string{"file1.jpg", "file2.jpg", "file3.jpg"}
	assert.Equal(t, expected, result)
	assert.Equal(t, 3, debugfCallCount) // Should log empty groups (2) + final summary (1)
}

func TestDefaultGroupFlattener_Flatten_NilInput(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	// Track logger calls
	debugCallCount := 0
	logger.DebugFunc = func(_ string) {
		debugCallCount++
	}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	result := flattener.Flatten(nil)

	assert.Empty(t, result)
	assert.Equal(t, 1, debugCallCount) // Should log flattening operation
}

func TestDefaultGroupFlattener_Flatten_PreservesOrder(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	// Track logger calls
	debugfCallCount := 0
	logger.DebugfFunc = func(_ string, _ ...interface{}) {
		debugfCallCount++
	}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	input := [][]string{
		{"z.jpg", "a.jpg"},
		{"m.jpg"},
		{"b.jpg", "y.jpg"},
	}

	result := flattener.Flatten(input)

	// Should preserve the order from input groups
	expected := []string{"z.jpg", "a.jpg", "m.jpg", "b.jpg", "y.jpg"}
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, debugfCallCount) // Should log flattening operation
}

func TestDefaultGroupFlattener_Flatten_LoggingBehavior(t *testing.T) {
	t.Parallel()
	logger := &testutils.MockLogger{}

	// Track logger calls with message content
	loggedMessages := []string{}
	logger.DebugfFunc = func(format string, _ ...interface{}) {
		loggedMessages = append(loggedMessages, format)
	}

	flattener, err := NewDefaultGroupFlattener(logger)
	require.NoError(t, err)

	input := [][]string{
		{"file1.jpg", "file2.jpg"},
	}

	flattener.Flatten(input)

	assert.Len(t, loggedMessages, 1)
}

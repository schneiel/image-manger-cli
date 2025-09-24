package log

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConsoleWriter_Stdout(t *testing.T) {
	t.Parallel()
	// Arrange
	writer := NewDefaultConsoleWriter()

	// Act
	stdout := writer.Stdout()

	// Assert
	if stdout != os.Stdout {
		t.Errorf("Stdout() should return os.Stdout, got %v", stdout)
	}
}

func TestDefaultConsoleWriter_Stderr(t *testing.T) {
	t.Parallel()
	// Arrange
	writer := NewDefaultConsoleWriter()

	// Act
	stderr := writer.Stderr()

	// Assert
	if stderr != os.Stderr {
		t.Errorf("Stderr() should return os.Stderr, got %v", stderr)
	}
}

// MockConsoleWriter is a mock implementation for testing.
type MockConsoleWriter struct {
	stdout io.Writer
	stderr io.Writer
}

func (m *MockConsoleWriter) Stdout() io.Writer {
	return m.stdout
}

func (m *MockConsoleWriter) Stderr() io.Writer {
	return m.stderr
}

func TestDefaultConsoleLoggerWithWriter_Integration(t *testing.T) {
	t.Parallel()
	// Arrange
	mockWriter := &MockConsoleWriter{
		stdout: &MockWriter{},
		stderr: &MockWriter{},
	}
	logger, err := NewDefaultConsoleLoggerWithWriter(INFO, mockWriter)
	require.NoError(t, err)

	// Act & Assert - should not panic
	logger.Info("Test message")
	logger.Error("Test error")
}

// MockWriter is a simple mock writer for testing.
type MockWriter struct {
	written []byte
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	m.written = append(m.written, p...)
	return len(p), nil
}

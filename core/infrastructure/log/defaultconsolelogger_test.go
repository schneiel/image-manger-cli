package log

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper methods for MockWriter.
func (m *MockWriter) String() string {
	return string(m.written)
}

func (m *MockWriter) Reset() {
	m.written = nil
}

// Helper function to create a mock console writer with MockWriter instances.
func newTestMockConsoleWriter() (*MockConsoleWriter, *MockWriter, *MockWriter) {
	stdoutMock := &MockWriter{}
	stderrMock := &MockWriter{}
	consoleWriter := &MockConsoleWriter{
		stdout: stdoutMock,
		stderr: stderrMock,
	}
	return consoleWriter, stdoutMock, stderrMock
}

func TestNewDefaultConsoleLogger(t *testing.T) {
	t.Parallel()

	logger, err := NewDefaultConsoleLogger(INFO)
	require.NoError(t, err)

	require.NotNil(t, logger)
	assert.Equal(t, INFO, logger.minLevel)
	assert.NotNil(t, logger.stdoutLogger)
	assert.NotNil(t, logger.stderrLogger)
	assert.NotNil(t, logger.consoleWriter)
}

func TestNewDefaultConsoleLoggerWithWriter(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, _, _ := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	require.NotNil(t, logger)
	assert.Equal(t, DEBUG, logger.minLevel)
	assert.Equal(t, mockConsoleWriter, logger.consoleWriter)
}

func TestNewDefaultConsoleLoggerWithWriter_NilWriter(t *testing.T) {
	t.Parallel()

	_, err := NewDefaultConsoleLoggerWithWriter(INFO, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "consoleWriter cannot be nil")
}

func TestDefaultConsoleLogger_SetLevel(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, _, _ := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(INFO, mockConsoleWriter)
	require.NoError(t, err)

	// Change level
	logger.SetLevel(ERROR)
	assert.Equal(t, ERROR, logger.minLevel)

	// Change level again
	logger.SetLevel(DEBUG)
	assert.Equal(t, DEBUG, logger.minLevel)
}

func TestDefaultConsoleLogger_Debug(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Debug("test debug message")

	output := stdoutMock.String()
	assert.Contains(t, output, "[DEBUG] test debug message")
	assert.Empty(t, stderrMock.String())
}

func TestDefaultConsoleLogger_Info(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Info("test info message")

	output := stdoutMock.String()
	assert.Contains(t, output, "[INFO] test info message")
	assert.Empty(t, stderrMock.String())
}

func TestDefaultConsoleLogger_Warn(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Warn("test warn message")

	output := stdoutMock.String()
	assert.Contains(t, output, "[WARN] test warn message")
	assert.Empty(t, stderrMock.String())
}

func TestDefaultConsoleLogger_Error(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Error("test error message")

	output := stderrMock.String()
	assert.Contains(t, output, "[ERROR] test error message")
	assert.Empty(t, stdoutMock.String())
}

func TestDefaultConsoleLogger_Debugf(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Debugf("test debug %s %d", "message", 42)

	output := stdoutMock.String()
	assert.Contains(t, output, "[DEBUG] test debug message 42")
	assert.Empty(t, stderrMock.String())
}

func TestDefaultConsoleLogger_Infof(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Infof("test info %s %d", "message", 42)

	output := stdoutMock.String()
	assert.Contains(t, output, "[INFO] test info message 42")
	assert.Empty(t, stderrMock.String())
}

func TestDefaultConsoleLogger_Warnf(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Warnf("test warn %s %d", "message", 42)

	output := stdoutMock.String()
	assert.Contains(t, output, "[WARN] test warn message 42")
	assert.Empty(t, stderrMock.String())
}

func TestDefaultConsoleLogger_Errorf(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Errorf("test error %s %d", "message", 42)

	output := stderrMock.String()
	assert.Contains(t, output, "[ERROR] test error message 42")
	assert.Empty(t, stdoutMock.String())
}

func TestDefaultConsoleLogger_LevelFiltering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		minLevel Level
		logLevel Level
		message  string
		should   bool
	}{
		{"DEBUG level allows DEBUG", DEBUG, DEBUG, "debug msg", true},
		{"DEBUG level allows INFO", DEBUG, INFO, "info msg", true},
		{"DEBUG level allows WARN", DEBUG, WARN, "warn msg", true},
		{"DEBUG level allows ERROR", DEBUG, ERROR, "error msg", true},
		{"INFO level blocks DEBUG", INFO, DEBUG, "debug msg", false},
		{"INFO level allows INFO", INFO, INFO, "info msg", true},
		{"INFO level allows WARN", INFO, WARN, "warn msg", true},
		{"INFO level allows ERROR", INFO, ERROR, "error msg", true},
		{"WARN level blocks DEBUG", WARN, DEBUG, "debug msg", false},
		{"WARN level blocks INFO", WARN, INFO, "info msg", false},
		{"WARN level allows WARN", WARN, WARN, "warn msg", true},
		{"WARN level allows ERROR", WARN, ERROR, "error msg", true},
		{"ERROR level blocks DEBUG", ERROR, DEBUG, "debug msg", false},
		{"ERROR level blocks INFO", ERROR, INFO, "info msg", false},
		{"ERROR level blocks WARN", ERROR, WARN, "warn msg", false},
		{"ERROR level allows ERROR", ERROR, ERROR, "error msg", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
			logger, err := NewDefaultConsoleLoggerWithWriter(tt.minLevel, mockConsoleWriter)
			require.NoError(t, err)

			// Log the message based on level
			switch tt.logLevel {
			case DEBUG:
				logger.Debug(tt.message)
			case INFO:
				logger.Info(tt.message)
			case WARN:
				logger.Warn(tt.message)
			case ERROR:
				logger.Error(tt.message)
			}

			// Check if message was logged
			stdoutOutput := stdoutMock.String()
			stderrOutput := stderrMock.String()
			totalOutput := stdoutOutput + stderrOutput

			if tt.should {
				assert.Contains(t, totalOutput, tt.message, "Message should be logged")
			} else {
				assert.NotContains(t, totalOutput, tt.message, "Message should be filtered out")
			}
		})
	}
}

func TestDefaultConsoleLogger_OutputRouting(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, stderrMock := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	// Test that ERROR goes to stderr
	logger.Error("error message")
	assert.Contains(t, stderrMock.String(), "error message")
	assert.Empty(t, stdoutMock.String())

	// Reset and test that other levels go to stdout
	stdoutMock.Reset()
	stderrMock.Reset()

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")

	stdoutOutput := stdoutMock.String()
	assert.Contains(t, stdoutOutput, "debug message")
	assert.Contains(t, stdoutOutput, "info message")
	assert.Contains(t, stdoutOutput, "warn message")
	assert.Empty(t, stderrMock.String())
}

func TestDefaultConsoleLogger_MessageFormat(t *testing.T) {
	t.Parallel()

	mockConsoleWriter, stdoutMock, _ := newTestMockConsoleWriter()
	logger, err := NewDefaultConsoleLoggerWithWriter(DEBUG, mockConsoleWriter)
	require.NoError(t, err)

	logger.Info("test message")

	output := stdoutMock.String()
	// Check that output contains timestamp, level, and message
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "test message")
	// Check that it contains date/time (basic check)
	assert.True(t, strings.Contains(output, "/") || strings.Contains(output, ":"))
}

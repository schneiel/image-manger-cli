// Package log provides SOLID-compliant logging components.
package log

import (
	"io"
	"os"
)

// ConsoleWriter defines the interface for console output operations.
type ConsoleWriter interface {
	// Stdout returns the standard output writer
	Stdout() io.Writer
	// Stderr returns the standard error writer
	Stderr() io.Writer
}

// DefaultConsoleWriter implements ConsoleWriter using the actual OS stdout/stderr.
type DefaultConsoleWriter struct{}

// NewDefaultConsoleWriter creates a new DefaultConsoleWriter.
func NewDefaultConsoleWriter() ConsoleWriter {
	return &DefaultConsoleWriter{}
}

// Stdout returns the standard output writer.
func (w *DefaultConsoleWriter) Stdout() io.Writer {
	return os.Stdout
}

// Stderr returns the standard error writer.
func (w *DefaultConsoleWriter) Stderr() io.Writer {
	return os.Stderr
}

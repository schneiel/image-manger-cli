// Package testutils provides mock implementations for testing.
package testutils

import (
	"github.com/spf13/cobra"
)

// CommandHandler interface duplicated here to avoid import cycles.
type CommandHandler interface {
	RunE(cmd *cobra.Command, args []string) error
	Execute(config interface{}) error
	Validate(config interface{}) error
}

// FlagSetup interface duplicated here to avoid import cycles.
type FlagSetup interface {
	SetupFlags(cmd *cobra.Command, cfg interface{})
}

// MockSortHandler implements CommandHandler for testing.
type MockSortHandler struct {
	RunEFunc     func(cmd *cobra.Command, args []string) error
	ExecuteFunc  func(config interface{}) error
	ValidateFunc func(config interface{}) error
}

// RunE implements CommandHandler interface.
func (m *MockSortHandler) RunE(cmd *cobra.Command, args []string) error {
	if m.RunEFunc != nil {
		return m.RunEFunc(cmd, args)
	}
	return nil
}

// Execute implements CommandHandler interface.
func (m *MockSortHandler) Execute(config interface{}) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(config)
	}
	return nil
}

// Validate implements CommandHandler interface.
func (m *MockSortHandler) Validate(config interface{}) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(config)
	}
	return nil
}

// MockDedupHandler implements CommandHandler for testing.
type MockDedupHandler struct {
	RunEFunc     func(cmd *cobra.Command, args []string) error
	ExecuteFunc  func(config interface{}) error
	ValidateFunc func(config interface{}) error
}

// RunE implements CommandHandler interface.
func (m *MockDedupHandler) RunE(cmd *cobra.Command, args []string) error {
	if m.RunEFunc != nil {
		return m.RunEFunc(cmd, args)
	}
	return nil
}

// Execute implements CommandHandler interface.
func (m *MockDedupHandler) Execute(config interface{}) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(config)
	}
	return nil
}

// Validate implements CommandHandler interface.
func (m *MockDedupHandler) Validate(config interface{}) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(config)
	}
	return nil
}

// MockFlagSetup implements FlagSetup for testing.
type MockFlagSetup struct {
	SetupFlagsFunc func(cmd *cobra.Command, cfg interface{})
}

// SetupFlags implements FlagSetup interface.
func (m *MockFlagSetup) SetupFlags(cmd *cobra.Command, cfg interface{}) {
	if m.SetupFlagsFunc != nil {
		m.SetupFlagsFunc(cmd, cfg)
	}
}

// MockCommandLocalizer implements command localization for testing.
type MockCommandLocalizer struct {
	TranslateFunc func(key string, args ...map[string]interface{}) string
}

// Translate implements command localization interface.
func (m *MockCommandLocalizer) Translate(key string, args ...map[string]interface{}) string {
	if m.TranslateFunc != nil {
		return m.TranslateFunc(key, args...)
	}
	return key
}

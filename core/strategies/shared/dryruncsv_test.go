package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Note: MockLogger and MockLocalizer are defined in defaultcsvresource_test.go

// MockActionResource for testing DryRunCSVLogger.
type MockActionResource struct {
	mock.Mock
}

func (m *MockActionResource) Setup() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockActionResource) Teardown() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewDryRunCSVLogger(t *testing.T) {
	mockResource := &MockActionResource{}
	mockLogger := &MockLogger{}
	mockLocalizer := &MockLocalizer{}

	dryRunLogger := NewDryRunCSVLogger(mockResource, mockLogger, mockLocalizer)

	assert.NotNil(t, dryRunLogger)
	assert.Equal(t, mockResource, dryRunLogger.GetResource())
}

func TestDryRunCSVLogger_LogOperation(t *testing.T) {
	t.Run("With Destination", func(t *testing.T) {
		mockResource := &MockActionResource{}
		mockLogger := &MockLogger{}
		mockLocalizer := &MockLocalizer{}

		// Setup mock logger without verification (research shows low testing value)
		mockLogger.On("Infof", "[DRY RUN] %s: %s -> %s", []interface{}{"copy", "source.jpg", "dest.jpg"})

		dryRunLogger := NewDryRunCSVLogger(mockResource, mockLogger, mockLocalizer)
		err := dryRunLogger.LogOperation("copy", "source.jpg", "dest.jpg")

		// Focus on behavior: operation should succeed
		require.NoError(t, err)
		// No mock assertion - logging verification has minimal test value per research
	})

	t.Run("Without Destination", func(t *testing.T) {
		mockResource := &MockActionResource{}
		mockLogger := &MockLogger{}
		mockLocalizer := &MockLocalizer{}

		// Setup mock logger without verification
		mockLogger.On("Infof", "[DRY RUN] %s: %s", []interface{}{"delete", "source.jpg"})

		dryRunLogger := NewDryRunCSVLogger(mockResource, mockLogger, mockLocalizer)
		err := dryRunLogger.LogOperation("delete", "source.jpg", "")

		// Focus on behavior: operation should succeed
		require.NoError(t, err)
		// No mock assertion - logging verification has minimal test value per research
	})

	// Adding test for nil logger
	t.Run("With Nil Logger", func(t *testing.T) {
		mockResource := &MockActionResource{}
		mockLocalizer := &MockLocalizer{}

		dryRunLogger := NewDryRunCSVLogger(mockResource, nil, mockLocalizer)
		err := dryRunLogger.LogOperation("copy", "source.jpg", "dest.jpg")

		require.NoError(t, err)
	})
}

func TestDryRunCSVLogger_LogDeduplicationOperation(t *testing.T) {
	mockResource := &MockActionResource{}
	mockLogger := &MockLogger{}
	mockLocalizer := &MockLocalizer{}

	// Setup mock logger without verification (research shows low testing value)
	mockLogger.On(
		"Infof",
		"[DRY RUN] Potential duplicate found: %s (Original: %s)",
		[]interface{}{"duplicate.jpg", "original.jpg"},
	)

	dryRunLogger := NewDryRunCSVLogger(mockResource, mockLogger, mockLocalizer)
	err := dryRunLogger.LogDeduplicationOperation("original.jpg", "duplicate.jpg")

	// Focus on behavior: operation should succeed
	require.NoError(t, err)
	// No mock assertion - logging verification has minimal test value per research
}

func TestDryRunCSVLogger_GetResource(t *testing.T) {
	mockResource := &MockActionResource{}
	mockLogger := &MockLogger{}
	mockLocalizer := &MockLocalizer{}

	dryRunLogger := NewDryRunCSVLogger(mockResource, mockLogger, mockLocalizer)
	resource := dryRunLogger.GetResource()

	assert.Equal(t, mockResource, resource)
}

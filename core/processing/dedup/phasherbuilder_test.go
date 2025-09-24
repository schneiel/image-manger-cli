package dedup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestWithWorkers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		count       int
		expectError bool
		errorMsg    string
	}{
		{"valid worker count", 4, false, ""},
		{"single worker", 1, false, ""},
		{"high worker count", 100, false, ""},
		{"zero workers", 0, true, "worker count must be greater than 0"},
		{"negative workers", -1, true, "worker count must be greater than 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config := &PHasherConfig{}
			option := WithWorkers(tt.count)
			err := option(config)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.count, config.numWorkers)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		logger      interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid logger", testutils.NewFakeLogger(), false, ""},
		{"nil logger", nil, true, "logger cannot be nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config := &PHasherConfig{}
			var option PHasherOption
			if tt.logger != nil {
				option = WithLogger(tt.logger.(*testutils.FakeLogger))
			} else {
				option = WithLogger(nil)
			}
			err := option(config)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.logger, config.logger)
			}
		})
	}
}

func TestWithFilesystem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		filesystem  interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid filesystem", testutils.NewFakeFileSystem(), false, ""},
		{"nil filesystem", nil, true, "filesystem cannot be nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config := &PHasherConfig{}
			var option PHasherOption
			if tt.filesystem != nil {
				option = WithFilesystem(tt.filesystem.(*testutils.FakeFileSystem))
			} else {
				option = WithFilesystem(nil)
			}
			err := option(config)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.filesystem, config.filesystem)
			}
		})
	}
}

func TestNewDefaultPHasherWithOptions_Success(t *testing.T) {
	t.Parallel()

	fakeLogger := testutils.NewFakeLogger()
	fakeFileSystem := testutils.NewFakeFileSystem()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithWorkers(8),
		WithLogger(fakeLogger),
		WithFilesystem(fakeFileSystem),
	)

	require.NoError(t, err)
	assert.NotNil(t, hasher)

	// Verify the hasher is properly configured
	defaultHasher, ok := hasher.(*DefaultPHasher)
	require.True(t, ok)
	assert.Equal(t, 8, defaultHasher.numWorkers)
	assert.Equal(t, fakeLogger, defaultHasher.logger)
	assert.Equal(t, fakeFileSystem, defaultHasher.fs)
	assert.Equal(t, fakeLocalizer, defaultHasher.localizer)
}

func TestNewDefaultPHasherWithOptions_DefaultWorkers(t *testing.T) {
	t.Parallel()

	fakeLogger := testutils.NewFakeLogger()
	fakeFileSystem := testutils.NewFakeFileSystem()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithLogger(fakeLogger),
		WithFilesystem(fakeFileSystem),
	)

	require.NoError(t, err)
	assert.NotNil(t, hasher)

	// Verify default worker count is used
	defaultHasher, ok := hasher.(*DefaultPHasher)
	require.True(t, ok)
	assert.Equal(t, 4, defaultHasher.numWorkers) // Default value
}

func TestNewDefaultPHasherWithOptions_NilLocalizer(t *testing.T) {
	t.Parallel()

	fakeLogger := testutils.NewFakeLogger()
	fakeFileSystem := testutils.NewFakeFileSystem()

	hasher, err := NewDefaultPHasherWithOptions(
		nil,
		WithLogger(fakeLogger),
		WithFilesystem(fakeFileSystem),
	)

	require.Error(t, err)
	assert.Nil(t, hasher)
	assert.Contains(t, err.Error(), "localizer cannot be nil")
}

func TestNewDefaultPHasherWithOptions_MissingLogger(t *testing.T) {
	t.Parallel()

	fakeFileSystem := testutils.NewFakeFileSystem()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithFilesystem(fakeFileSystem),
	)

	require.Error(t, err)
	assert.Nil(t, hasher)
	assert.Contains(t, err.Error(), "logger is required")
}

func TestNewDefaultPHasherWithOptions_MissingFilesystem(t *testing.T) {
	t.Parallel()

	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithLogger(fakeLogger),
	)

	require.Error(t, err)
	assert.Nil(t, hasher)
	assert.Contains(t, err.Error(), "filesystem is required")
}

func TestNewDefaultPHasherWithOptions_InvalidWorkerOption(t *testing.T) {
	t.Parallel()

	fakeLogger := testutils.NewFakeLogger()
	fakeFileSystem := testutils.NewFakeFileSystem()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithWorkers(-1), // Invalid worker count
		WithLogger(fakeLogger),
		WithFilesystem(fakeFileSystem),
	)

	require.Error(t, err)
	assert.Nil(t, hasher)
	assert.Contains(t, err.Error(), "worker count must be greater than 0")
}

func TestNewDefaultPHasherWithOptions_InvalidLoggerOption(t *testing.T) {
	t.Parallel()

	fakeFileSystem := testutils.NewFakeFileSystem()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithWorkers(4),
		WithLogger(nil), // Invalid logger
		WithFilesystem(fakeFileSystem),
	)

	require.Error(t, err)
	assert.Nil(t, hasher)
	assert.Contains(t, err.Error(), "logger cannot be nil")
}

func TestNewDefaultPHasherWithOptions_InvalidFilesystemOption(t *testing.T) {
	t.Parallel()

	fakeLogger := testutils.NewFakeLogger()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithWorkers(4),
		WithLogger(fakeLogger),
		WithFilesystem(nil), // Invalid filesystem
	)

	require.Error(t, err)
	assert.Nil(t, hasher)
	assert.Contains(t, err.Error(), "filesystem cannot be nil")
}

func TestNewDefaultPHasherWithOptions_MultipleOptions(t *testing.T) {
	t.Parallel()

	fakeLogger := testutils.NewFakeLogger()
	fakeFileSystem := testutils.NewFakeFileSystem()
	fakeLocalizer := testutils.NewFakeLocalizer()

	hasher, err := NewDefaultPHasherWithOptions(
		fakeLocalizer,
		WithWorkers(12),
		WithLogger(fakeLogger),
		WithFilesystem(fakeFileSystem),
		WithWorkers(6), // Should override previous worker count
	)

	require.NoError(t, err)
	assert.NotNil(t, hasher)

	// Verify the last worker option takes precedence
	defaultHasher, ok := hasher.(*DefaultPHasher)
	require.True(t, ok)
	assert.Equal(t, 6, defaultHasher.numWorkers)
}

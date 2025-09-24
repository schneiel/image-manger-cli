package di_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/time"
	"github.com/schneiel/ImageManagerGo/internal/di"
)

// MockConfigLoader implements config.Loader for testing.
type MockConfigLoader struct {
	LoadFunc func(path string) (*config.Config, error)
}

func (m *MockConfigLoader) Load(path string) (*config.Config, error) {
	if m.LoadFunc != nil {
		return m.LoadFunc(path)
	}
	return &config.Config{}, nil
}

// MockLocalizer is defined in defaultargumentparser_test.go

func TestNewConfigBuilder(t *testing.T) {
	t.Parallel()

	builder := di.NewConfigBuilder()
	require.NotNil(t, builder)
	assert.IsType(t, &di.ConfigBuilder{}, builder)
}

func TestConfigBuilder_Build_RequiresParameters(t *testing.T) {
	t.Parallel()

	builder := di.NewConfigBuilder()
	_, err := builder.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Build() requires parameters")
}

func TestConfigBuilder_BuildConfig_Success(t *testing.T) {
	t.Parallel()

	builder := di.NewConfigBuilder()
	args := []string{"program", "--config", "test.yaml", "command"}

	// Create mock core dependencies
	core := &di.CoreDependencies{
		Parser:       config.NewYAMLParser(),
		FileSystem:   filesystem.NewDefaultFileSystem(),
		TimeProvider: time.NewDefaultTimeProvider(),
		Localizer:    &MockLocalizer{},
	}

	// Create a mock file reader that returns valid config data
	mockFileReader := &MockFileReader{
		ReadFunc: func(_ string) ([]byte, error) {
			return []byte(`
sorter:
  source: "/test/source"
  destination: "/test/dest"
  log: "sorter.log"
deduplicator:
  source: "/test/source"
  log: "dedup.log"
`), nil
		},
	}
	core.FileReader = mockFileReader

	cfg, err := builder.BuildConfig(args, core)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "/test/source", cfg.Sorter.Source)
	assert.Equal(t, "/test/dest", cfg.Sorter.Destination)
}

func TestConfigBuilder_BuildConfig_NoConfigPath(t *testing.T) {
	t.Parallel()

	builder := di.NewConfigBuilder()
	args := []string{"program", "command"} // No config path

	// Create mock core dependencies
	core := &di.CoreDependencies{
		Parser:       config.NewYAMLParser(),
		FileSystem:   filesystem.NewDefaultFileSystem(),
		TimeProvider: time.NewDefaultTimeProvider(),
		Localizer:    &MockLocalizer{},
	}

	// When no config path is provided, it should return default config
	mockFileReader := &MockFileReader{
		ReadFunc: func(_ string) ([]byte, error) {
			// Should not be called when path is empty
			return nil, assert.AnError
		},
	}
	core.FileReader = mockFileReader

	cfg, err := builder.BuildConfig(args, core)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	// Should return default config when no path is provided
}

func TestConfigBuilder_BuildConfig_FileReadError(t *testing.T) {
	t.Parallel()

	builder := di.NewConfigBuilder()
	args := []string{"program", "--config", "nonexistent.yaml", "command"}

	// Create mock core dependencies
	core := &di.CoreDependencies{
		Parser:       config.NewYAMLParser(),
		FileSystem:   filesystem.NewDefaultFileSystem(),
		TimeProvider: time.NewDefaultTimeProvider(),
		Localizer:    &MockLocalizer{},
	}

	// Create a mock file reader that returns an error
	mockFileReader := &MockFileReader{
		ReadFunc: func(_ string) ([]byte, error) {
			return nil, assert.AnError
		},
	}
	core.FileReader = mockFileReader

	_, err := builder.BuildConfig(args, core)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load configuration")
	assert.Contains(t, err.Error(), "nonexistent.yaml")
}

func TestConfigBuilder_BuildConfig_InvalidYAML(t *testing.T) {
	t.Parallel()

	builder := di.NewConfigBuilder()
	args := []string{"program", "--config", "invalid.yaml", "command"}

	// Create mock core dependencies
	core := &di.CoreDependencies{
		Parser:       config.NewYAMLParser(),
		FileSystem:   filesystem.NewDefaultFileSystem(),
		TimeProvider: time.NewDefaultTimeProvider(),
		Localizer:    &MockLocalizer{},
	}

	// Create a mock file reader that returns invalid YAML
	mockFileReader := &MockFileReader{
		ReadFunc: func(_ string) ([]byte, error) {
			return []byte("invalid: yaml: content: ["), nil
		},
	}
	core.FileReader = mockFileReader

	_, err := builder.BuildConfig(args, core)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load configuration")
}

// MockFileReader implements config.FileReader for testing.
type MockFileReader struct {
	ReadFunc func(path string) ([]byte, error)
}

func (m *MockFileReader) ReadFile(path string) ([]byte, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(path)
	}
	return []byte{}, nil
}

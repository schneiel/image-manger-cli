package di

import (
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/config"
)

// ConfigBuilder handles configuration loading and validation.
type ConfigBuilder struct {
	argParser *ArgumentParser
}

// NewConfigBuilder creates a new configuration builder.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		argParser: NewArgumentParser(),
	}
}

// BuildConfig loads and validates the application configuration.
func (cb *ConfigBuilder) BuildConfig(args []string, core *CoreDependencies) (*config.Config, error) {
	configLoader := config.NewDefaultConfigLoader(
		core.FileReader,
		core.Parser,
		core.Localizer,
	)

	configPath := cb.argParser.ExtractConfigPath(args)
	cfg, err := configLoader.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration from path %q: %w", configPath, err)
	}

	return cfg, nil
}

// Build is an alias for BuildConfig to satisfy the standard builder pattern.
// This method provides a simplified interface but requires parameters.
func (cb *ConfigBuilder) Build() (*config.Config, error) {
	// For the standard Build() method, we need to handle the lack of parameters.
	return nil, errors.New("Build() requires parameters - use BuildConfig(args, core) instead")
}

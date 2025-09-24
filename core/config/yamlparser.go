package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// YAMLParser implements Parser for YAML files.
type YAMLParser struct{}

// NewYAMLParser creates a new YAML parser.
func NewYAMLParser() Parser {
	return &YAMLParser{}
}

// Parse parses YAML data into a Config struct.
func (p *YAMLParser) Parse(data []byte) (*Config, error) {
	// Start with default configuration
	config := *DefaultConfig()

	// Unmarshal YAML data on top of defaults
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML config: %w", err)
	}
	return &config, nil
}

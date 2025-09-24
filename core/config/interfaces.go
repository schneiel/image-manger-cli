package config

// FileReader defines the interface for reading files.
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// Parser defines the interface for parsing configuration data.
type Parser interface {
	Parse(data []byte) (*Config, error)
}

// Loader defines the interface for loading configuration.
type Loader interface {
	Load(filename string) (*Config, error)
}

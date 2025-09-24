package config

// DateConfig specifies the strategies for extracting a timestamp from an image file.
type DateConfig struct {
	StrategyOrder  []string     `yaml:"strategyOrder"`
	ExifStrategies []ExifConfig `yaml:"exifStrategies"`
}

// DefaultDateConfig returns a DateConfig instance with default values.
func DefaultDateConfig() DateConfig {
	return DateConfig{
		StrategyOrder:  []string{"exif", "modTime", "creationTime"},
		ExifStrategies: DefaultExifConfig(),
	}
}

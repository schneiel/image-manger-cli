package config

// ExifConfig specifies the strategies for extracting a timestamp from an image file.
type ExifConfig struct {
	FieldName string `yaml:"fieldName"`
	Layout    string `yaml:"layout"`
}

// DefaultExifConfig returns a ExifConfig instance with default values.
func DefaultExifConfig() []ExifConfig {
	return []ExifConfig{
		{FieldName: "DateTimeOriginal", Layout: "2006:01:02 15:04:05"},
		{FieldName: "DateTime", Layout: "2006:01:02 15:04:05"},
		{FieldName: "SubSecDateTimeOriginal", Layout: "2006:01:02 15:04:05.00"},
		{FieldName: "GPSDateStamp", Layout: "2006:01:02"},
	}
}

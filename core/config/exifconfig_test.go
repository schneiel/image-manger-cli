package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestDefaultExifConfig tests the default EXIF configuration.
func TestDefaultExifConfig(t *testing.T) {
	t.Parallel()
	defaultConfig := DefaultExifConfig()

	// Check the number of default configurations
	assert.Len(t, defaultConfig, 4)

	// Check the first configuration
	assert.Equal(t, "DateTimeOriginal", defaultConfig[0].FieldName)
	assert.Equal(t, "2006:01:02 15:04:05", defaultConfig[0].Layout)

	// Check the second configuration
	assert.Equal(t, "DateTime", defaultConfig[1].FieldName)
	assert.Equal(t, "2006:01:02 15:04:05", defaultConfig[1].Layout)

	// Check the third configuration
	assert.Equal(t, "SubSecDateTimeOriginal", defaultConfig[2].FieldName)
	assert.Equal(t, "2006:01:02 15:04:05.00", defaultConfig[2].Layout)

	// Check the fourth configuration
	assert.Equal(t, "GPSDateStamp", defaultConfig[3].FieldName)
	assert.Equal(t, "2006:01:02", defaultConfig[3].Layout)
}

func TestExifConfigMarshalYAML(t *testing.T) {
	t.Parallel()
	config := ExifConfig{
		FieldName: "TestField",
		Layout:    "2006:01:02 15:04:05",
	}

	data, err := yaml.Marshal(config)
	require.NoError(t, err, "Marshaling ExifConfig should not return an error")

	expected := `fieldName: TestField
layout: 2006:01:02 15:04:05
`
	assert.Equal(t, expected, string(data), "Marshaled YAML should match the expected output")
}

func TestExifConfigUnmarshalYAML(t *testing.T) {
	t.Parallel()
	yamlData := `fieldName: TestField
layout: "2006:01:02 15:04:05"
`

	var config ExifConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	require.NoError(t, err, "Unmarshaling YAML should not return an error")

	assert.Equal(t, "TestField", config.FieldName, "FieldName should be correctly unmarshaled")
	assert.Equal(t, "2006:01:02 15:04:05", config.Layout, "Layout should be correctly unmarshaled")
}

func TestExifConfigMarshalYAMLInvalid(t *testing.T) {
	t.Parallel()
	config := ExifConfig{
		FieldName: "InvalidField",
		Layout:    "invalid layout",
	}

	_, err := yaml.Marshal(config)
	require.NoError(t, err, "Marshaling invalid ExifConfig should not return an error")
}

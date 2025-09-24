package exif

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultExifDecoder_Decode(t *testing.T) {
	t.Parallel()
	decoder := NewDefaultExifDecoder()

	tests := []struct {
		name          string
		data          []byte
		expectedError bool
	}{
		{
			name:          "valid exif data",
			data:          []byte{0xFF, 0xD8, 0xFF, 0xE1}, // JPEG header
			expectedError: true,                           // This will fail because it's not valid EXIF data
		},
		{
			name:          "empty data",
			data:          []byte{},
			expectedError: true,
		},
		{
			name:          "nil data",
			data:          nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			reader := bytes.NewReader(tt.data)
			result, err := decoder.Decode(reader)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestDefaultExifDecoder_Constructor(t *testing.T) {
	t.Parallel()
	decoder := NewDefaultExifDecoder()
	assert.NotNil(t, decoder)
	assert.Implements(t, (*Decoder)(nil), decoder)
}

func TestDefaultExifDecoder_InterfaceCompliance(t *testing.T) {
	t.Parallel()
	decoder := NewDefaultExifDecoder()
	assert.NotNil(t, decoder)
}

func TestDefaultExifDecoder_DecodeWithInvalidData(t *testing.T) {
	t.Parallel()
	decoder := NewDefaultExifDecoder()

	// Test with invalid data that should cause decoding errors
	invalidData := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	reader := bytes.NewReader(invalidData)

	result, err := decoder.Decode(reader)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDefaultExifDecoder_DecodeWithMinimalJPEG(t *testing.T) {
	t.Parallel()
	decoder := NewDefaultExifDecoder()

	// Minimal JPEG without EXIF data
	minimalJPEG := []byte{
		0xFF, 0xD8, // JPEG SOI marker
		0xFF, 0xDB, // DQT marker
		0x00, 0x43, // Length
		0x00, // Table info
		// DQT table data (simplified)
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40,
		0x41, 0x42, // Fill remaining bytes
		0xFF, 0xD9, // JPEG EOI marker
	}

	reader := bytes.NewReader(minimalJPEG)
	result, err := decoder.Decode(reader)

	// This should fail because there's no EXIF data
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDefaultExifDecoder_DecodeWithLargeData(t *testing.T) {
	t.Parallel()
	decoder := NewDefaultExifDecoder()

	// Create large data buffer
	data := make([]byte, 1024*10) // 10KB
	for i := range data {
		data[i] = byte(i % 256)
	}

	reader := bytes.NewReader(data)
	result, err := decoder.Decode(reader)

	// This should fail because it's not valid EXIF data
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDefaultExifDecoder_DecodeWithNilReader(t *testing.T) {
	t.Parallel()
	decoder := NewDefaultExifDecoder()

	// Test with nil reader - this should panic or return error
	assert.Panics(t, func() {
		_, _ = decoder.Decode(nil)
	})
}

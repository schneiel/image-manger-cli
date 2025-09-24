package exif

import (
	"fmt"
	"io"

	"github.com/rwcarlsen/goexif/exif"
)

// DefaultExifDecoder implements Decoder using the GoExif package.
type DefaultExifDecoder struct{}

// NewDefaultExifDecoder creates a new DefaultExifDecoder.
func NewDefaultExifDecoder() Decoder {
	return DefaultExifDecoder{}
}

// Decode decodes EXIF data from an io.Reader.
func (d DefaultExifDecoder) Decode(r io.Reader) (*exif.Exif, error) {
	exifData, err := exif.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EXIF data: %w", err)
	}
	return exifData, nil
}

package exif

import (
	"io"

	"github.com/rwcarlsen/goexif/exif"
)

// Decoder interface for decoding EXIF data (dependency injection).
type Decoder interface {
	Decode(r io.Reader) (*exif.Exif, error)
}

// Reader interface for reading EXIF metadata from files.
type Reader interface {
	ReadExif(filePath string) (map[string]interface{}, error)
}

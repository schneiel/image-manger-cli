package hash

import (
	"errors"
	"image"

	"github.com/corona10/goimagehash"
)

// ErrNilImage is returned when a nil image is passed to hash functions.
var ErrNilImage = errors.New("image cannot be nil")

// HashProvider abstracts hash operations for better testability
//
//nolint:revive // HashProvider name is clear and intentional despite package stuttering
type HashProvider interface {
	// DifferenceHash computes the difference hash of an image
	DifferenceHash(img image.Image) (*goimagehash.ImageHash, error)
	// AverageHash computes the average hash of an image
	AverageHash(img image.Image) (*goimagehash.ImageHash, error)
	// PerceptionHash computes the perception hash of an image
	PerceptionHash(img image.Image) (*goimagehash.ImageHash, error)
}

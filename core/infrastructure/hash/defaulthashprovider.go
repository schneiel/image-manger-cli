// Package hash provides perceptual hashing functionality for images.
package hash

import (
	"fmt"
	stdimage "image"

	"github.com/corona10/goimagehash"
)

// DefaultHashProvider implements the HashProvider interface using the goimagehash library.
type DefaultHashProvider struct{}

// NewDefaultHashProvider creates a new instance of DefaultHashProvider.
func NewDefaultHashProvider() HashProvider {
	return &DefaultHashProvider{}
}

// DifferenceHash computes the difference hash of an image.
// The difference hash is computed by comparing adjacent pixels and is useful for finding similar images.
func (h *DefaultHashProvider) DifferenceHash(img stdimage.Image) (*goimagehash.ImageHash, error) {
	if img == nil {
		return nil, ErrNilImage
	}
	hash, err := goimagehash.DifferenceHash(img)
	if err != nil {
		return nil, fmt.Errorf("failed to compute difference hash: %w", err)
	}
	return hash, nil
}

// AverageHash computes the average hash of an image.
// The average hash is computed by comparing each pixel to the average and is useful for finding similar images.
func (h *DefaultHashProvider) AverageHash(img stdimage.Image) (*goimagehash.ImageHash, error) {
	if img == nil {
		return nil, ErrNilImage
	}
	hash, err := goimagehash.AverageHash(img)
	if err != nil {
		return nil, fmt.Errorf("failed to compute average hash: %w", err)
	}
	return hash, nil
}

// PerceptionHash computes the perception hash of an image.
// The perception hash uses DCT (Discrete Cosine Transform) and is more robust to image modifications.
func (h *DefaultHashProvider) PerceptionHash(img stdimage.Image) (*goimagehash.ImageHash, error) {
	if img == nil {
		return nil, ErrNilImage
	}
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, fmt.Errorf("failed to compute perception hash: %w", err)
	}
	return hash, nil
}

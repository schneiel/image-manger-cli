package hash

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestImage creates a simple test image for testing purposes.
func createTestImage(width, height int, fillColor color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, fillColor)
		}
	}
	return img
}

// createGradientImage creates a gradient test image.
func createGradientImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a simple gradient
			grayValue := (x + y) * 255 / (width + height)
			if grayValue > 255 {
				grayValue = 255
			}
			// Safe conversion: grayValue is already clamped to 0-255 range above
			gray := uint8(grayValue) // #nosec G115
			img.Set(x, y, color.RGBA{gray, gray, gray, 255})
		}
	}
	return img
}

func TestNewDefaultHashProvider(t *testing.T) {
	t.Parallel()

	provider := NewDefaultHashProvider()

	assert.NotNil(t, provider)
	assert.IsType(t, &DefaultHashProvider{}, provider)
}

func TestDefaultHashProvider_DifferenceHash(t *testing.T) {
	t.Parallel()

	provider := NewDefaultHashProvider()

	t.Run("successful_hash_calculation", func(t *testing.T) {
		t.Parallel()

		img := createTestImage(64, 64, color.RGBA{255, 0, 0, 255})

		hash, err := provider.DifferenceHash(img)

		require.NoError(t, err)
		assert.NotNil(t, hash)
	})

	t.Run("nil_image_error", func(t *testing.T) {
		t.Parallel()

		hash, err := provider.DifferenceHash(nil)

		require.Error(t, err)
		assert.Equal(t, ErrNilImage, err)
		assert.Nil(t, hash)
	})

	t.Run("different_images_produce_different_hashes", func(t *testing.T) {
		t.Parallel()

		// Use gradient images which will produce different hashes
		img1 := createGradientImage(64, 64)

		// Create a different gradient (reversed)
		img2 := image.NewRGBA(image.Rect(0, 0, 64, 64))
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				// Reversed gradient
				grayValue := 255 - (x+y)*255/(64+64)
				if grayValue < 0 {
					grayValue = 0
				} else if grayValue > 255 {
					grayValue = 255
				}
				// Safe conversion: grayValue is already clamped to 0-255 range above
				gray := uint8(grayValue) // #nosec G115
				img2.Set(x, y, color.RGBA{gray, gray, gray, 255})
			}
		}

		hash1, err1 := provider.DifferenceHash(img1)
		hash2, err2 := provider.DifferenceHash(img2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotNil(t, hash1)
		assert.NotNil(t, hash2)

		// Hashes should be different for different images
		distance, err := hash1.Distance(hash2)
		require.NoError(t, err)
		assert.Positive(t, distance)
	})

	t.Run("same_image_produces_same_hash", func(t *testing.T) {
		t.Parallel()

		img := createTestImage(64, 64, color.RGBA{255, 0, 0, 255})

		hash1, err1 := provider.DifferenceHash(img)
		hash2, err2 := provider.DifferenceHash(img)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotNil(t, hash1)
		assert.NotNil(t, hash2)

		// Same image should produce identical hashes
		distance, err := hash1.Distance(hash2)
		require.NoError(t, err)
		assert.Equal(t, 0, distance)
	})
}

func TestDefaultHashProvider_AverageHash(t *testing.T) {
	t.Parallel()

	provider := NewDefaultHashProvider()

	t.Run("successful_hash_calculation", func(t *testing.T) {
		t.Parallel()

		img := createTestImage(64, 64, color.RGBA{128, 128, 128, 255})

		hash, err := provider.AverageHash(img)

		require.NoError(t, err)
		assert.NotNil(t, hash)
	})

	t.Run("nil_image_error", func(t *testing.T) {
		t.Parallel()

		hash, err := provider.AverageHash(nil)

		require.Error(t, err)
		assert.Equal(t, ErrNilImage, err)
		assert.Nil(t, hash)
	})

	t.Run("gradient_image_hash", func(t *testing.T) {
		t.Parallel()

		img := createGradientImage(64, 64)

		hash, err := provider.AverageHash(img)

		require.NoError(t, err)
		assert.NotNil(t, hash)
	})

	t.Run("different_images_produce_different_hashes", func(t *testing.T) {
		t.Parallel()

		// Create two gradient images with different patterns
		img1 := createGradientImage(64, 64)
		img2 := createTestImage(64, 64, color.RGBA{255, 255, 255, 255}) // Solid white

		hash1, err1 := provider.AverageHash(img1)
		hash2, err2 := provider.AverageHash(img2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotNil(t, hash1)
		assert.NotNil(t, hash2)

		// Gradient vs solid should produce different hashes
		distance, err := hash1.Distance(hash2)
		require.NoError(t, err)
		assert.Positive(t, distance)
	})
}

func TestDefaultHashProvider_PerceptionHash(t *testing.T) {
	t.Parallel()

	provider := NewDefaultHashProvider()

	t.Run("successful_hash_calculation", func(t *testing.T) {
		t.Parallel()

		img := createGradientImage(64, 64)

		hash, err := provider.PerceptionHash(img)

		require.NoError(t, err)
		assert.NotNil(t, hash)
	})

	t.Run("nil_image_error", func(t *testing.T) {
		t.Parallel()

		hash, err := provider.PerceptionHash(nil)

		require.Error(t, err)
		assert.Equal(t, ErrNilImage, err)
		assert.Nil(t, hash)
	})

	t.Run("complex_image_hash", func(t *testing.T) {
		t.Parallel()

		// Create a more complex image with patterns
		img := image.NewRGBA(image.Rect(0, 0, 64, 64))
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				// Create a checkerboard pattern
				if (x/8+y/8)%2 == 0 {
					img.Set(x, y, color.RGBA{255, 255, 255, 255})
				} else {
					img.Set(x, y, color.RGBA{0, 0, 0, 255})
				}
			}
		}

		hash, err := provider.PerceptionHash(img)

		require.NoError(t, err)
		assert.NotNil(t, hash)
	})

	t.Run("perception_hash_robustness", func(t *testing.T) {
		t.Parallel()

		// Create two similar images (one slightly modified)
		img1 := createGradientImage(64, 64)

		// Create a slightly modified version
		img2 := image.NewRGBA(image.Rect(0, 0, 64, 64))
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				// Same gradient but with slight noise
				grayValue := (x + y) * 255 / (64 + 64)
				if grayValue > 255 {
					grayValue = 255
				}
				// Safe conversion: grayValue is already clamped to 0-255 range above
				gray := uint8(grayValue) // #nosec G115
				// Add small random variation
				if (x+y)%10 == 0 {
					gray += 5 // Small modification
				}
				img2.Set(x, y, color.RGBA{gray, gray, gray, 255})
			}
		}

		hash1, err1 := provider.PerceptionHash(img1)
		hash2, err2 := provider.PerceptionHash(img2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotNil(t, hash1)
		assert.NotNil(t, hash2)

		// Perception hash should be relatively robust to small changes
		distance, err := hash1.Distance(hash2)
		require.NoError(t, err)
		// Distance should be small for similar images
		assert.Less(t, distance, 20) // Threshold for "similar" images
	})
}

func TestDefaultHashProvider_AllHashTypes(t *testing.T) {
	t.Parallel()

	provider := NewDefaultHashProvider()
	img := createGradientImage(64, 64)

	t.Run("all_hash_types_work", func(t *testing.T) {
		t.Parallel()

		dhash, err1 := provider.DifferenceHash(img)
		ahash, err2 := provider.AverageHash(img)
		phash, err3 := provider.PerceptionHash(img)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		assert.NotNil(t, dhash)
		assert.NotNil(t, ahash)
		assert.NotNil(t, phash)
	})

	t.Run("different_hash_types_produce_different_results", func(t *testing.T) {
		t.Parallel()

		dhash, err1 := provider.DifferenceHash(img)
		ahash, err2 := provider.AverageHash(img)
		phash, err3 := provider.PerceptionHash(img)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)

		// Different hash algorithms should generally produce different results
		// (though they might occasionally be the same for certain images)
		dHashStr := dhash.ToString()
		aHashStr := ahash.ToString()
		pHashStr := phash.ToString()

		assert.NotEmpty(t, dHashStr)
		assert.NotEmpty(t, aHashStr)
		assert.NotEmpty(t, pHashStr)

		// At least one should be different (they usually all are)
		different := (dHashStr != aHashStr) || (aHashStr != pHashStr) || (dHashStr != pHashStr)
		assert.True(t, different, "Hash algorithms should produce different results for most images")
	})
}

func TestDefaultHashProvider_EdgeCases(t *testing.T) {
	t.Parallel()

	provider := NewDefaultHashProvider()

	t.Run("single_pixel_image", func(t *testing.T) {
		t.Parallel()

		img := createTestImage(1, 1, color.RGBA{255, 0, 0, 255})

		dhash, err1 := provider.DifferenceHash(img)
		ahash, err2 := provider.AverageHash(img)
		phash, err3 := provider.PerceptionHash(img)

		// All should work even with tiny images
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		assert.NotNil(t, dhash)
		assert.NotNil(t, ahash)
		assert.NotNil(t, phash)
	})

	t.Run("large_image", func(t *testing.T) {
		t.Parallel()

		img := createTestImage(512, 512, color.RGBA{128, 128, 128, 255})

		dhash, err1 := provider.DifferenceHash(img)
		ahash, err2 := provider.AverageHash(img)
		phash, err3 := provider.PerceptionHash(img)

		// Should handle larger images without issues
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		assert.NotNil(t, dhash)
		assert.NotNil(t, ahash)
		assert.NotNil(t, phash)
	})

	t.Run("transparent_image", func(t *testing.T) {
		t.Parallel()

		img := createTestImage(64, 64, color.RGBA{255, 0, 0, 128}) // Semi-transparent

		dhash, err1 := provider.DifferenceHash(img)
		ahash, err2 := provider.AverageHash(img)
		phash, err3 := provider.PerceptionHash(img)

		// Should handle transparency
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		assert.NotNil(t, dhash)
		assert.NotNil(t, ahash)
		assert.NotNil(t, phash)
	})
}

func BenchmarkDefaultHashProvider_DifferenceHash(b *testing.B) {
	provider := NewDefaultHashProvider()
	img := createGradientImage(64, 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = provider.DifferenceHash(img)
	}
}

func BenchmarkDefaultHashProvider_AverageHash(b *testing.B) {
	provider := NewDefaultHashProvider()
	img := createGradientImage(64, 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = provider.AverageHash(img)
	}
}

func BenchmarkDefaultHashProvider_PerceptionHash(b *testing.B) {
	provider := NewDefaultHashProvider()
	img := createGradientImage(64, 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = provider.PerceptionHash(img)
	}
}

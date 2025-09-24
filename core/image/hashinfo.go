package image

import "github.com/corona10/goimagehash"

// HashInfo stores the path and perceptual hash of an image.
type HashInfo struct {
	FilePath string
	Hash     *goimagehash.ImageHash
}

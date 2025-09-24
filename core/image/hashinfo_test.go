package image

import (
	"testing"

	"github.com/corona10/goimagehash"
	"github.com/stretchr/testify/assert"
)

func TestHashInfo_NewHashInfo(t *testing.T) {
	t.Parallel()
	filePath := "/path/to/image.jpg"
	// Create a mock ImageHash (we'll use nil for simplicity in tests)
	var imageHash *goimagehash.ImageHash

	hashInfo := HashInfo{
		FilePath: filePath,
		Hash:     imageHash,
	}

	assert.Equal(t, filePath, hashInfo.FilePath)
	assert.Equal(t, imageHash, hashInfo.Hash)
}

func TestHashInfo_ZeroValue(t *testing.T) {
	t.Parallel()
	var hashInfo HashInfo

	assert.Empty(t, hashInfo.FilePath)
	assert.Nil(t, hashInfo.Hash)
}

func TestHashInfo_FieldAssignment(t *testing.T) {
	t.Parallel()
	hashInfo := HashInfo{}

	// Test individual field assignment
	hashInfo.FilePath = "/new/path/test.png"
	hashInfo.Hash = nil // Using nil for simplicity in tests

	assert.Equal(t, "/new/path/test.png", hashInfo.FilePath)
	assert.Nil(t, hashInfo.Hash)
}

func TestHashInfo_HashComparison(t *testing.T) {
	t.Parallel()
	hash1 := HashInfo{
		FilePath: "/path1/image1.jpg",
		Hash:     nil,
	}

	hash2 := HashInfo{
		FilePath: "/path2/image2.jpg",
		Hash:     nil,
	}

	// Test that both nil hashes are equal
	assert.Equal(t, hash1.Hash, hash2.Hash)

	// Test different file paths
	assert.NotEqual(t, hash1.FilePath, hash2.FilePath)
}

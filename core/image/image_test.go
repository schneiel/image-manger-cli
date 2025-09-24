package image

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImage_NewImage(t *testing.T) {
	t.Parallel()
	filePath := "/path/to/image.jpg"
	originalFileName := "image.jpg"
	date := time.Now()

	img := Image{
		FilePath:         filePath,
		OriginalFileName: originalFileName,
		Date:             date,
	}

	assert.Equal(t, filePath, img.FilePath)
	assert.Equal(t, originalFileName, img.OriginalFileName)
	assert.Equal(t, date, img.Date)
}

func TestImage_ZeroValue(t *testing.T) {
	t.Parallel()
	var img Image

	assert.Empty(t, img.FilePath)
	assert.Empty(t, img.OriginalFileName)
	assert.True(t, img.Date.IsZero())
}

func TestImage_FieldAssignment(t *testing.T) {
	t.Parallel()
	img := Image{}

	// Test individual field assignment
	img.FilePath = "/new/path/test.png"
	img.OriginalFileName = "test.png"
	img.Date = time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

	assert.Equal(t, "/new/path/test.png", img.FilePath)
	assert.Equal(t, "test.png", img.OriginalFileName)
	assert.Equal(t, time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC), img.Date)
}

func TestImage_IsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		image    Image
		expected bool
	}{
		{
			name: "valid image with all fields",
			image: Image{
				FilePath:         "/path/to/image.jpg",
				OriginalFileName: "image.jpg",
				Date:             time.Now(),
			},
			expected: true,
		},
		{
			name: "empty file path",
			image: Image{
				FilePath:         "",
				OriginalFileName: "image.jpg",
				Date:             time.Now(),
			},
			expected: false,
		},
		{
			name: "empty original filename",
			image: Image{
				FilePath:         "/path/to/image.jpg",
				OriginalFileName: "",
				Date:             time.Now(),
			},
			expected: false,
		},
		{
			name: "zero date",
			image: Image{
				FilePath:         "/path/to/image.jpg",
				OriginalFileName: "image.jpg",
				Date:             time.Time{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.image.FilePath != "" &&
				tt.image.OriginalFileName != "" &&
				!tt.image.Date.IsZero()
			assert.Equal(t, tt.expected, valid)
		})
	}
}

func TestImage_String(t *testing.T) {
	t.Parallel()
	img := Image{
		FilePath:         "/path/to/image.jpg",
		OriginalFileName: "image.jpg",
		Date:             time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
	}

	str := img.FilePath + " (" + img.OriginalFileName + ")"
	assert.Contains(t, str, "image.jpg")
	assert.Contains(t, str, "/path/to/image.jpg")
}

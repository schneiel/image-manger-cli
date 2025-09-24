package dedup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicateGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		group    DuplicateGroup
		expected []string
	}{
		{
			name:     "empty group",
			group:    DuplicateGroup{},
			expected: []string{},
		},
		{
			name:     "single file group",
			group:    DuplicateGroup{"/path/to/image1.jpg"},
			expected: []string{"/path/to/image1.jpg"},
		},
		{
			name: "multiple files group",
			group: DuplicateGroup{
				"/path/to/image1.jpg",
				"/path/to/image2.jpg",
				"/path/to/image3.jpg",
			},
			expected: []string{
				"/path/to/image1.jpg",
				"/path/to/image2.jpg",
				"/path/to/image3.jpg",
			},
		},
		{
			name: "group with duplicate paths",
			group: DuplicateGroup{
				"/path/to/image1.jpg",
				"/path/to/image1.jpg",
				"/path/to/image2.jpg",
			},
			expected: []string{
				"/path/to/image1.jpg",
				"/path/to/image1.jpg",
				"/path/to/image2.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, []string(tt.group))
			assert.Len(t, tt.group, len(tt.expected))
		})
	}
}

func TestDuplicateGroup_Operations(t *testing.T) {
	t.Parallel()

	t.Run("append to group", func(t *testing.T) {
		t.Parallel()
		group := DuplicateGroup{"/path/to/image1.jpg"}
		group = append(group, "/path/to/image2.jpg")

		expected := DuplicateGroup{"/path/to/image1.jpg", "/path/to/image2.jpg"}
		assert.Equal(t, expected, group)
	})

	t.Run("access by index", func(t *testing.T) {
		t.Parallel()
		group := DuplicateGroup{"/path/to/image1.jpg", "/path/to/image2.jpg"}

		assert.Equal(t, "/path/to/image1.jpg", group[0])
		assert.Equal(t, "/path/to/image2.jpg", group[1])
	})

	t.Run("length check", func(t *testing.T) {
		t.Parallel()
		group := DuplicateGroup{"/path/to/image1.jpg", "/path/to/image2.jpg"}

		assert.Len(t, group, 2)
	})

	t.Run("iteration", func(t *testing.T) {
		t.Parallel()
		group := DuplicateGroup{"/path/to/image1.jpg", "/path/to/image2.jpg"}

		var paths []string
		for _, path := range group {
			paths = append(paths, path)
		}

		assert.Equal(t, []string{"/path/to/image1.jpg", "/path/to/image2.jpg"}, paths)
	})
}

func TestFileGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		group    FileGroup
		expected [][]string
	}{
		{
			name:     "empty file group",
			group:    FileGroup{},
			expected: [][]string{},
		},
		{
			name:     "single group with single file",
			group:    FileGroup{{"file1.jpg"}},
			expected: [][]string{{"file1.jpg"}},
		},
		{
			name: "single group with multiple files",
			group: FileGroup{
				{"file1.jpg", "file2.jpg", "file3.jpg"},
			},
			expected: [][]string{
				{"file1.jpg", "file2.jpg", "file3.jpg"},
			},
		},
		{
			name: "multiple groups",
			group: FileGroup{
				{"file1.jpg", "file2.jpg"},
				{"file3.jpg", "file4.jpg", "file5.jpg"},
				{"file6.jpg"},
			},
			expected: [][]string{
				{"file1.jpg", "file2.jpg"},
				{"file3.jpg", "file4.jpg", "file5.jpg"},
				{"file6.jpg"},
			},
		},
		{
			name: "empty groups within file group",
			group: FileGroup{
				{"file1.jpg"},
				{},
				{"file2.jpg"},
			},
			expected: [][]string{
				{"file1.jpg"},
				{},
				{"file2.jpg"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, [][]string(tt.group))
			assert.Len(t, tt.group, len(tt.expected))
		})
	}
}

func TestFileGroup_Operations(t *testing.T) {
	t.Parallel()

	t.Run("append to file group", func(t *testing.T) {
		t.Parallel()
		fileGroup := FileGroup{{"file1.jpg"}}
		fileGroup = append(fileGroup, []string{"file2.jpg", "file3.jpg"})

		expected := FileGroup{
			{"file1.jpg"},
			{"file2.jpg", "file3.jpg"},
		}
		assert.Equal(t, expected, fileGroup)
	})

	t.Run("access by index", func(t *testing.T) {
		t.Parallel()
		fileGroup := FileGroup{
			{"file1.jpg", "file2.jpg"},
			{"file3.jpg"},
		}

		assert.Equal(t, []string{"file1.jpg", "file2.jpg"}, fileGroup[0])
		assert.Equal(t, []string{"file3.jpg"}, fileGroup[1])
	})

	t.Run("nested access", func(t *testing.T) {
		t.Parallel()
		fileGroup := FileGroup{
			{"file1.jpg", "file2.jpg"},
			{"file3.jpg"},
		}

		assert.Equal(t, "file1.jpg", fileGroup[0][0])
		assert.Equal(t, "file2.jpg", fileGroup[0][1])
		assert.Equal(t, "file3.jpg", fileGroup[1][0])
	})

	t.Run("length operations", func(t *testing.T) {
		t.Parallel()
		fileGroup := FileGroup{
			{"file1.jpg", "file2.jpg"},
			{"file3.jpg"},
		}

		assert.Len(t, fileGroup, 2)
		assert.Len(t, fileGroup[0], 2)
		assert.Len(t, fileGroup[1], 1)
	})

	t.Run("iteration", func(t *testing.T) {
		t.Parallel()
		fileGroup := FileGroup{
			{"file1.jpg", "file2.jpg"},
			{"file3.jpg"},
		}

		var allGroups [][]string
		for _, group := range fileGroup {
			allGroups = append(allGroups, group)
		}

		expected := [][]string{
			{"file1.jpg", "file2.jpg"},
			{"file3.jpg"},
		}
		assert.Equal(t, expected, allGroups)
	})

	t.Run("nested iteration", func(t *testing.T) {
		t.Parallel()
		fileGroup := FileGroup{
			{"file1.jpg", "file2.jpg"},
			{"file3.jpg"},
		}

		var allFiles []string
		for _, group := range fileGroup {
			allFiles = append(allFiles, group...)
		}

		expected := []string{"file1.jpg", "file2.jpg", "file3.jpg"}
		assert.Equal(t, expected, allFiles)
	})
}

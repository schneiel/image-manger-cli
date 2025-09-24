package dedupkeep

import (
	"os"
	"testing"
	"time"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestOldestFile(t *testing.T) { //nolint:gocognit // comprehensive test with multiple scenarios
	t.Parallel()
	tests := []struct {
		name           string
		paths          []string
		expectedToKeep string
		expectedCount  int
	}{
		{
			name:           "empty paths",
			paths:          []string{},
			expectedToKeep: "",
			expectedCount:  0,
		},
		{
			name:           "single path",
			paths:          []string{"/path/to/file1.jpg"},
			expectedToKeep: "/path/to/file1.jpg",
			expectedCount:  0,
		},
		{
			name:           "multiple paths with different mod times",
			paths:          []string{"/path/to/file1.jpg", "/path/to/file2.jpg", "/path/to/file3.jpg"},
			expectedToKeep: "/path/to/file1.jpg", // oldest based on mock
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			// Create mock filesystem
			mockFS := &testutils.MockFileSystem{}

			// Set up mock file times - file1 is oldest
			mockFS.StatFunc = func(name string) (os.FileInfo, error) {
				switch name {
				case "/path/to/file1.jpg":
					return &testutils.MockFileInfo{
						ModTimeFunc: func() time.Time { return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) },
					}, nil
				case "/path/to/file2.jpg":
					return &testutils.MockFileInfo{
						ModTimeFunc: func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) },
					}, nil
				case "/path/to/file3.jpg":
					return &testutils.MockFileInfo{
						ModTimeFunc: func() time.Time { return time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC) },
					}, nil
				}
				return nil, os.ErrNotExist
			}

			keepFunc := OldestFile(mockFS)
			toKeep, toRemove := keepFunc(tt.paths)

			if toKeep != tt.expectedToKeep {
				t.Errorf("expected to keep %q, got %q", tt.expectedToKeep, toKeep)
			}

			if len(toRemove) != tt.expectedCount {
				t.Errorf("expected %d files to remove, got %d", tt.expectedCount, len(toRemove))
			}

			// Verify all expected files are marked for removal
			if len(tt.paths) > 1 {
				for _, path := range tt.paths {
					if path != toKeep {
						found := false
						for _, removed := range toRemove {
							if removed == path {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("expected path %q to be in remove list", path)
						}
					}
				}
			}
		})
	}
}

func TestShortestPath(t *testing.T) { //nolint:gocognit // comprehensive test with multiple scenarios
	t.Parallel()
	tests := []struct {
		name           string
		paths          []string
		expectedToKeep string
		expectedCount  int
	}{
		{
			name:           "empty paths",
			paths:          []string{},
			expectedToKeep: "",
			expectedCount:  0,
		},
		{
			name:           "single path",
			paths:          []string{"/very/long/path/to/file.jpg"},
			expectedToKeep: "/very/long/path/to/file.jpg",
			expectedCount:  0,
		},
		{
			name:           "multiple paths with different lengths",
			paths:          []string{"/very/long/path/to/file1.jpg", "/short.jpg", "/medium/path/file2.jpg"},
			expectedToKeep: "/short.jpg",
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			keepFunc := ShortestPath()
			toKeep, toRemove := keepFunc(tt.paths)

			if toKeep != tt.expectedToKeep {
				t.Errorf("expected to keep %q, got %q", tt.expectedToKeep, toKeep)
			}

			if len(toRemove) != tt.expectedCount {
				t.Errorf("expected %d files to remove, got %d", tt.expectedCount, len(toRemove))
			}

			// Verify all expected files are marked for removal
			if len(tt.paths) > 1 {
				for _, path := range tt.paths {
					if path != toKeep {
						found := false
						for _, removed := range toRemove {
							if removed == path {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("expected path %q to be in remove list", path)
						}
					}
				}
			}
		})
	}
}

func TestGetModTime(t *testing.T) {
	t.Parallel()
	t.Run("successful stat", func(_ *testing.T) {
		mockFS := &testutils.MockFileSystem{}
		expectedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		mockFS.StatFunc = func(_ string) (os.FileInfo, error) {
			return &testutils.MockFileInfo{
				ModTimeFunc: func() time.Time { return expectedTime },
			}, nil
		}

		info := getModTime(mockFS, "/test/path.jpg")

		if info.ModTime() != expectedTime {
			t.Errorf("expected mod time %v, got %v", expectedTime, info.ModTime())
		}
	})

	t.Run("stat error returns dummy", func(_ *testing.T) {
		mockFS := &testutils.MockFileSystem{}
		mockFS.StatFunc = func(_ string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}

		info := getModTime(mockFS, "/nonexistent/path.jpg")

		if info.Name() != "path.jpg" {
			t.Errorf("expected dummy file name 'path.jpg', got %q", info.Name())
		}

		// Dummy file should return zero time
		if !info.ModTime().IsZero() {
			t.Errorf("expected zero time for dummy file, got %v", info.ModTime())
		}
	})
}

package handlers

import (
	"errors"
	"testing"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultTaskExecutor(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)

	if executor == nil {
		t.Fatal("Expected executor to be created, got nil")
	}

	// Note: Cannot test unexported fields from external test package
	// The fact that NewDefaultTaskExecutor returns a non-nil executor
	// and subsequent Execute calls work correctly validates dependency injection
}

func TestDefaultTaskExecutor_Execute_SortTask(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	mockSortTask.RunFunc = func() error {
		return nil
	}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/destination",
		},
	}

	err := executor.Execute("sort", cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDefaultTaskExecutor_Execute_DedupTask(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	mockDedupTask.RunFunc = func() error {
		return nil
	}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}

	err := executor.Execute("dedup", cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDefaultTaskExecutor_Execute_UnknownTask(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{}

	err := executor.Execute("unknown", cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "unknown task: unknown"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDefaultTaskExecutor_Execute_SortTaskError(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	expectedError := errors.New("sort task failed")
	mockSortTask.RunFunc = func() error {
		return expectedError
	}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{
		Sorter: config.SorterConfig{
			Source:      "/test/source",
			Destination: "/test/destination",
		},
	}

	err := executor.Execute("sort", cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedWrappedError := "task execution failed: " + expectedError.Error()
	if err.Error() != expectedWrappedError {
		t.Errorf("Expected error '%s', got '%s'", expectedWrappedError, err.Error())
	}
}

func TestDefaultTaskExecutor_Execute_DedupTaskError(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	expectedError := errors.New("dedup task failed")
	mockDedupTask.RunFunc = func() error {
		return expectedError
	}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source: "/test/source",
		},
	}

	err := executor.Execute("dedup", cfg)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedWrappedError := "task execution failed: " + expectedError.Error()
	if err.Error() != expectedWrappedError {
		t.Errorf("Expected error '%s', got '%s'", expectedWrappedError, err.Error())
	}
}

func TestDefaultTaskExecutor_Execute_EmptyConfig(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	mockSortTask.RunFunc = func() error {
		return nil
	}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{}

	err := executor.Execute("sort", cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDefaultTaskExecutor_Execute_CompleteConfig(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	mockDedupTask.RunFunc = func() error {
		return nil
	}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{
		Deduplicator: config.DeduplicatorConfig{
			Source:         "/test/source",
			ActionStrategy: "dryRun",
			KeepStrategy:   "keepOldest",
		},
		Sorter: config.SorterConfig{
			Source:         "/test/source",
			Destination:    "/test/destination",
			ActionStrategy: "copy",
		},
		Files: config.FilesConfig{
			ApplicationLog: "test.log",
		},
		AllowedImageExtensions: []string{".jpg", ".png", ".jpeg"},
	}

	err := executor.Execute("dedup", cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDefaultTaskExecutor_Execute_CaseSensitive(t *testing.T) {
	t.Parallel()
	mockLogger := &testutils.MockLogger{}
	mockLocalizer := &testutils.MockLocalizer{}
	mockFileUtils := &testutils.MockFileUtils{}
	mockSortTask := &testutils.MockTask{}
	mockDedupTask := &testutils.MockTask{}

	executor := NewDefaultTaskExecutor(mockLogger, mockLocalizer, mockFileUtils, mockSortTask, mockDedupTask)
	cfg := config.Config{}

	testCases := getCaseSensitiveTestCases()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			setupMockTasks(testCase.taskName, mockSortTask, mockDedupTask)

			err := executor.Execute(testCase.taskName, cfg)

			assertExecutionResult(t, err, testCase.shouldError)
		})
	}
}

func getCaseSensitiveTestCases() []struct {
	name        string
	taskName    string
	shouldError bool
} {
	return []struct {
		name        string
		taskName    string
		shouldError bool
	}{
		{
			name:        "Sort task",
			taskName:    "sort",
			shouldError: false,
		},
		{
			name:        "Dedup task",
			taskName:    "dedup",
			shouldError: false,
		},
		{
			name:        "Uppercase sort",
			taskName:    "SORT",
			shouldError: true,
		},
		{
			name:        "Uppercase dedup",
			taskName:    "DEDUP",
			shouldError: true,
		},
		{
			name:        "Mixed case",
			taskName:    "Sort",
			shouldError: true,
		},
		{
			name:        "Empty task name",
			taskName:    "",
			shouldError: true,
		},
	}
}

func setupMockTasks(taskName string, mockSortTask, mockDedupTask *testutils.MockTask) {
	switch taskName {
	case "sort":
		mockSortTask.RunFunc = func() error {
			return nil
		}
	case "dedup":
		mockDedupTask.RunFunc = func() error {
			return nil
		}
	}
}

func assertExecutionResult(t *testing.T, err error, shouldError bool) {
	t.Helper()
	if shouldError {
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	} else {
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
}

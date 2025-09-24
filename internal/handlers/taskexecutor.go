package handlers

import (
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	"github.com/schneiel/ImageManagerGo/core/task"
)

// DefaultTaskExecutor implements TaskExecutor using a task factory for proper DI.
type DefaultTaskExecutor struct {
	logger    log.Logger
	localizer i18n.Localizer
	fileUtils filesystem.FileUtils
	sortTask  task.Task
	dedupTask task.Task
}

// NewDefaultTaskExecutor creates a new default task executor with proper dependency injection.
func NewDefaultTaskExecutor(
	logger log.Logger,
	localizer i18n.Localizer,
	fileUtils filesystem.FileUtils,
	sortTask, dedupTask task.Task,
) *DefaultTaskExecutor {
	return &DefaultTaskExecutor{
		logger:    logger,
		localizer: localizer,
		fileUtils: fileUtils,
		sortTask:  sortTask,
		dedupTask: dedupTask,
	}
}

// Execute executes a task with the given name and configuration using injected dependencies.
func (e *DefaultTaskExecutor) Execute(taskName string, _ config.Config) error {
	var taskToRun task.Task

	switch taskName {
	case "sort":
		taskToRun = e.sortTask
	case "dedup":
		taskToRun = e.dedupTask
	default:
		return fmt.Errorf("unknown task: %s", taskName)
	}

	err := taskToRun.Run()
	if err != nil {
		return fmt.Errorf("task execution failed: %w", err)
	}

	return nil
}

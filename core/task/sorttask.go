package task

import (
	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	"github.com/schneiel/ImageManagerGo/core/processing/sort"
)

// SortTask implements the Task interface for image sorting.
type SortTask struct {
	config         *config.Config
	logger         log.Logger
	localizer      i18n.Localizer
	imageProcessor sort.ImageProcessor
	fileUtils      filesystem.FileUtils
	
	// Lazy strategy creation
	strategyFactory func() sort.Strategy
	cachedStrategy  sort.Strategy
	lastConfigValue string // Track config changes
}

// getStrategy returns the action strategy, creating it lazily or rebuilding if config changed
func (t *SortTask) getStrategy() sort.Strategy {
	currentConfigValue := t.config.Sorter.ActionStrategy
	
	// Rebuild strategy if config changed or never built
	if t.cachedStrategy == nil || t.lastConfigValue != currentConfigValue {
		t.cachedStrategy = t.strategyFactory()
		t.lastConfigValue = currentConfigValue
	}
	
	return t.cachedStrategy
}

// Run executes the sort task using injected dependencies.
func (t *SortTask) Run() error {
	images := t.imageProcessor.Process(t.config.Sorter.Source)
	if len(images) == 0 {
		return nil
	}

	// Get strategy (lazy creation with config change detection)
	actionStrategy := t.getStrategy()
	
	if actionStrategy.GetResources() != nil {
		err := actionStrategy.GetResources().Setup()
		if err != nil {
			return err
		}
		defer func() { _ = actionStrategy.GetResources().Teardown() }()
	}

	for _, img := range images {
		_ = actionStrategy.Execute(img.FilePath, t.config.Sorter.Destination)
	}
	return nil
}

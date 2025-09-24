package di

import (
	"embed"
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	coretimepkg "github.com/schneiel/ImageManagerGo/core/infrastructure/time"
)

// CoreBuilder handles building core system dependencies.
type CoreBuilder struct {
	argParser *ArgumentParser
}

// NewCoreBuilder creates a new core dependencies builder.
func NewCoreBuilder() *CoreBuilder {
	return &CoreBuilder{
		argParser: NewArgumentParser(),
	}
}

// BuildCore initializes core system dependencies.
func (cb *CoreBuilder) BuildCore(args []string, localesFS *embed.FS) (*CoreDependencies, error) {
	core := &CoreDependencies{
		Parser:       config.NewYAMLParser(),
		FileSystem:   filesystem.NewDefaultFileSystem(),
		TimeProvider: coretimepkg.NewDefaultTimeProvider(),
	}

	// Initialize file reader
	core.FileReader = config.NewDefaultFileReaderWithFilesystem(core.FileSystem)

	// Initialize localizer
	language := cb.argParser.ExtractLanguage(args)
	localizerConfig := i18n.LocalizerConfig{Language: language, LocalesFS: localesFS}
	localizer, err := i18n.NewBundleLocalizer(localizerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create bundle localizer for language %q: %w", language, err)
	}
	core.Localizer = localizer

	// Initialize file utilities
	fileUtils, err := filesystem.NewFileUtils(core.FileSystem, core.Localizer)
	if err != nil {
		return nil, fmt.Errorf("failed to create file utils: %w", err)
	}
	core.FileUtils = fileUtils

	return core, nil
}

// Build is an alias for BuildCore to satisfy the standard builder pattern
// This method provides a simplified interface for building core dependencies with default parameters.
func (cb *CoreBuilder) Build() (*CoreDependencies, error) {
	// For the standard Build() method, we need to handle the lack of parameters
	// This would typically be used in contexts where args and localesFS are available elsewhere
	return nil, errors.New("Build() requires parameters - use BuildCore(args, localesFS) instead")
}

// CoreDependencies holds the basic system dependencies.
type CoreDependencies struct {
	Parser       config.Parser
	FileReader   config.FileReader
	FileSystem   filesystem.FileSystem
	TimeProvider coretimepkg.TimeProvider
	Localizer    i18n.Localizer
	FileUtils    filesystem.FileUtils
}

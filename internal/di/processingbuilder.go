package di

import (
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/date"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/exif"
	"github.com/schneiel/ImageManagerGo/core/processing/dedup"
	"github.com/schneiel/ImageManagerGo/core/processing/sort"
	datestrategies "github.com/schneiel/ImageManagerGo/core/strategies/date"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupaction"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupkeep"
	"github.com/schneiel/ImageManagerGo/core/strategies/shared"
	"github.com/schneiel/ImageManagerGo/core/strategies/sortaction"
)

// ProcessingBuilder handles image processing pipeline setup.
type ProcessingBuilder struct{}

// NewProcessingBuilder creates a new processing builder.
func NewProcessingBuilder() *ProcessingBuilder {
	return &ProcessingBuilder{}
}

// BuildProcessing sets up image processing dependencies.
func (pb *ProcessingBuilder) BuildProcessing(
	core *CoreDependencies,
	cfg *config.Config,
	logging *LoggingDependencies,
) (*ProcessingDependencies, error) {
	processing := &ProcessingDependencies{}

	err := pb.buildExifComponents(core, processing)
	if err != nil {
		return nil, err
	}

	err = pb.buildDateProcessing(core, cfg, processing)
	if err != nil {
		return nil, err
	}

	err = pb.buildSortComponents(cfg, logging, core, processing)
	if err != nil {
		return nil, err
	}

	err = pb.buildDedupComponents(cfg, logging, core, processing)
	if err != nil {
		return nil, err
	}

	return processing, nil
}

func (pb *ProcessingBuilder) buildExifComponents(
	core *CoreDependencies,
	processing *ProcessingDependencies,
) error {
	processing.ExifDecoder = exif.NewDefaultExifDecoder()
	exifReader, err := exif.NewDefaultReader(
		core.Localizer,
		core.FileSystem,
		processing.ExifDecoder,
	)
	if err != nil {
		return fmt.Errorf("failed to create EXIF reader: %w", err)
	}
	processing.ExifReader = exifReader
	return nil
}

func (pb *ProcessingBuilder) buildDateProcessing(
	core *CoreDependencies,
	cfg *config.Config,
	processing *ProcessingDependencies,
) error {
	dateProcessor, err := pb.buildDateProcessor(core, cfg)
	if err != nil {
		return err
	}
	processing.DateProcessor = dateProcessor
	return nil
}

func (pb *ProcessingBuilder) buildSortComponents(
	cfg *config.Config,
	logging *LoggingDependencies,
	core *CoreDependencies,
	processing *ProcessingDependencies,
) error {
	processing.ImageFinder = sort.NewDefaultImageFinder(cfg.AllowedImageExtensions)
	imageAnalyzer, err := sort.NewDefaultImageAnalyzer(
		processing.DateProcessor,
		processing.ExifReader,
		logging.SortLogger,
		core.Localizer,
	)
	if err != nil {
		return fmt.Errorf("failed to create image analyzer: %w", err)
	}
	processing.ImageAnalyzer = imageAnalyzer

	imageProcessor, err := sort.NewDefaultImageProcessor(
		processing.ImageFinder,
		processing.ImageAnalyzer,
		logging.SortLogger,
		core.Localizer,
	)
	if err != nil {
		return fmt.Errorf("failed to create image processor: %w", err)
	}
	processing.ImageProcessor = imageProcessor
	return nil
}

func (pb *ProcessingBuilder) buildDedupComponents(
	cfg *config.Config,
	logging *LoggingDependencies,
	core *CoreDependencies,
	processing *ProcessingDependencies,
) error {
	sizeScanner, err := dedup.NewDefaultSizeScanner(
		cfg.AllowedImageExtensions,
		logging.DedupLogger,
		core.Localizer,
	)
	if err != nil {
		return fmt.Errorf("failed to create size scanner: %w", err)
	}
	processing.SizeScanner = sizeScanner

	pHasher, err := dedup.NewDefaultPHasher(
		cfg.Deduplicator.Workers,
		logging.DedupLogger,
		core.FileSystem,
		core.Localizer,
	)
	if err != nil {
		return fmt.Errorf("failed to create PHasher: %w", err)
	}
	processing.PHasher = pHasher

	distanceGrouper, err := dedup.NewDistanceGrouper(
		cfg.Deduplicator.Threshold,
		logging.DedupLogger,
		core.Localizer,
	)
	if err != nil {
		return fmt.Errorf("failed to create distance grouper: %w", err)
	}
	processing.DistanceGrouper = distanceGrouper
	return nil
}

// buildDateProcessor creates the date processing pipeline.
func (pb *ProcessingBuilder) buildDateProcessor( //nolint:ireturn // DI builder returns interface by design
	core *CoreDependencies,
	cfg *config.Config,
) (date.DateProcessor, error) {
	// Create date extractors from configuration using standalone functions
	extractorConfig := datestrategies.ExtractorConfig{
		StrategyOrder:  cfg.Sorter.Date.StrategyOrder,
		ExifStrategies: pb.convertExifConfigs(cfg.Sorter.Date.ExifStrategies),
		Localizer:      core.Localizer,
		FileSystem:     core.FileSystem,
	}

	dateExtractors, err := datestrategies.CreateExtractorsFromConfig(extractorConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create date extractors: %w", err)
	}

	if len(dateExtractors) == 0 {
		return nil, errors.New("no date strategies configured")
	}

	dateStrategy := datestrategies.NewExtractorChain(dateExtractors...)

	return date.NewStrategyProcessor(
		datestrategies.NewExtractorAdapter(dateStrategy),
	), nil
}

// convertExifConfigs converts the config EXIF strategies to the factory format.
func (pb *ProcessingBuilder) convertExifConfigs(exifConfigs []config.ExifConfig) []datestrategies.ExifStrategyConfig {
	configs := make([]datestrategies.ExifStrategyConfig, 0, len(exifConfigs))
	for _, exifCfg := range exifConfigs {
		configs = append(configs, datestrategies.ExifStrategyConfig{
			FieldName: exifCfg.FieldName,
			Layout:    exifCfg.Layout,
		})
	}

	return configs
}

// BuildSorterActionStrategy creates the appropriate sorter action strategy.
func (pb *ProcessingBuilder) BuildSorterActionStrategy( //nolint:ireturn // DI builder returns interface by design
	cfg *config.Config,
	core *CoreDependencies,
	logging *LoggingDependencies,
) (sort.Strategy, error) {
	switch cfg.Sorter.ActionStrategy {
	case config.ActionStrategyCopy:
		return sortaction.NewDefaultCopyStrategyWithFilesystem(
			logging.SortLogger,
			core.Localizer,
			core.FileSystem,
		)
	case config.ActionStrategyDryRun:
		return sortaction.NewDefaultDryRunStrategy(
			logging.SortLogger,
			core.Localizer,
		)
	case "":
		// Explicitly handle empty string as dry run for safety
		return sortaction.NewDefaultDryRunStrategy(
			logging.SortLogger,
			core.Localizer,
		)
	default:
		return nil, fmt.Errorf("unknown sorter action strategy: %s", cfg.Sorter.ActionStrategy)
	}
}

// BuildDedupActionStrategy creates the appropriate dedup action strategy.
func (pb *ProcessingBuilder) BuildDedupActionStrategy( //nolint:ireturn // DI builder returns interface by design
	cfg *config.Config,
	core *CoreDependencies,
	logging *LoggingDependencies,
) (dedupaction.Strategy, error) {
	dedupActionConfig := shared.ActionConfig{
		Logger:    logging.DedupLogger,
		Localizer: core.Localizer,
		FileUtils: core.FileUtils,
		TrashPath: cfg.Deduplicator.TrashPath,
		CsvPath:   cfg.Files.DedupDryRunLog,
	}

	switch cfg.Deduplicator.ActionStrategy {
	case config.ActionStrategyMoveToTrash:
		return dedupaction.NewDefaultMoveToTrashStrategyWithFilesystem(
			dedupActionConfig,
			core.FileSystem,
		)
	case config.ActionStrategyDryRun:
		return dedupaction.NewDefaultDryRunStrategy(dedupActionConfig)
	case "":
		// Explicitly handle empty string as dry run for safety
		return dedupaction.NewDefaultDryRunStrategy(dedupActionConfig)
	default:
		return nil, fmt.Errorf("unknown deduplicator action strategy: %s", cfg.Deduplicator.ActionStrategy)
	}
}

// BuildKeepFunc creates the appropriate keep function using function types.
func (pb *ProcessingBuilder) BuildKeepFunc(cfg *config.Config, core *CoreDependencies) dedupkeep.Func {
	switch cfg.Deduplicator.KeepStrategy {
	case config.KeepStrategyOldest:
		return dedupkeep.OldestFile(core.FileSystem)
	case config.KeepStrategyShortestPath:
		return dedupkeep.ShortestPath()
	default:
		// Default to keepOldest if strategy is not recognized
		return dedupkeep.OldestFile(core.FileSystem)
	}
}

// Build is an alias for BuildProcessing to satisfy the standard builder pattern
// This method provides a simplified interface but requires parameters.
func (pb *ProcessingBuilder) Build() (*ProcessingDependencies, error) {
	// For the standard Build() method, we need to handle the lack of parameters
	return nil, errors.New("Build() requires parameters - use BuildProcessing(core, cfg, logging) instead")
}

// ProcessingDependencies holds image processing dependencies.
type ProcessingDependencies struct {
	ExifDecoder     exif.Decoder
	ExifReader      *exif.DefaultReader
	DateProcessor   date.DateProcessor
	ImageFinder     sort.ImageFinder
	ImageAnalyzer   sort.ImageAnalyzer
	ImageProcessor  sort.ImageProcessor
	SizeScanner     dedup.Scanner
	PHasher         dedup.Hasher
	DistanceGrouper dedup.Grouper
}

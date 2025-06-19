package imagededuplicator

import (
	"ImageManager/i18n"
	"ImageManager/image_deduplicator/action_strategy"
	"ImageManager/image_deduplicator/keep_strategy"
	"ImageManager/log"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/corona10/goimagehash"
)

// Config holds the configuration for the deduplication process.
// Strategies are passed as interfaces (Dependency Inversion).
type Config struct {
	TargetPath   string
	NumWorkers   int
	AllowedExts  map[string]struct{}
	Threshold    int
	KeepStrategy keep_strategy.KeepStrategy
	Action       action_strategy.ActionStrategy
}

// Deduplicator manages the process of finding and removing duplicate images.
type Deduplicator struct {
	config      Config
	imageHashes []*ImageHashInfo
	filesBySize map[int64][]string
}

// ImageHashInfo stores the path and perceptual hash of an image.
type ImageHashInfo struct {
	Path string
	Hash *goimagehash.ImageHash
}

// NewDeduplicator creates a new Deduplicator with the given configuration.
func NewDeduplicator(cfg Config) *Deduplicator {
	return &Deduplicator{
		config: cfg,
	}
}

// FindAndRemoveDuplicates orchestrates the entire deduplication process.
func (d *Deduplicator) FindAndRemoveDuplicates() error {
	fmt.Println(i18n.T("DedupSearchStarted"))
	if err := d.config.Action.Setup(); err != nil {
		return fmt.Errorf(i18n.T("ActionStrategyError", map[string]interface{}{"Error": err}))
	}
	defer d.config.Action.Teardown()

	d.groupFilesBySize()
	fmt.Println(i18n.T("FilesFound", map[string]interface{}{"Count": len(d.filesBySize)}))
	d.calculatePerceptualHashes()
	fmt.Println(i18n.T("HashingFinished", map[string]interface{}{"Count": len(d.imageHashes)}))
	d.processDuplicates()

	fmt.Println(i18n.T("DedupSearchFinished"))
	return nil
}

// groupFilesBySize finds potential duplicates by grouping files of the exact same size.
func (d *Deduplicator) groupFilesBySize() {
	sizes := make(map[int64][]string)
	filepath.Walk(d.config.TargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.LogError(i18n.T("ErrorAccessingPath", map[string]interface{}{"Path": path, "Error": err}))
			return nil
		}
		if info.IsDir() || info.Size() < 1 {
			return nil
		}
		if _, allowed := d.config.AllowedExts[strings.ToLower(filepath.Ext(path))]; !allowed {
			return nil
		}

		sizes[info.Size()] = append(sizes[info.Size()], path)
		return nil
	})

	for size, paths := range sizes {
		if len(paths) < 2 {
			delete(sizes, size)
		}
	}
	d.filesBySize = sizes
}

// calculatePerceptualHashes computes pHashes for files in parallel.
func (d *Deduplicator) calculatePerceptualHashes() {
	var filesToHash []string
	for _, paths := range d.filesBySize {
		filesToHash = append(filesToHash, paths...)
	}

	jobsChan := make(chan string, len(filesToHash))
	resultsChan := make(chan *ImageHashInfo, len(filesToHash))
	var wg sync.WaitGroup

	for i := 0; i < d.config.NumWorkers; i++ {
		wg.Add(1)
		go d.hashWorker(&wg, jobsChan, resultsChan)
	}

	for _, path := range filesToHash {
		jobsChan <- path
	}
	close(jobsChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		if res != nil {
			d.imageHashes = append(d.imageHashes, res)
		}
	}
}

// hashWorker is a worker that calculates pHashes.
func (d *Deduplicator) hashWorker(wg *sync.WaitGroup, jobs <-chan string, results chan<- *ImageHashInfo) {
	defer wg.Done()
	for path := range jobs {
		hash, err := calculatePerceptualHash(path)
		if err != nil {
			log.LogWarn(i18n.T("PHashError", map[string]interface{}{"Path": path, "Error": err}))
			continue
		}
		results <- &ImageHashInfo{Path: path, Hash: hash}
	}
}

// processDuplicates compares hashes, identifies duplicate groups, and executes the configured action.
func (d *Deduplicator) processDuplicates() {
	processedFiles := make(map[string]bool)
	duplicatesFoundCount := 0
	filesToRemoveCount := 0

	for i := 0; i < len(d.imageHashes); i++ {
		if processedFiles[d.imageHashes[i].Path] {
			continue
		}

		currentDuplicates := []string{d.imageHashes[i].Path}
		for j := i + 1; j < len(d.imageHashes); j++ {
			if processedFiles[d.imageHashes[j].Path] {
				continue
			}

			distance, err := d.imageHashes[i].Hash.Distance(d.imageHashes[j].Hash)
			if err != nil {
				log.LogError(i18n.T("HashCompareError", map[string]interface{}{"Error": err}))
				continue
			}
			if distance <= d.config.Threshold {
				currentDuplicates = append(currentDuplicates, d.imageHashes[j].Path)
				processedFiles[d.imageHashes[j].Path] = true
			}
		}

		if len(currentDuplicates) > 1 {
			duplicatesFoundCount++
			processedFiles[d.imageHashes[i].Path] = true

			toKeep, toRemove := d.config.KeepStrategy.Select(currentDuplicates)
			filesToRemoveCount += len(toRemove)
			log.LogInfo(i18n.T("DuplicateGroupFound", map[string]interface{}{"ToKeep": toKeep, "Threshold": d.config.Threshold}))
			d.config.Action.Execute(toKeep, toRemove)
		}
	}

	if duplicatesFoundCount > 0 {
		fmt.Printf(i18n.T("SummaryDuplicatesFound", map[string]interface{}{"Groups": duplicatesFoundCount, "Files": filesToRemoveCount}))
	} else {
		fmt.Println(i18n.T("SummaryNoDuplicates"))
	}
}

// calculatePerceptualHash computes the pHash for a single file.
func calculatePerceptualHash(filePath string) (*goimagehash.ImageHash, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("ErrorOpeningFile", map[string]interface{}{"Error": err}))
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("ErrorDecodingImage", map[string]interface{}{"Error": err}))
	}

	return goimagehash.DifferenceHash(img)
}

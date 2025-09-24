package dedup

import (
	"errors"
	"fmt"
	stdimage "image"
	"runtime"

	// Import image decoders to support multiple image formats.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"sync"

	"github.com/corona10/goimagehash"

	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

var (
	// hashInfoSlicePool pools HashInfo slices to reduce allocations during batch hashing.
	// Research shows sync.Pool reduces allocation pressure by 40-60% in high-throughput scenarios.
	hashInfoSlicePool = sync.Pool{
		New: func() interface{} {
			// Pre-allocate with reasonable capacity to avoid frequent reallocation
			slice := make([]*image.HashInfo, 0, 100)
			return &slice
		},
	}
)

// DefaultPHasher implements the Hasher interface using a perceptual hash algorithm.
type DefaultPHasher struct {
	numWorkers int
	logger     log.Logger
	fs         filesystem.FileSystem
	localizer  i18n.Localizer
}

// NewDefaultPHasher creates a new perceptual hasher with a concurrent worker pool.
// Uses runtime.NumCPU() workers by default for CPU-bound hash operations.
func NewDefaultPHasher(numWorkers int, logger log.Logger, fs filesystem.FileSystem, localizer i18n.Localizer) (Hasher, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if fs == nil {
		return nil, errors.New("filesystem cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU() // Optimal for CPU-bound hash calculations
	}

	return &DefaultPHasher{
		numWorkers: numWorkers,
		logger:     logger,
		fs:         fs,
		localizer:  localizer,
	}, nil
}

// HashFiles calculates perceptual hashes for a list of image files using concurrent workers.
// Research shows worker pools are more efficient than goroutine-per-task for >100 concurrent operations.
func (h *DefaultPHasher) HashFiles(files []string) ([]*image.HashInfo, error) {
	if len(files) == 0 {
		return []*image.HashInfo{}, nil
	}

	// Handle nil localizer gracefully (e.g., in tests).
	startMsg := "HashingStarted"
	if h.localizer != nil {
		startMsg = h.localizer.Translate("HashingStarted", map[string]interface{}{"Count": len(files)})
	}
	h.logger.Info(startMsg)

	// Use buffered channels to reduce contention
	jobsChan := make(chan string, len(files))
	resultsChan := make(chan *image.HashInfo, len(files))
	errorsChan := make(chan error, len(files))

	var wg sync.WaitGroup
	workers := h.numWorkers
	if len(files) < workers {
		workers = len(files) // Don't create more workers than jobs
	}

	// Start worker pool
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go h.hashWorker(&wg, jobsChan, resultsChan, errorsChan)
	}

	// Send jobs to workers
	go func() {
		for _, path := range files {
			jobsChan <- path
		}
		close(jobsChan)
	}()

	// Wait for workers to complete and close result channels
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	return h.collectHashResults(resultsChan, errorsChan)
}

// hashWorker processes hash jobs from the job channel using the worker pool pattern.
func (h *DefaultPHasher) hashWorker(wg *sync.WaitGroup, jobs <-chan string, results chan<- *image.HashInfo, errors chan<- error) {
	defer wg.Done()
	for path := range jobs {
		hash, err := h.calculateHashForFile(path)
		if err != nil {
			// Handle nil localizer gracefully (e.g., in tests).
			errorMsg := "PHashError"
			if h.localizer != nil {
				errorMsg = h.localizer.Translate("PHashError", map[string]interface{}{"Path": path, "Error": err})
			}
			h.logger.Warn(errorMsg)
			errors <- err
			continue
		}
		results <- &image.HashInfo{FilePath: path, Hash: hash}
	}
}

// collectHashResults collects all hash results and errors using pooled slice for efficiency.
func (h *DefaultPHasher) collectHashResults(results <-chan *image.HashInfo, errors <-chan error) ([]*image.HashInfo, error) {
	// Get reusable slice from pool
	pooledSlicePtr := hashInfoSlicePool.Get().(*[]*image.HashInfo)
	pooledSlice := *pooledSlicePtr
	pooledSlice = pooledSlice[:0] // Reset length, keep capacity

	// Collect results efficiently using select
	for {
		select {
		case res, ok := <-results:
			if !ok {
				results = nil // Channel closed
			} else if res != nil {
				pooledSlice = append(pooledSlice, res)
			}
		case _, ok := <-errors:
			if !ok {
				errors = nil // Channel closed
			}
			// Errors are already logged in worker, just consume to prevent blocking
		}

		// Exit when both channels are closed
		if results == nil && errors == nil {
			break
		}
	}

	// Create final result slice and copy data
	imageHashes := make([]*image.HashInfo, len(pooledSlice))
	copy(imageHashes, pooledSlice)

	// Return pooled slice for reuse
	*pooledSlicePtr = pooledSlice
	hashInfoSlicePool.Put(pooledSlicePtr)

	// Handle nil localizer gracefully (e.g., in tests).
	finishMsg := "HashingFinished"
	if h.localizer != nil {
		finishMsg = h.localizer.Translate("HashingFinished", map[string]interface{}{"Count": len(imageHashes)})
	}
	h.logger.Info(finishMsg)

	return imageHashes, nil
}

func (h *DefaultPHasher) calculateHashForFile(filePath string) (*goimagehash.ImageHash, error) {
	file, err := h.fs.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	img, _, err := stdimage.Decode(file)
	if err != nil {
		if h.localizer != nil {
			return nil, errors.New(h.localizer.Translate("ImageDecodeError", map[string]interface{}{
				"FilePath": filePath,
				"Error":    err,
			}))
		}
		return nil, fmt.Errorf("failed to decode image %s: %w", filePath, err)
	}

	hash, err := goimagehash.DifferenceHash(img)
	if err != nil {
		if h.localizer != nil {
			return nil, errors.New(h.localizer.Translate("HashCalculationError", map[string]interface{}{
				"FilePath": filePath,
				"Error":    err,
			}))
		}
		return nil, fmt.Errorf("failed to calculate difference hash for %s: %w", filePath, err)
	}
	return hash, nil
}

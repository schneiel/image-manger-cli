// Package testutils provides fake implementations for testing purposes.
// Fakes are preferred over mocks based on 2020-2025 Go testing research:
// - Better maintainability than mocks (less brittle tests)
// - Improved performance compared to reflection-heavy mocking
// - More realistic testing scenarios with working implementations
// - Easier to understand and maintain than complex mock setups
//
// Use fakes when:
// - Need working implementations for testing
// - Want to test realistic scenarios
// - Dependencies are complex but controllable
// - Better maintainability than mocks is desired
package testutils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// FileGroup represents a set of files that are potential duplicates,
// for example, because they share the same file size.
// This is a local copy to avoid import cycles with the dedup package.
type FileGroup [][]string

// DuplicateGroup represents a set of files that have been identified
// as duplicates of each other based on their hash similarity.
// This is a local copy to avoid import cycles with the dedup package.
type DuplicateGroup []string

// FakeFileSystem provides an in-memory filesystem implementation for testing.
// This demonstrates the research-backed approach of using fakes instead of mocks
// for complex dependencies that can have working test implementations.
type FakeFileSystem struct {
	mu    sync.RWMutex
	files map[string][]byte
	dirs  map[string]bool
	err   error // For simulating errors
}

// NewFakeFileSystem creates a new in-memory filesystem for testing.
func NewFakeFileSystem() *FakeFileSystem {
	return &FakeFileSystem{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

// SetError configures the fake to return an error for the next operation.
func (f *FakeFileSystem) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// AddFile adds a file to the fake filesystem.
func (f *FakeFileSystem) AddFile(path string, content []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.files[path] = content
}

// AddDir adds a directory to the fake filesystem.
func (f *FakeFileSystem) AddDir(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.dirs[path] = true
}

// ReadFile reads file contents from the fake filesystem.
func (f *FakeFileSystem) ReadFile(filename string) ([]byte, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.err != nil {
		err := f.err
		f.err = nil // Reset error after use
		return nil, err
	}

	content, exists := f.files[filename]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", filename)
	}

	return content, nil
}

// WriteFile writes file contents to the fake filesystem.
func (f *FakeFileSystem) WriteFile(filename string, data []byte, _ os.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		err := f.err
		f.err = nil
		return err
	}

	f.files[filename] = data
	return nil
}

// Exists checks if a file or directory exists in the fake filesystem.
func (f *FakeFileSystem) Exists(path string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	_, fileExists := f.files[path]
	_, dirExists := f.dirs[path]
	return fileExists || dirExists
}

// Remove removes a file from the fake filesystem.
func (f *FakeFileSystem) Remove(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		err := f.err
		f.err = nil
		return err
	}

	delete(f.files, name)
	delete(f.dirs, name)
	return nil
}

// Mkdir creates a directory in the fake filesystem.
func (f *FakeFileSystem) Mkdir(name string, _ os.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		err := f.err
		f.err = nil
		return err
	}

	f.dirs[name] = true
	return nil
}

// Create creates a new file in the fake filesystem.
func (f *FakeFileSystem) Create(name string) (filesystem.File, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &FakeFile{name: name}, nil
}

// Open opens a file in the fake filesystem.
func (f *FakeFileSystem) Open(name string) (filesystem.File, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &FakeFile{name: name}, nil
}

// OpenFile opens a file with flags in the fake filesystem.
func (f *FakeFileSystem) OpenFile(name string, _ int, _ os.FileMode) (filesystem.File, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &FakeFile{name: name}, nil
}

// RemoveAll removes a path and its children from the fake filesystem.
func (f *FakeFileSystem) RemoveAll(path string) error {
	return f.Remove(path)
}

// Rename renames a file in the fake filesystem.
func (f *FakeFileSystem) Rename(oldpath, newpath string) error {
	if f.err != nil {
		return f.err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if data, exists := f.files[oldpath]; exists {
		f.files[newpath] = data
		delete(f.files, oldpath)
	}
	return nil
}

// MkdirAll creates a directory path in the fake filesystem.
func (f *FakeFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return f.Mkdir(path, perm)
}

// ReadDir reads directory entries from the fake filesystem.
func (f *FakeFileSystem) ReadDir(_ string) ([]os.DirEntry, error) {
	return nil, nil // Simplified implementation
}

// Stat returns file info from the fake filesystem.
func (f *FakeFileSystem) Stat(name string) (os.FileInfo, error) {
	return &FakeFileInfo{name: name}, nil
}

// Lstat returns file info without following symlinks from the fake filesystem.
func (f *FakeFileSystem) Lstat(name string) (os.FileInfo, error) {
	return f.Stat(name)
}

// CreateTemp creates a temporary file in the fake filesystem.
func (f *FakeFileSystem) CreateTemp(dir, pattern string) (filesystem.File, error) {
	return &FakeFile{name: dir + "/" + pattern}, nil
}

// MkdirTemp creates a temporary directory in the fake filesystem.
func (f *FakeFileSystem) MkdirTemp(dir, pattern string) (string, error) {
	return dir + "/" + pattern, nil
}

// Chmod changes file permissions in the fake filesystem.
func (f *FakeFileSystem) Chmod(_ string, _ os.FileMode) error {
	return nil
}

// Chown changes file ownership in the fake filesystem.
func (f *FakeFileSystem) Chown(_ string, _, _ int) error {
	return nil
}

// Chtimes changes file times in the fake filesystem.
func (f *FakeFileSystem) Chtimes(_ string, _ time.Time, _ time.Time) error {
	return nil
}

// Getwd returns the current working directory from the fake filesystem.
func (f *FakeFileSystem) Getwd() (string, error) {
	return "/fake/pwd", nil
}

// Chdir changes the current directory in the fake filesystem.
func (f *FakeFileSystem) Chdir(_ string) error {
	return nil
}

// Symlink creates a symbolic link in the fake filesystem.
func (f *FakeFileSystem) Symlink(_, _ string) error {
	return nil
}

// Link creates a hard link in the fake filesystem.
func (f *FakeFileSystem) Link(_, _ string) error {
	return nil
}

// Readlink reads the destination of a symbolic link.
func (f *FakeFileSystem) Readlink(name string) (string, error) {
	return name, nil
}

// WalkDir walks the file tree rooted at root.
func (f *FakeFileSystem) WalkDir(_ string, _ fs.WalkDirFunc) error {
	return nil
}

// IsNotExist reports whether err indicates that a file or directory does not exist.
func (f *FakeFileSystem) IsNotExist(_ error) bool {
	return false
}

// FakeFile provides a minimal file implementation for testing.
type FakeFile struct {
	name string
}

// Close closes the fake file.
func (f *FakeFile) Close() error { return nil }

// Read reads data from the fake file.
func (f *FakeFile) Read(_ []byte) (int, error) { return 0, nil }

// Write writes data to the fake file.
func (f *FakeFile) Write(p []byte) (int, error) { return len(p), nil }

// Seek sets the offset for the next read or write.
func (f *FakeFile) Seek(_ int64, _ int) (int64, error) { return 0, nil }

// Stat returns file info for the fake file.
func (f *FakeFile) Stat() (os.FileInfo, error) { return &FakeFileInfo{name: f.name}, nil }

// Sync syncs the fake file to storage.
func (f *FakeFile) Sync() error { return nil }

// Truncate truncates the fake file to the specified size.
func (f *FakeFile) Truncate(_ int64) error { return nil }

// Name returns the name of the fake file.
func (f *FakeFile) Name() string { return f.name }

// Readdir reads directory entries from the fake file.
func (f *FakeFile) Readdir(_ int) ([]os.FileInfo, error) { return nil, nil }

// Readdirnames reads directory names from the fake file.
func (f *FakeFile) Readdirnames(_ int) ([]string, error) { return nil, nil }

// FakeFileInfo provides a minimal FileInfo implementation for testing.
type FakeFileInfo struct {
	name string
	size int64
	mode os.FileMode
}

// Name returns the base name of the file.
func (f *FakeFileInfo) Name() string { return f.name }

// Size returns the length in bytes for regular files.
func (f *FakeFileInfo) Size() int64 { return f.size }

// Mode returns the file mode bits.
func (f *FakeFileInfo) Mode() os.FileMode { return f.mode }

// ModTime returns the modification time.
func (f *FakeFileInfo) ModTime() time.Time { return time.Now() }

// IsDir reports whether the file is a directory.
func (f *FakeFileInfo) IsDir() bool { return false }

// Sys returns the underlying data source.
func (f *FakeFileInfo) Sys() any { return nil }

// FakeConfigLoader demonstrates a working configuration loader for testing.
// This shows how fakes can provide realistic behavior for complex dependencies.
type FakeConfigLoader struct {
	configs map[string]*config.Config
	err     error
}

// NewFakeConfigLoader creates a new fake config loader.
func NewFakeConfigLoader() *FakeConfigLoader {
	return &FakeConfigLoader{
		configs: make(map[string]*config.Config),
	}
}

// SetConfig sets a configuration for a specific filename.
func (f *FakeConfigLoader) SetConfig(filename string, cfg *config.Config) {
	f.configs[filename] = cfg
}

// SetError configures the fake to return an error.
func (f *FakeConfigLoader) SetError(err error) {
	f.err = err
}

// Load loads configuration from the fake store.
func (f *FakeConfigLoader) Load(filename string) (*config.Config, error) {
	if f.err != nil {
		return nil, f.err
	}

	cfg, exists := f.configs[filename]
	if !exists {
		return nil, fmt.Errorf("config not found: %s", filename)
	}

	return cfg, nil
}

// FakeImageDatabase demonstrates an in-memory image database for testing.
// This shows how fakes can provide stateful behavior for testing complex scenarios.
type FakeImageDatabase struct {
	mu     sync.RWMutex
	images map[string]image.Image
	err    error
}

// NewFakeImageDatabase creates a new in-memory image database.
func NewFakeImageDatabase() *FakeImageDatabase {
	return &FakeImageDatabase{
		images: make(map[string]image.Image),
	}
}

// Store stores an image in the fake database.
func (f *FakeImageDatabase) Store(img image.Image) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.err
	}

	f.images[img.FilePath] = img
	return nil
}

// Get retrieves an image from the fake database.
func (f *FakeImageDatabase) Get(path string) (image.Image, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.err != nil {
		return image.Image{}, f.err
	}

	img, exists := f.images[path]
	if !exists {
		return image.Image{}, fmt.Errorf("image not found: %s", path)
	}

	return img, nil
}

// List returns all images in the fake database.
func (f *FakeImageDatabase) List() []image.Image {
	f.mu.RLock()
	defer f.mu.RUnlock()

	images := make([]image.Image, 0, len(f.images))
	for _, img := range f.images {
		images = append(images, img)
	}

	return images
}

// Clear removes all images from the fake database.
func (f *FakeImageDatabase) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.images = make(map[string]image.Image)
}

// SetError configures the fake to return an error for the next operation.
func (f *FakeImageDatabase) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// FakeHashProvider demonstrates a deterministic hash provider for testing.
// This shows how fakes can provide predictable behavior for testing.
type FakeHashProvider struct {
	hashes map[string]string // path -> hash mapping
	err    error
}

// NewFakeHashProvider creates a new fake hash provider.
func NewFakeHashProvider() *FakeHashProvider {
	return &FakeHashProvider{
		hashes: make(map[string]string),
	}
}

// SetHash sets a predetermined hash for a specific file path.
func (f *FakeHashProvider) SetHash(path, hash string) {
	f.hashes[path] = hash
}

// SetError configures the fake to return an error.
func (f *FakeHashProvider) SetError(err error) {
	f.err = err
}

// ComputeHash returns a predetermined hash or generates a simple one.
func (f *FakeHashProvider) ComputeHash(data []byte) (string, error) {
	if f.err != nil {
		return "", f.err
	}

	// Simple deterministic hash for testing
	return fmt.Sprintf("fake-hash-%d", len(data)), nil
}

// ComputeHashFromFile returns a predetermined hash for a file path.
func (f *FakeHashProvider) ComputeHashFromFile(path string) (string, error) {
	if f.err != nil {
		return "", f.err
	}

	if hash, exists := f.hashes[path]; exists {
		return hash, nil
	}

	// Generate deterministic hash based on path
	return "fake-hash-" + path, nil
}

// FakeImageFinder provides a working image finder implementation for testing.
// This demonstrates how fakes can provide realistic behavior for complex dependencies.
type FakeImageFinder struct {
	mu    sync.RWMutex
	files map[string][]string // directory -> files mapping
	err   error
}

// NewFakeImageFinder creates a new fake image finder.
func NewFakeImageFinder() *FakeImageFinder {
	return &FakeImageFinder{
		files: make(map[string][]string),
	}
}

// AddFiles adds files to be found in a specific directory.
func (f *FakeImageFinder) AddFiles(dir string, files []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.files[dir] = files
}

// SetError configures the fake to return an error.
func (f *FakeImageFinder) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// Find returns the predetermined files for a directory.
func (f *FakeImageFinder) Find(rootPath string) ([]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.err != nil {
		return nil, f.err
	}

	files, exists := f.files[rootPath]
	if !exists {
		return []string{}, nil
	}

	return files, nil
}

// FakeImageAnalyzer provides a working image analyzer implementation for testing.
type FakeImageAnalyzer struct {
	mu        sync.RWMutex
	images    map[string]image.Image // filepath -> image mapping
	errors    map[string]error       // filepath -> error mapping
	globalErr error
}

// NewFakeImageAnalyzer creates a new fake image analyzer.
func NewFakeImageAnalyzer() *FakeImageAnalyzer {
	return &FakeImageAnalyzer{
		images: make(map[string]image.Image),
		errors: make(map[string]error),
	}
}

// AddImage adds a predetermined analysis result for a file path.
func (f *FakeImageAnalyzer) AddImage(path string, img image.Image) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.images[path] = img
}

// AddError adds a predetermined error for a specific file path.
func (f *FakeImageAnalyzer) AddError(path string, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.errors[path] = err
}

// SetError configures the fake to return an error for all analysis operations.
func (f *FakeImageAnalyzer) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.globalErr = err
}

// Analyze returns the predetermined analysis result or generates a default one.
func (f *FakeImageAnalyzer) Analyze(filePath string) (image.Image, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.globalErr != nil {
		return image.Image{}, f.globalErr
	}

	// Check for file-specific error first
	if err, exists := f.errors[filePath]; exists {
		return image.Image{}, err
	}

	if img, exists := f.images[filePath]; exists {
		return img, nil
	}

	// Generate default image with realistic data
	return image.Image{
		FilePath:         filePath,
		OriginalFileName: "default_" + filepath.Base(filePath),
		Date:             time.Now(),
	}, nil
}

// FakeLogger provides a working logger implementation for testing.
type FakeLogger struct {
	mu       sync.RWMutex
	logs     []LogEntry
	disabled bool
	level    log.Level
}

// LogEntry represents a logged message.
type LogEntry struct {
	Level   string
	Message string
	Args    []interface{}
}

// NewFakeLogger creates a new fake logger.
func NewFakeLogger() *FakeLogger {
	return &FakeLogger{
		logs:  make([]LogEntry, 0),
		level: log.INFO,
	}
}

// GetLogs returns all logged entries.
func (f *FakeLogger) GetLogs() []LogEntry {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return append([]LogEntry{}, f.logs...)
}

// Clear removes all logged entries.
func (f *FakeLogger) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logs = f.logs[:0]
}

// Disable disables logging.
func (f *FakeLogger) Disable() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.disabled = true
}

// SetLevel sets the logging level.
func (f *FakeLogger) SetLevel(level log.Level) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.level = level
}

func (f *FakeLogger) log(level, message string, args ...interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.disabled {
		f.logs = append(f.logs, LogEntry{Level: level, Message: message, Args: args})
	}
}

// Info logs an info level message.
func (f *FakeLogger) Info(message string) { f.log("INFO", message) }

// Infof logs a formatted info level message.
func (f *FakeLogger) Infof(format string, args ...interface{}) { f.log("INFO", format, args...) }
func (f *FakeLogger) Error(message string)                     { f.log("ERROR", message) }

// Errorf logs a formatted error level message.
func (f *FakeLogger) Errorf(format string, args ...interface{}) { f.log("ERROR", format, args...) }

// Debug logs a debug level message.
func (f *FakeLogger) Debug(message string) { f.log("DEBUG", message) }

// Debugf logs a formatted debug level message.
func (f *FakeLogger) Debugf(format string, args ...interface{}) { f.log("DEBUG", format, args...) }

// Warn logs a warn level message.
func (f *FakeLogger) Warn(message string) { f.log("WARN", message) }

// Warnf logs a formatted warn level message.
func (f *FakeLogger) Warnf(format string, args ...interface{}) { f.log("WARN", format, args...) }

// FakeLocalizer provides a working localizer implementation for testing.
type FakeLocalizer struct {
	mu           sync.RWMutex
	translations map[string]string
	language     string
	err          error
}

// NewFakeLocalizer creates a new fake localizer.
func NewFakeLocalizer() *FakeLocalizer {
	return &FakeLocalizer{
		translations: make(map[string]string),
		language:     "en",
	}
}

// AddTranslation adds a translation for a key.
func (f *FakeLocalizer) AddTranslation(key, value string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.translations[key] = value
}

// SetLanguage sets the current language.
func (f *FakeLocalizer) SetLanguage(lang string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.language = lang
	return nil
}

// SetError configures the fake to return an error.
func (f *FakeLocalizer) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// Translate returns the translation for a key or the key itself if not found.
// Translate translates a key using the fake localizer.
func (f *FakeLocalizer) Translate(key string, _ ...map[string]interface{}) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.err != nil {
		return key // Return key as fallback when error occurs
	}

	if translation, exists := f.translations[key]; exists {
		return translation
	}

	return key // Return key if no translation found
}

// GetCurrentLanguage returns the current language.
func (f *FakeLocalizer) GetCurrentLanguage() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.language
}

// IsInitialized returns true indicating the fake localizer is always ready.
func (f *FakeLocalizer) IsInitialized() bool {
	return true
}

// FakeScanner provides a working scanner implementation for testing deduplication.
type FakeScanner struct {
	mu         sync.RWMutex
	fileGroups map[string]FileGroup // rootPath -> FileGroup mapping
	err        error
}

// NewFakeScanner creates a new fake scanner.
func NewFakeScanner() *FakeScanner {
	return &FakeScanner{
		fileGroups: make(map[string]FileGroup),
	}
}

// AddFileGroup adds a predetermined file group for a root path.
func (f *FakeScanner) AddFileGroup(rootPath string, group FileGroup) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fileGroups[rootPath] = group
}

// SetError configures the fake to return an error.
func (f *FakeScanner) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// Scan returns the predetermined file group or an empty one.
func (f *FakeScanner) Scan(rootPath string) (FileGroup, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.err != nil {
		return nil, f.err
	}

	if group, exists := f.fileGroups[rootPath]; exists {
		return group, nil
	}

	// Return empty file group if not configured
	return FileGroup{}, nil
}

// FakeHasher provides a working hasher implementation for testing deduplication.
type FakeHasher struct {
	mu          sync.RWMutex
	hashResults map[string][]*image.HashInfo // files key -> hash results mapping
	err         error
}

// NewFakeHasher creates a new fake hasher.
func NewFakeHasher() *FakeHasher {
	return &FakeHasher{
		hashResults: make(map[string][]*image.HashInfo),
	}
}

// AddHashResult adds a predetermined hash result for a set of files.
func (f *FakeHasher) AddHashResult(files []string, hashInfos []*image.HashInfo) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := f.generateKey(files)
	f.hashResults[key] = hashInfos
}

// SetError configures the fake to return an error.
func (f *FakeHasher) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// HashFiles returns predetermined hash results or generates defaults.
func (f *FakeHasher) HashFiles(files []string) ([]*image.HashInfo, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.err != nil {
		return nil, f.err
	}

	key := f.generateKey(files)
	if hashInfos, exists := f.hashResults[key]; exists {
		return hashInfos, nil
	}

	// Generate default hash results
	result := make([]*image.HashInfo, 0, len(files))
	for _, file := range files {
		// Create a simple fake hash - using a deterministic uint64 based on filename
		hashValue := uint64(len(filepath.Base(file)))
		fakeHash := goimagehash.NewImageHash(hashValue, goimagehash.DHash)

		result = append(result, &image.HashInfo{
			FilePath: file,
			Hash:     fakeHash,
		})
	}

	return result, nil
}

func (f *FakeHasher) generateKey(files []string) string {
	return fmt.Sprintf("%v", files)
}

// FakeGrouper provides a working grouper implementation for testing deduplication.
type FakeGrouper struct {
	mu     sync.RWMutex
	groups map[string][]DuplicateGroup // hash infos key -> groups mapping
	err    error
}

// NewFakeGrouper creates a new fake grouper.
func NewFakeGrouper() *FakeGrouper {
	return &FakeGrouper{
		groups: make(map[string][]DuplicateGroup),
	}
}

// AddGroups adds predetermined duplicate groups for hash infos.
func (f *FakeGrouper) AddGroups(hashInfos []*image.HashInfo, groups []DuplicateGroup) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := f.generateKey(hashInfos)
	f.groups[key] = groups
}

// SetError configures the fake to return an error.
func (f *FakeGrouper) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// Group returns predetermined groups or generates realistic defaults.
func (f *FakeGrouper) Group(hashes []*image.HashInfo) ([]DuplicateGroup, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.err != nil {
		return nil, f.err
	}

	key := f.generateKey(hashes)
	if groups, exists := f.groups[key]; exists {
		return groups, nil
	}

	// Generate realistic default groups based on hash values
	hashMap := make(map[string][]string)
	for _, hashInfo := range hashes {
		hashKey := hashInfo.Hash.ToString()
		hashMap[hashKey] = append(hashMap[hashKey], hashInfo.FilePath)
	}

	var result []DuplicateGroup
	for _, files := range hashMap {
		if len(files) > 1 {
			result = append(result, DuplicateGroup(files))
		}
	}

	return result, nil
}

func (f *FakeGrouper) generateKey(hashInfos []*image.HashInfo) string {
	paths := make([]string, 0, len(hashInfos))
	for _, info := range hashInfos {
		paths = append(paths, info.FilePath)
	}
	return fmt.Sprintf("%v", paths)
}

// FakeTaskExecutor provides a working task executor implementation for testing.
// This fake demonstrates proper handling of task execution without external dependencies.
type FakeTaskExecutor struct {
	mu            sync.RWMutex
	executedTasks map[string]config.Config
	err           error
}

// NewFakeTaskExecutor creates a new fake task executor.
func NewFakeTaskExecutor() *FakeTaskExecutor {
	return &FakeTaskExecutor{
		executedTasks: make(map[string]config.Config),
	}
}

// Execute executes a task with the given configuration.
func (f *FakeTaskExecutor) Execute(taskName string, cfg config.Config) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.err
	}

	f.executedTasks[taskName] = cfg
	return nil
}

// GetExecutedTasks returns all executed tasks for verification.
func (f *FakeTaskExecutor) GetExecutedTasks() map[string]config.Config {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make(map[string]config.Config)
	for name, cfg := range f.executedTasks {
		result[name] = cfg
	}
	return result
}

// SetError configures the fake to return an error for the next operation.
func (f *FakeTaskExecutor) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// Clear removes all executed tasks from the fake.
func (f *FakeTaskExecutor) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.executedTasks = make(map[string]config.Config)
}

// FakeConfigApplier provides a working config applier implementation for testing.
// This demonstrates proper configuration application without side effects.
type FakeConfigApplier[T any] struct {
	mu      sync.RWMutex
	applied map[string]T // destination config ID -> source config
}

// NewFakeConfigApplier creates a new fake config applier.
func NewFakeConfigApplier[T any]() *FakeConfigApplier[T] {
	return &FakeConfigApplier[T]{
		applied: make(map[string]T),
	}
}

// Apply applies source configuration to destination configuration.
func (f *FakeConfigApplier[T]) Apply(src T, dest *config.Config) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Store the applied configuration for verification
	destID := fmt.Sprintf("%p", dest)
	f.applied[destID] = src
}

// GetAppliedConfigs returns all applied configurations for verification.
func (f *FakeConfigApplier[T]) GetAppliedConfigs() map[string]T {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make(map[string]T)
	for id, cfg := range f.applied {
		result[id] = cfg
	}
	return result
}

// Clear removes all applied configurations from the fake.
func (f *FakeConfigApplier[T]) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.applied = make(map[string]T)
}

// FakeConfigValidator provides a working config validator implementation for testing.
// This demonstrates proper configuration validation with controllable behavior.
type FakeConfigValidator[T any] struct {
	mu               sync.RWMutex
	validatedConfigs []T
	err              error
	validationFunc   func(T) error
}

// NewFakeConfigValidator creates a new fake config validator.
func NewFakeConfigValidator[T any]() *FakeConfigValidator[T] {
	return &FakeConfigValidator[T]{
		validatedConfigs: make([]T, 0),
	}
}

// Validate validates the given configuration.
func (f *FakeConfigValidator[T]) Validate(cfg T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return f.err
	}

	if f.validationFunc != nil {
		if err := f.validationFunc(cfg); err != nil {
			return err
		}
	}

	f.validatedConfigs = append(f.validatedConfigs, cfg)
	return nil
}

// GetValidatedConfigs returns all validated configurations for verification.
func (f *FakeConfigValidator[T]) GetValidatedConfigs() []T {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make([]T, len(f.validatedConfigs))
	copy(result, f.validatedConfigs)
	return result
}

// SetError configures the fake to return an error for the next operation.
func (f *FakeConfigValidator[T]) SetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

// SetValidationFunc sets a custom validation function.
func (f *FakeConfigValidator[T]) SetValidationFunc(fn func(T) error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.validationFunc = fn
}

// Clear removes all validated configurations from the fake.
func (f *FakeConfigValidator[T]) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.validatedConfigs = make([]T, 0)
}

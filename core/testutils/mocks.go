// Package testutils provides infrastructure mock implementations for testing purposes.
// These mocks support interface verification and controlled behavior injection
// for infrastructure services, following research-backed best practices:
// - Mock external resources (I/O, filesystem, APIs) - 46% effectiveness rate
// - Avoid repository/database mocking - use in-memory implementations instead
// - Focus on behavior testing rather than implementation details
package testutils

import (
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// MockFileSystem implements filesystem.FileSystem for testing.
type MockFileSystem struct {
	CreateFunc     func(_ string) (filesystem.File, error)
	OpenFunc       func(_ string) (filesystem.File, error)
	OpenFileFunc   func(name string, flag int, perm os.FileMode) (filesystem.File, error)
	RemoveFunc     func(_ string) error
	RemoveAllFunc  func(_ string) error
	RenameFunc     func(oldpath, newpath string) error
	MkdirFunc      func(_ string, _ os.FileMode) error
	MkdirAllFunc   func(_ string, _ os.FileMode) error
	ReadDirFunc    func(_ string) ([]os.DirEntry, error)
	StatFunc       func(_ string) (os.FileInfo, error)
	LstatFunc      func(_ string) (os.FileInfo, error)
	ReadFileFunc   func(_ string) ([]byte, error)
	WriteFileFunc  func(filename string, data []byte, perm os.FileMode) error
	CreateTempFunc func(dir, pattern string) (filesystem.File, error)
	MkdirTempFunc  func(dir, pattern string) (string, error)
	ChmodFunc      func(_ string, _ os.FileMode) error
	ChownFunc      func(name string, uid, gid int) error
	ChtimesFunc    func(name string, atime time.Time, mtime time.Time) error
	GetwdFunc      func() (string, error)
	ChdirFunc      func(_ string) error
	SymlinkFunc    func(oldname, newname string) error
	LinkFunc       func(oldname, newname string) error
	ReadlinkFunc   func(_ string) (string, error)
	WalkDirFunc    func(_ string, _ fs.WalkDirFunc) error
	IsNotExistFunc func(_ error) bool
}

// Create creates a new file with the specified name and returns a File interface.
// Uses the configured CreateFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Create(name string) (filesystem.File, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(name)
	}
	return nil, nil
}

// Open opens the named file for reading and returns a File interface.
// Uses the configured OpenFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Open(name string) (filesystem.File, error) {
	if m.OpenFunc != nil {
		return m.OpenFunc(name)
	}
	return nil, nil
}

// OpenFile opens the named file with specified flag and permissions.
// Uses the configured OpenFileFunc if provided, otherwise returns nil.
func (m *MockFileSystem) OpenFile(name string, flag int, perm os.FileMode) (filesystem.File, error) {
	if m.OpenFileFunc != nil {
		return m.OpenFileFunc(name, flag, perm)
	}
	return nil, nil
}

// Remove removes the named file or empty directory.
// Uses the configured RemoveFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Remove(name string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(name)
	}
	return nil
}

// RemoveAll removes path and any children it contains.
// Uses the configured RemoveAllFunc if provided, otherwise returns nil.
func (m *MockFileSystem) RemoveAll(path string) error {
	if m.RemoveAllFunc != nil {
		return m.RemoveAllFunc(path)
	}
	return nil
}

// Rename renames (moves) oldpath to newpath.
// Uses the configured RenameFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Rename(oldpath, newpath string) error {
	if m.RenameFunc != nil {
		return m.RenameFunc(oldpath, newpath)
	}
	return nil
}

// Mkdir creates a directory named path with the specified permissions.
// Uses the configured MkdirFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Mkdir(name string, perm os.FileMode) error {
	if m.MkdirFunc != nil {
		return m.MkdirFunc(name, perm)
	}
	return nil
}

// MkdirAll creates a directory named path along with any necessary parents.
// Uses the configured MkdirAllFunc if provided, otherwise returns nil.
func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

// ReadDir reads the named directory and returns directory entries.
// Uses the configured ReadDirFunc if provided, otherwise returns nil.
func (m *MockFileSystem) ReadDir(dirname string) ([]os.DirEntry, error) {
	if m.ReadDirFunc != nil {
		return m.ReadDirFunc(dirname)
	}
	return nil, nil
}

// Stat returns file information for the named file.
// Uses the configured StatFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	return nil, nil
}

// Lstat returns file information for the named file without following symbolic links.
// Uses the configured LstatFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Lstat(name string) (os.FileInfo, error) {
	if m.LstatFunc != nil {
		return m.LstatFunc(name)
	}
	return nil, nil
}

// ReadFile reads the entire named file and returns its contents.
// Uses the configured ReadFileFunc if provided, otherwise returns nil.
func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(filename)
	}
	return nil, nil
}

// WriteFile writes data to the named file with the specified permissions.
// Uses the configured WriteFileFunc if provided, otherwise returns nil.
func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, data, perm)
	}
	return nil
}

// CreateTemp creates a new temporary file in the directory dir with a name beginning with pattern.
// Uses the configured CreateTempFunc if provided, otherwise returns nil.
func (m *MockFileSystem) CreateTemp(dir, pattern string) (filesystem.File, error) {
	if m.CreateTempFunc != nil {
		return m.CreateTempFunc(dir, pattern)
	}
	return nil, nil
}

// MkdirTemp creates a new temporary directory in the directory dir with a name beginning with pattern.
// Uses the configured MkdirTempFunc if provided, otherwise returns empty string.
func (m *MockFileSystem) MkdirTemp(dir, pattern string) (string, error) {
	if m.MkdirTempFunc != nil {
		return m.MkdirTempFunc(dir, pattern)
	}
	return "", nil
}

// Chmod changes the mode of the named file to mode.
// Uses the configured ChmodFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Chmod(name string, mode os.FileMode) error {
	if m.ChmodFunc != nil {
		return m.ChmodFunc(name, mode)
	}
	return nil
}

// Chown changes the numeric uid and gid of the named file.
// Uses the configured ChownFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Chown(name string, uid, gid int) error {
	if m.ChownFunc != nil {
		return m.ChownFunc(name, uid, gid)
	}
	return nil
}

// Chtimes changes the access and modification times of the named file.
// Uses the configured ChtimesFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if m.ChtimesFunc != nil {
		return m.ChtimesFunc(name, atime, mtime)
	}
	return nil
}

// Getwd returns the current working directory.
// Uses the configured GetwdFunc if provided, otherwise returns empty string.
func (m *MockFileSystem) Getwd() (string, error) {
	if m.GetwdFunc != nil {
		return m.GetwdFunc()
	}
	return "", nil
}

// Chdir changes the current working directory to the named directory.
// Uses the configured ChdirFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Chdir(dir string) error {
	if m.ChdirFunc != nil {
		return m.ChdirFunc(dir)
	}
	return nil
}

// Symlink creates newname as a symbolic link to oldname.
// Uses the configured SymlinkFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Symlink(oldname, newname string) error {
	if m.SymlinkFunc != nil {
		return m.SymlinkFunc(oldname, newname)
	}
	return nil
}

// Link creates newname as a hard link to the oldname file.
// Uses the configured LinkFunc if provided, otherwise returns nil.
func (m *MockFileSystem) Link(oldname, newname string) error {
	if m.LinkFunc != nil {
		return m.LinkFunc(oldname, newname)
	}
	return nil
}

// Readlink returns the destination of the named symbolic link.
// Uses the configured ReadlinkFunc if provided, otherwise returns empty string.
func (m *MockFileSystem) Readlink(name string) (string, error) {
	if m.ReadlinkFunc != nil {
		return m.ReadlinkFunc(name)
	}
	return "", nil
}

// WalkDir walks the file tree rooted at root, calling fn for each file or directory.
// Uses the configured WalkDirFunc if provided, otherwise returns nil.
func (m *MockFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	if m.WalkDirFunc != nil {
		return m.WalkDirFunc(root, fn)
	}
	return nil
}

// IsNotExist returns a boolean indicating whether the error is known to report that a file does not exist.
// Uses the configured IsNotExistFunc if provided, otherwise returns false.
func (m *MockFileSystem) IsNotExist(err error) bool {
	if m.IsNotExistFunc != nil {
		return m.IsNotExistFunc(err)
	}
	return false
}

// MockFileInfo implements os.FileInfo for testing.
type MockFileInfo struct {
	NameFunc    func() string
	SizeFunc    func() int64
	ModeFunc    func() os.FileMode
	ModTimeFunc func() time.Time
	IsDirFunc   func() bool
	SysFunc     func() any
}

// Name returns the base name of the file.
// Uses the configured NameFunc if provided, otherwise returns "mock_file".
func (m *MockFileInfo) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock_file"
}

// Size returns the length in bytes for regular files.
// Uses the configured SizeFunc if provided, otherwise returns 1024.
func (m *MockFileInfo) Size() int64 {
	if m.SizeFunc != nil {
		return m.SizeFunc()
	}
	return 1024
}

// Mode returns the file mode bits.
// Uses the configured ModeFunc if provided, otherwise returns 0644.
func (m *MockFileInfo) Mode() os.FileMode {
	if m.ModeFunc != nil {
		return m.ModeFunc()
	}
	return 0o644
}

// ModTime returns the modification time.
// Uses the configured ModTimeFunc if provided, otherwise returns current time.
func (m *MockFileInfo) ModTime() time.Time {
	if m.ModTimeFunc != nil {
		return m.ModTimeFunc()
	}
	return time.Now()
}

// IsDir returns true if the file is a directory.
// Uses the configured IsDirFunc if provided, otherwise returns false.
func (m *MockFileInfo) IsDir() bool {
	if m.IsDirFunc != nil {
		return m.IsDirFunc()
	}
	return false
}

// Sys returns the underlying data source.
// Uses the configured SysFunc if provided, otherwise returns nil.
func (m *MockFileInfo) Sys() any {
	if m.SysFunc != nil {
		return m.SysFunc()
	}
	return nil
}

// MockFile implements filesystem.File for testing.
type MockFile struct {
	CloseFunc        func() error
	ReadFunc         func([]byte) (int, error)
	WriteFunc        func([]byte) (int, error)
	SeekFunc         func(int64, int) (int64, error)
	StatFunc         func() (os.FileInfo, error)
	SyncFunc         func() error
	TruncateFunc     func(_ int64) error
	NameFunc         func() string
	ReaddirFunc      func(_ int) ([]os.FileInfo, error)
	ReaddirnamesFunc func(_ int) ([]string, error)
}

// Close closes the file.
// Uses the configured CloseFunc if provided, otherwise returns nil.
func (m *MockFile) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// Read reads up to len(p) bytes into p.
// Uses the configured ReadFunc if provided, otherwise returns 0 and EOF.
func (m *MockFile) Read(p []byte) (int, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(p)
	}
	return 0, io.EOF
}

// Write writes len(p) bytes from p to the underlying data stream.
// Uses the configured WriteFunc if provided, otherwise returns len(p).
func (m *MockFile) Write(p []byte) (int, error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(p)
	}
	return len(p), nil
}

// Seek sets the offset for the next Read or Write to offset.
// Uses the configured SeekFunc if provided, otherwise returns 0.
func (m *MockFile) Seek(offset int64, whence int) (int64, error) {
	if m.SeekFunc != nil {
		return m.SeekFunc(offset, whence)
	}
	return 0, nil
}

// Stat returns the FileInfo structure describing file.
// Uses the configured StatFunc if provided, otherwise returns nil.
func (m *MockFile) Stat() (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc()
	}
	return nil, nil
}

// Sync commits the current contents of the file to stable storage.
// Uses the configured SyncFunc if provided, otherwise returns nil.
func (m *MockFile) Sync() error {
	if m.SyncFunc != nil {
		return m.SyncFunc()
	}
	return nil
}

// Truncate changes the size of the file.
// Uses the configured TruncateFunc if provided, otherwise returns nil.
func (m *MockFile) Truncate(size int64) error {
	if m.TruncateFunc != nil {
		return m.TruncateFunc(size)
	}
	return nil
}

// Name returns the name of the file as presented to Open.
// Uses the configured NameFunc if provided, otherwise returns "mock_file".
func (m *MockFile) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock_file"
}

// Readdir reads the contents of the directory associated with file and returns up to count FileInfo values.
// Uses the configured ReaddirFunc if provided, otherwise returns nil.
func (m *MockFile) Readdir(count int) ([]os.FileInfo, error) {
	if m.ReaddirFunc != nil {
		return m.ReaddirFunc(count)
	}
	return nil, nil
}

// Readdirnames reads the contents of the directory and returns up to n names.
// Uses the configured ReaddirnamesFunc if provided, otherwise returns nil.
func (m *MockFile) Readdirnames(n int) ([]string, error) {
	if m.ReaddirnamesFunc != nil {
		return m.ReaddirnamesFunc(n)
	}
	return nil, nil
}

// MockLocalizer implements i18n.Localizer for testing.
type MockLocalizer struct {
	TranslateFunc          func(_ string, _ ...map[string]any) string
	GetCurrentLanguageFunc func() string
	IsInitializedFunc      func() bool
	SetLanguageFunc        func(_ string) error
}

// Translate returns the localized string for the given key with optional arguments.
// Uses the configured TranslateFunc if provided, otherwise returns the key.
func (m *MockLocalizer) Translate(key string, args ...map[string]any) string {
	if m.TranslateFunc != nil {
		return m.TranslateFunc(key, args...)
	}
	return key
}

// GetCurrentLanguage returns the current language code.
// Uses the configured GetCurrentLanguageFunc if provided, otherwise returns "en".
func (m *MockLocalizer) GetCurrentLanguage() string {
	if m.GetCurrentLanguageFunc != nil {
		return m.GetCurrentLanguageFunc()
	}
	return "en"
}

// IsInitialized returns true if the localizer has been initialized.
// Uses the configured IsInitializedFunc if provided, otherwise returns true.
func (m *MockLocalizer) IsInitialized() bool {
	if m.IsInitializedFunc != nil {
		return m.IsInitializedFunc()
	}
	return true
}

// SetLanguage sets the current language for localization.
// Uses the configured SetLanguageFunc if provided, otherwise returns nil.
func (m *MockLocalizer) SetLanguage(lang string) error {
	if m.SetLanguageFunc != nil {
		return m.SetLanguageFunc(lang)
	}
	return nil
}

// MockLogger implements log.Logger for testing.
type MockLogger struct {
	SetLevelFunc func(_ log.Level)
	DebugFunc    func(_ string)
	InfoFunc     func(_ string)
	WarnFunc     func(_ string)
	ErrorFunc    func(_ string)
	DebugfFunc   func(_ string, _ ...any)
	InfofFunc    func(_ string, _ ...any)
	WarnfFunc    func(_ string, _ ...any)
	ErrorfFunc   func(_ string, _ ...any)
}

// SetLevel sets the logging level.
// Uses the configured SetLevelFunc if provided, otherwise does nothing.
func (m *MockLogger) SetLevel(level log.Level) {
	if m.SetLevelFunc != nil {
		m.SetLevelFunc(level)
	}
}

// Debug logs a debug message.
// Uses the configured DebugFunc if provided, otherwise does nothing.
func (m *MockLogger) Debug(msg string) {
	if m.DebugFunc != nil {
		m.DebugFunc(msg)
	}
}

// Info logs an info message.
// Uses the configured InfoFunc if provided, otherwise does nothing.
func (m *MockLogger) Info(msg string) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg)
	}
}

// Warn logs a warning message.
// Uses the configured WarnFunc if provided, otherwise does nothing.
func (m *MockLogger) Warn(msg string) {
	if m.WarnFunc != nil {
		m.WarnFunc(msg)
	}
}

// Error logs an error message.
// Uses the configured ErrorFunc if provided, otherwise does nothing.
func (m *MockLogger) Error(msg string) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg)
	}
}

// Debugf logs a formatted debug message.
// Uses the configured DebugfFunc if provided, otherwise does nothing.
func (m *MockLogger) Debugf(format string, v ...any) {
	if m.DebugfFunc != nil {
		m.DebugfFunc(format, v...)
	}
}

// Infof logs a formatted info message.
// Uses the configured InfofFunc if provided, otherwise does nothing.
func (m *MockLogger) Infof(format string, v ...any) {
	if m.InfofFunc != nil {
		m.InfofFunc(format, v...)
	}
}

// Warnf logs a formatted warning message.
// Uses the configured WarnfFunc if provided, otherwise does nothing.
func (m *MockLogger) Warnf(format string, v ...any) {
	if m.WarnfFunc != nil {
		m.WarnfFunc(format, v...)
	}
}

// Errorf logs a formatted error message.
// Uses the configured ErrorfFunc if provided, otherwise does nothing.
func (m *MockLogger) Errorf(format string, v ...any) {
	if m.ErrorfFunc != nil {
		m.ErrorfFunc(format, v...)
	}
}

// MockTimeProvider implements coretime.TimeProvider for testing.
type MockTimeProvider struct {
	NowFunc func() time.Time
}

// Now returns the current time.
// Uses the configured NowFunc if provided, otherwise returns the actual current time.
func (m *MockTimeProvider) Now() time.Time {
	if m.NowFunc != nil {
		return m.NowFunc()
	}
	return time.Now()
}

// MockFileReader implements config.FileReader for testing.
type MockFileReader struct {
	ReadFileFunc func(_ string) ([]byte, error)
}

// ReadFile reads the contents of the file at the specified path.
// Uses the configured ReadFileFunc if provided, otherwise returns nil.
func (m *MockFileReader) ReadFile(path string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}
	return nil, nil
}

// MockConfigParser implements config.Parser for testing.
type MockConfigParser struct {
	ParseFunc func(_ []byte) (*config.Config, error)
}

// Parse parses configuration data into a Config struct.
// Uses the configured ParseFunc if provided, otherwise returns nil.
func (m *MockConfigParser) Parse(data []byte) (*config.Config, error) {
	if m.ParseFunc != nil {
		return m.ParseFunc(data)
	}
	return nil, nil
}

// MockDateProcessor implements date.DateProcessor for testing.
type MockDateProcessor struct {
	GetBestAvailableDateFunc func(_ map[string]any, _ string) (time.Time, error)
}

// GetBestAvailableDate extracts the best available date from metadata fields.
// Uses the configured GetBestAvailableDateFunc if provided, otherwise returns zero time.
func (m *MockDateProcessor) GetBestAvailableDate(fields map[string]any, filePath string) (time.Time, error) {
	if m.GetBestAvailableDateFunc != nil {
		return m.GetBestAvailableDateFunc(fields, filePath)
	}
	return time.Time{}, nil
}

// MockImageFinder implements imagesorter.ImageFinder for testing.
type MockImageFinder struct {
	FindFunc func(_ string) ([]string, error)
}

// Find discovers image files in the specified root path.
// Uses the configured FindFunc if provided, otherwise returns nil.
func (m *MockImageFinder) Find(rootPath string) ([]string, error) {
	if m.FindFunc != nil {
		return m.FindFunc(rootPath)
	}
	return nil, nil
}

// MockImageAnalyzer implements imagesorter.ImageAnalyzer for testing.
type MockImageAnalyzer struct {
	AnalyzeFunc func(_ string) (image.Image, error)
}

// Analyze extracts metadata from an image file.
// Uses the configured AnalyzeFunc if provided, otherwise returns empty Image.
func (m *MockImageAnalyzer) Analyze(filePath string) (image.Image, error) {
	if m.AnalyzeFunc != nil {
		return m.AnalyzeFunc(filePath)
	}
	return image.Image{}, nil
}

// MockImageProcessor implements imagesorter.ImageProcessor for testing.
type MockImageProcessor struct {
	ProcessFunc func(_ string) []image.Image
}

// Process analyzes all images in the specified directory.
// Uses the configured ProcessFunc if provided, otherwise returns empty slice.
func (m *MockImageProcessor) Process(dirPath string) []image.Image {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(dirPath)
	}
	return []image.Image{}
}

// MockExifReader implements exif.Reader for testing.
type MockExifReader struct {
	ReadExifFunc func(_ string) (map[string]any, error)
}

// ReadExif reads EXIF metadata from an image file.
// Uses the configured ReadExifFunc if provided, otherwise returns nil.
func (m *MockExifReader) ReadExif(filePath string) (map[string]any, error) {
	if m.ReadExifFunc != nil {
		return m.ReadExifFunc(filePath)
	}
	return nil, nil
}

// MockDateStrategy implements date.strategy.Strategy for testing.
type MockDateStrategy struct {
	ExtractFunc func(_ map[string]any, _ string) (time.Time, error)
}

// Extract extracts date information using a specific strategy.
// Uses the configured ExtractFunc if provided, otherwise returns zero time.
func (m *MockDateStrategy) Extract(fields map[string]any, filePath string) (time.Time, error) {
	if m.ExtractFunc != nil {
		return m.ExtractFunc(fields, filePath)
	}
	return time.Time{}, nil
}

// MockTask implements task.Task for testing.
type MockTask struct {
	RunFunc func() error
}

// Run executes the task.
// Uses the configured RunFunc if provided, otherwise returns nil.
func (m *MockTask) Run() error {
	if m.RunFunc != nil {
		return m.RunFunc()
	}
	return nil
}

// MockFileUtils implements filesystem.FileUtils for testing.
type MockFileUtils struct {
	CopyFileFunc  func(sourcePath, destinationPath string) error
	ExistsFunc    func(_ string) bool
	EnsureDirFunc func(_ string) error
}

// CopyFile copies a file from source to destination path.
// Uses the configured CopyFileFunc if provided, otherwise returns nil.
func (m *MockFileUtils) CopyFile(sourcePath, destinationPath string) error {
	if m.CopyFileFunc != nil {
		return m.CopyFileFunc(sourcePath, destinationPath)
	}
	return nil
}

// Exists checks if a file or directory exists at the specified path.
// Uses the configured ExistsFunc if provided, otherwise returns false.
func (m *MockFileUtils) Exists(path string) bool {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(path)
	}
	return false
}

// EnsureDir creates the directory if it doesn't exist.
// Uses the configured EnsureDirFunc if provided, otherwise returns nil.
func (m *MockFileUtils) EnsureDir(path string) error {
	if m.EnsureDirFunc != nil {
		return m.EnsureDirFunc(path)
	}
	return nil
}

// MockConfigApplier implements handlers.ConfigApplier for testing.
type MockConfigApplier[T any] struct {
	ApplyFunc func(_ T, _ *config.Config)
}

// Apply applies configuration from source to destination config.
// Uses the configured ApplyFunc if provided, otherwise does nothing.
func (m *MockConfigApplier[T]) Apply(src T, dest *config.Config) {
	if m.ApplyFunc != nil {
		m.ApplyFunc(src, dest)
	}
}

// MockConfigValidator implements handlers.ConfigValidator for testing.
type MockConfigValidator[T any] struct {
	ValidateFunc func(_ T) error
}

// Validate validates the configuration.
// Uses the configured ValidateFunc if provided, otherwise returns nil.
func (m *MockConfigValidator[T]) Validate(cfg T) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(cfg)
	}
	return nil
}

// MockTaskExecutor implements handlers.TaskExecutor for testing.
type MockTaskExecutor struct {
	ExecuteFunc func(_ string, _ config.Config) error
}

// Execute executes a named task with the provided configuration.
// Uses the configured ExecuteFunc if provided, otherwise returns nil.
func (m *MockTaskExecutor) Execute(taskName string, cfg config.Config) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(taskName, cfg)
	}
	return nil
}

// MockResource implements strategies.ActionResource for testing.
type MockResource struct {
	SetupFunc    func() error
	TeardownFunc func() error
}

// Setup initializes the resource.
// Uses the configured SetupFunc if provided, otherwise returns nil.
func (m *MockResource) Setup() error {
	if m.SetupFunc != nil {
		return m.SetupFunc()
	}
	return nil
}

// Teardown cleans up the resource.
// Uses the configured TeardownFunc if provided, otherwise returns nil.
func (m *MockResource) Teardown() error {
	if m.TeardownFunc != nil {
		return m.TeardownFunc()
	}
	return nil
}

// MockLoggerFactory implements log.LoggerFactory for testing.
type MockLoggerFactory struct {
	CreateLoggerFunc func(logFile string) (log.Logger, error)
}

// CreateLogger creates a new logger instance with the specified log file.
// Uses the configured CreateLoggerFunc if provided, otherwise returns a MockLogger.
func (m *MockLoggerFactory) CreateLogger(logFile string) (log.Logger, error) {
	if m.CreateLoggerFunc != nil {
		return m.CreateLoggerFunc(logFile)
	}
	return &MockLogger{}, nil
}

// MockConfig implements config.Config for testing.
type MockConfig struct {
	// Add fields as needed for testing
}

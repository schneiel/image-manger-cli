package di_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/schneiel/ImageManagerGo/internal/di"
)

// testCase represents a common test case structure for argument parser tests.
type testCase struct {
	name     string
	args     []string
	expected string
}

// runArgumentParserTests is a shared helper for testing argument parser functions.
func runArgumentParserTests(t *testing.T, testCases []testCase, testFunc func([]string) string) {
	t.Helper()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// t.Parallel() removed due to race condition in i18n.NewBundleLocalizer().
			result := testFunc(testCase.args)
			if result != testCase.expected {
				t.Errorf("got %v, want %v", result, testCase.expected)
			}
		})
	}
}

func TestDefaultArgumentParser_GetLanguage(t *testing.T) {
	t.Parallel()
	// t.Parallel() removed due to race condition in i18n.NewBundleLocalizer().

	mockLocalizer := &MockLocalizer{}
	parser, err := di.NewDefaultArgumentParser(mockLocalizer)
	require.NoError(t, err)

	tests := []testCase{
		{
			name:     "should return en when no language flag",
			args:     []string{"program", "command"},
			expected: "en",
		},
		{
			name:     "should return language when --lang flag present",
			args:     []string{"program", "--lang", "de", "command"},
			expected: "de",
		},
		{
			name:     "should return language when -l flag present",
			args:     []string{"program", "-l", "fr", "command"},
			expected: "fr",
		},
		{
			name:     "should return en when flag has no value",
			args:     []string{"program", "--lang", "command"},
			expected: "en",
		},
		{
			name:     "should return en when flag has no value with -l",
			args:     []string{"program", "-l", "command"},
			expected: "en",
		},
	}

	runArgumentParserTests(t, tests, parser.GetLanguage)
}

func TestDefaultArgumentParser_GetConfigPath(t *testing.T) {
	t.Parallel()
	// t.Parallel() removed due to race condition in i18n.NewBundleLocalizer().

	mockLocalizer := &MockLocalizer{}
	parser, err := di.NewDefaultArgumentParser(mockLocalizer)
	require.NoError(t, err)

	tests := []testCase{
		{
			name:     "should return empty when no config flag",
			args:     []string{"program", "command"},
			expected: "",
		},
		{
			name:     "should return path when --config flag present",
			args:     []string{"program", "--config", "/path/to/config.yaml", "command"},
			expected: "/path/to/config.yaml",
		},
		{
			name:     "should return path when -c flag present",
			args:     []string{"program", "-c", "config.yaml", "command"},
			expected: "config.yaml",
		},
		{
			name:     "should return empty when flag has no value",
			args:     []string{"program", "--config", "command"},
			expected: "",
		},
	}

	runArgumentParserTests(t, tests, parser.GetConfigPath)
}

// MockLocalizer is a simple mock for testing.
type MockLocalizer struct{}

func (m *MockLocalizer) Translate(key string, _ ...map[string]interface{}) string {
	return key
}

func (m *MockLocalizer) GetCurrentLanguage() string {
	return "en"
}

func (m *MockLocalizer) IsInitialized() bool {
	return true
}

func (m *MockLocalizer) SetLanguage(_ string) error {
	return nil
}

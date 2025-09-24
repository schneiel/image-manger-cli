// Package di provides dependency injection container and related utilities.
package di

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultArgumentParser provides functionality for parsing command line arguments.
type DefaultArgumentParser struct {
	localizer i18n.Localizer
}

// NewDefaultArgumentParser creates a new DefaultArgumentParser with injected dependencies.
func NewDefaultArgumentParser(localizer i18n.Localizer) (*DefaultArgumentParser, error) {
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}

	return &DefaultArgumentParser{localizer: localizer}, nil
}

// GetLanguage extracts the language argument from command line arguments.
// Returns the language code or "en" as default if not specified.
func (p *DefaultArgumentParser) GetLanguage(args []string) string {
	for i, arg := range args {
		if (arg == "--lang" || arg == "-l") && i+1 < len(args) {
			// Check if the next argument is not another flag and is a valid language code
			nextArg := args[i+1]
			if !isFlag(nextArg) && isValidLanguageCode(nextArg) {
				return nextArg
			}
		}
	}

	return "en" // Default fallback
}

// GetConfigPath extracts the config path argument from command line arguments.
// Returns the config path or empty string if not specified.
func (p *DefaultArgumentParser) GetConfigPath(args []string) string {
	for i, arg := range args {
		if (arg == "--config" || arg == "-c") && i+1 < len(args) {
			// Check if the next argument is not another flag and looks like a path
			nextArg := args[i+1]
			if !isFlag(nextArg) && looksLikePath(nextArg) {
				return nextArg
			}
		}
	}

	return "" // Default fallback
}

// isFlag checks if a string is a command line flag.
func isFlag(arg string) bool {
	return arg != "" && arg[0] == '-'
}

// isValidLanguageCode checks if a string is a valid language code.
func isValidLanguageCode(code string) bool {
	// Simple validation for common language codes
	validCodes := map[string]bool{
		"en": true,
		"de": true,
		"fr": true,
		"es": true,
		"it": true,
		"pt": true,
		"ru": true,
		"ja": true,
		"ko": true,
		"zh": true,
	}

	return validCodes[code]
}

// looksLikePath checks if a string looks like a file path.
func looksLikePath(arg string) bool {
	// Simple heuristic: paths typically contain slashes, dots, or have file extensions
	if arg == "" {
		return false
	}

	if arg[0] == '/' || arg[0] == '.' || arg[0] == '~' {
		return true
	}

	// Check file extensions with bounds checking
	argLen := len(arg)
	if argLen >= 4 && arg[argLen-4:] == ".yml" {
		return true
	}
	if argLen >= 5 && arg[argLen-5:] == ".yaml" {
		return true
	}
	if argLen >= 5 && arg[argLen-5:] == ".json" {
		return true
	}
	if argLen >= 5 && arg[argLen-5:] == ".toml" {
		return true
	}

	return false
}

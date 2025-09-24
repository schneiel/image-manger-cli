// Package log provides SOLID-compliant logging components.
package log

import (
	"errors"
	"fmt"
	"strings"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultLevelParser provides functionality for parsing log levels.
type DefaultLevelParser struct {
	localizer i18n.Localizer
}

// NewDefaultLevelParser creates a new DefaultLevelParser with injected dependencies.
func NewDefaultLevelParser(localizer i18n.Localizer) (*DefaultLevelParser, error) {
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	return &DefaultLevelParser{localizer: localizer}, nil
}

// Parse converts a string to a Level with localized error messages.
func (p *DefaultLevelParser) Parse(s string) (Level, error) {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return DEBUG, fmt.Errorf("%s", p.localizer.Translate("InvalidLogLevel",
			map[string]interface{}{"Level": s}))
	}
}

package time

import "time"

// DefaultTimeProvider implements TimeProvider using the standard time package.
type DefaultTimeProvider struct{}

// NewDefaultTimeProvider creates a new DefaultTimeProvider.
func NewDefaultTimeProvider() TimeProvider {
	return &DefaultTimeProvider{}
}

// Now returns the current time.
func (p *DefaultTimeProvider) Now() time.Time {
	return time.Now()
}

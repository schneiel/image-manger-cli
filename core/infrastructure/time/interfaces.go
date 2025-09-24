// Package time provides time-related interfaces and implementations for better testability.
package time

import "time"

// TimeProvider abstracts time operations for better testability
//
//nolint:revive // TimeProvider name is clear and intentional despite package stuttering
type TimeProvider interface {
	// Now returns the current time
	Now() time.Time
}

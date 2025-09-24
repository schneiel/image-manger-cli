package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultTimeProvider(t *testing.T) {
	t.Parallel()

	provider := NewDefaultTimeProvider()

	assert.NotNil(t, provider)
	assert.IsType(t, &DefaultTimeProvider{}, provider)
}

func TestDefaultTimeProvider_Now(t *testing.T) {
	t.Parallel()

	provider := NewDefaultTimeProvider()

	// Get the current time
	currentTime := provider.Now()

	// Verify that the returned time is not zero
	assert.False(t, currentTime.IsZero(), "Now() should return a non-zero time")

	// Verify that the returned time is close to the actual current time
	// (allowing for a small time difference due to execution time)
	actualTime := time.Now()
	timeDiff := actualTime.Sub(currentTime)
	assert.Less(t, timeDiff, time.Second, "Now() should return a time close to the actual current time")
}

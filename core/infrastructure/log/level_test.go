package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		level Level
		want  string
	}{
		{"debug", DEBUG, "DEBUG"},
		{"info", INFO, "INFO"},
		{"warn", WARN, "WARN"},
		{"error", ERROR, "ERROR"},
		{"unknown", Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			got := tt.level.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLevel_Constants(t *testing.T) {
	t.Parallel()
	assert.Equal(t, DEBUG, Level(0))
	assert.Equal(t, INFO, Level(1))
	assert.Equal(t, WARN, Level(2))
	assert.Equal(t, ERROR, Level(3))
}

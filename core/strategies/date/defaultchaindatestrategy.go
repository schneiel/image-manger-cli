package date

import "time"

// DefaultChainDateStrategy tries multiple date strategies in order.
type DefaultChainDateStrategy struct {
	strategies []Strategy
}

// NewDefaultChainDateStrategy creates a new DefaultChainDateStrategy.
func NewDefaultChainDateStrategy(strategies ...Strategy) *DefaultChainDateStrategy {
	return &DefaultChainDateStrategy{strategies: strategies}
}

// Extract extracts date from image metadata fields and file path.
func (c *DefaultChainDateStrategy) Extract(fields map[string]interface{}, filePath string) (time.Time, error) {
	for _, strategy := range c.strategies {
		if date, err := strategy.Extract(fields, filePath); err == nil && !date.IsZero() {
			return date, nil
		}
	}
	return time.Time{}, nil
}

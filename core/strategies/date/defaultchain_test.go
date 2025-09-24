package date

import (
	"errors"
	"testing"
	"time"

	"github.com/schneiel/ImageManagerGo/core/testutils"
)

func TestNewDefaultChainDateStrategy(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)

	if chain == nil {
		t.Fatal("Expected chain to be created, got nil")
	}

	if len(chain.strategies) != 2 {
		t.Errorf("Expected 2 strategies, got %d", len(chain.strategies))
	}

	if chain.strategies[0] != mockStrategy1 {
		t.Error("Expected first strategy to be injected")
	}
	if chain.strategies[1] != mockStrategy2 {
		t.Error("Expected second strategy to be injected")
	}
}

func TestNewDefaultChainDateStrategy_EmptyStrategies(t *testing.T) {
	t.Parallel()
	chain := NewDefaultChainDateStrategy()

	if chain == nil {
		t.Fatal("Expected chain to be created, got nil")
	}

	if len(chain.strategies) != 0 {
		t.Errorf("Expected 0 strategies, got %d", len(chain.strategies))
	}
}

func TestDefaultChainDateStrategy_Extract_FirstStrategySuccess(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, errors.New("should not be called")
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDefaultChainDateStrategy_Extract_SecondStrategySuccess(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, nil // First strategy returns zero time
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDefaultChainDateStrategy_Extract_AllStrategiesFail(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, errors.New("strategy 1 failed")
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, errors.New("strategy 2 failed")
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return zero time when all strategies fail
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultChainDateStrategy_Extract_AllStrategiesReturnZero(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, nil
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, nil
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return zero time when all strategies return zero
	zeroTime := time.Time{}
	if !result.Equal(zeroTime) {
		t.Errorf("Expected zero time, got %v", result)
	}
}

func TestDefaultChainDateStrategy_Extract_MixedResults(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}
	mockStrategy3 := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, errors.New("strategy 1 failed")
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, nil // Strategy 2 returns zero time
	}
	mockStrategy3.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2, mockStrategy3)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDefaultChainDateStrategy_Extract_NilFields(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, errors.New("should not be called")
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(nil, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDefaultChainDateStrategy_Extract_EmptyFields(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, errors.New("should not be called")
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	fields := map[string]interface{}{}
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, result)
	}
}

func TestDefaultChainDateStrategy_Extract_ComplexFilePath(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	expectedDate := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return expectedDate, nil
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		return time.Time{}, errors.New("should not be called")
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}

	testCases := []string{
		"/path/to/test/image.jpg",
		"/path/with spaces/file name.png",
		"/path/with/unicode/测试图片.jpeg",
		"/path/with/dots/file.name.with.dots.tiff",
		"/path/with/numbers/IMG_2023_001.jpg",
		"/path/with/special/chars/file@#$%^&*().png",
	}

	for _, filePath := range testCases {
		t.Run(filePath, func(_ *testing.T) {
			result, err := chain.Extract(fields, filePath)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !result.Equal(expectedDate) {
				t.Errorf("Expected %v, got %v", expectedDate, result)
			}
		})
	}
}

func TestDefaultChainDateStrategy_Extract_StrategyOrder(t *testing.T) {
	t.Parallel()
	mockStrategy1 := &testutils.MockDateStrategy{}
	mockStrategy2 := &testutils.MockDateStrategy{}

	date1 := time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC)
	date2 := time.Date(2023, 1, 16, 14, 20, 0, 0, time.UTC)

	callCount := 0
	mockStrategy1.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		callCount++
		return date1, nil
	}
	mockStrategy2.ExtractFunc = func(_ map[string]interface{}, _ string) (time.Time, error) {
		callCount++
		return date2, nil
	}

	chain := NewDefaultChainDateStrategy(mockStrategy1, mockStrategy2)
	fields := map[string]interface{}{
		"DateTimeOriginal": "2023:01:15 12:30:00",
	}
	filePath := "/path/to/test/image.jpg"

	result, err := chain.Extract(fields, filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return first strategy's result
	if !result.Equal(date1) {
		t.Errorf("Expected %v, got %v", date1, result)
	}

	// Should only call first strategy
	if callCount != 1 {
		t.Errorf("Expected 1 strategy call, got %d", callCount)
	}
}

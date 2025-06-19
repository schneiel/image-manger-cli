package date

import (
	"fmt"
	"strings"
	"time"
)

// ExifStrategy extracts a date from a specific EXIF field.
type ExifStrategy struct {
	fieldName FieldName
	layout    string
}

// GetDate attempts to parse the date from the EXIF data provided in the fields map.
func (s *ExifStrategy) GetDate(fields map[string]interface{}, _ string) (time.Time, error) {
	dateValue, ok := fields[string(s.fieldName)]
	if !ok {
		return time.Time{}, fmt.Errorf("field %s not found", s.fieldName)
	}

	dateString, ok := dateValue.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("field %s is not a string", s.fieldName)
	}

	dateString = strings.TrimRightFunc(dateString, func(r rune) bool {
		return r < 32 || r == 0
	})

	if strings.HasSuffix(dateString, "Z") {
		return time.Parse("2006:01:02 15:04:05Z", dateString)
	}
	return time.Parse(s.layout, dateString)
}

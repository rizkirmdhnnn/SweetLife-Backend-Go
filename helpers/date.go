package helper

import (
	"time"
)

// ParsedDate parses a string date into a time.Time object
func ParsedDate(value string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}
	return parsedDate, nil
}

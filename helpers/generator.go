package helper

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// GenerateFileName generates a unique file name with the given extension
func GenerateFileName(extension string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x%s", b, strings.ToLower(extension))
}

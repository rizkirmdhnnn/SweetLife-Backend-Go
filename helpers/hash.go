package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateHash generates a hash from the given email and key
func GenerateHash(email, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

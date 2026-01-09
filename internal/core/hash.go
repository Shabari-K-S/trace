package core

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashContent computes the SHA256 hash of the given content.
func HashContent(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// HashString computes the SHA256 hash of a string.
func HashString(s string) string {
	return HashContent([]byte(s))
}

// ShortHash returns the first 7 characters of a hash.
func ShortHash(hash string) string {
	if len(hash) >= 7 {
		return hash[:7]
	}
	return hash
}

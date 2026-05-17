package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomHex generates a cryptographically random hex string of n bytes.
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Ptr returns a pointer to the given value — useful for optional fields.
func Ptr[T any](v T) *T { return &v }

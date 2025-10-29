package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeSHA256 returns the SHA-256 digest of the input bytes.
func ComputeSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

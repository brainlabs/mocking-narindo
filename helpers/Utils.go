package helpers

import (
	"crypto/sha1"
	"encoding/hex"
)

// HashSHA1 hashing sha256 plaintext
func HashSHA1(plainText string) string {

	h := sha1.New()
	h.Write([]byte(plainText))

	return hex.EncodeToString(h.Sum(nil))
}

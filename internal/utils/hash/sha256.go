package hash

import (
	"crypto/sha256"
	"fmt"
)

func HashSha256(plainText string) (hashStr string) {
	h := sha256.New()
	h.Write([]byte(plainText))

	return fmt.Sprintf("%x", h.Sum(nil))
}

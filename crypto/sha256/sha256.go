package sha256

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	hashString := hex.EncodeToString(hashValue)
	return hashString
}

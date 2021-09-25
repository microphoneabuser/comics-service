package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func GenHash(password string) string {
	salt := os.Getenv("salt")
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

package auth

import (
	"crypto/sha256"
	"fmt"
)

// Hash function generates hash with SHA256 algorithm.
func Hash(password string) string {
	hash := sha256.Sum256([]byte(password))

	return fmt.Sprintf("%x", hash[:])
}

// CheckPassword checks password and hash to see if they are equal.
func CheckPassword(password, hash string) bool {
	passwordHash := sha256.Sum256([]byte(password))

	return fmt.Sprintf("%x", passwordHash[:]) == hash
}

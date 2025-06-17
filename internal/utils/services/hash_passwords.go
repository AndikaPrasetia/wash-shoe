package utils

import "golang.org/x/crypto/bcrypt"

// Hashing password
func HashPassword(password string) (string, error) {
	// DefaultCost = 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

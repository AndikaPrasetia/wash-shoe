// Package utils: for reuseable utilities
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword is for hasing password
func HashPassword(password string) (string, error) {
	// DefaultCost = 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash is for checking hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

package utils

import (
	"errors"
	"fmt"
	"unicode"
)

// ValidatePassword is for password security
// min = panjang minimal (bisa 8, 10, 12 – terserah kebijakan)
// max = panjang maksimal (umumnya 64 untuk batasi DoS regex)
func ValidatePassword(password string, min, max int) error {
	n := len(password)
	if n < min {
		return fmt.Errorf("password must be at least %d characters", min)
	}
	if n > max {
		return fmt.Errorf("password must be at most %d characters", max)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	switch {
	case !hasUpper:
		return errors.New("password must contain at least one uppercase letter")
	case !hasLower:
		return errors.New("password must contain at least one lowercase letter")
	case !hasDigit:
		return errors.New("password must contain at least one number")
	case !hasSpecial:
		return errors.New("password must contain at least one special character")
	}

	return nil
}

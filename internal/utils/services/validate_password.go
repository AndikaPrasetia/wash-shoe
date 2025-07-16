package utils

import (
	"errors"
	"unicode"
)

// ValidatePassword is for password security
// min = panjang minimal (bisa 8, 10, 12 – terserah kebijakan)
// max = panjang maksimal (umumnya 64 untuk batasi DoS regex)
func ValidatePassword(password string, min, max int) error {
	if len(password) < min || len(password) > max {
		return errors.New("password must be between " + string(rune(min)) + "-" + string(rune(max)) + " characters")
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

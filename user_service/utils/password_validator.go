package utils

import "unicode"

// ValidatePassword checks if a password meets the following rules:
// - Between 8 and 72 characters
// - At least one uppercase letter
// - At least one digit
// - At least one symbol
// - ASCII characters only (no unicode)
func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 72 {
		return false
	}

	var hasUpper, hasDigit, hasSymbol bool
	for _, ch := range password {
		if ch > 127 {
			return false // non-ASCII unicode not allowed
		}
		if unicode.IsUpper(ch) {
			hasUpper = true
		}
		if unicode.IsDigit(ch) {
			hasDigit = true
		}
		if unicode.IsPunct(ch) || unicode.IsSymbol(ch) {
			hasSymbol = true
		}
	}

	return hasUpper && hasDigit && hasSymbol
}

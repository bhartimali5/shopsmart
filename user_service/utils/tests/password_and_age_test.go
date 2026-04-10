package tests

import (
	"strings"
	"testing"

	"example.com/rest-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"ShouldReturnFalse_WhenPasswordIsEmpty", "", false},
		{"ShouldReturnFalse_WhenPasswordIsTooShort", "Abc1!", false},
		{"ShouldReturnFalse_WhenPasswordHasNoDigits", "Abcdefg!", false},
		{"ShouldReturnFalse_WhenPasswordHasNoSymbols", "Abcdefg1", false},
		{"ShouldReturnFalse_WhenPasswordHasNoUppercase", "abcdefg1!", false},
		{"ShouldReturnFalse_WhenPasswordExceedsMaxLength", strings.Repeat("A1!", 25), false}, // 75 chars
		{"ShouldReturnFalse_WhenPasswordContainsUnicode", "Abcdéf1!", false},
		{"ShouldReturnTrue_WhenPasswordMeetsAllRules", "Abcdef1!", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.ValidatePassword(tc.password)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsEligible(t *testing.T) {
	// Boundary value analysis:
	// - 17: just below the boundary (must be rejected)
	// - 18: exact boundary (must be accepted)
	// - 19: just above the boundary (must be accepted)
	// - 0:  zero age (invalid input, must be rejected)
	// - -1: negative age (invalid input, must be rejected)
	tests := []struct {
		name     string
		age      int
		expected bool
	}{
		{"ShouldReturnFalse_WhenAgeIsJustBelowBoundary", 17, false},  // lower boundary - 1
		{"ShouldReturnTrue_WhenAgeIsExactlyAtBoundary", 18, true},    // exact boundary
		{"ShouldReturnTrue_WhenAgeIsJustAboveBoundary", 19, true},    // lower boundary + 1
		{"ShouldReturnFalse_WhenAgeIsZero", 0, false},                // invalid input
		{"ShouldReturnFalse_WhenAgeIsNegative", -1, false},           // invalid input
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.IsEligible(tc.age)
			assert.Equal(t, tc.expected, result)
		})
	}
}

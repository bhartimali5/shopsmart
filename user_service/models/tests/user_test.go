// Unit testcases for the User model.
package models_test

import (
	"testing"

	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidatePassword_ShouldReturnTrue_WhenCorrectPasswordIsProvided(t *testing.T) {
	test_user := models.User{
		Email:    "test_user@example.com",
		Password: "secret123",
	}

	hashedPassword, err := utils.HashPassword(test_user.Password)
	result, err := test_user.ValidatePassword(hashedPassword)
	assert.NoError(t, err)
	assert.True(t, *result)
}

func TestValidatePassword_ShouldReturnNil_WhenIncorrectPasswordIsProvided(t *testing.T) {
	test_user := models.User{
		Email:    "test_user@example.com",
		Password: "secret123",
	}

	hashedPassword, err := utils.HashPassword("wrongpassword")
	result, err := test_user.ValidatePassword(hashedPassword)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidatePassword_ShouldReturnError_WhenGivenEmptyHashedPassword(t *testing.T) {
	test_user := models.User{
		Email:    "test_user@example.com",
		Password: "secret123",
	}

	// Simulate an error by passing an empty hashed password
	result, err := test_user.ValidatePassword("")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidatePassword_ShouldReturnNil_WhenGivenSimilarPasswordButNotExact(t *testing.T) {
	test_user := models.User{
		Email:    "test_user@example.com",
		Password: "secret123",
	}

	hashedPassword, err := utils.HashPassword("Secret123") // Note the capital 'S'
	assert.NoError(t, err)
	result, err := test_user.ValidatePassword(hashedPassword)
	assert.Error(t, err)
	assert.Nil(t, result) // Should be false because of case sensitivity
}

func TestValidatePassword_ShouldReturnError_WhenGivenInvalidHashedPassword(t *testing.T) {
	test_user := models.User{
		Email:    "test_user@example.com",
		Password: "secret123",
	}

	// Simulate an error by passing an invalid hashed password
	result, err := test_user.ValidatePassword("invalid_hashed_password")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidatePassword_ShouldReturnNilWhenGivenEmptyPassword(t *testing.T) {
	test_user := models.User{
		Email:    "test_user@example.com",
		Password: "secret123",
	}
	hashedPassword, err := utils.HashPassword("secret123 ") // Note the trailing space
	result, err := test_user.ValidatePassword(hashedPassword)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestValidatePassword_ShouldReturnError_WhenPasswordHasTrailingSpace(t *testing.T) {
	test_user := models.User{
		Email:    "test_user@example.com",
		Password: "",
	}
	// Simulate an error by passing an empty password
	result, err := test_user.ValidatePassword("")
	assert.Error(t, err)
	assert.Nil(t, result)
}

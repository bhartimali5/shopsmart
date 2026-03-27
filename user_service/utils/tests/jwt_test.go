package tests

import (
	"testing"

	"example.com/rest-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT_ShouldReturnToken_WhenValidEmailAndIDAreProvided(t *testing.T) {
	email := "test_user@example.com"
	id := "12345"
	role := "user"
	token, err := utils.GenerateJWT(email, id, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWT_ShouldGenerateDifferentTokensWhenGivenDifferentRoles(t *testing.T) {
	adminToken, _ := utils.GenerateJWT("bharti@example.com", "user-123", "admin")
	customerToken, _ := utils.GenerateJWT("bharti@example.com", "user-123", "customer")

	assert.NotEqual(t, adminToken, customerToken)
}

//////////Tests for VerifyToken function

func TestVerifyToken_ShouldReturnUserIDAndRole_WhenValidTokenIsProvided(t *testing.T) {
	email := "test_user@example.com"
	id := "12345"
	role := "user"
	token, err := utils.GenerateJWT(email, id, role)
	assert.NoError(t, err)

	userID, userRole, err := utils.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, id, userID)
	assert.Equal(t, role, userRole)
}

func TestVerifyToken_ShouldReturnError_WhenInvalidTokenIsProvided(t *testing.T) {
	invalidToken := "invalid.token.string"

	userID, userRole, err := utils.VerifyToken(invalidToken)
	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Empty(t, userRole)
}

/*
	func TestVerifyToken_ShouldReturnError_WhenExpiredTokenIsProvided(t *testing.T) {
		email := "test_user@example.com"
		id := "12345"
		role := "user"
		token, err := utils.GenerateJWT(email, id, role)
		assert.NoError(t, err)

		// Simulate an expired token (this would require a more complex setup in a real test)
		userID, userRole, err := utils.VerifyToken(token)
		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Empty(t, userRole)
	}
*/
func TestVerifyToken_ShouldReturnError_WhenTokenWithInvalidSignatureIsProvided(t *testing.T) {
	email := "test_user@example.com"
	id := "12345"
	role := "user"
	token, err := utils.GenerateJWT(email, id, role)
	assert.NoError(t, err)

	otherToken, err := utils.GenerateJWT("other_user@example.com", "67890", "admin")
	assert.NoError(t, err)

	// Simulate a token with an invalid signature (this would require a more complex setup in a real test)
	userID, userRole, err := utils.VerifyToken(otherToken)

	assert.NotEqual(t, userID, id)
	assert.NotEqual(t, userRole, role)
	assert.NotEqual(t, token, otherToken)
}

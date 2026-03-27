package tests

import (
	"database/sql"
	"testing"

	"example.com/rest-api/db"
	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	sqliteDB, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	db.DB = sqliteDB
	db.CreateTables() // ← reuse your existing function, don't duplicate SQL
	t.Cleanup(func() { sqliteDB.Close() })
}

func SeedTestData(t *testing.T, email, password, role string) models.User {
	t.Helper()
	// Seed test data
	testUser := models.User{
		Email:    email,
		Password: password,
		Role:     role,
	}
	err := testUser.Save()
	require.NoError(t, err)
	return testUser
}

//Happy path test for user login

func TestUserLoginIntegration_whenValidCredentialsAreProvided(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "password123"
	role := "user"
	SeedTestData(t, email, password, role)
	// Simulate login
	loginData := models.User{
		Email:    email,
		Password: password,
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.NoError(t, err)
	isValid, err := loginData.ValidatePassword(hashedPassword)
	require.NoError(t, err)
	assert.True(t, *isValid)

	user_role, err := loginData.FetchUserRoleFromEmail()
	require.NoError(t, err)
	assert.Equal(t, role, user_role)

	token, err := utils.GenerateJWT(loginData.Email, loginData.ID, user_role)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotContains(t, token, "password123") // Ensure password is not in the token
	// To ensure the token is not plain text, we can check if it contains a dot (.) which is typical for JWTs
	assert.Contains(t, token, ".")

	// assert role is not empty
	assert.NotEmpty(t, user_role)
	assert.Equal(t, role, user_role)
}

// Additional edge case tests for user login

func TestUserLoginIntegration_whenValidEmailButInvalidPasswordIsProvided(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "password123"
	role := "user"
	SeedTestData(t, email, password, role)

	// Simulate login with invalid password
	invalidPassword := "wrongpassword"
	loginData := models.User{
		Email:    email,
		Password: invalidPassword,
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.NoError(t, err)
	isValid, err := loginData.ValidatePassword(hashedPassword)
	require.Error(t, err)
	assert.Nil(t, isValid)
	assert.Error(t, err)
}

func TestUserLoginIntegration_whenInvalidEmailIsProvided(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "password123"
	role := "user"
	SeedTestData(t, email, password, role)

	// Simulate login with invalid email
	loginData := models.User{
		Email:    "nonexistent@example.com",
		Password: password,
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.Error(t, err)
	assert.Empty(t, hashedPassword) // No hashed password should be returned for non-existent email
	isValid, err := loginData.ValidatePassword(hashedPassword)
	assert.Error(t, err)   // Should return an error since the email doesn't exist
	assert.Nil(t, isValid) // isValid should be false due to the error

}

func TestUserLoginIntegration_whenEmptyEmailAndPasswordAreProvided(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "secret123"
	role := "user"
	SeedTestData(t, email, password, role)
	// Simulate login with empty email and password
	loginData := models.User{
		Email:    "",
		Password: "",
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.Error(t, err)
	assert.Empty(t, hashedPassword) // No hashed password should be returned for empty email
	isValid, err := loginData.ValidatePassword(hashedPassword)
	assert.Error(t, err)   // Should return an error since the email is empty
	assert.Nil(t, isValid) // isValid should be false due to the error
}

func TestUserLoginIntegration_whenValidEmailButEmptyPasswordIsProvided(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "secret123"
	role := "user"
	SeedTestData(t, email, password, role)

	// Simulate login with valid email but empty password
	loginData := models.User{
		Email:    email,
		Password: "",
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword) // Hashed password should be returned for valid email
	isValid, err := loginData.ValidatePassword(hashedPassword)
	assert.Error(t, err)   // Should return an error since the password is empty
	assert.Nil(t, isValid) // isValid should be false due to the error
}

func TestUserLoginIntegration_whenEmptyEmailButValidPasswordIsProvided(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "secret123"
	role := "user"
	SeedTestData(t, email, password, role)

	// Simulate login with empty email but valid password
	loginData := models.User{
		Email:    "",
		Password: password,
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.Error(t, err)
	assert.Empty(t, hashedPassword) // No hashed password should be returned for empty email
	isValid, err := loginData.ValidatePassword(hashedPassword)
	assert.Error(t, err)   // Should return an error since the email is empty
	assert.Nil(t, isValid) // isValid should be false due to the error
}

func TestUserLoginIntegration_WhenValidCredentialsAreProvided_ShouldGenerateJWTWithCorrectClaims(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "secret123"
	role := "user"
	SeedTestData(t, email, password, role)
	// Simulate login
	loginData := models.User{
		Email:    email,
		Password: password,
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.NoError(t, err)
	isValid, err := loginData.ValidatePassword(hashedPassword)
	require.NoError(t, err)
	assert.True(t, *isValid)

	user_role, err := loginData.FetchUserRoleFromEmail()
	require.NoError(t, err)
	assert.Equal(t, role, user_role)

	token, err := utils.GenerateJWT(loginData.Email, loginData.ID, user_role)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Decode the token to verify claims
	user_id, user_role, err := utils.VerifyToken(token)
	require.NoError(t, err)
	assert.Equal(t, loginData.ID, user_id)
	assert.Equal(t, role, user_role)
}

func TestUserLoginIntegration_WhenGivenUserCredentials_ShouldNotReturnAdminRoleForRegularUser(t *testing.T) {
	setupTestDB(t)
	email := "test@example.com"
	password := "secret123"
	role := "user"
	SeedTestData(t, email, password, role)
	// Simulate login
	loginData := models.User{
		Email:    email,
		Password: password,
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.NoError(t, err)
	isValid, err := loginData.ValidatePassword(hashedPassword)
	require.NoError(t, err)
	assert.True(t, *isValid)

	user_role, err := loginData.FetchUserRoleFromEmail()
	require.NoError(t, err)
	assert.NotEqual(t, user_role, "admin") // Ensure regular user does not get admin role
	assert.Equal(t, role, user_role)       // Ensure the correct role is returned
}

func TestUserLoginIntegration_WhenGivenValidCredentials_ShouldNotFetchDetailsOfOtherUser(t *testing.T) {
	setupTestDB(t)
	email1 := "test1@example.com"
	email2 := "test2@example.com"
	password := "secret123"
	role := "user"
	SeedTestData(t, email1, password, role)
	SeedTestData(t, email2, password, "admin")
	// Simulate login with first user's credentials
	loginData := models.User{
		Email:    email1,
		Password: password,
	}
	hashedPassword, err := loginData.FetchPasswordHashFromEmail()
	require.NoError(t, err)
	isValid, err := loginData.ValidatePassword(hashedPassword)
	require.NoError(t, err)
	assert.True(t, *isValid)

	user_role, err := loginData.FetchUserRoleFromEmail()
	require.NoError(t, err)
	assert.NotEqual(t, user_role, "admin") // Ensure first user does not get admin role
	assert.Equal(t, role, user_role)
}

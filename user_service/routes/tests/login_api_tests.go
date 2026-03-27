package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/rest-api/db"
	"example.com/rest-api/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// 1. Point db.DB to an in-memory SQLite instance (no file, no leftover state)
	var err error
	db.DB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic("Could not connect to test database: " + err.Error())
	}

	// 2. Create the tables your handlers depend on
	db.CreateTables()

	// 3. Run all tests, then exit
	os.Exit(m.Run())
}

func setupRouter(t *testing.T) *gin.Engine {
	// No DB init here — TestMain already did it once for the whole suite
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.RegisterRoutes(router) // or however you wire up your routes
	return router
}

func doPost(router http.Handler, path string, body map[string]string) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func parseBody(w *httptest.ResponseRecorder) map[string]string {
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	return body
}

func TestLogin_WhenGivenValidCredentials_ReturnsOKStatus(t *testing.T) {
	router := setupRouter(t)
	// First, create a user to login with
	signUpBody := map[string]string{
		"email":    "test_user@example.com",
		"password": "password123",
	}
	w := doPost(router, "/signup", signUpBody)

	// Now, attempt to login with the same credentials
	loginBody := map[string]string{
		"email":    "test_user@example.com",
		"password": "password123",
	}
	w = doPost(router, "/login", loginBody)
	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(w)
	assert.Contains(t, body["message"], "Login successful")
	assert.NotEmpty(t, body["token"])

}

func TestLogin_WhenGivenInvalidCredentials_ReturnsUnauthorizedStatus(t *testing.T) {
	router := setupRouter(t)
	loginBody := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}
	w := doPost(router, "/login", loginBody)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_WhenGivenInvalidJSON_ReturnsBadRequestStatus(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("{invalid-json}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_WhenGivenValidCredentials_ReturnsJWTToken(t *testing.T) {
	router := setupRouter(t)
	// First, create a user to login with
	signUpBody := map[string]string{
		"email":    "test_user@example.com",
		"password": "password123",
	}
	w := doPost(router, "/signup", signUpBody)

	// Now, attempt to login with the same credentials
	loginBody := map[string]string{
		"email":    "test_user@example.com",
		"password": "password123",
	}
	w = doPost(router, "/login", loginBody)
	assert.Equal(t, http.StatusOK, w.Code)
	body := parseBody(w)
	assert.Contains(t, body["message"], "Login successful")
	assert.NotEmpty(t, body["token"])
}

func TestLogin_WhenGivenEmptyEmailAndPassword_ReturnsBadRequestStatus(t *testing.T) {
	router := setupRouter(t)
	loginBody := map[string]string{
		"email":    "",
		"password": "",
	}
	w := doPost(router, "/login", loginBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_WhenGivenEmptyEmail_ReturnsBadRequestStatus(t *testing.T) {
	router := setupRouter(t)
	loginBody := map[string]string{
		"email":    "",
		"password": "password123",
	}
	w := doPost(router, "/login", loginBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_WhenGivenEmptyPassword_ReturnsBadRequestStatus(t *testing.T) {
	router := setupRouter(t)
	loginBody := map[string]string{
		"email":    "test_user@example.com",
		"password": "",
	}
	w := doPost(router, "/login", loginBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_WhenGivenInvalidEmail_ReturnsUnauthorizedStatus(t *testing.T) {
	router := setupRouter(t)
	loginBody := map[string]string{
		"email":    "12345",
		"password": "password123",
	}
	w := doPost(router, "/login", loginBody)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

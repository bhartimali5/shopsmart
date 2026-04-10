package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"example.com/rest-api/db"
	"example.com/rest-api/routes"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// SETUP
// ============================================================

var router *gin.Engine
var validToken string

func TestMain(m *testing.M) {
	var err error
	db.DB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic("could not open test db: " + err.Error())
	}
	db.CreateTables()

	gin.SetMode(gin.TestMode)
	router = gin.Default()
	routes.RegisterRoutes(router)

	validToken, err = utils.GenerateJWT("test@example.com", "user-test-123", "user")
	if err != nil {
		panic("could not generate test token: " + err.Error())
	}

	os.Exit(m.Run())
}

func doRequest(method, path, token string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func startMockCartService(userID string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"cart": map[string]interface{}{
				"user_id":     userID,
				"total_price": 250.00,
			},
			"items": []map[string]interface{}{
				{"id": "item-001", "cart_id": "cart-001", "product_id": "prod-001", "quantity": 2, "price": 125.00},
			},
		})
	}))
}

// ============================================================
// TIMEOUT TEST
//
// Simulates cart service hanging (delayed response).
// Asserts the order service returns a meaningful error — not a raw 500.
// ============================================================

func TestCreateOrder_ShouldReturnMeaningfulError_WhenCartServiceTimesOut(t *testing.T) {
	// Arrange — cart server that hangs longer than our timeout
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second) // hangs
		w.WriteHeader(http.StatusOK)
	}))
	defer slowServer.Close()

	utils.CartServiceBaseURL = slowServer.URL
	utils.HTTPTimeout = 100 * time.Millisecond // short timeout so test doesn't wait 3s

	// Act
	w := doRequest(http.MethodPost, "/orders", validToken, nil, nil)

	// Assert — must return 400 with a human-readable message, not a raw stack trace
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err, "response must be valid JSON, not a raw stack trace")
	assert.Contains(t, body, "error", "response must contain an 'error' field")
	assert.NotContains(t, w.Body.String(), "goroutine", "must not expose Go stack trace")

	// reset timeout for other tests
	utils.HTTPTimeout = 5 * time.Second
}

// ============================================================
// IDEMPOTENCY TESTS
//
// POST the same order twice with the same Idempotency-Key.
// First call → 201 (created).
// Second call → 200 (already exists, return cached response).
// No duplicate record in DB.
// ============================================================

func TestCreateOrder_ShouldReturn201_OnFirstCallWithIdempotencyKey(t *testing.T) {
	// Arrange
	cartServer := startMockCartService("user-idempotent-1")
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL
	utils.HTTPTimeout = 5 * time.Second

	token, _ := utils.GenerateJWT("idem@example.com", "user-idempotent-1", "user")

	// Act
	w := doRequest(http.MethodPost, "/orders", token, nil, map[string]string{
		"Idempotency-Key": "unique-key-abc",
	})

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateOrder_ShouldReturn200_OnSecondCallWithSameIdempotencyKey(t *testing.T) {
	// Arrange — same key, same user as above test
	cartServer := startMockCartService("user-idempotent-2")
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	token, _ := utils.GenerateJWT("idem2@example.com", "user-idempotent-2", "user")
	headers := map[string]string{"Idempotency-Key": "unique-key-xyz"}

	// Act — first call creates the order
	first := doRequest(http.MethodPost, "/orders", token, nil, headers)
	require.Equal(t, http.StatusCreated, first.Code)

	// Act — second call with same key
	second := doRequest(http.MethodPost, "/orders", token, nil, headers)

	// Assert — returns 200, not 201
	assert.Equal(t, http.StatusOK, second.Code)
}

func TestCreateOrder_ShouldNotCreateDuplicateRecord_WhenSameIdempotencyKeyUsedTwice(t *testing.T) {
	// Arrange
	cartServer := startMockCartService("user-idempotent-3")
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	token, _ := utils.GenerateJWT("idem3@example.com", "user-idempotent-3", "user")
	headers := map[string]string{"Idempotency-Key": "unique-key-dedup"}

	// Act — call twice with same key
	doRequest(http.MethodPost, "/orders", token, nil, headers)
	doRequest(http.MethodPost, "/orders", token, nil, headers)

	// Assert — only one record exists in DB for this idempotency key
	var count int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM orders WHERE idempotency_key = ?`, "unique-key-dedup").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "duplicate order must not be created for the same idempotency key")
}

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
	// in-memory DB — no file, no leftover state
	var err error
	db.DB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic("could not open test db: " + err.Error())
	}
	db.CreateTables()

	gin.SetMode(gin.TestMode)
	router = gin.Default()
	routes.RegisterRoutes(router)

	// generate a valid JWT for a test user
	validToken, err = utils.GenerateJWT("test@example.com", "user-test-123", "user")
	if err != nil {
		panic("could not generate test token: " + err.Error())
	}

	os.Exit(m.Run())
}

// startMockCartService spins up a test HTTP server that mimics cart service
func startMockCartService(userID string, hasItems bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasItems {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"cart": map[string]interface{}{
				"user_id":     userID,
				"total_price": 250.00,
			},
			"items": []map[string]interface{}{
				{
					"id":         "item-001",
					"cart_id":    "cart-001",
					"product_id": "prod-001",
					"quantity":   2,
					"price":      125.00,
				},
			},
		})
	}))
}

func doRequest(method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ============================================================
// API TESTS — /orders
// ============================================================

func TestCreateOrder_ShouldReturn201_WhenCartHasItems(t *testing.T) {
	// Arrange
	cartServer := startMockCartService("user-test-123", true)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	// Act
	w := doRequest(http.MethodPost, "/orders", validToken, nil)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateOrder_ShouldReturnApplicationJSON_OnEveryResponse(t *testing.T) {
	// Arrange
	cartServer := startMockCartService("user-test-123", true)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	// Act
	w := doRequest(http.MethodPost, "/orders", validToken, nil)

	// Assert — Content-Type must always be application/json
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestCreateOrder_ShouldReturn400_WhenCartIsEmpty(t *testing.T) {
	// Arrange — cart service returns empty cart (missing required data)
	cartServer := startMockCartService("user-test-123", false)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	// Act
	w := doRequest(http.MethodPost, "/orders", validToken, nil)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestCreateOrder_ShouldReturn401_WhenTokenIsMissing(t *testing.T) {
	// Act — no token provided
	w := doRequest(http.MethodPost, "/orders", "", nil)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestGetOrders_ShouldReturn404_WhenUserHasNoOrders(t *testing.T) {
	// Arrange — fresh token for a user with no orders
	token, err := utils.GenerateJWT("noorders@example.com", "user-no-orders", "user")
	require.NoError(t, err)

	// Act
	w := doRequest(http.MethodGet, "/orders", token, nil)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestUpdateOrderStatus_ShouldReturn404_WhenOrderIDDoesNotExist(t *testing.T) {
	// Act
	w := doRequest(http.MethodPatch, "/orders/nonexistent-id/status", validToken,
		map[string]string{"status": "DELIVERED"})

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// ============================================================
// RESPONSE BODY SHAPE TEST
// Verifies the response body has the correct structure and fields
// ============================================================

func TestCreateOrder_ShouldReturnCorrectBodyShape_WhenOrderIsCreated(t *testing.T) {
	// Arrange
	cartServer := startMockCartService("user-test-123", true)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	// Act
	w := doRequest(http.MethodPost, "/orders", validToken, nil)
	require.Equal(t, http.StatusCreated, w.Code)

	// Assert body shape
	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Contains(t, body, "user_id", "response must contain user_id")
	assert.Contains(t, body, "order_date", "response must contain order_date")
	assert.Contains(t, body, "status", "response must contain status")
	assert.Contains(t, body, "total_price", "response must contain total_price")
}

// ============================================================
// CONSUMER CONTRACT TESTS
//
// These tests define what fields the notification_service and
// payment_service (consumers of order events) expect to always
// be present in the order response. If any field is removed or
// renamed, these tests catch the breaking change immediately.
// ============================================================

func TestOrderContract_ShouldAlwaysContainUserID(t *testing.T) {
	cartServer := startMockCartService("user-test-123", true)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	w := doRequest(http.MethodPost, "/orders", validToken, nil)
	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)

	// CONTRACT: notification_service needs user_id to fetch user email
	assert.NotEmpty(t, body["user_id"], "contract violation: user_id must always be present and non-empty")
}

func TestOrderContract_ShouldAlwaysContainOrderDate(t *testing.T) {
	cartServer := startMockCartService("user-test-123", true)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	w := doRequest(http.MethodPost, "/orders", validToken, nil)
	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)

	// CONTRACT: notification_service includes order_date in confirmation email
	assert.NotEmpty(t, body["order_date"], "contract violation: order_date must always be present and non-empty")
}

func TestOrderContract_ShouldAlwaysContainStatus(t *testing.T) {
	cartServer := startMockCartService("user-test-123", true)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	w := doRequest(http.MethodPost, "/orders", validToken, nil)
	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)

	// CONTRACT: payment_service uses status to determine if payment is needed
	assert.NotEmpty(t, body["status"], "contract violation: status must always be present and non-empty")
}

func TestOrderContract_ShouldAlwaysContainTotalPrice(t *testing.T) {
	cartServer := startMockCartService("user-test-123", true)
	defer cartServer.Close()
	utils.CartServiceBaseURL = cartServer.URL

	w := doRequest(http.MethodPost, "/orders", validToken, nil)
	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)

	// CONTRACT: payment_service uses total_price to process the payment amount
	assert.NotNil(t, body["total_price"], "contract violation: total_price must always be present")
}

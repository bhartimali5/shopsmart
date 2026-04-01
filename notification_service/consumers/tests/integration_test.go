package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/rest-api/consumers"
	"example.com/rest-api/dto"
	"example.com/rest-api/utils"
	"github.com/stretchr/testify/assert"
)

// capturedEmail records what was passed to SendEmail during the test
var capturedEmail string

// testEmailSender captures the recipient instead of logging
type testEmailSender struct{}

func (t *testEmailSender) Send(to, subject, body string) {
	capturedEmail = to
}

// startMockUserService spins up a real HTTP test server that mimics user service internal endpoint
func startMockUserService(email string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify internal key is present
		if r.Header.Get("X-Internal-Key") != os.Getenv("INTERNAL_SERVICE_KEY") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			json.NewEncoder(w).Encode(map[string]string{"email": email})
		}
	}))
}

func TestHandleOrderCreatedIntegration_ShouldSendEmail_WhenUserServiceReturnsEmail(t *testing.T) {
	os.Setenv("INTERNAL_SERVICE_KEY", "test-secret")

	server := startMockUserService("user@example.com", http.StatusOK)
	defer server.Close()

	// Point GetUserEmail at the test server
	utils.UserServiceBaseURL = server.URL

	capturedEmail = ""
	event := dto.OrderCreatedEvent{
		ID:          "order-001",
		UserID:      "user-123",
		OrderDate:   "2026-03-30",
		TotalAmount: 300.00,
		Status:      "PENDING_PAYMENT",
	}

	err := consumers.HandleOrderCreated(event, utils.NewUserEmailFetcher(), &testEmailSender{})

	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", capturedEmail)
}

func TestHandleOrderCreatedIntegration_ShouldReturnError_WhenUserServiceReturnsNotFound(t *testing.T) {
	os.Setenv("INTERNAL_SERVICE_KEY", "test-secret")

	server := startMockUserService("", http.StatusNotFound)
	defer server.Close()

	utils.UserServiceBaseURL = server.URL

	capturedEmail = ""
	event := dto.OrderCreatedEvent{UserID: "nonexistent-user"}

	err := consumers.HandleOrderCreated(event, utils.NewUserEmailFetcher(), &testEmailSender{})

	assert.Error(t, err)
	assert.Empty(t, capturedEmail)
}

func TestHandleOrderCreatedIntegration_ShouldReturnError_WhenInternalKeyIsMissing(t *testing.T) {
	os.Unsetenv("INTERNAL_SERVICE_KEY") // no key set

	server := startMockUserService("user@example.com", http.StatusOK)
	defer server.Close()

	utils.UserServiceBaseURL = server.URL

	capturedEmail = ""
	event := dto.OrderCreatedEvent{UserID: "user-123"}

	err := consumers.HandleOrderCreated(event, utils.NewUserEmailFetcher(), &testEmailSender{})

	assert.Error(t, err)
	assert.Empty(t, capturedEmail)
}

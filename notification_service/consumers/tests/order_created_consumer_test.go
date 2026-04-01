package tests

import (
	"errors"
	"testing"

	"example.com/rest-api/consumers"
	"example.com/rest-api/dto"
	"github.com/stretchr/testify/assert"
)

// --- mocks ---

type mockEmailSender struct {
	capturedTo      string
	capturedSubject string
	capturedBody    string
}

func (m *mockEmailSender) Send(to, subject, body string) {
	m.capturedTo = to
	m.capturedSubject = subject
	m.capturedBody = body
}

type mockUserEmailFetcher struct {
	email string
	err   error
}

func (m *mockUserEmailFetcher) GetEmail(_ string) (string, error) {
	return m.email, m.err
}

// --- tests ---

func TestHandleOrderCreated_ShouldSendEmail_WhenValidEventIsProvided(t *testing.T) {
	event := dto.OrderCreatedEvent{
		ID:          "order-123",
		UserID:      "user-456",
		OrderDate:   "2026-03-30",
		TotalAmount: 250.00,
		Status:      "PENDING_PAYMENT",
	}
	fetcher := &mockUserEmailFetcher{email: "user@example.com"}
	sender := &mockEmailSender{}

	err := consumers.HandleOrderCreated(event, fetcher, sender)

	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", sender.capturedTo)
	assert.Equal(t, "Order Confirmed!", sender.capturedSubject)
}

func TestHandleOrderCreated_ShouldReturnError_WhenUserEmailFetchFails(t *testing.T) {
	event := dto.OrderCreatedEvent{UserID: "user-456"}
	fetcher := &mockUserEmailFetcher{err: errors.New("user service unavailable")}
	sender := &mockEmailSender{}

	err := consumers.HandleOrderCreated(event, fetcher, sender)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not fetch user email")
	assert.Empty(t, sender.capturedTo) // email should not be sent
}

func TestHandleOrderCreated_ShouldIncludeOrderDetails_InEmailBody(t *testing.T) {
	event := dto.OrderCreatedEvent{
		ID:          "order-789",
		UserID:      "user-456",
		OrderDate:   "2026-03-30",
		TotalAmount: 499.99,
		Status:      "PENDING_PAYMENT",
	}
	fetcher := &mockUserEmailFetcher{email: "user@example.com"}
	sender := &mockEmailSender{}

	err := consumers.HandleOrderCreated(event, fetcher, sender)

	assert.NoError(t, err)
	assert.Contains(t, sender.capturedBody, "order-789")
	assert.Contains(t, sender.capturedBody, "499.99")
	assert.Contains(t, sender.capturedBody, "2026-03-30")
}

func TestHandleOrderCreated_ShouldNotSendEmail_WhenUserIDIsEmpty(t *testing.T) {
	event := dto.OrderCreatedEvent{UserID: ""}
	fetcher := &mockUserEmailFetcher{err: errors.New("user not found")}
	sender := &mockEmailSender{}

	err := consumers.HandleOrderCreated(event, fetcher, sender)

	assert.Error(t, err)
	assert.Empty(t, sender.capturedTo)
}

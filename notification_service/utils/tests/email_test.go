package tests

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"example.com/rest-api/utils"
	"github.com/stretchr/testify/assert"
)

func captureLog(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	return buf.String()
}

func TestSendEmail_ShouldLogEmailDetails_WhenCalled(t *testing.T) {
	output := captureLog(func() {
		utils.SendEmail("user@example.com", "Order Confirmed!", "Your order has been placed.")
	})

	assert.True(t, strings.Contains(output, "user@example.com"))
	assert.True(t, strings.Contains(output, "Order Confirmed!"))
}

func TestSendEmail_ShouldLogRecipient_WhenEmailIsSent(t *testing.T) {
	output := captureLog(func() {
		utils.SendEmail("another@example.com", "Subject", "Body")
	})

	assert.Contains(t, output, "another@example.com")
}

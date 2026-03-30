package utils

import "log"

// SendEmail is a mock email sender. Replace with real SMTP implementation later.
func SendEmail(to, subject, body string) {
	log.Printf("[EMAIL] To: %s | Subject: %s | Body: %s", to, subject, body)
}

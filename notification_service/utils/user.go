package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// UserServiceBaseURL can be overridden in tests to point at a mock server
var UserServiceBaseURL = "http://localhost:8001"

type userEmailResponse struct {
	Email string `json:"email"`
}

// UserEmailFetcher is the real implementation of consumers.UserEmailFetcher
type UserEmailFetcher struct{}

func NewUserEmailFetcher() *UserEmailFetcher {
	return &UserEmailFetcher{}
}

func (u *UserEmailFetcher) GetEmail(userID string) (string, error) {
	url := fmt.Sprintf("%s/internal/users/%s/email", UserServiceBaseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Internal-Key", os.Getenv("INTERNAL_SERVICE_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call user service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("user service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result userEmailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Email, nil
}

// GetUserEmail is kept for backward compatibility
func GetUserEmail(userID string) (string, error) {
	return NewUserEmailFetcher().GetEmail(userID)
}

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type userEmailResponse struct {
	Email string `json:"email"`
}

// GetUserEmail fetches user email from user_service internal endpoint.
func GetUserEmail(userID string) (string, error) {
	url := fmt.Sprintf("http://localhost:8001/internal/users/%s/email", userID)

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

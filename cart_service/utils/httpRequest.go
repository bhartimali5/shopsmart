// Code to make http request to other services
package utils

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

func MakeGetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get valid response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func MakePostRequest(url string, payload []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", io.NopCloser(bytes.NewReader(payload)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get valid response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

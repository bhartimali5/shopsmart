package utils

import (
	"errors"
	"io"
	"net/http"
)

func MakeHTTPGETRequest(url string, auth_token string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New("failed to create GET request to " + url + ": " + err.Error())
	}
	req.Header.Set("Authorization", auth_token) // forward token
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", errors.New("failed to make GET request to " + url + ": " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to make GET request to " + url + ": status " + resp.Status)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(resBody)

	return bodyString, nil
}

func MakeHTTPDELETERequest(url string, auth_token string) error {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.New("failed to create DELETE request to " + url + ": " + err.Error())
	}
	req.Header.Set("Authorization", auth_token) // forward token
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return errors.New("failed to make DELETE request to " + url + ": " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to make DELETE request to " + url + ": status " + resp.Status)
	}

	return nil
}

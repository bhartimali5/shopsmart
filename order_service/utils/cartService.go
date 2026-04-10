package utils

import (
	"encoding/json"
	"fmt"
)

// CartServiceBaseURL can be overridden in tests to point at a mock server
var CartServiceBaseURL = "http://localhost:8080"

type CartSummary struct {
	Cart struct {
		UserID     string  `json:"user_id"`
		TotalPrice float64 `json:"total_price"`
	} `json:"cart"`
	Items []struct {
		ID        string  `json:"id"`
		CartId    string  `json:"cart_id"`
		ProductID string  `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	} `json:"items"`
}

func GetCartItemDetails(userId string, auth_token string) CartSummary {
	var cartDetails CartSummary

	url := CartServiceBaseURL + "/cart/items"

	respBody, err := MakeHTTPGETRequest(url, auth_token)
	if err != nil {
		fmt.Println("Error making HTTP GET request:", err)
		return cartDetails
	}

	err = json.Unmarshal([]byte(respBody), &cartDetails)
	if err != nil {
		fmt.Println("Error unmarshalling cart details:", err)
		return cartDetails
	}
	return cartDetails
}

func ClearUserCart(userId string, auth_token string) error {
	url := CartServiceBaseURL + "/cart/clear/"
	err := MakeHTTPDELETERequest(url, auth_token)
	if err != nil {
		fmt.Println("Error making HTTP DELETE request:", err)
		return err
	}
	return nil
}

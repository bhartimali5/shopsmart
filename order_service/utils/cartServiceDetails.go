package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CartSummary struct {
	Cart struct {
		ID         string  `json:"id"`
		TotalPrice float64 `json:"total_price"`
	} `json:"cart"`
}

func GetCartItemDetails(userId string, auth_token string) CartSummary {
	// Make request to Cart Service to get active cart ID for the user
	var cartDetails CartSummary

	url := "http://localhost:8080/cart/items"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", auth_token) // forward token
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return cartDetails
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cartDetails
	}

	err = json.Unmarshal(body, &cartDetails)
	if err != nil {
		fmt.Println("Error unmarshalling cart details:", err)
		return cartDetails
	}
	fmt.Println("Cart Details:", cartDetails)
	return cartDetails
}

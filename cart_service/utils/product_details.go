// To call product service and fetch product details
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ProductDetails struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	// Add other product fields as needed
}

func FetchProductPrice(productId string) (float64, error) {
	url := fmt.Sprintf("http://localhost:8000/products/%s", productId)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to get valid response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	res := string(body)
	var s ProductDetails
	err = json.Unmarshal([]byte(res), &s)
	if err != nil {
		return 0, err
	}
	return s.Price, nil
}

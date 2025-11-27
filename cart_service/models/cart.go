package models

import (
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type Cart struct {
	ID        string `json:"id" swaggerignore:"true"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func CreateCart(cartItem Cart) (Cart, error) {
	query := `INSERT INTO carts (id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Cart{}, err
	}
	defer stmt.Close()
	// use time utils to get current timestamp
	current_time := utils.GetCurrentTimestamp()
	cartItem.CreatedAt = current_time
	cartItem.UpdatedAt = current_time
	cartItem.ID = utils.GenerateUUID()
	_, err = stmt.Exec(cartItem.ID, cartItem.UserId, cartItem.CreatedAt, cartItem.UpdatedAt)
	if err != nil {
		return Cart{}, err
	}

	return cartItem, nil

}

func GetCartByUserId(userId string) (Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = ?`
	row := db.DB.QueryRow(query, userId)
	var cartItem Cart
	err := row.Scan(&cartItem.ID, &cartItem.UserId, &cartItem.CreatedAt, &cartItem.UpdatedAt)
	if err != nil {
		return Cart{}, err
	}
	return cartItem, nil
}

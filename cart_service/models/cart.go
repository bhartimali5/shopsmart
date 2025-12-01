package models

import (
	"fmt"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type Cart struct {
	ID        string `json:"id" swaggerignore:"true"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IsActive  bool   `json:"is_active"`
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

func GetCartByUserId(userId string) ([]Cart, error) {
	query := `SELECT * FROM carts WHERE user_id = ?`
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return []Cart{}, err
	}
	var carts []Cart
	for rows.Next() {
		var cart Cart
		err := rows.Scan(&cart.ID, &cart.UserId, &cart.CreatedAt, &cart.UpdatedAt, &cart.IsActive)
		if err != nil {
			return []Cart{}, err
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func GetActiveCartByUserId(userId string) (Cart, error) {
	fmt.Println(userId)
	query := `SELECT * FROM carts WHERE user_id = ? AND is_active = 1 LIMIT 1`
	row := db.DB.QueryRow(query, userId)
	var cart Cart
	err := row.Scan(&cart.ID, &cart.UserId, &cart.CreatedAt, &cart.UpdatedAt, &cart.IsActive)
	if err != nil {
		return Cart{}, err
	}
	return cart, nil
}

func DeleteUserCart(cartId string) error {
	query := `DELETE FROM carts WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cartId)
	if err != nil {
		panic(err)
	}
	return nil

}

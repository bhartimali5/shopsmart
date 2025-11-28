package models

import (
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type CartItem struct {
	ID        string  `json:"id"`
	CartId    string  `json:"cart_id"`
	ProductId string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func AddItemToCart(item CartItem) (CartItem, error) {
	query := `INSERT INTO cart_items (id, cart_id, product_id, quantity, price) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return CartItem{}, err
	}
	defer stmt.Close()
	item.ID = utils.GenerateUUID()
	_, err = stmt.Exec(item.ID, item.CartId, item.ProductId, item.Quantity, item.Price)
	if err != nil {
		panic(err)
	}

	return item, nil
}

func RemoveItemFromCart(itemId string) error {
	query := `DELETE FROM cart_items WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(itemId)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetCartItemsByCartId(cartId string) []CartItem {
	query := `SELECT id, cart_id, product_id, quantity, price FROM cart_items WHERE cart_id = ?`
	rows, err := db.DB.Query(query, cartId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(&item.ID, &item.CartId, &item.ProductId, &item.Quantity, &item.Price)
		if err != nil {
			panic(err)
		}
		items = append(items, item)
	}
	return items
}

func GetCartItemById(itemId string) (CartItem, error) {
	query := `SELECT id, cart_id, product_id, quantity, price FROM cart_items WHERE id = ?`
	row := db.DB.QueryRow(query, itemId)
	var item CartItem
	err := row.Scan(&item.ID, &item.CartId, &item.ProductId, &item.Quantity, &item.Price)
	if err != nil {
		return CartItem{}, err
	}
	return item, nil
}

func ClearCart(cartId string) error {
	query := `DELETE FROM cart_items WHERE cart_id = ?`
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

func UpdateCartItem(itemId string, quantity int, price float64) error {
	query := `UPDATE cart_items SET quantity = ?, price = ? WHERE id = ? `
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(quantity, price, itemId)
	if err != nil {
		panic(err)
	}
	return nil
}

package models

import (
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type Order struct {
	ID          string  `json:"id"`
	CartID      string  `json:"cart_id"`
	OrderDate   string  `json:"order_date"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	UserID      string  `json:"user_id"`
}

func (o *Order) Save() error {
	query := `INSERT INTO orders (id, user_id, order_date, cart_id, status, total_amount) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	o.ID = utils.GenerateUUID()
	_, err := db.DB.Exec(query, o.ID, o.UserID, o.OrderDate, o.CartID, o.Status, o.TotalAmount)
	return err
}

func GetOrdersByUserID(userID string) ([]Order, error) {
	query := `SELECT id, user_id, order_date, cart_id, status, total_amount 
			  FROM orders WHERE user_id = ?`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.OrderDate, &order.CartID, &order.Status, &order.TotalAmount); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

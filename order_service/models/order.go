package models

type Order struct {
	ID          string  `json:"id"`
	CartID      string  `json:"cart_id"`
	OrderDate   string  `json:"order_date"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	UserID      string  `json:"user_id"`
}

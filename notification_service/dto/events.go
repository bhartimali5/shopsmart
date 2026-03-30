package dto

// OrderCreatedEvent is the payload received from order.created RabbitMQ event
type OrderCreatedEvent struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	CartID      string  `json:"cart_id"`
	OrderDate   string  `json:"order_date"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_price"`
}

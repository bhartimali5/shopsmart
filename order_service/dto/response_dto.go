package dto

type CreateOrderResponseDTO struct {
	UserID      string  `json:"user_id" binding:"required"`
	OrderDate   string  `json:"order_date" binding:"required"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_price" binding:"required"`
}

type GetOrdersResponseDTO struct {
	Orders []OrderResponseDTO `json:"orders"`
}

type OrderResponseDTO struct {
	ID          string  `json:"id"`
	CartID      string  `json:"cart_id"`
	OrderDate   string  `json:"order_date"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_price"`
	UserID      string  `json:"user_id"`
}

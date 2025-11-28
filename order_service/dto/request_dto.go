package dto

type CreateOrderRequestDTO struct {
	OrderDate string `json:"order_date" binding:"required"`
}

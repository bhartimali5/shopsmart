package dto

type CreateOrderRequestDTO struct {
	OrderDate string `json:"order_date" binding:"required"`
}

type UpdateOrderStatusRequestDTO struct {
	Status *string `json:"status" binding:"required"`
}

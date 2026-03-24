package dto

type PaymentListResponseDTO struct {
	Payments []PaymentResponseDTO `json:"payments_list"`
}

type PaymentResponseDTO struct {
	OrderId       string `json:"order_id"`
	UserId        string `json:"user_id"`
	CartId        string `json:"cart_id"`
	PaymentStatus string `json:"payment_status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

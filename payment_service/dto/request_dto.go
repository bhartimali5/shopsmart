package dto

type OrderEvent struct {
	OrderId string `json:"order_id"`
	UserId  string `json:"user_id"`
	CartId  string `json:"cart_id"`
}

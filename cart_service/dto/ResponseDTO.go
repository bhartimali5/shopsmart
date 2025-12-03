package dto

type AddItemResponse struct {
	Message string   `json:"message"`
	Item    CartItem `json:"item"`
}

type CartItem struct {
	ID        string  `json:"id"`
	ProductId string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type RemoveItemResponse struct {
	Message string `json:"message"`
	ItemId  string `json:"item_id"`
}

type ClearCartResponse struct {
	Message string `json:"message"`
}

type OrderEvent struct {
	OrderId string `json:"order_id"`
	UserId  string `json:"user_id"`
	CartId  string `json:"cart_id"`
}

package dto

type AddItemRequest struct {
	ProductId string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type RemoveItemRequest struct {
	ItemId string `json:"item_id" binding:"required"`
}

type ClearCartRequest struct {
	UserId string `json:"user_id" binding:"required"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

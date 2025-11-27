package dto

type CreateProductRequestDTO struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Price         float32 `json:"price" binding:"required"`
	CategoryID    string  `json:"category_id" binding:"required"`
	StockQuantity int     `json:"stock_quantity" binding:"required"`
}

type UpdateProductRequestDTO struct {
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Price         *float32 `json:"price,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	StockQuantity *int     `json:"stock_quantity,omitempty"`
}

type CreateCategoryRequestDTO struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryRequestDTO struct {
	Name *string `json:"name"`
}

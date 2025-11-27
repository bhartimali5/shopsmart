package dto

type GetProductsResponseDTO struct {
	Products []ProductResponseDTO `json:"products"`
}

type ProductResponseDTO struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float32 `json:"price"`
	CategoryID    string  `json:"category_id"`
	StockQuantity int     `json:"stock_quantity"`
}

type CreateProductResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
type UpdateProductResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type DeleteProductResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type CategoryResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetCategoriesResponseDTO struct {
	Categories []CategoryResponseDTO `json:"categories"`
}

type CreateCategoryResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type DeleteCategoryResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type UpdateCategoryResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

package routes

import (
	"net/http"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Summary      Get all Products
// @Description  get Products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.GetProductsResponseDTO
// @Router       /products [get]
// @Security BearerAuth
func getProducts(context *gin.Context) {
	Products, err := models.GetAllproducts()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch Products"})
		return
	}
	response := dto.GetProductsResponseDTO{}
	products := []dto.ProductResponseDTO{}
	for _, product := range Products {
		products = append(products, dto.ProductResponseDTO{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price,
			CategoryID:    product.CategoryID,
			StockQuantity: product.StockQuantity,
		})
	}
	response.Products = products
	context.JSON(http.StatusOK, response)
}

// GetProductByID godoc
// @Summary      Get a Product by ID
// @Description  get Product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  dto.ProductResponseDTO
// @Router       /products/{id} [get]
// @Security BearerAuth
func getProductByID(context *gin.Context) {

	ProductId := context.Param("id")
	product, err := models.GetProductByID(ProductId)
	if product == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product"})
		return
	}

	context.JSON(http.StatusOK, product)
}

// CreateProducts godoc
// @Summary      Create a new Product
// @Description  create Product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      dto.CreateProductRequestDTO  true  "Product to create"
// @Success      201  {object}  dto.CreateProductResponse
// @Router       /products [post]
// @Security BearerAuth
func createProducts(context *gin.Context) {
	var newProduct models.Product
	if err := context.ShouldBindJSON(&newProduct); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check if category_id exists
	_, err := models.GetCategoryByID(newProduct.CategoryID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Category doesn't exists"})
		return
	}

	err = newProduct.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := dto.CreateProductResponse{
		ID:      newProduct.ID,
		Message: "Product created successfully!",
	}
	context.JSON(http.StatusCreated, response)
}

// UpdateProduct godoc
// @Summary      Update an existing Product
// @Description  update Product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Product ID"
// @Param        product  body      dto.UpdateProductRequestDTO  true  "Product to update"
// @Router       /products/{id} [patch]
// @Security BearerAuth
func updateProduct(context *gin.Context) {
	ProductId := context.Param("id")

	var updatedProduct dto.UpdateProductRequestDTO
	if err := context.ShouldBindJSON(&updatedProduct); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := models.GetProductByID(ProductId)
	if product == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product"})
		return
	}
	product.ID = ProductId
	if updatedProduct.Name != nil {
		product.Name = *updatedProduct.Name
	}
	if updatedProduct.Description != nil {
		product.Description = *updatedProduct.Description
	}
	if updatedProduct.Price != nil {
		product.Price = *updatedProduct.Price
	}
	if updatedProduct.StockQuantity != nil {
		product.StockQuantity = *updatedProduct.StockQuantity
	}
	if updatedProduct.CategoryID != nil {
		_, err := models.GetCategoryByID(*updatedProduct.CategoryID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Category doesn't exists"})
			return
		}
		product.CategoryID = *updatedProduct.CategoryID
	}
	err = product.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Product updated succesfully!", "product": updatedProduct})
}

// DeleteProduct godoc
// @Summary      Delete a Product
// @Description  delete Product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Router       /products/{id} [delete]
// @Security BearerAuth
func deleteProduct(context *gin.Context) {
	ProductId := context.Param("id")
	product, err := models.GetProductByID(ProductId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product"})
		return
	}

	err = product.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete product"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Product deleted succesfully!"})
}

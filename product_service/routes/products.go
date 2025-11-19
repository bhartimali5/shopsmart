package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Summary      Get all Products
// @Description  get Products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Product
// @Router       /products [get]
func getProducts(context *gin.Context) {
	Products, err := models.GetAllproducts()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch Products"})
		return
	}
	context.JSON(http.StatusOK, Products)
}

// GetProductByID godoc
// @Summary      Get a Product by ID
// @Description  get Product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Router       /products/{id} [get]
func getProductByID(context *gin.Context) {

	ProductId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
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

	context.JSON(http.StatusOK, product)
}

// CreateProducts godoc
// @Summary      Create a new Product
// @Description  create Product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Product to create"
// @Router       /products [post]
func createProducts(context *gin.Context) {
	var newProduct models.Product
	if err := context.ShouldBindJSON(&newProduct); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category_id := context.GetInt64("category_id")
	newProduct.CategoryID = int(category_id)
	err := newProduct.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Product created succesfully!", "product": newProduct})
}

// UpdateProduct godoc
// @Summary      Update an existing Product
// @Description  update Product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Product ID"
// @Param        product  body      models.Product  true  "Product to update"
// @Router       /products/{id} [put]
func updateProduct(context *gin.Context) {
	ProductId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	category_id := context.GetInt64("category_id")
	product, err := models.GetProductByID(ProductId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product"})
		return
	}
	if int64(product.CategoryID) != category_id {
		context.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this product"})
		return
	}
	var updatedProduct models.Product
	if err := context.ShouldBindJSON(&updatedProduct); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedProduct.ID = int(ProductId)
	err = updatedProduct.Update()
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
// @Param        id   path      int  true  "Product ID"
// @Router       /products/{id} [delete]
func deleteProduct(context *gin.Context) {
	ProductId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	categoryID := context.GetInt64("category_id")
	product, err := models.GetProductByID(ProductId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product"})
		return
	}

	if product.CategoryID != int(categoryID) {
		context.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this product"})
		return
	}

	err = product.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete product"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Product deleted succesfully!"})
}

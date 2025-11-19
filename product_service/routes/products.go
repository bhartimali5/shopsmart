package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func getProducts(context *gin.Context) {
	Products, err := models.GetAllproducts()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch Products"})
		return
	}
	context.JSON(http.StatusOK, Products)
}

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

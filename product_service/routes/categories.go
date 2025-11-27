package routes

import (
	"net/http"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category  body      dto.CreateCategoryRequestDTO  true  "Category to create"
// @Success      201  {object}  dto.CategoryResponseDTO
// @Router       /categories [post]
// @Security BearerAuth
func createCategory(context *gin.Context) {
	var newCategory dto.CreateCategoryRequestDTO
	if err := context.ShouldBindJSON(&newCategory); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	category := models.Category{
		Name: newCategory.Name,
	}

	err := category.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.CategoryResponseDTO{
		ID:   category.ID,
		Name: category.Name,
	}

	context.JSON(http.StatusCreated, response)
}

// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.CategoryResponseDTO
// @Router       /categories [get]
func getCategories(context *gin.Context) {
	categories, err := models.GetAllCategories()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch categories"})
		return
	}

	var response []dto.CategoryResponseDTO
	for _, category := range categories {
		response = append(response, dto.CategoryResponseDTO{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	context.JSON(http.StatusOK, response)
}

// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  dto.CategoryResponseDTO
// @Router       /categories/{id} [get]
func getCategoryByID(context *gin.Context) {
	categoryId := context.Param("id")
	category, err := models.GetCategoryByID(categoryId)
	if category == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch category"})
		return
	}

	response := dto.CategoryResponseDTO{
		ID:   category.ID,
		Name: category.Name,
	}

	context.JSON(http.StatusOK, response)
}

// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Category ID"
// @Param        category  body      dto.UpdateCategoryRequestDTO  true  "Category to update"
// @Router       /categories/{id} [patch]
// @Security BearerAuth
func updateCategory(context *gin.Context) {
	categoryId := context.Param("id")

	var updateData dto.UpdateCategoryRequestDTO
	if err := context.ShouldBindJSON(&updateData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	category, err := models.GetCategoryByID(categoryId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch category"})
		return
	}
	category.Name = *updateData.Name
	err = category.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update category"})
		return
	}

	response := dto.CategoryResponseDTO{
		ID:   category.ID,
		Name: category.Name,
	}

	context.JSON(http.StatusOK, response)
}

// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Router       /categories/{id} [delete]
// @Security BearerAuth
func DeleteCategory(context *gin.Context) {
	categoryId := context.Param("id")
	category := models.Category{ID: categoryId}
	err := category.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete category"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.AuthMiddleware)

	authenticated.GET("/products", getProducts)
	authenticated.GET("/products/:id", getProductByID)
	authenticated.GET("/categories", getCategories)
	authenticated.GET("/categories/:id", getCategoryByID)

	// Protected routes

	authenticated_admin := server.Group("/")
	authenticated_admin.Use(middlewares.AuthMiddleware, middlewares.AdminOnly())
	authenticated_admin.POST("/products", createProducts)
	authenticated_admin.PATCH("/products/:id", updateProduct)
	authenticated_admin.DELETE("/products/:id", deleteProduct)
	authenticated_admin.POST("/categories", createCategory)
	authenticated_admin.DELETE("/categories/:id", DeleteCategory)
	authenticated_admin.PATCH("/categories/:id", updateCategory)

}

package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/products", getProducts)
	server.GET("/products/:id", getProductByID)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.AuthMiddleware)
	authenticated.POST("/products", createProducts)
	authenticated.PUT("/products/:id", updateProduct)
	authenticated.DELETE("/products/:id", deleteProduct)

}

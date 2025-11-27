package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	authenticated := server.Group("/")
	authenticated.Use(middlewares.AuthMiddleware)
	{
		authenticated.POST("/cart/items", AddItemToCart)
		authenticated.DELETE("/cart/items/:itemId", RemoveItemFromCart)
		authenticated.GET("/cart/items", ViewCart)
		authenticated.DELETE("/cart/clear", ClearCart)
		authenticated.PATCH("/cart/items/:itemId", UpdateCartItemQuantity)
	}

}

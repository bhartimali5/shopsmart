package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.AuthMiddleware)

	{
		authenticated.GET("/orders", GetCurrentUserOrders)
		authenticated.POST("/orders", CreateUserOrder)
		authenticated.PATCH("/orders/:order_id/status", UpdateOrderStatus)
	}
}

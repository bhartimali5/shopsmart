package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.AuthMiddleware)

	userRoutes := authenticated.Group("/users")
	{
		userRoutes.GET("/:user_id/orders", GetUserOrders)
		userRoutes.POST("/:user_id/orders", CreateUserOrder)
	}
}

package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("signup", signUp)
	server.POST("login", login)
	server.GET("/user/profile", middlewares.AuthMiddleware, getUserProfile)
}

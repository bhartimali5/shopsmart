package routes

import (
	"example.com/rest-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("signup", signUp)
	server.POST("login", login)
	server.GET("/getUserProfile", middlewares.AuthMiddleware, getUserProfile)
}

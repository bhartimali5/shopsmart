package main

import (
	"example.com/rest-api/consumers"
	"github.com/gin-gonic/gin"
)

func main() {
	go consumers.OrderCreatedConsumer()

	server := gin.Default()
	server.Run(":8004")
}

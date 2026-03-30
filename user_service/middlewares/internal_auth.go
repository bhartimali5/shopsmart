package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func InternalAuthMiddleware(c *gin.Context) {
	key := c.GetHeader("X-Internal-Key")
	if key == "" || key != os.Getenv("INTERNAL_SERVICE_KEY") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

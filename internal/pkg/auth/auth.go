package auth

import (
	"expenses_tracker/internal/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAuthMiddleware(jwtService *jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValue := c.GetHeader("Authorization")
		if headerValue == "" {
			c.Next()
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "This route requires authorization"})
			return
		}

		token := strings.TrimPrefix(headerValue, "Bearer ")

		userId, err := jwtService.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
			return
		}

		c.Set("UserId", userId)
		c.Next()
	}
}

func GetUserId(c *gin.Context) (int64, bool) {
	userId, exists := c.Get("UserId")
	if !exists {
		return 0, false
	}

	id, ok := userId.(int64)
	return id, ok
}

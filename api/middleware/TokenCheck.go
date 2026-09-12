package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

const RequiredToken = "super-big-balls-key-62169"

func CheckTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Access forbidden: Missing Authorization token.",
				Error:   "TOKEN_MISSING",
			})
			c.Abort()
			return
		}

		if token != RequiredToken {
			c.JSON(http.StatusForbidden, APIResponse{
				Success: false,
				Message: "Access forbidden: Invalid authentication token.",
				Error:   "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

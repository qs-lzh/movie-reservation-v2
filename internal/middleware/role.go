package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qs-lzh/movie-reservation/internal/dto"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := c.Get("user_role")
		if !ok {
			dto.InternalServerError(c, "Failed to get user role from claims")
			return
		}
		if userRole != "admin" {
			dto.Forbidden(c, "Not permitted to use")
			return
		}

		c.Next()
	}
}

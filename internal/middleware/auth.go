package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/qs-lzh/movie-reservation/internal/dto"
	"github.com/qs-lzh/movie-reservation/internal/security"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("jwt")
		if err != nil {
			dto.Unauthorized(c, "Failed to get jwt token from cookie")
			return
		}
		claims, err := security.VerifyToken(tokenStr)
		if err != nil {
			if errors.Is(err, security.ErrInvalidToken) {
				dto.Unauthorized(c, "Invalid token")
				return
			}
			dto.InternalServerError(c, "Failed to verify token")
			return
		}

		// Extract user_id from claims and convert from float64 to uint
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			dto.InternalServerError(c, "Invalid user_id in token")
			return
		}
		userID := uint(userIDFloat)

		c.Set("user_id", userID)
		c.Set("user_role", claims["user_role"])

		c.Next()
	}
}

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLogger(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		zapLogger.Info("http request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", cost),
			zap.Int("size", c.Writer.Size()),
		)
	}
}

func ErrorLogger(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		for _, e := range c.Errors {
			zapLogger.Error("http request error",
				zap.Error(e.Err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
		}
	}
}

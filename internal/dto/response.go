package dto

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Message string     `json:"message,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, statusCode int, data any, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Error(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, 400, "BAD_REQUEST", message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, "UNAUTHORIZED", message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, 403, "FORBIDDEN", message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, 404, "NOT_FOUND", message)
}

func Conflict(c *gin.Context, code string, message string) {
	Error(c, 409, code, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, "INTERNAL_SERVER_ERROR", message)
}

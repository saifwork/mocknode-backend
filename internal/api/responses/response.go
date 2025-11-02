package responses

import "github.com/gin-gonic/gin"

// SuccessResponse defines the standard structure for successful responses
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    any `json:"data,omitempty"`
}

// ErrorResponse defines the standard structure for failed responses
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// JSONSuccess sends a JSON success response
func JSONSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// JSONError sends a JSON error response
func JSONError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Code:    statusCode,
	})
}

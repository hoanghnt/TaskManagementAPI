package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Validation failed",
		"details": errors,
	})
}

// PaginatedResponse sends a paginated response with message
func PaginatedResponse(c *gin.Context, statusCode int, message string, data interface{}, totalItems int64, page, pageSize int) {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"pagination": gin.H{
			"total":       totalItems,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

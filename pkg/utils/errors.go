package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError sends a JSON response with the specified status code and error message.
func HandleError(context *gin.Context, err error, statusCode int) {
	context.JSON(statusCode, gin.H{"error": err.Error()})
	context.Abort()
}

func BadReqError(c *gin.Context, err error) {
	HandleError(c, err, http.StatusBadRequest)
}

func InternalError(c *gin.Context, err error) {
	// TODO: Log Error
	HandleError(c, err, http.StatusInternalServerError)
}

func UnauthorizedError(c *gin.Context, err error) {
	HandleError(c, err, http.StatusUnauthorized)
}

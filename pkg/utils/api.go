package utils

import (
	"errors"
	"goatsync/internal/codec"
	"net/http"

	"github.com/gin-gonic/gin"
)

const MSG_PACK_CONTENT_TYPE = "application/msgpack"

func getContentType(context *gin.Context) string {
	contentType := context.GetHeader("Content-Type")
	if contentType == "" {
		return "application/json"
	}

	return contentType
}

func ParseRequest(context *gin.Context, data any) error {
	contentType := getContentType(context)

	if contentType == MSG_PACK_CONTENT_TYPE {
		if err := codec.NewDecoder(context.Request.Body).Decode(data); err != nil {
			return errors.New("failed to decode msgpack")
		}
	} else {
		if err := context.ShouldBindJSON(data); err != nil {
			return err
		}
	}

	return nil
}

func SendResponse(context *gin.Context, response any) error {
	contentType := getContentType(context)

	if contentType != MSG_PACK_CONTENT_TYPE {
		// Fallback to JSON response
		context.JSON(http.StatusOK, response)
		return nil
	}

	packed, err := codec.Marshal(response)
	if err != nil {
		return errors.New("failed to encode response")
	}

	context.Data(http.StatusOK, MSG_PACK_CONTENT_TYPE, packed)
	return nil
}

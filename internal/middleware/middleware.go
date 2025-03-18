package middleware

import (
	"errors"
	"goatsync/internal/repository"
	"goatsync/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	PERMISSION_READ      = "read"
	PERMISSION_READWRITE = "readwrite"
)

func VerifyCollectionAccess(requiredAccess string) gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenStr := context.GetHeader("Authorization")
		if tokenStr == "" {
			utils.BadReqError(context, errors.New("missing token from Header"))
			return
		}

		tokens := repository.GetTokens()
		token, exists := tokens[tokenStr]
		if !exists {
			utils.UnauthorizedError(context, errors.New("invalid token"))
			return
		}

		// collectionUID := context.Param("collection_uid")
		// if collectionUID != "" {
		// 	// TODO: Check collection access level from database
		// 	// Example:
		// 	// var collection Collection
		// 	// err := db.Where("uid = ?", collectionUID).First(&collection).Error
		// 	// if err != nil || collection.Owner != username {
		// 	//     c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		// 	//     c.Abort()
		// 	//     return
		// 	// }
		// }

		context.Set("username", token.User.Username)
		context.Next()
	}
}

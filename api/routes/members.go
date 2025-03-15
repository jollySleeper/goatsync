package routes

import (
	"goatsync/internal/handlers"
	"goatsync/internal/middleware"

	"github.com/gin-gonic/gin"
)

const (
	PERMISSION_READ      = "read"
	PERMISSION_READWRITE = "readwrite"
)

// Router setup
func setupMemberRoutes(rg *gin.RouterGroup) {
	members := rg.Group("/member")
	members.Use(middleware.VerifyCollectionAccess(PERMISSION_READ))
	{
		members.GET("/", verifyCollectionAdmin(), handlers.ListMembers)
		members.DELETE("/:username", verifyCollectionAdmin(), handlers.DeleteMember)
		members.PATCH("/:username", verifyCollectionAdmin(), handlers.UpdateMemberAccess)
		members.POST("/leave", handlers.LeaveMember)
	}
}

// Check if the user is an admin of the collection
func verifyCollectionAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		// username := context.GetString("username")
		// collectionID := context.GetString("collection_id")

		// TODO: Check if user is admin
		// Example:
		// var member CollectionMember
		// err := db.Where("collection_id = ? AND username = ? AND access_level = ?",
		//     collectionID, username, AccessLevelAdmin).First(&member).Error

		context.Next()
	}
}

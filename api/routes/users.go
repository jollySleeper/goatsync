package routes

import (
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
}

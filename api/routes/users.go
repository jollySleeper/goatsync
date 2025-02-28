package routes

import (
	"github.com/gin-gonic/gin"

    "goatsync/internal/handlers"
)

func setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
    {
        user.POST("/login_challenge", handlers.LoginChallenge)
    }
}

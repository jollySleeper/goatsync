package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

    
const banner = `
--------------------------- Welcome To ------------------------------

 ██████╗  ██████╗  █████╗ ████████╗███████╗██╗   ██╗███╗   ██╗ ██████╗
██╔════╝ ██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝╚██╗ ██╔╝████╗  ██║██╔════╝
██║  ███╗██║   ██║███████║   ██║   ███████╗ ╚████╔╝ ██╔██╗ ██║██║     
██║   ██║██║   ██║██╔══██║   ██║   ╚════██║  ╚██╔╝  ██║╚██╗██║██║     
╚██████╔╝╚██████╔╝██║  ██║   ██║   ███████║   ██║   ██║ ╚████║╚██████╗
 ╚═════╝  ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═══╝ ╚═════╝
`

func main () {
    fmt.Println(banner)

	engine := gin.Default()
	engine.GET("/is_etebase", isEtebase)
	
    // Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Using default port as dint find PORT in env")
		port = "8080"
	}

	// Start the server
	log.Printf("Server starting on port %s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func isEtebase(context *gin.Context) {
	context.Status(http.StatusOK)
}

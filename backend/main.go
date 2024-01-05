package main

import (
	"fmt"
	"main/handlers"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define a route to serve static files (e.g., your React frontend)
	r.StaticFS("/static", http.Dir(filepath.Join(".", "frontend", "build", "static")))

	// Define a route for card-related operations
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World")
	})
	r.GET("/api/cards/upgrades", handlers.GetCards)

	// Run the server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	r.Run(port)
}

package main

import (
	"fmt"
	"main/handlers"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Change to frontend domain in PROD
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// Define a route to serve static files (e.g., your React frontend)
	r.StaticFS("/static", http.Dir(filepath.Join(".", "frontend", "build", "static")))

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World")
	})
	r.GET("/api/cards/upgrades", handlers.GetCards)

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	r.Run(port)
}

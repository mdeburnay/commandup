package main

import (
	"log"
	"main/handlers"
	"main/internal"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	conn := internal.ConnectToDB()
	if conn == nil {
		log.Panic("Database not connecting. Exiting.")
	}

	r := gin.Default()
	setupMiddleware(r)
	handlers.SetupRoutes(r, conn)

	port := "localhost:8080"
	r.Run(port)
}

func setupMiddleware(r *gin.Engine) {
	// Apply global middleware and configurations here
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.StaticFS("/static", http.Dir(filepath.Join(".", "frontend", "build", "static")))
	r.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hello World") })
}

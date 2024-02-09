package main

import (
	"log"

	"commandup/models"
	"commandup/routers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	models.Init()
}

func main() {
	routersInit := routers.InitRouter()
	port := "localhost:8080"

	server := &http.Server{
		Addr:    port,
		Handler: routersInit,
	}

	log.Printf("Server started at %s", port)
	server.ListenAndServe()
}

func setupMiddleware(r *gin.Engine) {
	// Apply global middleware and configurations here
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
}

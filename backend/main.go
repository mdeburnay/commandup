package main

import (
	"fmt"
	"log"
	"main/cards"
	"main/internal"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn := internal.ConnectToDB()
	if conn == nil {
		log.Panic("Database not connecting. Exiting.")
	}

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

	r.GET("cards/upgrades", cards.GetCards(conn))
	r.POST("cards/upload-card-collection", cards.UploadCardCollection(conn))

	port := "localhost:8080"
	fmt.Printf("Server is running on port %s\n", port)
	r.Run(port)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"main/handlers"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

func main() {

	conn := connectToDB()
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
	r.GET("/api/cards/upgrades", handlers.GetCards)

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	r.Run(port)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=username dbname=commandup sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if counts > 5 {
			log.Println("Database not connecting. Exiting...")
			os.Exit(1)
		}

		if err != nil {
			log.Println("Database not connecting. Retrying...")
			counts++
		} else {
			log.Println("Connected to database")
			return connection
		}

		log.Println("Backing off for " + fmt.Sprint(counts) + " seconds")
		time.Sleep(time.Duration(counts) * time.Second)
		continue
	}
}

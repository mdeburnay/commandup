package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

var db *sql.DB

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	conn := connectToDB()
	if conn == nil {
		log.Panic("Database not connecting. Exiting.")
	}

	log.Println("Connected to database")
}

func CloseDB() {
	defer db.Close()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
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

	log.Println("DSN: ", dsn)

	if dsn == "" {
		log.Panic("DSN not set. Exiting.")
	}

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

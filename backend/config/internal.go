package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var counts int64

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

func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

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

package models

import (
	"database/sql"
	"log"
)

type CardView struct {
	Name string `json:"name"`
}

type CardList struct {
	CardViews []CardView `json:"cardViews"`
	Header    string     `json:"header"`
	Tag       string     `json:"tag"`
}

type JsonDict struct {
	CardLists []CardList `json:"cardlists"`
}

type Container struct {
	JsonDict JsonDict `json:"json_dict"`
}

type ApiResponse struct {
	Header      string    `json:"header"`
	Description string    `json:"description"`
	Container   Container `json:"container"`
}

func GetUserCards() (rows *sql.Rows, err error) {
	rows, err = db.Query("SELECT name FROM cards")
	if err != nil {
		log.Printf("Database error: %v", err)
	}
	return rows, err
}

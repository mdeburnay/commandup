package models

import (
	"database/sql"
	"log"
	"strconv"
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

type Card struct {
	Count           int
	TradelistCount  int
	Name            string
	Edition         string
	Condition       string
	Language        string
	Foil            bool
	Tags            string
	LastModified    string
	CollectorNumber string
	Alter           bool
	Proxy           bool
	PurchasePrice   float64
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

func UploadUserCards(records [][]string) error {
	for _, record := range records[1:] {
		var exists bool
		foilBool := record[6] == "foil"

		err := db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM cards 
            WHERE name = $1 AND edition = $2 AND foil = $3
        )`, record[2], record[3], foilBool).Scan(&exists) // Use foilBool, which is now correctly a boolean

		if err != nil {
			log.Printf("Error checking if record exists: %v", err)
			continue
		}

		if exists {
			log.Printf("Record already exists: %v", record)
			continue
		}

		var card Card
		card.Count, _ = strconv.Atoi(record[0])
		card.TradelistCount, _ = strconv.Atoi(record[1])
		card.Name = record[2]
		card.Edition = record[3]
		card.Condition = record[4]
		card.Language = record[5]
		card.Foil = record[6] == "foil"
		card.Tags = record[7]
		card.CollectorNumber = record[9]
		card.Alter = record[10] == "Yes"
		card.Proxy = record[11] == "Yes"
		card.PurchasePrice, _ = strconv.ParseFloat(record[12], 64)

		_, err = db.Exec("INSERT INTO cards (name, edition, condition, language, foil, tags, collector_number, alter, proxy, purchase_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", card.Name, card.Edition, card.Condition, card.Language, card.Foil, card.Tags, card.CollectorNumber, card.Alter, card.Proxy, card.PurchasePrice)
		if err != nil {
			if err != nil {
				log.Printf("Error inserting record: %v", err)
				continue
			}
		}
	}

	return nil
}

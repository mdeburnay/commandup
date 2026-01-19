package models

import (
	"database/sql"
	"fmt"
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

func GetUserCards(userID int) (rows *sql.Rows, err error) {
	rows, err = db.Query(`
		SELECT c.name 
		FROM user_cards uc
		INNER JOIN cards c ON uc.card_id = c.id
		WHERE uc.user_id = $1
	`, userID)
	if err != nil {
		log.Printf("Database error: %v", err)
		return nil, err
	}
	return rows, err
}

func UploadUserCards(userID int, records [][]string) error {
	successCount := 0
	errorCount := 0

	log.Printf("Starting upload for user_id: %d, total records: %d", userID, len(records)-1)

	for i, record := range records[1:] {
		if len(record) < 13 {
			log.Printf("Skipping record %d: insufficient columns (%d)", i+1, len(record))
			errorCount++
			continue
		}

		card := parseCardRecord(record)

		// 1. Ensure the card template exists in the master list
		cardID, err := upsertCardTemplate(card)
		if err != nil {
			log.Printf("Error with card template '%s': %v", card.Name, err)
			errorCount++
			continue
		}

		// 2. Link this specific card instance to the user
		err = upsertUserOwnership(userID, cardID, card)
		if err != nil {
			log.Printf("Error with user ownership (user_id=%d, card_id=%d): %v", userID, cardID, err)
			errorCount++
			continue
		}

		successCount++
	}

	log.Printf("Upload complete: %d successful, %d errors", successCount, errorCount)
	if errorCount > 0 && successCount == 0 {
		return fmt.Errorf("all %d records failed to upload", errorCount)
	}
	return nil
}

func parseCardRecord(record []string) Card {
	count, _ := strconv.Atoi(record[0])
	tradelistCount, _ := strconv.Atoi(record[1])
	foil := record[6] == "foil"
	alter := record[10] == "Yes"
	proxy := record[11] == "Yes"
	price, _ := strconv.ParseFloat(record[12], 64)

	return Card{
		Count:           count,
		TradelistCount:  tradelistCount,
		Name:            record[2],
		Edition:         record[3],
		Condition:       record[4],
		Language:        record[5],
		Foil:            foil,
		Tags:            record[7],
		CollectorNumber: record[9],
		Alter:           alter,
		Proxy:           proxy,
		PurchasePrice:   price,
	}
}

func upsertCardTemplate(card Card) (int, error) {
	var cardID int
	err := db.QueryRow(`
		INSERT INTO cards (name, edition, collector_number, foil)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name, edition, collector_number, foil) 
		DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, card.Name, card.Edition, card.CollectorNumber, card.Foil).Scan(&cardID)
	return cardID, err
}

func upsertUserOwnership(userID int, cardID int, card Card) error {
	_, err := db.Exec(`
		INSERT INTO user_cards (
			user_id, card_id, count, tradelist_count, purchase_price, 
			condition, language, "alter", proxy, tags
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id, card_id, condition, language, "alter", proxy) 
		DO UPDATE SET 
			count = EXCLUDED.count,
			tradelist_count = EXCLUDED.tradelist_count,
			purchase_price = EXCLUDED.purchase_price,
			tags = EXCLUDED.tags,
			updated_at = CURRENT_TIMESTAMP
	`, userID, cardID, card.Count, card.TradelistCount, card.PurchasePrice,
		card.Condition, card.Language, card.Alter, card.Proxy, card.Tags)
	return err
}

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

		// check if card exists in cards table
		var cardID int
		err := db.QueryRow(`
			SELECT id FROM cards 
			WHERE name = $1 AND edition = $2 AND foil = $3
		`, card.Name, card.Edition, card.Foil).Scan(&cardID)

		if err != nil {
			if err == sql.ErrNoRows {
				// create card if it doesn't exist
				err = db.QueryRow(`
					INSERT INTO cards (name, edition, condition, language, foil, tags, collector_number, alter, proxy)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING id
				`, card.Name, card.Edition, card.Condition, card.Language, card.Foil, card.Tags, card.CollectorNumber, card.Alter, card.Proxy).Scan(&cardID)

				if err != nil {
					log.Printf("Error creating card '%s' (%s, foil=%v): %v", card.Name, card.Edition, card.Foil, err)
					errorCount++
					continue
				}
				log.Printf("Created new card: %s (id=%d)", card.Name, cardID)
			} else {
				log.Printf("Error checking if card exists '%s': %v", card.Name, err)
				errorCount++
				continue
			}
		}

		// Check if user already has this card
		var userCardExists bool
		err = db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM user_cards 
				WHERE user_id = $1 AND card_id = $2
			)
		`, userID, cardID).Scan(&userCardExists)

		if err != nil {
			log.Printf("Error checking if user has card (user_id=%d, card_id=%d): %v", userID, cardID, err)
			errorCount++
			continue
		}


		if userCardExists {
			// update existing user_card entry
			_, err = db.Exec(`
				UPDATE user_cards 
				SET count = $1, tradelist_count = $2, purchase_price = $3, updated_at = CURRENT_TIMESTAMP
				WHERE user_id = $4 AND card_id = $5
			`, card.Count, card.TradelistCount, card.PurchasePrice, userID, cardID)
			if err != nil {
				log.Printf("Error updating user_card (user_id=%d, card_id=%d): %v", userID, cardID, err)
				errorCount++
				continue
			}
			successCount++
		} else {
			// create new user_card entry
			_, err = db.Exec(`
				INSERT INTO user_cards (user_id, card_id, count, tradelist_count, purchase_price)
				VALUES ($1, $2, $3, $4, $5)
			`, userID, cardID, card.Count, card.TradelistCount, card.PurchasePrice)
			if err != nil {
				log.Printf("Error inserting user_card (user_id=%d, card_id=%d): %v", userID, cardID, err)
				errorCount++
				continue
			}
			successCount++
		}
	}

	log.Printf("Upload complete: %d successful, %d errors", successCount, errorCount)
	
	if errorCount > 0 && successCount == 0 {
		return fmt.Errorf("all %d records failed to upload", errorCount)
	}
	
	return nil
}

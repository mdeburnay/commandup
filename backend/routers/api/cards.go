package routers

import (
	"commandup/models"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
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

type Titles struct {
	CardsYouHave string `json:"cardsYouHave"`
	CardsYouNeed string `json:"cardsYouNeed"`
	CardsToCut   string `json:"cardsToCut"`
}

type CardCategory struct {
	Title string   `json:"title"`
	Cards []string `json:"cards"`
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

type CardListResponse []CardCategory

func GetCardUpgrades(c *gin.Context) {
	rows, err := models.GetUserCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch cards from database"})
		return
	}

	var userCardCollection []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows of users cards from database"})
			return
		}

		userCardCollection = append(userCardCollection, name)
	}

	apiURL := "https://json.edhrec.com/pages/precon/revenant-recon/mirko-obsessive-theorist.json"

	cardList, err := fetchApiResponse(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var userCardMap map[string]bool

	var response CardListResponse
	var cardsToCut []string   // Accumulate all cards to cut here
	var cardsYouHave []string // Accumulate all cards you have here
	var cardsYouNeed []string // Accumulate all cards you need here

	userCardMap = make(map[string]bool)
	for _, cardName := range userCardCollection {
		userCardMap[cardName] = true
	}

	for _, cardListData := range cardList.Container.JsonDict.CardLists {
		tag := cardListData.Tag
		cardViews := cardListData.CardViews

		// Process cards to cut separately to ensure they're always included
		if tag == "cardstocut" || tag == "landstocut" {
			for _, cardView := range cardViews {
				cardsToCut = append(cardsToCut, cardView.Name)
			}
			continue
		}

		for _, cardView := range cardViews {
			if _, exists := userCardMap[cardView.Name]; exists {
				cardsYouHave = append(cardsYouHave, cardView.Name)
			} else {
				cardsYouNeed = append(cardsYouNeed, cardView.Name)
			}
		}
	}

	// After processing all card lists, create the categories
	if len(cardsYouHave) > 0 {
		response = append(response, CardCategory{
			Title: "Cards You Have",
			Cards: uniqueStrings(cardsYouHave), // Ensure uniqueness
		})
	}
	if len(cardsYouNeed) > 0 {
		response = append(response, CardCategory{
			Title: "Cards You Need",
			Cards: uniqueStrings(cardsYouNeed), // Ensure uniqueness
		})
	}
	if len(cardsToCut) > 0 {
		response = append(response, CardCategory{
			Title: "Cards To Cut",
			Cards: uniqueStrings(cardsToCut), // Ensure uniqueness
		})
	}

	responseDataJSON, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", responseDataJSON)
}

func UploadCardCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Default().Println("Uploading card collection")
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}

		log.Default().Println("Saving file")

		tempFilePath := "temp_card_collection.csv"
		err = c.SaveUploadedFile(file, tempFilePath)

		log.Default().Println("File saved")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		log.Default().Println("Opening file")

		f, err := os.Open(tempFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the file"})
			return
		}
		defer f.Close()

		log.Default().Println("Reading file")

		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the CSV file"})
			return
		}

		log.Default().Println("Inserting records")

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
				continue // Skip this record due to error
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

		log.Default().Println("Removing temp file")

		os.Remove(tempFilePath)

		log.Default().Println("File removed")

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded and processed successfully"})
	}
}

func fetchApiResponse(apiURL string) (ApiResponse, error) {
	var apiResponse ApiResponse

	response, err := http.Get(apiURL)
	if err != nil {
		log.Default().Println("Error fetching API response")
		return apiResponse, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return apiResponse, fmt.Errorf("Received non-OK response: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {

		return apiResponse, err
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return apiResponse, err
	}

	return apiResponse, nil
}

func readCSVFile(filePath string) []string {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return nil
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the remaining CSV data into a slice of strings
	var cards []string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break // End of file
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			return nil
		}

		// Add the row to the slice of cards
		cards = append(cards, row...)
	}

	return cards
}

func compareCardCollections(cardViews []CardView, userCardCollection []string) (matchingCards []string, nonMatchingCards []string) {
	cardMap := make(map[string]bool)

	for _, card := range userCardCollection {
		cardMap[card] = true
	}

	for _, cardView := range cardViews {
		if _, exists := cardMap[cardView.Name]; exists {
			matchingCards = append(matchingCards, cardView.Name)
		} else {
			nonMatchingCards = append(nonMatchingCards, cardView.Name)
		}
	}

	return
}

func extractCardNames(cardViews []CardView) []string {
	var cardNames []string
	for _, card := range cardViews {
		cardNames = append(cardNames, card.Name)
	}
	return cardNames
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, value := range input {
		if _, ok := seen[value]; !ok {
			seen[value] = true
			result = append(result, value)
		}
	}

	return result
}

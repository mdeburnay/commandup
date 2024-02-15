package routers

import (
	"commandup/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func UploadCardCollection(c *gin.Context) {
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

	err = models.UploadUserCards(records)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert records into the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and processed successfully"})
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
		return apiResponse, fmt.Errorf("received non-ok response: %s", response.Status)
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

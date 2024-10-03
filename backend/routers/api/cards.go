package routers

import (
	"commandup/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type CardView struct {
	Name      string   `json:"name"`
	Synergy   *float64 `json:"synergy, omitempty"`
	Inclusion *int     `json:"inclusion, omitempty"`
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
	Header    string    `json:"header"`
	Container Container `json:"container"`
}

type Titles struct {
	CardsYouHave string `json:"cardsYouHave"`
	CardsYouNeed string `json:"cardsYouNeed"`
	CardsToCut   string `json:"cardsToCut"`
}

type CardCategory struct {
	Title string     `json:"title"`
	Cards []CardView `json:"cards"`
}

type CardListResponse []CardCategory

type CommanderPrecon struct {
	Commander string `json:"commander"`
	Precon    string `json:"precon"`
}

var userCardCollection []string

func GetCardUpgrades(c *gin.Context) {
	log.Default().Println("Getting user cards...")
	rows, err := models.GetUserCards()

	var commanderPrecon CommanderPrecon

	if err := c.ShouldBindJSON(&commanderPrecon); err != nil {
		log.Default().Println("Error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		log.Default().Print("Could not fetch cards from database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch cards from database"})
		return
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Default().Panicln("Error scanning user card rows from database")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows of users cards from database"})
			return
		}

		userCardCollection = append(userCardCollection, name)
	}

	apiUrl := generateApiUrl(&commanderPrecon.Precon, commanderPrecon.Commander)

	log.Default().Println("Search for precon", commanderPrecon.Precon, " and commander", commanderPrecon.Commander)

	cardList, err := fetchApiResponse(apiUrl)

	if err != nil {
		log.Default().Println("Search params did not match any results")
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	// Depending on whether the user specifies for a precon or not, the response will be different
	cardListResponse := formatCardListResponse(cardList, userCardCollection, &commanderPrecon.Precon)

	responseDataJSON, err := json.Marshal(cardListResponse)
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

	cardCollectionFilePath := "temp_card_collection.csv"
	err = c.SaveUploadedFile(file, cardCollectionFilePath)

	log.Default().Println("File saved")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	log.Default().Println("Opening file")

	cardCollection, err := os.Open(cardCollectionFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the file"})
		return
	}
	defer cardCollection.Close()

	log.Default().Println("Reading file")

	csvReader := csv.NewReader(cardCollection)

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

	log.Default().Println("Fetching API response")

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

func uniqueCardViews(input []CardView) []CardView {
	seen := make(map[string]bool)
	var result []CardView

	for _, card := range input {
		if _, ok := seen[card.Name]; !ok {
			seen[card.Name] = true
			result = append(result, card)
		}
	}

	return result
}

func formatString(input string) string {
	lowercaseStr := strings.ToLower(input)

	hypenStr := strings.Replace(lowercaseStr, " ", "-", -1)

	sanitiseStr := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || r == '-' {
			return r
		}
		return -1
	}, hypenStr)

	return sanitiseStr
}

func generateApiUrl(precon *string, commander string) string {
	baseUrl := "https://json.edhrec.com/pages"

	if precon != nil && *precon != "" { // Check also if precon is not an empty string
		formattedPreconName := formatString(*precon)
		formattedCommanderName := formatString(commander)
		return baseUrl + "/precon/" + formattedPreconName + "/" + formattedCommanderName + ".json"
	}

	// If no precon is provided, only format the commander name.
	formattedCommanderName := formatString(commander)
	return baseUrl + "/commanders/" + formattedCommanderName + ".json"
}

func formatCardListResponse(cardList ApiResponse, userCardCollection []string, precon *string) CardListResponse {
	if precon != nil && *precon != "" {
		return formatPreconCardListResponse(cardList, userCardCollection)
	}

	return formatPreconCardListResponse(cardList, userCardCollection)
}

func formatPreconCardListResponse(cardList ApiResponse, userCardCollection []string) CardListResponse {
	var userCardMap map[string]bool

	var response CardListResponse
	var cardsToCut []CardView   // Accumulate all cards to cut here
	var cardsYouHave []CardView // Accumulate all cards you have here
	var cardsYouNeed []CardView // Accumulate all cards you need here

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
				cardsToCut = append(cardsToCut, CardView{
					Name: cardView.Name,
				})
			}
			continue
		}

		for _, cardView := range cardViews {
			if _, exists := userCardMap[cardView.Name]; exists {
				synergyValue := math.Round(*cardView.Synergy * 100)
				inclusionValue := *cardView.Inclusion

				cardsYouHave = append(cardsYouHave, CardView{
					Name:      cardView.Name,
					Synergy:   &synergyValue,
					Inclusion: &inclusionValue,
				})
			} else {
				synergyValue := math.Round(*cardView.Synergy * 100)
				inclusionValue := *cardView.Inclusion

				cardsYouNeed = append(cardsYouNeed, CardView{
					Name:      cardView.Name,
					Synergy:   &synergyValue,
					Inclusion: &inclusionValue,
				})
			}
		}
	}

	// After processing all card lists, create the categories
	if len(cardsYouHave) > 0 {
		response = append(response, CardCategory{
			Title: "Cards You Have",
			Cards: uniqueCardViews(cardsYouHave), // Ensure uniqueness
		})
	}
	if len(cardsYouNeed) > 0 {
		response = append(response, CardCategory{
			Title: "Cards You Need",
			Cards: uniqueCardViews(cardsYouNeed), // Ensure uniqueness
		})
	}
	if len(cardsToCut) > 0 {
		response = append(response, CardCategory{
			Title: "Cards To Cut",
			Cards: uniqueCardViews(cardsToCut), // Ensure uniqueness
		})
	}

	return response
}

// func formatCommanderCardListResponse(cardList ApiResponse, userCardCollection []string) CardListResponse {
// 	log.Default().Println("Formatting commander card list response")
// 	return cardList
// }

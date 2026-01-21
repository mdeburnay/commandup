package routers

import (
	"commandup/middleware"
	"commandup/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type CardView struct {
	Name      string   `json:"name"`
	Synergy   *float64 `json:"synergy,omitempty"`
	Inclusion *int     `json:"inclusion,omitempty"`
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

type Precon struct {
	Precon string `json:"precon"`
}

func GetCardUpgrades(c *gin.Context) {
	log.Default().Println("Request to fetch card upgrades - GetCardUpgrades")

	var precon Precon

	if err := c.ShouldBindJSON(&precon); err != nil {
		log.Default().Println("Precon not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a precon name."})
		return
	}

	apiUrl := generateApiUrl(&precon.Precon)

	// fetch user cards if authenticated
	var userCardCollection []string = []string{}
	isAuthenticated := false

	if userID, ok := middleware.GetUserID(c); ok {
		isAuthenticated = true
		rows, err := models.GetUserCards(userID)

		log.Default().Println("User ID: ", userID)

		if err != nil {
			log.Printf("Could not fetch cards from database: %v", err)
		} else {
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err != nil {
					log.Printf("Error scanning user card rows: %v", err)
					continue
				}
				userCardCollection = append(userCardCollection, name)
			}
			rows.Close()
		}
	}

	cardList, err := fetchApiResponse(apiUrl)

	if err != nil {
		log.Default().Println("Precon not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Precon not found. Please check the precon name and try again."})
		return
	}

	cardListResponse := formatCardListResponse(cardList, userCardCollection, &precon.Precon, isAuthenticated)

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

	// Get user ID from context (RequireAuth middleware should have already checked this)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		log.Printf("ERROR: Upload attempted without authentication")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	log.Printf("Uploading cards for user_id: %d", userID)
	err = models.UploadUserCards(userID, records)

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

func generateApiUrl(precon *string) string {
	baseUrl := "https://json.edhrec.com/pages"

	if precon == nil || *precon == "" {
		return ""
	}

	formattedPreconName := formatString(*precon)
	return baseUrl + "/precon/" + formattedPreconName + ".json"
}

func formatCardListResponse(cardList ApiResponse, userCardCollection []string, precon *string, isAuthenticated bool) CardListResponse {
	if precon != nil && *precon != "" {
		return formatPreconCardListResponse(cardList, userCardCollection, isAuthenticated)
	}

	return formatPreconCardListResponse(cardList, userCardCollection, isAuthenticated)
}

func formatPreconCardListResponse(cardList ApiResponse, userCardCollection []string, isAuthenticated bool) CardListResponse {
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

		// Process cards to add (cardstoadd, landstoadd)
		for _, cardView := range cardViews {
			// only process cards that have inclusion data
			if cardView.Inclusion == nil {
				continue
			}

			inclusionValue := *cardView.Inclusion

			if _, exists := userCardMap[cardView.Name]; exists {
				cardsYouHave = append(cardsYouHave, CardView{
					Name:      cardView.Name,
					Inclusion: &inclusionValue,
				})
			} else {
				cardsYouNeed = append(cardsYouNeed, CardView{
					Name:      cardView.Name,
					Inclusion: &inclusionValue,
				})
			}
		}
	}

	// After processing all card lists, create the categories
	// Only show "Cards You Have" if user is authenticated
	if isAuthenticated && len(cardsYouHave) > 0 {
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

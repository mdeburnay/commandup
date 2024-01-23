package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

/********
* TYPES *
*********/

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

type CardListResponse struct {
	CardsToAdd []string `json:"cardstoadd"`
	LandsToAdd []string `json:"landstoadd"`
	CardsToCut []string `json:"cardstocut"`
	LandsToCut []string `json:"landstocut"`
}

/************
* FUNCTIONS *
*************/

func GetCards(c *gin.Context) {

	log.Default().Println("Getting cards")

	userCardCollection := readCSVFile("card_collection.csv")

	apiURL := "https://json.edhrec.com/pages/precon/eldrazi-unbound/zhulodok-void-gorger.json"

	log.Default().Println("Fetching API response")

	cardList, err := fetchApiResponse(apiURL)

	log.Default().Println("Got API response")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	log.Default().Println("Got card list")

	container := cardList.Container
	returnedCardLists := CardListResponse{}

	for _, cardListData := range container.JsonDict.CardLists {

		log.Default().Println("Processing card list")

		tag := cardListData.Tag

		cardViews := cardListData.CardViews

		matchingCards := compareCardCollections(cardViews, userCardCollection)

		log.Default().Println("Got matching cards")

		switch tag {
		case "cardstoadd":
			returnedCardLists.CardsToAdd = append(returnedCardLists.CardsToAdd, matchingCards...)
			log.Default().Println("Got cards to add")
		case "landstoadd":
			returnedCardLists.LandsToAdd = append(returnedCardLists.LandsToAdd, matchingCards...)
			log.Default().Println("Got lands to add")
		case "cardstocut":
			returnedCardLists.CardsToCut = append(returnedCardLists.CardsToCut, matchingCards...)
			log.Default().Println("Got cards to cut")
		case "landstocut":
			returnedCardLists.LandsToCut = append(returnedCardLists.LandsToCut, matchingCards...)
			log.Default().Println("Got lands to cut")
		}

		// Serialize the response data to JSON
		responseDataJSON, err := json.Marshal(returnedCardLists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		log.Default().Println("Got response data")

		// Return the JSON response to the frontend
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseDataJSON)
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
		log.Default().Println("Error reading API response")
		return apiResponse, err
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Default().Println("Error unmarshalling API response")
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

func compareCardCollections(upgradeCardList []CardView, userCardCollection []string) []string {
	var matchingCards []string

	for _, cardData := range upgradeCardList {
		cardName := cardData.Name

		for _, userCard := range userCardCollection {
			if userCard == cardName {
				matchingCards = append(matchingCards, cardName)
				break
			}
		}
	}

	return matchingCards
}

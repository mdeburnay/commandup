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

type Titles struct {
	CardsYouHave string `json:"cardsYouHave"`
	CardsYouNeed string `json:"cardsYouNeed"`
	CardsToCut   string `json:"cardsToCut"`
}

type CardCategory struct {
	Title string   `json:"title"`
	Cards []string `json:"cards"`
}

type CardListResponse []CardCategory

/************
* FUNCTIONS *
*************/

func GetCards(c *gin.Context) {
	userCardCollection := readCSVFile("card_collection.csv")

	apiURL := "https://json.edhrec.com/pages/precon/chaos-incarnate/kardur-doomscourge.json"

	cardList, err := fetchApiResponse(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var response CardListResponse
	var cardsToCut []string   // Accumulate all cards to cut here
	var cardsYouHave []string // Accumulate all cards you have here
	var cardsYouNeed []string // Accumulate all cards you need here

	for _, cardListData := range cardList.Container.JsonDict.CardLists {
		tag := cardListData.Tag
		cardViews := cardListData.CardViews

		switch tag {
		case "cardstoadd", "landstoadd":
			matchingCards, nonMatchingCards := compareCardCollections(cardViews, userCardCollection)
			cardsYouHave = append(cardsYouHave, matchingCards...)
			cardsYouNeed = append(cardsYouNeed, nonMatchingCards...)
		case "cardstocut", "landstocut":
			cardsToCut = append(cardsToCut, extractCardNames(cardViews)...)
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

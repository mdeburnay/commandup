package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type ResponseData struct {
	MatchingCards []string `json:"matchingCards"`
}

// GetCards handles the GET request to fetch cards
func GetCards(c *gin.Context) {

	userCardCollection := readCSVFile("user_card_collection.csv")

	apiURL := "https://json.edhrec.com/pages/precon/chaos-incarnate/kardur-doomscourge.json"

	// Send an HTTP GET request to the API
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return
	}
	defer response.Body.Close()

	// Check if the response status code is not 200 (OK)
	if response.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", response.Status)
		return
	}

	// Read the response body into a byte slice
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Define a struct that matches the structure of the JSON data
	var data map[string]interface{} // You can define a struct that matches your JSON data here

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON data:", err)
		return
	}

	// Access the "Cards to Add" section within the JSON data
	container := data["container"].(map[string]interface{})
	cardLists := container["json_dict"].(map[string]interface{})["cardlists"].([]interface{})

	var upgradeCards []interface{}
	// var landsToAdd []interface{}
	// var cardsToCut []interface{}
	// var landsToCut []interface{}

	for _, item := range cardLists {
		cardListData, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Error: Unable to parse card list item")
			continue
		}

		tag, ok := cardListData["tag"].(string)
		if !ok {
			fmt.Println("Error: Unable to extract 'tag' from JSON")
			continue
		}

		cardViews, ok := cardListData["cardviews"].([]interface{})
		if !ok {
			fmt.Println("Error: Unable to extract 'cardviews' from JSON")
			continue
		}

		switch tag {
		case "cardstoadd":
			upgradeCards = cardViews
		// case "landstoadd":
		// 	landsToAdd = cardViews
		// case "cardstocut":
		// 	cardsToCut = cardViews
		// case "landstocut":
		// landsToCut = cardViews
		default:
			fmt.Println("Unknown tag:", tag)
		}
	}

	// Compare JSON data with user's card collection and store results
	matchingCards := compareCardCollections(upgradeCards, userCardCollection)

	// Create a ResponseData object
	responseData := ResponseData{
		MatchingCards: matchingCards,
	}

	// Serialize the response data to JSON
	responseDataJSON, err := json.Marshal(responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the JSON response to the frontend
	c.Data(http.StatusOK, "application/json; charset=utf-8", responseDataJSON)
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

func compareCardCollections(upgradeCardList []interface{}, userCardCollection []string) []string {
	var matchingCards []string

	// Iterate through the "cardlist" array
	for _, item := range upgradeCardList {
		cardData, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Error: Unable to parse 'cardlist' item")
			continue
		}

		// Navigate the nested structure to access the card name
		cardName, ok := cardData["name"].(string)
		if !ok {
			// If the "name" field is nested within another map, you need to access it accordingly.
			nameField, nameFieldOk := cardData["name"].(map[string]interface{})
			if nameFieldOk {
				cardName, ok = nameField["name"].(string)
			}
		}

		if !ok {
			fmt.Println("Error: Unable to extract 'name' from JSON")
			continue
		}

		// Check if the cardName is in the user's collection
		for _, userCard := range userCardCollection {
			if cardName == userCard {
				matchingCards = append(matchingCards, fmt.Sprintf(`"%s"`, cardName))
				break // Break out of the inner loop once a match is found
			}
		}
	}

	return matchingCards
}

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

/********
* TYPES *
*********/

type CardView struct {
	Name string `json:"name"`
	// Add other fields as required
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
	Container Container `json:"container"`
}

type ResponseData struct {
	MatchingCards []string `json:"matchingCards"`
}

/************
* FUNCTIONS *
*************/

// GetCards handles the GET request to fetch cards
func GetCards(c *gin.Context) {

	userCardCollection := readCSVFile("card_collection.csv")

	apiURL := "https://json.edhrec.com/pages/precon/eldrazi-unbound/zhulodok-void-gorger.json"

	cardList, err := fetchApiResponse(apiURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Access the "Cards to Add" section within the JSON data
	container := cardList["container"].(map[string]interface{})
	cardLists := container["json_dict"].(map[string]interface{})["cardlists"].([]interface{})

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

		// Compare JSON data with user's card collection and store results
		matchingCards := compareCardCollections(cardViews, userCardCollection)

		// Create a return object like so:
		// 1. Cards owned
		// 2. Lands owned
		// 3. Cards to add
		// 4. Lands to add
		// 5. Cards to cut
		// 6. Lands to cut

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
}

func fetchApiResponse(apiURL string) (ApiResponse, error) {
	var apiResponse ApiResponse

	response, err := http.Get(apiURL)
	if err != nil {
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

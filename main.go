package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// Read and parse the user's CSV card collection
	userCardCollection := readCSVFile("user_card_collection.csv")

	// Define the URL for the JSON API (replace with the actual URL)
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
	cardsToAddSection := container["json_dict"].(map[string]interface{})["cardlists"].([]interface{})[0].(map[string]interface{})
	upgradeCardList := cardsToAddSection["cardviews"].([]interface{})

	// Compare JSON data with user's card collection and store results
	results := compareCardCollections(upgradeCardList, userCardCollection)

	// Store or display the results as needed
	fmt.Printf("Matching Cards: %+v\n", results)
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

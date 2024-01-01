package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	body, err := ioutil.ReadAll(response.Body)
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

	// Compare JSON data with user's card collection and store results
	results := compareCardCollections(data, userCardCollection)

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

	// Read the CSV data into a slice of strings
	cards, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return nil
	}

	// Trim leading and trailing white spaces from card names
	for i, card := range cards {
		cards[i] = strings.TrimSpace(card)
	}

	return cards
}

func compareCardCollections(jsonData map[string]interface{}, userCardCollection []string) []string {
	var matchingCards []string

	// Iterate through the JSON data and check if each card is in the user's collection
	for _, card := range jsonData {
		cardName := card.(string) // Assuming the card name is stored as a string in the JSON data

		// Trim leading and trailing white spaces from the card name
		cardName = strings.TrimSpace(cardName)

		// Check if the cardName is in the user's collection
		for _, userCard := range userCardCollection {
			if cardName == userCard {
				matchingCards = append(matchingCards, cardName)
				break // Break out of the inner loop once a match is found
			}
		}
	}

	return matchingCards
}

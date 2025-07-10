package main

import (
	"fmt"
	"log"
	"os"

	"crowdstrike-data-collector/api" // Import our custom api package

	"github.com/joho/godotenv" // Import godotenv for .env file handling
)

func main() {
	// Load environment variables from .env file
	// godotenv.Load() looks for a .env file in the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the API URL from environment variables
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		log.Fatalf("API_URL environment variable not set in .env file")
	}

	fmt.Println("Fetching data from:", apiURL)

	// Call the function from the 'api' package to fetch products
	responseBody, err := api.FetchRandomProducts(apiURL)
	if err != nil {
		log.Fatalf("Error fetching products: %v", err)
	}

	fmt.Println("\nAPI Response:")
	fmt.Println(string(responseBody))
}

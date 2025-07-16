package main

import (
	"fmt"
	"log"
	"time"

	rtr "crowdstrike-data-collector/api" // Import the rtr package

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create a new CrowdStrikeRTRClient instance
	rtrClient, err := rtr.NewCrowdStrikeRTRClient()
	if err != nil {
		log.Fatalf("Configuration Error: %v", err)
	}

	// 1. Get Authentication Token
	fmt.Println("--- Step 1: Getting Authentication Token ---")
	if !rtrClient.GetAuthToken() {
		log.Fatal("Failed to get authentication token. Exiting.")
	}
	fmt.Println("Authentication token obtained successfully.")

	// 2. Initialize RTR Session
	fmt.Println("\n--- Step 2: Initializing RTR Session ---")
	if !rtrClient.InitializeRTRSession() {
		log.Fatal("Failed to initialize RTR session. Exiting.")
	}
	fmt.Printf("RTR Session ID: %s\n", rtrClient.SessionID)

	// 3. Run the RTR Script
	// Replace "test-omkar.ps1" with the actual name of your cloud-stored script if different.
	fmt.Println("\n--- Step 3: Running RTR Script ---")
	if !rtrClient.RunRTRScript("test-omkar.ps1") {
		log.Fatal("Failed to run RTR script. Exiting.")
	}
	fmt.Printf("Cloud Request ID for command: %s\n", rtrClient.CloudRequestID)

	// Give some time for the command to execute and status to update
	fmt.Println("\nWaiting 5 seconds for command execution...")
	time.Sleep(5 * time.Second)

	// 4. Get Status of the executed RTR command
	fmt.Println("\n--- Step 4: Getting RTR Command Status ---")
	status, err := rtrClient.GetRTRCommandStatus()
	if err != nil {
		log.Fatalf("Failed to get command status: %v", err)
	}
	if status != nil {
		fmt.Println("RTR Command Status retrieved successfully.")
	} else {
		fmt.Println("RTR Command Status could not be retrieved.")
	}

	fmt.Println("\n--- Application Finished ---")
}

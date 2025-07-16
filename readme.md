# **CrowdStrike RTR Go Client**

This project provides a Go application to interact with the CrowdStrike Real-time Response (RTR) API. It demonstrates how to authenticate with the CrowdStrike API, initialize an RTR session, run a cloud-stored script on an endpoint, and retrieve the command's status.

## **Table of Contents**

- Features
- Prerequisites
- Project Structure
- Setup
  - Environment Variables (.env)
- Installation
- Usage
- Error Handling
- Important Notes

## **Features**

- **Authentication:** Obtains an OAuth2 access token from the CrowdStrike API.
- **RTR Session Management:** Initializes a Real-time Response session with a specified endpoint.
- **Script Execution:** Runs a cloud-stored RTR script on the active session.
- **Command Status Retrieval:** Fetches the status and output of the executed RTR command.
- **Modular Design:** Separates API client logic, models, and main application flow into distinct packages.
- **Environment Variable Loading:** Uses .env files for secure credential management.

## **Prerequisites**

Before running this application, ensure you have:

- **Go (Golang) installed:** Version 1.18 or higher is recommended. You can check your Go version using go version. If not installed, please refer to the [official Go installation guide](https://go.dev/doc/install).
- **CrowdStrike Falcon API Credentials:**
  - CLIENT_ID
  - CLIENT_SECRET
  - These can be generated in the CrowdStrike Falcon console under **Support & Resources > API Clients and Keys**. Ensure the API client has the necessary permissions for Real-time Response (e.g., Real-time Response -> Read and Write).
  - Click [here](https://falcon.crowdstrike.com/documentation/page/a2a7fc0e/crowdstrike-oauth2-based-apis) to check official documentation for more information.
- **CrowdStrike Device ID (AID):** The Host ID of the endpoint you wish to target with RTR commands. This can be found in the Falcon console under **Host management**.

## **Project Structure**

The project is organized as follows:

```
crowdstrike-data-collector/
├── .env # Environment variables (CLIENT_ID, CLIENT_SECRET, DEVICE_ID)
├── .gitignore # Specifies files/directories to ignore in Git
├── go.mod # Defines the module path and direct dependencies
├── go.sum # Stores cryptographic checksums for module dependencies
├── main.go # Main application entry point
└── api/ # Package for CrowdStrike RTR client logic
├── api.go # Implements the CrowdStrikeRTRClient and API interaction methods (Manager Class)
```

## **Setup**

1. Clone the repository (or create the files manually):
   If you're starting from scratch, create the crowdstrike-data-collector directory and the sub-directory rtr.
2. Initialize Go Module:
   Navigate to the root of your project (crowdstrike-data-collector) in your terminal and initialize the Go module:
   go mod init crowdstrike-data-collector
   <br/>(You can replace crowdstrike-data-collector with your desired module path, e.g., github.com/yourusername/yourproject).

### **Environment Variables (.env)**

Create a file named .env in the root of your crowdstrike-data-collector directory. Populate it with your CrowdStrike API credentials and the target device ID:

CLIENT_ID="YOUR_CROWDSTRIKE_CLIENT_ID"
CLIENT_SECRET="YOUR_CROWDSTRIKE_CLIENT_SECRET"
DEVICE_ID="YOUR_CROWDSTRIKE_DEVICE_ID"

**Replace the placeholder values with your actual credentials and device ID.**

## **Installation**

After setting up the .env file and project structure, you need to download the Go dependencies. From the project root, run:

go mod tidy

This command will download the github.com/joho/godotenv package and update your go.mod and go.sum files.

## **Usage**

To run the application, navigate to the root of your crowdstrike-data-collector directory and execute:

go run .

The application will perform the following steps:

1. **Get Authentication Token:** Attempts to obtain an OAuth2 access token.
2. **Initialize RTR Session:** Attempts to establish an RTR session with the DEVICE_ID specified in your .env file.
3. **Run RTR Script:** Attempts to execute the test-omkar.ps1 (or your specified script name) on the active RTR session.
4. **Get RTR Command Status:** Waits for 5 seconds, then retrieves and prints the status of the executed command.

You will see output in your console detailing each step, including API responses.

## **Error Handling**

The application includes robust error handling for API calls, network issues, and JSON parsing. Any critical errors will cause the program to exit with a descriptive message. Warnings are printed if DEVICE_ID is not found in the .env file.

## **Important Notes**

- **API Permissions:** Ensure your CrowdStrike API client has the necessary Real-time Response permissions (both Read and Write) to perform all actions.
- **Device Online Status:** The target DEVICE_ID must be online and reachable for RTR sessions to be successfully initialized and commands to be executed.
- **Script Name:** The run_rtr_script function currently attempts to run a script named "test-omkar.ps1". This script must exist as a cloud-stored script in your CrowdStrike Falcon environment. Adjust the script_name argument in main.go if you need to run a different script.
- **go.mod and go.sum:**
  - go.mod defines your module and its direct dependencies. It's the primary configuration for Go's module system.
  - go.sum contains cryptographic checksums of all direct and indirect dependencies. It ensures the integrity and authenticity of downloaded modules, preventing tampering. Both files are critical for reproducible builds and should always be committed to version control.

package rtr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// CrowdStrikeRTRClient holds the necessary credentials, API endpoints,
// and session information for interacting with the CrowdStrike RTR API.
type CrowdStrikeRTRClient struct {
	ClientID          string
	ClientSecret      string
	BaseURL           string
	AuthTokenURL      string
	RTRSessionURL     string
	RTRAdminCommandURL string

	AccessToken   string
	DeviceID      string
	SessionID     string
	CloudRequestID string

	HTTPClient *http.Client // Reusable HTTP client
}

// NewCrowdStrikeRTRClient initializes and returns a new CrowdStrikeRTRClient.
// It loads credentials from environment variables and sets up API endpoints.
func NewCrowdStrikeRTRClient() (*CrowdStrikeRTRClient, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	deviceID := os.Getenv("DEVICE_ID")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("CLIENT_ID and CLIENT_SECRET must be set in the .env file")
	}
	if deviceID == "" {
		fmt.Println("Warning: DEVICE_ID not found in .env. Please set it or provide it programmatically.")
	}

	baseURL := "https://api.crowdstrike.com"
	return &CrowdStrikeRTRClient{
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		DeviceID:          deviceID,
		BaseURL:           baseURL,
		AuthTokenURL:      fmt.Sprintf("%s/oauth2/token", baseURL),
		RTRSessionURL:     fmt.Sprintf("%s/real-time-response/entities/sessions/v1", baseURL),
		RTRAdminCommandURL: fmt.Sprintf("%s/real-time-response/entities/admin-command/v1", baseURL),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second, // Set a default timeout for HTTP requests
		},
	}, nil
}

// getHeaders constructs HTTP headers based on content type and authentication status.
func (c *CrowdStrikeRTRClient) getHeaders(contentType string, includeAuth bool) map[string]string {
	headers := map[string]string{
		"accept": "application/json",
		"Content-Type": contentType,
	}
	if includeAuth && c.AccessToken != "" {
		headers["authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)
	}
	return headers
}

// makeAPICall is a generic helper to perform HTTP requests and handle responses.
func (c *CrowdStrikeRTRClient) makeAPICall(
	method string,
	url string,
	headers map[string]string,
	params map[string]string,
	jsonPayload interface{}, // Use interface{} for generic JSON payload
	formData url.Values,    // Use url.Values for form data
) (map[string]interface{}, error) { // Return map[string]interface{} for generic JSON response
	var reqBody []byte
	var err error

	if jsonPayload != nil {
		reqBody, err = json.Marshal(jsonPayload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
		}
	} else if formData != nil {
		reqBody = []byte(formData.Encode())
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add query parameters
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w. Response: %s", err, string(bodyBytes))
	}

	return result, nil
}

// GetAuthToken obtains an authentication token from the CrowdStrike API.
func (c *CrowdStrikeRTRClient) GetAuthToken() bool {
	headers := c.getHeaders("application/x-www-form-urlencoded", false)
	formData := url.Values{}
	formData.Set("client_id", c.ClientID)
	formData.Set("client_secret", c.ClientSecret)

	tokenInfo, err := c.makeAPICall("POST", c.AuthTokenURL, headers, nil, nil, formData)
	if err != nil {
		fmt.Printf("Failed to get authentication token: %v\n", err)
		return false
	}

	if accessToken, ok := tokenInfo["access_token"].(string); ok {
		c.AccessToken = accessToken
		return true
	}

	fmt.Println("Failed to get access token from response.")
	return false
}

// InitializeRTRSession initializes a new Real-time Response session.
func (c *CrowdStrikeRTRClient) InitializeRTRSession() bool {
	if c.DeviceID == "" {
		fmt.Println("Device ID not provided. Cannot initialize RTR session.")
		return false
	}

	headers := c.getHeaders("application/json", true)
	params := map[string]string{"timeout": "30", "timeout_duration": "30s"}
	payload := map[string]interface{}{"device_id": c.DeviceID, "queue_offline": false}

	fmt.Printf("Attempting to initialize RTR session for device: %s...\n", c.DeviceID)
	sessionInfo, err := c.makeAPICall("POST", c.RTRSessionURL, headers, params, payload, nil)
	if err != nil {
		fmt.Printf("Failed to initialize RTR session: %v\n", err)
		return false
	}

	// Assuming the response structure is `{"resources": [{"session_id": "..."}]}`
	if resources, ok := sessionInfo["resources"].([]interface{}); ok && len(resources) > 0 {
		if resourceMap, ok := resources[0].(map[string]interface{}); ok {
			if sessionID, ok := resourceMap["session_id"].(string); ok {
				c.SessionID = sessionID
				return true
			}
		}
	}
	fmt.Println("Failed to get session_id from RTR session initialization response.")
	return false
}

// RunRTRScript runs an RTR script on a host.
func (c *CrowdStrikeRTRClient) RunRTRScript(scriptName string) bool {
	if c.DeviceID == "" || c.SessionID == "" {
		fmt.Println("Device ID or Session ID not available. Cannot run RTR script.")
		return false
	}

	headers := c.getHeaders("application/json", true)
	payload := map[string]interface{}{
		"base_command":   "runscript",
		"command_string": fmt.Sprintf(`runscript -CloudFile="%s"`, scriptName),
		"device_id":      c.DeviceID,
		"id":             0, // This ID might be an internal counter, often 0 for new commands
		"persist":        true,
		"session_id":     c.SessionID,
	}

	fmt.Printf("Attempting to run RTR script '%s' for session: %s on device: %s...\n",
		scriptName, c.SessionID, c.DeviceID)
	commandResponse, err := c.makeAPICall("POST", c.RTRAdminCommandURL, headers, nil, payload, nil)
	if err != nil {
		fmt.Printf("Failed to run RTR script: %v\n", err)
		return false
	}

	// Assuming the response structure is `{"resources": [{"cloud_request_id": "..."}]}`
	if resources, ok := commandResponse["resources"].([]interface{}); ok && len(resources) > 0 {
		if resourceMap, ok := resources[0].(map[string]interface{}); ok {
			if cloudRequestID, ok := resourceMap["cloud_request_id"].(string); ok {
				c.CloudRequestID = cloudRequestID
				return true
			}
		}
	}
	fmt.Println("Failed to get cloud_request_id from run script response.")
	return false
}

// GetRTRCommandStatus gets the status of a single executed RTR administrator command.
func (c *CrowdStrikeRTRClient) GetRTRCommandStatus() (map[string]interface{}, error) {
	if c.CloudRequestID == "" {
		return nil, fmt.Errorf("Cloud Request ID not available. Cannot get command status.")
	}

	headers := c.getHeaders("application/json", true)
	params := map[string]string{
		"cloud_request_id": c.CloudRequestID,
		"sequence_id":      "0", // Typically 0 for the initial command status
	}

	fmt.Printf("Attempting to get status for command with Cloud Request ID: %s...\n", c.CloudRequestID)
	statusResponse, err := c.makeAPICall("GET", c.RTRAdminCommandURL, headers, params, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get RTR command status: %w", err)
	}

	// You can add more specific parsing here if you want to extract command output, errors, etc.
	fmt.Println("RTR Command Status Response (Raw):")
	prettyJSON, _ := json.MarshalIndent(statusResponse, "", "  ")
	fmt.Println(string(prettyJSON))

	return statusResponse, nil
}

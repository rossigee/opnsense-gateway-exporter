package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GatewayStatus struct {
	Name            string `json:"name"`
	Address         string `json:"address"`
	Status          string `json:"status"`
	Loss            string `json:"loss"`
	Delay           string `json:"delay"`
	StandardDev     string `json:"stddev"`
	StatusTranslated string `json:"status_translated"`
}

type GatewayStatusResponse struct {
	Items  []GatewayStatus `json:"items"`
	Status string          `json:"status"`
}

func main() {
	err := getGatewayStatus()
	if err != nil {
		log.Fatal(err)
	}
}

func getGatewayStatus() error {
	apiURL := os.Getenv("API_URL") + "/api/routes/gateway/status"

	// Create an HTTP client
	client := &http.Client{}

	// Create a GET request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Fetch the API_KEY and API_SECRET from environment variables
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	// Set Basic Authentication using the API_KEY as username and API_SECRET as password
	req.SetBasicAuth(apiKey, apiSecret)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	// Print the response body
	fmt.Println(string(body))

	// Parse the JSON response
	var gatewayStatusResponse GatewayStatusResponse
	err = json.Unmarshal(body, &gatewayStatusResponse)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Print the gateway status
	for _, gateway := range gatewayStatusResponse.Items {
		fmt.Println("Gateway:", gateway.Name)
		fmt.Println("Address:", gateway.Address)
		fmt.Println("Status:", gateway.Status)
		fmt.Println("Loss:", gateway.Loss)
		fmt.Println("Delay:", gateway.Delay)
		fmt.Println("Standard Deviation:", gateway.StandardDev)
		fmt.Println("Status Translated:", gateway.StatusTranslated)
		fmt.Println()
	}

	return nil
}

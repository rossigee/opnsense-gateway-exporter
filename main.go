package main

import (
	"os"
)

func main() {
	apiUrl := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	go func() {
		start_prometheus(apiUrl, apiKey, apiSecret)
	}()

	// Wait for end signal
	<-make(chan struct{})
}

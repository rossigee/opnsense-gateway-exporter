package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type GatewayStatus struct {
	Name             string  `json:"name"`
	Address          string  `json:"address"`
	Loss             string  `json:"loss"`
	LossValue        float64 // New attribute to hold parsed loss value
	Delay            string  `json:"delay"`
	DelayValue       float64 // New attribute to hold parsed delay value
	StandardDev      string  `json:"stddev"`
	StandardDevValue float64 // New attribute to hold parsed standard deviation value
	StatusTranslated string  `json:"status_translated"`
	StatusValue      float64 // New attribute to hold parsed status value
}

type GatewayStatusResponse struct {
	Items  []GatewayStatus `json:"items"`
	Status string          `json:"status"`
}

type opnSenseExporter struct {
	apiURL    string
	apiKey    string
	apiSecret string
}

func (e *opnSenseExporter) getGatewayStatus() (*GatewayStatusResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", e.apiURL+"/api/routes/gateway/status", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.SetBasicAuth(e.apiKey, e.apiSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var gatewayStatusResponse GatewayStatusResponse
	err = json.Unmarshal(body, &gatewayStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Parse provided values
	for i := range gatewayStatusResponse.Items {
		loss, err := strconv.ParseFloat(strings.TrimSuffix(gatewayStatusResponse.Items[i].Loss, " %"), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing loss: %v", err)
		}
		gatewayStatusResponse.Items[i].LossValue = loss

		delay, err := strconv.ParseFloat(strings.TrimSuffix(gatewayStatusResponse.Items[i].Delay, " ms"), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing delay: %v", err)
		}
		gatewayStatusResponse.Items[i].DelayValue = delay
		
		stddev, err := strconv.ParseFloat(strings.TrimSuffix(gatewayStatusResponse.Items[i].StandardDev, " ms"), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing stddev: %v", err)
		}
		gatewayStatusResponse.Items[i].StandardDevValue = stddev

		online := gatewayStatusResponse.Items[i].StatusTranslated == "Online"
		if online {
			gatewayStatusResponse.Items[i].StatusValue = 1.0
		} else {
			gatewayStatusResponse.Items[i].StatusValue = 0.0
		}
	}

	return &gatewayStatusResponse, nil
}

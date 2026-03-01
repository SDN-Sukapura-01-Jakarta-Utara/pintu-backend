package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Grafana helper functions
func createDashboardRequest(dashboardData Dashboard) ([]byte, error) {
	var bodyReader io.Reader
	jsonBody, err := json.Marshal(dashboardData)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", "http://localhost:3011/api/dashboards/db", bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin01")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

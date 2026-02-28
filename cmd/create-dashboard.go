package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	grafanaURL   = "http://localhost:3011"
	grafanaUser  = "admin"
	grafanaPass  = "admin01"
)

type Dashboard struct {
	Dashboard DashboardBody `json:"dashboard"`
	Overwrite bool          `json:"overwrite"`
}

type DashboardBody struct {
	Title    string    `json:"title"`
	Tags     []string  `json:"tags"`
	Panels   []Panel   `json:"panels"`
	Refresh  string    `json:"refresh"`
	Time     TimeRange `json:"time"`
	Timezone string    `json:"timezone"`
}

type Panel struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	GridPos   GridPos   `json:"gridPos"`
	Targets   []Target  `json:"targets"`
	Datasource string   `json:"datasource"`
}

type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type Target struct {
	RefID        string `json:"refId"`
	Expr         string `json:"expr"`
	Interval     string `json:"interval,omitempty"`
	LegendFormat string `json:"legendFormat,omitempty"`
}

type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func makeRequest(method, url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(grafanaUser, grafanaPass)

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

func main() {
	fmt.Println("Creating Grafana Dashboard...")
	fmt.Println("============================\n")

	dashboard := Dashboard{
		Overwrite: true,
		Dashboard: DashboardBody{
			Title:    "PINTU Backend Monitoring",
			Tags:     []string{"prometheus", "pintu-backend"},
			Refresh:  "30s",
			Timezone: "browser",
			Time: TimeRange{
				From: "now-1h",
				To:   "now",
			},
			Panels: []Panel{
				// Panel 1: Request Rate
				{
					ID:         1,
					Title:      "Request Rate (req/sec)",
					Type:       "timeseries",
					Datasource: "Prometheus",
					GridPos:    GridPos{X: 0, Y: 0, W: 12, H: 8},
					Targets: []Target{
						{
							RefID:        "A",
							Expr:         "rate(http_requests_total[1m])",
							Interval:     "30s",
							LegendFormat: "{{method}} {{endpoint}}",
						},
					},
				},
				// Panel 2: Response Time (95th percentile)
				{
					ID:         2,
					Title:      "Response Time p95 (seconds)",
					Type:       "timeseries",
					Datasource: "Prometheus",
					GridPos:    GridPos{X: 12, Y: 0, W: 12, H: 8},
					Targets: []Target{
						{
							RefID:        "A",
							Expr:         "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
							LegendFormat: "p95 {{method}} {{endpoint}}",
						},
					},
				},
				// Panel 3: Error Rate
				{
					ID:         3,
					Title:      "Error Rate (errors/sec)",
					Type:       "timeseries",
					Datasource: "Prometheus",
					GridPos:    GridPos{X: 0, Y: 8, W: 12, H: 8},
					Targets: []Target{
						{
							RefID:        "A",
							Expr:         "rate(http_errors_total[1m])",
							LegendFormat: "{{method}} {{endpoint}} {{status}}",
						},
					},
				},
				// Panel 4: Total Requests
				{
					ID:         4,
					Title:      "Total Requests (5min)",
					Type:       "stat",
					Datasource: "Prometheus",
					GridPos:    GridPos{X: 12, Y: 8, W: 12, H: 8},
					Targets: []Target{
						{
							RefID: "A",
							Expr:  "sum(increase(http_requests_total[5m]))",
						},
					},
				},
			},
		},
	}

	_, err := makeRequest("POST", grafanaURL+"/api/dashboards/db", dashboard)
	if err != nil {
		fmt.Printf("✗ Failed: %v\n", err)
		return
	}

	fmt.Println("✓ Dashboard created successfully!")
	fmt.Printf("\nOpen: %s/dashboards\n", grafanaURL)
	fmt.Println("Look for: PINTU Backend Monitoring")
}

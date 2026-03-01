package main

// Grafana types and helpers for setup
type DataSource struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	Access    string `json:"access"`
	IsDefault bool   `json:"isDefault"`
}

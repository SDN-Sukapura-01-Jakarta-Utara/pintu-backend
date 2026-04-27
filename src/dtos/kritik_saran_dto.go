package dtos

import "time"

// KritikSaranCreateRequest represents the request payload for creating KritikSaran
type KritikSaranCreateRequest struct {
	Nama        string `json:"nama" binding:"required"`
	KritikSaran string `json:"kritik_saran" binding:"required"`
}

// KritikSaranGetAllRequest represents the request payload for getting all kritik saran with filters
type KritikSaranGetAllRequest struct {
	Search struct {
		StartDate string `json:"start_date"` // Format: YYYY-MM-DD
		EndDate   string `json:"end_date"`   // Format: YYYY-MM-DD
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// KritikSaranResponse represents the response payload for KritikSaran
type KritikSaranResponse struct {
	ID          uint      `json:"id"`
	Nama        string    `json:"nama"`
	KritikSaran string    `json:"kritik_saran"`
	CreatedAt   time.Time `json:"created_at"`
}

// KritikSaranListResponse represents the response payload for listing KritikSaran
type KritikSaranListResponse struct {
	Data  []KritikSaranResponse `json:"data"`
	Total int64                 `json:"total"`
}

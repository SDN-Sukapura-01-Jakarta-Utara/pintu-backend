package dtos

import "time"

// JumbotronCreateRequest represents the request payload for creating Jumbotron
type JumbotronCreateRequest struct {
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// JumbotronUpdateRequest represents the request payload for updating Jumbotron
type JumbotronUpdateRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// JumbotronResponse represents the response payload for Jumbotron
type JumbotronResponse struct {
	ID          uint      `json:"id"`
	File        string    `json:"file"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// JumbotronListResponse represents the response payload for listing Jumbotron
type JumbotronListResponse struct {
	Data   []JumbotronResponse `json:"data"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
	Total  int64               `json:"total"`
}

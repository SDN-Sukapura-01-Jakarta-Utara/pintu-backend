package dtos

import "time"

// BidangStudiCreateRequest represents the request payload for creating BidangStudi
type BidangStudiCreateRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=100"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// BidangStudiUpdateRequest represents the request payload for updating BidangStudi
type BidangStudiUpdateRequest struct {
	ID     uint    `json:"id" binding:"required"`
	Name   *string `json:"name" binding:"omitempty,min=2,max=100"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// BidangStudiResponse represents the response payload for BidangStudi
type BidangStudiResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// BidangStudiListResponse represents the response payload for listing BidangStudi
type BidangStudiListResponse struct {
	Data []BidangStudiResponse `json:"data"`
}

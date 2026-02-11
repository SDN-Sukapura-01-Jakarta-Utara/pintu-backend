package dtos

import "time"

// SaranaPrasaranaCreateRequest represents the request payload for creating SaranaPrasarana
type SaranaPrasaranaCreateRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=255"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// SaranaPrasaranaUpdateRequest represents the request payload for updating SaranaPrasarana
type SaranaPrasaranaUpdateRequest struct {
	ID     uint    `json:"id" binding:"required"`
	Name   *string `json:"name" binding:"omitempty,min=2,max=255"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// SaranaPrasaranaResponse represents the response payload for SaranaPrasarana
type SaranaPrasaranaResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Foto        string    `json:"foto"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID *uint     `json:"created_by_id"`
	UpdatedByID *uint     `json:"updated_by_id"`
}

// SaranaPrasaranaListResponse represents the response payload for listing SaranaPrasarana
type SaranaPrasaranaListResponse struct {
	Data   []SaranaPrasaranaResponse `json:"data"`
	Limit  int                       `json:"limit"`
	Offset int                       `json:"offset"`
	Total  int64                     `json:"total"`
}

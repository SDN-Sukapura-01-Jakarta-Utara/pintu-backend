package dtos

import "time"

// ApplicationCreateRequest represents the request payload for creating Application
type ApplicationCreateRequest struct {
	Nama             string `json:"nama" binding:"required"`
	Link             string `json:"link" binding:"required"`
	ShowInJumbotron  bool   `json:"show_in_jumbotron"`
	Status           string `json:"status"`
}

// ApplicationUpdateRequest represents the request payload for updating Application
type ApplicationUpdateRequest struct {
	ID               uint    `json:"id" binding:"required"`
	Nama             *string `json:"nama"`
	Link             *string `json:"link"`
	ShowInJumbotron  *bool   `json:"show_in_jumbotron"`
	Status           *string `json:"status"`
}

// ApplicationGetAllRequest represents the request payload for getting all applications
type ApplicationGetAllRequest struct {
	Search struct {
		Nama             string `json:"nama"`
		Status           string `json:"status"`
		ShowInJumbotron  *bool  `json:"show_in_jumbotron"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// ApplicationResponse represents the response payload for Application
type ApplicationResponse struct {
	ID               uint       `json:"id"`
	Nama             string     `json:"nama"`
	Link             string     `json:"link"`
	ShowInJumbotron  bool       `json:"show_in_jumbotron"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedByID      *uint      `json:"created_by_id,omitempty"`
	UpdatedByID      *uint      `json:"updated_by_id,omitempty"`
}

// ApplicationListResponse represents the response payload for listing Application
type ApplicationListResponse struct {
	Data []ApplicationResponse `json:"data"`
}

// ApplicationListWithPaginationResponse represents the response with pagination
type ApplicationListWithPaginationResponse struct {
	Data       []ApplicationResponse `json:"data"`
	Pagination PaginationInfo        `json:"pagination"`
}

// ApplicationPublicRequest represents the request payload for public application list
type ApplicationPublicRequest struct {
	Filter struct {
		ShowInJumbotron *bool `json:"show_in_jumbotron"`
	} `json:"filter"`
}

// ApplicationPublicResponse represents the public response for applications
type ApplicationPublicResponse struct {
	ID              uint   `json:"id"`
	Nama            string `json:"nama"`
	Link            string `json:"link"`
	ShowInJumbotron bool   `json:"show_in_jumbotron"`
}

// ApplicationPublicListResponse represents the public list response
type ApplicationPublicListResponse struct {
	Data  []ApplicationPublicResponse `json:"data"`
	Total int64                       `json:"total"`
}

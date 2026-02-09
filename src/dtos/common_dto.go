package dtos

// IDRequest represents a simple request with only ID
type IDRequest struct {
	ID uint `json:"id" binding:"required"`
}

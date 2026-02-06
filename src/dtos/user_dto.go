package dtos

// UserCreateRequest represents the request payload for creating User
type UserCreateRequest struct {
	// Add fields here
}

// UserUpdateRequest represents the request payload for updating User
type UserUpdateRequest struct {
	// Add fields here
}

// UserResponse represents the response payload for User
type UserResponse struct {
	ID uint `json:"id"`
	// Add fields here
}

// UserListResponse represents the response payload for listing User
type UserListResponse struct {
	Data []UserResponse `json:"data"`
	Total int64 `json:"total"`
}

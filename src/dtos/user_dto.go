package dtos

import (
	"time"
)

// UserCreateRequest represents the request payload for creating User
type UserCreateRequest struct {
	Nama    string `json:"nama" binding:"required"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	RoleIDs []uint `json:"role_ids" binding:"required,min=1"`
	Status  string `json:"status" binding:"required,oneof=active inactive suspended"`
}

// UserUpdateRequest represents the request payload for updating User
type UserUpdateRequest struct {
	Nama    string `json:"nama" binding:"omitempty"`
	Username string `json:"username" binding:"omitempty,min=3"`
	RoleIDs []uint `json:"role_ids" binding:"omitempty"`
	Status  string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UserUpdatePasswordRequest represents the request payload for updating password
type UserUpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserResponse represents the response payload for User
type UserResponse struct {
	ID        uint            `json:"id"`
	Nama      string          `json:"nama"`
	Username  string          `json:"username"`
	Roles     []RoleResponse  `json:"roles"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	CreatedByID *uint         `json:"created_by_id"`
	UpdatedByID *uint         `json:"updated_by_id"`
}

// UserListResponse represents the response payload for listing User
type UserListResponse struct {
	Data  []UserResponse `json:"data"`
	Total int64          `json:"total"`
}

// UserGetAllRequest represents the request payload for getting all users with filters
type UserGetAllRequest struct {
	Search struct {
		Nama     string `json:"nama"`
		Username string `json:"username"`
		RoleIDs  []uint `json:"role_ids"`
		Status   string `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	Page       int   `json:"page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// UserResponseDetail represents the detailed response payload for User
type UserResponseDetail struct {
	ID          uint            `json:"id"`
	Nama        string          `json:"nama"`
	Username    string          `json:"username"`
	Roles       []RoleResponse  `json:"roles"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
}

// UserListWithPaginationResponse represents the response with pagination
type UserListWithPaginationResponse struct {
	Data       []UserResponseDetail `json:"data"`
	Pagination PaginationInfo       `json:"pagination"`
}

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse represents the response payload for user login
type UserLoginResponse struct {
	ID        uint            `json:"id"`
	Nama      string          `json:"nama"`
	Username  string          `json:"username"`
	Roles     []RoleResponse  `json:"roles"`
	Status    string          `json:"status"`
	Token     string          `json:"token"`
	CreatedAt time.Time       `json:"created_at"`
}

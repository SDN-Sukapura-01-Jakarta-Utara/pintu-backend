package dtos

import "time"

// PelatihCreateRequest represents request to create new pelatih
type PelatihCreateRequest struct {
	Nama               string   `json:"nama" binding:"required"`
	Username           *string  `json:"username"`
	Password           *string  `json:"password"`
	Telepon            string   `json:"telepon"`
	Alamat             string   `json:"alamat"`
	Keahlian           string   `json:"keahlian"`
	Status             string   `json:"status" binding:"required,oneof=active inactive"`
	EkstrakurikulerIDs []uint   `json:"ekstrakurikuler_ids" binding:"required,min=1"`
	RoleIDs            []uint   `json:"role_ids" binding:"required,min=1"` // Required: must provide at least one role
}

// PelatihGetAllRequest represents request to get all pelatih with filters and pagination
type PelatihGetAllRequest struct {
	Search struct {
		Nama              string  `json:"nama"`
		EkstrakurikulerID *uint   `json:"ekstrakurikuler_id"`
		Status            string  `json:"status"`
	} `json:"search"`
	Pagination struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	} `json:"pagination"`
}

// PelatihResponse represents pelatih data in response
type PelatihResponse struct {
	ID             uint                        `json:"id"`
	Nama           string                      `json:"nama"`
	Username       *string                     `json:"username"`
	Telepon        string                      `json:"telepon"`
	Alamat         string                      `json:"alamat"`
	FotoProfil     *string                     `json:"foto_profil"`
	Keahlian       string                      `json:"keahlian"`
	Sertifikat     interface{}                 `json:"sertifikat"`
	Status         string                      `json:"status"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
	Ekstrakurikuler []EkstrakurikulerResponse  `json:"ekstrakurikuler,omitempty"`
	Roles          []RoleResponse              `json:"roles,omitempty"`
}

// PelatihListWithPaginationResponse represents list with pagination
type PelatihListWithPaginationResponse struct {
	Data       []PelatihResponse `json:"data"`
	Pagination PaginationInfo    `json:"pagination"`
}

// PelatihUpdateRequest represents request to update pelatih
type PelatihUpdateRequest struct {
	ID                  uint     `json:"id" binding:"required"`
	Nama                *string  `json:"nama"`
	Username            *string  `json:"username"`
	Password            *string  `json:"password"`
	Telepon             *string  `json:"telepon"`
	Alamat              *string  `json:"alamat"`
	Keahlian            *string  `json:"keahlian"`
	Status              *string  `json:"status"`
	EkstrakurikulerIDs  []uint   `json:"ekstrakurikuler_ids"`  // Optional: update ekstrakurikuler assignments
	RoleIDs             []uint   `json:"role_ids"`             // Optional: update role assignments
	SertifikatToDelete  []string `json:"sertifikat_to_delete"` // Optional: URLs of sertifikat to delete
	// Note: FotoProfil and Sertifikat are handled via multipart file upload in controller
}

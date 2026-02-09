package dtos

// TahunPelajaranCreateRequest represents the request payload for creating TahunPelajaran
type TahunPelajaranCreateRequest struct {
	// Add fields here
}

// TahunPelajaranUpdateRequest represents the request payload for updating TahunPelajaran
type TahunPelajaranUpdateRequest struct {
	// Add fields here
}

// TahunPelajaranResponse represents the response payload for TahunPelajaran
type TahunPelajaranResponse struct {
	ID uint `json:"id"`
	// Add fields here
}

// TahunPelajaranListResponse represents the response payload for listing TahunPelajaran
type TahunPelajaranListResponse struct {
	Data []TahunPelajaranResponse `json:"data"`
	Total int64 `json:"total"`
}

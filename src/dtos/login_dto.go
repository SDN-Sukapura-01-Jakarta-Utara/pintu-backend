package dtos

import "time"

// LoginRequest represents the request payload for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response payload for successful login
type LoginResponse struct {
	Token       string              `json:"token"`
	User        UserLoginResponse   `json:"user"`
	Permissions []string            `json:"permissions"`
	ExpiresAt   time.Time           `json:"expires_at"`
}

// StudentLoginResponse represents the response payload for successful student login
type StudentLoginResponse struct {
	Token       string                     `json:"token"`
	Student     StudentDetailLoginResponse `json:"student"`
	Permissions []string                   `json:"permissions"`
	ExpiresAt   time.Time                  `json:"expires_at"`
}

// StudentDetailLoginResponse represents student details in login response
type StudentDetailLoginResponse struct {
	ID                 uint                          `json:"id"`
	Nama               string                        `json:"nama"`
	NIS                string                        `json:"nis"`
	NISN               string                        `json:"nisn"`
	JenisKelamin       string                        `json:"jenis_kelamin"`
	TempatLahir        string                        `json:"tempat_lahir"`
	TanggalLahir       *time.Time                    `json:"tanggal_lahir"`
	Username           string                        `json:"username"`
	Status             string                        `json:"status"`
	Photo              string                        `json:"photo,omitempty"`
	Roles              []RoleResponse                `json:"roles"`
	Rombel             []StudentRombelLoginResponse  `json:"rombel"`
	CreatedAt          time.Time                     `json:"created_at"`
}

// StudentRombelLoginResponse represents rombel data in student login response
type StudentRombelLoginResponse struct {
	ID               uint                 `json:"id"`
	RombelID         uint                 `json:"rombel_id"`
	RombelName       string               `json:"rombel_name"`
	KelasID          uint                 `json:"kelas_id"`
	KelasName        string               `json:"kelas_name"`
	TahunPelajaranID uint                 `json:"tahun_pelajaran_id"`
	TahunPelajaran   string               `json:"tahun_pelajaran"`
	Status           string               `json:"status"`
}

package models

import "time"

// PelatihRole represents the many-to-many relationship between Pelatih and Role
type PelatihRole struct {
	PelatihID uint      `gorm:"primaryKey" json:"pelatih_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for PelatihRole
func (m *PelatihRole) TableName() string {
	return "pelatih_roles"
}

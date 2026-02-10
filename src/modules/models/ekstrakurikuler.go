package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// KelasIDs custom type for JSONB array
type KelasIDs []uint

// Value implements the driver.Valuer interface
func (k KelasIDs) Value() (driver.Value, error) {
	return json.Marshal(k)
}

// Scan implements the sql.Scanner interface
func (k *KelasIDs) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return gorm.ErrInvalidData
	}
	return json.Unmarshal(bytes, &k)
}

// Ekstrakurikuler represents the Ekstrakurikuler model
type Ekstrakurikuler struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"uniqueIndex;not null" json:"name"`
	KelasIDs    KelasIDs        `gorm:"type:jsonb;default:'[]'" json:"kelas_ids"`
	Kategori    string          `gorm:"index;not null" json:"kategori"`
	Status      string          `gorm:"default:active" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID *uint           `json:"created_by_id"`
	UpdatedByID *uint           `json:"updated_by_id"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Ekstrakurikuler
func (m *Ekstrakurikuler) TableName() string {
	return "ekstrakurikulers"
}

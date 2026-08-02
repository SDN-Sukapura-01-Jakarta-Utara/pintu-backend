package models

import (
	"time"

	"gorm.io/gorm"
)

// Formulir represents the Formulir model
type Formulir struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	Judul                 string         `gorm:"type:varchar(255);not null" json:"judul"`
	Slug                  string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Deskripsi             string         `gorm:"type:text" json:"deskripsi"`
	CreatedByUserID       uint           `gorm:"not null" json:"created_by_user_id"`
	CreatedByRole         *string        `gorm:"type:varchar(50)" json:"created_by_role"` // "pendidik", "tendik", "murid", or NULL
	IsActive              bool           `gorm:"default:true" json:"is_active"`
	MaxResponses          *int           `json:"max_responses"`
	StartDate             *time.Time     `gorm:"type:timestamp" json:"start_date"`
	EndDate               *time.Time     `gorm:"type:timestamp" json:"end_date"`
	AccessType            string         `gorm:"type:varchar(50);default:public" json:"access_type"`
	TargetUserTypes       []byte         `gorm:"type:jsonb" json:"target_user_types"`
	RombelIDs             []byte         `gorm:"type:jsonb;column:rombel_ids" json:"rombel_ids"` // Array of rombel IDs (for murid filtering)
	AllowMultipleResponses bool          `gorm:"default:false" json:"allow_multiple_responses"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Pertanyaan            []FormulirPertanyaan `gorm:"foreignKey:FormulirID;constraint:OnDelete:CASCADE" json:"pertanyaan,omitempty"`
	CreatedBy             *User                `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
}

// TableName specifies the table name for Formulir
func (m *Formulir) TableName() string {
	return "forms"
}

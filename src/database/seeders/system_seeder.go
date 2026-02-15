package seeders

import (
	"pintu-backend/src/modules/models"

	"gorm.io/gorm"
)

// SystemSeeder seeds system data
func SystemSeeder(db *gorm.DB) error {
	system := []models.System{
		{
			Nama:        "PINTU",
			Description: "Portal Informasi Terpadu",
			Status:      "active",
		},
		{
			Nama:        "SIEKSA",
			Description: "Sistem Informasi Ekstrakurikuler Sukapura Satu",
			Status:      "active",
		},
		{
			Nama:        "SIPERSA",
			Description: "Sistem Informasi Perpustakaan Sukapura Satu",
			Status:      "active",
		},
	}

	return db.CreateInBatches(system, 100).Error
}

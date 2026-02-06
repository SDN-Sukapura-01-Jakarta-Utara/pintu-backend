package seeders

import (
	"fmt"
	"gorm.io/gorm"
)

// UserSeeder handles seeding data for User
type UserSeeder struct {
	db *gorm.DB
}

// NewUserSeeders creates a new User seeder
func NewUserSeeders(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

// Run executes the seeder
func (s *UserSeeder) Run() error {
	fmt.Println("Seeding User...")

	// Add your seeding logic here

	fmt.Println("%!s(MISSING) seeding completed")
	return nil
}

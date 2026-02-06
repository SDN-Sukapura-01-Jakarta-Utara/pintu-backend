package seeders

import (
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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

// UserData represents user record for seeding
type UserData struct {
	Nama             string
	Username         string
	Password         string
	RoleName         string
	AccessibleSystem string
	Status           string
}

// Run executes the seeder
func (s *UserSeeder) Run() error {
	fmt.Println("Seeding users...")

	users := []UserData{
		{
			Nama:             "Administrator",
			Username:         "admin",
			Password:         "admin01", // Change this in production
			RoleName:         "Administrator",
			AccessibleSystem: `["PINTU"]`,
			Status:           "active",
		},
		{
			Nama:             "Kepala Sekolah",
			Username:         "kepsek",
			Password:         "kepsek01", // Change this in production
			RoleName:         "Kepala Sekolah",
			AccessibleSystem: `["PINTU"]`,
			Status:           "active",
		},
	}

	// Check if users already exist
	var count int64
	s.db.Table("users").Count(&count)
	if count > 0 {
		fmt.Println("✓ Users already seeded")
		return nil
	}

	// Insert users
	for _, user := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", user.Username, err)
		}

		// Get role ID by name
		var roleID uint
		if err := s.db.Table("roles").Where("name = ?", user.RoleName).Select("id").Row().Scan(&roleID); err != nil {
			return fmt.Errorf("failed to find role %s: %w", user.RoleName, err)
		}

		// Validate accessible_system is valid JSON
		var validJSON interface{}
		if err := json.Unmarshal([]byte(user.AccessibleSystem), &validJSON); err != nil {
			return fmt.Errorf("invalid JSON for accessible_system: %w", err)
		}

		// Create user
		result := s.db.Table("users").Create(map[string]interface{}{
			"nama":              user.Nama,
			"username":          user.Username,
			"password":          string(hashedPassword),
			"role_id":           roleID,
			"accessible_system": user.AccessibleSystem,
			"status":            user.Status,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Username, result.Error)
		}
	}

	fmt.Printf("✓ %d users seeded successfully\n", len(users))
	fmt.Println("⚠️  Remember to change default passwords in production!")
	return nil
}

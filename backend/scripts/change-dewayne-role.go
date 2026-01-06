package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Get school_admin role ID
	var schoolAdminRoleID uint
	err = db.Raw("SELECT id FROM roles WHERE name = 'school_admin'").Scan(&schoolAdminRoleID).Error
	if err != nil {
		log.Fatal("Failed to get school_admin role:", err)
	}

	// Update DeWayne's role
	result := db.Exec("UPDATE users SET role_id = ? WHERE email = 'dewayne.frazier@aun.edu.ng'", schoolAdminRoleID)
	if result.Error != nil {
		log.Fatal("Failed to update user:", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Fatal("User not found")
	}

	fmt.Printf("✓ Updated DeWayne to school_admin (role_id: %d)\n", schoolAdminRoleID)
}

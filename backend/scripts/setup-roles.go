package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
}

type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	Category    string
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Run migrations
	db.AutoMigrate(&Role{}, &Permission{}, &RolePermission{})
	fmt.Println("✓ Tables created")

	// Create roles
	roles := []Role{
		{Name: "platform_admin", Description: "Platform Administrator"},
		{Name: "school_admin", Description: "School Administrator"},
		{Name: "teacher", Description: "Teacher"},
		{Name: "student", Description: "Student"},
	}

	for _, r := range roles {
		db.FirstOrCreate(&r, Role{Name: r.Name})
		fmt.Printf("✓ Created role: %s\n", r.Name)
	}

	fmt.Println("\n✓ All roles created")
}

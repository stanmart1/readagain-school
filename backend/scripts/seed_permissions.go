package main

import (
	"log"

	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/models"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("✅ Database connected")

	// Define all permissions
	permissions := []models.Permission{
		// Analytics
		{Name: "analytics.view", Description: "View analytics dashboard", Category: "analytics"},
		{Name: "analytics.export", Description: "Export analytics data", Category: "analytics"},

		// Users
		{Name: "users.view", Description: "View users", Category: "users"},
		{Name: "users.create", Description: "Create users", Category: "users"},
		{Name: "users.edit", Description: "Edit users", Category: "users"},
		{Name: "users.delete", Description: "Delete users", Category: "users"},
		{Name: "users.manage", Description: "Full user management", Category: "users"},

		// Roles
		{Name: "roles.view", Description: "View roles", Category: "roles"},
		{Name: "roles.create", Description: "Create roles", Category: "roles"},
		{Name: "roles.edit", Description: "Edit roles", Category: "roles"},
		{Name: "roles.delete", Description: "Delete roles", Category: "roles"},
		{Name: "roles.manage", Description: "Full role management", Category: "roles"},

		// Audit Logs
		{Name: "audit_logs.view", Description: "View audit logs", Category: "audit"},
		{Name: "audit_logs.export", Description: "Export audit logs", Category: "audit"},

		// Books
		{Name: "books.view", Description: "View books", Category: "books"},
		{Name: "books.create", Description: "Create books", Category: "books"},
		{Name: "books.edit", Description: "Edit books", Category: "books"},
		{Name: "books.delete", Description: "Delete books", Category: "books"},
		{Name: "books.manage", Description: "Full book management", Category: "books"},

		// Reviews
		{Name: "reviews.view", Description: "View reviews", Category: "reviews"},
		{Name: "reviews.moderate", Description: "Moderate reviews", Category: "reviews"},
		{Name: "reviews.delete", Description: "Delete reviews", Category: "reviews"},
		{Name: "reviews.manage", Description: "Full review management", Category: "reviews"},

		// Orders
		{Name: "orders.view", Description: "View orders", Category: "orders"},
		{Name: "orders.edit", Description: "Edit orders", Category: "orders"},
		{Name: "orders.delete", Description: "Delete orders", Category: "orders"},
		{Name: "orders.manage", Description: "Full order management", Category: "orders"},

		// Shipping
		{Name: "shipping.view", Description: "View shipping", Category: "shipping"},
		{Name: "shipping.manage", Description: "Manage shipping", Category: "shipping"},

		// Reading Analytics
		{Name: "reading.view_analytics", Description: "View reading analytics", Category: "reading"},
		{Name: "reading.manage", Description: "Manage reading data", Category: "reading"},

		// Reports
		{Name: "reports.view", Description: "View reports", Category: "reports"},
		{Name: "reports.generate", Description: "Generate reports", Category: "reports"},
		{Name: "reports.export", Description: "Export reports", Category: "reports"},

		// Email Templates
		{Name: "email_templates.view", Description: "View email templates", Category: "email"},
		{Name: "email_templates.create", Description: "Create email templates", Category: "email"},
		{Name: "email_templates.edit", Description: "Edit email templates", Category: "email"},
		{Name: "email_templates.delete", Description: "Delete email templates", Category: "email"},

		// Blog
		{Name: "blog.view", Description: "View blog posts", Category: "blog"},
		{Name: "blog.create", Description: "Create blog posts", Category: "blog"},
		{Name: "blog.edit", Description: "Edit blog posts", Category: "blog"},
		{Name: "blog.delete", Description: "Delete blog posts", Category: "blog"},
		{Name: "blog.publish", Description: "Publish blog posts", Category: "blog"},

		// Works/Portfolio
		{Name: "works.view", Description: "View works", Category: "works"},
		{Name: "works.create", Description: "Create works", Category: "works"},
		{Name: "works.edit", Description: "Edit works", Category: "works"},
		{Name: "works.delete", Description: "Delete works", Category: "works"},

		// About
		{Name: "about.view", Description: "View about page", Category: "about"},
		{Name: "about.edit", Description: "Edit about page", Category: "about"},

		// Contact
		{Name: "contact.view", Description: "View contact messages", Category: "contact"},
		{Name: "contact.reply", Description: "Reply to contact messages", Category: "contact"},
		{Name: "contact.delete", Description: "Delete contact messages", Category: "contact"},

		// FAQ
		{Name: "faq.view", Description: "View FAQs", Category: "faq"},
		{Name: "faq.create", Description: "Create FAQs", Category: "faq"},
		{Name: "faq.edit", Description: "Edit FAQs", Category: "faq"},
		{Name: "faq.delete", Description: "Delete FAQs", Category: "faq"},

		// Settings
		{Name: "settings.view", Description: "View settings", Category: "settings"},
		{Name: "settings.edit", Description: "Edit settings", Category: "settings"},
		{Name: "settings.manage", Description: "Full settings management", Category: "settings"},
	}

	// Seed permissions
	log.Println("\n📝 Seeding permissions...")
	for _, perm := range permissions {
		var existing models.Permission
		err := database.DB.Where("name = ?", perm.Name).First(&existing).Error
		if err != nil {
			if err := database.DB.Create(&perm).Error; err != nil {
				log.Printf("❌ Failed to create permission %s: %v", perm.Name, err)
			} else {
				log.Printf("✓ Created permission: %s", perm.Name)
			}
		} else {
			log.Printf("- Permission already exists: %s", perm.Name)
		}
	}

	// Get all permissions from database
	var allPermissions []models.Permission
	database.DB.Find(&allPermissions)

	// Get SuperAdmin role (ID: 1)
	var superAdminRole models.Role
	if err := database.DB.First(&superAdminRole, 1).Error; err != nil {
		log.Fatal("❌ SuperAdmin role not found")
	}

	// Get Admin role (ID: 2)
	var adminRole models.Role
	if err := database.DB.First(&adminRole, 2).Error; err != nil {
		log.Fatal("❌ Admin role not found")
	}

	// Assign ALL permissions to SuperAdmin
	log.Println("\n🔐 Assigning permissions to SuperAdmin...")
	if err := database.DB.Model(&superAdminRole).Association("Permissions").Replace(allPermissions); err != nil {
		log.Fatal("❌ Failed to assign permissions to SuperAdmin:", err)
	}
	log.Printf("✅ Assigned %d permissions to SuperAdmin", len(allPermissions))

	// Assign most permissions to Admin (exclude some sensitive ones)
	log.Println("\n🔐 Assigning permissions to Admin...")
	adminPermissions := []models.Permission{}
	excludeForAdmin := map[string]bool{
		"roles.delete":    true,
		"users.delete":    true,
		"settings.manage": true,
	}

	for _, perm := range allPermissions {
		if !excludeForAdmin[perm.Name] {
			adminPermissions = append(adminPermissions, perm)
		}
	}

	if err := database.DB.Model(&adminRole).Association("Permissions").Replace(adminPermissions); err != nil {
		log.Fatal("❌ Failed to assign permissions to Admin:", err)
	}
	log.Printf("✅ Assigned %d permissions to Admin", len(adminPermissions))

	log.Println("\n✅ Permission seeding completed!")
	log.Printf("   Total permissions: %d", len(allPermissions))
	log.Printf("   SuperAdmin permissions: %d", len(allPermissions))
	log.Printf("   Admin permissions: %d", len(adminPermissions))
}

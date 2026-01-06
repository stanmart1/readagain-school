package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	"readagain/internal/models"
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

	var count int64
	db.Model(&models.Permission{}).Count(&count)
	fmt.Printf("Current permissions count: %d\n", count)

	if count > 0 {
		fmt.Println("Permissions already exist")
		return
	}

	// Create comprehensive permissions
	permissions := []models.Permission{
		// Users
		{Name: "users.view", Description: "View users", Category: "Users"},
		{Name: "users.create", Description: "Create users", Category: "Users"},
		{Name: "users.edit", Description: "Edit users", Category: "Users"},
		{Name: "users.delete", Description: "Delete users", Category: "Users"},
		{Name: "users.activate", Description: "Activate/deactivate users", Category: "Users"},
		{Name: "users.export", Description: "Export users", Category: "Users"},
		
		// Roles
		{Name: "roles.view", Description: "View roles", Category: "Roles"},
		{Name: "roles.create", Description: "Create roles", Category: "Roles"},
		{Name: "roles.edit", Description: "Edit roles", Category: "Roles"},
		{Name: "roles.delete", Description: "Delete roles", Category: "Roles"},
		
		// Permissions
		{Name: "permissions.view", Description: "View permissions", Category: "Permissions"},
		{Name: "permissions.assign", Description: "Assign permissions", Category: "Permissions"},
		
		// Books
		{Name: "books.view", Description: "View books", Category: "Books"},
		{Name: "books.create", Description: "Create books", Category: "Books"},
		{Name: "books.edit", Description: "Edit books", Category: "Books"},
		{Name: "books.delete", Description: "Delete books", Category: "Books"},
		{Name: "books.publish", Description: "Publish books", Category: "Books"},
		{Name: "books.feature", Description: "Feature books", Category: "Books"},
		{Name: "books.export", Description: "Export books", Category: "Books"},
		
		// Authors
		{Name: "authors.view", Description: "View authors", Category: "Authors"},
		{Name: "authors.create", Description: "Create authors", Category: "Authors"},
		{Name: "authors.edit", Description: "Edit authors", Category: "Authors"},
		{Name: "authors.delete", Description: "Delete authors", Category: "Authors"},
		{Name: "authors.approve", Description: "Approve authors", Category: "Authors"},
		
		// Categories
		{Name: "categories.view", Description: "View categories", Category: "Categories"},
		{Name: "categories.create", Description: "Create categories", Category: "Categories"},
		{Name: "categories.edit", Description: "Edit categories", Category: "Categories"},
		{Name: "categories.delete", Description: "Delete categories", Category: "Categories"},
		
		// Library
		{Name: "library.view", Description: "View library", Category: "Library"},
		{Name: "library.assign", Description: "Assign books to users", Category: "Library"},
		{Name: "library.remove", Description: "Remove books from library", Category: "Library"},
		{Name: "library.export", Description: "Export library data", Category: "Library"},
		
		// Reviews
		{Name: "reviews.view", Description: "View reviews", Category: "Reviews"},
		{Name: "reviews.create", Description: "Create reviews", Category: "Reviews"},
		{Name: "reviews.edit", Description: "Edit reviews", Category: "Reviews"},
		{Name: "reviews.delete", Description: "Delete reviews", Category: "Reviews"},
		{Name: "reviews.approve", Description: "Approve reviews", Category: "Reviews"},
		{Name: "reviews.feature", Description: "Feature reviews", Category: "Reviews"},
		
		// Reading Analytics
		{Name: "reading.view_own", Description: "View own reading data", Category: "Reading"},
		{Name: "reading.view_all", Description: "View all reading data", Category: "Reading"},
		{Name: "reading.view_analytics", Description: "View reading analytics", Category: "Reading"},
		{Name: "reading.export", Description: "Export reading data", Category: "Reading"},
		
		// Analytics
		{Name: "analytics.view", Description: "View analytics dashboard", Category: "Analytics"},
		{Name: "analytics.export", Description: "Export analytics", Category: "Analytics"},
		
		// Audit Logs
		{Name: "audit_logs.view", Description: "View audit logs", Category: "Audit"},
		{Name: "audit_logs.export", Description: "Export audit logs", Category: "Audit"},
		
		// Reports
		{Name: "reports.view", Description: "View reports", Category: "Reports"},
		{Name: "reports.create", Description: "Create reports", Category: "Reports"},
		{Name: "reports.export", Description: "Export reports", Category: "Reports"},
		
		// Blog
		{Name: "blog.view", Description: "View blog posts", Category: "Blog"},
		{Name: "blog.create", Description: "Create blog posts", Category: "Blog"},
		{Name: "blog.edit", Description: "Edit blog posts", Category: "Blog"},
		{Name: "blog.delete", Description: "Delete blog posts", Category: "Blog"},
		{Name: "blog.publish", Description: "Publish blog posts", Category: "Blog"},
		
		// FAQ
		{Name: "faq.view", Description: "View FAQs", Category: "FAQ"},
		{Name: "faq.create", Description: "Create FAQs", Category: "FAQ"},
		{Name: "faq.edit", Description: "Edit FAQs", Category: "FAQ"},
		{Name: "faq.delete", Description: "Delete FAQs", Category: "FAQ"},
		
		// About
		{Name: "about.view", Description: "View about page", Category: "About"},
		{Name: "about.edit", Description: "Edit about page", Category: "About"},
		
		// Contact
		{Name: "contact.view", Description: "View contact messages", Category: "Contact"},
		{Name: "contact.reply", Description: "Reply to contact messages", Category: "Contact"},
		{Name: "contact.delete", Description: "Delete contact messages", Category: "Contact"},
		
		// Testimonials
		{Name: "testimonials.view", Description: "View testimonials", Category: "Testimonials"},
		{Name: "testimonials.create", Description: "Create testimonials", Category: "Testimonials"},
		{Name: "testimonials.edit", Description: "Edit testimonials", Category: "Testimonials"},
		{Name: "testimonials.delete", Description: "Delete testimonials", Category: "Testimonials"},
		
		// Settings
		{Name: "settings.view", Description: "View settings", Category: "Settings"},
		{Name: "settings.manage", Description: "Manage settings", Category: "Settings"},
		
		// Notifications
		{Name: "notifications.view", Description: "View notifications", Category: "Notifications"},
		{Name: "notifications.send", Description: "Send notifications", Category: "Notifications"},
		{Name: "notifications.delete", Description: "Delete notifications", Category: "Notifications"},
		
		// Groups
		{Name: "groups.view", Description: "View groups", Category: "Groups"},
		{Name: "groups.create", Description: "Create groups", Category: "Groups"},
		{Name: "groups.edit", Description: "Edit groups", Category: "Groups"},
		{Name: "groups.delete", Description: "Delete groups", Category: "Groups"},
		{Name: "groups.manage_members", Description: "Manage group members", Category: "Groups"},
		
		// Achievements
		{Name: "achievements.view", Description: "View achievements", Category: "Achievements"},
		{Name: "achievements.create", Description: "Create achievements", Category: "Achievements"},
		{Name: "achievements.edit", Description: "Edit achievements", Category: "Achievements"},
		{Name: "achievements.delete", Description: "Delete achievements", Category: "Achievements"},
		{Name: "achievements.award", Description: "Award achievements", Category: "Achievements"},
		
		// Chat
		{Name: "chat.view", Description: "View chat rooms", Category: "Chat"},
		{Name: "chat.create", Description: "Create chat rooms", Category: "Chat"},
		{Name: "chat.send", Description: "Send messages", Category: "Chat"},
		{Name: "chat.moderate", Description: "Moderate chat", Category: "Chat"},
		{Name: "chat.delete", Description: "Delete messages", Category: "Chat"},
	}

	for _, p := range permissions {
		db.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}

	db.Model(&models.Permission{}).Count(&count)
	fmt.Printf("\n✓ Created %d permissions\n", count)
}

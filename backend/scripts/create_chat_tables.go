package main

import (
	"fmt"
	"log"

	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/models"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create chat tables
	fmt.Println("Creating chat tables...")

	if err := database.DB.AutoMigrate(
		&models.ChatRoom{},
		&models.ChatMessage{},
		&models.ChatMember{},
		&models.ChatReaction{},
	); err != nil {
		log.Fatalf("Failed to create chat tables: %v", err)
	}

	// Add indexes for performance
	fmt.Println("Adding indexes...")

	// ChatRoom indexes
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_rooms_type ON chat_rooms(type)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_rooms_group_id ON chat_rooms(group_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_rooms_book_id ON chat_rooms(book_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_rooms_created_by ON chat_rooms(created_by)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_rooms_is_active ON chat_rooms(is_active)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_rooms_last_message_at ON chat_rooms(last_message_at)")

	// ChatMessage indexes
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_messages_room_id ON chat_messages(room_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_messages_user_id ON chat_messages(user_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_messages_reply_to_id ON chat_messages(reply_to_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_messages_is_deleted ON chat_messages(is_deleted)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at)")

	// ChatMember indexes
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_members_room_id ON chat_members(room_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_members_user_id ON chat_members(user_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_members_room_user ON chat_members(room_id, user_id)")

	// ChatReaction indexes
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_reactions_message_id ON chat_reactions(message_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_reactions_user_id ON chat_reactions(user_id)")
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_chat_reactions_unique ON chat_reactions(message_id, user_id, emoji)")

	fmt.Println("✅ Chat tables created successfully!")
}

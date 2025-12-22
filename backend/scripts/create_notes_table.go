package main

import (
	"log"
	"readagain/internal/config"
	"readagain/internal/database"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg.Database.URL)

	// Create notes table
	err := database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE,
			user_library_id INTEGER NOT NULL REFERENCES user_libraries(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			page_number INTEGER,
			position TEXT,
			color VARCHAR(50) DEFAULT 'yellow'
		);

		CREATE INDEX IF NOT EXISTS idx_notes_user_library_id ON notes(user_library_id);
		CREATE INDEX IF NOT EXISTS idx_notes_deleted_at ON notes(deleted_at);
	`).Error

	if err != nil {
		log.Fatalf("Failed to create notes table: %v", err)
	}

	// Add book_id to reading_goals table
	err = database.DB.Exec(`
		ALTER TABLE reading_goals ADD COLUMN IF NOT EXISTS book_id INTEGER REFERENCES books(id) ON DELETE CASCADE;
		CREATE INDEX IF NOT EXISTS idx_reading_goals_book_id ON reading_goals(book_id);
	`).Error

	if err != nil {
		log.Fatalf("Failed to add book_id to reading_goals: %v", err)
	}

	log.Println("✅ Notes table created and book_id added to reading_goals successfully")
}

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment variable (Coolify sets this)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	fmt.Println("Initializing production database...")
	fmt.Printf("Database URL: %s\n", maskPassword(dbURL))

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✓ Database connection successful")

	// Get the directory of the current executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	scriptDir := filepath.Dir(execPath)
	
	// Look for SQL file in the same directory as the executable
	sqlFile := filepath.Join(scriptDir, "production-db-backup.sql")
	
	// If not found, try relative to current working directory
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		sqlFile = "scripts/production-db-backup.sql"
	}
	
	// If still not found, try in the scripts directory
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		sqlFile = "production-db-backup.sql"
	}

	fmt.Printf("Reading SQL file: %s\n", sqlFile)
	
	// Read SQL file
	sqlContent, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// Clean and prepare SQL content
	sqlString := string(sqlContent)
	
	// Remove comments and empty lines
	lines := strings.Split(sqlString, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "--") {
			cleanLines = append(cleanLines, line)
		}
	}
	
	cleanSQL := strings.Join(cleanLines, "\n")
	
	// Split into individual statements
	statements := strings.Split(cleanSQL, ";")
	
	fmt.Println("Executing SQL statements...")
	
	// Execute each statement
	successCount := 0
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		fmt.Printf("Executing statement %d...\n", i+1)
		
		_, err := db.Exec(stmt)
		if err != nil {
			// Log error but continue with other statements
			fmt.Printf("Warning: Failed to execute statement %d: %v\n", i+1, err)
			fmt.Printf("Statement: %s\n", truncateString(stmt, 100))
		} else {
			successCount++
		}
	}
	
	fmt.Printf("✓ Database initialization completed!\n")
	fmt.Printf("✓ Successfully executed %d statements\n", successCount)
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("1. Verify the database schema and data")
	fmt.Println("2. Update any application configuration as needed")
	fmt.Println("3. Start the application server")
}

// maskPassword masks the password in database URL for logging
func maskPassword(dbURL string) string {
	if strings.Contains(dbURL, "@") {
		parts := strings.Split(dbURL, "@")
		if len(parts) == 2 {
			userPart := parts[0]
			hostPart := parts[1]
			
			if strings.Contains(userPart, ":") {
				userParts := strings.Split(userPart, ":")
				if len(userParts) >= 2 {
					userParts[len(userParts)-1] = "****"
					userPart = strings.Join(userParts, ":")
				}
			}
			
			return userPart + "@" + hostPart
		}
	}
	return dbURL
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

package main

import (
	"log"
	"readagain/internal/config"
	"readagain/internal/database"
)

func main() {
	cfg := config.Load()
	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	tables := []string{
		"order_items",
		"orders",
		"shipping_addresses",
		"payments",
		"bank_transfers",
	}

	for _, table := range tables {
		if err := database.DB.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			log.Printf("Failed to drop %s table: %v", table, err)
		} else {
			log.Printf("✓ Dropped %s table", table)
		}
	}

	log.Println("\n✅ All order and shipping tables dropped successfully!")
}

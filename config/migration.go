package config

import (
	"fmt"
	"log"

	"github.com/ichtrojan/thoth"
)

// MigrateDB handles database migrations for adding timestamp fields
func MigrateDB() {
	logger, _ := thoth.Init("logs")
	
	database := Database()
	
	// Add created_at column
	_, err := database.Exec(`ALTER TABLE todos ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`)
	if err != nil {
		fmt.Println("Error adding created_at column:", err)
		logger.Log(err)
	} else {
		fmt.Println("Successfully added created_at column")
	}
	
	// Add updated_at column
	_, err = database.Exec(`ALTER TABLE todos ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`)
	if err != nil {
		fmt.Println("Error adding updated_at column:", err)
		logger.Log(err)
	} else {
		fmt.Println("Successfully added updated_at column")
	}
	
	// Add completed_at column (allows NULL)
	_, err = database.Exec(`ALTER TABLE todos ADD COLUMN completed_at TIMESTAMP NULL`)
	if err != nil {
		fmt.Println("Error adding completed_at column:", err)
		logger.Log(err)
	} else {
		fmt.Println("Successfully added completed_at column")
	}
	
	// Set default values for existing rows
	_, err = database.Exec(`UPDATE todos SET created_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE created_at IS NULL`)
	if err != nil {
		fmt.Println("Error setting default timestamps for existing rows:", err)
		logger.Log(err)
	} else {
		fmt.Println("Successfully updated timestamps for existing rows")
	}
	
	// Close the database connection
	if err := database.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
		logger.Log(err)
	}
}
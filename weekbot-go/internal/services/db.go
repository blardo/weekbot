package services

import (
	"fmt"
	"os"
	"strings"
	"weekbot-go/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GetDB creates or returns a database connection for the given guild ID
// If RESET_DB=true, it will delete the existing database file first
func GetDB(gid string) (*gorm.DB, error) {
	gdbName := gid + ".db"
	if strings.EqualFold(os.Getenv("RESET_DB"), "true") {
		if err := os.Remove(gdbName); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	
	db, err := gorm.Open(sqlite.Open(gdbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Suggestion{},
		&models.Ballot{},
		&models.Voter{},
		&models.Poll{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

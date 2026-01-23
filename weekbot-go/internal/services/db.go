package services

import (
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GetDB creates or returns a database connection for the given guild ID
// If RESET_DB=true, it will delete the existing database file first
// Note: Migrations are handled by models/bot.go to avoid import cycles
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

	return db, nil
}

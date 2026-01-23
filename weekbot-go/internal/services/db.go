package services

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GetDB creates or returns a database connection for the given guild ID
// If Config.ResetDB is true, it will delete the existing database file first
// Note: Migrations are handled by models/bot.go to avoid import cycles
func GetDB(gid string) (*gorm.DB, error) {
	gdbName := gid + ".db"
	config := GetConfig()
	if config.ResetDB {
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

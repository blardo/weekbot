package db

import (
	"weekbot/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)



func GetDB(gid string) (*gorm.DB, error) {
	gdbName := gid + ".db"
	db, err := gorm.Open(sqlite.Open(gdbName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

	

func GetCurrentSuggestions(guildId string) ([]models.Suggestion, error) {
	db, err := GetDB(guildId)
	if err != nil {
		return nil, err
	}
	
	var suggestions []models.Suggestion
	db.Where("guild_id = ?", guildId).Find(&suggestions, "Updicks > 3" )
	
	return suggestions, nil
}

func GetAllSuggestions(guildId string) ([]models.Suggestion, error) {
	println("Getting all suggestions")
	db, err := GetDB(guildId)
	if err != nil {
		println("Error getting db")
		return nil, err
	}
	
	var suggestions []models.Suggestion
	db.Where("guild_id = ?", guildId).Find(&suggestions, "Updicks >= 0" )
	
	
	
	return suggestions, nil
}

func GetTables(guildId string) ([]string, error) {
	db, err := GetDB(guildId)
	if err != nil {
		return nil, err
	}
	
	var tables []string
	db.Raw("SELECT name FROM sqlite_master WHERE type='table';").Scan(&tables)
	
	return tables, nil
}

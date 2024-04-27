package db

import (
	"weekbot/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
 
func Start()
	
func connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("weekbot.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}


func GetCurrentSuggestions(guildId *string) ([]models.Suggestion, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	
	var suggestions []models.Suggestion
	db.Where("guild_id = ?", guildId).Find(&suggestions, "Updicks > 3" )
	
	return suggestions, nil
}
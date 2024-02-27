package weekbot

import (
	"os"

	"gorm.io/gorm"
)

type Config struct {
	DiscordToken       string
	DataBaseConnection *gorm.DB
}

// GetConfig returns the configuration for the bot from the environment
func GetConfig(db *gorm.DB) Config {
	config := Config{
		DiscordToken:       os.Getenv("DISCORD_TOKEN"),
		DataBaseConnection: db,
	}
	return config
}

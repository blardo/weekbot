package services

import (
	"os"

	"gorm.io/gorm"
)

type Config struct {
	DiscordToken string
	DB           *gorm.DB
}

// GetConfig returns the configuration for the bot from the environment
func GetConfig(db *gorm.DB) *Config {
	config := Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		DB:           db,
	}
	return &config
}

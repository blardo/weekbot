package services

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	AppID        string
	ResetDB      bool
	LogFormat    string
	Env          string
	OpenAIAPIKey string
}

var globalConfig *Config
var configOnce sync.Once

// GetConfig returns the configuration for the bot from the environment
// Loads from .env file if present, then overrides with actual environment variables
func GetConfig() *Config {
	configOnce.Do(func() {
		// Load .env file (ignores errors if file doesn't exist)
		_ = godotenv.Load()
		
		// Also try loading from weekbot-go/.env if we're in that directory
		_ = godotenv.Load(".env")
		
		resetDB := strings.EqualFold(os.Getenv("RESET_DB"), "true")
		
		globalConfig = &Config{
			DiscordToken: os.Getenv("DISCORD_TOKEN"),
			AppID:        os.Getenv("APP_ID"),
			ResetDB:      resetDB,
			LogFormat:    os.Getenv("LOG_FORMAT"),
			Env:          os.Getenv("ENV"),
			OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		}
	})
	return globalConfig
}

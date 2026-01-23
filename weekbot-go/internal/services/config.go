package services

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken          string
	AppID                 string
	ResetDB               bool
	LogFormat             string
	Env                   string
	GeminiAPIKey          string
	MinSuggestionsToStart int
	MinUpdicksToQualify   int
	MinBallotsToEndPoll   int
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
		
		// Parse threshold values with defaults
		minSuggestionsToStart := parseIntEnv("MIN_SUGGESTIONS_TO_START_POLL", 3)
		minUpdicksToQualify := parseIntEnv("MIN_UPDICKS_TO_QUALIFY", 3)
		minBallotsToEndPoll := parseIntEnv("MIN_BALLOTS_TO_END_POLL", 5)
		
		globalConfig = &Config{
			DiscordToken:          os.Getenv("DISCORD_TOKEN"),
			AppID:                 os.Getenv("APP_ID"),
			ResetDB:               resetDB,
			LogFormat:             os.Getenv("LOG_FORMAT"),
			Env:                   os.Getenv("ENV"),
			GeminiAPIKey:          os.Getenv("GEMINI_API_KEY"),
			MinSuggestionsToStart: minSuggestionsToStart,
			MinUpdicksToQualify:   minUpdicksToQualify,
			MinBallotsToEndPoll:   minBallotsToEndPoll,
		}
	})
	return globalConfig
}

// parseIntEnv parses an integer environment variable, returning defaultValue if not set or invalid
func parseIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

package weekbot

import "os"

type Config struct {
	DiscordToken string
}

// GetConfig returns the configuration for the bot from the environment
func GetConfig() Config {
	config := Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
	}
	return config
}

package weekbot

import (
	"fmt"
	"os"
	"os/signal"

	"weekbot/services/discord"
)

func NewBot(config Config) *Bot {
	discordService, err := discord.NewDiscordService(config.DiscordToken)
	if err != nil {
		fmt.Println("Error connecting to Discord:", err)
		return nil
	}

	return &Bot{
		config: &config,
		dsc:    discordService,
	}
}

// Bot is the main struct for the weekbot package
type Bot struct {
	config *Config
	dsc    *discord.DiscordService
}

// Run starts the bot
func (b *Bot) Run() {
	defer b.dsc.Disconnect()

	// Connect to Discord using the token from the config
	fmt.Println("Connecting to Discord...")
	err := b.dsc.Connect()
	if err != nil {
		fmt.Println("Error connecting to Discord:", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	fmt.Println("Shutting down...")
}

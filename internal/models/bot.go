package models

import (
	"fmt"
	"weekbot/internal/router"
	"weekbot/internal/services"
	"weekbot/internal/services/discord"
)

// Bot is the main struct for the weekbot package
type Bot struct {
	config *services.Config
	dsc    *discord.DiscordService
	router *router.Router
}

func NewBot(config *services.Config, router *router.Router) (*Bot, error) {
	discordService, err := discord.NewDiscordService(config.DiscordToken)
	if err != nil {
		fmt.Println("Error Configuring Discord Client:", err)
		return nil, err
	}

	bot := &Bot{
		config: config,
		dsc:    discordService,
		router: router,
	}

	return bot, nil
}

// Run starts the bot, connects to Discord, and sets up the router
func (b *Bot) Start() {
	// Connect to Discord using the token from the config
	fmt.Println("Connecting to Discord...")
	go func() {
		err := b.dsc.Connect()
		if err != nil {
			fmt.Println("Error connecting to Discord:", err)
			return
		}
		b.router.Setup(b.dsc)
	}()
}

// Stop is a scaffolded function to stop the bot
func (b *Bot) Stop() {
	b.dsc.Disconnect()
}

// Connected is a shorthand function to check if the bot is connected to Discord
func (b *Bot) Connected() bool {
	return b.dsc.Connected()
}

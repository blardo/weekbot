package weekbot

import (
	"fmt"
	"weekbot/internal/handlers/commands"
	"weekbot/internal/services/discord"

	"github.com/bwmarrin/discordgo"
)

// Bot is the main struct for the weekbot package
type Bot struct {
	config *Config
	dsc    *discord.DiscordService
}

func NewBot(config Config) (*Bot, error) {
	discordService, err := discord.NewDiscordService(config.DiscordToken)
	if err != nil {
		fmt.Println("Error Configuring Discord Client:", err)
		return nil, err
	}

	bot := &Bot{
		config: &config,
		dsc:    discordService,
	}

	return bot, nil
}

// Run starts the bot
func (b *Bot) Run() {
	// Connect to Discord using the token from the config
	fmt.Println("Connecting to Discord...")
	go func() {
		err := b.dsc.Connect()
		if err != nil {
			fmt.Println("Error connecting to Discord:", err)
			return
		}
	}()

	fmt.Println("Shutting down...")
}

func (b *Bot) Stop() {
	b.dsc.Disconnect()
}

func (b *Bot) AddHandler(handler interface{}) {
	b.dsc.AddHandler(handler)
}

func (b *Bot) SetupCommands() {
	// Setup the bot's commands
	commandlist := []*discordgo.ApplicationCommand{
		commands.GetPingCommand(),
		commands.GetPollCommand(),
	}

	for _, command := range commandlist {
		err := b.dsc.AddSlashCommand(command)
		if err != nil {
			fmt.Println("Error adding command:", err)
		}
	}

	// Bind the router to the bot's session
	b.dsc.AddHandler(commands.ParseCommand)
}

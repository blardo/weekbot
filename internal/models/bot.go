package models

import (
	"fmt"
	"weekbot/internal/services"
	"weekbot/internal/services/discord"

	"gorm.io/gorm"
)

var botInstances = make(map[string]*Bot)

// Holder of all basic state for each bot instance. This is the primary interface
// for the bot to interact with Discord.

type Bot struct {
	GuildID string
	Config  *services.Config
	DB      *gorm.DB
}

func NewBot(config *services.Config, gid string) (*Bot, error) {
	db, err := services.GetDB(gid)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Config:  config,
		GuildID: gid,
		DB:      db,
	}

	configureSchema(db)

	botInstances[gid] = bot

	fmt.Println("Connected to guild", gid)

	return bot, nil
}

func GetBot(guildID string) *Bot {
	return botInstances[guildID]
}

func GetBotInstances() map[string]*Bot {
	return botInstances
}

func (b *Bot) SendMessage(channelID, message string) error {
	discordService, err := discord.GetDiscordService()
	if err != nil {
		return err
	}

	return discordService.SendMessage(channelID, message)
}

func configureSchema(db *gorm.DB) {
	db.AutoMigrate(&Suggestion{})
	db.AutoMigrate(&Poll{})
}

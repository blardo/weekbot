package models

import (
	"fmt"
	"weekbot-go/internal/logger"
	"weekbot-go/internal/services"

	"gorm.io/gorm"
)

var botInstances = make(map[string]*Bot)

// Holder of all basic state for each bot instance.
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

	// Auto-migrate all models
	err = db.AutoMigrate(
		&Suggestion{},
		&Ballot{},
		&Voter{},
		&Poll{},
		&ServerConfig{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	bot := &Bot{
		Config:  config,
		GuildID: gid,
		DB:      db,
	}

	botInstances[gid] = bot

	logger.Info("Connected to guild", "guild_id", gid)

	return bot, nil
}

func GetBot(guildID string) *Bot {
	return botInstances[guildID]
}

func GetBotInstances() map[string]*Bot {
	return botInstances
}

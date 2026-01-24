package models

import (
	"weekbot-go/internal/logger"

	"gorm.io/gorm"
)

// ServerConfig stores per-guild configuration settings
// Nil pointer values indicate "use global default"
type ServerConfig struct {
	gorm.Model
	GuildID               string  `gorm:"uniqueIndex"`
	MinSuggestionsToStart *int    // nil = use default (3)
	MinUpdicksToQualify   *int    // nil = use default (3)
	MinBallotsToEndPoll   *int    // nil = use default (5)
	ReactionEmoji         *string // nil = use default ("👍")
}

// GetServerConfig retrieves the server config for a guild, or nil if not found
func GetServerConfig(db *gorm.DB, guildID string) *ServerConfig {
	var config ServerConfig
	// Use Take instead of First to avoid "record not found" being logged by GORM
	result := db.Where("guild_id = ?", guildID).Take(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// This is expected when no config exists yet - not an error
			return nil
		}
		logger.Error("Error fetching server config", "error", result.Error, "guild_id", guildID)
		return nil
	}
	return &config
}

// GetOrCreateServerConfig retrieves the server config for a guild, creating one if it doesn't exist
func GetOrCreateServerConfig(db *gorm.DB, guildID string) *ServerConfig {
	config := GetServerConfig(db, guildID)
	if config != nil {
		return config
	}

	// Create new config with nil values (use defaults)
	config = &ServerConfig{
		GuildID: guildID,
	}
	result := db.Create(config)
	if result.Error != nil {
		logger.Error("Error creating server config", "error", result.Error, "guild_id", guildID)
		return nil
	}
	logger.Info("Created server config", "guild_id", guildID)
	return config
}

// UpdateServerConfigField updates a specific integer field in the server config
func UpdateServerConfigField(db *gorm.DB, guildID string, field string, value int) error {
	config := GetOrCreateServerConfig(db, guildID)
	if config == nil {
		return gorm.ErrRecordNotFound
	}

	switch field {
	case "min-suggestions-to-start-poll":
		config.MinSuggestionsToStart = &value
	case "min-updicks":
		config.MinUpdicksToQualify = &value
	case "min-ballots-to-end-poll":
		config.MinBallotsToEndPoll = &value
	default:
		return gorm.ErrInvalidField
	}

	result := db.Save(config)
	if result.Error != nil {
		logger.Error("Error updating server config", "error", result.Error, "guild_id", guildID, "field", field, "value", value)
		return result.Error
	}
	logger.Info("Updated server config", "guild_id", guildID, "field", field, "value", value)
	return nil
}

// UpdateServerConfigStringField updates a specific string field in the server config
func UpdateServerConfigStringField(db *gorm.DB, guildID string, field string, value string) error {
	config := GetOrCreateServerConfig(db, guildID)
	if config == nil {
		return gorm.ErrRecordNotFound
	}

	switch field {
	case "reaction-emoji":
		config.ReactionEmoji = &value
	default:
		return gorm.ErrInvalidField
	}

	result := db.Save(config)
	if result.Error != nil {
		logger.Error("Error updating server config", "error", result.Error, "guild_id", guildID, "field", field, "value", value)
		return result.Error
	}
	logger.Info("Updated server config", "guild_id", guildID, "field", field, "value", value)
	return nil
}

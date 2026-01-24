package config

import (
	"weekbot-go/internal/services"

	"gorm.io/gorm"
)

// serverConfigGetter is an interface to avoid circular imports with models package
type serverConfigGetter interface {
	GetMinSuggestionsToStart() *int
	GetMinUpdicksToQualify() *int
	GetMinBallotsToEndPoll() *int
}

// serverConfig holds all server-specific settings from the database
type serverConfig struct {
	MinSuggestionsToStart *int
	MinUpdicksToQualify   *int
	MinBallotsToEndPoll   *int
	ReactionEmoji         *string
}

// getServerConfig retrieves server config without importing models directly
// This uses a query approach to avoid circular imports
func getServerConfig(db *gorm.DB, guildID string) *serverConfig {
	if db == nil {
		return nil
	}

	var result serverConfig
	db.Table("server_configs").
		Select("min_suggestions_to_start, min_updicks_to_qualify, min_ballots_to_end_poll, reaction_emoji").
		Where("guild_id = ?", guildID).
		First(&result)

	return &result
}

// MinSuggestionsToStartPoll returns the minimum number of suggestions needed to start a poll
// Checks server-specific config first, then falls back to environment variable (default: 3)
func MinSuggestionsToStartPoll(db *gorm.DB, guildID string) int {
	if cfg := getServerConfig(db, guildID); cfg != nil && cfg.MinSuggestionsToStart != nil {
		return *cfg.MinSuggestionsToStart
	}
	return services.GetConfig().MinSuggestionsToStart
}

// MinUpdicksToQualify returns the minimum number of reactions needed for a suggestion to qualify
// Checks server-specific config first, then falls back to environment variable (default: 3)
func MinUpdicksToQualify(db *gorm.DB, guildID string) int {
	if cfg := getServerConfig(db, guildID); cfg != nil && cfg.MinUpdicksToQualify != nil {
		return *cfg.MinUpdicksToQualify
	}
	return services.GetConfig().MinUpdicksToQualify
}

// MinBallotsToEndPoll returns the minimum number of cast ballots needed to end a poll
// Checks server-specific config first, then falls back to environment variable (default: 5)
func MinBallotsToEndPoll(db *gorm.DB, guildID string) int {
	if cfg := getServerConfig(db, guildID); cfg != nil && cfg.MinBallotsToEndPoll != nil {
		return *cfg.MinBallotsToEndPoll
	}
	return services.GetConfig().MinBallotsToEndPoll
}

// QualifyingEmojiForGuild returns the emoji that qualifies suggestions for a specific guild
// Checks server-specific config first, then falls back to default (👍)
func QualifyingEmojiForGuild(db *gorm.DB, guildID string) string {
	if cfg := getServerConfig(db, guildID); cfg != nil && cfg.ReactionEmoji != nil {
		return *cfg.ReactionEmoji
	}
	return DefaultQualifyingEmoji
}

// Discord channel and emoji constants
const (
	// WeekNameChannelName is the name of the channel where week suggestions are posted
	WeekNameChannelName = "week-name"
	// DefaultQualifyingEmoji is the default emoji that qualifies suggestions
	DefaultQualifyingEmoji = "👍"
	// ConfirmationEmoji is the emoji added when a suggestion qualifies
	ConfirmationEmoji = "✅"
)

// AcceptableWeekSuffixes are the valid endings for week name suggestions
var AcceptableWeekSuffixes = []string{"week", "week.", "week!", "week?"}

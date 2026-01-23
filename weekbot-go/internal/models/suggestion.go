package models

import (
	"strings"
	"weekbot-go/internal/config"
	"weekbot-go/internal/logger"

	"gorm.io/gorm"
)

type Suggestion struct {
	gorm.Model
	Content string
	Used    bool
	GuildID string
	Updicks int
	PollID  uint
}

func NewSuggestion(db *gorm.DB, content string, guildID string) *Suggestion {
	content = strings.TrimSpace(content)
	s := &Suggestion{
		Content: content,
		GuildID: guildID,
		Used:    false,
		Updicks: 0,
	}
	result := db.Create(s)
	if result.Error != nil {
		logger.Error("Failed to create suggestion", "error", result.Error, "content", content, "guild_id", guildID)
		panic(result.Error)
	} else {
		logger.Debug("Suggestion created", "rows_affected", result.RowsAffected, "content", content, "guild_id", guildID)
	}

	// Debug: verify the saved record
	var savedSuggestion Suggestion
	db.First(&savedSuggestion, "content = ? AND guild_id = ?", content, guildID)
	logger.Debug("Saved suggestion verified", "content", savedSuggestion.Content, "guild_id", savedSuggestion.GuildID, "updicks", savedSuggestion.Updicks)

	return s
}

func NormalizeSuggestionContent(content string) string {
	return strings.TrimSpace(strings.ToLower(content))
}

func FindActiveSuggestionByContent(db *gorm.DB, content string, guildID string) (*Suggestion, bool) {
	normalized := NormalizeSuggestionContent(content)
	var suggestion Suggestion
	result := db.Where("guild_id = ? AND used = ? AND lower(content) = ?", guildID, false, normalized).First(&suggestion)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			logger.Error("Error finding suggestion", "error", result.Error, "content", content, "guild_id", guildID)
		}
		return nil, false
	}
	if suggestion.ID == 0 {
		return nil, false
	}
	return &suggestion, true
}

func UpdateSuggestion(db *gorm.DB, content string, guildID string, updicks int) {
	var suggestion Suggestion
	normalized := NormalizeSuggestionContent(content)
	result := db.Where("guild_id = ? AND lower(content) = ?", guildID, normalized).First(&suggestion)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Suggestion not found, create a new one
			suggestion = *NewSuggestion(db, content, guildID)
		} else {
			// Some other error occurred
			logger.Error("Error finding suggestion for update", "error", result.Error, "content", content, "guild_id", guildID)
		}
	}

	suggestion.Updicks = updicks
	result = db.Save(&suggestion)
	if result.Error != nil {
		logger.Error("Error saving suggestion update", "error", result.Error, "content", content, "guild_id", guildID, "updicks", updicks)
		return
	}

	logger.Debug("Suggestion updated", "rows_affected", result.RowsAffected, "content", content, "guild_id", guildID, "updicks", updicks)

}

func GetMostRecentUnusedSuggestions(db *gorm.DB) []Suggestion {
	var suggestions []Suggestion
	db.Where("used = ? AND updicks >= ?", false, config.MinUpdicksToQualify).Find(&suggestions)
	return suggestions
}

// GetAllSuggestions gets all suggestions for a guild
func GetAllSuggestions(db *gorm.DB, guildID string) ([]Suggestion, error) {
	var suggestions []Suggestion
	err := db.Where("guild_id = ?", guildID).Find(&suggestions).Error
	return suggestions, err
}

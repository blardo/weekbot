package models

import (
	"strings"
	"weekbot-go/internal/config"

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
		panic(result.Error)
	} else {
		println("Suggestion created ", result.RowsAffected)
	}

	// Debug print to verify the saved record
	var savedSuggestion Suggestion
	db.First(&savedSuggestion, "content = ? AND guild_id = ?", content, guildID)
	println("Saved Suggestion: ", savedSuggestion.Content, savedSuggestion.GuildID, savedSuggestion.Updicks)

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
			println("Error finding suggestion: ", result.Error)
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
			println("Error updating suggestion: ", result.Error)
		}
	}

	suggestion.Updicks = updicks
	result = db.Save(&suggestion)
	if result.Error != nil {
		println("Error updating suggestion: ", result.Error)
	}

	println("Suggestion updated ", result.RowsAffected)

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

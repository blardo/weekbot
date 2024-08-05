package models

import (
	"errors"
	"fmt"

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

func UpdateSuggestion(db *gorm.DB, content string, guildID string, updicks int) error {
	var suggestion Suggestion
	result := db.First(&suggestion, "content = ? AND guild_id = ?", content, guildID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("suggestion not found")
		}
		return result.Error
	}

	suggestion.Updicks = updicks
	result = db.Save(&suggestion)
	if result.Error != nil {
		return result.Error
	}

	println("Suggestion updated ", result.RowsAffected)
	return nil
}

func GetMostRecentUnusedSuggestions(db *gorm.DB) []Suggestion {
	var suggestions []Suggestion
	db.Where("used = ? AND updicks >= ?", false, 3).Find(&suggestions)
	return suggestions
}

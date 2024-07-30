package models

import (
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
	return s
}

func GetMostRecentUnusedSuggestions(db *gorm.DB) []Suggestion {
	var suggestions []Suggestion
	db.Where("used = ? AND updicks >= ?", false, 3).Find(&suggestions)
	return suggestions
}

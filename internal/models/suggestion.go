package models

import (
	"gorm.io/gorm"
)

type Suggestion struct {
	gorm.Model
	Content string
	Used    bool
	PollID  uint
}

func NewSuggestion(db *gorm.DB, content string) *Suggestion {
	s := &Suggestion{
		Content: content,
		Used:    false,
	}
	db.Create(s)
	return s
}

func GetMostRecentUnusedSuggestions(db *gorm.DB) []Suggestion {
	var suggestions []Suggestion
	db.Where("used = ?", false).Find(&suggestions)
	return suggestions
}

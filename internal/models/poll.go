package models

import (
	"gorm.io/gorm"
)

// Poll is a struct that represents a poll
type Poll struct {
	gorm.Model
	Suggestions []Suggestion
	InProgress  bool
}

// StartPoll creates a new poll from the most recent suggestions
func StartPoll(bot *Bot) *Poll {
	current := GetCurrentPoll(bot.DB)
	if current != nil {
		return current
	}
	suggestions := GetMostRecentUnusedSuggestions(bot.DB)
	poll := &Poll{
		Suggestions: suggestions,
		InProgress:  true,
	}
	bot.DB.Create(poll)
	return poll
}

// IsComplete returns true if the poll is complete
func (p *Poll) IsComplete() bool {
	return !p.InProgress
}

// GetCurrentPoll returns the current poll from the database
func GetCurrentPoll(db *gorm.DB) *Poll {
	var poll Poll
	db.Where("in_progress = ?", true).First(&poll)
	if poll.ID == 0 {
		return nil
	}

	return &poll
}

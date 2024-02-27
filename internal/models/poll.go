package models

import "gorm.io/gorm"

// Poll is a struct that represents a poll
type Poll struct {
	gorm.Model
	Suggestions []Suggestion
	InProgress  bool
	UpNext      bool
}

// NewPoll creates a new poll from the most recent suggestions
func NewPoll(suggestions []Suggestion) *Poll {
	return &Poll{
		Suggestions: suggestions,
		InProgress:  true,
		UpNext:      true,
	}
}

// AddSuggestion adds a suggestion to the poll
func (p *Poll) AddSuggestion(suggestion Suggestion) {
	p.Suggestions = append(p.Suggestions, suggestion)
}

// IsComplete returns true if the poll is complete
func (p *Poll) IsComplete() bool {
	return !p.InProgress
}

// GetCurrentPoll returns the current poll from the database
func GetCurrentPoll(db gorm.DB) *Poll {
	var poll Poll
	db.Where("in_progress = ?", true).First(&poll)
	return &poll
}

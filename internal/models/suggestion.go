package models

import "gorm.io/gorm"

// Suggestion is a struct that represents a suggestion
type Suggestion struct {
	gorm.Model
	WeekName string
	Selected bool
	Updicks  int
	PollID   uint
}

// NewSuggestion creates a new suggestion
func NewSuggestion(weekName string) *Suggestion {
	return &Suggestion{
		WeekName: weekName,
		Selected: false,
		Updicks:  0, 
	}
}

package models

import "gorm.io/gorm"

// Suggestion is a struct that represents a suggestion
type Suggestion struct {
	gorm.Model
	weekName string
	selected bool
	updicks  int
	PollID   uint
}

// Add adds a suggestion to the list
func (s *Suggestion) NewSuggestion(weekName string) {
	s.weekName = weekName
	s.selected = false
	s.updicks = 0
}

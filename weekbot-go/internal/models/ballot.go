package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Ballot struct {
	gorm.Model
	PollID       uint
	VoterId      string
	FirstChoice  string `gorm:"column:first_choice"`  // Deprecated, kept for backward compatibility
	SecondChoice string `gorm:"column:second_choice"` // Deprecated, kept for backward compatibility
	ThirdChoice  string `gorm:"column:third_choice"` // Deprecated, kept for backward compatibility
	Choices      string `gorm:"type:text"`            // JSON array of ranked choices
	Date         time.Time
	Cast         bool
}

// GetSubmitterID returns the SubmitterID of the ballot
func (b *Ballot) GetSubmitterID() string {
	return b.VoterId
}

// SetSubmitterID sets the SubmitterID of the ballot
func (b *Ballot) SetSubmitterID(submitterID string) {
	b.VoterId = submitterID
}

// GetFirstChoice returns the FirstChoice of the ballot (for backward compatibility)
func (b *Ballot) GetFirstChoice() string {
	return b.GetChoiceAt(1)
}

// SetFirstChoice sets the FirstChoice of the ballot (for backward compatibility)
func (b *Ballot) SetFirstChoice(firstChoice string) {
	choices := b.GetChoices()
	if len(choices) > 0 {
		choices[0] = firstChoice
	} else {
		choices = []string{firstChoice}
	}
	b.SetChoices(choices)
}

// GetSecondChoice returns the SecondChoice of the ballot (for backward compatibility)
func (b *Ballot) GetSecondChoice() string {
	return b.GetChoiceAt(2)
}

// SetSecondChoice sets the SecondChoice of the ballot (for backward compatibility)
func (b *Ballot) SetSecondChoice(secondChoice string) {
	choices := b.GetChoices()
	for len(choices) < 2 {
		choices = append(choices, "")
	}
	choices[1] = secondChoice
	b.SetChoices(choices)
}

// GetThirdChoice returns the ThirdChoice of the ballot (for backward compatibility)
func (b *Ballot) GetThirdChoice() string {
	return b.GetChoiceAt(3)
}

// SetThirdChoice sets the ThirdChoice of the ballot (for backward compatibility)
func (b *Ballot) SetThirdChoice(thirdChoice string) {
	choices := b.GetChoices()
	for len(choices) < 3 {
		choices = append(choices, "")
	}
	choices[2] = thirdChoice
	b.SetChoices(choices)
}

// GetDate returns the Date of the ballot
func (b *Ballot) GetDate() time.Time {
	return b.Date
}

// SetDate sets the Date of the ballot
func (b *Ballot) SetDate(date time.Time) {
	b.Date = date
}

func (b *Ballot) SetCast(cast bool) {
	b.Cast = cast
}

// GetBallotByID gets a ballot by its ID
func GetBallotByVoterID(db *gorm.DB, voterID string) *Ballot {
	var ballot Ballot
	db.Where("voter_id = ?", voterID).First(&ballot)
	return &ballot
}

// GetAllBallots gets all ballots
func GetAllBallots(db *gorm.DB) []Ballot {
	var ballots []Ballot
	db.Find(&ballots)
	return ballots
}

// GetChoices returns the ranked choices as a slice
// If Choices is empty, migrates from legacy First/Second/ThirdChoice fields
func (b *Ballot) GetChoices() []string {
	if b.Choices != "" {
		var choices []string
		if err := json.Unmarshal([]byte(b.Choices), &choices); err == nil && len(choices) > 0 {
			return choices
		}
	}
	
	// Fallback: migrate from legacy fields
	var choices []string
	if b.FirstChoice != "" {
		choices = append(choices, b.FirstChoice)
	}
	if b.SecondChoice != "" {
		choices = append(choices, b.SecondChoice)
	}
	if b.ThirdChoice != "" {
		choices = append(choices, b.ThirdChoice)
	}
	
	// Auto-migrate: save to Choices field if we found legacy data
	if len(choices) > 0 && b.Choices == "" {
		b.SetChoices(choices)
	}
	
	return choices
}

// SetChoices sets the ranked choices from a slice
func (b *Ballot) SetChoices(choices []string) error {
	data, err := json.Marshal(choices)
	if err != nil {
		return err
	}
	b.Choices = string(data)
	return nil
}

// AddChoice adds a choice to the ranked list
func (b *Ballot) AddChoice(choice string) error {
	choices := b.GetChoices()
	// Check if already selected
	for _, c := range choices {
		if c == choice {
			return nil // Already selected, ignore
		}
	}
	choices = append(choices, choice)
	return b.SetChoices(choices)
}

// GetChoiceAt returns the choice at the given rank (1-indexed)
func (b *Ballot) GetChoiceAt(rank int) string {
	choices := b.GetChoices()
	if rank < 1 || rank > len(choices) {
		return ""
	}
	return choices[rank-1]
}

// Save saves the ballot to the database
func (b *Ballot) Save(db *gorm.DB) {
	db.Save(b)
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Ballot struct {
    gorm.Model
	PollID 	     uint
    VoterId  string
    FirstChoice  string
    SecondChoice string
    ThirdChoice  string
    Date         time.Time
}

// GetSubmitterID returns the SubmitterID of the ballot
func (b *Ballot) GetSubmitterID() string {
    return b.VoterId
}

// SetSubmitterID sets the SubmitterID of the ballot
func (b *Ballot) SetSubmitterID(submitterID string) {
    b.VoterId = submitterID
}

// GetFirstChoice returns the FirstChoice of the ballot
func (b *Ballot) GetFirstChoice() string {
    return b.FirstChoice
}

// SetFirstChoice sets the FirstChoice of the ballot
func (b *Ballot) SetFirstChoice(firstChoice string) {
    b.FirstChoice = firstChoice
}

// GetSecondChoice returns the SecondChoice of the ballot
func (b *Ballot) GetSecondChoice() string {
    return b.SecondChoice
}

// SetSecondChoice sets the SecondChoice of the ballot
func (b *Ballot) SetSecondChoice(secondChoice string) {
    b.SecondChoice = secondChoice
}

// GetThirdChoice returns the ThirdChoice of the ballot
func (b *Ballot) GetThirdChoice() string {
    return b.ThirdChoice
}

// SetThirdChoice sets the ThirdChoice of the ballot
func (b *Ballot) SetThirdChoice(thirdChoice string) {
    b.ThirdChoice = thirdChoice
}

// GetDate returns the Date of the ballot
func (b *Ballot) GetDate() time.Time {
    return b.Date
}

// SetDate sets the Date of the ballot
func (b *Ballot) SetDate(date time.Time) {
    b.Date = date
}

// GetBallotByID gets a ballot by its ID
func GetBallotByID(db *gorm.DB, id uint) *Ballot {
    var ballot Ballot
    db.First(&ballot, id)
    return &ballot
}

// GetAllBallots gets all ballots
func GetAllBallots(db *gorm.DB) []Ballot {
    var ballots []Ballot
    db.Find(&ballots)
    return ballots
}

// Save saves the ballot to the database
func (b *Ballot) Save(db *gorm.DB) {
    db.Save(b)
}
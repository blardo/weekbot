package models

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// Poll is a struct that represents a poll
type Poll struct {
    gorm.Model
    ID          uint `sql:"AUTO_INCREMENT" gorm:"primary_key"`
    Suggestions []Suggestion
    InProgress  bool
    IsComplete  bool
    Ballots     []Ballot `gorm:"foreignKey:PollID"`
    
}
// NewOrCurrentPoll creates a new poll or returns the current one
func NewOrCurrentPoll(bot *Bot) *Poll {
    poll := GetCurrentPoll(bot.DB)
    if poll != nil {
        println("Current poll found")
        return poll
    }

    suggestions := GetMostRecentUnusedSuggestions(bot.DB)
    if len(suggestions) == 0 {
        fmt.Println("No suggestions found")
        return nil
    }

    poll = &Poll{
        Suggestions: suggestions,
        InProgress:  true,
    }

    bot.DB.Create(poll)
    return poll
}

// GetCurrentPoll gets the current poll from the database
func GetCurrentPoll(db *gorm.DB) *Poll {
    var poll Poll
    db.Preload("Suggestions").Where("in_progress = ? and is_complete = ?", true, false).First(&poll)
    if poll.ID == 0 {
        fmt.Println("No current poll found")
        return nil
    }

    return &poll
}

// GetSelectOptions gets the select options for the poll
func (p *Poll) GetSelectOptions() []discordgo.SelectMenuOption {
    var options []discordgo.SelectMenuOption
    var filter []Suggestion
    for _, suggestion := range p.Suggestions {
        isDuplicate := false
        for _, f := range filter {
            if strings.EqualFold(suggestion.Content, f.Content){
                isDuplicate = true
                break
            }
        }

        if !isDuplicate {
            filter = append(filter, suggestion)
        }
    }
    // discords options limit is 25
    if len(filter) > 25 {
        filter = filter[:25]
    } 
    for _, suggestion := range filter {
        options = append(options, discordgo.SelectMenuOption{
            Label:   suggestion.Content,
            Value:   suggestion.Content,
            Default: false,
            Emoji:   discordgo.ComponentEmoji{Name: "📅"},
        })
    }

    return options
}

// IsVoter checks if a user is a voter
func (p *Poll) HasBallot(voterID string) bool {
    for _, voter := range p.Ballots {
        if voter.VoterId == voterID {
            return true
        }
    }
    return false
}



// GetBallots gets the ballots of the poll
func (p *Poll) GetBallots() []Ballot {
    return p.Ballots
}

// AddBallot adds a ballot to the poll
func (p *Poll) AddBallot(ballot Ballot) {
    p.Ballots = append(p.Ballots, ballot)
}

// EndPoll ends the poll
func (p *Poll) EndPoll() {
    p.InProgress = false
    p.IsComplete = true

    for _, suggestion := range p.Suggestions {
        suggestion.Used = true
    }
}

func (p *Poll) AddBallotForVoter(bot *Bot, userID string, ballot Ballot) {
    var voter Voter
    bot.DB.Where("user_id = ?", userID).First(&voter)

    if voter.UserID == "" {
        // Voter does not exist, create new voter
        voter = Voter{
            UserID: userID,
        }
        bot.DB.Create(&voter)
    }

    // Add the voter's ID to the ballot and save it
    ballot.VoterId = userID
    bot.DB.Create(&ballot)

    // Add the ballot to the poll
    p.AddBallot(ballot)
}
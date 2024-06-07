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
	Suggestions []Suggestion
	InProgress  bool
	IsComplete  bool
	VoterIDs    []string `gorm:"type:text[]"`
}

func NewOrCurrentPoll(bot *Bot) *Poll {
	poll := GetCurrentPoll(bot.DB)
	if poll != nil {
		return poll
	}

	suggestions := GetMostRecentUnusedSuggestions(bot.DB)
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found")
		return nil
	}

	for i := range suggestions {
		suggestions[i].Used = true
		bot.DB.Save(&suggestions[i])
	}

	voterIDs := getVoterIDs(poll)
	if voterIDs == nil {
		voterIDs = []string{}
	}

	poll = &Poll{
		Suggestions: suggestions,
		InProgress:  true,
		VoterIDs:   voterIDs,
	}

	bot.DB.Create(poll)
	return poll
}

func GetCurrentPoll(db *gorm.DB) *Poll {
	var poll Poll
	db.Preload("Suggestions").Where("in_progress = ? and is_complete = ?", true, false).First(&poll)
	if poll.ID == 0 {
		fmt.Println("No current poll found")
		return nil
	}

	return &poll
}

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

func (p *Poll) AddVoter(voterID string) {
	p.VoterIDs = append(p.VoterIDs, voterID)
}

func (p *Poll) IsVoter(voterID string) bool {
	for _, id := range p.VoterIDs {
		if id == voterID {
			return true
		}
	}
	return false
}

func getVoterIDs(poll *Poll) []string {
	return poll.VoterIDs
}


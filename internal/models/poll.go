package models

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// Poll is a struct that represents a poll
type Poll struct {
	gorm.Model
	Suggestions []Suggestion
	InProgress  bool
	IsComplete  bool
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

	poll = &Poll{
		Suggestions: suggestions,
		InProgress:  true,
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
	for _, suggestion := range p.Suggestions {
		options = append(options, discordgo.SelectMenuOption{
			Label:   suggestion.Content,
			Value:   suggestion.Content,
			Default: false,
			Emoji:   discordgo.ComponentEmoji{Name: "📅"},
		})
	}

	return options
}

package commands

import (
	"fmt"
	"weekbot/models"

	"github.com/bwmarrin/discordgo"
)

// HandleWeekSuggestion adds a week suggestion to the list
func HandleWeekSuggestion(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check that the message is not from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is just the word week, ignore it
	if m.Content == "week" || m.Content == "Week" {
		return
	}

	// Create a new suggestion
	suggestion := &models.Suggestion{}
	suggestion.NewSuggestion(m.Content)
	fmt.Println("Adding suggestion:", suggestion)
	s.ChannelMessageSend(m.ChannelID, "Week suggestion added: "+m.Content)
}

package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ParseCommand parses a command from a message
func ParseCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// React to only messages not sent by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Check if the message is a command
	if m.Content[0] != '/' {
		return
	}
	// Split the message into its parts
	parts := strings.Split(m.Content, " ")
	// Get the command name
	command := parts[0][1:]

	// Handle the command
	switch command {
	case "ping":
		HandlePing(s, m)
	case "poll":
		HandleWeekPoll(s, m)
	}
}

package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// ParseCommand parses a command from a message
func ParseInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData().Name
	switch command {
	case "ping":
		HandlePing(s, i)
	case "poll":
		HandleWeekPoll(s, i)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func ParseChatCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// React to only messages not sent by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	// TODO: Parse the message into a command
}

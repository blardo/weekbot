package commands

import (
	"github.com/bwmarrin/discordgo"
)

// HandleWeekPoll handles the /poll command
func HandleWeekPoll(s *discordgo.Session, m *discordgo.MessageCreate) {
	// React to only messages not sent by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Send a message to the channel
	s.ChannelMessageSend(m.ChannelID, "Starting the Week Name Poll")
}

func GetPollCommand() *discordgo.ApplicationCommand {
	pollCommand := &discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Create a poll",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "question",
				Description: "The question to ask",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "options",
				Description: "The options to choose from",
				Required:    true,
			},
		},
	}
	return pollCommand
}

package commands

import (
	"github.com/bwmarrin/discordgo"
)

// HandlePing handles the /ping command
func HandlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	// React to only messages not sent by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func GetPingCommand() *discordgo.ApplicationCommand {
	pingCommand := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping the bot",
		Type:        discordgo.ChatApplicationCommand,
	}
	return pingCommand
}

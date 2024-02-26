package commands

import (
	"github.com/bwmarrin/discordgo"
)

// HandlePing handles the /ping command
func HandlePing(s *discordgo.Session, m *discordgo.InteractionCreate) {
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

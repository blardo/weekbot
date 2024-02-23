package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// DiscordService is your primary interface to Discord from the bot
type DiscordService struct {
	session *discordgo.Session
}

func NewDiscordService(token string) (*DiscordService, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		fmt.Println("Connected to Discord as", m.User.Username)
	})

	return &DiscordService{session}, nil
}

// Connect connects to Discord
func (d *DiscordService) Connect() error {
	return d.session.Open()
}

// Disconnect disconnects from Discord
func (d *DiscordService) Disconnect() error {
	return d.session.Close()
}

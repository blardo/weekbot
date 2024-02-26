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
		fmt.Println("Invite URL: https://discordapp.com/oauth2/authorize?client_id=" + m.User.ID + "&scope=bot&permissions=0")
		fmt.Println("Currently on", len(m.Guilds), "servers")
		// Print out the servers we're on
		for _, guild := range m.Guilds {
			fmt.Println("  ", guild.Name, "(", guild.ID, ")")
		}
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

// SendMessage sends a message to a channel
func (d *DiscordService) SendMessage(channelID, message string) error {
	_, err := d.session.ChannelMessageSend(channelID, message)
	return err
}

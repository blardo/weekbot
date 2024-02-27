package router

import (
	"fmt"
	"strings"
	"weekbot/internal/handlers/commands"
	"weekbot/internal/services/discord"

	"github.com/bwmarrin/discordgo"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

// Setup adds the bot's handlers to the Discord client
func (r *Router) Setup(ds *discord.DiscordService) {
	ds.AddHandler(ParseInteraction)
	ds.AddHandler(ParseChatCommand)
	ds.AddHandler(HandleReactions)
}

// ParseCommand parses a command from a message
func ParseInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.Type

	switch command {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "ping":
			s.ChannelMessageSend(i.ChannelID, "Pong!")
		case "poll":
			commands.HandleWeekPoll(s, i)
		default:
			fmt.Println("Unknown command:", i.ApplicationCommandData().Name)
		}
	default: // Ignore other types of interactions
		return
	}
}

func ParseChatCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// React to only messages not sent by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message ends in the word week, add it to the list of suggestions for the poll
	message := strings.Split(m.Content, " ")
	if message[len(message)-1] == "Week" || message[len(message)-1] == "week" {
		commands.HandleWeekSuggestion(s, m)

	}
}

func HandleReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// React to only messages not sent by the bot
	if r.UserID == s.State.User.ID {
		return
	}

	if r.Emoji.Name == "bd" {
		s.MessageReactionAdd(r.ChannelID, r.MessageID, "👍")
	}
}

package handlers

import (
	"fmt"
	"strings"

	actions "weekbot/internal/actions"
	commands "weekbot/internal/commands"

	"github.com/bwmarrin/discordgo"
)

// ParseCommand parses a command from a message
func ParseInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.Type

	switch command {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "ping":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Pong!"},
			})
		case "poll":
			commands.HandleWeekPoll(s, i)
		case "endpoll":
			commands.HandleEndPoll(s, i)
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
	acceptableWeeks := []string{"week", "week.", "week!", "week?"} // move to constants file
	for _, week := range acceptableWeeks {
		if strings.ToLower(message[len(message)-1]) == week {
			actions.HandleWeekSuggestion(s, m)
			break
		}
	}
}

func HandleReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// React to only messages not sent by the bot
	if r.UserID == s.State.User.ID {
		return
	}

	println(r.Emoji.Name)
	if r.Emoji.Name == "bd" { // bd 👍
	
		s.MessageReactionAdd(r.ChannelID, r.MessageID, "👍")
	}
}

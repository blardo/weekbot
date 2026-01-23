package handlers

import (
	"fmt"
	"strings"
	"weekbot-go/internal/config"
	"weekbot-go/internal/logger"
	"weekbot-go/internal/models"

	"weekbot-go/internal/actions"
	"weekbot-go/internal/commands"

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
		case "listweeks":
			commands.HandleListWeeks(s, i)
		default:
			logger.Warn("Unknown command", "command", i.ApplicationCommandData().Name)
		}
	case discordgo.InteractionMessageComponent:
		commands.HandlePollComponent(s, i)
	default: // Ignore other types of interactions
		return
	}
}

func ParseChatCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Early return if message is from bot or has no author
	if m.Author == nil || m.Author.ID == "" || m.Author.ID == s.State.User.ID {
		return
	}

	// Use cached channel lookup
	channelCache := GetChannelCache()
	channelID, err := channelCache.GetChannelIDByName(m.GuildID, config.WeekNameChannelName)
	if err != nil {
		// Silently ignore if channel not found (bot might not be in that guild yet)
		return
	}

	// If the message isn't in the week-name channel, ignore it
	if m.ChannelID != channelID {
		return
	}

	// Check if message should be processed (deduplication)
	tracker := GetMessageTracker()
	if !tracker.ShouldProcessMessage(m.ID) {
		return
	}

	// If the message ends in the word week, add it to the list of suggestions for the poll
	message := strings.Split(m.Content, " ")
	for _, week := range config.AcceptableWeekSuffixes {
		if len(message) > 0 && strings.ToLower(message[len(message)-1]) == week {
			actions.HandleWeekSuggestion(s, m)
			break
		}
	}
}

func HandleReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Early return if reaction is from bot
	if r.UserID == "" || r.UserID == s.State.User.ID {
		return
	}

	// Check if reaction should be processed (deduplication)
	tracker := GetMessageTracker()
	emoji := r.Emoji.Name
	if r.Emoji.ID != "" {
		emoji = fmt.Sprintf(":%s:%s", r.Emoji.Name, r.Emoji.ID)
	}
	
	if !tracker.ShouldProcessReaction(r.ChannelID, r.MessageID, r.UserID, emoji) {
		return
	}

	bot := models.GetBot(r.GuildID)
	if bot == nil {
		return
	}

	m, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		logger.Error("Error retrieving message", "error", err, "channel_id", r.ChannelID, "message_id", r.MessageID)
		return
	}
	
	reaction, err := s.MessageReactions(r.ChannelID, r.MessageID, emoji, 100, "", "")
	if err != nil {
		logger.Error("Error getting reactions", "error", err, "channel_id", r.ChannelID, "message_id", r.MessageID, "emoji", emoji)
		return
	}

	if r.Emoji.Name == config.QualifyingEmoji && len(reaction) >= config.MinUpdicksToQualify {
		models.UpdateSuggestion(bot.DB, m.Content, r.GuildID, len(reaction))
		s.MessageReactionAdd(r.ChannelID, r.MessageID, config.ConfirmationEmoji)
	}
}

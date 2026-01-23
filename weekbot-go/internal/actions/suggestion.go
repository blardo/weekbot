package actions

import (
	"log"
	"strings"
	"weekbot-go/internal/models"

	"github.com/bwmarrin/discordgo"
)

// HandleWeekSuggestion adds a week suggestion to the list
func HandleWeekSuggestion(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Additional defensive check (though ParseChatCommand should have already filtered)
	if m.Author == nil || m.Author.ID == "" || m.Author.ID == s.State.User.ID {
		log.Printf("Skipping message from bot or nil author (message ID: %s)", m.ID)
		return
	}

	// Note: Message deduplication is now handled in ParseChatCommand
	// This function assumes it's only called for unique messages

	// If the message is just the word week, ignore it
	if m.Content == "week" || m.Content == "Week" {
		return
	}

	suggestion := strings.TrimSpace(m.Content)
	normalized := models.NormalizeSuggestionContent(suggestion)
	log.Printf("Processing suggestion: %s (normalized: %s, message ID: %s)", suggestion, normalized, m.ID)

	// Create a new suggestion from the message
	bot := models.GetBot(m.GuildID)
	if bot == nil {
		log.Printf("Bot not found for guild %s", m.GuildID)
		return
	}

	existing, exists := models.FindActiveSuggestionByContent(bot.DB, suggestion, bot.GuildID)
	if exists {
		log.Printf("Suggestion already exists: %s (normalized: %s, existing ID: %d, existing content: %s, message ID: %s, author: %s)", 
			suggestion, normalized, existing.ID, existing.Content, m.ID, m.Author.ID)
		_, err := s.ChannelMessageSend(m.ChannelID, "Week suggestion already exists: "+suggestion)
		if err != nil {
			log.Printf("Error sending 'already exists' message: %v", err)
		} else {
			log.Printf("Sent 'already exists' message for suggestion: %s", suggestion)
		}
		return
	}

	models.NewSuggestion(bot.DB, suggestion, bot.GuildID)
	log.Printf("Added new suggestion: %s (message ID: %s)", suggestion, m.ID)

	_, err := s.ChannelMessageSend(m.ChannelID, "Week suggestion added: "+m.Content+". Weekbot will add a 👍 when this suggestion has enough votes to be added to the next poll.")
	if err != nil {
		log.Printf("Error sending 'added' message: %v", err)
	} else {
		log.Printf("Sent 'added' message for suggestion: %s", suggestion)
	}
	list, err := models.GetAllSuggestions(bot.DB, m.GuildID)
	if err != nil {
		log.Printf("Error getting suggestions: %v", err)
	} else {
		for _, suggestion := range list {
			log.Printf("Existing suggestion: %s", suggestion.Content)
		}
	}
}

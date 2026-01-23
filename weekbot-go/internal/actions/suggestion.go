package actions

import (
	"log"
	"strings"
	"weekbot-go/internal/models"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
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
	log.Printf("=== HandleWeekSuggestion called === Processing suggestion: %s (normalized: %s, message ID: %s, author: %s)", 
		suggestion, normalized, m.ID, m.Author.ID)

	// Create a new suggestion from the message
	bot := models.GetBot(m.GuildID)
	if bot == nil {
		log.Printf("Bot not found for guild %s", m.GuildID)
		return
	}

	// Use a transaction to atomically check and create the suggestion
	// This prevents race conditions where the same message is processed twice
	var newSuggestion *models.Suggestion
	var alreadyExists bool
	
	err := bot.DB.Transaction(func(tx *gorm.DB) error {
		// Check if suggestion already exists
		existing, exists := models.FindActiveSuggestionByContent(tx, suggestion, bot.GuildID)
		if exists {
			alreadyExists = true
			log.Printf("Suggestion already exists in transaction: %s (normalized: %s, existing ID: %d, existing content: %s, message ID: %s, author: %s)", 
				suggestion, normalized, existing.ID, existing.Content, m.ID, m.Author.ID)
			return nil // Return nil to commit (we're just checking)
		}
		
		// Create the suggestion within the transaction
		content := strings.TrimSpace(suggestion)
		sugg := &models.Suggestion{
			Content: content,
			GuildID: bot.GuildID,
			Used:    false,
			Updicks: 0,
		}
		if err := tx.Create(sugg).Error; err != nil {
			log.Printf("Error creating suggestion in transaction: %v", err)
			return err
		}
		newSuggestion = sugg
		log.Printf("Created suggestion in transaction: %s (ID: %d, message ID: %s)", suggestion, sugg.ID, m.ID)
		return nil
	})
	
	if err != nil {
		log.Printf("Transaction error: %v", err)
		return
	}
	
	if alreadyExists {
		_, err := s.ChannelMessageSend(m.ChannelID, "Week suggestion already exists: "+suggestion)
		if err != nil {
			log.Printf("Error sending 'already exists' message: %v", err)
		} else {
			log.Printf("Sent 'already exists' message for suggestion: %s", suggestion)
		}
		return
	}
	
	if newSuggestion == nil {
		log.Printf("Unexpected: newSuggestion is nil after transaction")
		return
	}
	
	// Send success message
	_, err = s.ChannelMessageSend(m.ChannelID, "Week suggestion added: "+m.Content+". Weekbot will add a 👍 when this suggestion has enough votes to be added to the next poll.")
	if err != nil {
		log.Printf("Error sending 'added' message: %v", err)
	} else {
		log.Printf("Sent 'added' message for suggestion: %s (ID: %d)", suggestion, newSuggestion.ID)
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

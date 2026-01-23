package commands

import (
	"fmt"
	"strings"
	"weekbot-go/internal/logger"
	"weekbot-go/internal/models"

	"github.com/bwmarrin/discordgo"
)

// HandleListWeeks lists all current week suggestions for debugging
func HandleListWeeks(s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot := models.GetBot(i.GuildID)
	if bot == nil {
		respondEphemeral(s, i, "Bot not available in this server.")
		return
	}

	var response strings.Builder
	response.WriteString("**Current Week Suggestions:**\n\n")

	// Get current poll if one exists
	currentPoll := models.GetCurrentPoll(bot.DB)
	if currentPoll != nil {
		response.WriteString(fmt.Sprintf("**Active Poll (ID: %d):**\n", currentPoll.ID))
		if len(currentPoll.Suggestions) > 0 {
			for idx, suggestion := range currentPoll.Suggestions {
				status := "❌ Used"
				if !suggestion.Used {
					status = "✅ Active"
				}
				response.WriteString(fmt.Sprintf("%d. %s (%d 👍) - %s\n", 
					idx+1, suggestion.Content, suggestion.Updicks, status))
			}
		} else {
			response.WriteString("No suggestions in current poll.\n")
		}
		response.WriteString("\n")
	} else {
		response.WriteString("**No active poll currently.**\n\n")
	}

	// Get all unused suggestions
	unusedSuggestions := models.GetMostRecentUnusedSuggestions(bot.DB)
	response.WriteString(fmt.Sprintf("**Unused Suggestions (%d):**\n", len(unusedSuggestions)))
	if len(unusedSuggestions) > 0 {
		for idx, suggestion := range unusedSuggestions {
			response.WriteString(fmt.Sprintf("%d. %s (%d 👍)\n", 
				idx+1, suggestion.Content, suggestion.Updicks))
		}
	} else {
		response.WriteString("No unused suggestions available.\n")
	}

	// Get all suggestions for stats
	allSuggestions, err := models.GetAllSuggestions(bot.DB, i.GuildID)
	if err == nil {
		usedCount := 0
		for _, s := range allSuggestions {
			if s.Used {
				usedCount++
			}
		}
		response.WriteString(fmt.Sprintf("\n**Stats:** %d total, %d used, %d unused\n", 
			len(allSuggestions), usedCount, len(allSuggestions)-usedCount))
	}

	// Discord has a 2000 character limit for messages
	content := response.String()
	if len(content) > 2000 {
		content = content[:1997] + "..."
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error("Error responding to listweeks command", "error", err)
	}
}

// GetListWeeksCommand returns the listweeks command
func GetListWeeksCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "listweeks",
		Description: "List all current week suggestions (debug command)",
		Type:        discordgo.ChatApplicationCommand,
	}
}

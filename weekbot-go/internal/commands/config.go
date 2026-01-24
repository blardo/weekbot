package commands

import (
	"fmt"
	"strconv"
	"weekbot-go/internal/config"
	"weekbot-go/internal/logger"
	"weekbot-go/internal/models"

	"github.com/bwmarrin/discordgo"
)

// GetConfigCommand returns the /config slash command definition
func GetConfigCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "config",
		Description: "View or modify server poll settings",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "view",
				Description: "View current server settings",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "set",
				Description: "Update a server setting (admin only)",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "setting",
						Description: "The setting to update",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "min-suggestions-to-start-poll", Value: "min-suggestions-to-start-poll"},
							{Name: "min-updicks", Value: "min-updicks"},
							{Name: "min-ballots-to-end-poll", Value: "min-ballots-to-end-poll"},
							{Name: "reaction-emoji", Value: "reaction-emoji"},
						},
					},
					{
						Name:        "value",
						Description: "The new value (number for thresholds, emoji name for reaction-emoji)",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
		},
	}
}

func floatPtr(f float64) *float64 {
	return &f
}

// HandleConfigCommand handles the /config command
func HandleConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot := models.GetBot(i.GuildID)
	if bot == nil {
		respondEphemeral(s, i, "Bot is not available in this server.")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		respondEphemeral(s, i, "Please specify a subcommand: view or set")
		return
	}

	switch options[0].Name {
	case "view":
		handleConfigView(s, i, bot)
	case "set":
		handleConfigSet(s, i, bot)
	default:
		respondEphemeral(s, i, "Unknown subcommand")
	}
}

// handleConfigView shows current server settings
func handleConfigView(s *discordgo.Session, i *discordgo.InteractionCreate, bot *models.Bot) {
	serverConfig := models.GetServerConfig(bot.DB, i.GuildID)

	// Get current effective values
	minSuggestions := config.MinSuggestionsToStartPoll(bot.DB, i.GuildID)
	minUpdicks := config.MinUpdicksToQualify(bot.DB, i.GuildID)
	minBallots := config.MinBallotsToEndPoll(bot.DB, i.GuildID)
	reactionEmoji := config.QualifyingEmojiForGuild(bot.DB, i.GuildID)

	// Build the response showing current values and whether they're custom or default
	var minSuggestionsSource, minUpdicksSource, minBallotsSource, emojiSource string
	if serverConfig != nil && serverConfig.MinSuggestionsToStart != nil {
		minSuggestionsSource = "(custom)"
	} else {
		minSuggestionsSource = "(default)"
	}
	if serverConfig != nil && serverConfig.MinUpdicksToQualify != nil {
		minUpdicksSource = "(custom)"
	} else {
		minUpdicksSource = "(default)"
	}
	if serverConfig != nil && serverConfig.MinBallotsToEndPoll != nil {
		minBallotsSource = "(custom)"
	} else {
		minBallotsSource = "(default)"
	}
	if serverConfig != nil && serverConfig.ReactionEmoji != nil {
		emojiSource = "(custom)"
	} else {
		emojiSource = "(default)"
	}

	content := fmt.Sprintf("**Server Poll Settings**\n\n"+
		"`min-suggestions-to-start-poll`: %d %s\n"+
		"  Minimum suggestions needed to start a poll\n\n"+
		"`min-updicks`: %d %s\n"+
		"  Minimum reactions needed for a suggestion to qualify\n\n"+
		"`min-ballots-to-end-poll`: %d %s\n"+
		"  Minimum votes needed to end a poll\n\n"+
		"`reaction-emoji`: %s %s\n"+
		"  Emoji that marks suggestions for the poll\n\n"+
		"_Use `/config set <setting> <value>` to customize (admin only)_",
		minSuggestions, minSuggestionsSource,
		minUpdicks, minUpdicksSource,
		minBallots, minBallotsSource,
		reactionEmoji, emojiSource)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error("Error responding to config view", "error", err, "guild_id", i.GuildID)
	}
}

// handleConfigSet updates a server setting (admin only)
func handleConfigSet(s *discordgo.Session, i *discordgo.InteractionCreate, bot *models.Bot) {
	// Check for admin permission
	if !hasAdminPermission(i) {
		respondEphemeral(s, i, "You need Administrator or Manage Server permission to change settings.")
		return
	}

	options := i.ApplicationCommandData().Options[0].Options
	if len(options) < 2 {
		respondEphemeral(s, i, "Please specify both setting and value.")
		return
	}

	var setting string
	var value string
	for _, opt := range options {
		switch opt.Name {
		case "setting":
			setting = opt.StringValue()
		case "value":
			value = opt.StringValue()
		}
	}

	if setting == "" {
		respondEphemeral(s, i, "Invalid setting specified.")
		return
	}

	// Handle string settings (emoji)
	if setting == "reaction-emoji" {
		err := models.UpdateServerConfigStringField(bot.DB, i.GuildID, setting, value)
		if err != nil {
			logger.Error("Error updating server config", "error", err, "guild_id", i.GuildID, "setting", setting)
			respondEphemeral(s, i, "Failed to update setting. Please try again.")
			return
		}
		content := fmt.Sprintf("Updated `%s` to **%s**", setting, value)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Handle integer settings
	intValue, err := strconv.Atoi(value)
	if err != nil {
		respondEphemeral(s, i, "Invalid value. Please enter a number for this setting.")
		return
	}

	if intValue < 1 || intValue > 100 {
		respondEphemeral(s, i, "Value must be between 1 and 100.")
		return
	}

	err = models.UpdateServerConfigField(bot.DB, i.GuildID, setting, intValue)
	if err != nil {
		logger.Error("Error updating server config", "error", err, "guild_id", i.GuildID, "setting", setting)
		respondEphemeral(s, i, "Failed to update setting. Please try again.")
		return
	}

	content := fmt.Sprintf("Updated `%s` to **%s**", setting, value)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error("Error responding to config set", "error", err, "guild_id", i.GuildID)
	}
}

// hasAdminPermission checks if the user has Administrator or Manage Server permission
func hasAdminPermission(i *discordgo.InteractionCreate) bool {
	if i.Member == nil {
		return false
	}
	// Check for Administrator or Manage Server permission using Member.Permissions
	// which is populated by Discord in interaction events
	return i.Member.Permissions&discordgo.PermissionAdministrator != 0 ||
		i.Member.Permissions&discordgo.PermissionManageServer != 0
}

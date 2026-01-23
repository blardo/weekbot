package commands

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"
	"weekbot-go/internal/config"
	"weekbot-go/internal/logger"
	"weekbot-go/internal/models"
	"weekbot-go/internal/services"
	"weekbot-go/internal/services/imagegen"

	"github.com/bwmarrin/discordgo"
)

// HandleWeekPoll handles the /poll command
func HandleWeekPoll(s *discordgo.Session, m *discordgo.InteractionCreate) {
	bot := models.GetBot(m.GuildID)

	poll := models.NewOrCurrentPoll(bot)

	if poll == nil {
		logger.Warn("Poll is nil, cannot start poll", "guild_id", m.GuildID)
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Not Enough Suggestions to Start Poll",
			},
		})
		return
	}

	/// Single Use DB Functions Here ----- /////

	// // Get all ballots
	// ballots := models.GetAllBallots(bot.DB)

	// // // delete all ballots
	// for _, ballot := range ballots {
	// 	bot.DB.Delete(&ballot)
	// 	println("Ballot: " + ballot.VoterId + " Deleted")
	// }

	// end poll

	//poll.EndPoll()
	// Delete all balllots
	//

	////// -------------- ////////////

	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Poll has started @everyone",
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							CustomID: "poll_button",
							Label:    "Vote Here",
							Style:    discordgo.PrimaryButton,
							Emoji: &discordgo.ComponentEmoji{
								Name: "🗳️",
							},
						},
					},
				},
			},
		},
	})
}

func HandlePollComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	switch i.MessageComponentData().CustomID {
	case "poll_button":
		handlePollButton(s, i)
	case "rank_choice":
		handleRankChoice(s, i)
	case "first_choice", "second_choice", "third_choice":
		// Legacy support - redirect to new handler
		handleRankChoice(s, i)
	case "submit_button":
		handlePollSubmit(s, i)
	default:
		return
	}
}

func handlePollButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot := models.GetBot(i.GuildID)
	if bot == nil {
		respondEphemeral(s, i, "Poll is not available in this server.")
		return
	}

	userID, ok := interactionUserID(i)
	if !ok {
		respondEphemeral(s, i, "Unable to identify your user.")
		return
	}

	poll := models.NewOrCurrentPoll(bot)
	if poll == nil {
		respondEphemeral(s, i, "No active poll right now.")
		return
	}

	if poll.BallotCast(userID) {
		respondEphemeral(s, i, "You have already voted.")
		return
	}

	if !poll.HasBallot(userID) {
		ballot := models.Ballot{
			VoterId: userID,
			PollID:  poll.ID,
			Date:    time.Now(),
			Cast:    false,
		}
		ballot.SetChoices([]string{}) // Initialize empty choices array
		poll.AddBallotForVoter(bot, ballot)
	}

	// Get current ballot to show existing selections (must be for this poll)
	ballot := models.GetBallotByVoterIDAndPollID(bot.DB, userID, poll.ID)
	if ballot.ID == 0 {
		// Ballot should have been created above, but if not, create it now
		ballot = &models.Ballot{
			VoterId: userID,
			PollID:  poll.ID,
			Date:    time.Now(),
			Cast:    false,
		}
		ballot.SetChoices([]string{})
		poll.AddBallotForVoter(bot, *ballot)
		// Reload to get the created ballot
		ballot = models.GetBallotByVoterIDAndPollID(bot.DB, userID, poll.ID)
	}
	selectedChoices := ballot.GetChoices()
	remainingOptions := poll.GetSelectOptionsExcluding(selectedChoices)

	// Build content showing current selections
	content := "**Rank your week name preferences:**\n\n"
	if len(selectedChoices) > 0 {
		content += "**Your current ranking:**\n"
		for i, choice := range selectedChoices {
			content += fmt.Sprintf("%d. %s\n", i+1, choice)
		}
		content += "\n"
	}
	
	if len(remainingOptions) == 0 {
		content += "✅ You've ranked all available options!"
	} else {
		content += fmt.Sprintf("Select your #%d choice:", len(selectedChoices)+1)
	}

	var components []discordgo.MessageComponent
	
	// Only show select menu if there are remaining options
	if len(remainingOptions) > 0 {
		components = append(components, &discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "rank_choice",
					Placeholder: fmt.Sprintf("Select choice #%d", len(selectedChoices)+1),
					Options:     remainingOptions,
				},
			},
		})
	}

	// Always show submit button
	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				CustomID: "submit_button",
				Label:    "Submit Vote",
				Style:    discordgo.PrimaryButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: "🗳️",
				},
			},
		},
	})

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: components,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func handleRankChoice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot := models.GetBot(i.GuildID)
	if bot == nil {
		respondEphemeral(s, i, "Poll is not available in this server.")
		return
	}

	userID, ok := interactionUserID(i)
	if !ok {
		respondEphemeral(s, i, "Unable to identify your user.")
		return
	}

	poll := models.GetCurrentPoll(bot.DB)
	if poll == nil {
		respondEphemeral(s, i, "No active poll right now.")
		return
	}

	ballot := models.GetBallotByVoterIDAndPollID(bot.DB, userID, poll.ID)
	if ballot.ID == 0 {
		respondEphemeral(s, i, "Please click \"Vote Here\" first.")
		return
	}

	if len(i.MessageComponentData().Values) == 0 {
		respondEphemeral(s, i, "No selection made.")
		return
	}

	selectedValue := i.MessageComponentData().Values[0]
	
	// Add the choice to the ballot
	if err := ballot.AddChoice(selectedValue); err != nil {
		logger.Error("Error adding choice to ballot", "error", err, "user_id", userID)
		respondEphemeral(s, i, "Error saving your choice. Please try again.")
		return
	}

	bot.DB.Save(ballot)

	// Get updated choices and remaining options
	selectedChoices := ballot.GetChoices()
	remainingOptions := poll.GetSelectOptionsExcluding(selectedChoices)

	// Build updated content
	content := "**Rank your week name preferences:**\n\n"
	if len(selectedChoices) > 0 {
		content += "**Your current ranking:**\n"
		for i, choice := range selectedChoices {
			content += fmt.Sprintf("%d. %s\n", i+1, choice)
		}
		content += "\n"
	}
	
	if len(remainingOptions) == 0 {
		content += "✅ You've ranked all available options! Click Submit to finalize your vote."
	} else {
		content += fmt.Sprintf("Select your #%d choice:", len(selectedChoices)+1)
	}

	var components []discordgo.MessageComponent
	
	// Only show select menu if there are remaining options
	if len(remainingOptions) > 0 {
		components = append(components, &discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "rank_choice",
					Placeholder: fmt.Sprintf("Select choice #%d", len(selectedChoices)+1),
					Options:     remainingOptions,
				},
			},
		})
	}

	// Always show submit button
	components = append(components, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				CustomID: "submit_button",
				Label:    "Submit Vote",
				Style:    discordgo.PrimaryButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: "🗳️",
				},
			},
		},
	})

	// Update the message with new components
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: components,
		},
	})
	if err != nil {
		logger.Error("Error updating interaction", "error", err, "user_id", userID)
	}
}

func handlePollSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot := models.GetBot(i.GuildID)
	if bot == nil {
		respondEphemeral(s, i, "Poll is not available in this server.")
		return
	}

	userID, ok := interactionUserID(i)
	if !ok {
		respondEphemeral(s, i, "Unable to identify your user.")
		return
	}

	poll := models.GetCurrentPoll(bot.DB)
	if poll == nil {
		respondEphemeral(s, i, "No active poll right now.")
		return
	}

	ballot := models.GetBallotByVoterIDAndPollID(bot.DB, userID, poll.ID)
	if ballot.ID == 0 {
		respondEphemeral(s, i, "Please click \"Vote Here\" first.")
		return
	}

	choices := ballot.GetChoices()
	if len(choices) == 0 {
		respondEphemeral(s, i, "Please rank at least one option before submitting.")
		return
	}

	ballot.Cast = true
	bot.DB.Save(ballot)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Your vote has been submitted",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func interactionUserID(i *discordgo.InteractionCreate) (string, bool) {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID, true
	}
	if i.User != nil {
		return i.User.ID, true
	}
	return "", false
}

func respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func HandleEndPoll(s *discordgo.Session, m *discordgo.InteractionCreate) {
	bot := models.GetBot(m.GuildID)

	poll := models.GetCurrentPoll(bot.DB)
	eligibleBallots := 0

	if poll == nil {
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No poll is currently in progress",
			},
		})
		return
	}
	for _, ballot := range poll.Ballots {
		if ballot.Cast {
			eligibleBallots++
		}
	}
	if eligibleBallots < config.MinBallotsToEndPoll {
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have enough votes to end the poll",
			},
		})
		return
	}

	// change server name

	newName := poll.PerformRankedChoiceVoting()

	logger.Info("Ranked choice voting completed", "winner", newName, "poll_id", poll.ID, "guild_id", m.GuildID)
	
	// Generate image for the new week name
	config := services.GetConfig()
	imageGen := imagegen.NewImageGenerator(config.OpenAIAPIKey)
	
	var iconData string
	if imageGen != nil {
		logger.Info("Generating image for week name", "week_name", newName, "guild_id", m.GuildID)
		imageBytes, err := imageGen.GenerateImageForWeek(newName)
		if err != nil {
			logger.Error("Error generating image", "error", err, "week_name", newName, "guild_id", m.GuildID)
			// Continue without image if generation fails
		} else {
			// Convert image to base64 data URI for Discord
			iconData = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(imageBytes))
			logger.Info("Image generated successfully", "size_bytes", len(imageBytes), "week_name", newName, "guild_id", m.GuildID)
		}
	} else {
		logger.Debug("Image generation not configured (no OpenAI API key)", "guild_id", m.GuildID)
	}
	
	// Update server name and icon
	guildParams := &discordgo.GuildParams{
		Name: newName,
	}
	if iconData != "" {
		guildParams.Icon = iconData
	}
	
	_, err := s.GuildEdit(m.GuildID, guildParams)
	if err != nil {
		logger.Error("Error changing server name/icon", "error", err, "new_name", newName, "guild_id", m.GuildID)
		return
	}
	logger.Info("Server name changed", "new_name", newName, "icon_updated", iconData != "", "guild_id", m.GuildID)

	// end poll (this will mark only the winning suggestion as used)
	poll.EndPoll(bot.DB, newName)
	bot.DB.Save(poll)

	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The poll has been ended",
		},
	})

	// end listener for select menu
}

func GetPollCommand() *discordgo.ApplicationCommand {
	pollCommand := &discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Run the Week Name Poll",
		Type:        discordgo.ChatApplicationCommand,
	}
	return pollCommand
}

func EndPollCommand() *discordgo.ApplicationCommand {
	pollCommand := &discordgo.ApplicationCommand{
		Name:        "endpoll",
		Description: "End the Week Name Poll",
		Type:        discordgo.ChatApplicationCommand,
	}
	return pollCommand
}

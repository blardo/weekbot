package commands

import (
	"weekbot/internal/models"

	"github.com/bwmarrin/discordgo"
)

// HandleWeekPoll handles the /poll command
func HandleWeekPoll(s *discordgo.Session, m *discordgo.InteractionCreate) {
	bot := models.GetBot(m.GuildID)

	poll := models.NewOrCurrentPoll(bot)

	if poll == nil {
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No suggestions found",
			},
		})
		return
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "Poll has started",
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "first_choice",
						Placeholder: "Select a week",
						Options:     poll.GetSelectOptions(),
					},
				},
			},
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "second_choice",
						Placeholder: "Select a week",
						Options:     poll.GetSelectOptions(),
					},
				},
			},
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "third_choice",
						Placeholder: "Select a week",
						Options:     poll.GetSelectOptions(),
					},
				},
			},
		},
	})
}

func HandleEndPoll(s *discordgo.Session, m *discordgo.InteractionCreate) {
	bot := models.GetBot(m.GuildID)

	poll := models.GetCurrentPoll(bot.DB)

	if poll == nil {
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No poll is currently in progress",
			},
		})
		return
	}

	poll.IsComplete = true
	bot.DB.Save(poll)

	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The poll has been ended",
		},
	})
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

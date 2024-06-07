package commands

import (
	"log"
	"weekbot/internal/models"

	"github.com/bwmarrin/discordgo"
)

// HandleWeekPoll handles the /poll command
func HandleWeekPoll(s *discordgo.Session, m *discordgo.InteractionCreate) {
	bot := models.GetBot(m.GuildID)

	poll := models.NewOrCurrentPoll(bot)
	
	if poll == nil {
		println("poll is nil")
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "No suggestions found",
		})
		return
	}
	println("poll is not nil")


	////
	////
	/// Idea - Lets send a button that shows modal with all the suggestions
	//// so we can process the information independantly of the message itself.
	////
	
	options := poll.GetSelectOptions()
	for _, option := range options {
		println("poll options: " + option.Label)
	}
	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Poll has started",
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							CustomID: "poll_button",
							Label:    "Vote Here",
							Style:    discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🗳️",
							},
						},
					},
				},
			},
		},
	})


	// click button, check if userId in array, if no show modal, select option, submit, add userId to array, add ballot to poll


	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
		if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "poll_button" {
	
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Select your options",
					Flags: discordgo.MessageFlagsEphemeral,
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
						&discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: "submit_button",
									Label:    "Submit",
									Style:    discordgo.PrimaryButton,
									Emoji: discordgo.ComponentEmoji{
										Name: "🗳️",
									},
								},
							},
						},
					},
				},
			})

			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}
		}
	})
	// create listener for select menu
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




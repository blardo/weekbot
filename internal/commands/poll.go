package commands

import (
	"log"
	"time"
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
	

	
	///
	/// Idea - Lets send a button that shows modal with all the suggestions
	///so we can process the information independantly of the message itself.
	///
	
	
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
	// done =======, done ===================, done ===========, 


	// poll button handler

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
		if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "poll_button" {
			bot := models.GetBot(i.GuildID)
			if bot == nil{
				println("Bot is nil")
				return
			}
			ballots := models.GetAllBallots(bot.DB)

			for _, ballot := range ballots {
				println("Ballot: " + ballot.VoterId)
			}
	
			poll := models.NewOrCurrentPoll(bot)
			if poll == nil {
				println("Poll is nil")
				return
			}
			println("handler poll pull ")
			
			if i.Member == nil{
				println("Member is nil")
				return
			}
			
			pollBallots := poll.GetBallots()
			for _, ballot := range pollBallots {
				println("Poll Ballot: " + ballot.VoterId)
			}
			
		
			if poll.HasBallot(bot.DB, i.Member.User.ID) {
	
				println("isVoter poll pull ")
				println(len(poll.Ballots))
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You have already voted",
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					println("Error responding to interaction: %v", err)
					return
				}
			} else {
				log.Println(i.Member.User.ID)
				ballot := models.Ballot{
					VoterId: i.Member.User.ID,
					PollID: poll.ID,
					Date: time.Now(),
				}
				println("Ballot created")
				poll.AddBallotForVoter(bot, i.Member.User.ID, ballot)
				pollBallots := poll.GetBallots()
				for _, ballot := range pollBallots {
					println("Poll Ballot: " + ballot.VoterId)
				}
			}
	
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

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			// This is a slash command interaction
		} else if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "first_choice" { {
			// This is a select menu interaction
			// You can access the selected options with i.MessageComponent.Values
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}

			ballot := models.GetBallotByVoterID(bot.DB, i.Member.User.ID)
			ballot.FirstChoice = i.MessageComponentData().Values[0]
			println("First choice: " + i.MessageComponentData().Values[0])
		}
		} else if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "second_choice" { {
			// This is a select menu interaction
			// You can access the selected options with i.MessageComponent.Values
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}
			println("Second choice: " + i.MessageComponentData().Values[0])
		} 
		} else if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "third_choice" { {
			// This is a select menu interaction
			// You can access the selected options with i.MessageComponent.Values
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}
			println("Third choice: " + i.MessageComponentData().Values[0])
		}
		}
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			// This is a slash command interaction
		} else if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == "submit_button" {

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your vote has been submitted",
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}
		}
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
	
	poll.EndPoll()
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




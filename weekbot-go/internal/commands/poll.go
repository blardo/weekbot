package commands

import (
	"log"
	"time"
	"weekbot-go/internal/models"

	"github.com/bwmarrin/discordgo"
)

// HandleWeekPoll handles the /poll command
func HandleWeekPoll(s *discordgo.Session, m *discordgo.InteractionCreate) {
	bot := models.GetBot(m.GuildID)

	poll := models.NewOrCurrentPoll(bot)

	if poll == nil {
		println("poll is nil")
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
							Emoji: discordgo.ComponentEmoji{
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
	case "first_choice", "second_choice", "third_choice":
		handlePollChoice(s, i)
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
		poll.AddBallotForVoter(bot, ballot)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Select your options",
			Flags:   discordgo.MessageFlagsEphemeral,
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
	}
}

func handlePollChoice(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	ballot := models.GetBallotByVoterID(bot.DB, userID)
	if ballot.ID == 0 || ballot.PollID != poll.ID {
		respondEphemeral(s, i, "Please click \"Vote Here\" first.")
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	switch i.MessageComponentData().CustomID {
	case "first_choice":
		ballot.FirstChoice = i.MessageComponentData().Values[0]
	case "second_choice":
		ballot.SecondChoice = i.MessageComponentData().Values[0]
	case "third_choice":
		ballot.ThirdChoice = i.MessageComponentData().Values[0]
	default:
		return
	}

	bot.DB.Save(ballot)
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

	ballot := models.GetBallotByVoterID(bot.DB, userID)
	if ballot.ID == 0 || ballot.PollID != poll.ID {
		respondEphemeral(s, i, "Please click \"Vote Here\" first.")
		return
	}

	if ballot.FirstChoice == "" || ballot.SecondChoice == "" || ballot.ThirdChoice == "" {
		respondEphemeral(s, i, "Please select all options.")
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
	ellibleBallots := 0

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
			ellibleBallots++
		}
	}
	if ellibleBallots < 5 {
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

	println("New name is: " + newName)
	_, err := s.GuildEdit(m.GuildID, &discordgo.GuildParams{
		Name: newName,
	})
	if err != nil {
		log.Printf("Error changing server name: %v", err)

	}
	log.Printf("Server name changed to: %s", newName)

	// end poll
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

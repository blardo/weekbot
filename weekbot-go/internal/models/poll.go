package models

import (
	"math/rand"
	"strings"
	"weekbot-go/internal/config"
	"weekbot-go/internal/logger"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// Poll is a struct that represents a poll
type Poll struct {
	gorm.Model
	ID          uint `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Suggestions []Suggestion
	InProgress  bool
	IsComplete  bool
	Ballots     []Ballot `gorm:"foreignKey:PollID"`
}

// NewOrCurrentPoll creates a new poll or returns the current one
func NewOrCurrentPoll(bot *Bot) *Poll {
	poll := GetCurrentPoll(bot.DB)
	if poll != nil {
		logger.Debug("Current poll found", "poll_id", poll.ID, "guild_id", bot.GuildID)
		return poll
	}

	suggestions := GetMostRecentUnusedSuggestions(bot.DB)
	logger.Debug("Checking suggestions for new poll", "count", len(suggestions), "guild_id", bot.GuildID)
	if len(suggestions) < config.MinSuggestionsToStartPoll() {
		logger.Info("Not enough suggestions to start poll", "count", len(suggestions), "required", config.MinSuggestionsToStartPoll(), "guild_id", bot.GuildID)
		return nil
	}

	poll = &Poll{
		Suggestions: suggestions,
		InProgress:  true,
	}
	bot.DB.Create(poll)
	logger.Info("Created new poll", "poll_id", poll.ID, "suggestions", len(suggestions), "guild_id", bot.GuildID)
	return poll
}

// GetCurrentPoll gets the current poll from the database
func GetCurrentPoll(db *gorm.DB) *Poll {
	var poll Poll
	db.Preload("Suggestions").Preload("Ballots").Where("in_progress = ? and is_complete = ?", true, false).First(&poll)
	if poll.ID == 0 {
		logger.Debug("No current poll found")
		return nil
	}

	return &poll
}

// GetSelectOptions gets the select options for the poll
func (p *Poll) GetSelectOptions() []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	var filter []Suggestion
	for _, suggestion := range p.Suggestions {
		isDuplicate := false
		for _, f := range filter {
			if strings.EqualFold(suggestion.Content, f.Content) {
				isDuplicate = true
				break
			}
		}

		// Filter out bot messages and invalid suggestions
		if !isDuplicate && suggestion.Updicks >= config.MinUpdicksToQualify() {
			// Skip suggestions that look like bot messages or are too long
			content := strings.ToLower(suggestion.Content)
			if strings.Contains(content, "week suggestion added") ||
				strings.Contains(content, "weekbot will add") ||
				strings.Contains(content, "already exists") ||
				len(suggestion.Content) > 100 {
				continue
			}
			filter = append(filter, suggestion)
		}
	}
	// discords options limit is 25
	if len(filter) > 25 {
		filter = filter[:25]
	}
	for _, suggestion := range filter {
		// Discord requires label and value to be 1-100 characters
		// Truncate if necessary, but prefer keeping the full content
		label := suggestion.Content
		value := suggestion.Content
		if len(label) > 100 {
			label = label[:97] + "..."
		}
		if len(value) > 100 {
			// For value, we can use a hash or ID, but for now truncate
			// In the future, we might want to use suggestion.ID as value
			value = value[:100]
		}
		options = append(options, discordgo.SelectMenuOption{
			Label:   label,
			Value:   value,
			Default: false,
			Emoji:   &discordgo.ComponentEmoji{Name: "📅"},
		})
	}

	return options
}

// GetSelectOptionsExcluding returns select options excluding the provided choices
func (p *Poll) GetSelectOptionsExcluding(excludeChoices []string) []discordgo.SelectMenuOption {
	allOptions := p.GetSelectOptions()
	var filtered []discordgo.SelectMenuOption
	
	excludeMap := make(map[string]bool)
	for _, choice := range excludeChoices {
		excludeMap[choice] = true
	}
	
	for _, option := range allOptions {
		if !excludeMap[option.Value] {
			filtered = append(filtered, option)
		}
	}
	
	return filtered
}

// IsVoter checks if a user is a voter
func (p *Poll) HasBallot(voterID string) bool {
	ballots := p.GetBallots()
	for _, ballot := range ballots {
		logger.Debug("Checking ballot", "voter_id", ballot.VoterId, "poll_id", ballot.PollID, "cast", ballot.Cast)
		if ballot.VoterId == voterID && ballot.PollID == p.ID {
			return true
		}
	}
	return false
}

func (p *Poll) BallotCast(voterID string) bool {
	var ballot Ballot
	ballots := p.GetBallots()
	for _, b := range ballots {
		if b.VoterId == voterID {
			ballot = b
			break
		}
	}
	return ballot.Cast

}

// GetBallots gets the ballots of the poll
func (p *Poll) GetBallots() []Ballot {
	return p.Ballots
}

// AddBallot adds a ballot to the poll
func (p *Poll) AddBallotToPoll(bot *Bot, ballot Ballot) {
	for index, b := range p.Ballots {
		if b.VoterId == ballot.VoterId {
			p.Ballots[index] = ballot
			break
		}
	}
	p.Ballots = append(p.Ballots, ballot)
	bot.DB.Save(&p)

}

func (p *Poll) AddBallotForVoter(bot *Bot, ballot Ballot) {
	var voter Voter
	bot.DB.Where("user_id = ?", ballot.VoterId).First(&voter)

	if voter.UserID == "" {
		// Voter does not exist, create new voter
		voter = Voter{
			UserID: ballot.VoterId,
		}
		bot.DB.Create(&voter)
		logger.Debug("Created new voter", "user_id", ballot.VoterId, "guild_id", bot.GuildID)
	}
	bot.DB.Create(&ballot)

	// Add the ballot to the poll
	logger.Debug("Adding ballot to poll", "voter_id", ballot.VoterId, "poll_id", p.ID)
	p.AddBallotToPoll(bot, ballot)

}

func (p *Poll) PerformRankedChoiceVoting() string {
	voteCounts := make(map[string]int)
	totalBallots := len(p.Ballots)
	totalEligibleBallots := 0
	// Initialize vote counts for each suggestion
	for _, suggestion := range p.Suggestions {
		voteCounts[suggestion.Content] = 0
	}

	// First round: count first-choice votes
	for _, ballot := range p.Ballots {
		if !ballot.Cast {
			continue
		}
		totalEligibleBallots++
		choices := ballot.GetChoices()
		if len(choices) > 0 {
			firstChoice := choices[0]
			voteCounts[firstChoice]++
			logger.Debug("First choice vote", "choice", firstChoice, "count", voteCounts[firstChoice])
		} else {
			// Fallback to legacy FirstChoice field for backward compatibility
			if ballot.FirstChoice != "" {
				voteCounts[ballot.FirstChoice]++
				logger.Debug("First choice vote (legacy)", "choice", ballot.FirstChoice, "count", voteCounts[ballot.FirstChoice])
			}
		}
	}

	// Remove suggestions that did not receive any votes in the first round
	for suggestion, count := range voteCounts {
		logger.Debug("Suggestion vote count", "suggestion", suggestion, "count", count)
	}
	for suggestion, count := range voteCounts {
		if count == 0 {
			logger.Debug("Removing suggestion with zero votes", "suggestion", suggestion)
			delete(voteCounts, suggestion)
		}
	}

	if totalEligibleBallots < 4 {
		allOneVote := true
		for _, count := range voteCounts {
			if count != 1 {
				allOneVote = false
				break
			}
		}
		if allOneVote {
			// Collect all suggestions with one vote
			var oneVoteSuggestions []string
			for suggestion, count := range voteCounts {
				if count == 1 {
					oneVoteSuggestions = append(oneVoteSuggestions, suggestion)
				}
			}
			// Randomly select a winner from the suggestions with one vote
			return oneVoteSuggestions[rand.Intn(len(oneVoteSuggestions))]
		}
	}

	// Check for majority
	for {
		// Find the suggestion with the highest votes
		var maxVotes int
		var maxSuggestion string
		for suggestion, count := range voteCounts {
			if count > maxVotes {
				maxVotes = count
				maxSuggestion = suggestion
			}
		}

		// Check if the suggestion has more than 50% of the votes
		if maxVotes > (totalBallots+1)/2 {
			return maxSuggestion
		}

		// Find the suggestion with the fewest votes
		var minVotes int = totalBallots + 1
		var minSuggestion string
		for suggestion, count := range voteCounts {
			if count > 0 && count < minVotes {
				minVotes = count
				minSuggestion = suggestion
			}
		}

		// Eliminate the suggestion with the fewest votes
		delete(voteCounts, minSuggestion)

		// Redistribute votes
		for _, ballot := range p.Ballots {
			choices := ballot.GetChoices()
			var firstChoice string
			if len(choices) > 0 {
				firstChoice = choices[0]
			} else {
				// Fallback to legacy FirstChoice field
				firstChoice = ballot.FirstChoice
			}
			
			if firstChoice == minSuggestion {
				// Find next available choice in the voter's ranking
				redistributed := false
				if len(choices) > 1 {
					// Try each subsequent choice in order
					for i := 1; i < len(choices); i++ {
						nextChoice := choices[i]
						if voteCounts[nextChoice] > 0 {
							voteCounts[nextChoice]++
							redistributed = true
							break
						}
					}
				}
				
				// Fallback to legacy fields if no choices array
				if !redistributed && len(choices) == 0 {
					if ballot.SecondChoice != "" && voteCounts[ballot.SecondChoice] > 0 {
						voteCounts[ballot.SecondChoice]++
					} else if ballot.ThirdChoice != "" && voteCounts[ballot.ThirdChoice] > 0 {
						voteCounts[ballot.ThirdChoice]++
					}
				}
			}
		}
	}
}

// EndPoll ends the poll and marks only the winning suggestion as used
// Non-winning suggestions can be suggested again in future polls
func (p *Poll) EndPoll(db *gorm.DB, winningSuggestion string) {
	logger.Info("Ending poll", "poll_id", p.ID, "winner", winningSuggestion)
	p.InProgress = false
	p.IsComplete = true

	// Only mark the winning suggestion as used
	// Other suggestions can be suggested again in future polls
	winningNormalized := NormalizeSuggestionContent(winningSuggestion)
	for _, suggestion := range p.Suggestions {
		suggestionNormalized := NormalizeSuggestionContent(suggestion.Content)
		if strings.EqualFold(suggestionNormalized, winningNormalized) {
			suggestion.Used = true
			if err := db.Save(&suggestion).Error; err != nil {
				logger.Error("Error marking winning suggestion as used", "error", err, "suggestion_id", suggestion.ID, "poll_id", p.ID)
			} else {
				logger.Debug("Marked winning suggestion as used", "suggestion_id", suggestion.ID, "content", suggestion.Content, "poll_id", p.ID)
			}
		} else {
			// Reset the suggestion so it can be used in future polls
			// Clear the PollID so it's not associated with this poll anymore
			suggestion.PollID = 0
			if err := db.Save(&suggestion).Error; err != nil {
				logger.Error("Error resetting suggestion for future polls", "error", err, "suggestion_id", suggestion.ID, "poll_id", p.ID)
			} else {
				logger.Debug("Reset suggestion for future polls", "suggestion_id", suggestion.ID, "content", suggestion.Content, "poll_id", p.ID)
			}
		}
	}
}

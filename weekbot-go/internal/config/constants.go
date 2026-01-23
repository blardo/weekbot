package config

import "weekbot-go/internal/services"

// MinSuggestionsToStartPoll returns the minimum number of suggestions needed to start a poll
// This value is configurable via MIN_SUGGESTIONS_TO_START_POLL environment variable (default: 3)
func MinSuggestionsToStartPoll() int {
	return services.GetConfig().MinSuggestionsToStart
}

// MinUpdicksToQualify returns the minimum number of "bd" reactions needed for a suggestion to qualify
// This value is configurable via MIN_UPDICKS_TO_QUALIFY environment variable (default: 3)
func MinUpdicksToQualify() int {
	return services.GetConfig().MinUpdicksToQualify
}

// MinBallotsToEndPoll returns the minimum number of cast ballots needed to end a poll
// This value is configurable via MIN_BALLOTS_TO_END_POLL environment variable (default: 5)
func MinBallotsToEndPoll() int {
	return services.GetConfig().MinBallotsToEndPoll
}

// Discord channel and emoji constants
const (
	// WeekNameChannelName is the name of the channel where week suggestions are posted
	WeekNameChannelName = "week-name"
	// QualifyingEmoji is the emoji that qualifies suggestions (needs MinUpdicksToQualify reactions)
	QualifyingEmoji = "bd"
	// ConfirmationEmoji is the emoji added when a suggestion qualifies
	ConfirmationEmoji = "👍"
)

// AcceptableWeekSuffixes are the valid endings for week name suggestions
var AcceptableWeekSuffixes = []string{"week", "week.", "week!", "week?"}

package config

// Poll configuration constants
const (
	// MinSuggestionsToStartPoll is the minimum number of suggestions needed to start a poll
	MinSuggestionsToStartPoll = 1 // Set to 1 for testing, should be 3 in production
	// MinUpdicksToQualify is the minimum number of "bd" reactions needed for a suggestion to qualify
	MinUpdicksToQualify = 1 // Set to 1 for testing, should be 3 in production
	// MinBallotsToEndPoll is the minimum number of cast ballots needed to end a poll
	MinBallotsToEndPoll = 1 // Set to 1 for testing, should be 5 in production
)

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

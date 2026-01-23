package config

// Poll configuration constants
const (
	// MinSuggestionsToStartPoll is the minimum number of suggestions needed to start a poll
	MinSuggestionsToStartPoll = 3
	// MinUpdicksToQualify is the minimum number of "bd" reactions needed for a suggestion to qualify
	MinUpdicksToQualify = 3
	// MinBallotsToEndPoll is the minimum number of cast ballots needed to end a poll
	MinBallotsToEndPoll = 5
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

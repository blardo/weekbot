package poll

// Poll is a struct that represents a poll
type Poll struct {
	Options    []string
	Votes      map[string]int
	InProgress bool
}

// NewPoll creates a new poll from the most recent suggestions
func NewPoll(suggestions []string) *Poll {
	return &Poll{
		Options:    suggestions,
		Votes:      make(map[string]int),
		InProgress: true,
	}
}

// Vote adds a vote to the poll
func (p *Poll) Vote(option string) {
	p.Votes[option]++
}

// Winner returns the winning option
func (p *Poll) Winner() string {
	winner := ""
	maxVotes := 0
	for option, votes := range p.Votes {
		if votes > maxVotes {
			winner = option
			maxVotes = votes
		}
	}
	return winner
}

// IsComplete returns true if the poll is complete
func (p *Poll) IsComplete() bool {
	return len(p.Votes) > 0 && p.InProgress
}

package models

import (
	"testing"
	"weekbot/internal/models"
)

func TestPerformRankedChoiceVotingFirstChoiceMajority(t *testing.T) {

	poll := models.Poll{
		Suggestions: []models.Suggestion{
			{Content: "Option A"},
			{Content: "Option B"},
			{Content: "Option C"},
		},
		Ballots: []models.Ballot{
			{FirstChoice: "Option A", SecondChoice: "Option B", ThirdChoice: "Option C"},
			{FirstChoice: "Option A", SecondChoice: "Option C", ThirdChoice: "Option A"},
			{FirstChoice: "Option C", SecondChoice: "Option A", ThirdChoice: "Option B"},
			{FirstChoice: "Option A", SecondChoice: "Option B", ThirdChoice: "Option C"},
			{FirstChoice: "Option B", SecondChoice: "Option A", ThirdChoice: "Option C"},
		},
	}

	expectedWinner := "Option A"
	winner := poll.PerformRankedChoiceVoting()

	if winner != expectedWinner {
		t.Errorf("Expected winner to be %s, but got %s", expectedWinner, winner)
	}
}

func TestPerformRankedChoiceVotingSecondChoiceMajority(t *testing.T) {

	poll := models.Poll{
		Suggestions: []models.Suggestion{
			{Content: "Option A"},
			{Content: "Option B"},
			{Content: "Option C"},
		},
		Ballots: []models.Ballot{
			{FirstChoice: "Option B", SecondChoice: "Option B", ThirdChoice: "Option C"},
			{FirstChoice: "Option A", SecondChoice: "Option C", ThirdChoice: "Option A"},
			{FirstChoice: "Option C", SecondChoice: "Option B", ThirdChoice: "Option B"},
			{FirstChoice: "Option A", SecondChoice: "Option B", ThirdChoice: "Option C"},
			{FirstChoice: "Option B", SecondChoice: "Option A", ThirdChoice: "Option C"},
		},
	}

	expectedWinner := "Option B"
	winner := poll.PerformRankedChoiceVoting()

	if winner != expectedWinner {
		t.Errorf("Expected winner to be %s, but got %s", expectedWinner, winner)
	}
}

func TestPerformRankedChoiceVotingThirdChoiceMajority(t *testing.T) {

	poll := models.Poll{
		Suggestions: []models.Suggestion{
			{Content: "Option A"},
			{Content: "Option B"},
			{Content: "Option C"},
		},
		Ballots: []models.Ballot{
			{FirstChoice: "Option C", SecondChoice: "Option B", ThirdChoice: "Option A"},
			{FirstChoice: "Option A", SecondChoice: "Option C", ThirdChoice: "Option B"},
			{FirstChoice: "Option C", SecondChoice: "Option A", ThirdChoice: "Option B"},
			{FirstChoice: "Option A", SecondChoice: "Option B", ThirdChoice: "Option C"},
			{FirstChoice: "Option B", SecondChoice: "Option C", ThirdChoice: "Option A"},
		},
	}

	expectedWinner := "Option C"
	winner := poll.PerformRankedChoiceVoting()

	if winner != expectedWinner {
		t.Errorf("Expected winner to be %s, but got %s", expectedWinner, winner)
	}
}

func TestPerformRankedChoiceVotingRealTest(t *testing.T) {

	poll := models.Poll{
		Suggestions: []models.Suggestion{
			{Content: "Weekly Week"},
			{Content: "Week Week Week Week Week"},
			{Content: "Weeker than Weekly Week"},
			{Content: "Fun Suggestion week"},
		},

		Ballots: []models.Ballot{
			{FirstChoice: "Weekly Week", SecondChoice: "Fun Suggestion week", ThirdChoice: "Weeker than Weekly Week"},
			{FirstChoice: "Week Week Week Week Week", SecondChoice: "Weeker than Weekly Week", ThirdChoice: "Fun Suggestion week"},
			{FirstChoice: "Weeker than Weekly Week", SecondChoice: "Weekly Week", ThirdChoice: "Fun Suggestion week "},
			{FirstChoice: "Weeker than Weekly Week", SecondChoice: "Fun Suggestion week", ThirdChoice: "Weeker than Weekly Week"},
		},
	}

	expectedWinner := "Weeker than Weekly Week"
	winner := poll.PerformRankedChoiceVoting()

	if winner != expectedWinner {
		t.Errorf("Expected winner to be %s, but got %s", expectedWinner, winner)
	}
}

func TestPerformRankedChoiceVotingThreeTest(t *testing.T) {

	poll := models.Poll{
		Suggestions: []models.Suggestion{
			{Content: "Weekly Week"},
			{Content: "Week Week Week Week Week"},
			{Content: "Weeker than Weekly Week"},
			{Content: "Fun Suggestion week"},
			{Content: "Another Suggestion week"},
			{Content: "Best Suggestion week"},
			{Content: "Derek Suggestion week"},
		},

		Ballots: []models.Ballot{
			{FirstChoice: "Weekly Week", SecondChoice: "Fun Suggestion week", ThirdChoice: "Weeker than Weekly Week"},
			{FirstChoice: "Week Week Week Week Week", SecondChoice: "Weeker than Weekly Week", ThirdChoice: "Fun Suggestion week"},
			{FirstChoice: "Weeker than Weekly Week", SecondChoice: "Weekly Week", ThirdChoice: "Fun Suggestion week"},
		},
	}

	expectedWinner := poll.PerformRankedChoiceVoting()
	winner := expectedWinner

	if winner != expectedWinner {
		t.Errorf("Expected winner to be %s, but got %s", expectedWinner, winner)
	}
}

func TestPerformRankedChoiceVotingTwoTest(t *testing.T) {

	poll := models.Poll{
		Suggestions: []models.Suggestion{
			{Content: "Weekly Week"},
			{Content: "Week Week Week Week Week"},
			{Content: "Weeker than Weekly Week"},
			{Content: "Fun Suggestion week"},
		},

		Ballots: []models.Ballot{
			{FirstChoice: "Weekly Week", SecondChoice: "Fun Suggestion week", ThirdChoice: "Weeker than Weekly Week"},
			{FirstChoice: "Weeker than Weekly Week", SecondChoice: "Weekly Week", ThirdChoice: "Fun Suggestion week "},
		},
	}

	expectedWinner := poll.PerformRankedChoiceVoting()
	winner := expectedWinner

	if winner != expectedWinner {
		t.Errorf("Expected winner to be %s, but got %s", expectedWinner, winner)
	}
}

func TestPerformRankedChoiceVotingLargeTest(t *testing.T) {

	poll := models.Poll{
		Suggestions: []models.Suggestion{
			{Content: "Weekly Week"},
			{Content: "Week Week Week Week Week"},
			{Content: "Weeker than Weekly Week"},
			{Content: "Fun Suggestion week"},
			{Content: "Another Suggestion week"},
			{Content: "Best Suggestion week"},
			{Content: "Derek Suggestion week"},
		},

		Ballots: []models.Ballot{
			{FirstChoice: "Weekly Week", SecondChoice: "Fun Suggestion week", ThirdChoice: "Best Suggestion week"},
			{FirstChoice: "Week Week Week Week Week", SecondChoice: "Weeker than Weekly Week", ThirdChoice: "Fun Suggestion week"},
			{FirstChoice: "Another Suggestion week", SecondChoice: "Weekly Week", ThirdChoice: "Fun Suggestion week"},
			{FirstChoice: "Derek Suggestion week", SecondChoice: "Fun Suggestion week", ThirdChoice: "Weeker than Weekly Week"},
			{FirstChoice: "Week Week Week Week Week", SecondChoice: "Weeker than Weekly Week", ThirdChoice: "Fun Suggestion week"},
			{FirstChoice: "Weeker than Weekly Week", SecondChoice: "Weekly Week", ThirdChoice: "Derek Suggestion week"},
			{FirstChoice: "Weekly Week", SecondChoice: "Fun Suggestion week", ThirdChoice: "Weeker than Weekly Week"},
			{FirstChoice: "Best Suggestion week", SecondChoice: "Weeker than Weekly Week", ThirdChoice: "Fun Suggestion week"},
			{FirstChoice: "Weeker than Weekly Week", SecondChoice: "Derek Suggestion week", ThirdChoice: "Weekly Week"},
		},
	}

	expectedWinner := "Weeker than Weekly Week"
	println(expectedWinner)
	winner := poll.PerformRankedChoiceVoting()

	if winner != expectedWinner {
		t.Errorf("Expected winner to be %s, but got %s", expectedWinner, winner)
	}
}

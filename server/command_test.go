package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatLeaderboardMarkdownEmpty(t *testing.T) {
	text := formatLeaderboardMarkdown("🏆 Classement du channel", PeriodWeek, Leaderboard{})

	assert.Contains(t, text, "7 derniers jours")
	assert.Contains(t, text, "Aucun message compté")
}

func TestFormatLeaderboardMarkdownTable(t *testing.T) {
	leaderboard := Leaderboard{
		Entries: []LeaderboardEntry{
			{UserID: "u1", Username: "alice", Count: 42, Rank: 1},
			{UserID: "u2", Username: "bob", Count: 37, Rank: 2},
		},
		Me: &LeaderboardEntry{UserID: "u9", Username: "dave", Count: 3, Rank: 17},
	}

	text := formatLeaderboardMarkdown("🏆 Classement global", PeriodAll, leaderboard)

	assert.Contains(t, text, "depuis le début")
	assert.Contains(t, text, "| 1 | @alice | 42 |")
	assert.Contains(t, text, "| 2 | @bob | 37 |")
	assert.Contains(t, text, "| 17 | @dave (toi) | 3 |")
}

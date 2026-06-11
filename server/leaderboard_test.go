package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriodDays(t *testing.T) {
	now := time.Date(2026, 6, 11, 15, 0, 0, 0, time.UTC)

	week := periodDays(PeriodWeek, now)
	require.Len(t, week, 7)
	assert.Equal(t, "2026-06-11", week[0])
	assert.Equal(t, "2026-06-05", week[6])

	month := periodDays(PeriodMonth, now)
	require.Len(t, month, 30)
	assert.Equal(t, "2026-06-11", month[0])
	assert.Equal(t, "2026-05-13", month[29])

	assert.Nil(t, periodDays(PeriodAll, now))
}

func TestPeriodDaysCrossesMonthBoundary(t *testing.T) {
	now := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	week := periodDays(PeriodWeek, now)
	assert.Equal(t, "2026-02-24", week[6])
}

func TestValidPeriod(t *testing.T) {
	assert.True(t, validPeriod(PeriodWeek))
	assert.True(t, validPeriod(PeriodMonth))
	assert.True(t, validPeriod(PeriodAll))
	assert.False(t, validPeriod("year"))
	assert.False(t, validPeriod(""))
}

func TestMergeCounts(t *testing.T) {
	dst := map[string]int64{"alice": 2, "bob": 1}
	mergeCounts(dst, map[string]int64{"alice": 3, "carol": 5})

	assert.Equal(t, map[string]int64{"alice": 5, "bob": 1, "carol": 5}, dst)
}

func identityResolver(userID string) string { return userID }

func TestBuildLeaderboardSortsAndRanks(t *testing.T) {
	counts := map[string]int64{"alice": 10, "bob": 30, "carol": 20}

	leaderboard := buildLeaderboard(counts, 10, "alice", identityResolver)

	require.Len(t, leaderboard.Entries, 3)
	assert.Equal(t, "bob", leaderboard.Entries[0].UserID)
	assert.Equal(t, 1, leaderboard.Entries[0].Rank)
	assert.Equal(t, "carol", leaderboard.Entries[1].UserID)
	assert.Equal(t, "alice", leaderboard.Entries[2].UserID)
	assert.Equal(t, 3, leaderboard.Entries[2].Rank)
	assert.Nil(t, leaderboard.Me)
}

func TestBuildLeaderboardTiesAreDeterministic(t *testing.T) {
	counts := map[string]int64{"zed": 5, "ann": 5}

	leaderboard := buildLeaderboard(counts, 10, "", identityResolver)

	require.Len(t, leaderboard.Entries, 2)
	assert.Equal(t, "ann", leaderboard.Entries[0].UserID)
	assert.Equal(t, "zed", leaderboard.Entries[1].UserID)
}

func TestBuildLeaderboardTruncatesAndFillsMe(t *testing.T) {
	counts := map[string]int64{"a": 50, "b": 40, "c": 30, "me": 10}

	leaderboard := buildLeaderboard(counts, 2, "me", identityResolver)

	require.Len(t, leaderboard.Entries, 2)
	require.NotNil(t, leaderboard.Me)
	assert.Equal(t, "me", leaderboard.Me.UserID)
	assert.Equal(t, 4, leaderboard.Me.Rank)
	assert.Equal(t, int64(10), leaderboard.Me.Count)
}

func TestBuildLeaderboardMeInTopHasNoMeEntry(t *testing.T) {
	counts := map[string]int64{"me": 50, "b": 40}

	leaderboard := buildLeaderboard(counts, 10, "me", identityResolver)

	assert.Nil(t, leaderboard.Me)
}

func TestBuildLeaderboardEmpty(t *testing.T) {
	leaderboard := buildLeaderboard(map[string]int64{}, 10, "me", identityResolver)

	assert.Empty(t, leaderboard.Entries)
	assert.Nil(t, leaderboard.Me)
}

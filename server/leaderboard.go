package main

import (
	"sort"
	"time"

	"github.com/pkg/errors"
)

const (
	PeriodWeek  = "week"
	PeriodMonth = "month"
	PeriodAll   = "all"
)

type LeaderboardEntry struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Count    int64  `json:"count"`
	Rank     int    `json:"rank"`
}

type Leaderboard struct {
	Entries []LeaderboardEntry `json:"entries"`
	// Me is the requesting user's own entry when ranked outside the top list.
	Me *LeaderboardEntry `json:"me,omitempty"`
}

func validPeriod(period string) bool {
	return period == PeriodWeek || period == PeriodMonth || period == PeriodAll
}

// periodDays returns the rolling list of days (YYYY-MM-DD, UTC) covered by a
// period, most recent first. Returns nil for PeriodAll.
func periodDays(period string, now time.Time) []string {
	var n int
	switch period {
	case PeriodWeek:
		n = 7
	case PeriodMonth:
		n = 30
	default:
		return nil
	}
	days := make([]string, 0, n)
	for i := 0; i < n; i++ {
		days = append(days, now.UTC().AddDate(0, 0, -i).Format("2006-01-02"))
	}
	return days
}

// mergeCounts adds src counts into dst.
func mergeCounts(dst, src map[string]int64) {
	for userID, count := range src {
		dst[userID] += count
	}
}

// channelCounts returns userID → count for one channel over a period.
func (p *Plugin) channelCounts(channelID, period string, now time.Time) (map[string]int64, error) {
	if period == PeriodAll {
		return p.kvstore.GetChannelTotals(channelID)
	}
	counts := map[string]int64{}
	for _, day := range periodDays(period, now) {
		dayCounts, err := p.kvstore.GetDayCounts(day, channelID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to aggregate channel counts")
		}
		mergeCounts(counts, dayCounts)
	}
	return counts, nil
}

// globalCounts returns userID → count across every channel over a period.
func (p *Plugin) globalCounts(period string, now time.Time) (map[string]int64, error) {
	counts := map[string]int64{}
	if period == PeriodAll {
		byChannel, err := p.kvstore.GetAllChannelTotals()
		if err != nil {
			return nil, errors.Wrap(err, "failed to aggregate global totals")
		}
		for _, channelCounts := range byChannel {
			mergeCounts(counts, channelCounts)
		}
		return counts, nil
	}
	for _, day := range periodDays(period, now) {
		byChannel, err := p.kvstore.GetAllDayCounts(day)
		if err != nil {
			return nil, errors.Wrap(err, "failed to aggregate global day counts")
		}
		for _, channelCounts := range byChannel {
			mergeCounts(counts, channelCounts)
		}
	}
	return counts, nil
}

// buildLeaderboard sorts counts, keeps the top size entries, resolves
// usernames, and fills Me when meID ranks outside the top list.
func buildLeaderboard(counts map[string]int64, size int, meID string, resolveUsername func(userID string) string) Leaderboard {
	entries := make([]LeaderboardEntry, 0, len(counts))
	for userID, count := range counts {
		entries = append(entries, LeaderboardEntry{UserID: userID, Count: count})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].UserID < entries[j].UserID
	})

	leaderboard := Leaderboard{Entries: []LeaderboardEntry{}}
	for i := range entries {
		entries[i].Rank = i + 1
		if i < size {
			entries[i].Username = resolveUsername(entries[i].UserID)
			leaderboard.Entries = append(leaderboard.Entries, entries[i])
		} else if entries[i].UserID == meID {
			entries[i].Username = resolveUsername(entries[i].UserID)
			me := entries[i]
			leaderboard.Me = &me
		}
	}
	return leaderboard
}

package kvstore

// KVStore exposes the persistence layer of the gamification plugin.
//
// Counts are stored as maps of userID → count:
//   - chan_total:{channelID}        cumulative counts for one channel
//   - day:{YYYY-MM-DD}:{channelID}  counts for one channel on one day (TTL ~40 days)
type KVStore interface {
	// IncrementCount adds one message for the user in the channel, both in the
	// channel cumulative counter and in the per-day counter.
	IncrementCount(channelID, userID, day string) error

	// GetChannelTotals returns the cumulative counts for one channel.
	GetChannelTotals(channelID string) (map[string]int64, error)

	// GetAllChannelTotals returns the cumulative counts of every channel, keyed by channelID.
	GetAllChannelTotals() (map[string]map[string]int64, error)

	// GetDayCounts returns the counts of one channel for one day.
	GetDayCounts(day, channelID string) (map[string]int64, error)

	// GetAllDayCounts returns the counts of every channel for one day, keyed by channelID.
	GetAllDayCounts(day string) (map[string]map[string]int64, error)
}

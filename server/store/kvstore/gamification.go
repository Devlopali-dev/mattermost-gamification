package kvstore

import (
	"maps"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"
)

const (
	chanTotalPrefix = "chan_total:"
	dayPrefix       = "day:"

	// dayKeyTTL keeps per-day counters slightly longer than the 30-day
	// aggregation window so the KV store stays bounded without a purge job.
	dayKeyTTL = 40 * 24 * time.Hour

	casRetries   = 5
	listPageSize = 200
)

type Client struct {
	client *pluginapi.Client
}

func NewKVStore(client *pluginapi.Client) KVStore {
	return Client{
		client: client,
	}
}

func chanTotalKey(channelID string) string {
	return chanTotalPrefix + channelID
}

func dayKey(day, channelID string) string {
	return dayPrefix + day + ":" + channelID
}

func (kv Client) IncrementCount(channelID, userID, day string) error {
	if err := kv.incrementMapValue(chanTotalKey(channelID), userID, 0); err != nil {
		return errors.Wrap(err, "failed to increment channel total")
	}
	if err := kv.incrementMapValue(dayKey(day, channelID), userID, dayKeyTTL); err != nil {
		return errors.Wrap(err, "failed to increment day count")
	}
	return nil
}

// incrementMapValue performs an atomic compare-and-set increment of one entry
// in the map stored under key. Safe with concurrent writers (HA cluster).
func (kv Client) incrementMapValue(key, userID string, ttl time.Duration) error {
	for range casRetries {
		var counts map[string]int64
		if err := kv.client.KV.Get(key, &counts); err != nil {
			return errors.Wrap(err, "failed to read counts")
		}

		var oldValue any
		if counts == nil {
			counts = map[string]int64{}
		} else {
			old := make(map[string]int64, len(counts))
			maps.Copy(old, counts)
			oldValue = old
		}
		counts[userID]++

		options := []pluginapi.KVSetOption{pluginapi.SetAtomic(oldValue)}
		if ttl > 0 {
			options = append(options, pluginapi.SetExpiry(ttl))
		}

		saved, err := kv.client.KV.Set(key, counts, options...)
		if err != nil {
			return errors.Wrap(err, "failed to write counts")
		}
		if saved {
			return nil
		}
	}
	return errors.Errorf("failed to increment %q after %d attempts", key, casRetries)
}

func (kv Client) GetChannelTotals(channelID string) (map[string]int64, error) {
	return kv.getCounts(chanTotalKey(channelID))
}

func (kv Client) GetAllChannelTotals() (map[string]map[string]int64, error) {
	return kv.getCountsByPrefix(chanTotalPrefix)
}

func (kv Client) GetDayCounts(day, channelID string) (map[string]int64, error) {
	return kv.getCounts(dayKey(day, channelID))
}

func (kv Client) GetAllDayCounts(day string) (map[string]map[string]int64, error) {
	return kv.getCountsByPrefix(dayPrefix + day + ":")
}

func (kv Client) getCounts(key string) (map[string]int64, error) {
	var counts map[string]int64
	if err := kv.client.KV.Get(key, &counts); err != nil {
		return nil, errors.Wrapf(err, "failed to get counts for %q", key)
	}
	if counts == nil {
		counts = map[string]int64{}
	}
	return counts, nil
}

// getCountsByPrefix scans every key matching prefix and returns the stored
// maps keyed by the key suffix (the channelID).
func (kv Client) getCountsByPrefix(prefix string) (map[string]map[string]int64, error) {
	result := map[string]map[string]int64{}
	for page := 0; ; page++ {
		keys, err := kv.client.KV.ListKeys(page, listPageSize, pluginapi.WithPrefix(prefix))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list keys with prefix %q", prefix)
		}
		for _, key := range keys {
			counts, err := kv.getCounts(key)
			if err != nil {
				return nil, err
			}
			result[strings.TrimPrefix(key, prefix)] = counts
		}
		if len(keys) < listPageSize {
			return result, nil
		}
	}
}

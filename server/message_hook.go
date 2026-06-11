package main

import (
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// MessageHasBeenPosted counts every human message posted in a team channel.
func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
	if !p.shouldCount(post) {
		return
	}

	day := time.Now().UTC().Format("2006-01-02")
	if err := p.kvstore.IncrementCount(post.ChannelId, post.UserId, day); err != nil {
		p.client.Log.Error("Failed to increment message count", "user_id", post.UserId, "channel_id", post.ChannelId, "error", err.Error())
	}
}

func (p *Plugin) shouldCount(post *model.Post) bool {
	user, err := p.client.User.Get(post.UserId)
	if err != nil {
		p.client.Log.Warn("Failed to get user for counting", "user_id", post.UserId, "error", err.Error())
		return false
	}

	channel, err := p.client.Channel.Get(post.ChannelId)
	if err != nil {
		p.client.Log.Warn("Failed to get channel for counting", "channel_id", post.ChannelId, "error", err.Error())
		return false
	}

	return shouldCountPost(post, user, channel, p.getConfiguration().excludedChannelSet())
}

// shouldCountPost holds the pure counting rules: human messages only (threads
// included), in team channels only, excluding configured channels.
func shouldCountPost(post *model.Post, user *model.User, channel *model.Channel, excludedChannels map[string]bool) bool {
	if post == nil || user == nil || channel == nil {
		return false
	}
	if post.IsSystemMessage() {
		return false
	}
	if post.GetProp("from_webhook") == "true" || post.GetProp("from_bot") == "true" {
		return false
	}
	if user.IsBot {
		return false
	}
	if channel.Type != model.ChannelTypeOpen && channel.Type != model.ChannelTypePrivate {
		return false
	}
	if excludedChannels[channel.Id] {
		return false
	}
	return true
}

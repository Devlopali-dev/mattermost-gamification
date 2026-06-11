package main

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
)

func humanPost() *model.Post {
	return &model.Post{Id: "post1", UserId: "user1", ChannelId: "chan1", Message: "hello"}
}

func humanUser() *model.User {
	return &model.User{Id: "user1"}
}

func openChannel() *model.Channel {
	return &model.Channel{Id: "chan1", Type: model.ChannelTypeOpen}
}

func TestShouldCountPost(t *testing.T) {
	noExclusions := map[string]bool{}

	t.Run("counts a human message in an open channel", func(t *testing.T) {
		assert.True(t, shouldCountPost(humanPost(), humanUser(), openChannel(), noExclusions))
	})

	t.Run("counts a message in a private channel", func(t *testing.T) {
		channel := openChannel()
		channel.Type = model.ChannelTypePrivate
		assert.True(t, shouldCountPost(humanPost(), humanUser(), channel, noExclusions))
	})

	t.Run("counts a thread reply", func(t *testing.T) {
		post := humanPost()
		post.RootId = "root1"
		assert.True(t, shouldCountPost(post, humanUser(), openChannel(), noExclusions))
	})

	t.Run("ignores system messages", func(t *testing.T) {
		post := humanPost()
		post.Type = model.PostTypeJoinChannel
		assert.False(t, shouldCountPost(post, humanUser(), openChannel(), noExclusions))
	})

	t.Run("ignores webhook posts", func(t *testing.T) {
		post := humanPost()
		post.AddProp("from_webhook", "true")
		assert.False(t, shouldCountPost(post, humanUser(), openChannel(), noExclusions))
	})

	t.Run("ignores bot-flagged posts", func(t *testing.T) {
		post := humanPost()
		post.AddProp("from_bot", "true")
		assert.False(t, shouldCountPost(post, humanUser(), openChannel(), noExclusions))
	})

	t.Run("ignores bot users", func(t *testing.T) {
		user := humanUser()
		user.IsBot = true
		assert.False(t, shouldCountPost(humanPost(), user, openChannel(), noExclusions))
	})

	t.Run("ignores direct messages", func(t *testing.T) {
		channel := openChannel()
		channel.Type = model.ChannelTypeDirect
		assert.False(t, shouldCountPost(humanPost(), humanUser(), channel, noExclusions))
	})

	t.Run("ignores group messages", func(t *testing.T) {
		channel := openChannel()
		channel.Type = model.ChannelTypeGroup
		assert.False(t, shouldCountPost(humanPost(), humanUser(), channel, noExclusions))
	})

	t.Run("ignores excluded channels", func(t *testing.T) {
		excluded := map[string]bool{"chan1": true}
		assert.False(t, shouldCountPost(humanPost(), humanUser(), openChannel(), excluded))
	})

	t.Run("ignores nil inputs", func(t *testing.T) {
		assert.False(t, shouldCountPost(nil, humanUser(), openChannel(), noExclusions))
		assert.False(t, shouldCountPost(humanPost(), nil, openChannel(), noExclusions))
		assert.False(t, shouldCountPost(humanPost(), humanUser(), nil, noExclusions))
	})
}

func TestExcludedChannelSet(t *testing.T) {
	config := &configuration{ExcludedChannels: "chan1, chan2 ,,chan3"}
	assert.Equal(t, map[string]bool{"chan1": true, "chan2": true, "chan3": true}, config.excludedChannelSet())

	empty := &configuration{}
	assert.Empty(t, empty.excludedChannelSet())
}

func TestLeaderboardSizeDefault(t *testing.T) {
	assert.Equal(t, 10, (&configuration{}).leaderboardSize())
	assert.Equal(t, 25, (&configuration{LeaderboardSize: 25}).leaderboardSize())
}

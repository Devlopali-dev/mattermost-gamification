package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTPRequiresAuthentication(t *testing.T) {
	plugin := Plugin{}
	plugin.router = plugin.initRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard", nil)

	plugin.ServeHTTP(nil, w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestServeHTTPLeaderboardRequiresChannelID(t *testing.T) {
	plugin := Plugin{}
	plugin.router = plugin.initRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard", nil)
	r.Header.Set("Mattermost-User-ID", "test-user-id")

	plugin.ServeHTTP(nil, w, r)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestServeHTTPLeaderboardRejectsInvalidPeriod(t *testing.T) {
	plugin := Plugin{}
	plugin.router = plugin.initRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/leaderboard?channel_id=chan1&period=year", nil)
	r.Header.Set("Mattermost-User-ID", "test-user-id")

	plugin.ServeHTTP(nil, w, r)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type leaderboardResponse struct {
	Period  string       `json:"period"`
	Channel *Leaderboard `json:"channel,omitempty"`
	Global  Leaderboard  `json:"global"`
}

// initRouter initializes the HTTP router for the plugin.
func (p *Plugin) initRouter() *mux.Router {
	router := mux.NewRouter()

	// Middleware to require that the user is logged in
	router.Use(p.MattermostAuthorizationRequired)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.HandleFunc("/leaderboard", p.handleLeaderboard).Methods(http.MethodGet)

	return router
}

// ServeHTTP handles requests under <siteUrl>/plugins/com.devlopali.gamification/.
func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) MattermostAuthorizationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleLeaderboard serves GET /api/v1/leaderboard?period=week|month|all[&channel_id=...]
// Without channel_id, only the global leaderboard is returned.
func (p *Plugin) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")

	channelID := r.URL.Query().Get("channel_id")

	period := r.URL.Query().Get("period")
	if period == "" {
		period = PeriodWeek
	}
	if !validPeriod(period) {
		http.Error(w, "period must be week, month or all", http.StatusBadRequest)
		return
	}

	now := time.Now()
	size := p.getConfiguration().leaderboardSize()
	resolve := p.usernameResolver()

	var channelBoard *Leaderboard
	if channelID != "" {
		// Only channel members may read a channel's stats.
		if _, err := p.client.Channel.GetMember(channelID, userID); err != nil {
			http.Error(w, "Not a member of this channel", http.StatusForbidden)
			return
		}

		channelCounts, err := p.channelCounts(channelID, period, now)
		if err != nil {
			p.client.Log.Error("Failed to compute channel leaderboard", "channel_id", channelID, "error", err.Error())
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		board := buildLeaderboard(channelCounts, size, userID, resolve)
		channelBoard = &board
	}

	globalCounts, err := p.globalCounts(period, now)
	if err != nil {
		p.client.Log.Error("Failed to compute global leaderboard", "error", err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	response := leaderboardResponse{
		Period:  period,
		Channel: channelBoard,
		Global:  buildLeaderboard(globalCounts, size, userID, resolve),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.client.Log.Error("Failed to encode leaderboard response", "error", err.Error())
	}
}

// usernameResolver returns a memoizing userID → username lookup.
func (p *Plugin) usernameResolver() func(userID string) string {
	cache := map[string]string{}
	return func(userID string) string {
		if username, ok := cache[userID]; ok {
			return username
		}
		username := userID
		if user, err := p.client.User.Get(userID); err == nil {
			username = user.Username
		}
		cache[userID] = username
		return username
	}
}

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

const leaderboardCommandTrigger = "leaderboard"

func (p *Plugin) registerCommands() error {
	if err := p.client.SlashCommand.Register(&model.Command{
		Trigger:          leaderboardCommandTrigger,
		AutoComplete:     true,
		AutoCompleteDesc: "Affiche le classement des messages (channel courant ou global)",
		AutoCompleteHint: "[global] [week|month|all]",
		DisplayName:      "Leaderboard",
	}); err != nil {
		return errors.Wrap(err, "failed to register leaderboard command")
	}
	return nil
}

// ExecuteCommand handles /leaderboard [global] [week|month|all].
func (p *Plugin) ExecuteCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	fields := strings.Fields(args.Command)
	if len(fields) == 0 || strings.TrimPrefix(fields[0], "/") != leaderboardCommandTrigger {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Commande inconnue : %s", args.Command),
		}, nil
	}

	global := false
	period := PeriodWeek
	for _, arg := range fields[1:] {
		switch strings.ToLower(arg) {
		case "global":
			global = true
		case PeriodWeek, PeriodMonth, PeriodAll:
			period = strings.ToLower(arg)
		default:
			return &model.CommandResponse{
				ResponseType: model.CommandResponseTypeEphemeral,
				Text:         "Usage : `/leaderboard [global] [week|month|all]`",
			}, nil
		}
	}

	text, err := p.leaderboardMarkdown(args.ChannelId, args.UserId, period, global)
	if err != nil {
		p.client.Log.Error("Failed to build leaderboard command response", "error", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Erreur lors du calcul du classement.",
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text,
	}, nil
}

func (p *Plugin) leaderboardMarkdown(channelID, userID, period string, global bool) (string, error) {
	now := time.Now()
	size := p.getConfiguration().leaderboardSize()

	var counts map[string]int64
	var err error
	var title string
	if global {
		title = "🏆 Classement global"
		counts, err = p.globalCounts(period, now)
	} else {
		title = "🏆 Classement du channel"
		counts, err = p.channelCounts(channelID, period, now)
	}
	if err != nil {
		return "", err
	}

	leaderboard := buildLeaderboard(counts, size, userID, p.usernameResolver())
	return formatLeaderboardMarkdown(title, period, leaderboard), nil
}

func periodLabel(period string) string {
	switch period {
	case PeriodWeek:
		return "7 derniers jours"
	case PeriodMonth:
		return "30 derniers jours"
	default:
		return "depuis le début"
	}
}

func formatLeaderboardMarkdown(title, period string, leaderboard Leaderboard) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "#### %s — %s\n\n", title, periodLabel(period))

	if len(leaderboard.Entries) == 0 {
		sb.WriteString("Aucun message compté pour le moment.")
		return sb.String()
	}

	sb.WriteString("| # | Utilisateur | Messages |\n|---|---|---|\n")
	for _, entry := range leaderboard.Entries {
		fmt.Fprintf(&sb, "| %d | @%s | %d |\n", entry.Rank, entry.Username, entry.Count)
	}
	if leaderboard.Me != nil {
		fmt.Fprintf(&sb, "| … | | |\n| %d | @%s (toi) | %d |\n", leaderboard.Me.Rank, leaderboard.Me.Username, leaderboard.Me.Count)
	}
	return sb.String()
}

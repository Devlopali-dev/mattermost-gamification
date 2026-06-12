import manifest from 'manifest';

export type LeaderboardEntry = {
    user_id: string;
    username: string;
    count: number;
    rank: number;
};

export type Board = {
    entries: LeaderboardEntry[];
    me?: LeaderboardEntry;
};

export type LeaderboardResponse = {
    period: string;
    channel?: Board;
    global: Board;
};

export type Period = 'week' | 'month' | 'all';

export function leaderboardUrl(channelId: string | null, period: Period): string {
    const params = new URLSearchParams({period});
    if (channelId) {
        params.set('channel_id', channelId);
    }
    return `/plugins/${manifest.id}/api/v1/leaderboard?${params.toString()}`;
}

export async function fetchLeaderboard(channelId: string | null, period: Period): Promise<LeaderboardResponse> {
    const response = await fetch(leaderboardUrl(channelId, period), {
        headers: {'X-Requested-With': 'XMLHttpRequest'},
    });
    if (!response.ok) {
        throw new Error(`Leaderboard request failed (${response.status})`);
    }
    return response.json();
}

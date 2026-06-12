import manifest from 'manifest';

import {leaderboardUrl} from './client';

describe('leaderboardUrl', () => {
    it('builds the plugin API url with channel and period', () => {
        const url = leaderboardUrl('chan1', 'week');

        expect(url).toBe(`/plugins/${manifest.id}/api/v1/leaderboard?period=week&channel_id=chan1`);
    });

    it('omits channel_id for the global leaderboard', () => {
        const url = leaderboardUrl(null, 'all');

        expect(url).toBe(`/plugins/${manifest.id}/api/v1/leaderboard?period=all`);
    });

    it('escapes query parameters', () => {
        const url = leaderboardUrl('a&b', 'month');

        expect(url).toContain('channel_id=a%26b');
        expect(url).toContain('period=month');
    });
});

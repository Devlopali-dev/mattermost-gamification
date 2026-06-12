import React, {useCallback, useEffect, useState} from 'react';
import {useSelector} from 'react-redux';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import BoardSection, {boardStyles} from './board_section';
import PeriodTabs from './period_tabs';

import type {LeaderboardResponse, Period} from '../client';
import {fetchLeaderboard} from '../client';

const styles: Record<string, React.CSSProperties> = {
    panel: {
        display: 'flex',
        flexDirection: 'column',
        height: '100%',
        overflowY: 'auto',
        padding: '12px 16px',
    },
    refresh: {
        marginTop: '16px',
        padding: '6px 12px',
        border: 'none',
        borderRadius: '4px',
        background: 'rgba(var(--center-channel-color-rgb), 0.08)',
        color: 'var(--center-channel-color)',
        cursor: 'pointer',
        alignSelf: 'flex-start',
    },
};

export default function LeaderboardPanel() {
    const channelId = useSelector(getCurrentChannelId);
    const [period, setPeriod] = useState<Period>('week');
    const [data, setData] = useState<LeaderboardResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const load = useCallback(async () => {
        if (!channelId) {
            return;
        }
        setLoading(true);
        setError(null);
        try {
            setData(await fetchLeaderboard(channelId, period));
        } catch (e) {
            setError('Impossible de charger le classement.');
        } finally {
            setLoading(false);
        }
    }, [channelId, period]);

    useEffect(() => {
        load();
    }, [load]);

    return (
        <div style={styles.panel}>
            <PeriodTabs
                period={period}
                onChange={setPeriod}
            />

            {loading && !data && <div style={boardStyles.muted}>{'Chargement…'}</div>}
            {error && <div style={boardStyles.muted}>{error}</div>}

            {data && (
                <>
                    <BoardSection
                        title={'Ce channel'}
                        board={data.channel}
                    />
                    <BoardSection
                        title={'Global'}
                        board={data.global}
                    />
                </>
            )}

            <button
                style={styles.refresh}
                onClick={load}
                disabled={loading}
            >
                {loading ? 'Chargement…' : 'Rafraîchir'}
            </button>
        </div>
    );
}

import React, {useCallback, useEffect, useState} from 'react';
import {useSelector} from 'react-redux';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import type {Board, LeaderboardResponse, Period} from '../client';
import {fetchLeaderboard} from '../client';

const PERIODS: Array<{key: Period; label: string}> = [
    {key: 'week', label: 'Semaine'},
    {key: 'month', label: 'Mois'},
    {key: 'all', label: 'Total'},
];

const styles: Record<string, React.CSSProperties> = {
    panel: {
        display: 'flex',
        flexDirection: 'column',
        height: '100%',
        overflowY: 'auto',
        padding: '12px 16px',
    },
    tabs: {
        display: 'flex',
        gap: '4px',
        marginBottom: '16px',
    },
    tab: {
        flex: 1,
        padding: '6px 0',
        border: '1px solid rgba(var(--center-channel-color-rgb), 0.16)',
        borderRadius: '4px',
        background: 'transparent',
        color: 'var(--center-channel-color)',
        cursor: 'pointer',
        fontWeight: 600,
    },
    tabActive: {
        background: 'var(--button-bg)',
        color: 'var(--button-color)',
        border: '1px solid var(--button-bg)',
    },
    sectionTitle: {
        margin: '12px 0 8px',
        fontSize: '14px',
        fontWeight: 700,
        textTransform: 'uppercase',
        opacity: 0.7,
    },
    row: {
        display: 'flex',
        alignItems: 'center',
        padding: '6px 4px',
        borderBottom: '1px solid rgba(var(--center-channel-color-rgb), 0.08)',
    },
    rank: {
        width: '28px',
        fontWeight: 700,
        opacity: 0.7,
    },
    username: {
        flex: 1,
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap',
    },
    count: {
        fontVariantNumeric: 'tabular-nums',
        fontWeight: 600,
    },
    meRow: {
        background: 'rgba(var(--center-channel-color-rgb), 0.06)',
        borderRadius: '4px',
        marginTop: '4px',
    },
    muted: {
        opacity: 0.7,
        padding: '8px 4px',
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

function medal(rank: number): string {
    switch (rank) {
    case 1:
        return '🥇';
    case 2:
        return '🥈';
    case 3:
        return '🥉';
    default:
        return String(rank);
    }
}

function BoardSection({title, board}: {title: string; board?: Board}) {
    if (!board || board.entries.length === 0) {
        return (
            <>
                <div style={styles.sectionTitle}>{title}</div>
                <div style={styles.muted}>{'Aucun message compté pour le moment.'}</div>
            </>
        );
    }
    return (
        <>
            <div style={styles.sectionTitle}>{title}</div>
            {board.entries.map((entry) => (
                <div
                    key={entry.user_id}
                    style={styles.row}
                >
                    <span style={styles.rank}>{medal(entry.rank)}</span>
                    <span style={styles.username}>{'@' + entry.username}</span>
                    <span style={styles.count}>{entry.count}</span>
                </div>
            ))}
            {board.me && (
                <div style={{...styles.row, ...styles.meRow}}>
                    <span style={styles.rank}>{board.me.rank}</span>
                    <span style={styles.username}>{'@' + board.me.username + ' (toi)'}</span>
                    <span style={styles.count}>{board.me.count}</span>
                </div>
            )}
        </>
    );
}

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
            <div style={styles.tabs}>
                {PERIODS.map((tab) => (
                    <button
                        key={tab.key}
                        style={{...styles.tab, ...(period === tab.key ? styles.tabActive : {})}}
                        onClick={() => setPeriod(tab.key)}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>

            {loading && !data && <div style={styles.muted}>{'Chargement…'}</div>}
            {error && <div style={styles.muted}>{error}</div>}

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

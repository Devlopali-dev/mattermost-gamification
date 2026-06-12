import React, {useCallback, useEffect, useState} from 'react';

import BoardSection, {boardStyles} from './board_section';
import PeriodTabs from './period_tabs';

import type {LeaderboardResponse, Period} from '../client';
import {fetchLeaderboard} from '../client';

const styles: Record<string, React.CSSProperties> = {
    page: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        height: '100%',
        overflowY: 'auto',
        background: 'var(--center-channel-bg)',
        color: 'var(--center-channel-color)',
    },
    content: {
        width: '100%',
        maxWidth: '640px',
        padding: '32px 24px',
    },
    title: {
        fontSize: '22px',
        fontWeight: 700,
        marginBottom: '4px',
    },
    subtitle: {
        opacity: 0.7,
        marginBottom: '24px',
    },
    refresh: {
        marginTop: '16px',
        padding: '6px 12px',
        border: 'none',
        borderRadius: '4px',
        background: 'rgba(var(--center-channel-color-rgb), 0.08)',
        color: 'var(--center-channel-color)',
        cursor: 'pointer',
    },
};

// Full-page product view shown from the product switcher menu.
export default function GamificationApp() {
    const [period, setPeriod] = useState<Period>('week');
    const [data, setData] = useState<LeaderboardResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            setData(await fetchLeaderboard(null, period));
        } catch (e) {
            setError('Impossible de charger le classement.');
        } finally {
            setLoading(false);
        }
    }, [period]);

    useEffect(() => {
        load();
    }, [load]);

    return (
        <div style={styles.page}>
            <div style={styles.content}>
                <div style={styles.title}>{'🏆 Leaderboard'}</div>
                <div style={styles.subtitle}>{'Classement global des messages sur tout le serveur.'}</div>

                <PeriodTabs
                    period={period}
                    onChange={setPeriod}
                />

                {loading && !data && <div style={boardStyles.muted}>{'Chargement…'}</div>}
                {error && <div style={boardStyles.muted}>{error}</div>}

                {data && (
                    <BoardSection
                        title={'Global'}
                        board={data.global}
                    />
                )}

                <button
                    style={styles.refresh}
                    onClick={load}
                    disabled={loading}
                >
                    {loading ? 'Chargement…' : 'Rafraîchir'}
                </button>
            </div>
        </div>
    );
}

import React from 'react';

import type {Board} from '../client';

export const boardStyles: Record<string, React.CSSProperties> = {
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
};

export function medal(rank: number): string {
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

export default function BoardSection({title, board}: {title: string; board?: Board}) {
    if (!board || board.entries.length === 0) {
        return (
            <>
                <div style={boardStyles.sectionTitle}>{title}</div>
                <div style={boardStyles.muted}>{'Aucun message compté pour le moment.'}</div>
            </>
        );
    }
    return (
        <>
            <div style={boardStyles.sectionTitle}>{title}</div>
            {board.entries.map((entry) => (
                <div
                    key={entry.user_id}
                    style={boardStyles.row}
                >
                    <span style={boardStyles.rank}>{medal(entry.rank)}</span>
                    <span style={boardStyles.username}>{'@' + entry.username}</span>
                    <span style={boardStyles.count}>{entry.count}</span>
                </div>
            ))}
            {board.me && (
                <div style={{...boardStyles.row, ...boardStyles.meRow}}>
                    <span style={boardStyles.rank}>{board.me.rank}</span>
                    <span style={boardStyles.username}>{'@' + board.me.username + ' (toi)'}</span>
                    <span style={boardStyles.count}>{board.me.count}</span>
                </div>
            )}
        </>
    );
}

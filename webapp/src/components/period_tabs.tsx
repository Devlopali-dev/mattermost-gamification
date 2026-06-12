import React from 'react';

import type {Period} from '../client';

export const PERIODS: Array<{key: Period; label: string}> = [
    {key: 'week', label: 'Semaine'},
    {key: 'month', label: 'Mois'},
    {key: 'all', label: 'Total'},
];

const styles: Record<string, React.CSSProperties> = {
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
};

type Props = {
    period: Period;
    onChange: (period: Period) => void;
};

export default function PeriodTabs({period, onChange}: Props) {
    return (
        <div style={styles.tabs}>
            {PERIODS.map((tab) => (
                <button
                    key={tab.key}
                    style={{...styles.tab, ...(period === tab.key ? styles.tabActive : {})}}
                    onClick={() => onChange(tab.key)}
                >
                    {tab.label}
                </button>
            ))}
        </div>
    );
}

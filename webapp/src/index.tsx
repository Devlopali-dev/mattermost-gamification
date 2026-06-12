// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import manifest from 'manifest';
import React from 'react';
import type {Store} from 'redux';

import type {GlobalState} from '@mattermost/types/store';

import type {PluginRegistry} from 'types/mattermost-webapp';

import GamificationApp from './components/gamification_app';
import LeaderboardPanel from './components/leaderboard_panel';

const trophyIcon = (
    <span
        aria-label='Leaderboard'
        style={{fontSize: '16px'}}
    >
        {'🏆'}
    </span>
);

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState>) {
        const {showRHSPlugin} = registry.registerRightHandSidebarComponent(
            LeaderboardPanel,
            'Leaderboard',
        );

        registry.registerChannelHeaderButtonAction(
            trophyIcon,
            () => store.dispatch(showRHSPlugin),
            'Leaderboard',
            'Voir le classement des messages',
        );

        // Entry in the product switcher menu (alongside Channels, Playbooks, …),
        // pointing to a full-page global leaderboard.
        registry.registerProduct(
            '/gamification',
            'trophy-outline',
            'Gamification',
            '/gamification',
            GamificationApp,
            () => null,
            undefined,
            false,
        );
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());

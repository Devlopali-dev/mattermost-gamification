# Mattermost Gamification

Plugin Mattermost qui compte les messages des utilisateurs (par channel et global) et affiche un leaderboard directement dans l'interface.

## Fonctionnalités

- **Comptage automatique** : chaque message humain est compté (réponses de threads incluses). Bots, webhooks, messages système et messages privés (DM/GM) sont exclus.
- **Leaderboard dans l'interface** : bouton 🏆 dans le header de channel → panneau latéral droit avec deux classements (channel courant + global).
- **Périodes** : Semaine (7 jours glissants), Mois (30 jours glissants), Total.
- **Rang personnel** : si vous êtes hors du top 10, votre rang s'affiche quand même.
- **Slash command** : `/leaderboard [global] [week|month|all]` — réponse éphémère en markdown (pratique sur mobile).

## Installation

1. Télécharger le bundle `com.devlopali.gamification-x.y.z.tar.gz` (ou le construire, voir ci-dessous).
2. **System Console → Plugins → Plugin Management → Upload Plugin**.
3. Activer le plugin.

Le comptage démarre à zéro à l'activation (pas de backfill de l'historique).

Requiert un serveur Mattermost self-hosted ≥ 9.0.

## Configuration

**System Console → Plugins → Gamification** :

| Paramètre | Description | Défaut |
|---|---|---|
| Channels exclus | IDs de channels à exclure du comptage (séparés par virgules) | vide |
| Taille du classement | Nombre d'utilisateurs affichés | 10 |

## Développement

Prérequis : Go ≥ 1.25, Node ≥ 20, make.

```bash
make dist          # build complet → dist/com.devlopali.gamification-x.y.z.tar.gz
make test          # tests Go + webapp
make check-style   # lint Go + ESLint
```

Déploiement direct sur un serveur local :

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=<token>
make deploy
```

## Architecture

- `server/` — plugin Go : hook `MessageHasBeenPosted` (comptage), API REST `/api/v1/leaderboard`, slash command, stockage KV (clés par channel + par jour, TTL 40 jours sur les clés journalières).
- `webapp/` — React/TypeScript : panneau RHS (onglets Semaine/Mois/Total), bouton de header de channel.
- Design détaillé : [docs/plans/2026-06-11-mattermost-gamification-design.md](docs/plans/2026-06-11-mattermost-gamification-design.md)

Basé sur le [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template).
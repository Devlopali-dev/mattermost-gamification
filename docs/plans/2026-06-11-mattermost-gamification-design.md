# Design — Plugin Mattermost de gamification (compteur de messages)

Date : 2026-06-11
Statut : validé

## Objectif

Compter les messages des utilisateurs (par channel et global) et afficher un
leaderboard dans un endroit visible de l'interface Mattermost.

## Décisions (validées avec l'utilisateur)

| Sujet | Décision |
|---|---|
| Architecture | Plugin Mattermost (serveur Go + webapp React), pas de bot externe |
| Périmètre v1 | Compteur + classement uniquement (pas de niveaux/badges) |
| Comptage | Tout message humain, threads inclus ; exclut bots, webhooks, messages système, DM/GM ; éditions et suppressions ignorées |
| Périodes | Semaine (7 jours glissants), Mois (30 jours glissants), Total |
| Historique | Départ à zéro à l'activation, pas de backfill |

## Architecture

### Serveur (Go)

- Hook `MessageHasBeenPosted` : filtre puis incrémente les compteurs.
- Filtres : `post.IsSystemMessage()`, prop `from_webhook`, user bot,
  channels de type `D`/`G`, channels exclus par config.
- Stockage : KV store du plugin (pas de table SQL custom).
- API REST `/plugins/{id}/api/v1/leaderboard` pour le webapp.
- Slash command `/leaderboard [global]` → réponse éphémère markdown.
- Job quotidien : purge des clés jour > 31 jours.

### Webapp (React/TypeScript)

- Bouton 🏆 dans le header de channel → ouvre le panneau RHS.
- RHS : onglets Semaine / Mois / Total ; bloc « Ce channel » + bloc « Global » ;
  top 10 + ligne rang perso si hors top ; refresh manuel.

### Build

Structure du mattermost-plugin-starter-template, `make dist` → bundle
`.tar.gz` à uploader dans System Console → Plugins. Self-hosted requis.

## Modèle de données (KV store)

```
total:global:{userID}                 → int
total:chan:{channelID}:{userID}       → int
day:{YYYY-MM-DD}:{channelID}:{userID} → int
```

- Écriture : 3 incréments par post, compare-and-set atomique
  (`KVSetWithOptions{Atomic: true}`, retry ×3). Sûr en cluster HA.
- Lecture Total : scan préfixe + tri desc + top N.
- Lecture Semaine/Mois : somme des clés `day:` sur 7/30 dates glissantes.
- Rétention : purge quotidienne des clés `day:` > 31 jours.
- Limite assumée : scan préfixe à chaque lecture, OK jusqu'à quelques
  milliers d'utilisateurs actifs ; v2 = table SQL si besoin.

## Sécurité

Endpoints REST : vérifient `Mattermost-User-Id` et l'appartenance au channel
demandé. Pas de stats d'un channel privé sans en être membre.

## Configuration (System Console)

- Channels exclus (liste d'IDs).
- Taille du top (défaut 10).

## Tests

- Go : filtres de comptage, agrégation périodes, CAS retry (plugintest).
- Webapp : rendu RHS (top 10, onglets) avec Jest.

## Hors périmètre v1

Niveaux/badges/streaks, backfill historique, live-update WebSocket,
anti-spam (longueur min, cooldown).
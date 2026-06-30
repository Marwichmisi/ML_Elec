# Phase 1: Core Foundation — Specification

**Created:** 2026-06-30
**Ambiguity score:** 0.182 (gate: ≤ 0.20)
**Requirements:** 6 locked

## Goal

Un binaire Go unique démarre un bus NATS embedded, gère des plugins via child processes (JSON-RPC stdin/stdout), stocke des données capteurs en SQLite WAL, expose une API REST, et charge une configuration centralisée — le tout sous 5000 LOC de code core.

## Background

Projet greenfield — aucun code Go n'existe encore. Le projet ML_Elec vise la collecte de données capteurs en temps réel et la détection d'anomalies pour la maintenance prédictive d'équipements électriques. La Phase 1 construit le noyau infrastructurel sur lequel les phases suivantes (Plugin SDK, MQTT, Détection, Dashboard) s'appuient.

Le noyau doit être un micro-noyau strict : le core ne contient que l'infrastructure (bus, stockage, API, config, lifecycle plugins) — aucune logique métier dans le core.

## Requirements

1. **NATS Embedded**: Le core intègre un bus NATS embedded pour la communication interne entre composants.
   - Current: Aucun bus de message n'existe
   - Target: NATS embedded démarre avec le core, fournit pub/sub interne, configuration par défaut (pas de port réseau exposé)
   - Acceptance: Le core démarre et NATS est opérationnel — un publish/subscribe entre deux goroutines internes fonctionne sans erreur

2. **Plugin Manager**: Le core gère des plugins Python via child processes avec isolation, communiquant via JSON-RPC sur stdin/stdout (pattern HashiCorp go-plugin).
   - Current: Aucun gestionnaire de plugins
   - Target: Plugin manager lance un mock Python en child process, échange des messages JSON-RPC sur stdin/stdout, isole les défaillances (un plugin crash ne fait pas crasher le core)
   - Acceptance: Un plugin mock Python démarre, envoie un message JSON-RPC, reçoit une réponse, et si le plugin crash, le core continue de fonctionner

3. **REST API**: Le core expose une API REST pour l'interaction externe (health check, données capteurs).
   - Current: Aucune API
   - Target: API REST sur un port configurable avec au minimum un endpoint health check et un endpoint de query de données capteurs
   - Acceptance: `GET /health` retourne 200 OK ; les requêtes concurrentes sont gérées sans data race (testé avec `-race`)

4. **SQLite WAL**: Le core stocke les données capteurs en SQLite avec mode WAL pour les séries temporelles.
   - Current: Aucun stockage
   - Target: SQLite en mode WAL, écriture et lecture de données capteurs avec timestamps, gestion propre de la corruption
   - Acceptance: Écriture et relecture d'un enregistrement capteur fonctionne ; si le fichier DB est corrompu, le core exit avec code 1 et log d'erreur

5. **Configuration centralisée**: Un fichier de configuration contrôle le core et les plugins (activables/désactivables).
   - Current: Aucune configuration
   - Target: Fichier de config (YAML ou TOML) chargé au démarrage, contrôle les ports, les plugins activés/désactivés, et les paramètres du core
   - Acceptance: Le core démarre avec le fichier de config ; un plugin désactivé dans la config ne est pas lancé ; si le fichier est absent, le core démarre avec des valeurs par défaut

6. **Core < 5000 LOC**: Le code source du core (hors tests et code généré) reste sous 5000 lignes de code.
   - Current: Aucun code
   - Target: Toute la logique infrastructure du core (NATS, plugins, REST, SQLite, config) tient dans < 5000 LOC
   - Acceptance: `find . -name '*.go' ! -name '*_test.go' ! -path './generated/*' | xargs wc -l` retourne < 5000

## Boundaries

**In scope:**
- Binaire Go unique contenant NATS embedded, REST API, SQLite WAL, plugin manager, configuration
- Plugin mock Python pour validation (JSON-RPC stdin/stdout)
- Endpoints REST : health check + query données capteurs
- Fichier de configuration (YAML ou TOML)
- Gestion d'erreur propre (exit codes, logs stderr)
- Tests unitaires et d'intégration

**Out of scope:**
- SDK de plugin (CORE-06) — Phase 2
- Plugin MQTT réel (ACQ-01 à ACQ-04) — Phase 2
- Détection d'anomalies (DET-01 à DET-04) — Phase 3
- Dashboard frontend — Phase 4
- Authentification API (OAuth, JWT) — pas critique pour MVP industriel
- Déploiement Docker/Kubernetes — Phase 7
- Monitoring/observabilité du core lui-même — Phase 4

## Constraints

- **Langage** : Go uniquement pour le core (pas de dépendances C si possible)
- **NATS** : Configuration par défaut, pas de port réseau exposé en Phase 1
- **SQLite** : Mode WAL obligatoire pour les performances d'écriture concurrente
- **Plugin communication** : JSON-RPC sur stdin/stdout uniquement (pas de socket réseau)
- **LOC budget** : < 5000 LOC pour le code source core (hors tests, hors code généré)
- **Dépendances** : Utiliser les packages Go standard + NATS + SQLite driver uniquement pour le core
- **Go skills** : L'agent développeur DOIT utiliser les skills Go (golang-*) disponibles pour suivre les bonnes pratiques (error handling, naming, structs/interfaces, testing, lint)

## Acceptance Criteria

- [ ] Le binaire Go compile et démarre sans erreur
- [ ] NATS embedded est opérationnel (pub/sub interne fonctionne)
- [ ] Plugin mock Python est lancé en child process et communique via JSON-RPC
- [ ] Un plugin crash ne fait pas crasher le core (isolation testée)
- [ ] `GET /health` retourne 200 OK
- [ ] L'API REST supporte des requêtes concurrentes sans data race (`go test -race`)
- [ ] SQLite stocke et relit des données capteurs en mode WAL
- [ ] Si le fichier DB est corrompu, le core exit avec code 1
- [ ] Le fichier de configuration est chargé au démarrage
- [ ] Un plugin désactivé dans la config n'est pas lancé
- [ ] Si le fichier de config est absent, le core démarre avec des valeurs par défaut
- [ ] Le core fait < 5000 LOC (hors tests et code généré)
- [ ] Tous les tests passent (`go test ./...`)
- [ ] Pas de data race détecté (`go test -race ./...`)

## Edge Coverage

**Coverage:** 6/6 applicable edges resolved · 0 unresolved

| Category | Requirement | Status | Resolution / Reason |
|----------|-------------|--------|---------------------|
| startup-failure | R1 | ✅ covered | NATS failure → exit code 1 + stderr message (Acceptance: core démarre et NATS opérationnel) |
| crash-isolation | R2 | ✅ covered | Plugin crash → core continue, log warning (Acceptance: un plugin crash ne fait pas crasher le core) |
| concurrency | R3 | ✅ covered | Requêtes concurrentes → pas de data race (Acceptance: go test -race passe) |
| corruption | R4 | ✅ covered | DB corrompue → exit code 1 + log (Acceptance: Si fichier DB corrompu, core exit code 1) |
| missing-config | R5 | ✅ covered | Config absente → valeurs par défaut + log warning (Acceptance: core démarre avec defaults) |
| loc-counting | R6 | ✅ covered | Tests et code généré exclus du comptage LOC (Acceptance: wc -l < 5000 hors tests) |

## Prohibitions (must-NOT)

**Coverage:** 6/6 applicable prohibitions resolved · 0 unresolved

| Prohibition (must-NOT statement) | Requirement | Status | Verification / Reason |
|----------------------------------|-------------|--------|------------------------|
| Le bus NATS embedded ne doit PAS persister les messages au-delà de la durée de vie du processus | R1 | resolved | verification: test — test vérifie que NATS est éphémère |
| Les plugins ne doivent PAS pouvoir communiquer entre eux directement (tout passe par le core) | R2 | resolved | verification: test — test vérifie l'isolation des plugins |
| L'API REST ne doit PAS exposer les sujets NATS internes directement aux clients externes | R3 | resolved | verification: test — test vérifie que les réponses API ne contiennent pas de NATS internals |
| SQLite ne doit PAS utiliser de requêtes SQL non paramétrisées | R4 | resolved | verification: lint — rule `no-raw-sql` détecte les requêtes non paramétrisées |
| La configuration ne doit PAS stocker de secrets en texte brut | R5 | resolved | verification: test — test scanne le fichier config pour les patterns "password", "secret", "key" |
| Les fichiers de test et code généré ne comptent PAS dans la limite de 5000 LOC | R6 | resolved | verification: judgment — vérification manuelle du comptage LOC |

## Ambiguity Report

| Dimension          | Score | Min  | Status | Notes                                    |
|--------------------|-------|------|--------|------------------------------------------|
| Goal Clarity       | 0.85  | 0.75 | ✓      | Architecture claire : binaire unique, NATS embedded |
| Boundary Clarity   | 0.80  | 0.70 | ✓      | In/out scope explicite dans REQUIREMENTS |
| Constraint Clarity | 0.75  | 0.65 | ✓      | LOC budget, NATS config, SQLite WAL      |
| Acceptance Criteria| 0.85  | 0.70 | ✓      | 14 critères pass/fail                    |
| **Ambiguity**      | 0.182 | ≤0.20| ✓      |                                          |

## Interview Log

| Round | Perspective     | Question summary                                      | Decision locked                                          |
|-------|-----------------|-------------------------------------------------------|----------------------------------------------------------|
| 1     | Researcher      | État du code existant ? Déclencheur ?                  | Projet greenfield, zéro code existant                    |
| 2     | Simplifier      | Version minimale ? NATS embedded obligatoire ?         | Un seul binaire Go, NATS embedded obligatoire            |
| 3     | Boundary Keeper | Hors scope ? Critère de 'done' ?                      | Roadmap suffit pour exclusions, 'done' = tout fonctionne ensemble |
| 4     | Failure Analyst | Risque principal ? Config NATS ?                      | NATS crash = risque #1, config par défaut                |

---

*Phase: 01-core-foundation*
*Spec created: 2026-06-30*
*Next step: /gsd-discuss-phase 1 — implementation decisions (how to build what's specified above)*

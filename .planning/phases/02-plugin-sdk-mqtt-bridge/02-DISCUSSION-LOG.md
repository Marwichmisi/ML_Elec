# Phase 02: Plugin SDK & MQTT Bridge - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-07-03
**Phase:** 02-plugin-sdk-mqtt-bridge
**Areas discussed:** Broker MQTT, Structure proto & SDK, Format binaire MQTT, Stratégie de tests, Gestion erreurs MQTT, Structure config.yaml, API REST assets, Migration config

---

## Broker MQTT

| Option | Description | Selected |
|--------|-------------|----------|
| aler9/mqtt (Recommended) | Broker MQTT complet en Go, QoS 0/1/2, léger | |
| mellium/mqtt | Implémentation MQTT Go plus bas-niveau | |
| Autre broker | Libre choix | ✓ (Eclipse Paho Go) |

| Port local fixe (Recommended) | Broker écoute sur port fixe (ex: 1883) | ✓ |
| Port configurable | Port défini dans config.yaml | |

| Même processus (Recommended) | Broker + plugin = un seul binaire Go | |
| Processus séparé | Broker dans binaire séparé | ✓ |

| Plugin lance le broker (Recommended) | Plugin démarre broker comme sous-processus | |
| Broker externe fixe | Broker installé séparément (Mosquitto) | |
| Tu décides | L'agent choisit | ✓ |

**User's choice:** Eclipse Paho Go, port fixe, processus séparé, cycle de vie à décider par l'agent
**Notes:** L'utilisateur a choisi un broker séparé malgré la recommandation "même processus" — plus d'isolation mais plus de complexité

---

## Structure proto & SDK

| Option | Description | Selected |
|--------|-------------|----------|
| Un par service (Recommended) | lifecycle.proto + sensor.proto | ✓ |
| Un seul fichier | plugin.proto avec tous les services | |
| Tu décides | Libre choix | |

| Commit (Recommended) | Code généré versionné dans le repo | ✓ |
| Généré au build | Pas de code généré dans le repo | |

| Go plugin only (Recommended) | Plugin MQTT = binaire Go uniquement | |
| Go + Python plugin | Support Python dès maintenant | ✓ |

| Wrapper simplifié (Recommended) | Client Python avec contrat gRPC unique | ✓ (variant: wrapper simplifié + contrat unique) |
| Mêmes protos | Plugin Python implémente mêmes interfaces | |
| Tu décides | Libre choix | |

| Package unique v1 (Recommended) | Un seul package pkg/sdk/ | |
| Versioning dès v1 | pkg/sdk/v1/ | ✓ |
| Tu décides | Libre choix | |

| Package séparé (Recommended) | pkg/sdk/v1/ et pkg/sdk/v2/ coexistent | |
| Tags Git | Un seul package, versioning via tags | ✓ (Tags Git + nouveau module /v2 sur breaking change) |

**User's choice:** Un proto par service, commit, Go + Python (wrapper simplifié), versioning dès v1, tags Git
**Notes:** L'utilisateur veut un wrapper Python simplifié avec contrat gRPC unique maintenu par la plateforme. Versioning: tags Git pour toutes les versions, nouveau module /v2 uniquement sur breaking change.

---

## Format binaire MQTT

| Option | Description | Selected |
|--------|-------------|----------|
| 1 kHz (Recommended) | 1000 échantillons/seconde | |
| 10 kHz | Plus précis mais charge lourde | |
| Configurable | Fréquence définie dans config.yaml | ✓ |
| Tu décides | Libre choix | ✓ (Configurable + profils) |

| Header fixe + payload (Recommended) | 4-8 octets en-tête + données float32 | ✓ (24-32 octets) |
| Structure complexe | En-tête avec checksum, version format | |
| Tu décides | Libre choix | |

| Format simple (Recommended) | {temperature, vibration, timestamp} | |
| Format niché | Structure avec device_id, readings array | |
| Tu décides | Libre choix | ✓ (2 types: telemetry + status) |

| Hiérarchique ISA-95 (Recommended) | factory/{site}/{area}/{line}/{asset}/{type} | ✓ |
| Plus simple | esp32/{device_id}/telemetry | |
| Tu décides | Libre choix | |

**User's choice:** Configurable avec profils prédéfinis (1-25.6 kHz), header 24-32 octets avec version, 2 types JSON (telemetry + status), topics ISA-95/UNS
**Notes:** L'utilisateur a fourni une structure de topics complète inspirée ISA-95/UNS avec exemples concrets. Format binaire: header fixe, version 1er octet, encoding int16 défaut, pas de CRC v1. Profils: low, standard, advanced, bearing.

---

## Stratégie de tests

| Option | Description | Selected |
|--------|-------------|----------|
| Broker embarqué en test (Recommended) | Vrai broker Paho en integration tests | |
| Mock broker | Simuler le broker en unit tests | |
| Les deux | Unit + integration | ✓ |
| Tu décides | Libre choix | |

| Test Go (Recommended) | Test Go qui lance plugin, tue, vérifie core | ✓ |
| Script shell | Script bash pour kill manuel | |
| Les deux | Test Go + script shell | ✓ |
| Tu décides | Libre choix | |

| Benchmark Go (Recommended) | Benchmark Go avec timer | ✓ |
| Outil externe | mqtt-stresser ou similaire | |
| Les deux | Benchmark + outil externe | |
| Tu décides | Libre choix | |

**User's choice:** Les deux (mock + vrai broker), test Go principal + script shell, benchmark Go
**Notes:** Architecture de tests bien définie: mock en unit, vrai broker en integration, crash isolation en Go, perf en benchmark Go

---

## Gestion erreurs MQTT

| Option | Description | Selected |
|--------|-------------|----------|
| Backoff exponentiel + LWT (Recommended) | Reconnexion auto avec backoff (1s→30s), LWT, clean session=false | ✓ |
| Backoff simple | Intervalle fixe (5s) | |
| Tu décides | Libre choix | |

**User's choice:** Backoff exponentiel + LWT + clean session=false
**Notes:** Utilisateur a demandé d'appliquer toutes les recommandations du skill MQTT development (topics hiérarchiques, wildcards, QoS, sessions, keep-alive, LWT)

---

## Structure config.yaml

| Option | Description | Selected |
|--------|-------------|----------|
| Sections séparées (Recommended) | mqtt: {}, validation: {}, assets: {} | ✓ |
| Imbrication | mqtt: {broker: {}, validation: {}} | |
| Tu décides | Libre choix | |

| Port, topics, QoS, profils (Recommended) | Les essentials | |
| Plus complet | Ajouter keep_alive, clean_session, session_expiry | |
| Tu décides | Libre choix | ✓ (Configuration par couches) |

| Température, vibration, courant (Recommended) | min/max pour 3 capteurs | |
| Plus étendu | Ajouter humidité, pression, etc. | |
| Tu décides | Libre choix | ✓ (Validation 3 niveaux) |

**User's choice:** Sections séparées, configuration par couches, validation 3 niveaux
**Notes:** L'utilisateur veut un système de validation robuste: plages valeurs + qualité données + timestamps/complétude/taille/santé capteurs. Configuration par couches: essentiels visibles, avancés regroupés.

---

## API REST assets

| Option | Description | Selected |
|--------|-------------|----------|
| JSON standard (Recommended) | {data: ..., error: ...} cohérent Phase 1 | |
| Format industriel | Structure verbeuse avec metadata, HATEOAS | |
| Tu décides | Libre choix | ✓ (JSON + pagination via skill REST) |

| JSON + pagination (Recommended) | {data: [...], pagination: {page, limit, total}} | ✓ |
| JSON simple | {data: [...]} sans pagination | |
| Tu décides | Libre choix | |

| Endpoint dédié (Recommended) | GET /api/v1/assets/{id}/sensors | ✓ (Les deux: dédié + collection globale) |
| 嵌套 dans GET /assets | GET /api/v1/assets inclut sensors | |

**User's choice:** JSON + pagination, endpoints imbriqués + collection globale avec pagination
**Notes:** L'utilisateur a inspiré du skill REST API design. Recommande les deux: endpoint relation (assets/{id}/sensors) et collection globale (/sensors) avec filtrage et pagination sur les deux.

---

## Migration config

| Option | Description | Selected |
|--------|-------------|----------|
| Valeurs par défaut (Recommended) | Code gère champs manquants avec defaults | ✓ |
| Script de migration | Script transforme ancien → nouveau config | |
| Tu décides | Libre choix | |

| Alias + warning (Recommended) | Ancien nom fonctionne + warning dépréciation | ✓ |
|Erreur immédiate | Démarrage échoue si ancien champ détecté | |
| Tu décides | Libre choix | |

**User's choice:** Valeurs par défaut, alias + warning pour rétrocompatibilité
**Notes:** L'utilisateur a explicitement choisi l'approche semver: alias fonctionne + warning en v1, suppression en v2. Cohérent avec le versioning sémantique.

---

## the agent's Discretion

- Cycle de vie exact du broker MQTT (comment le plugin le démarre/arrête)
- Paramètres de backoff exponentiel (factor, max retries)
- Structure exacte des headers binaires (24 vs 32 octets)
- Organisation des fichiers dans chaque package
- Configuration des paramètres de logging
- Choix du logger (slog standard, zerolog, zap)

## Deferred Ideas

None — discussion stayed within phase scope

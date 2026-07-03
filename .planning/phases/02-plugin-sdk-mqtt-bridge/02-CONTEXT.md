# Phase 02: Plugin SDK & MQTT Bridge - Context

**Gathered:** 2026-07-03
**Status:** Ready for planning

<domain>
## Phase Boundary

SDK gRPC versionné pour plugins (lifecycle + sensor services), plugin MQTT externe avec broker Paho séparé, validation 3 niveaux des données capteurs, registry d'assets ISA-95/UNS avec API REST, et support des payloads JSON + binaires MQTT — le tout réspectant l'architecture micro-noyau et les skills Go obligatoires.

</domain>

<spec_lock>
## Requirements (locked via SPEC.md)

**5 requirements are locked.** See `02-SPEC.md` for full requirements, boundaries, and acceptance criteria.

Downstream agents MUST read `02-SPEC.md` before planning or implementing. Requirements are not duplicated here.

**In scope (from SPEC.md):**
- `pkg/sdk/` — gRPC plugin SDK with `.proto` files and generated Go code
- `cmd/mqtt-plugin/` — MQTT bridge plugin (external child process)
- Embedded MQTT broker (QoS 0+1) within the MQTT plugin
- Data validation (configurable ranges, timestamp monotonicity)
- Asset registry (SQLite tables + REST API)
- JSON + binary MQTT payload parsing
- Dynamic MQTT topic discovery (wildcard subscriptions)
- Plugin manager upgrade to support gRPC transport
- Config.yaml additions for MQTT, validation rules, asset configuration

**Out of scope (from SPEC.md):**
- Dashboard UI (Phase 4) — no frontend work in this phase
- Anomaly detection engine (Phase 3) — no detection logic
- Alert notifications (Phase 3/5) — no email/in-app alerts
- QoS 2 support — QoS 0+1 sufficient for industrial sensors
- Modbus/OPC UA plugins (v2 requirements) — MQTT only for now
- Python plugin SDK — gRPC proto supports it but no Python SDK wrapper in this phase
- MQTT authentication/TLS — plaintext for v1, security deferred
- Plugin auto-discovery — plugins listed in config.yaml

</spec_lock>

<decisions>
## Implementation Decisions

### Broker MQTT
- **D-01:** Lib broker: Eclipse Paho Go — broker MQTT embarqué dans le plugin
- **D-02:** Port local fixe (1883) — les ESP32 se connectent directement
- **D-03:** Broker dans un processus séparé du plugin MQTT — le plugin MQTT démarre le broker comme sous-processus
- **D-04:** Cycle de vie broker: l'agent décide de la meilleure approche (plugin lance le broker)

### Structure proto & SDK
- **D-05:** Un fichier .proto par service: `lifecycle.proto` (Init/Start/Stop) + `sensor.proto` (Collect)
- **D-06:** Code généré commit dans le repo — les plugins n'ont pas besoin de protoc pour compiler
- **D-07:** Support Go + Python plugin via go-plugin dès v1
- **D-08:** Plugin Python: wrapper simplifié avec contrat gRPC unique maintenu par la plateforme
- **D-09:** Versioning dès v1: `pkg/sdk/v1/` — package unique pour v1
- **D-10:** Breaking changes futurs: tags Git + nouveau module `/v2` uniquement sur breaking change

### Format binaire MQTT
- **D-11:** Fréquence d'échantillonnage vibration: configurable avec profils prédéfinis (low, standard, advanced, bearing), défaut 1 kHz, max 25.6 kHz
- **D-12:** Structure trame binaire v1: header fixe 24-32 octets, payload brut, version en 1er octet, encoding configurable (int16 défaut), pas de CRC v1
- **D-13:** Payloads JSON: telemetry `{ts, values:{temperature, humidity, current}}` + status `{ts, status, battery, rssi, firmware}`
- **D-14:** Topics MQTT: hiérarchie ISA-95/UNS `factory/{site}/{area}/{line}/{asset}/{type}` + `sys/{component}/{type}`, wildcard-friendly

### Stratégie de tests
- **D-15:** Tests plugin MQTT: mock broker en unit tests + vrai broker Paho en integration tests
- **D-16:** Tests isolation crash: test Go principal (lance plugin, tue, vérifie core continue) + script shell pour contrôle manuel
- **D-17:** Tests performance: benchmark Go avec timer, mesure broker→NATS (100 msg/s, <100ms)

### Gestion erreurs MQTT
- **D-18:** Reconnexion: backoff exponentiel (1s→30s) + LWT pour status offline + clean session=false
- **D-19:** Toutes les bonnes pratiques MQTT skill appliquées: topics hiérarchiques, wildcards, QoS, sessions, keep-alive

### Structure config.yaml
- **D-20:** Sections séparées: `mqtt: {}`, `validation: {}`, `assets: {}` — propre, modulaire
- **D-21:** Configuration par couches: paramètres essentiels visibles, paramètres avancés regroupés
- **D-22:** Validation 3 niveaux: plages valeurs + qualité données + timestamps/complétude/taille/santé capteurs

### API REST assets
- **D-23:** Réponses JSON + pagination: `{data: [...], pagination: {page, limit, total}}`
- **D-24:** Endpoints: `GET /api/v1/assets/{id}/sensors` (relation) + `GET /api/v1/sensors` (collection globale) avec pagination sur les deux
- **D-25:** Conventions REST: noms pluriels, status codes standard (201, 400, 404, 409), DELETE → 405

### Migration config
- **D-26:** Valeurs par défaut: le code gère les champs manquants avec des defaults
- **D-27:** Rétrocompatibilité: alias + warning pour champs renommés, suppression en v2 uniquement

### the agent's Discretion
- Cycle de vie exact du broker MQTT (comment le plugin le démarre/arrête)
- Paramètres de backoff exponentiel (factor, max retries)
- Structure exacte des headers binaires (24 vs 32 octets)
- Organisation des fichiers dans chaque package
- Configuration des paramètres de logging
- Choix du logger (slog standard, zerolog, zap)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Exigences et spécifications
- `.planning/phases/02-plugin-sdk-mqtt-bridge/02-SPEC.md` — 5 exigences verrouillées, critères d'acceptation, boundaries
- `.planning/REQUIREMENTS.md` — 24 requirements v1, traceability matrix
- `.planning/PROJECT.md` — Contexte projet, key decisions, constraints
- `.planning/ROADMAP.md` — 7 phases, dependency graph, pitfall mitigations

### Architecture et patterns
- `.planning/STATE.md` — État actuel du projet, décisions accumulées
- `.planning/phases/01-core-foundation/01-CONTEXT.md` — Décisions Phase 1 (layout, REST, SQLite, DI)

### Skills obligatoires
- `.agents/skills/golang-error-handling/SKILL.md` — Gestion des erreurs Go
- `.agents/skills/golang-naming/SKILL.md` — Conventions de nommage Go
- `.agents/skills/golang-structs-interfaces/SKILL.md` — Structs et interfaces
- `.agents/skills/golang-testing/SKILL.md` — Tests Go
- `.agents/skills/golang-lint/SKILL.md` — Configuration lint Go
- `.agents/skills/golang-how-to/SKILL.md` — Orchestrateur de skills Go
- `.agents/skills/golang-grpc/SKILL.md` — Patterns gRPC production-ready
- `.agents/skills/sqlite-database-expert/SKILL.md` — SQLite embedded, migrations, FTS
- `.agents/skills/mqtt-development/SKILL.md` — Bonnes pratiques MQTT IoT
- `.agents/skills/rest-api-design/SKILL.md` — Conventions REST API
- `.agents/skills/validate-data/SKILL.md` — Validation et QA des données

Aucun ADR externe n'existe — les exigences sont complètement capturées dans SPEC.md et les décisions ci-dessus.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/plugin/manager.go` — plugin manager go-plugin existant, à étendre pour gRPC
- `internal/storage/storage.go` — SQLite CRUD, tables `sensor_readings` existantes
- `internal/config/config.go` — YAML config avec `plugins.enabled`, à étendre pour MQTT/validation/assets
- `cmd/mock-plugin/main.go` — mock plugin à mettre à jour pour nouveau SDK gRPC
- Core ~1400 LOC — marge pour ajouter SDK + validation (limite 5000 LOC)

### Established Patterns
- Architecture micro-noyau: core = infrastructure uniquement, pas de logique métier
- Communication JSON-RPC stdin/stdout: pattern HashiCorp go-plugin (à étendre pour gRPC)
- NATS embedded: bus interne éphémère, pas de port réseau exposé
- DI avec google/wire: initialisation compile-time
- REST: net/http standard, pas de framework
- SQLite: modernc.org/sqlite (pure Go), WAL mode

### Integration Points
- Plugin manager → upgrade pour support gRPC transport (coexistence net/rpc + gRPC)
- Config → nouvelles sections mqtt, validation, assets
- Storage → nouvelles tables assets, asset_sensors
- REST API → nouveaux endpoints /api/v1/assets, /api/v1/sensors
- NATS → plugin MQTT publie sur `sensor.*` subjects

</code_context>

<specifics>
## Specific Ideas

- Hiérarchie topics ISA-95/UNS: `factory/{site}/{area}/{line}/{asset}/{type}` — conçue pour l'ensemble de la plateforme (ESP32, PLC, OPC UA, dashboards, IA)
- Profils de fréquence vibration: low, standard, advanced, bearing — évite les erreurs de configuration
- Validation 3 niveaux: pas juste des seuils mais qualité données, complétude, taille, santé capteurs
- Configuration par couches: essentiels visibles, avancés regroupés — compromis simplicité/flexibilité
- Skills Go obligatoires pour l'agent dev — mentionnés dans CONTEXT.md ET dans AGENTS.md
- L'installation Go 1.26.4 est disponible sur la machine
- Skills golang, sqlite, mqtt-development, golang-grpc, validate-data, rest-api-design installés et disponibles

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 02-plugin-sdk-mqtt-bridge*
*Context gathered: 2026-07-03*

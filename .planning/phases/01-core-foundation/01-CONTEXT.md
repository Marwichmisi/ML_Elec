# Phase 01: Core Foundation - Context

**Gathered:** 2026-06-30
**Status:** Ready for planning

<domain>
## Phase Boundary

Un binaire Go unique démarre un bus NATS embedded, gère des plugins via child processes (JSON-RPC stdin/stdout), stocke des données capteurs en SQLite WAL, expose une API REST, et charge une configuration centralisée — le tout sous 5000 LOC de code core.

</domain>

<spec_lock>
## Requirements (locked via SPEC.md)

**6 requirements are locked.** See `01-SPEC.md` for full requirements, boundaries, and acceptance criteria.

Downstream agents MUST read `01-SPEC.md` before planning or implementing. Requirements are not duplicated here.

**In scope (from SPEC.md):**
- Binaire Go unique contenant NATS embedded, REST API, SQLite WAL, plugin manager, configuration
- Plugin mock Python pour validation (JSON-RPC stdin/stdout)
- Endpoints REST : health check + query données capteurs
- Fichier de configuration (YAML ou TOML)
- Gestion d'erreur propre (exit codes, logs stderr)
- Tests unitaires et d'intégration

**Out of scope (from SPEC.md):**
- SDK de plugin (CORE-06) — Phase 2
- Plugin MQTT réel (ACQ-01 à ACQ-04) — Phase 2
- Détection d'anomalies (DET-01 à DET-04) — Phase 3
- Dashboard frontend — Phase 4
- Authentification API (OAuth, JWT) — pas critique pour MVP industriel
- Déploiement Docker/Kubernetes — Phase 7
- Monitoring/observabilité du core lui-même — Phase 4

</spec_lock>

<decisions>
## Implementation Decisions

### Layout du projet Go
- **D-01:** Structure `cmd/ + internal/` — binaire dans `cmd/ml-elec/main.go`, composants dans `internal/`
- **D-02:** 5 packages dans internal/: `nats`, `plugin`, `api`, `storage`, `config` — un par composant du SPEC
- **D-03:** Config par défaut: `./config.yaml` — recherche dans répertoire courant, puis `~/.config/ml-elec/config.yaml`, puis valeurs par défaut
- **D-04:** Tests unitaires à côté des packages (`_test.go` dans le même dossier)
- **D-05:** DI avec google/wire pour l'initialisation — code généré à compile-time, pas de反射
- **D-06:** Module Go unique dans la racine (pas de go.work pour v1)

### Framework REST
- **D-07:** net/http standard pur — pas de framework externe, zéro dépendance
- **D-08:** Handlers dans `internal/api/` — fichiers par endpoint (health.go, sensors.go)
- **D-09:** Documentation Swagger/OpenAPI via swaggo — annotations Go → spec OpenAPI
- **D-10:** Réponses JSON standard avec structure `{"data": ..., "error": ...}`

### SQLite driver
- **D-11:** modernc.org/sqlite — driver pure Go, pas de CGO, cross-compilation facile
- **D-12:** Builder de requêtes (squirrel) pour la construction SQL paramétrée
- **D-13:** Migrations avec bibliothèque externe (golang-migrate/migrate)
- **D-14:** Vér corruption: au démarrage (PRAGMA integrity_check) + périodique

### Cycle de vie & shutdown
- **D-15:** context.Context + signal SIGINT/SIGTERM pour l'arrêt gracieux
- **D-16:** Ordre d'arrêt LIFO: API → Plugins → NATS → Storage
- **D-17:** Timeout d'arrêt: 30 secondes avant forçage
- **D-18:** Retry avec backoff au démarrage des composants

### Skills Go obligatoires
- **D-19:** L'agent développeur DOIT charger les skills golang-* pour: error handling, naming, structs/interfaces, testing, lint
- **D-20:** Les skills Go seront mentionnés dans CONTEXT.md ET dans AGENTS.md pour double sécurité

### the agent's Discretion
- Organisation des fichiers dans chaque package (l'agent dev choisit)
- Logique de retry的具体参数 (backoff exponentiel, max retries)
- Configuration des paramètres de logging
- Structure exacte des réponses JSON d'erreur
- Choix du logger (slog standard, zerolog, zap)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Exigences et spécifications
- `.planning/phases/01-core-foundation/01-SPEC.md` — 6 exigences verrouillées, critères d'acceptation, boundaries
- `.planning/REQUIREMENTS.md` — 24 requirements v1, traceability matrix
- `.planning/PROJECT.md` — Contexte projet, key decisions, constraints
- `.planning/ROADMAP.md` — 7 phases, dependency graph, pitfall mitigations

### Architecture et patterns
- `.planning/STATE.md` — État actuel du projet, décisions accumulées

### Skills obligatoires
- `.agents/skills/golang-error-handling/SKILL.md` — Gestion des erreurs Go
- `.agents/skills/golang-naming/SKILL.md` — Conventions de nommage Go
- `.agents/skills/golang-structs-interfaces/SKILL.md` — Structs et interfaces
- `.agents/skills/golang-testing/SKILL.md` — Tests Go
- `.agents/skills/golang-lint/SKILL.md` — Configuration lint Go
- `.agents/skills/golang-how-to/SKILL.md` — Orchestrateur de skills Go

Aucun ADR externe n'existe — les exigences sont complètement capturées dans SPEC.md et les décisions ci-dessus.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Aucun code existant — projet greenfield

### Established Patterns
- Architecture micro-noyau: core = infrastructure uniquement, pas de logique métier
- Communication JSON-RPC stdin/stdout: pattern HashiCorp go-plugin
- NATS embedded: bus interne éphémère, pas de port réseau exposé

### Integration Points
- Phase 2 ajoutera le Plugin SDK et le plugin MQTT
- Phase 4 ajoutera le dashboard React
- Le core est la fondation pour toutes les phases suivantes

</code_context>

<specifics>
## Specific Ideas

- L'utilisateur a Go 1.26.4 installé sur sa machine
- Les skills golang sont installés et disponibles dans `.agents/skills/`
- L'agent dev doit obligatoirement utiliser les skills Go pour suivre les bonnes pratiques
- Le core doit rester sous 5000 LOC (hors tests et code généré)
- Le projet est pour une soutenance de 3ème année avec démo live

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 01-core-foundation*
*Context gathered: 2026-06-30*

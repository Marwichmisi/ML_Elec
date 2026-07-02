# Phase 01: Core Foundation - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-06-30
**Phase:** 01-core-foundation
**Areas discussed:** Layout du projet Go, Framework REST, SQLite driver, Cycle de vie & shutdown

---

## Layout du projet Go

| Option | Description | Selected |
|--------|-------------|----------|
| cmd/ + internal/ | Structure standard Go avec binaire dans cmd/ et composants dans internal/ | ✓ |
| cmd/ + pkg/ | Plus ouvert mais risque d'exposer trop d'API interne | |
| Structure plate | Tout dans le dossier racine, chaos garanti au-delà de 1000 LOC | |

**User's choice:** cmd/ + internal/ (Recommandé)
**Notes:** Structure standard Go, isole le core, facile à maintenir

| Option | Description | Selected |
|--------|-------------|----------|
| 5 packages | internal/nats, internal/plugin, internal/api, internal/storage, internal/config | ✓ |
| 3 packages groupés | internal/core, internal/api, internal/storage — moins de packages mais couplage | |
| Tu décides | L'agent dev choisit selon les besoins | |

**User's choice:** 5 packages (Recommandé)
**Notes:** Un par composant du SPEC, clair et maintenable

| Option | Description | Selected |
|--------|-------------|----------|
| ./config.yaml | Dans le répertoire courant, puis ~/.config/ml-elec/config.yaml | ✓ |
| ./ml-elec.yaml | Nom spécifique au projet | |
| Tu décides | L'agent dev choisit | |

**User's choice:** ./config.yaml (Recommandé)
**Notes:** Simple pour le dev local

| Option | Description | Selected |
|--------|-------------|----------|
| Tests unitaires à côté | _test.go dans le même dossier que le code | ✓ |
| Dossier test/ séparé | Tous les tests dans test/ à la racine | |
| Les deux | Tests unitaires à côté + tests d'intégration dans test/ | |

**User's choice:** Tests unitaires à côté (Recommandé)
**Notes:** Pattern standard Go, les tests ont accès aux internals du package

| Option | Description | Selected |
|--------|-------------|----------|
| google/wire | DI à compile-time, génère du code, pas de反射 | ✓ |
| uber-go/dig | DI à runtime avec反射, plus flexible mais overhead | |
| uber-go/fx | Framework complet avec lifecycle hooks, trop lourd | |
| samber/do | Alternative moderne avec service containers | |

**User's choice:** google/wire (Recommandé)
**Notes:** Simple, performant, bien documenté, idéal pour un micro-noyau

| Option | Description | Selected |
|--------|-------------|----------|
| Module unique | Un seul go.mod dans la racine | |
| Monorepo avec go.work | go.work pour workspace local | |
| Tu décides | L'agent dev choisit | |

**User's choice:** Tu décides
**Notes:** L'agent dev choisit selon les besoins

| Option | Description | Selected |
|--------|-------------|----------|
| Mentionner dans CONTEXT.md | Section 'Development Constraints' dans CONTEXT.md | |
| Mettre dans AGENTS.md | Règle dans AGENTS.md pour forcer le chargement | |
| Les deux | Mentionner dans CONTEXT.md ET dans AGENTS.md | ✓ |

**User's choice:** Les deux
**Notes:** Double sécurité pour s'assurer que les skills Go sont utilisés

---

## Framework REST

| Option | Description | Selected |
|--------|-------------|----------|
| net/http standard pur | Pas de dépendance externe, plus de boilerplate mais zéro dépendance | ✓ |
| net/http + Chi router | Léger, performant, compatible stdlib | |
| Gin | Populaire mais utilise反射 | |
| Tu décides | L'agent dev choisit | |

**User's choice:** net/http standard pur
**Notes:** Zéro dépendance externe

| Option | Description | Selected |
|--------|-------------|----------|
| Handlers dans internal/api/ | Chaque endpoint a sa fonction handler | ✓ |
| Interface Handler | Définir une interface Handler avec method ServeHTTP | |
| Tu décides | L'agent dev choisit | |

**User's choice:** Handlers dans internal/api/ (Recommandé)
**Notes:** Simple, clair, maintenable

| Option | Description | Selected |
|--------|-------------|----------|
| Swagger/OpenAPI | Générer spec OpenAPI à partir des annotations Go (swaggo) | ✓ |
| README uniquement | Documenter les endpoints dans le README | |
| Pas de doc pour v1 | Juste les endpoints fonctionnels | |

**User's choice:** Swagger/OpenAPI (Recommandé)
**Notes:** Doc automatique, testable, standard industriel

| Option | Description | Selected |
|--------|-------------|----------|
| JSON standard | Toutes les réponses en JSON avec structure {"data": ..., "error": ...} | ✓ |
| JSON + content negotiation | Supporter JSON et XML via Accept header | |
| Tu décides | L'agent dev choisit | |

**User's choice:** JSON standard (Recommandé)
**Notes:** Standard, simple à consommer

---

## SQLite driver

| Option | Description | Selected |
|--------|-------------|----------|
| modernc.org/sqlite | Driver pure Go, pas de CGO, cross-compilation facile | ✓ |
| mattn/go-sqlite3 | Driver CGO, plus mature mais nécessite GCC | |
| Tu décides | L'agent dev choisit | |

**User's choice:** modernc.org/sqlite (Recommandé)
**Notes:** Pas de dépendance C, performance comparable

| Option | Description | Selected |
|--------|-------------|----------|
| Paramétrées uniquement | Toujours utiliser db.Query('SELECT ... WHERE id = ?', id) | |
| Avec builder de requêtes | Utiliser squirrel pour builder les requêtes | ✓ |
| Tu décides | L'agent dev choisit | |

**User's choice:** Avec builder de requêtes
**Notes:** Plus de sécurité avec squirrel

| Option | Description | Selected |
|--------|-------------|----------|
| Migrations embeddées | Fichiers .sql dans internal/storage/migrations/, embeddés avec Go embed | |
| Bibliothèque de migrations | Utiliser golang-migrate/migrate | ✓ |
| Schéma en dur dans le code | Créer les tables dans le code Go | |

**User's choice:** Bibliothèque de migrations
**Notes:** Plus de features avec golang-migrate/migrate

| Option | Description | Selected |
|--------|-------------|----------|
| Vérification au démarrage | PRAGMA integrity_check au démarrage | |
| Vérification périodique | Vérifier l'intégrité toutes les N minutes | |
| Les deux | Vérification au démarrage + périodique | ✓ |

**User's choice:** Les deux
**Notes:** Maximum de robustesse

---

## Cycle de vie & shutdown

| Option | Description | Selected |
|--------|-------------|----------|
| context.Context + signal | Écouter SIGINT/SIGTERM, annuler le contexte root | ✓ |
| sync.WaitGroup | Lancer chaque composant dans une goroutine | |
| Tu décides | L'agent dev choisit | |

**User's choice:** context.Context + signal (Recommandé)
**Notes:** Pattern standard Go

| Option | Description | Selected |
|--------|-------------|----------|
| LIFO | Arrêter dans l'ordre inverse du démarrage | ✓ |
| Parallèle | Arrêter tous les composants en parallèle | |
| Tu décides | L'agent dev choisit | |

**User's choice:** LIFO (Recommandé)
**Notes:** Évite les erreurs de connexion pendant l'arrêt

| Option | Description | Selected |
|--------|-------------|----------|
| 30 secondes | Suffisant pour SQLite et NATS embedded | ✓ |
| 10 secondes | Plus agressif | |
| Configurable | Timeout défini dans le fichier de config | |

**User's choice:** 30 secondes (Recommandé)
**Notes:** Suffisant pour les composants du core

| Option | Description | Selected |
|--------|-------------|----------|
| Fail-fast | Si un composant échoue, exit immédiatement | |
| Retry avec backoff | Réessayer N fois avant d'abandonner | ✓ |
| Tu décides | L'agent dev choisit | |

**User's choice:** Retry avec backoff
**Notes:** Plus résilient

---

## the agent's Discretion

- Organisation des fichiers dans chaque package (l'agent dev choisit)
- Paramètres de retry (backoff exponentiel, max retries)
- Configuration des paramètres de logging
- Structure exacte des réponses JSON d'erreur
- Choix du logger (slog standard, zerolog, zap)

## Deferred Ideas

None — discussion stayed within phase scope

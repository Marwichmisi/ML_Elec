---
project_name: 'ML_Elec'
user_name: 'Marwane'
date: '2026-03-05'
sections_completed:
  ['technology_stack', 'language_rules', 'framework_rules', 'testing_rules', 'quality_rules', 'workflow_rules', 'anti_patterns']
status: 'complete'
rule_count: 54
optimized_for_llm: true
---

# Project Context for AI Agents

_This file contains critical rules and patterns that AI agents must follow when implementing code in this project. Focus on unobvious details that agents might otherwise miss._

---

## Technology Stack & Versions

### Runtime & Langage
- **Node.js** 20 LTS — runtime edge principal (ARM64 + x86_64 natif)
- **TypeScript** 5.x — `strict: true` obligatoire sur TOUS les packages sans exception
- **pnpm** 9.x — gestionnaire de paquets (workspaces monorepo)

### Monorepo Structure
```
ml-elec/
├── apps/web/              # React 18 + Vite 6 + React Flow
├── packages/edge-agent/   # Node.js + DuckDB + Fastify
├── packages/shared/       # Types TypeScript partagés + schémas Zod
└── packages/ui/           # Composants React Flow custom (nodes ML_Elec)
```

### Backend / Edge Agent
- **Fastify** 4.x — API REST (pas Express — 3× plus rapide, TypeScript-first)
- **@duckdb/node-api** 1.4.4 — base de données analytique embarquée (API Neo)
- **node-opcua** 2.163.x — client OPC-UA complet (MIT, reconnexion built-in)

### Frontend
- **React** 18.x — UI principale
- **Vite** 6.x — build tool (requis par Vitest 4.x)
- **@xyflow/react** 12.x — éditeur visuel node-based (React Flow)

### Tests
- **Vitest** 4.x — tests unitaires + intégration (Vite-compatible, TS natif)
- **Playwright** 1.50+ — tests E2E headless

### Déploiement
- **Docker + Buildx** 26.x — images multi-arch `linux/amd64` + `linux/arm64`
- **Turborepo** — builds incrémentaux monorepo

## Critical Implementation Rules

### Language-Specific Rules (TypeScript)

#### Configuration TypeScript
- `strict: true` est OBLIGATOIRE dans tous les `tsconfig.json` — ne jamais désactiver
- `packages/edge-agent` : `target: "ES2022"`, `module: "NodeNext"`, `moduleResolution: "NodeNext"` — OBLIGATOIRE pour la compatibilité ESM Node.js 20
- `apps/web` : `module: "ESNext"`, `moduleResolution: "Bundler"` (Vite)
- Les types partagés vivent UNIQUEMENT dans `packages/shared` — jamais dupliqués entre packages

#### Imports / Exports
- Dans `packages/edge-agent`, les imports DOIVENT inclure l'extension `.js` (même pour les fichiers `.ts`) — règle ESM NodeNext stricte :
  - ✅ `import { MetricsWriter } from './metrics-writer.js'`
  - ❌ `import { MetricsWriter } from './metrics-writer'`
- Utiliser des `import type` pour les imports uniquement utilisés en typage — réduit le bundle et évite les effets de bord circulaires
- Les exports de `packages/shared` doivent passer par `index.ts` — jamais d'imports directs vers des sous-fichiers depuis un autre package

#### Gestion des Erreurs
- TOUJOURS typer les erreurs capturées : `catch (err) { if (err instanceof Error)... }` — TypeScript strict ne permet pas `catch (err: any)`
- Les erreurs OPC-UA et DuckDB doivent être enveloppées dans des types d'erreurs domaine (ex: `OpcUaConnectionError`, `StorageWriteError`) définis dans `packages/shared/errors.ts`
- Utiliser `Result<T, E>` ou `never throw` dans les modules edge — les exceptions non catchées plantent le processus Node.js sans recovery

#### Zod — Validation des Données
- Tous les inputs externes (payloads API REST, configs JSON, données OPC-UA) DOIVENT être validés via des schémas Zod définis dans `packages/shared`
- Ne jamais `as unknown as MyType` — utiliser `schema.parse()` ou `schema.safeParse()`
- Les schémas Zod servent aussi à générer les types TypeScript : `z.infer<typeof schema>`

### Framework-Specific Rules

#### DuckDB — @duckdb/node-api (API Neo)
- Utiliser EXCLUSIVEMENT l'API Neo (promesses) — l'ancienne API callback est dépréciée :
  - ✅ `const db = await Database.create('metrics.duckdb')`
  - ❌ `new duckdb.Database('metrics.duckdb', (err) => {...})`
- Pour l'ingestion haute fréquence, utiliser `Appender` (pas de SQL string parsing) :
  ```ts
  const appender = await conn.createAppender('main', 'metrics')
  appender.appendVarchar(nodeId)
  appender.appendDouble(value)
  appender.appendTimestamp(timestamp)
  appender.endRow()
  await appender.flush()
  ```
- Sur Raspberry Pi 4 (4 GB RAM), TOUJOURS configurer au démarrage :
  ```sql
  SET memory_limit = '1.5GB'; SET threads = 2;
  ```
- Utiliser `TIMESTAMPTZ` (pas `TIMESTAMP`) pour toutes les colonnes de séries temporelles
- Le partitionnement Hive est OBLIGATOIRE pour les données metrics : `data/year=YYYY/month=MM/day=DD/metrics.parquet`
- Ne JAMAIS modifier rétroactivement les fichiers Parquet existants — append-only

#### node-opcua
- Toujours créer le client avec `keepSessionAlive: true` et `connectionStrategy` avec backoff exponentiel — la reconnexion automatique est built-in mais doit être explicitement activée
- Le `NodeId` OPC-UA est sensible à la casse et au namespace : `"ns=1;s=Temperature"` ≠ `"ns=2;s=Temperature"`
- En production, utiliser `SecurityPolicy.Basic256Sha256` + `MessageSecurityMode.SignAndEncrypt` — jamais `None` hors dev local
- node-opcua génère les certificats X509 automatiquement au premier démarrage dans `./pki/` — ne pas commiter ce dossier dans git
- L'approche recommandée est **hybride** : polling `Read()` pour compatibilité PLCs anciens + subscription `CreateMonitoredItems()` pour variables critiques haute fréquence

#### React Flow (@xyflow/react v12)
- Utiliser le pattern **Computing Flows** pour la propagation des données : `updateNodeData(nodeId, data)` → `useNodeConnections(handleId)` → `useNodesData(nodeIds)`
- Les nœuds ML_Elec custom sont définis dans `packages/ui/` et enregistrés via `nodeTypes` map — ne jamais définir les types inline dans le composant parent (cause des re-renders infinis)
- Le graphe est exporté en `ReactFlowJsonObject` et interprété côté edge-agent comme un DAG — les `node.data` doivent TOUJOURS correspondre aux interfaces TypeScript définies dans `packages/shared/pipeline-nodes.ts`
- `null`/`undefined` dans `node.data.output` = arrêt du flux dans la branche (conditional branching) — convention à respecter dans tous les nœuds custom
- Limite soft : 100 nœuds max par workflow — au-delà activer `useVirtualizer` pour la performance

#### Fastify (edge-agent API)
- Tous les handlers de routes DOIVENT déclarer leurs schémas JSON (request + response) pour la sérialisation rapide et la validation automatique
- Utiliser les plugins Fastify (`fastify.register()`) pour la modularité — ne pas tout mettre dans le fichier principal
- JWT obligatoire sur toutes les routes sauf `/health` et `/metrics` (Prometheus) — utiliser `@fastify/jwt`

### Testing Rules

#### Organisation des Tests
- Les tests unitaires sont co-localisés avec le code source : `src/duckdb/metrics-writer.ts` → `src/duckdb/metrics-writer.test.ts`
- Les tests d'intégration utilisent le suffixe `.integration.test.ts`
- Les tests E2E vivent dans `apps/web/e2e/` (Playwright uniquement)
- Ne jamais mélanger tests unitaires et d'intégration dans le même fichier

#### Seuils de Couverture (non négociables)
- `packages/edge-agent` : **85% statements**, **80% branches**
- `packages/shared` : **95% statements** (types + utilitaires critiques)
- `apps/web` : **70% statements**

#### Vitest — Règles Spécifiques
- Toujours utiliser `@vitest/coverage-v8` (provider V8, pas istanbul)
- DuckDB doit être testé avec une base **in-memory** : `Database.create(':memory:')` — jamais de fichier `.duckdb` dans les tests
- Pattern `beforeEach` + `afterEach` OBLIGATOIRE pour créer/fermer la connexion DuckDB dans chaque test — les connexions non fermées causent des fuites mémoire

#### Tests d'Intégration OPC-UA
- Utiliser `OPCUAServer` de `node-opcua` comme mock server en mémoire — ne jamais dépendre d'un PLC réel en CI
- Le mock server doit exposer au minimum : `ns=1;s=Temperature`, `ns=1;s=Current`, `ns=1;s=Vibration`
- Tester OBLIGATOIREMENT le scénario de reconnexion : shutdown du mock server → attente → restart → vérifier que le collecteur se reconnecte automatiquement sans intervention

#### Tests E2E Playwright
- Scénario minimum obligatoire : créer un pipeline complet (Source OPC-UA → Filtre → Alerte)
- Tests en mode headless en CI, headed en développement local
- Lancer avec : `pnpm --filter web exec playwright test --reporter=html`
- Les snapshots visuels (screenshots) sont interdits en CI — utiliser des assertions sur le DOM uniquement

#### Mocks & Données de Test
- Les fixtures de données Parquet de test vivent dans `packages/edge-agent/src/__fixtures__/`
- Ne jamais mocker `@duckdb/node-api` — utiliser `:memory:` à la place
- Mocker `node-opcua` uniquement pour les tests unitaires purs ; les tests d'intégration utilisent le vrai `OPCUAServer`

### Code Quality & Style Rules

#### Linting & Formatage
- **ESLint** v9 avec `@typescript-eslint/parser` — configuration partagée depuis `packages/shared/eslint.config.ts`, héritée par tous les packages
- **Prettier** v3 — formatage uniforme, intégré au pre-commit hook (Husky + lint-staged) — ne jamais formater manuellement
- Les règles critiques activées : `@typescript-eslint/no-explicit-any`, `@typescript-eslint/no-floating-promises`, `no-console` (sauf dans scripts)
- `no-floating-promises` est CRITIQUE — toute promesse non awaitée dans edge-agent peut causer des corruptions silencieuses de données

#### Conventions de Nommage
- **Fichiers** : `kebab-case.ts` pour tous les fichiers source
- **Classes & Interfaces** : `PascalCase` — ex: `MetricsWriter`, `IOpcUaCollector`
- **Fonctions & Variables** : `camelCase`
- **Constantes globales** : `SCREAMING_SNAKE_CASE` — ex: `DEFAULT_POLLING_MS`
- **Types Zod** : préfixe `z` + PascalCase — ex: `zMetricPayload`
- **Fichiers de test** : `*.test.ts` (unitaires), `*.integration.test.ts`
- **Nœuds React Flow custom** : suffixe `Node` — ex: `OpcUaSourceNode`, `AnomalyDetectorNode`

#### Organisation du Code (Architecture Hexagonale)
- `packages/edge-agent/src/` suit la structure des ports & adapters :
  ```
  src/
  ├── domain/        # Logique métier pure — zéro import de lib externe
  ├── ports/         # Interfaces TypeScript (IOpcUaCollector, IStorageEngine)
  ├── adapters/      # Implémentations concrètes (node-opcua, DuckDB)
  ├── api/           # Routes Fastify (handlers uniquement, pas de logique)
  └── scheduler/     # Tâches cron (flush Parquet, compaction)
  ```
- Le dossier `domain/` ne doit JAMAIS importer depuis `adapters/` — uniquement depuis `ports/`
- Les handlers Fastify ne contiennent pas de logique métier — ils délèguent au domaine via les ports

#### Documentation
- JSDoc obligatoire sur toutes les interfaces publiques de `packages/shared`
- Les fonctions du domaine métier doivent avoir un commentaire expliquant le **pourquoi** (pas le quoi)
- Le format de commit suit **Conventional Commits** : `feat:`, `fix:`, `chore:`, `test:`, `docs:`, `refactor:`

### Development Workflow Rules

#### Commandes Monorepo (pnpm)
- Toujours exécuter les commandes depuis la racine via pnpm workspaces :
  - `pnpm -r run build` — build tous les packages
  - `pnpm -r run test` — tests dans tous les packages
  - `pnpm --filter edge-agent run dev` — dev d'un package spécifique
- Ne JAMAIS `cd packages/edge-agent && npm install` — utiliser uniquement `pnpm add <pkg> --filter edge-agent` depuis la racine
- `pnpm install --frozen-lockfile` en CI — le lockfile ne doit jamais diverger silencieusement

#### Git & Branches
- Convention de nommage des branches :
  - `feat/<ticket>-description-courte`
  - `fix/<ticket>-description-courte`
  - `chore/description`
- La branche `main` est protégée — aucun push direct, uniquement via PR
- Chaque PR doit passer le pipeline CI complet (lint + typecheck + tests) avant merge

#### Docker — Build Multi-Arch
- Toujours builder avec `docker buildx build --platform linux/amd64,linux/arm64`
- Ne JAMAIS builder une image amd64-only pour la production — le Raspberry Pi (ARM64) est la cible principale
- Tester l'image ARM64 localement avec QEMU avant de pusher : `docker run --platform linux/arm64 ghcr.io/ml-elec/edge-agent:latest`
- L'image Docker de base est `node:20-alpine` — ne pas utiliser `node:20` (trop lourde pour l'edge)

#### Variables d'Environnement
- Toutes les variables d'environnement sont définies dans `packages/edge-agent/src/config.ts` via un schéma Zod validé au démarrage
- Le processus s'arrête au démarrage si une variable requise est manquante (fail-fast) — ne jamais utiliser de valeurs par défaut silencieuses pour les configs critiques (OPCUA_ENDPOINT, DUCKDB_DATA_DIR)
- Les secrets (clés API LLM, JWT_SECRET) ne transitent JAMAIS dans le code — uniquement via variables d'environnement ou Docker secrets

#### CI/CD Pipeline
- **Workflow 1** `ci.yml` : lint + typecheck + tests unitaires → déclenché sur chaque push et PR
- **Workflow 2** `build-docker.yml` : build image multi-arch + push GHCR → déclenché uniquement sur les tags `v*`
- Le cache GitHub Actions (`cache-from: type=gha`) est OBLIGATOIRE pour les builds Docker — sans lui, chaque build prend 15+ minutes

### Critical Don't-Miss Rules

#### Contrainte Absolue — Release Blocker
- **NFR01 : Onboarding ≤ 15 minutes** est un release blocker absolu — tout choix d'implémentation qui ajoute des étapes manuelles à l'onboarding est INTERDIT en V1
- Chaque nouvelle feature doit être évaluée : "Est-ce que ça casse l'onboarding en 15 min ?" — si oui, reporter en V2

#### Anti-Patterns à Éviter Absolument
- ❌ **Ne jamais utiliser `any` en TypeScript** — utiliser `unknown` + type guard ou `z.infer<typeof schema>` via Zod
- ❌ **Ne jamais bloquer l'event loop Node.js** — toutes les opérations DuckDB et OPC-UA sont asynchrones ; ne jamais utiliser les variantes sync (`execSync`, etc.)
- ❌ **Ne jamais exposer l'API Fastify sur 0.0.0.0 sans authentification JWT** — en edge, l'API doit être accessible uniquement sur le réseau local (LAN)
- ❌ **Ne jamais écrire directement dans un fichier Parquet existant** — uniquement via DuckDB Appender ou `COPY ... TO` vers un nouveau fichier
- ❌ **Ne jamais stocker de clé API LLM en clair** dans le code, les logs, ou les variables d'environnement non chiffrées
- ❌ **Ne jamais ajouter de dépendance cloud obligatoire** pour les fonctions core (ingestion, détection, alertes, dashboards) — le mode offline-first est non négociable (NFR17)

#### Sécurité Edge — Règles Non Négociables
- L'audit log (`AuditLog` entity) est **append-only et inaltérable** — aucun endpoint DELETE/UPDATE ne doit exister sur cette table (NFR11)
- Les données Parquet contenant des mesures de production doivent être chiffrées avec `ENCRYPTION_CONFIG` DuckDB (AES-256-GCM) — pas optionnel
- Le dossier `./pki/` (certificats OPC-UA auto-générés) doit être dans `.gitignore` et `.dockerignore`
- RBAC : vérifier les permissions AVANT d'exécuter la logique métier — jamais après (fail-fast security)

#### Offline-First — Règles d'Implémentation
- Le module LLM DOIT être désactivable via config sans redémarrage du système
- Quand le LLM est désactivé ou indisponible, l'interface affiche les alertes ML sans explication textuelle — jamais de spinner infini ou d'erreur bloquante
- Toutes les features core (ingestion, détection, alertes, dashboards) doivent passer les tests fonctionnels avec `NETWORK=none` (simulation offline)

#### Performance Edge — Gotchas Critiques
- DuckDB et node-opcua tournent dans le **même processus Node.js** — une requête analytique lourde peut bloquer le polling OPC-UA ; utiliser `SET threads = 2` et planifier les requêtes lourdes hors des fenêtres de collecte critique
- Les fichiers Parquet se fragmentent avec le temps (nombreux petits fichiers) — la tâche de **compaction nocturne** (cron 02:00) est OBLIGATOIRE en production
- React Flow avec >50 nœuds sans `useVirtualizer` cause des freezes UI — implémenter la virtualisation dès le premier rendu, pas en optimisation tardive

#### Modèle de Données — Règles d'Intégrité
- Les 8 entités du modèle conceptuel (`Equipement`, `Capteur`, `Pipeline`, `Alerte`, `ModeleML`, `Utilisateur`, `AuditLog`, `Plugin`) sont définies dans `packages/shared/domain/entities.ts` — ne jamais les redéfinir dans un package
- `AuditLog` est **toujours écrit avant** la modification de configuration (write-ahead) — jamais après
- L'entité `Alerte` a TOUJOURS un niveau de sévérité parmi exactement 3 valeurs : `'info' | 'warning' | 'critical'` — pas de valeur libre

#### Internationalisation
- Toutes les chaînes de l'interface sont externalisées dans des fichiers de traduction (`fr.json`, `en.json`) dès le premier commit — jamais de texte hardcodé dans les composants React
- Les explications LLM sont générées dans la langue de l'interface utilisateur — passer le paramètre `locale` à toutes les requêtes LLM

---

## Usage Guidelines

**Pour les agents IA :**
- Lire ce fichier AVANT d'implémenter tout code dans ce projet
- Suivre TOUTES les règles exactement telles que documentées
- En cas de doute, préférer l'option la plus restrictive
- Mettre à jour ce fichier si de nouveaux patterns émergent

**Pour les humains :**
- Garder ce fichier lean et focalisé sur les besoins des agents
- Mettre à jour lors de changements de stack technologique
- Réviser trimestriellement pour supprimer les règles obsolètes
- Supprimer les règles qui deviennent évidentes avec le temps

_Dernière mise à jour : 2026-03-05_

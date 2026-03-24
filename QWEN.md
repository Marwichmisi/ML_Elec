# QWEN.md — Contexte Projet ML_Elec

**Dernière mise à jour :** 2026-03-24  
**Statut :** Specs complètes — Prêt pour implémentation  
**Utilisateur :** Marwane  
**Langue de communication :** Français

---

## 📋 Vue d'Ensemble du Projet

**ML_Elec** est une plateforme open-source de maintenance prédictive pour équipements électriques, conçue pour les PMEs industrielles (50–500 employés). 

### Proposition de Valeur Unique
- **Offline-first** : Fonctionne sans cloud, sans dépendance réseau, en souveraineté totale
- **Double couche IA** : ML prédictif (Isolation Forest / Autoencoder / RUL) + LLM explicatif en langage naturel
- **No-code accessible** : Éditeur visuel de pipelines, onboarding en ≤ 15 minutes

### Public Cible
| Persona | Rôle | Besoins |
|---------|------|---------|
| **Karim** | Technicien de maintenance | Voir en 3s, décider en 60s, sans jargon |
| **Sofiane** | Ingénieur électrotechnicien | Configurer sans code, déployer seul |
| **M. Benali** | Responsable production | Vue ROI synthétique (pannes évitées, économies) |
| **Yasmine** | Étudiante/contributrice | Documentation claire, architecture plugin ouverte |

---

## 🏗️ Architecture Technique

### Stack Technologique Validée

| Couche | Technologie | Version | Usage |
|--------|-------------|---------|-------|
| **Runtime** | Node.js | 20 LTS | Runtime edge principal (ARM64 + x86_64) |
| **Langage** | TypeScript | 5.x | `strict: true` obligatoire |
| **Package Manager** | pnpm | 9.x | Workspaces monorepo |
| **Backend Edge** | Fastify | 4.x | API REST (3× plus rapide qu'Express) |
| **Base de Données** | DuckDB | 1.4.4+ | Stockage analytique embarqué (API Neo) |
| **Protocoles** | node-opcua | 2.163.x | Client OPC-UA complet |
| | paho-mqtt | - | Ingestion MQTT |
| **Frontend** | React | 18.x | UI principale |
| | Vite | 6.x | Build tool |
| | React Flow | 12.x | Éditeur visuel no-code |
| | shadcn/ui | - | Composants UI |
| | Zustand | - | State management |
| | TanStack Query | - | Data fetching |
| **Tests** | Vitest | 4.x | Tests unitaires + intégration |
| | Playwright | 1.50+ | Tests E2E |
| **Déploiement** | Docker + Buildx | 26.x | Images multi-arch (AMD64 + ARM64) |
| **Monorepo** | Turborepo | - | Builds incrémentaux |

### Structure Monorepo

```
ml-elec/
├── apps/web/                    # React 18 + Vite 6 + React Flow
│   ├── src/
│   │   ├── components/          # shadcn/ui + composants custom
│   │   ├── stores/              # Zustand stores
│   │   ├── hooks/               # TanStack Query + custom hooks
│   │   └── router.tsx           # React Router v6
│   └── package.json
│
├── packages/edge-agent/         # Node.js + DuckDB + Fastify
│   ├── src/
│   │   ├── domain/              # Logique métier pure
│   │   ├── ports/               # Interfaces (IOpcUaCollector, IStorageEngine)
│   │   ├── adapters/            # Implémentations (node-opcua, DuckDB)
│   │   ├── api/                 # Routes Fastify + RBAC
│   │   ├── ml/                  # Isolation Forest + ONNX
│   │   ├── llm/                 # Providers API (OpenAI, Anthropic, Mistral)
│   │   └── storage/             # DuckDB + Parquet + migrations SQL
│   └── package.json
│
├── packages/shared/             # Code partagé
│   ├── domain/                  # Entités TypeScript
│   ├── schemas/                 # Schémas Zod
│   └── errors/                  # Classes d'erreurs métier
│
└── packages/ui/                 # Composants React Flow custom
    ├── components/              # NodeField, NodeHeader, StatusBadge
    └── nodes/                   # MqttSourceNode, AnomalyDetectorNode, etc.
```

---

## 📊 État Actuel du Projet

### ✅ Phase de Spécifications — TERMINÉE

| Document | Statut | Chemin |
|----------|--------|--------|
| **Product Brief** | ✅ Complet | `_bmad-output/planning-artifacts/product-brief-ML_Elec-2026-03-04.md` |
| **PRD** | ✅ Complet (670 lignes) | `_bmad-output/planning-artifacts/prd.md` |
| **Architecture** | ✅ Complète (2674 lignes) | `_bmad-output/planning-artifacts/architecture.md` |
| **UX Design** | ✅ Complet (1480 lignes) | `_bmad-output/planning-artifacts/ux-design-specification.md` |
| **Épics & Stories** | ✅ Complet (978 lignes) | `_bmad-output/planning-artifacts/epics.md` |
| **Project Context** | ✅ Complet (54 règles) | `_bmad-output/project-context.md` |

### 📈 Couverture des Exigences

| Type | Total | Couvert | Pourcentage |
|------|-------|---------|-------------|
| **Functional Requirements (FRs)** | 52 | 52 | 100% |
| **Non-Functional Requirements (NFRs)** | 29 | 29 | 100% |
| **UX Design Requirements (UX-DRs)** | 30 | 30 | 100% |

### 📦 Structure des 8 Épics (35 Stories)

| # | Epic | Stories | FRs | Valeur Utilisateur |
|---|------|---------|-----|-------------------|
| **1** | Edge Agent - Ingestion & Stockage | 5 | 6 | Connecter équipements + acquérir données temps-réel |
| **2** | Détection & Intelligence Augmentée | 8 | 13 | Alertes prédictives + explications LLM + suggestions |
| **3** | Éditeur No-Code - Pipelines Visuels | 4 | 5 | Configurer pipelines par drag & drop sans code |
| **4** | Dashboard - Visualisation & Décision | 6 | 8 | Voir en 3s, décider en 60s |
| **5** | Notifications & Partage | 2 | 3 | Recevoir alertes + partager résumés |
| **6** | Administration & Sécurité | 4 | 9 | RBAC + Audit log + Mises à jour |
| **7** | Onboarding & Déploiement | 3 | 2 | Déployer + premier signal ≤ 15 min |
| **8** | Internationalisation | 1 | 1 | Interface FR/EN + architecture i18n |

---

## 🎯 Prochaines Étapes Recommandées

### Option 1 : Commencer l'Implémentation
Commencer par **Epic 1 Story 1** : Setup Monorepo & Infrastructure de Build

```bash
# Commands de départ (à documenter dans le README)
mkdir -p ml-elec/{apps/web,packages/{edge-agent,shared,ui}}
cd ml-elec
pnpm init
pnpm add -D turbo typescript eslint prettier
```

### Option 2 : Vérifier la Prêtude à l'Implémentation
Utiliser la compétence BMad : `bmad-check-implementation-readiness`

### Option 3 : Planifier les Sprints
Utiliser la compétence BMad : `bmad-sprint-planning`

---

## 📁 Structure des Fichiers Clés

### Répertoires de Planification (BMAD)

```
_bmad-output/
├── planning-artifacts/
│   ├── prd.md                          # Product Requirements Document
│   ├── architecture.md                 # Architecture Decision Document
│   ├── ux-design-specification.md      # UX Design Specification
│   ├── epics.md                        # Epic & Story Breakdown (ACTIF)
│   ├── product-brief-ML_Elec-*.md      # Product Brief
│   └── research/                       # Rapports de recherche
│       ├── market-supervision-industrielle-hybride-research-*.md
│       ├── domain-genie-electrique-electrotechnique-research-*.md
│       └── technical-faisabilite-ml-elec-research-*.md
│
├── implementation-artifacts/
│   └── sprint-status.yaml              # Statut des sprints
│
├── brainstorming/                      # Sessions de brainstorming
├── dialog/                             # Historique des conversations
├── project-context.md                  # Règles techniques pour AI agents
└── wds-workflow-status.yaml            # Statut workflow WDS
```

### BMAD Core (Système d'Agents)

```
_bmad/
├── _config/
│   └── agent-manifest.csv              # Liste des agents BMad disponibles
├── _memory/                            # Mémoire et connaissances
├── bmm/                                # Modules BMad
│   └── agents/                         # Définitions des agents
│       ├── analyst.md                  # Mary - Business Analyst
│       ├── architect.md                # Winston - Architecte
│       ├── dev.md                      # Amelia - Developer
│       ├── pm.md                       # John - Product Manager
│       ├── qa.md                       # Quinn - QA Engineer
│       ├── sm.md                       # Bob - Scrum Master
│       ├── tech-writer.md              # Paige - Technical Writer
│       └── ux-designer.md              # Sally - UX Designer
├── core/                               # Core BMad
└── wds/                                # WDS Workflow System
```

---

## 🔧 Commandes de Développement (À Implémenter)

> **NOTE :** Le code d'implémentation a été supprimé volontairement. Ces commandes sont des placeholders pour la future implémentation.

### Installation

```bash
cd ml-elec
pnpm install
```

### Build

```bash
pnpm -r run build
# ou avec Turborepo
pnpm run build
```

### Tests

```bash
# Tests unitaires + intégration
pnpm -r run test

# Tests E2E
pnpm -r run test:e2e
```

### Lint & Format

```bash
pnpm -r run lint
pnpm -r run format
```

### Déploiement

```bash
# Build Docker multi-arch
docker buildx build --platform linux/amd64,linux/arm64 -t ml-elec/edge-agent:latest .

# Docker Compose (dev local)
docker compose up -d
```

---

## 📐 Règles Techniques Critiques (Project Context)

### TypeScript — Configuration Stricte

- `strict: true` est **OBLIGATOIRE** dans tous les `tsconfig.json`
- `packages/edge-agent` : `target: "ES2022"`, `module: "NodeNext"`, `moduleResolution: "NodeNext"`
- `apps/web` : `module: "ESNext"`, `moduleResolution: "Bundler"` (Vite)

### Imports ESM — Règles Strictes

```typescript
// ✅ Bon (packages/edge-agent)
import { MetricsWriter } from './metrics-writer.js'

// ❌ Interdit
import { MetricsWriter } from './metrics-writer'
```

### DuckDB — API Neo (Promesses)

```typescript
// ✅ Bon — API Neo
const db = await Database.create('metrics.duckdb')
const conn = await db.connect()
await conn.query('SELECT * FROM metrics')

// ❌ Interdit — Ancienne API callback
const db = new duckdb.Database('metrics.duckdb', (err) => {...})
```

### Ingestion Haute Fréquence — Appender

```typescript
// ✅ Bon — Appender (pas de SQL string parsing)
const appender = await conn.createAppender('main', 'metrics')
appender.appendVarchar(nodeId)
appender.appendDouble(value)
appender.appendTimestamp(timestamp)
appender.endRow()
await appender.flush()
```

### Configuration Raspberry Pi 4 (4 GB RAM)

```sql
SET memory_limit = '1.5GB';
SET threads = 2;
```

### Schémas Zod — Validation Obligatoire

```typescript
// ✅ Bon — Validation Zod
const result = zEquipementSchema.safeParse(data)
if (!result.success) {
  throw new ValidationError(result.error)
}

// ❌ Interdit — Type assertion brute
const equipement = data as unknown as Equipement
```

---

## 🎨 Design System & UX

### Design Tokens (Basés sur le Logo)

```typescript
const tokens = {
  brand: {
    primary: '#F28A00',      // Orange énergie
    primaryHover: '#D97706',
    ink: '#2F3B46',          // Anthracite technique
    steel: '#5F6B77',        // Gris acier
  },
  surface: {
    base: '#F4F6F8',
    card: '#FFFFFF',
  },
  semantic: {
    success: '#2E7D32',
    warning: '#B7791F',
    error: '#C0392B',
    info: '#1F6FEB',
  }
}
```

### Typographie

- **UI** : IBM Plex Sans
- **Technique** : IBM Plex Mono (logs, IDs, valeurs capteurs)

### Principes UX Clés

1. **Voir en 3s, décider en 60s** — Glanceability absolue
2. **Offline invisible** — Fonctions critiques toujours disponibles
3. **Progressive disclosure par rôle** — Karim (simple) vs Sofiane (profond)
4. **Langage technicien, pas data scientist** — "Vibration anormale roulement gauche" vs "anomalie canal vibratoire axial Z, score 0.87"

---

## 🧪 Agents BMad Disponibles

| Agent | Nom | Rôle | Icone |
|-------|-----|------|-------|
| `bmad-master` | BMad Master | Orchestrateur de workflow | 🧙 |
| `analyst` | Mary | Business Analyst | 📊 |
| `architect` | Winston | Architecte Système | 🏗️ |
| `dev` | Amelia | Developer | 💻 |
| `pm` | John | Product Manager | 📋 |
| `qa` | Quinn | QA Engineer | 🧪 |
| `sm` | Bob | Scrum Master | 🏃 |
| `tech-writer` | Paige | Technical Writer | 📚 |
| `ux-designer` | Sally | UX Designer | 🎨 |

### Compétences BMad Utilisables

- `bmad-create-prd` — Créer un PRD
- `bmad-validate-prd` — Valider un PRD
- `bmad-edit-prd` — Éditer un PRD
- `bmad-create-epics-and-stories` — Créer Épics & Stories ✅ (DÉJÀ FAIT)
- `bmad-create-architecture` — Créer l'architecture
- `bmad-create-ux-design` — Créer le design UX
- `bmad-check-implementation-readiness` — Vérifier prêtitude implémentation
- `bmad-sprint-planning` — Planifier les sprints
- `bmad-sprint-status` — Vérifier statut sprint
- `bmad-party-mode` — Discussion multi-agents
- `bmad-help` — Obtenir de l'aide sur la prochaine étape

---

## 📊 Métriques de Succès (NFR Critiques)

| NFR | Description | Cible | Vérification |
|-----|-------------|-------|--------------|
| **NFR01** | Onboarding complet | ≤ 15 min | Test utilisateur (5 techniciens) |
| **NFR02** | Dashboard loading | < 2s | Test de charge |
| **NFR03** | Latence détection → notification | < 5s | Injection anomalie simulée |
| **NFR05** | Inférence ML | < 500ms | Benchmark Raspberry Pi 5 |
| **NFR06** | Taux de faux positifs | < 5% | Datasets de référence (CWRU, NASA, FEMTO) |
| **NFR28** | Machine prioritaire identifiable | ≤ 3s (≥ 90%) | Instrumentation UX |
| **NFR29** | Action utilisateur enregistrée | ≤ 60s (≥ 80%) | Instrumentation UX |

---

## 🔒 Sécurité & Conformité

### NFRs Sécurité

- **NFR07** : Données chiffrées AES-256 au repos
- **NFR08** : Communications TLS 1.2+
- **NFR09** : RBAC sur 100% des endpoints API
- **NFR10** : Conformité RGPD (minimisation données)
- **NFR11** : Audit log inaltérable (append-only, 12 mois)

### RBAC — 4 Rôles

| Rôle | Permissions | Vue par défaut |
|------|-------------|----------------|
| **Viewer** (Karim) | Lecture seule | Dashboard santé machines |
| **Manager** (M. Benali) | Lecture seule | Vue synthétique ROI |
| **Engineer** (Sofiane) | Configuration complète | Éditeur + dashboard |
| **Admin** | Gestion utilisateurs, système | Console administration |

---

## 📝 Notes Importantes

### Code Supprimé Volontairement

Le code d'implémentation initial (`apps/web/*`, `packages/*`, `docker-compose.yml`, etc.) a été supprimé volontairement le **2026-03-24** pour repartir sur une base saine avec les spécifications complètes.

**Commit :** `4dc9e4a` — _"refactor: Nettoyage code + Finalisation specs Épics & Stories"_

### Fichiers Conservés

Seuls les fichiers de **planification** et de **spécifications** sont conservés :
- `_bmad-output/planning-artifacts/*` (PRD, Architecture, UX, Épics)
- `_bmad-output/project-context.md` (Règles techniques)
- `_bmad/*` (Système d'agents BMad)

---

## 🚀 Démarrage Rapide (Pour Future Implémentation)

### 1. Initialiser Monorepo

```bash
mkdir -p ml-elec/{apps/web,packages/{edge-agent,shared,ui}}
cd ml-elec
pnpm init
pnpm add -D turbo typescript eslint prettier
```

### 2. Créer Structure de Base

Suivre le guide dans `_bmad-output/planning-artifacts/epics.md` → **Epic 1 Story 1**

### 3. Valider l'Initialisation

```bash
# Test structure
test -d apps/web && test -d packages/edge-agent && echo "✅ Structure OK"

# Test TypeScript strict
grep -r '"strict": true' packages/*/tsconfig.json apps/*/tsconfig.json && echo "✅ TypeScript strict OK"

# Test build
pnpm -r run build && echo "✅ Build OK"
```

---

## 📞 Besoin d'Aide ?

Utilise la compétence **`bmad-help`** pour obtenir des recommandations sur la prochaine étape :

```
Invoque la compétence `bmad-help` avec :
"Par où commencer l'implémentation de ML_Elec ?"
```

Ou consulte les documents de spécifications dans `_bmad-output/planning-artifacts/` pour plus de détails sur chaque Epic et Story.

---

**Fin du fichier QWEN.md**

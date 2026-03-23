# 🚀 ML_Elec - Plateforme de Maintenance Prédictive

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Node.js](https://img.shields.io/badge/Node.js-20+-green.svg)](https://nodejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)](https://www.typescriptlang.org/)

**Plateforme de maintenance prédictive pour l'électrotechnique** - Détection d'anomalies, prédiction RUL, et dashboards industriels en temps réel.

---

## 📋 Table des Matières

- [Fonctionnalités](#-fonctionnalités)
- [Architecture](#-architecture)
- [Démarrage Rapide](#-démarrage-rapide)
- [Documentation](#-documentation)
- [Développement](#-développement)
- [Déploiement](#-déploiement)
- [Contributing](#-contributing)

---

## ✨ Fonctionnalités

### Core
- 🔌 **Acquisition de données** - OPC-UA, MQTT, CSV
- 🧠 **Détection d'anomalies** - Algorithmes ML non supervisés
- 📈 **Prédiction RUL** - Remaining Useful Life des équipements
- 🎨 **Éditeur No-Code** - Pipelines de traitement visuels avec React Flow
- 📊 **Dashboards** - Visualisation temps réel et historique
- 🚨 **Alertes** - Notifications email, SMS, webhooks

### Edge-First
- 🍇 **DuckDB** - Base de données analytique embarquée
- 🔐 **Offline-First** - Fonctionne sans connexion cloud
- 📦 **Docker Multi-Arch** - Linux amd64 + ARM64 (Raspberry Pi)

---

## 🏗️ Architecture

```
ml-elec/
├── apps/
│   └── web/              # React 18 + Vite 6 + React Flow (@xyflow/react)
├── packages/
│   ├── edge-agent/       # Node.js 20 + Fastify 4 + DuckDB + node-opcua
│   ├── shared/           # Types TypeScript + schémas Zod
│   └── ui/               # Composants React Flow custom
└── docker-compose.yml    # Orchestration complète
```

### Stack Technique

| Couche | Technologie | Version |
|--------|-------------|---------|
| **Runtime** | Node.js | 20 LTS |
| **Langage** | TypeScript | 5.x (strict: true) |
| **Monorepo** | pnpm + Turborepo | 9.x + 2.x |
| **Backend** | Fastify | 4.x |
| **Database** | DuckDB | 1.5.x (@duckdb/node-api) |
| **OPC-UA** | node-opcua | 2.163.x |
| **MQTT** | mqtt.js | 5.x |
| **Frontend** | React + Vite | 18.x + 6.x |
| **React Flow** | @xyflow/react | 12.x |
| **Validation** | Zod | 3.x |
| **Tests** | Vitest + Playwright | 4.x + 1.50.x |

---

## 🚀 Démarrage Rapide

### Prérequis

- Node.js 20+ : `nvm use 20`
- pnpm 9+ : `corepack enable`
- Docker (optionnel) : `docker --version`

### Installation Locale (15 minutes)

```bash
# 1. Cloner le repository
git clone https://github.com/ml-elec/ml-elec.git
cd ml-elec

# 2. Installer les dépendances
pnpm install

# 3. Configurer les variables d'environnement
cp .env.example .env
# Éditer .env avec vos configurations

# 4. Démarrer les services (Docker)
docker compose up -d

# 5. Démarrer le développement
pnpm run dev
```

### Accès aux Services

| Service | URL | Description |
|---------|-----|-------------|
| Web Dashboard | http://localhost:3000 | Interface utilisateur |
| Edge Agent API | http://localhost:3001 | API REST |
| Health Check | http://localhost:3001/health | État du système |
| MQTT Broker | mqtt://localhost:1883 | Broker de messages |

---

## 📚 Documentation

### Guides

- [Architecture Technique](./vision_architecture.md)
- [Configuration](./docs/configuration.md)
- [API Reference](./docs/api.md)
- [Développement Plugins](./docs/plugins.md)

### Contexte Projet

- [Project Context](_bmad-output/project-context.md) - Règles et patterns pour AI agents
- [Code Review Report](_bmad-output/code-review-report.md) - Audit complet

---

## 🛠️ Développement

### Commandes Monorepo

```bash
# Build tous les packages
pnpm run build

# Tests tous les packages
pnpm run test

# Linting
pnpm run lint

# Type checking
pnpm run typecheck

# Développement (hot reload)
pnpm run dev
```

### Commandes par Package

```bash
# Edge Agent
pnpm --filter edge-agent run dev
pnpm --filter edge-agent run test
pnpm --filter edge-agent run build

# Web Dashboard
pnpm --filter web run dev
pnpm --filter web run test:e2e

# Shared (types)
pnpm --filter shared run build
```

### Tests

```bash
# Tests unitaires (Vitest)
pnpm run test:coverage

# Tests E2E (Playwright)
pnpm --filter web exec playwright test

# Tests E2E avec UI
pnpm --filter web exec playwright test --ui
```

---

## 📦 Déploiement

### Docker Multi-Arch

```bash
# Build pour amd64 + arm64
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/ml-elec/edge-agent:latest \
  --push \
  ./packages/edge-agent
```

### Docker Compose (Production)

```bash
# Démarrer tous les services
docker compose -f docker-compose.prod.yml up -d

# Voir les logs
docker compose logs -f edge-agent

# Arrêter
docker compose down
```

### Variables d'Environnement (Production)

```bash
# .env.production
NODE_ENV=production
PORT=3001
JWT_SECRET=<secret-généré-avec-openssl-rand-base64-32>
DUCKDB_DATA_DIR=/data/duckdb
MQTT_BROKER_URL=mqtt://mqtt-broker:1883
OPCUA_ENDPOINT=opc.tcp://<votre-plc>:4840
LOG_LEVEL=warn
```

---

## 🤝 Contributing

### Workflow

1. Fork le repository
2. Créer une branche feature : `git checkout -b feat/ma-fonctionnalite`
3. Commit selon Conventional Commits : `git commit -m 'feat: ajout fonctionnalité'`
4. Push : `git push origin feat/ma-fonctionnalite`
5. Ouvrir une Pull Request

### Conventional Commits

- `feat:` Nouvelle fonctionnalité
- `fix:` Correction de bug
- `docs:` Documentation
- `test:` Tests
- `refactor:` Refactoring
- `chore:` Tâches diverses

---

## 📄 License

MIT License - voir [LICENSE](LICENSE) pour plus de détails.

---

## 👥 Équipe

- **Marwane** - Lead Developer
- Basé sur l'architecture BMad (Business Model driven development)

---

## 🙏 Remerciements

- [Fastify](https://fastify.dev/) - Framework web performant
- [DuckDB](https://duckdb.org/) - Base de données analytique
- [node-opcua](https://node-opcua.github.io/) - Stack OPC-UA
- [React Flow](https://reactflow.dev/) - Éditeur de nœuds
- [Zod](https://zod.dev/) - Validation TypeScript-first

---

**🎯 North Star Metric :** Onboarding ≤ 15 minutes

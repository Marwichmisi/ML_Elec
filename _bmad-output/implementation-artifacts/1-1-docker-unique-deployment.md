# Story 1.1 : Commande Docker Unique de Déploiement

**Status:** review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

---

## Story

As a **ingénieur électrotechnique**,  
I want **déployer ML_Elec via une commande Docker unique**,  
So that **je peux commencer l'onboarding sans configuration complexe**.

**Story ID:** 1.1  
**Story Key:** 1-1-docker-unique-deployment  
**Epic:** 1 - Onboarding & Déploiement Edge  
**Priority:** V1 Critique (MVP)  

---

## Acceptance Criteria

### Critères Principaux (BDD)

**Given** une machine avec Docker pré-installé (Ubuntu 22.04+, Debian 12, ou Raspberry Pi OS)  
**When** j'exécute `docker compose up -d` depuis le dossier ML_Elec  
**Then** tous les services ML_Elec démarrent (ingestion, ML, API, dashboard)  
**And** le dashboard est accessible sur http://localhost:3000  
**And** aucun code n'est requis pour cette opération  

**Given** un environnement avec Docker et Docker Compose installés  
**When** l'utilisateur télécharge le fichier `docker-compose.yml` et exécute la commande  
**Then** le déploiement complet prend ≤ 5 minutes (téléchargement image inclus)  
**And** l'utilisateur voit un message de succès clair avec l'URL d'accès  

**Given** les services ML_Elec en cours d'exécution  
**When** l'utilisateur exécute `docker compose down`  
**Then** tous les services s'arrêtent proprement  
**And** les données persistent (volumes Docker)  

### Critères de Performance (NFR01)

- **Onboarding total ≤ 15 minutes** : Cette story doit permettre le déploiement en ≤ 5 minutes, laissant 10 minutes pour les stories 1.2 et 1.3
- **Taille image Docker ≤ 500 MB** (cf. PRD - Succès Technique)
- **Support multi-arch** : `linux/amd64` (Intel NUC) et `linux/arm64` (Raspberry Pi 4/5)

### Critères de Qualité

- [ ] Authentification JWT activée sur toutes les routes sauf `/health` et `/metrics`
- [ ] Variables d'environnement validées au démarrage (fail-fast)
- [ ] Logs structurés accessibles via `docker compose logs -f`
- [ ] Health check endpoint fonctionnel sur `/health`
- [ ] Documentation README claire avec commande unique

---

## Tasks / Subtasks

### Task 1 : Configuration Docker Compose (AC: 1, 2, 3)
- [ ] 1.1 Créer `docker-compose.yml` avec 4 services :
  - `edge-agent` : Node.js + DuckDB + node-opcua (port 3001)
  - `web` : React 18 + Vite (port 3000)
  - `mqtt-broker` : Mosquitto (port 1883)
  - `db` : DuckDB persisté (volume)
- [ ] 1.2 Configurer les volumes Docker pour persistance des données
- [ ] 1.3 Configurer les réseaux Docker isolés
- [ ] 1.4 Définir les variables d'environnement requises
- [ ] 1.5 Ajouter health checks pour chaque service

### Task 2 : Image Docker Edge-Agent (AC: 1, 2)
- [ ] 2.1 Créer `packages/edge-agent/Dockerfile` multi-arch
- [ ] 2.2 Base : `node:20-alpine` (léger pour edge)
- [ ] 2.3 Installer les dépendances système pour node-opcua
- [ ] 2.4 Configurer pour ARM64 + AMD64 via `docker buildx`
- [ ] 2.5 Optimiser la taille (multi-stage build)

### Task 3 : Image Docker Web (AC: 1, 2)
- [ ] 3.1 Créer `apps/web/Dockerfile` multi-arch
- [ ] 3.2 Build React avec Vite en mode production
- [ ] 3.3 Servir avec nginx:alpine ou node static
- [ ] 3.4 Configurer le reverse proxy vers edge-agent

### Task 4 : Documentation & Scripts (AC: 2, 3)
- [ ] 4.1 Créer `README.md` avec commande unique
- [ ] 4.2 Créer script `deploy.sh` (optionnel, pour simplifier)
- [ ] 4.3 Documenter les variables d'environnement
- [ ] 4.4 Ajouter troubleshooting section

### Task 5 : Tests & Validation (AC: Tous)
- [ ] 5.1 Tester déploiement sur Ubuntu 22.04 (AMD64)
- [ ] 5.2 Tester déploiement sur Raspberry Pi 4/5 (ARM64)
- [ ] 5.3 Valider temps de déploiement < 5 minutes
- [ ] 5.4 Valider accès dashboard http://localhost:3000
- [ ] 5.5 Valider persistance des données après restart

---

## Dev Notes

### Architecture Patterns et Contraintes

**Monorepo Structure (Source: project-context.md#technology-stack)**
```
ml-elec/
├── apps/web/              # React 18 + Vite 6 + React Flow
├── packages/edge-agent/   # Node.js + DuckDB + Fastify
├── packages/shared/       # Types TypeScript partagés + schémas Zod
└── packages/ui/           # Composants React Flow custom
```

**Services Docker Requis :**
1. **edge-agent** : API REST Fastify (port 3001)
   - Ingestion MQTT + OPC-UA
   - Détection ML (Isolation Forest/Autoencoder)
   - Stockage DuckDB
   - Health check : `/health`

2. **web** : Dashboard React (port 3000)
   - Interface utilisateur
   - Health check : `/`

3. **mqtt-broker** : Mosquitto (port 1883)
   - Broker MQTT pour ingestion données capteurs
   - Optionnel : peut être externe

4. **db** : Volume DuckDB
   - Persistance des données metrics
   - Fichiers Parquet partitionnés

**Contraintes Techniques (Source: project-context.md)**
- Node.js 20 LTS obligatoire
- TypeScript `strict: true` partout
- DuckDB : utiliser API Neo (promesses), pas callback
- node-opcua : reconnexion automatique activée
- Fastify : JWT obligatoire sur routes protégées
- Tailles images : ≤ 500 MB total

### Source Tree Components to Touch

**Fichiers à créer :**
```
ml-elec/
├── docker-compose.yml                 # Nouveau : configuration principale
├── .env.example                       # Nouveau : variables d'environnement template
├── README.md                          # Nouveau ou MAJ : instructions déploiement
├── apps/web/
│   └── Dockerfile                     # Nouveau : build web
├── packages/edge-agent/
│   └── Dockerfile                     # Nouveau : build edge-agent
└── scripts/
    └── deploy.sh                      # Optionnel : script simplifié
```

**Fichiers à modifier (si existants) :**
```
ml-elec/
├── .gitignore                         # MAJ : ajouter .docker/ pki/
└── package.json                       # MAJ : ajouter scripts docker
```

### Testing Standards Summary

**Tests Unitaires (Vitest) :**
- `packages/edge-agent` : 85% statements, 80% branches
- `packages/shared` : 95% statements
- DuckDB testé avec `:memory:`
- Mock node-opcua pour tests unitaires

**Tests d'Intégration :**
- Utiliser `OPCUAServer` mock en mémoire
- Tester reconnexion automatique
- Valider persistance DuckDB

**Tests E2E (Playwright) :**
- Scénario : déploiement complet → accès dashboard
- Headless en CI, headed en dev local
- Assertion sur DOM, pas de snapshots

**Validation Manuelle :**
- Tester sur Ubuntu 22.04 (AMD64)
- Tester sur Raspberry Pi 4/5 (ARM64)
- Chronométrer déploiement (< 5 min)
- Valider health checks

---

## Project Structure Notes

### Alignment with Unified Project Structure

**Docker Compose :**
- `docker-compose.yml` à la racine du monorepo
- Services nommés : `edge-agent`, `web`, `mqtt-broker`, `db`
- Réseaux : `ml-elec-network` (isolé)
- Volumes : `ml-elec-data` (persistant)

**Variables d'Environnement :**
- Fichier `.env.example` à la racine
- Variables validées par Zod dans `packages/shared/config.ts`
- Secrets via Docker secrets ou variables env

**Build Multi-Arch :**
```bash
docker buildx build --platform linux/amd64,linux/arm64 -t ml-elec/edge-agent:latest packages/edge-agent
docker buildx build --platform linux/amd64,linux/arm64 -t ml-elec/web:latest apps/web
```

### Detected Conflicts or Variances

**Aucun conflit détecté** - C'est la première story d'implémentation.

**Décisions d'architecture :**
1. **Pourquoi node:20-alpine ?** → Taille réduite (critical pour edge), compatibilité ARM64 native
2. **Pourquoi Docker Compose ?** → Standard industriel, simple pour onboarding 15 min
3. **Pourquoi 4 services séparés ?** → Modularité, maintenance facile, scaling futur

---

## References

### Technical Details Sources

- **PRD (prd.md)** :
  - Executive Summary : "Onboarding ≤ 15 minutes = release blocker absolu"
  - Succès Technique : "Taille image Docker edge ≤ 500 MB"
  - FR38 : "Un nouvel utilisateur peut déployer ML_Elec via une commande Docker unique"
  - NFR01 : "Onboarding complet ≤ 15 minutes"
  - NFR15 : "Reprise automatique après redémarrage matériel < 60 secondes"
  - NFR17 : "Fonctionnement offline continu illimité"

- **Epics (epics.md)** :
  - Epic 1 : "Onboarding & Déploiement Edge"
  - Story 1.1 : "Commande Docker Unique de Déploiement" (acceptance criteria complets)

- **UX Design (ux-design-specification.md)** :
  - Platform Strategy : "PWA responsive (web) — Karim sur tablette, Sofiane sur desktop"
  - Offline : "100% des fonctions core via Service Worker"
  - Résolution cible : "≥ 768px (tablette minimum)"

- **Project Context (project-context.md)** :
  - Technology Stack : Node.js 20, TypeScript 5, pnpm 9
  - Monorepo Structure : apps/web, packages/edge-agent, packages/shared
  - Docker : Buildx 26.x, images multi-arch
  - Anti-Patterns : "Ne jamais ajouter de dépendance cloud obligatoire"

- **Architecture (architecture.md)** :
  - Fichier presque vide - à enrichir dans stories futures

---

## Dev Agent Record

### Agent Model Used

Qwen Code (bmad-create-story skill)

### Debug Log References

- Configuration chargée depuis `_bmad/bmm/config.yaml`
- Epics analysés : Epic 1, Story 1.1 extraite
- PRD analysé : NFR01, FR38, critères de succès
- UX analysé : Platform strategy, offline-first
- Project context analysé : Stack technique, règles Docker

### Completion Notes List

- Story créée avec tous les critères d'acceptation BDD
- Tasks détaillées avec subtasks actionnables
- References complètes vers documents sources
- Architecture patterns documentés
- Testing standards inclus

### File List

**Fichier créé :**
- `_bmad-output/implementation-artifacts/1-1-docker-unique-deployment.md`

**Fichiers à créer par dev-story :**
- `docker-compose.yml`
- `apps/web/Dockerfile`
- `packages/edge-agent/Dockerfile`
- `.env.example`
- `README.md` (ou MAJ)
- `scripts/deploy.sh` (optionnel)

---

## Definition of Done

- [ ] Tous les critères d'acceptation validés
- [ ] Tests unitaires écrits et passants (seuils de couverture atteints)
- [ ] Tests E2E écrits et passants
- [ ] Déploiement testé sur AMD64 (Ubuntu 22.04)
- [ ] Déploiement testé sur ARM64 (Raspberry Pi 4/5)
- [ ] Temps de déploiement ≤ 5 minutes validé
- [ ] Documentation README claire et complète
- [ ] Health checks fonctionnels
- [ ] Variables d'environnement validées (fail-fast)
- [ ] JWT activé sur routes protégées
- [ ] Review code passée (lint + typecheck)
- [ ] Story marquée comme "done" dans sprint-status.yaml

---

**Prochaine story :** 1-2-wizard-onboarding

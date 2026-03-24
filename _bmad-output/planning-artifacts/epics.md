---
stepsCompleted: ["step-01-validate-prerequisites", "step-02-design-epics", "step-03-create-stories", "step-04-final-validation"]
inputDocuments:
  - "_bmad-output/planning-artifacts/prd.md"
  - "_bmad-output/planning-artifacts/architecture.md"
  - "_bmad-output/planning-artifacts/ux-design-specification.md"
workflowType: 'epics-and-stories'
project_name: 'ML_Elec'
user_name: 'Marwane'
date: '2026-03-24'
epicStructureVersion: '2.0-fusion-llm'
validationStatus: 'complete'
validationDate: '2026-03-24'
coverageSummary:
  FRs: '52/52 (100%)'
  NFRs: '29/29 (100%)'
  UXDRs: '30/30 (100%)'
  totalEpics: 8
  totalStories: 35
---

# ML_Elec - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for ML_Elec, decomposing the requirements from the PRD, UX Design, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

**1. Acquisition & Ingestion de Données**

- **FR01** : L'ingénieur peut connecter une source de données MQTT à un pipeline de traitement
- **FR02** : L'ingénieur peut connecter une source de données OPC-UA à un pipeline de traitement
- **FR03** : Le système peut détecter et signaler dans l'interface les signaux de données manquants, hors plage ou corrompus
- **FR04** : L'ingénieur peut configurer les paramètres de connexion (broker, topic, fréquence) sans modifier de code
- **FR44** : L'ingénieur peut importer un fichier CSV de données historiques pour valider le modèle ML avant connexion au matériel réel
- **FR05** : Le système peut ingérer des données en continu (24h/7j) avec reprise automatique après interruption, sans perte de données

**2. Détection d'Anomalies & Prédiction**

- **FR06** : Le système peut détecter des anomalies sur les signaux d'un équipement via le modèle généraliste par défaut
- **FR46** : Le système est livré avec un modèle ML généraliste pré-entraîné (bundled dans l'image Docker), fonctionnel dès l'installation sans configuration ni entraînement préalable
- **FR07** : Le système peut calculer une durée de vie résiduelle (RUL) et un horizon de défaillance estimé
- **FR08** : Le système peut générer une alerte avec niveau de sévérité (info / warning / critique)
- **FR09** : L'ingénieur peut configurer des seuils d'alerte et des règles de corrélation pour réduire les faux positifs
- **FR10** : Le technicien peut soumettre un feedback sur une alerte en 1 tap (utile / fausse alerte), avec note de contexte optionnelle et confirmation immédiate
- **FR11** : Le système peut afficher un indicateur de confiance du modèle ML sur chaque alerte
- **FR50** : Le système peut signaler un contexte ambigu avant action (température ambiante élevée, capteur récemment ajouté, modèle en phase d'apprentissage)
- **FR45** : Le système peut générer une suggestion d'action concrète pour chaque alerte (commander la pièce, planifier l'arrêt, surveiller)

**3. Explication & Intelligence (LLM)**

- **FR12** : Le système peut générer une explication en langage naturel pour chaque anomalie détectée, compréhensible par un technicien BTS sans formation ML
- **FR13** : Le technicien peut poser une question en langage naturel sur l'état d'un équipement et recevoir une réponse en < 10s
- **FR14** : L'administrateur peut activer ou désactiver le module LLM
- **FR15** : L'administrateur peut configurer le fournisseur LLM et la clé API

**4. Éditeur No-Code & Pipelines**

- **FR16** : L'ingénieur peut créer un pipeline de traitement par drag & drop sans code
- **FR17** : L'ingénieur peut modifier un pipeline existant sans redéploiement
- **FR18** : L'ingénieur peut installer un plugin depuis la bibliothèque en un clic
- **FR19** : Le développeur peut créer et publier un plugin selon l'architecture documentée
- **FR47** : Le système supporte 6 types de plugins avec contrats d'interface documentés : (1) connecteurs protocole, (2) transformateurs de données, (3) modèles ML, (4) blocs de visualisation, (5) canaux de notification, (6) exportateurs
- **FR20** : L'ingénieur peut tester un pipeline en mode simulation sur des données historiques avant déploiement en production

**5. Dashboards & Visualisation**

- **FR21** : Le technicien peut visualiser la courbe de santé temporelle de chaque équipement surveillé
- **FR22** : Le technicien peut consulter l'historique des alertes et leur statut
- **FR23** : Le manager peut accéder à une vue synthétique ROI (pannes évitées, économies, disponibilité)
- **FR24** : L'ingénieur peut accéder à une vue technique affichant : liste des pipelines et leur statut, modèles ML actifs et indicateur de confiance, 1 000 dernières entrées logs avec filtrage
- **FR25** : Tout utilisateur peut exporter des données en CSV/JSON
- **FR26** : L'ingénieur peut partager un lien de vue en lecture seule
- **FR48** : Le système affiche un bandeau d'état discret (online / offline / resynchronisation) avec statut de la file locale des actions utilisateur
- **FR49** : Le technicien peut voir la fraîcheur des données (timestamp de la dernière mesure fiable) sur la vue machine et chaque alerte

**6. Gestion des Équipements**

- **FR27** : L'ingénieur peut créer, nommer et archiver des équipements surveillés dans le système

**7. Notifications**

- **FR28** : Le technicien peut recevoir une notification (push, email ou SMS) lors d'une alerte sur un équipement surveillé
- **FR28b** : Le technicien peut configurer par machine les canaux de notification, les niveaux de sévérité déclencheurs et les plages horaires de réception
- **FR51** : Le technicien peut partager un résumé décisionnel (machine, risque, action recommandée, horizon) via WhatsApp/email/copie

**8. Administration & Sécurité**

- **FR29** : L'administrateur peut créer, modifier et supprimer des comptes utilisateurs avec attribution de rôle
- **FR30** : Chaque utilisateur accède à une interface adaptée à son rôle dès la connexion
- **FR31** : Le système enregistre un audit log de toute modification de configuration avec auteur, horodatage et valeur avant/après
- **FR32** : L'administrateur peut déclencher une mise à jour du logiciel depuis le dashboard
- **FR33** : L'administrateur peut effectuer un rollback vers la version précédente
- **FR34** : Le système notifie l'administrateur lorsqu'une nouvelle version est disponible
- **FR35** : Le système distingue les plugins open-source des plugins sous licence propriétaire
- **FR36** : L'administrateur peut consulter l'état de santé du système (services actifs, espace disque, connexions protocoles, uptime)
- **FR37** : Le système peut générer une alerte système lorsqu'une ressource critique atteint un seuil

**9. Onboarding & Déploiement**

- **FR38** : Un nouvel utilisateur peut déployer ML_Elec via une commande Docker unique
- **FR39** : Un nouvel utilisateur peut connecter un premier équipement et visualiser un premier signal actif sur le dashboard dans un délai ≤ 15 min après déploiement
- **FR40** : Le système fournit un guide d'onboarding en 5 étapes max, avec indicateur de progression et fallback données simulées
- **FR41** : L'ingénieur peut déployer ML_Elec dans un environnement totalement isolé (air-gap)
- **FR42** : Le développeur peut accéder à une documentation de contribution et créer un plugin dans un environnement local Docker

**10. Internationalisation**

- **FR43** : L'utilisateur peut sélectionner la langue de l'interface parmi les langues supportées (FR + EN minimum V1)

### Non-Functional Requirements

**Performance**

- **NFR01** : Onboarding complet ≤ 15 minutes (release blocker absolu)
- **NFR02** : Dashboard loading < 2 secondes sur hardware cible
- **NFR03** : Latence détection → notification < 5 secondes
- **NFR04** : Ingestion ≥ 1 000 points/seconde par équipement
- **NFR05** : Inférence ML < 500 ms par cycle
- **NFR06** : Taux de faux positifs < 5% sur datasets de référence
- **NFR28** : Machine prioritaire identifiable en ≤ 3 secondes (≥ 90% sessions)
- **NFR29** : Action utilisateur enregistrée en ≤ 60 secondes (≥ 80% cas)

**Sécurité**

- **NFR07** : Données stockées chiffrées AES-256
- **NFR08** : Communications chiffrées TLS 1.2+
- **NFR09** : RBAC appliqué sur 100% des endpoints API
- **NFR10** : Conformité RGPD (minimisation données, secrets chiffrés)
- **NFR11** : Audit log inaltérable (append-only, 12 mois minimum)

**Fiabilité**

- **NFR12** : Disponibilité ≥ 99,5% (hors maintenance planifiée)
- **NFR13** : Zéro perte de données (write-ahead logging)
- **NFR14** : Dégradation gracieuse LLM (fonctions ML intactes)
- **NFR15** : Reprise automatique < 60 secondes après redémarrage
- **NFR16** : Backup/restore en une commande CLI
- **NFR17** : Fonctionnement offline continu illimité

**Scalabilité**

- **NFR18** : Support 1-50 équipements par instance sans dégradation > 10%
- **NFR19** : Support 1-20 utilisateurs concurrents
- **NFR20** : Rotation automatique des données (rétention 12 mois par défaut)
- **NFR21** : Architecture plugin extensible (intégration < 4 heures)

**Accessibilité**

- **NFR22** : Interface WCAG 2.1 AA (contraste ≥ 4.5:1, navigation clavier)
- **NFR23** : Responsive tablette (≥ 768px), cibles tactiles 44×44 px
- **NFR24** : Architecture i18n (FR + EN V1, ajout langue < 2 heures)

**Intégration**

- **NFR25** : MQTT 3.1.1+ et OPC-UA natifs avec configuration guidée
- **NFR26** : Export CSV/JSON < 10 secondes pour 100 000 lignes
- **NFR27** : API REST documentée OpenAPI 3.0 (endpoints lecture V1)

### Additional Requirements

**Architecture Techniques**

- **Starter Template** : Architecture sur mesure (monorepo pnpm + Turborepo) — aucun starter template standard ne correspond aux spécificités ML_Elec
- **Monorepo Structure** : `apps/web`, `packages/edge-agent`, `packages/shared`, `packages/ui`
- **Backend Edge** : Fastify 4.x + DuckDB + node-opcua + paho-mqtt
- **Frontend** : Vite 6 + React 18 + React Flow 12 + shadcn/ui + Zustand + TanStack Query
- **Tests** : Vitest 4 (unitaires/intégration) + Playwright 1.50+ (E2E)
- **Déploiement** : Docker multi-arch (AMD64 + ARM64)
- **RBAC** : 4 rôles (Viewer/technician, Manager, Engineer, Admin) avec middleware Fastify
- **Stockage** : DuckDB + Parquet partitionné Hive, chiffrement AES-256-GCM
- **Cache** : Hybride Map in-memory (< 1 min) + DuckDB (moyen/long terme)
- **Migrations** : SQL versionné (001_initial_schema.sql, etc.)
- **Logging** : Pino (JSON structuré, niveaux configurables)
- **Monitoring** : Prometheus + Grafana (optionnel V2)

### UX Design Requirements

**Design Tokens & Système Visuel**

- **UX-DR1** : Système de couleurs basé sur le logo ML_Elec : `brand.primary: #F28A00` (orange énergie), `brand.ink: #2F3B46` (anthracite), `brand.steel: #5F6B77` (gris acier)
- **UX-DR2** : Couleurs sémantiques système distinctes de la marque : `success: #2E7D32`, `warning: #B7791F`, `error: #C0392B`, `info: #1F6FEB`
- **UX-DR3** : Typographie IBM Plex Sans (UI) + IBM Plex Mono (technique) avec échelle définie (h1: 32/40 à micro: 12/16)
- **UX-DR4** : Système d'espacement base 8px (scale: 4/8/12/16/24/32/48/64)
- **UX-DR5** : Règle d'usage : orange réservé aux actions/CTA, jamais pour texte fin sur fond clair

**Composants Réutilisables à Créer**

- **UX-DR6** : `HealthTimelineCard` — Afficher état santé machine avec courbe, confiance, horizon, CTA (states: loading, nominal, warning, critical, stale-data, offline)
- **UX-DR7** : `MachinePriorityList` — Trier et prioriser les machines par risque immédiat avec chips filtres, badges sévérité, horizon
- **UX-DR8** : `FocusDecisionPanel` — Mode critique "une machine, une décision" avec CTA unique et résumé partageable
- **UX-DR9** : `PipelineNode` — Composant React Flow pour éditeur no-code (types: source, transformation, ML, sortie)
- **UX-DR10** : `AlertBadge` — Badge de sévérité avec indicateur de confiance ML et contexte
- **UX-DR11** : `OfflineBanner` — Bandeau d'état réseau discret (online/offline/resync) avec file locale visible
- **UX-DR12** : `ROICard` — Vue synthétique direction (pannes évitées, économies, disponibilité)
- **UX-DR13** : `RunbookStep` — Composant d'intervention terrain avec étapes guidées et validation

**Accessibilité & Responsive**

- **UX-DR14** : Contraste WCAG AA minimum 4.5:1 (texte normal), 3:1 (grand texte) avec audit automatisé (score ≥ 90%)
- **UX-DR15** : Focus visible systématique (anneau 2px contrasté) sur tous éléments interactifs
- **UX-DR16** : Responsive tablette (≥ 768px) avec cibles tactiles minimales 44×44 px sur actions primaires
- **UX-DR17** : Support `prefers-reduced-motion` pour toutes animations
- **UX-DR18** : États critiques doublés par icône + texte explicite (jamais couleur seule)

**User Journeys & Flows**

- **UX-DR19** : Journey 1 (Karim - première anomalie) : Dashboard → fiche machine → courbe + confiance + horizon → action planifiée en < 60s
- **UX-DR20** : Journey 2 (Karim - fausse alerte) : Feedback 1 tap → note contexte → escalade Sofiane → analyse → correction éditeur no-code → validation 2 semaines
- **UX-DR21** : Journey 3 (Sofiane - déploiement) : Éditeur pipeline → simulation historique → déploiement Docker → machines visibles → partage vue direction
- **UX-DR22** : Journey 6 (Karim - offline) : Ouverture PWA → badge offline → fonctions core intactes → Q&A désactivé → actions en file locale → resync silencieuse
- **UX-DR23** : Flow de récupération fausse alerte : Feedback → analyse contexte → tuning règle → test simulation → déploiement correctif → mesure réduction faux positifs

**Expérience Utilisateur & Principes**

- **UX-DR24** : Principe "Voir en 3s, décider en 60s" : hiérarchie visuelle stricte, machine prioritaire en haut, CTA unique visible
- **UX-DR25** : Progressive disclosure par rôle RBAC : Karim (simplicité dashboard), Sofiane (profondeur éditeur), M. Benali (vue ROI)
- **UX-DR26** : Offline invisible : fonctions critiques toujours disponibles, LLM désactivable avec message clair, resync silencieuse
- **UX-DR27** : Feedback visible à chaque action : confirmation explicite, statut tracé (lu/en cours/planifié/synchronisé)
- **UX-DR28** : Langage technicien (pas data scientist) : "Vibration anormale roulement gauche" vs "anomalie canal vibratoire axial Z, score 0.87"
- **UX-DR29** : Instrumentation UX pour NFR28/NFR29 : timestamps (ouverture dashboard → machine prioritaire), (ouverture alerte → action)

**Internationalisation**

- **UX-DR30** : Architecture i18n dès V1 : fichiers de traduction FR/EN séparés, ajout langue < 2 heures sans recompilation

### FR Coverage Map

| FR | Epic(s) | Stories |
|----|---------|---------|
| FR01-FR05, FR27, FR38, FR41, FR42, FR44 | Epic 1: Edge Agent - Ingestion & Stockage | 1.1, 1.2, 1.3, 1.4, 1.5 |
| FR06-FR15, FR45-FR46, FR50 | **Epic 2: Détection & Intelligence Augmentée** (ML + LLM fusionnés) | 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7 |
| FR16-FR20, FR47 | Epic 3: Éditeur No-Code - Pipelines Visuels | 3.1, 3.2, 3.3, 3.4 |
| FR21-FR27, FR48-FR49 | Epic 4: Dashboard - Visualisation & Décision | 4.1, 4.2, 4.3, 4.4, 4.5, 4.6 |
| FR28-FR28b, FR51 | Epic 5: Notifications & Partage | 5.1, 5.2 |
| FR29-FR37 | Epic 6: Administration & Sécurité | 6.1, 6.2, 6.3, 6.4 |
| FR39-FR40 | Epic 7: Onboarding & Déploiement (validation NFR01) | 7.1, 7.2, 7.3 |
| FR43 | Epic 8: Internationalisation | 8.1 |

### NFR Coverage Map

| NFR | Category | Epic(s) Impactées |
|-----|----------|-------------------|
| NFR01, NFR39, NFR40 | Performance/Onboarding | Epic 8 |
| NFR02, NFR28, NFR29 | Performance/UX | Epic 5 |
| NFR03, NFR04, NFR05 | Performance/Backend | Epic 1, Epic 2 |
| NFR06 | Performance/ML | Epic 2 |
| NFR07-NFR11 | Sécurité | Epic 7 |
| NFR12-NFR17 | Fiabilité | Epic 1, Epic 7 |
| NFR18-NFR21 | Scalabilité | Epic 1, Epic 4 |
| NFR22-NFR24 | Accessibilité | Epic 5, Epic 9 |
| NFR25-NFR27 | Intégration | Epic 1, Epic 5 |

## Epic List

1. **Epic 1: Edge Agent - Ingestion & Stockage** — Infrastructure backend pour acquisition données temps-réel (MQTT, OPC-UA), stockage DuckDB/Parquet, cache hybride, API REST de base
2. **Epic 2: Détection & Intelligence Augmentée** — Modèles ML (Isolation Forest, Autoencoder, RUL), inférence edge, indicateur de confiance, **explications LLM en langage naturel, Q&A conversationnel**, détection contexte ambigu, suggestions d'actions
3. **Epic 3: Éditeur No-Code - Pipelines Visuels** — Éditeur drag & drop (React Flow), 6 types de plugins, simulation historique, validation avant déploiement
4. **Epic 4: Dashboard - Visualisation & Décision** — Courbe de santé temporelle, triage machines par risque, alertes actionnables, vue ROI direction, export CSV/JSON
5. **Epic 5: Notifications & Partage** — Notifications push/email/SMS, configuration par machine, partage résumé WhatsApp/email, file locale offline
6. **Epic 6: Administration & Sécurité** — RBAC 4 rôles, audit log inaltérable, gestion utilisateurs, mises à jour 1-clic, monitoring système
7. **Epic 7: Onboarding & Déploiement** — Wizard interactif 5 étapes, modèle ML bundled, déploiement Docker unique, onboarding ≤ 15 min (NFR01), air-gap
8. **Epic 8: Internationalisation** — Architecture i18n, fichiers FR/EN, sélection langue interface

---

## Epic Dependencies

```
Epic 1 (Edge Agent)
    │
    ↓
Epic 2 (Détection & Intelligence) ←─┐
    │                                │
    ↓                                │
Epic 3 (Éditeur No-Code) ────────────┤
    │                                │
    ↓                                │
Epic 4 (Dashboard) ←─────────────────┘
    │
    ↓
Epic 5 (Notifications)
    │
    ↓
Epic 6 (Administration)
    │
    ↓
Epic 7 (Onboarding) ← Valide NFR01 sur tous les Épics 1-6
    │
    ↓
Epic 8 (Internationalisation)
```

**Ordre de Sprint Recommandé :**

| Sprint | Epic | Objectif de Validation |
|--------|------|----------------------|
| Sprint 1-2 | Epic 1 | Edge Agent fonctionnel (ingestion + stockage) |
| Sprint 3-5 | Epic 2 | Détection + Explications LLM + Suggestions d'actions |
| Sprint 6-7 | Epic 3 | Éditeur No-Code avec plugins |
| Sprint 8-9 | Epic 4 | Dashboard avec courbe de santé + triage |
| Sprint 10 | Epic 5 | Notifications + Partage |
| Sprint 11 | Epic 6 | Admin + Sécurité + Audit log |
| Sprint 12 | Epic 7 | Onboarding ≤ 15 min (validation NFR01 avec 5 testeurs) |
| Sprint 13 | Epic 8 | i18n FR/EN |

---

## Epic 1: Edge Agent - Ingestion & Stockage

**Goal** : Établir l'infrastructure backend core pour l'acquisition de données temps-réel via MQTT et OPC-UA, le stockage analytique DuckDB/Parquet, et l'API REST de base avec authentification.

### Story 1.1: Setup Monorepo & Infrastructure de Build

**As a** développeur full-stack,
**I want** une structure monorepo pnpm + Turborepo configurée avec TypeScript strict et outils de qualité de code,
**So that** je puisse développer les 4 packages (edge-agent, web, shared, ui) avec des builds incrémentaux et une cohérence architecturale.

**Acceptance Criteria:**

**Given** une machine avec Node.js 20 LTS et pnpm 9.x installés
**When** j'exécute les commandes d'initialisation monorepo
**Then** la structure `apps/web`, `packages/edge-agent`, `packages/shared`, `packages/ui` est créée avec workspaces pnpm configurés
**And** TypeScript strict est activé dans tous les tsconfig.json
**And** ESLint v9 + Prettier v3 sont configurés avec configuration partagée
**And** Turbo build fonctionne sur tous les packages sans erreur

---

### Story 1.2: Entités & Schémas de Validation

**As a** architecte logiciel,
**I want** des interfaces TypeScript et schémas Zod pour toutes les entités métier,
**So that** le typage statique compile-time et la validation runtime soient cohérents sur tout le projet.

**Acceptance Criteria:**

**Given** les 8 entités du modèle de données (Équipement, Capteur, Pipeline, Alerte, ModèleML, Utilisateur, AuditLog, Plugin)
**When** je définis les interfaces TypeScript dans `packages/shared/domain/entities.ts`
**Then** chaque entité a une interface TypeScript avec typage fort
**And** chaque entité a un schéma Zod correspondant dans `packages/shared/schemas/*.schema.ts`
**And** les types TypeScript sont inférés depuis les schémas Zod (`z.infer<>`)
**And** tous les schémas sont exportés depuis `packages/shared/schemas/index.ts`

---

### Story 1.3: Migrations DuckDB & Stockage Parquet

**As a** backend developer,
**I want** des migrations SQL versionnées pour DuckDB et un système de stockage Parquet partitionné,
**So that** le schéma de base de données soit auditable et que les données historiques soient stockées efficacement.

**Acceptance Criteria:**

**Given** DuckDB installé avec chiffrement AES-256-GCM activé
**When** le migration runner exécute les fichiers SQL dans `packages/edge-agent/src/storage/migrations/`
**Then** les tables principales sont créées : `equipements`, `capteurs`, `metrics`, `alertes`, `modeles_ml`, `utilisateurs`, `audit_logs`, `plugins`
**And** les index de performance sont créés sur `timestamp`, `equipement_id`, `severity`
**And** le stockage Parquet est configuré avec partitionnement Hive par date
**And** le write-ahead logging est activé pour garantir zéro perte de données (NFR13)

---

### Story 1.4: Ingestion MQTT & OPC-UA

**As a** ingénieur électrotechnicien,
**I want** connecter des sources de données MQTT et OPC-UA à des pipelines de traitement,
**So that** je puisse acquérir des signaux capteurs en temps-réel depuis mes équipements industriels.

**Acceptance Criteria:**

**Given** un broker MQTT (Mosquitto) et un serveur OPC-UA (open62541) disponibles
**When** je configure une source MQTT (broker URL, topic, fréquence) via l'API
**Then** les données sont ingérées en continu (≥ 1 000 pts/s par NFR04)
**And** les données sont écrites dans DuckDB avec timestamp
**And** la reconnexion automatique fonctionne après coupure réseau (FR05)
**When** je configure une source OPC-UA (endpoint, node IDs, polling interval)
**Then** les données sont pollées et ingérées avec latence < 5s (NFR03)
**And** le cache Map in-memory (< 1 min) est utilisé pour les dernières valeurs
**And** les signaux manquants/hors plage/corrompus sont détectés et signalés (FR03)

---

### Story 1.5: API REST de Base & RBAC Middleware

**As a** administrateur système,
**I want** une API REST documentée avec authentification et RBAC appliqué sur tous les endpoints,
**So that** seuls les utilisateurs autorisés puissent accéder aux données selon leur rôle.

**Acceptance Criteria:**

**Given** 4 rôles définis (Viewer/technician, Manager, Engineer, Admin)
**When** le middleware RBAC Fastify est configuré avec décorateurs TypeScript
**Then** 100% des endpoints API ont une vérification de rôle (NFR09)
**And** la documentation OpenAPI 3.0 est générée automatiquement via @fastify/swagger (FR27, Decision 3.1)
**And** les endpoints de lecture (GET /equipements, /alertes, /metrics) sont accessibles selon RBAC
**And** la gestion des erreurs suit RFC 7807 (Decision 3.2)
**And** le logging JSON structuré avec Pino est activé avec niveaux configurables

---

## Epic 2: Détection & Intelligence Augmentée

**Goal** : Implémenter les modèles de Machine Learning pour la détection d'anomalies (Isolation Forest, Autoencoder), le calcul de RUL, **ET** intégrer un module LLM pour générer des explications en langage naturel, des suggestions d'actions concrètes et permettre un dialogue conversationnel avec les données des équipements.

### Story 2.1: Modèle ML Généraliste Bundled

**As a** nouvel utilisateur,
**I want** un modèle ML généraliste pré-entraîné inclus dans l'image Docker,
**So that** je puisse démarrer la détection d'anomalies immédiatement sans configuration ni entraînement préalable.

**Acceptance Criteria:**

**Given** l'image Docker ML_Elec fraîchement installée
**When** je démarre le système pour la première fois
**Then** le modèle généraliste (Isolation Forest + Autoencoder) est disponible immédiatement
**And** le modèle exécute sur 3 datasets de référence (CWRU Bearing, NASA Turbofan, FEMTO) avec détection ≥ 70% anomalies et faux positifs < 5% (FR46, NFR06)
**And** aucun entraînement préalable n'est requis
**And** le modèle est exporté au format ONNX pour portabilité

---

### Story 2.2: Détection d'Anomalies & Alertes

**As a** technicien de maintenance,
**I want** que le système détecte automatiquement des anomalies sur les signaux de mes équipements,
**So that** je sois alerté avant une défaillance critique.

**Acceptance Criteria:**

**Given** des données capteurs ingérées en continu (vibration, température, courant)
**When** le modèle Isolation Forest exécute un cycle d'inférence (< 500ms par NFR05)
**Then** un score d'anomalie est calculé pour chaque point de données
**And** une alerte est générée avec niveau de sévérité (info/warning/critique) quand le score dépasse un seuil (FR08)
**And** l'alerte inclut un indicateur de confiance du modèle (FR11)
**And** le contexte ambigu est signalé si applicable (température ambiante élevée, capteur récent, modèle en apprentissage) (FR50)
**And** la latence bout-en-bout détection → notification est < 5 secondes (NFR03)

---

### Story 2.3: Calcul RUL & Horizon de Défaillance

**As a** ingénieur de maintenance,
**I want** connaître la durée de vie résiduelle (RUL) et l'horizon de défaillance estimé de mes équipements,
**So that** je puisse planifier les interventions de manière préventive.

**Acceptance Criteria:**

**Given** un historique de données de 7+ jours pour un équipement
**When** le modèle RUL exécute
**Then** une estimation de durée de vie résiduelle est calculée (en jours/heures)
**And** un horizon de défaillance est affiché sur l'interface (ex: "défaillance probable dans 4-6 jours")
**And** l'incertitude du modèle est affichée (intervalle de confiance)
**And** le calcul RUL s'exécute en < 500ms (NFR05)

---

### Story 2.4: Configuration Seuils & Règles de Corrélation

**As a** ingénieur électrotechnicien,
**I want** configurer des seuils d'alerte et des règles de corrélation,
**So that** je puisse réduire les faux positifs et adapter le modèle à mon contexte.

**Acceptance Criteria:**

**Given** des alertes générées par le modèle ML
**When** j'accède à la configuration des seuils via l'interface Engineer
**Then** je peux ajuster les seuils de sévérité (info/warning/critique) par type d'équipement
**And** je peux créer des règles de corrélation (ex: "Ignorer alerte thermique si température ambiante > 35°C")
**And** les modifications de seuils critiques nécessitent une confirmation explicite avec avertissement (MOC simplifié)
**And** toutes les modifications sont tracées dans l'audit log avec auteur, horodatage, valeur avant/après (FR31, NFR11)
**And** je peux tester les règles en mode simulation sur données historiques (FR20)

---

### Story 2.5: Feedback Utilisateur & Boucle de Confiance

**As a** technicien de maintenance,
**I want** soumettre un feedback sur une alerte en 1 tap (utile / fausse alerte),
**So that** le système s'améliore et que je puisse signaler une erreur sans friction.

**Acceptance Criteria:**

**Given** une alerte affichée sur le dashboard
**When** je clique sur le bouton de feedback "Utile" ou "Fausse alerte"
**Then** le feedback est enregistré immédiatement avec confirmation visuelle
**And** je peux ajouter une note de contexte optionnelle (ex: "Machine à l'arrêt", "Maintenance en cours")
**And** le feedback est visible par l'ingénieur pour analyse (Journey 2)
**And** le ratio utile/fausse alerte est affiché dans le tableau de bord de qualité des alarmes
**And** le feedback contribue à l'ajustement des seuils (via règles de corrélation)

---

### Story 2.6: Intégration Providers LLM & Configuration

**As a** administrateur système,
**I want** configurer le fournisseur LLM (OpenAI, Anthropic, Mistral) et activer/désactiver le module,
**So that** je puisse choisir le provider et contrôler la dépendance API.

**Acceptance Criteria:**

**Given** des clés API pour OpenAI, Anthropic ou Mistral
**When** j'accède à la configuration LLM dans le dashboard Admin
**Then** je peux sélectionner le provider et saisir la clé API
**And** la clé API est stockée chiffrée (NFR10)
**And** je peux activer ou désactiver le module LLM (FR14)
**And** quand le module est désactivé, toutes les fonctions ML restent opérationnelles (NFR14)
**And** le statut LLM (actif/inactif) est visible dans le bandeau d'état système

---

### Story 2.7: Explications LLM & Q&A Conversationnel

**As a** technicien de maintenance,
**I want** lire une explication en langage naturel pour chaque anomalie détectée et poser des questions sur l'état des équipements,
**So that** je comprenne la cause probable, l'action à entreprendre, et puisse dialoguer avec mes données sans expertise ML.

**Acceptance Criteria:**

**Given** une alerte générée avec score ML, contexte et données capteurs
**When** le module LLM est actif
**Then** une explication est générée en langage naturel (ex: "Vibration anormale détectée sur roulement gauche — dégradation progressive, probabilité de défaillance dans 4 à 6 jours") (FR12)
**And** l'explication utilise le vocabulaire technicien (pas de jargon data science) (UX-DR28)
**And** le temps de lecture moyen est < 30 secondes
**And** le taux de compréhension est ≥ 80% (mesuré par test utilisateur avec 5 techniciens BTS)
**And** l'explication inclut cause probable, sévérité et horizon temporel
**When** je pose une question en langage naturel (ex: "Pourquoi le score de M-03 baisse ?") (FR13)
**Then** une réponse est générée en < 10 secondes basée sur les données de l'équipement
**And** le système indique clairement lorsqu'une question dépasse sa capacité
**And** le Q&A est désactivé automatiquement quand le module LLM est offline ou désactivé
**And** un fallback vers les courbes et métriques est affiché en mode offline (UX-DR26)

---

### Story 2.8: Suggestion d'Action Concrète

**As a** technicien de maintenance,
**I want** voir une suggestion d'action concrète pour chaque alerte,
**So that** je sache quoi faire immédiatement sans hésitation.

**Acceptance Criteria:**

**Given** une alerte générée avec explication LLM
**When** j'ouvre le détail de l'alerte
**Then** je vois une suggestion d'action concrète (ex: "Planifier une intervention jeudi", "Commander la pièce X", "Surveiller 48h") (FR45)
**And** la suggestion est copyable en 1 tap pour partage WhatsApp/email
**And** je peux sélectionner l'action choisie (planifier / surveiller / partager)
**And** l'action est tracée avec timestamp (NFR29 : décision ≤ 60 secondes)
**And** le contexte de l'alerte est affiché (données capteurs, confiance ML, facteurs environnementaux)

---

## Epic 3: Éditeur No-Code - Pipelines Visuels

**Goal** : Créer un éditeur visuel de pipelines par drag & drop permettant de configurer des flux de traitement de données sans écrire de code, avec un système de plugins extensible.

### Story 3.1: Éditeur React Flow - Structure de Base

**As a** ingénieur électrotechnicien,
**I want** créer un pipeline de traitement par drag & drop de nœuds,
**So that** je puisse configurer l'acquisition et le traitement des données sans code.

**Acceptance Criteria:**

**Given** l'éditeur React Flow 12 installé dans `apps/web`
**When** j'ouvre l'éditeur de pipelines
**Then** je vois une palette de nœuds typés par couleur (sources en bleu, transformations en violet, ML en orange, sorties en vert)
**And** je peux drag & drop des nœuds depuis la palette vers le canvas
**And** les connexions entre nœuds sont visuelles et logiques (connexions invalides impossibles)
**And** chaque nœud a un type sémantique clair : "Source MQTT", "Filtre harmonique", "Modèle Isolation Forest", "Seuil d'alerte", "Sortie dashboard"
**And** le pipeline est sauvegardé automatiquement avec validation de cohérence

---

### Story 3.2: Système de Plugins - 6 Types de Contrats

**As a** développeur externe,
**I want** créer et publier un plugin selon l'architecture documentée,
**So that** je puisse étendre ML_Elec avec de nouveaux protocoles, modèles ou visualisations.

**Acceptance Criteria:**

**Given** la documentation des 6 types de plugins (FR47)
**When** je développe un plugin dans un environnement local Docker
**Then** chaque type de plugin a un contrat d'interface documenté : (1) connecteurs protocole, (2) transformateurs de données, (3) modèles ML, (4) blocs de visualisation, (5) canaux de notification, (6) exportateurs
**And** un template de plugin est fourni pour chaque type
**And** les plugins open-source sont distingués des plugins propriétaires lors de l'installation (FR35)
**And** un plugin peut être installé en un clic depuis la bibliothèque (FR18)
**And** le temps d'intégration d'un nouveau plugin est < 4 heures pour un développeur familier avec l'API (NFR21)

---

### Story 3.3: Simulation & Validation avant Déploiement

**As a** ingénieur électrotechnicien,
**I want** tester un pipeline en mode simulation sur des données historiques,
**So that** je puisse valider le comportement avant déploiement en production.

**Acceptance Criteria:**

**Given** un pipeline configuré avec des nœuds connectés
**When** je clique sur le bouton "Test Simulation"
**Then** le pipeline s'exécute sur des données historiques (import CSV ou données existantes)
**And** chaque nœud affiche son statut en temps réel (OK / erreur / données traitées)
**And** je vois les données transformées à chaque étape
**And** un rapport de validation est généré (taux de détection, faux positifs, performance)
**And** je peux modifier le pipeline et re-tester sans redéploiement (FR17)
**And** une fois validé, je peux déployer en production en un clic

---

### Story 3.4: Import CSV & Fallback Protocole Non Supporté

**As a** ingénieur électrotechnicien,
**I want** importer un fichier CSV de données historiques pour valider le modèle ML,
**So that** je puisse tester ML_Elec avant connexion au matériel réel.

**Acceptance Criteria:**

**Given** un fichier CSV de données capteurs (timestamp, node_id, value)
**When** j'importe le CSV via l'interface
**Then** les données sont validées (colonnes requises, format timestamp, valeurs numériques)
**And** les données sont chargées dans DuckDB
**And** je peux exécuter le pipeline ML sur ces données historiques (FR44)
**And** je peux visualiser les alertes générées rétroactivement
**And** si un protocole n'est pas supporté nativement, je peux utiliser un fallback adaptateur CSV ou créer une issue pour un plugin communautaire (Journey 3)

---

## Epic 4: Dashboard - Visualisation & Décision

**Goal** : Créer l'interface utilisateur de décision avec courbe de santé temporelle, triage des machines par risque, et vue ROI pour la direction.

### Story 4.1: Courbe de Santé Temporelle (HealthTimelineCard)

**As a** technicien de maintenance,
**I want** visualiser la courbe de santé temporelle de chaque équipement surveillé,
**So that** je puisse voir l'évolution de la santé de ma machine dans le temps.

**Acceptance Criteria:**

**Given** un équipement avec historique de données et scores ML
**When** j'ouvre la fiche machine
**Then** je vois une courbe temporelle de santé (score sur 7-30 jours)
**And** la courbe montre la trajectoire (dégradation/stable/amélioration)
**And** les alertes passées sont marquées sur la courbe
**And** le RUL et l'horizon de défaillance sont affichés
**And** la fraîcheur des données est visible (timestamp dernière mesure fiable) (FR49)
**And** le composant `HealthTimelineCard` supporte les states : loading, nominal, warning, critical, stale-data, offline (UX-DR6)

---

### Story 4.2: Triage Machines par Risque (MachinePriorityList)

**As a** technicien de maintenance,
**I want** voir la liste de mes machines triées par risque immédiat,
**So that** je puisse identifier la machine prioritaire en moins de 3 secondes.

**Acceptance Criteria:**

**Given** plusieurs équipements surveillés avec scores de santé
**When** j'ouvre le dashboard principal
**Then** les machines sont triées par risque (critique → warning → info → nominal)
**And** la machine prioritaire est identifiable en ≤ 3 secondes (NFR28)
**And** chaque ligne machine affiche : nom, score de santé, badge de sévérité, horizon de défaillance
**And** des chips de filtres permettent de filtrer par sévérité ou type d'équipement
**And** le composant `MachinePriorityList` est responsive tablette (≥ 768px) avec cibles tactiles 44×44 px (UX-DR7, NFR23)

---

### Story 4.3: Alertes Actionnables & Suggestion d'Action

**As a** technicien de maintenance,
**I want** voir une suggestion d'action concrète pour chaque alerte,
**So that** je sache quoi faire immédiatement sans hésitation.

**Acceptance Criteria:**

**Given** une alerte générée avec explication LLM
**When** j'ouvre le détail de l'alerte
**Then** je vois une suggestion d'action concrète (ex: "Planifier une intervention jeudi", "Commander la pièce X", "Surveiller 48h") (FR45)
**And** la suggestion est copyable en 1 tap pour partage WhatsApp/email
**And** je peux sélectionner l'action choisie (planifier / surveiller / partager)
**And** l'action est tracée avec timestamp (NFR29 : décision ≤ 60 secondes)
**And** le contexte de l'alerte est affiché (données capteurs, confiance ML, facteurs environnementaux)

---

### Story 4.4: Vue Synthétique ROI (ROICard)

**As a** responsable de production (M. Benali),
**I want** accéder à une vue synthétique avec pannes évitées et économies réalisées,
**So that** je puisse voir le ROI de ML_Elec en 10 secondes.

**Acceptance Criteria:**

**Given** des alertes traitées et interventions planifiées
**When** j'ouvre la vue "Direction"
**Then** je vois un résumé avec : pannes évitées (nombre), économies estimées (€), disponibilité machines (+%)
**And** les chiffres sont lisibles en gros, sans jargon technique
**And** la vue est conçue pour être montrée en réunion (partage écran)
**And** le composant `ROICard` calcule automatiquement le ROI à chaque panne évitée (UX-DR12)
**And** un export PDF / partage lecture seule est disponible (FR23, FR26)

---

### Story 4.5: Export Données & Partage Vue Lecture Seule

**As a** ingénieur électrotechnicien,
**I want** exporter des données en CSV/JSON et partager un lien de vue en lecture seule,
**So that** je puisse analyser les données hors ligne et collaborer avec mes collègues.

**Acceptance Criteria:**

**Given** des données d'équipements, alertes ou metrics
**When** je clique sur "Exporter"
**Then** je peux choisir le format (CSV ou JSON)
**And** l'export de 100 000 lignes se fait en < 10 secondes (NFR26)
**And** je peux générer un lien de vue en lecture seule (FR26)
**And** le lien a une durée de validité configurable
**And** le lien peut être révoqué à tout moment
**And** la vue partagée affiche les données mais masque les actions de configuration

---

### Story 4.6: Gestion des Équipements

**As a** ingénieur électrotechnicien,
**I want** créer, nommer et archiver des équipements surveillés,
**So that** je puisse organiser mon parc de machines dans ML_Elec.

**Acceptance Criteria:**

**Given** accès à l'interface Engineer
**When** je crée un nouvel équipement
**Then** je peux saisir le nom, le type (moteur/pompe/compresseur), et les métadonnées
**And** l'équipement est enregistré dans DuckDB
**And** je peux archiver un équipement (masqué du dashboard mais historique conservé)
**And** je peux assigner des capteurs à l'équipement
**And** l'équipement apparaît dans la liste des machines surveillées

---

## Epic 5: Notifications & Partage

**Goal** : Mettre en place le système de notifications push/email/SMS et le partage de résumés décisionnels via WhatsApp/email.

### Story 5.1: Configuration Notifications par Machine

**As a** technicien de maintenance,
**I want** configurer les canaux de notification et les niveaux de sévérité par machine,
**So that** je reçoive les alertes importantes sans être submergé.

**Acceptance Criteria:**

**Given** des équipements surveillés
**When** j'accède à la configuration des notifications
**Then** je peux sélectionner les canaux (push, email, SMS) par machine
**And** je peux définir les niveaux de sévérité déclencheurs (ex: warning et critique uniquement)
**And** je peux configurer les plages horaires de réception (ex: pas de notifications la nuit)
**And** les notifications incluent : machine, sévérité, horizon, lien direct vers la fiche (FR28)
**And** je peux activer/désactiver les notifications globalement

---

### Story 5.2: Partage Résumé Décisionnel & File Locale Offline

**As a** technicien de maintenance,
**I want** partager un résumé décisionnel via WhatsApp/email,
**So that** je puisse collaborer avec mon équipe et ma hiérarchie.

**Acceptance Criteria:**

**Given** une alerte avec action recommandée
**When** je clique sur "Partager"
**Then** un résumé est généré avec : machine, risque, action recommandée, horizon de défaillance
**And** je peux copier le résumé ou l'envoyer via WhatsApp/email (FR51)
**And** en mode offline, l'envoi est mis en file locale
**And** la resynchronisation est silencieuse au retour du réseau
**And** le bandeau d'état offline affiche le statut de la file locale (FR48, UX-DR11)

---

## Epic 6: Administration & Sécurité

**Goal** : Implémenter la gestion des utilisateurs RBAC, l'audit log inaltérable, les mises à jour logicielles et le monitoring système.

### Story 6.1: Gestion Utilisateurs & RBAC

**As a** administrateur système,
**I want** créer, modifier et supprimer des comptes utilisateurs avec attribution de rôle,
**So that** chaque utilisateur accède à l'interface adaptée à son profil.

**Acceptance Criteria:**

**Given** le rôle Admin
**When** j'accède à la gestion des utilisateurs
**Then** je peux créer un utilisateur avec username, password et rôle (Viewer, Manager, Engineer, Admin)
**And** je peux modifier le rôle d'un utilisateur existant
**And** je peux supprimer un utilisateur (sauf le dernier Admin)
**And** chaque utilisateur accède à une interface adaptée à son rôle dès la connexion (FR30)
**And** les mots de passe sont stockés chiffrés (bcrypt)

---

### Story 6.2: Audit Log Inaltérable

**As a** administrateur système,
**I want** que toute modification de configuration soit tracée dans un audit log inaltérable,
**So that** je puisse auditer les changements et assurer la traçabilité OT.

**Acceptance Criteria:**

**Given** des modifications de configuration (seuils, règles, pipelines, utilisateurs)
**When** une modification est effectuée
**Then** une entrée est créée dans la table `audit_logs` avec : auteur, horodatage, entité, valeur avant, valeur après, justification optionnelle (FR31)
**And** l'audit log est append-only (non modifiable, non supprimable, même par Admin) (NFR11)
**And** la rétention est de 12 mois minimum
**And** je peux consulter l'audit log avec filtrage par utilisateur, entité, période
**And** l'intégrité de l'audit log est vérifiable par checksum

---

### Story 6.3: Mises à Jour 1-Clic & Rollback

**As a** administrateur système,
**I want** déclencher une mise à jour du logiciel depuis le dashboard et effectuer un rollback,
**So that** je puisse maintenir le système à jour sans interruption de service.

**Acceptance Criteria:**

**Given** une nouvelle version Docker disponible
**When** je clique sur "Mettre à jour" dans le dashboard Admin
**Then** la nouvelle image Docker est téléchargée
**And** le conteneur est redémarré avec la nouvelle version
**And** un rollback vers la version précédente est possible en un clic (FR33)
**And** une notification est affichée quand une nouvelle version est disponible (FR34)
**And** le mécanisme de mise à jour fonctionne offline depuis un registre local (air-gap) (FR41)

---

### Story 6.4: Monitoring Système & Alertes Ressources

**As a** administrateur système,
**I want** consulter l'état de santé du système et recevoir des alertes ressources critiques,
**So that** je puisse anticiper les problèmes d'infrastructure.

**Acceptance Criteria:**

**Given** le système en production
**When** j'ouvre la vue "Administration"
**Then** je vois : services actifs, espace disque, mémoire, connexions protocoles, uptime (FR36)
**And** des alertes système sont générées quand une ressource critique atteint un seuil (ex: disque > 90%, mémoire > 85%, déconnexion broker) (FR37)
**And** les métriques système sont exposées via endpoint `/metrics` (pour Prometheus V2) (Decision 5.1)
**And** le logging Pino est configuré avec niveaux (info, warn, error) et transport JSON structuré

---

## Epic 7: Onboarding & Déploiement

**Goal** : Créer un parcours d'onboarding interactif en 5 étapes permettant à un nouvel utilisateur de déployer ML_Elec et de visualiser un premier signal en ≤ 15 minutes.

### Story 7.1: Wizard d'Onboarding Interactif 5 Étapes

**As a** nouvel utilisateur,
**I want** un guide d'onboarding en 5 étapes avec indicateur de progression,
**So that** je puisse compléter l'installation sans aide externe.

**Acceptance Criteria:**

**Given** une installation fraîche de ML_Elec
**When** je me connecte pour la première fois
**Then** un wizard interactif se lance avec 5 étapes max : (1) Bienvenue, (2) Connexion source de données, (3) Configuration pipeline, (4) Visualisation dashboard, (5) Première alerte (FR40)
**And** un indicateur de progression est affiché
**And** un fallback avec données simulées est disponible si aucune source réelle n'est connectée
**And** le taux de complétion onboarding sans aide externe est ≥ 90% (mesuré par test utilisateur)

---

### Story 7.2: Déploiement Docker Unique & Air-Gap

**As a** ingénieur électrotechnicien,
**I want** déployer ML_Elec via une commande Docker unique et dans un environnement isolé,
**So that** je puisse installer le système rapidement sans dépendance internet.

**Acceptance Criteria:**

**Given** une machine avec Docker pré-installé
**When** j'exécute la commande de déploiement (`docker compose up -d`)
**Then** ML_Elec est déployé avec tous les services (edge-agent, web, DuckDB) (FR38)
**And** le déploiement fonctionne en environnement air-gap (totally isolated) depuis un registre local (FR41)
**And** la taille de l'image Docker edge est ≤ 500 MB
**And** l'image est multi-arch (AMD64 + ARM64) pour Raspberry Pi et Intel NUC

---

### Story 7.3: Modèle ML Bundled & Premier Signal ≤ 15 min

**As a** nouvel utilisateur,
**I want** connecter un premier équipement et visualiser un premier signal actif en ≤ 15 minutes,
**So that** je puisse valider la valeur de ML_Elec immédiatement.

**Acceptance Criteria:**

**Given** une installation fraîche avec modèle ML généraliste bundled (FR46)
**When** je suis le wizard d'onboarding
**Then** je peux connecter une source de données (MQTT/OPC-UA/CSV) en ≤ 5 minutes
**And** je peux visualiser un premier signal actif sur le dashboard en ≤ 15 minutes (FR39, NFR01)
**And** le modèle ML généraliste fonctionne immédiatement sans entraînement
**And** le test utilisateur avec 5 ingénieurs électrotechniques confirme que 5/5 complètent en ≤ 15 min

---

## Epic 8: Internationalisation

**Goal** : Mettre en place l'architecture i18n permettant d'ajouter une nouvelle langue en < 2 heures sans recompilation.

### Story 8.1: Architecture i18n & Fichiers de Traduction FR/EN

**As a** utilisateur international,
**I want** sélectionner la langue de l'interface parmi les langues supportées,
**So that** je puisse utiliser ML_Elec dans ma langue maternelle.

**Acceptance Criteria:**

**Given** l'architecture i18n configurée (react-i18next ou équivalent)
**When** je me connecte à l'application
**Then** je peux sélectionner la langue (FR + EN minimum V1) depuis un menu déroulant (FR43, FR24)
**And** les fichiers de traduction sont séparés (`locales/fr/translation.json`, `locales/en/translation.json`)
**And** l'ajout d'une nouvelle langue prend < 2 heures pour un traducteur non-technique (NFR24)
**And** la langue sélectionnée est persistée dans les préférences utilisateur
**And** tous les composants UI (dashboard, alertes, éditeur, Q&A) sont traduits
**And** le langage technique (IDs, logs, erreurs) reste en anglais pour cohérence

---

## Appendix: Requirements Coverage

### FR Coverage Summary

| FR | Status | Epic/Stories |
|----|--------|--------------|
| FR01-FR05, FR27, FR38, FR41, FR42, FR44 | ✅ Covered | Epic 1 (1.1-1.5) |
| FR06-FR15, FR45-FR46, FR50 | ✅ Covered | **Epic 2 (2.1-2.8)** — Détection & Intelligence Augmentée |
| FR16-FR20, FR47 | ✅ Covered | Epic 3 (3.1-3.4) — Éditeur No-Code |
| FR21-FR27, FR48-FR49 | ✅ Covered | Epic 4 (4.1-4.6) — Dashboard |
| FR28-FR28b, FR51 | ✅ Covered | Epic 5 (5.1-5.2) — Notifications |
| FR29-FR37 | ✅ Covered | Epic 6 (6.1-6.4) — Administration |
| FR39-FR40 | ✅ Covered | Epic 7 (7.1-7.3) — Onboarding |
| FR43 | ✅ Covered | Epic 8 (8.1) — Internationalisation |

### NFR Coverage Summary

| NFR | Status | Epic/Stories |
|-----|--------|--------------|
| NFR01, NFR39, NFR40 | ✅ Covered | Epic 7 (7.1, 7.3) — Onboarding |
| NFR02, NFR28, NFR29 | ✅ Covered | Epic 4 (4.1, 4.2, 4.3) — Dashboard |
| NFR03-NFR06 | ✅ Covered | Epic 1 (1.4), Epic 2 (2.2, 2.3) — Edge + ML |
| NFR07-NFR11 | ✅ Covered | Epic 6 (6.1, 6.2) — Administration |
| NFR12-NFR17 | ✅ Covered | Epic 1 (1.3, 1.5), Epic 6 (6.3) — Edge + Admin |
| NFR18-NFR21 | ✅ Covered | Epic 1 (1.4), Epic 3 (3.2) — Edge + Plugins |
| NFR22-NFR24 | ✅ Covered | Epic 4 (4.1, 4.2), Epic 8 (8.1) — Dashboard + i18n |
| NFR25-NFR27 | ✅ Covered | Epic 1 (1.4, 1.5), Epic 4 (4.5) — Edge + Dashboard |

### UX-DR Coverage Summary

| UX-DR | Status | Epic/Stories |
|-------|--------|--------------|
| UX-DR1-UX-DR5 (Design Tokens) | ✅ Covered | Epic 4 (4.1-4.4) — Dashboard |
| UX-DR6-UX-DR13 (Custom Components) | ✅ Covered | Epic 4 (4.1-4.4, 4.6) — Dashboard |
| UX-DR14-UX-DR18 (Accessibility) | ✅ Covered | Epic 4 (4.1, 4.2), Epic 8 (8.1) — Dashboard + i18n |
| UX-DR19-UX-DR23 (User Journeys) | ✅ Covered | Epic 4 (4.1-4.3), Epic 3 (3.3), Epic 5 (5.2) — Dashboard + Éditeur + Notifications |
| UX-DR24-UX-DR29 (UX Principles) | ✅ Covered | Epic 4 (4.1-4.3), Epic 2 (2.7) — Dashboard + LLM |
| UX-DR30 (i18n) | ✅ Covered | Epic 8 (8.1) — Internationalisation |

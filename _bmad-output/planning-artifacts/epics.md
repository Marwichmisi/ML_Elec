---
stepsCompleted: ["step-01-validate-prerequisites", "step-02-design-epics", "step-03-create-stories", "step-04-final-validation"]
inputDocuments:
  - "_bmad-output/planning-artifacts/prd.md"
  - "_bmad-output/planning-artifacts/architecture.md"
  - "_bmad-output/planning-artifacts/ux-design-specification.md"
workflowStatus: 'complete'
validationDate: '2026-03-23'
---

# ML_Elec - Epic Breakdown

## 🎯 Workflow Status : COMPLET ✅

**Document prêt pour le développement.**

- **49 stories** réparties sur **9 epics**
- **52 FRs** couvertes (100%)
- **29 NFRs** assignés
- **26 UX-DRs** implémentées
- **0 dépendance future** détectée
- **Validation finale** : 2026-03-23

## Requirements Inventory

### Functional Requirements

**1. Acquisition & Ingestion de Données**

FR01 : L'ingénieur peut connecter une source de données MQTT à un pipeline de traitement
FR02 : L'ingénieur peut connecter une source de données OPC-UA à un pipeline de traitement
FR03 : Le système peut détecter et signaler dans l'interface les signaux de données manquants, hors plage ou corrompus
FR04 : L'ingénieur peut configurer les paramètres de connexion (broker, topic, fréquence) sans modifier de code
FR44 : L'ingénieur peut importer un fichier CSV de données historiques pour valider le modèle ML avant connexion au matériel réel
FR05 : Le système peut ingérer des données en continu (24h/7j) avec reprise automatique après interruption, sans perte de données

**2. Détection d'Anomalies & Prédiction**

FR06 : Le système peut détecter des anomalies sur les signaux d'un équipement via le modèle généraliste par défaut
FR46 : Le système est livré avec un modèle ML généraliste pré-entraîné (bundled dans l'image Docker), fonctionnel dès l'installation sans configuration ni entraînement préalable
FR07 : Le système peut calculer une durée de vie résiduelle (RUL) et un horizon de défaillance estimé
FR08 : Le système peut générer une alerte avec niveau de sévérité (info / warning / critique)
FR09 : L'ingénieur peut configurer des seuils d'alerte et des règles de corrélation pour réduire les faux positifs
FR10 : Le technicien peut soumettre un feedback sur une alerte en 1 tap (utile / fausse alerte), avec note de contexte optionnelle et confirmation immédiate
FR11 : Le système peut afficher un indicateur de confiance du modèle ML sur chaque alerte
FR50 : Le système peut signaler un contexte ambigu avant action (température ambiante élevée, capteur récemment ajouté, modèle en phase d'apprentissage)
FR45 : Le système peut générer une suggestion d'action concrète pour chaque alerte (commander la pièce, planifier l'arrêt, surveiller)

**3. Explication & Intelligence (LLM)**

FR12 : Le système peut générer une explication en langage naturel pour chaque anomalie détectée, compréhensible par un technicien BTS sans formation ML
FR13 : Le technicien peut poser une question en langage naturel sur l'état d'un équipement et recevoir une réponse en < 10s
FR14 : L'administrateur peut activer ou désactiver le module LLM
FR15 : L'administrateur peut configurer le fournisseur LLM et la clé API

**4. Éditeur No-Code & Pipelines**

FR16 : L'ingénieur peut créer un pipeline de traitement par drag & drop sans code
FR17 : L'ingénieur peut modifier un pipeline existant sans redéploiement
FR18 : L'ingénieur peut installer un plugin depuis la bibliothèque en un clic
FR19 : Le développeur peut créer et publier un plugin selon l'architecture documentée
FR47 : Le système supporte 6 types de plugins avec contrats d'interface documentés
FR20 : L'ingénieur peut tester un pipeline en mode simulation sur des données historiques avant déploiement en production

**5. Dashboards & Visualisation**

FR21 : Le technicien peut visualiser la courbe de santé temporelle de chaque équipement surveillé
FR22 : Le technicien peut consulter l'historique des alertes et leur statut
FR23 : Le manager peut accéder à une vue synthétique ROI (pannes évitées, économies, disponibilité)
FR24 : L'ingénieur peut accéder à une vue technique affichant : liste des pipelines, modèles ML actifs, logs système
FR25 : Tout utilisateur peut exporter des données en CSV/JSON
FR26 : L'ingénieur peut partager un lien de vue en lecture seule
FR48 : Le système affiche un bandeau d'état discret (online / offline / resynchronisation) avec statut de la file locale
FR49 : Le technicien peut voir la fraîcheur des données (timestamp de la dernière mesure fiable) sur la vue machine et chaque alerte

**6. Gestion des Équipements**

FR27 : L'ingénieur peut créer, nommer et archiver des équipements surveillés dans le système

**7. Notifications**

FR28 : Le technicien peut recevoir une notification (push, email ou SMS) lors d'une alerte sur un équipement surveillé
FR28b : Le technicien peut configurer par machine les canaux de notification, les niveaux de sévérité et les plages horaires
FR51 : Le technicien peut partager un résumé décisionnel via WhatsApp/email/copie avec file locale offline

**8. Administration & Sécurité**

FR29 : L'administrateur peut créer, modifier et supprimer des comptes utilisateurs avec attribution de rôle
FR30 : Chaque utilisateur accède à une interface adaptée à son rôle dès la connexion
FR31 : Le système enregistre un audit log de toute modification de configuration avec auteur, horodatage et valeur avant/après
FR32 : L'administrateur peut déclencher une mise à jour du logiciel depuis le dashboard
FR33 : L'administrateur peut effectuer un rollback vers la version précédente
FR34 : Le système notifie l'administrateur lorsqu'une nouvelle version est disponible
FR35 : Le système distingue les plugins open-source des plugins sous licence propriétaire
FR36 : L'administrateur peut consulter l'état de santé du système (services actifs, espace disque, connexions)
FR37 : Le système peut générer une alerte système lorsqu'une ressource critique atteint un seuil

**9. Onboarding & Déploiement**

FR38 : Un nouvel utilisateur peut déployer ML_Elec via une commande Docker unique
FR39 : Un nouvel utilisateur peut connecter un premier équipement et visualiser un premier signal actif en ≤ 15 min
FR40 : Le système fournit un guide d'onboarding en 5 étapes max avec indicateur de progression
FR41 : L'ingénieur peut déployer ML_Elec dans un environnement totalement isolé (air-gap)
FR42 : Le développeur peut accéder à une documentation de contribution et créer un plugin en local Docker

**10. Internationalisation**

FR43 : L'utilisateur peut sélectionner la langue de l'interface parmi les langues supportées (FR + EN minimum V1)

### NonFunctional Requirements

**Performance**

NFR01 : L'onboarding complet (déploiement Docker → premier signal actif affiché) est réalisable en ≤ 15 minutes — release blocker absolu
NFR02 : Les pages dashboard se chargent en < 2 secondes sur hardware cible (mini-PC entrée de gamme)
NFR03 : Latence entre détection d'anomalie et notification utilisateur < 5 secondes
NFR04 : Ingestion de données continue 24h/7j à ≥ 1 000 points/seconde par équipement sans perte
NFR05 : Inférence ML exécutée en < 500 ms par cycle d'analyse sur hardware cible
NFR06 : Taux de faux positifs du modèle généraliste par défaut < 5% sur les datasets de référence
NFR28 : Sur le dashboard opérationnel, la machine prioritaire est identifiable en ≤ 3 secondes dans ≥ 90% des sessions
NFR29 : Après ouverture d'une alerte warning/critique, une action utilisateur est enregistrée en ≤ 60 secondes dans ≥ 80% des cas

**Sécurité**

NFR07 : Données stockées chiffrées au repos via AES-256
NFR08 : Toutes les communications chiffrées en transit : TLS 1.2+ pour API REST, MQTT over TLS, OPC-UA Security Mode
NFR09 : RBAC appliqué sur 100% des endpoints API — aucune escalade de privilèges possible
NFR10 : Conformité RGPD : aucune donnée personnelle ni industrielle transmise hors instance locale sans consentement
NFR11 : Audit log inaltérable (append-only), conservé minimum 12 mois

**Fiabilité**

NFR12 : Disponibilité système ≥ 99,5% (hors maintenance planifiée)
NFR13 : Zéro perte de données en cas de coupure électrique brutale — write-ahead logging
NFR14 : Dégradation gracieuse si LLM indisponible — toutes les fonctions ML restent opérationnelles
NFR15 : Reprise automatique après redémarrage matériel sans intervention en < 60 secondes
NFR16 : Backup et restore de la configuration complète en une seule commande CLI
NFR17 : Fonctionnement offline continu illimité — aucune fonctionnalité critique ne dépend d'internet

**Scalabilité**

NFR18 : Support de 1 à 50 équipements surveillés par instance sans dégradation > 10%
NFR19 : Support de 1 à 20 utilisateurs concurrents par instance
NFR20 : Rotation automatique des données historiques — rétention par défaut 12 mois
NFR21 : Architecture plugin extensible — temps d'intégration d'un nouveau plugin < 4 heures

**Accessibilité**

NFR22 : Interface conforme WCAG 2.1 niveau AA minimum
NFR23 : Interface responsive et touch-first — utilisable sur tablette (≥ 768px), cibles tactiles 44×44 px
NFR24 : Architecture i18n dès V1 — français et anglais livrés, ajout d'une langue en < 2 heures

**Intégration**

NFR25 : Protocoles industriels MQTT 3.1.1+ et OPC-UA supportés nativement avec configuration guidée
NFR26 : Export de données en CSV et JSON depuis tout dashboard en < 10 secondes pour 100 000 lignes
NFR27 : API REST documentée (OpenAPI 3.0) exposant les endpoints de lecture

### Additional Requirements

**Architecture & Infrastructure**

- Runtime edge déployé sur site industriel via Docker + Docker Compose
- Hardware cible : Raspberry Pi 5 (8GB) minimum, Intel NUC Gen 12+ standard
- Système hôte : Linux uniquement (Ubuntu 22.04 LTS, Debian 12)
- Mise à jour en un clic depuis le dashboard avec rollback
- Instance unique par site — une installation Docker = un réseau OT
- Authentification locale en V1 (username/password) — SSO/LDAP en V2
- Pas d'exposition internet par défaut

**Modèle de Permissions (RBAC)**

- Viewer (Karim) : Lecture seule — dashboards, alertes, historiques
- Manager (M. Benali) : Lecture seule — vue synthétique ROI
- Engineer (Sofiane) : Configuration complète — pipelines, modèles, règles, plugins
- Admin : Gestion utilisateurs, mise à jour système, export

**Modèle de Données Conceptuel**

8 entités principales : Équipement, Capteur, Pipeline, Alerte, ModèleML, Utilisateur, AuditLog, Plugin

### UX Design Requirements

**Design Tokens & Système Visuel**

UX-DR01 : Implémenter le système de couleurs basé sur le logo ML_Elec (orange #F28A00, anthracite #2F3B46, gris acier #5F6B77)
UX-DR02 : Définir les design tokens (couleurs, typographie IBM Plex Sans/Mono, spacing 8px-based)
UX-DR03 : Appliquer les règles d'accessibilité WCAG 2.1 AA (contraste 4.5:1, focus visible 2px, targets tactiles 44x44px)

**Composants Custom Métier**

UX-DR04 : HealthTimelineCard — afficher état santé machine avec trajectoire et risque actionnable (courbe + confiance + horizon + CTA)
UX-DR05 : MachinePriorityList — trier et prioriser les machines par risque immédiat avec lecture en ≤ 3 secondes
UX-DR06 : AlertDecisionCard — transformer alerte en décision explicite avec action recommandée
UX-DR07 : ConfidenceBadge — rendre la certitude ML lisible et comparable (high/medium/low/unknown)
UX-DR08 : RecommendedActionBar — maintenir l'action prioritaire visible en permanence (sticky)
UX-DR09 : FocusDecisionPanel (D8) — réduire le bruit pour décision critique immédiate (mode fullscreen)
UX-DR10 : InterventionRunbook (D6) — guider l'exécution terrain pas-à-pas avec étapes validables
UX-DR11 : OfflineStatusBar — informer discrètement de l'état réseau/LLM/sync sans anxiété
UX-DR12 : SyncQueueIndicator — montrer les actions en file locale et leur statut de synchronisation
UX-DR13 : FalseAlertFeedbackControl — capturer feedback "utile / fausse alerte" en 1 tap
UX-DR14 : PipelineNodePalette — permettre construction no-code des pipelines (drag & drop React Flow)
UX-DR15 : PipelineValidationPanel — valider pipeline avant déploiement avec simulation et checks
UX-DR16 : ShareActionSheet — partager rapidement la décision terrain (WhatsApp/email/copie) avec file offline

**Patterns d'Interaction**

UX-DR17 : Implémenter la hiérarchie des boutons (Primary/Secondary/Tertiary/Destructive) avec une seule action primaire par vue
UX-DR18 : Feedback universel structuré (Signal → Contexte → Prochaine action) pour tous les états système
UX-DR19 : Navigation dashboard-first par rôle RBAC (Karim → liste machines, Sofiane → éditeur pipelines)
UX-DR20 : Progressive disclosure — complexité cachée jusqu'à ce qu'on en ait besoin

**Responsive & Offline**

UX-DR21 : Stratégie tablet-first (768px-1023px) comme cible principale V1, desktop second, mobile compagnon
UX-DR22 : Fonctionnement offline transparent — mêmes parcours mentaux online/offline, sync silencieuse
UX-DR23 : Dégradation gracieuse LLM — afficher clairement "LLM indisponible, fonctions core actives"

**Accessibilité**

UX-DR24 : Navigation clavier complète sur les parcours critiques (J1, J2, J6, J3)
UX-DR25 : Support screen readers (NVDA, VoiceOver) avec annonces d'état via aria-live
UX-DR26 : Support prefers-reduced-motion pour les animations

### FR Coverage Map

| Epic | FRs Couvertes | NFRs Couvert | UX-DRs Couvertes |
|------|---------------|--------------|------------------|
| Epic 1 : Onboarding & Déploiement Edge | FR38, FR39, FR40, FR41, FR42 | NFR01, NFR15, NFR17 | - |
| Epic 2 : Acquisition de Données & Équipements | FR01, FR02, FR03, FR04, FR05, FR27, FR44 | NFR03, NFR04, NFR25 | - |
| Epic 3 : Détection ML & Prédiction | FR06, FR46, FR07, FR08, FR09, FR10, FR11, FR50, FR45 | NFR05, NFR06 | UX-DR07 |
| Epic 4 : Explication LLM & Q&A | FR12, FR13, FR14, FR15 | NFR14 | UX-DR06, UX-DR23 |
| Epic 5 : Éditeur No-Code Pipelines | FR16, FR17, FR18, FR19, FR47, FR20 | NFR21 | UX-DR14, UX-DR15 |
| Epic 6 : Dashboards & Visualisation | FR21, FR22, FR23, FR24, FR25, FR26, FR48, FR49 | NFR02, NFR26, NFR27 | UX-DR01, UX-DR02, UX-DR04, UX-DR05, UX-DR08, UX-DR09, UX-DR11, UX-DR12, UX-DR16, UX-DR17, UX-DR18, UX-DR19, UX-DR20, UX-DR21 |
| Epic 7 : Notifications & Alertes | FR28, FR28b, FR51 | NFR03 | UX-DR06, UX-DR16 |
| Epic 8 : Administration & RBAC | FR29, FR30, FR31, FR32, FR33, FR34, FR35, FR36, FR37 | NFR07, NFR08, NFR09, NFR10, NFR11, NFR12, NFR16, NFR20 | UX-DR18 |
| Epic 9 : NFRs Transverses (i18n & Accessibilité) | FR43 | NFR22, NFR23, NFR24 | UX-DR03, UX-DR24, UX-DR25, UX-DR26 |

**📌 Note :** Les NFRs transverses (NFR22-NFR24) et FR43 (i18n) s'appliquent à TOUS les epics. Chaque story doit inclure ces exigences dans sa Definition of Done.

## Epic List

1. **Epic 1 : Onboarding & Déploiement Edge** — Permettre déploiement Docker unique et onboarding ≤ 15 min
2. **Epic 2 : Acquisition de Données & Équipements** — Connecter sources MQTT/OPC-UA, importer CSV, gérer équipements
3. **Epic 3 : Détection ML & Prédiction** — Détecter anomalies, calculer RUL, gérer feedback et confiance
4. **Epic 4 : Explication LLM & Q&A** — Générer explications langage naturel et permettre Q&A conversationnel
5. **Epic 5 : Éditeur No-Code Pipelines** — Créer, modifier, tester et déployer pipelines par drag & drop
6. **Epic 6 : Dashboards & Visualisation** — Afficher courbes de santé, alertes, ROI et permettre export/partage
7. **Epic 7 : Notifications & Alertes** — Envoyer notifications push/email/SMS et permettre partage décisionnel
8. **Epic 8 : Administration & RBAC** — Gérer utilisateurs, rôles, audit log, mises à jour et santé système
9. **Epic 9 : NFRs Transverses (i18n & Accessibilité)** — WCAG 2.1 AA, responsive tablette, support FR/EN

**🎯 Priorisation MVP :**
- **V1 Critique** : Epics 1, 2, 3, 6, 8 (5 epics pour Journey 1 complet)
- **V1.x Post-MVP** : Epics 4, 5, 7, 9 (différenciateurs et config avancée)

---

## Epic 1 : Onboarding & Déploiement Edge

**Goal :** Permettre à un ingénieur électrotechnique sans expérience Docker de déployer ML_Elec et visualiser un premier signal actif en ≤ 15 minutes, avec fonctionnement offline total et support air-gap.

**FRs Covered :** FR38, FR39, FR40, FR41, FR42
**NFRs Covered :** NFR01, NFR15, NFR17

### Story 1.1 : Commande Docker Unique de Déploiement

As a **ingénieur électrotechnique**,
I want **déployer ML_Elec via une commande Docker unique**,
So that **je peux commencer l'onboarding sans configuration complexe**.

**Acceptance Criteria:**

**Given** une machine avec Docker pré-installé (Ubuntu 22.04+)
**When** j'exécute `docker compose up -d`
**Then** tous les services ML_Elec démarrent (ingestion, ML, API, dashboard)
**And** le dashboard est accessible sur http://localhost:3000

---

### Story 1.2 : Wizard d'Onboarding en 5 Étapes

As a **nouvel utilisateur**,
I want **un guide d'onboarding interactif en 5 étapes maximum**,
So that **je sais exactement quoi faire pour atteindre la première valeur**.

**Acceptance Criteria:**

**Given** le premier accès au dashboard après déploiement
**When** le wizard s'ouvre
**Then** il affiche 5 étapes maximum avec indicateur de progression
**And** il propose des données simulées si aucun équipement réel n'est connecté
**And** le taux de complétion cible est ≥ 90%

---

### Story 1.3 : Connexion Premier Équipement en ≤ 15 Min

As a **ingénieur électrotechnique**,
I want **connecter un premier équipement et visualiser un signal actif rapidement**,
So that **je valide la valeur produit dans le temps imparti**.

**Acceptance Criteria:**

**Given** le wizard d'onboarding complété
**When** je connecte une source de données (MQTT/OPC-UA ou CSV)
**Then** le premier signal actif est visible sur le dashboard en ≤ 15 minutes total
**And** aucun code n'est requis pour cette connexion

---

### Story 1.4 : Déploiement Air-Gap (Offline Total)

As a **ingénieur en environnement isolé**,
I want **déployer ML_Elec sans aucune connexion internet**,
So that **le produit fonctionne dans des sites industriels isolés**.

**Acceptance Criteria:**

**Given** un environnement sans accès internet (air-gap)
**When** j'installe ML_Elec depuis un registre Docker local ou fichier tar
**Then** toutes les fonctions core (ingestion, ML, alertes, dashboards) fonctionnent
**And** seul le module LLM est désactivé avec indication claire

---

### Story 1.5 : Reprise Automatique Après Redémarrage

As a **utilisateur en environnement industriel**,
I want **que le système redémarre automatiquement après coupure**,
So that **je n'ai pas à intervenir manuellement**.

**Acceptance Criteria:**

**Given** une coupure électrique brutale ou redémarrage forcé
**When** le hardware redémarre
**Then** ML_Elec reprend automatiquement en < 60 secondes
**And** aucune donnée n'est perdue (write-ahead logging)
**And** la configuration est restaurée automatiquement

---

## Epic 2 : Acquisition de Données & Équipements

**Goal :** Permettre la connexion de sources de données industrielles (MQTT, OPC-UA), l'import de données historiques CSV, et la gestion des équipements surveillés.

**FRs Covered :** FR01, FR02, FR03, FR04, FR05, FR27, FR44
**NFRs Covered :** NFR03, NFR04, NFR25

### Story 2.1 : Connexion Source MQTT

As a **ingénieur**,
I want **connecter une source de données MQTT à un pipeline**,
So that **je peux ingérer des données capteurs en temps réel**.

**Acceptance Criteria:**

**Given** un broker MQTT disponible (ex: Mosquitto)
**When** je configure la connexion (broker, topic, fréquence) via l'interface
**Then** les données sont ingérées en continu sans perte
**And** la connexion est validée visuellement avant déploiement

---

### Story 2.2 : Connexion Source OPC-UA

As a **ingénieur**,
I want **connecter une source de données OPC-UA**,
So that **je peux me connecter aux automates industriels standards**.

**Acceptance Criteria:**

**Given** un serveur OPC-UA disponible (ex: open62541)
**When** je configure la connexion (endpoint, credentials, nodes) via l'interface
**Then** les données sont ingérées avec sécurité TLS 1.2+
**And** la configuration est guidée avec validation de connexion

---

### Story 2.3 : Import CSV Données Historiques

As a **ingénieur**,
I want **importer un fichier CSV de données historiques**,
So that **je peux valider le modèle ML avant connexion au matériel réel**.

**Acceptance Criteria:**

**Given** un fichier CSV avec timestamps et valeurs capteurs
**When** j'importe le fichier via l'interface
**Then** les données sont chargées et visualisables immédiatement
**And** je peux exécuter le modèle ML sur ces données historiques

---

### Story 2.4 : Détection Signaux Corrompus

As a **ingénieur**,
I want **que le système détecte et signale les données manquantes ou corrompues**,
So that **je peux agir rapidement sur les problèmes capteurs**.

**Acceptance Criteria:**

**Given** des données entrantes avec gaps, hors plage ou corrompues
**When** le système ingère les données
**Then** les signaux problématiques sont signalés dans l'interface
**And** la qualité des données est visible sur chaque vue machine

---

### Story 2.5 : Ingestion Continue avec Reprise Auto

As a **utilisateur**,
I want **une ingestion de données 24h/7j avec reprise automatique**,
So that **je ne perds aucune donnée même après interruption**.

**Acceptance Criteria:**

**Given** une interruption réseau ou système
**When** la connexion reprend
**Then** l'ingestion reprend automatiquement sans perte de données
**And** le système supporte ≥ 1 000 points/seconde par équipement

---

### Story 2.6 : Créer, Nommer et Archiver des Équipements

As a **ingénieur**,
I want **créer, nommer et archiver des équipements surveillés**,
So that **je peux organiser mon parc de machines**.

**Acceptance Criteria:**

**Given** la page de gestion des équipements
**When** je crée un nouvel équipement
**Then** je peux le nommer, lui assigner des capteurs, le catégoriser
**And** je peux l'archiver quand il n'est plus surveillé
**And** la modification est tracée dans l'audit log

---

## Epic 3 : Détection ML & Prédiction

**Goal :** Détecter automatiquement les anomalies via le modèle généraliste, calculer la RUL, et fournir un système de feedback utilisateur pour améliorer la confiance.

**FRs Covered :** FR06, FR46, FR07, FR08, FR09, FR10, FR11, FR50, FR45
**NFRs Covered :** NFR05, NFR06

### Story 3.1 : Modèle ML Généraliste Pré-entraîné

As a **nouvel utilisateur**,
I want **un modèle ML généraliste fonctionnel dès l'installation**,
So that **je n'ai pas besoin d'entraînement préalable pour démarrer**.

**Acceptance Criteria:**

**Given** une installation fraîche de ML_Elec
**When** les données sont ingérées
**Then** le modèle généraliste détecte les anomalies sans configuration
**And** il atteint ≥ 70% de détection avec < 20% faux positifs sur datasets de référence

---

### Story 3.2 : Détection d'Anomalies en Temps Réel

As a **système**,
I want **détecter automatiquement les anomalies sur les signaux**,
So that **les utilisateurs sont alertés des dérives équipements**.

**Acceptance Criteria:**

**Given** des données capteurs en temps réel
**When** une anomalie est détectée (Isolation Forest/Autoencoder)
**Then** une alerte est générée avec niveau de sévérité (info/warning/critique)
**And** l'inférence s'exécute en < 500 ms par cycle

---

### Story 3.3 : Calcul RUL avec Horizon de Défaillance

As a **technicien**,
I want **connaître la durée de vie résiduelle (RUL) d'un équipement**,
So that **je peux planifier l'intervention au bon moment**.

**Acceptance Criteria:**

**Given** un équipement avec données historiques suffisantes
**When** le modèle RUL est exécuté
**Then** l'horizon de défaillance estimé est affiché (ex: "4-6 jours")
**And** la tendance de dégradation est visible sur la courbe de santé

---

### Story 3.4 : Configuration Seuils et Règles de Corrélation

As a **ingénieur**,
I want **configurer des seuils d'alerte et règles de corrélation**,
So that **je peux réduire les faux positifs**.

**Acceptance Criteria:**

**Given** des alertes avec faux positifs récurrents
**When** je configure des seuils ou règles de corrélation via l'éditeur
**Then** les fausses alertes sont filtrées
**And** la modification est tracée dans l'audit log

---

### Story 3.5 : Feedback Utilisateur 1 Tap

As a **technicien**,
I want **soumettre un feedback sur une alerte en 1 tap**,
So that **je peux signaler les fausses alertes facilement**.

**Acceptance Criteria:**

**Given** une alerte affichée dans l'interface
**When** je tape "Utile" ou "Fausse alerte"
**Then** le feedback est enregistré immédiatement avec confirmation
**And** je peux ajouter une note de contexte optionnelle

---

### Story 3.6 : Indicateur de Confiance ML

As a **utilisateur**,
I want **voir le niveau de confiance du modèle sur chaque alerte**,
So que **je peux juger de la crédibilité de l'alerte**.

**Acceptance Criteria:**

**Given** une alerte générée
**When** l'alerte est affichée
**Then** le ConfidenceBadge montre le niveau (high/medium/low/unknown) avec score %
**And** un tooltip explique les facteurs de confiance

---

### Story 3.7 : Signalisation Contexte Ambigu

As a **technicien**,
I want **être averti des contextes ambigus avant d'agir**,
So that **je ne fais pas d'intervention inutile**.

**Acceptance Criteria:**

**Given** une alerte dans un contexte ambigu (chaleur ambiante élevée, modèle en apprentissage, données manquantes)
**When** l'alerte est affichée
**Then** un indicateur de contexte ambigu est visible avant l'action
**And** des facteurs explicites sont listés (ex: "Température ambiante > 35°C")

---

### Story 3.8 : Suggestion d'Action Concrète

As a **technicien**,
I want **une suggestion d'action concrète pour chaque alerte**,
So that **je sais exactement quoi faire maintenant**.

**Acceptance Criteria:**

**Given** une alerte détectée
**When** l'alerte est affichée
**Then** une action est recommandée textuellement (ex: "Planifier intervention jeudi", "Surveiller 48h", "Commander la pièce")
**And** la suggestion est copyable pour partage WhatsApp/email

---

## Epic 4 : Explication LLM & Q&A

**Goal :** Générer des explications en langage naturel compréhensibles et permettre le Q&A conversationnel sur l'état des équipements, avec possibilité de désactiver le LLM.

**FRs Covered :** FR12, FR13, FR14, FR15
**NFRs Covered :** NFR14

### Story 4.1 : Explication LLM en Langage Naturel

As a **technicien**,
I want **une explication en langage naturel pour chaque anomalie**,
So that **je comprends la cause sans expertise ML**.

**Acceptance Criteria:**

**Given** une anomalie détectée
**When** l'alerte est générée
**Then** une explication est affichée en français technique (ex: "Vibration anormale sur roulement gauche — dégradation progressive, probabilité de défaillance dans 4 à 6 jours")
**And** 80% des techniciens BTS reformulent correctement l'action après lecture

---

### Story 4.2 : Q&A Conversationnel sur l'État Équipement

As a **technicien**,
I want **poser une question en langage naturel sur un équipement**,
So que **j'obtiens une réponse rapide sans analyser les courbes**.

**Acceptance Criteria:**

**Given** un équipement sélectionné
**When** je pose une question (ex: "Pourquoi le score de M-03 baisse ?")
**Then** une réponse basée sur les données est fournie en < 10s
**And** le système indique clairement les limites de ses capacités

---

### Story 4.3 : Activation/Désactivation Module LLM

As a **administrateur**,
I want **activer ou désactiver le module LLM**,
So que **je peux contrôler la dépendance aux API externes**.

**Acceptance Criteria:**

**Given** le panneau d'administration
**When** je bascule le switch LLM
**Then** le module est activé/désactivé immédiatement
**And** toutes les autres fonctions (ML, alertes, dashboards) restent opérationnelles
**And** l'état est visible dans OfflineStatusBar

---

### Story 4.4 : Configuration Fournisseur LLM et Clé API

As a **administrateur**,
I want **configurer le fournisseur LLM et la clé API**,
So que **je peux choisir entre OpenAI, Anthropic, Mistral, etc.**.

**Acceptance Criteria:**

**Given** le module LLM activé
**When** je configure un fournisseur et une clé API
**Then** la configuration est testée et validée
**And** la clé est stockée chiffrée (AES-256)

---

## Epic 5 : Éditeur No-Code Pipelines

**Goal :** Permettre la création, modification, test et déploiement de pipelines de traitement par drag & drop sans code, avec support de 6 types de plugins.

**FRs Covered :** FR16, FR17, FR18, FR19, FR47, FR20
**NFRs Covered :** NFR21

### Story 5.1 : Création Pipeline par Drag & Drop

As a **ingénieur**,
I want **créer un pipeline de traitement par drag & drop**,
So que **je peux configurer le traitement sans écrire de code**.

**Acceptance Criteria:**

**Given** l'éditeur de pipelines ouvert
**When** je glisse-dépose des nœuds (Source → Filtre → Anomalie → RUL → Sortie)
**Then** les connexions se font avec snap intelligent
**And** les nœuds sont colorés par type (bleu=violet=orange=vert)
**And** les labels utilisent le vocabulaire électrotechnique

---

### Story 5.2 : Modification Pipeline sans Redéploiement

As a **ingénieur**,
I want **modifier un pipeline existant sans redéploiement**,
So que **je peux ajuster la configuration rapidement**.

**Acceptance Criteria:**

**Given** un pipeline déployé en production
**When** je modifie un nœud ou une connexion
**Then** la modification est appliquée sans interruption du service
**And** la modification est tracée dans l'audit log

---

### Story 5.3 : Installation Plugin en 1 Clic

As a **ingénieur**,
I want **installer un plugin depuis la bibliothèque**,
So que **je peux étendre les capacités du système**.

**Acceptance Criteria:**

**Given** la bibliothèque de plugins ouverte
**When** je clique "Installer" sur un plugin
**Then** le plugin est installé et disponible dans la palette en 1 clic
**And** la licence (open-source/propriétaire) est clairement indiquée

---

### Story 5.4 : Architecture 6 Types de Plugins

As a **développeur**,
I want **6 types de plugins avec contrats documentés**,
So que **je peux créer des extensions compatibles**.

**Acceptance Criteria:**

**Given** la documentation développeur
**When** je développe un plugin
**Then** il suit le contrat d'un des 6 types : (1) connecteurs protocole, (2) transformateurs, (3) modèles ML, (4) blocs visualisation, (5) canaux notification, (6) exportateurs
**And** un développeur familier peut intégrer un nouveau plugin en < 4 heures

---

### Story 5.5 : Test Pipeline sur Données Historiques

As a **ingénieur**,
I want **tester un pipeline en simulation avant déploiement**,
So que **je peux valider le fonctionnement sans risque**.

**Acceptance Criteria:**

**Given** un pipeline construit
**When** je lance le mode test sur données historiques
**Then** chaque nœud affiche son statut (OK/erreur/données traitées) en temps réel
**And** un panneau de validation affiche les checks (latence, qualité données, cohérence) avec verdict
**And** le déploiement est bloqué si fail critique

---

## Epic 6 : Dashboards & Visualisation

**Goal :** Fournir des dashboards adaptés à chaque rôle (technicien, manager, ingénieur) avec courbes de santé, alertes, vue ROI, et permettre export/partage de données.

**FRs Covered :** FR21, FR22, FR23, FR24, FR25, FR26, FR48, FR49
**NFRs Covered :** NFR02, NFR26, NFR27

### Story 6.1 : Courbe de Santé Temporelle par Équipement

As a **technicien**,
I want **visualiser la courbe de santé temporelle de chaque équipement**,
So que **je vois l'évolution et la trajectoire, pas juste une valeur instantanée**.

**Acceptance Criteria:**

**Given** un équipement sélectionné
**When** la fiche machine s'ouvre
**Then** la courbe de santé avec historique (4+ jours) est affichée
**And** la tendance, le score actuel et l'horizon de défaillance sont visibles
**And** le chargement se fait en < 2 secondes

---

### Story 6.2 : Historique des Alertes et Statuts

As a **technicien**,
I want **consulter l'historique des alertes et leur statut**,
So que **je peux suivre les événements passés**.

**Acceptance Criteria:**

**Given** l'historique alertes ouvert
**When** je consulte la liste
**Then** chaque alerte affiche : date, sévérité, machine, explication, statut (nouvelle/reconnue/planifiée/résolue/fausse)
**And** je peux filtrer par machine, sévérité, période

---

### Story 6.3 : Vue Synthétique ROI pour Manager

As a **manager (M. Benali)**,
I want **accéder à une vue synthétique ROI**,
So que **je vois les pannes évitées et économies sans détails techniques**.

**Acceptance Criteria:**

**Given** la vue ROI ouverte
**When** le manager consulte l'écran
**Then** les KPIs suivants sont affichés en gros : pannes évitées, économie estimée (€), disponibilité machines (%)
**And** un export PDF est disponible pour présentation en réunion
**And** la lecture se fait en ≤ 10 secondes

---

### Story 6.4 : Vue Technique pour Ingénieur

As a **ingénieur**,
I want **accéder à une vue technique complète**,
So que **je peux diagnostiquer les problèmes système**.

**Acceptance Criteria:**

**Given** la vue technique ouverte
**When** je consulte l'écran
**Then** je vois : liste des pipelines et statuts, modèles ML actifs avec confiance, 1000 derniers logs avec filtrage
**And** je peux filtrer par composant et sévérité

---

### Story 6.5 : Export Données CSV/JSON

As a **utilisateur**,
I want **exporter des données en CSV ou JSON**,
So que **je peux les analyser dans d'autres outils**.

**Acceptance Criteria:**

**Given** un dashboard ou une liste d'alertes
**When** je clique "Exporter"
**Then** le fichier CSV/JSON est généré en < 10 secondes pour 100 000 lignes
**And** je peux choisir les colonnes/champs à exporter

---

### Story 6.6 : Partage Lien Lecture Seule

As a **ingénieur**,
I want **partager un lien de vue en lecture seule**,
So que **je peux montrer des données sans donner l'accès complet**.

**Acceptance Criteria:**

**Given** une vue dashboard ouverte
**When** je clique "Partager"
**Then** un lien lecture seule est généré avec expiration configurable
**And** le lien peut être envoyé par email/WhatsApp

---

### Story 6.7 : Bandeau d'État Offline/Online

As a **utilisateur**,
I want **voir l'état de connexion réseau discrètement**,
So que **je sais si je suis offline sans anxiété**.

**Acceptance Criteria:**

**Given** n'importe quelle page de l'application
**When** l'état réseau change
**Then** OfflineStatusBar affiche : online/offline/resynchronisation
**And** le statut LLM (actif/indisponible) est indiqué
**And** la barre est discrète et non anxiogène

---

### Story 6.8 : Fraîcheur des Données Visible

As a **technicien**,
I want **voir la fraîcheur des données (timestamp dernière mesure)**,
So que **je peux juger de la fiabilité de l'information**.

**Acceptance Criteria:**

**Given** une vue machine ou une alerte
**When** les données sont affichées
**Then** le timestamp de la dernière mesure fiable est visible
**And** un indicateur "données obsolètes" s'affiche si > seuil configurable

---

## Epic 7 : Notifications & Alertes

**Goal :** Envoyer des notifications (push, email, SMS) lors d'alertes et permettre le partage de résumés décisionnels via WhatsApp/email avec file offline.

**FRs Covered :** FR28, FR28b, FR51
**NFRs Covered :** NFR03

### Story 7.1 : Réception Notifications Push/Email/SMS

As a **technicien**,
I want **recevoir une notification lors d'une alerte**,
So that **je suis alerté immédiatement même hors de l'application**.

**Acceptance Criteria:**

**Given** une alerte critique générée
**When** la notification est envoyée
**Then** je reçois push/email/SMS selon configuration
**And** la notification inclut : machine, sévérité, horizon, lien vers fiche
**And** la latence notification < 5 secondes

---

### Story 7.2 : Configuration Canaux de Notification par Machine

As a **technicien**,
I want **configurer les canaux de notification par machine**,
So that **je reçois uniquement les alertes pertinentes**.

**Acceptance Criteria:**

**Given** la page de configuration d'une machine
**When** je configure les notifications
**Then** je peux choisir : canaux (push/email/SMS), niveaux de sévérité déclencheurs, plages horaires
**And** la configuration est spécifique à chaque machine

---

### Story 7.3 : Partage Résumé Décisionnel avec File Offline

As a **technicien**,
I want **partager un résumé décisionnel via WhatsApp/email/copie**,
So that **je peux coordonner avec mon équipe même offline**.

**Acceptance Criteria:**

**Given** une décision prise (planifier/surveiller)
**When** je clique "Partager"
**Then** un résumé est généré (machine, risque, action, créneau)
**And** je peux partager via WhatsApp/email ou copier le texte
**And** en offline, l'envoi est mis en file locale puis synchronisé

---

## Epic 8 : Administration & RBAC

**Goal :** Gérer les utilisateurs, les rôles, l'audit log, les mises à jour système et la santé globale de ML_Elec.

**FRs Covered :** FR29, FR30, FR31, FR32, FR33, FR34, FR35, FR36, FR37
**NFRs Covered :** NFR07, NFR08, NFR09, NFR10, NFR11, NFR12, NFR16, NFR20

### Story 8.1 : Gestion Utilisateurs et Rôles RBAC

As a **administrateur**,
I want **créer, modifier et supprimer des comptes utilisateurs avec rôles**,
So that **je peux contrôler les accès**.

**Acceptance Criteria:**

**Given** le panneau d'administration
**When** je crée un utilisateur
**Then** je peux assigner un rôle (Viewer/Manager/Engineer/Admin)
**And** l'interface s'adapte au rôle dès la connexion (FR30)
**And** RBAC est appliqué sur 100% des endpoints API

---

### Story 8.2 : Audit Log Inaltérable des Modifications

As a **administrateur**,
I want **un audit log de toute modification de configuration**,
So que **je peux tracer qui a fait quoi et quand**.

**Acceptance Criteria:**

**Given** une modification de configuration (seuil, règle, pipeline, utilisateur)
**When** la modification est effectuée
**Then** l'audit log enregistre : auteur, horodatage, valeur avant/après
**And** le log est inaltérable (append-only, non supprimable même par Admin)
**And** la rétention est ≥ 12 mois

---

### Story 8.3 : Mise à Jour Logiciel en 1 Clic

As a **administrateur**,
I want **déclencher une mise à jour depuis le dashboard**,
So that **je peux maintenir le système à jour facilement**.

**Acceptance Criteria:**

**Given** une nouvelle version disponible
**When** je clique "Mettre à jour"
**Then** la nouvelle image Docker est téléchargée et déployée
**And** le rollback est possible vers la version précédente
**And** je suis notifié des nouvelles versions disponibles

---

### Story 8.4 : Consultation État de Santé Système

As a **administrateur**,
I want **consulter l'état de santé du système**,
So that **je peux diagnostiquer les problèmes**.

**Acceptance Criteria:**

**Given** la page santé système ouverte
**When** je consulte l'écran
**Then** je vois : services actifs, espace disque, connexions protocoles, uptime
**And** des alertes système sont générées si seuils critiques atteints

---

### Story 8.5 : Backup & Restore Configuration

As a **administrateur**,
I want **backup et restore de la configuration en une commande**,
So that **je peux récupérer rapidement après incident**.

**Acceptance Criteria:**

**Given** la commande CLI de backup
**When** j'exécute `mlelec backup`
**Then** la configuration complète (pipelines, utilisateurs, seuils, modèles) est sauvegardée
**And** `mlelec restore` restaure l'instance avec validation fonctionnelle

---

## Epic 9 : NFRs Transverses (i18n & Accessibilité)

**Goal :** Appliquer les exigences transverses d'internationalisation (FR/EN) et d'accessibilité (WCAG 2.1 AA) à tous les epics.

**FRs Covered :** FR43
**NFRs Covered :** NFR22, NFR23, NFR24

**📌 Note :** Ces exigences s'appliquent à TOUS les epics. Chaque story doit inclure :
- **i18n (NFR24, FR43)** : Interface traduisible FR/EN, ajout langue < 2h
- **Accessibilité (NFR22)** : WCAG 2.1 AA, contraste 4.5:1, navigation clavier
- **Responsive (NFR23)** : Touch-first, cibles 44×44 px, tablette ≥ 768px

### Story 9.1 : Sélection Langue Interface (FR/EN)

As a **utilisateur**,
I want **sélectionner la langue de l'interface**,
So that **je peux utiliser ML_Elec dans ma langue**.

**Acceptance Criteria:**

**Given** les paramètres utilisateur
**When** je sélectionne une langue (FR ou EN)
**Then** toute l'interface est traduite immédiatement
**And** l'ajout d'une nouvelle langue est possible en < 2 heures par un traducteur non-technique

---

### Story 9.2 : Conformité WCAG 2.1 Niveau AA

As a **utilisateur avec besoins d'accessibilité**,
I want **une interface conforme WCAG 2.1 AA**,
So that **je peux utiliser ML_Elec sans barrières**.

**Acceptance Criteria:**

**Given** l'application ML_Elec
**When** un audit automatisé (axe-core/Lighthouse) est exécuté
**Then** le score accessibilité est ≥ 90%
**And** le contraste est ≥ 4.5:1 (texte), ≥ 3:1 (UI large)
**And** la navigation clavier est complète sur les parcours critiques (J1, J2, J6, J3)

---

### Story 9.3 : Responsive Tablet-First

As a **technicien sur tablette**,
I want **une interface optimisée pour tablette (≥ 768px)**,
So that **je peux l'utiliser sur le terrain facilement**.

**Acceptance Criteria:**

**Given** une tablette (iPad/Samsung Tab)
**When** j'utilise ML_Elec
**Then** toutes les actions sont accessibles avec cibles tactiles ≥ 44×44 px
**And** la stratégie tablet-first est appliquée (768px-1023px cible V1)
**And** le layout s'adapte desktop (12 colonnes) et tablette (8 colonnes)

---

### Story 9.4 : Support Screen Readers

As a **utilisateur avec screen reader**,
I want **une compatibilité NVDA/VoiceOver**,
So that **je peux naviguer avec un lecteur d'écran**.

**Acceptance Criteria:**

**Given** un screen reader activé (NVDA/VoiceOver)
**When** je navigue dans l'application
**Then** la structure HTML est sémantique avec ARIA approprié
**And** les états critiques sont annoncés via aria-live="polite"
**And** les parcours critiques (J1, J2, J6, J3) sont navigables entièrement

---

### Story 9.5 : Support prefers-reduced-motion

As a **utilisateur sensible aux animations**,
I want **désactiver les animations**,
So that **je peux utiliser l'application sans inconfort**.

**Acceptance Criteria:**

**Given** le système avec prefers-reduced-motion activé
**When** j'utilise ML_Elec
**Then** les animations sont réduites ou désactivées
**And** l'expérience reste fonctionnelle sans motion

---

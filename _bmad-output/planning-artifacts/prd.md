---
stepsCompleted: ["step-01-init", "step-02-discovery", "step-02b-vision", "step-02c-executive-summary", "step-03-success", "step-04-journeys", "step-05-domain", "step-06-innovation", "step-07-project-type", "step-08-scoping", "step-09-functional", "step-10-nonfunctional", "step-11-polish", "step-e-01-discovery", "step-e-02-review", "step-e-03-edit"]
inputDocuments:
  - "_bmad-output/planning-artifacts/product-brief-ML_Elec-2026-03-04.md"
  - "_bmad-output/planning-artifacts/research/market-supervision-industrielle-hybride-research-2026-03-03.md"
  - "_bmad-output/planning-artifacts/research/domain-genie-electrique-electrotechnique-research-2026-03-03.md"
  - "_bmad-output/planning-artifacts/research/technical-faisabilite-ml-elec-research-2026-03-03.md"
  - "_bmad-output/brainstorming/brainstorming-session-2026-03-03-00-13-24.md"
  - "_bmad-output/planning-artifacts/prd-validation-report.md"
workflowType: 'prd'
workflow: 'edit'
briefCount: 1
researchCount: 3
brainstormingCount: 1
projectDocsCount: 0
classification:
  projectType: "IoT/Edge + SaaS B2B (hybride)"
  domain: "Process Control / Predictive Maintenance"
  complexity: "high"
  projectContext: "greenfield"
  keyConstraint: "15 min onboarding = release blocker absolu"
lastEdited: '2026-03-06'
editHistory:
  - date: '2026-03-06'
    changes: "Troisième édition — alignement UX: KPI décisionnels (3s/60s), enrichissement journeys J1/J2/J3/J6, ajout FR48-FR51 (offline/sync, fraîcheur, contexte, partage), ajout NFR28-NFR29, clarification autorité seuils critiques"
  - date: '2026-03-05'
    changes: "Seconde édition — 3 axes : consolidation FRs 44-47 dans catégories existantes (suppression sections dupliquées 11-14), ajout critères d'acceptation FR12/FR39/FR46, ajout modèle de données conceptuel (8 entités)"
  - date: '2026-03-04'
    changes: "Édition complète basée sur rapport de validation — 8 axes : reconciliation business targets avec Brief, North Star Metric ajouté, 27 NFRs enrichis (méthodes de vérification), 3 fuites d'implémentation corrigées, FR28 splitté, 4 FRs ajoutés (FR44-FR47), 2 journeys ajoutés (J6 offline, J7 Q&A), données marché ajoutées, glossaire ajouté, scope V1-V3 enrichi, domaine ISA-18.2 & MOC ajoutés"
---

# Product Requirements Document — ML_Elec

**Author:** Marwane
**Date:** 2026-03-04

### Table des matières

1. [Executive Summary](#executive-summary)
2. [Classification du Projet](#classification-du-projet)
3. [Critères de Succès](#critères-de-succès)
4. [User Journeys](#user-journeys)
5. [Domain-Specific Requirements](#domain-specific-requirements)
6. [Innovation & Patterns Novateurs](#innovation--patterns-novateurs)
7. [Exigences Spécifiques IoT/Edge + SaaS B2B](#exigences-spécifiques-iotedge--saas-b2b)
8. [Scope & Développement Phasé](#scope--développement-phasé)
9. [Exigences Fonctionnelles](#exigences-fonctionnelles)
10. [Exigences Non-Fonctionnelles](#exigences-non-fonctionnelles)

## Executive Summary

ML_Elec est une plateforme open-source de maintenance prédictive pour équipements électriques, conçue pour les PMEs industrielles (50–500 employés). Elle s'adresse en priorité aux techniciens de maintenance et ingénieurs électrotechniques qui opèrent sans data scientist ni infrastructure cloud — et qui perdent aujourd'hui des milliers d'euros en pannes non anticipées faute d'outils accessibles.

### Contexte Marché

Le segment PME industrielle (50–500 employés) représente un TAM estimé à 2–3 milliards € en Europe. 73% des responsables industriels identifient la maintenance prédictive comme priorité n°1, mais les solutions actuelles restent inaccessibles : coût d'entrée > 50k€/an, intégration par prestataire, dépendance cloud. Chaque arrêt non planifié coûte entre 1 000 € et 1 000 000 €/heure selon le secteur. ML_Elec exploite cette fenêtre de 18–24 mois avant que les acteurs établis (Siemens, Schneider) ne ferment le marché avec des solutions propriétaires.

Le produit s'inscrit dans la vague d'émergence de l'IA embarquée : les modèles ML légers et les LLMs d'inférence rapide rendent désormais possible ce qui était réservé aux grandes industries.

### Ce qui rend ML_Elec spécial

Le moment "waouh" de l'utilisateur : brancher ML_Elec sur un moteur et voir, pour la première fois, **l'évolution de la santé de sa machine dans le temps** — une mémoire vivante et lisible, sans expertise data requise.

Trois différenciateurs combinés qu'aucun concurrent n'assemble aujourd'hui :
1. **Offline-first** — fonctionne sans cloud, sans dépendance réseau, en souveraineté totale
2. **Double couche IA** — ML prédictif (Isolation Forest / Autoencoder / RUL) + LLM explicatif en langage naturel
3. **No-code accessible** — éditeur visuel de pipelines, onboarding en 15 minutes chrono (cf. NFR01)

L'insight de fond : les PMEs ne manquent pas de données — elles manquent d'un outil qui traduit ces données en décisions compréhensibles pour les personnes sur le terrain, sans intermédiaire.

## Classification du Projet

| Dimension | Valeur |
|---|---|
| **Type** | IoT/Edge + SaaS B2B (hybride) |
| **Domaine** | Process Control / Predictive Maintenance |
| **Complexité** | HIGH — intégration multi-couches (acquisition → ML → LLM → visualisation → plugins) |
| **Contexte** | Greenfield |
| **Contrainte clé** | Onboarding ≤ 15 min = release blocker absolu (cf. NFR01) |

## Critères de Succès

### North Star Metric

**Nombre de moteurs surveillés en production active**
_(moteurs avec au moins une alerte générée dans les 30 derniers jours)_

Ce chiffre capture la valeur créée : des machines réellement surveillées, des alertes réelles générées — pas juste un logiciel installé.

### Succès Utilisateur

| Critère | Cible | Mesure |
|---|---|---|
| Onboarding complet (premier signal actif visible) | ≤ 15 minutes (cf. NFR01) | Timestamp premier signal − timestamp déploiement |
| Repérage machine prioritaire (glanceability UX) | ≤ 3 secondes (cf. NFR28) | Timestamp affichage dashboard → timestamp ouverture machine prioritaire |
| Décision actionnable sur alerte (boucle UX) | ≤ 60 secondes (cf. NFR29) | Timestamp ouverture alerte → timestamp action choisie/partagée |
| Time-to-First-Alert | < 24h post-installation (idéal < 1h) | Timestamp première alerte − timestamp installation |
| Taux de complétion onboarding sans aide externe | ≥ 90% | Test utilisateur (5 techniciens terrain) |
| Explication LLM compréhensible par technicien non-expert | Taux compréhension ≥ 80% | Test utilisateur (techniciens BTS sans formation ML) |
| Déploiement nouvelle machine par ingénieur seul | ≤ 30 minutes | Test utilisateur chronométré |
| Ratio Préventif/Réactif | > 50% alertes → intervention avant panne | (Pannes évitées / Pannes totales) × 100 |
| Rétention 30 jours | > 60% des déploiements actifs après 30j | Déploiements actifs J30 / Déploiements totaux |
| NPS utilisateur | > 40 | In-app survey opt-in |
| **Définition utilisateur actif** | ≥ 1 session toutes les 48h | Logs sessions |

**Signal d'alarme :** Fréquence d'ouverture < 3×/semaine = problème de valeur perçue. Taux de faux positifs > 30% = churn silencieux à venir.

### Succès Business

| Métrique | Cible 3 mois | Cible 12 mois |
|---|---|---|
| Déploiements pilotes actifs en production | 5 | 50 |
| Moteurs surveillés en production active (North Star) | 15 | 250 |
| Stars GitHub | 50 | 500 |
| Contributeurs externes | — | 10 |
| Case study avec ROI chiffré publiée | — | 1 |

**Note :** Ces cibles sont calibrées pour un MVP greenfield avec équipe réduite. Elles reflètent les objectifs du Product Brief.

### Succès Technique

| Critère | Cible |
|---|---|
| Latence détection anomalie (acquisition → alerte) | ≤ 5 secondes (cf. NFR03) |
| Taux de faux positifs modèle anomalie | < 20% à 3 mois, < 10% à 12 mois (cf. NFR06) |
| Disponibilité offline totale | 100% des features core (cf. NFR17) |
| Couverture protocoles V1 | MQTT + OPC-UA (cf. NFR25) |
| Taille image Docker edge | ≤ 500 MB |

**Note :** La cible faux positifs est alignée avec le Product Brief (< 20% KPI). Le seuil < 5% du NFR06 s'applique comme objectif technique sur datasets de référence, pas comme critère de succès en conditions réelles.

### Outcomes Mesurables

- Une PME déployant ML_Elec **réduit ses pannes non planifiées d'au moins 30%** sur les équipements surveillés dans les 30 jours suivant le déploiement
- Le technicien de maintenance **consulte le dashboard au moins 3 fois par semaine** (≥ 1 session toutes les 48h) sans relance externe
- Le modèle RUL génère une **alerte préventive au moins 72h avant** une défaillance détectable
- Le ratio préventif/réactif atteint **> 50%** sur les équipements surveillés dans les 90 jours suivant le déploiement

## User Journeys

### Journey 1 — Karim : La première anomalie 🟢 *(Chemin heureux)*

**Lundi matin, 7h43. Usine de textile, Tlemcen.**

Karim arrive avant les autres. Depuis 6 mois, le moteur M-07 de la ligne 3 lui donne des sueurs froides — une vibration irrégulière que personne d'autre n'entend, mais lui oui. La semaine dernière, il a encore manqué la réunion de famille parce que M-07 a lâché un vendredi soir.

Sofiane a installé ML_Elec le jeudi précédent. Karim a suivi l'onboarding en 12 minutes. Ce matin, il ouvre le dashboard pour la première fois en conditions réelles.

Il voit une courbe. Pas des chiffres bruts — **une courbe de santé**. M-07 y est représenté sur les 4 derniers jours. Et la courbe descend doucement depuis 48h.

Il clique sur l'alerte orange. Une phrase s'affiche : *"Vibration anormale détectée sur roulement gauche — dégradation progressive, probabilité de défaillance dans 4 à 6 jours si non traité. Cause probable : défaut de lubrification ou usure prématurée."*

Karim relit deux fois. Ce n'est pas du jargon — c'est exactement ce qu'il aurait dit à son chef, mais en mieux formulé. Il prend une photo de l'écran, la partage sur le groupe WhatsApp de l'équipe : *"M-07, intervention jeudi, pas vendredi soir."*

**Nouvelle réalité :** Karim a évité une panne critique, planifié une intervention préventive, et rentre chez lui à l'heure le vendredi.

**Exigences révélées :** dashboard santé visuel avec courbe temporelle, indicateur de confiance ML et fraîcheur des données visibles, alertes en langage naturel (LLM), action recommandée sélectionnable/partageable en ≤ 60s, RUL avec horizon de défaillance explicite.

---

### Journey 2 — Karim : La fausse alerte 🔴 *(Cas limite — récupération)*

**Mercredi, 14h20.**

ML_Elec génère une alerte rouge sur M-12 : *"Anomalie critique — risque de surchauffe."* Karim intervient, arrête la ligne, inspecte pendant 45 minutes. Rien. Le moteur est parfait. Karim est furieux.

Le lendemain, Sofiane analyse : M-12 était en zone de forte chaleur ambiante. Le modèle n'avait pas de contexte environnemental. Sofiane ajoute un capteur de température ambiante dans l'éditeur no-code en 8 minutes — drag & drop, pas de code. Il crée une règle : *"Ignorer alerte thermique si température ambiante > 35°C."*

Deux semaines plus tard, zéro fausse alerte. Karim reprend confiance.

**Nouvelle réalité :** L'erreur est récupérable — parce que Sofiane a pu corriger le modèle sans prestataire externe.

**Exigences révélées :** logs de contexte d'alerte, éditeur no-code de règles de corrélation, mécanisme de feedback en 1 tap (utile / fausse alerte) avec note contextuelle optionnelle, boucle de correction visible (analyse → ajustement → validation).

---

### Journey 3 — Sofiane : Déploiement nouvelle machine 🔧 *(Admin/déployeur)*

**Mardi matin. Nouvelle commande : surveiller 3 compresseurs dans l'atelier B.**

Sofiane ouvre l'éditeur visuel. Pipeline en drag & drop : source MQTT → filtre fréquence → bloc anomalie → bloc RUL → sortie dashboard. 22 minutes, zéro ligne de code.

Il bute sur un format propriétaire non standard. Il cherche dans la bibliothèque de plugins — un contributeur communautaire l'a publié 3 semaines plus tôt. Installation en un clic. Déploiement Docker sur le NUC de l'atelier B : 4 minutes. Les 3 compresseurs apparaissent dans le dashboard. Sofiane envoie un lien lecture seule à M. Benali.

**Nouvelle réalité :** Ce qui aurait nécessité 2 jours avec un prestataire SCADA est fait en une matinée — seul.

**Exigences révélées :** éditeur visuel no-code complet, validation pipeline sur données historiques avant déploiement, fallback protocole non supporté via plugin communautaire ou import CSV, déploiement Docker simplifié, partage de vue en lecture seule.

---

### Journey 4 — M. Benali : La décision de renouvellement 💼 *(Décideur budget)*

**Fin de trimestre. Revue de performance production.**

M. Benali ouvre la vue "résumé direction" configurée par Sofiane. Il voit : 3 pannes évitées, économie estimée 14 000€, disponibilité machines +4,2% vs trimestre précédent. Il n'a pas besoin de comprendre comment ça marche. Il voit le ROI. Il valide le renouvellement en 30 secondes et demande à Sofiane d'étendre la surveillance à l'atelier C.

**Nouvelle réalité :** M. Benali est passé de sceptique à sponsor interne.

**Exigences révélées :** vue synthétique direction (non-technique), calcul ROI automatique, KPIs de disponibilité machines, export rapport PDF / partage.

---

### Journey 5 — Yasmine : Première contribution open-source 🌱 *(Communauté)*

**Yasmine, étudiante M2 génie électrique, Alger.**

Elle découvre ML_Elec sur GitHub via un post Reddit. Elle clone le repo, lance Docker en 10 minutes, connecte un dataset simulé. Elle identifie un bug dans le plugin Modbus RTU, ouvre une issue, soumet une PR avec fix + test. Merge en 2 jours.

Elle revient deux semaines plus tard pour développer un plugin IEC 61850 — sujet de son mémoire. La communauté l'aide. Il sera utilisé en production par une PME tunisienne dans le mois suivant.

**Nouvelle réalité :** Yasmine a une contribution open-source réelle sur son CV. ML_Elec gagne un nouveau protocole gratuitement.

**Exigences révélées :** documentation développeur claire, architecture plugin ouverte et documentée, processus de contribution frictionless, environnement de dev local simple (Docker), tests automatisés pour les plugins.

---

### Journey 6 — Karim : Coupure réseau, zéro impact 🔌 *(Différenciateur offline-first)*

**Jeudi, 10h15. Même usine, Tlemcen. Coupure internet depuis 6h.**

Le FAI régional est en panne — comme tous les 2 mois. Les collègues de Karim ne peuvent plus accéder à leurs emails ni au cloud. Mais Karim ouvre ML_Elec sur sa tablette, connectée au réseau local de l'usine.

Tout fonctionne. Les courbes de santé se mettent à jour en temps réel. M-07 est stable. M-12 montre une légère dérive depuis hier soir — une alerte orange s'affiche : *"Augmentation progressive de la température de bobinage — surveiller, pas d'intervention immédiate."*

Karim note l'alerte, vérifie visuellement M-12 lors de sa ronde. Il ouvre le Q&A et tape : *"Quelle est la tendance de M-12 sur les 7 derniers jours ?"* — le module LLM est désactivé (pas d'internet), mais le système affiche les courbes et les métriques directement.

À 16h, internet revient. Rien n'a changé pour ML_Elec — aucune donnée perdue, aucun gap dans les historiques, aucune reconfiguration nécessaire.

**Nouvelle réalité :** ML_Elec est le seul outil de l'usine qui n'a pas été impacté par la coupure internet.

**Exigences révélées :** fonctionnement offline total sans dégradation des fonctions core (ingestion, détection, alertes, dashboards), badge offline discret, file locale des actions utilisateur avec resynchronisation silencieuse, LLM désactivable avec fallback vers données structurées, zéro dépendance cloud pour les fonctions critiques.

---

### Journey 7 — Karim : Le dialogue avec la machine 💬 *(Q&A conversationnel)*

**Vendredi, 8h30. Karim inspecte le dashboard.**

M-03 affiche un score de santé en baisse depuis 3 jours, mais pas d'alerte. Karim veut comprendre. Il ouvre le Q&A et tape : *"Pourquoi le score de M-03 baisse ?"*

En 6 secondes, le système répond : *"Le score de santé de M-03 a baissé de 12% sur les 3 derniers jours. La cause principale est une augmentation de 8% de l'amplitude vibratoire sur l'axe vertical. Ce pattern est cohérent avec un début d'usure de roulement, mais le niveau reste sous le seuil d'alerte. Prochaine alerte estimée dans 10-15 jours si la tendance se maintient."*

Karim comprend instantanément. Il note de surveiller M-03 et planifie une vérification visuelle la semaine prochaine. Sans le Q&A, il aurait ignoré la baisse de score ou appelé Sofiane pour interpréter les courbes.

**Nouvelle réalité :** Karim pose des questions en français et obtient des réponses actionnables sans être data scientist.

**Exigences révélées :** Q&A en langage naturel basé sur les données de l'équipement (historique, alertes, RUL), réponse en < 10s, indication claire des limites du système.

### Synthèse des exigences révélées par les journeys

| Capacité | Journeys concernés |
|---|---|
| Triage machine prioritaire en ≤ 3 secondes | J1, J6 |
| Décision actionnable en ≤ 60 secondes | J1, J6 |
| Dashboard santé avec courbe temporelle | J1, J4, J6 |
| Alertes en langage naturel (LLM) avec confiance visible | J1, J2, J7 |
| RUL avec horizon de défaillance | J1, J7 |
| Mécanisme de feedback sur alertes (1 tap + note optionnelle) | J2 |
| Éditeur no-code visuel (drag & drop) | J2, J3 |
| Validation pipeline sur historique avant déploiement | J3 |
| Fallback protocole non supporté (plugin / CSV) | J3 |
| Bibliothèque de plugins installables | J3, J5 |
| Déploiement Docker simplifié | J3, J5 |
| Vue synthétique direction + calcul ROI | J4 |
| Export rapport / partage lecture seule | J3, J4 |
| Documentation développeur + contribution | J5 |
| Fonctionnement offline total sans dégradation | J6 |
| File locale + resynchronisation silencieuse | J6 |
| Q&A conversationnel en langage naturel | J7 |
| Suggestion d'action concrète | J1, J7 |

## Domain-Specific Requirements

Les journeys ci-dessus définissent le *quoi* du point de vue utilisateur. Les sections suivantes encadrent le *comment* en tenant compte des contraintes réglementaires, techniques et d'intégration propres au domaine électrotechnique industriel.

### Compliance & Réglementaire

| Exigence | Statut ML_Elec |
|---|---|
| **IEC 62443** (cybersécurité OT) | Adjacent — ML_Elec accède aux données OT en lecture seule. Isolation réseau recommandée. |
| **RGPD** | S'applique si données liées à des opérateurs identifiables (logs nominatifs). Minimisation des données. |
| **Directive Machines / NIS2** | Non applicable en V1 — outil de monitoring, pas de contrôle actif des process. |
| **Certification matérielle** | Non requise — logiciel edge sur hardware standard (NUC, Raspberry Pi). |

**Principe directeur :** ML_Elec est en **lecture seule** sur les process industriels — il n'actionne rien. Cela élimine la majorité des contraintes réglementaires OT critiques.

### Contraintes Techniques

- **Isolation OT/IT** : Le runtime edge doit pouvoir s'exécuter dans un VLAN isolé sans accès internet
- **Latence** : Acquisition et détection < 5s (cf. NFR03) — pas de contrainte temps-réel strict (< 100ms), mais pas de batch différé
- **Offline total** : Zéro dépendance réseau pour les fonctions core (cf. NFR17) — LLM = optionnel/désactivable
- **Intégrité des données** : Stockage DuckDB + Parquet avec checksums — pas de modification rétroactive des historiques
- **Authentification** : Accès dashboard protégé (auth locale minimum), pas d'exposition internet par défaut

### Intégrations Requises

- **Protocoles V1** : MQTT (broker local), OPC-UA (serveur embarqué)
- **Protocoles V2** : Modbus RTU/TCP, IEC 61850 (via plugins communautaires)
- **GMAO** : Intégration native en V2 (SAP PM, Maximo, CARL) — en V1 : export CSV/JSON manuel
- **LLM providers** : OpenAI, Anthropic, Mistral (API) — désactivable, local (Ollama) en V2

### Risques Domaine & Mitigations

| Risque | Mitigation |
|---|---|
| Fausse alerte → perte de confiance | Mécanisme de feedback, seuils ajustables, règles de corrélation no-code |
| Données capteurs corrompues → modèle dégradé | Validation des données à l'ingestion, alertes qualité données |
| Dépendance API LLM → indisponibilité | LLM désactivable, fallback vers alertes ML seules (cf. NFR14) |
| Accès non autorisé au réseau OT | Déploiement en VLAN isolé, pas d'exposition WAN par défaut |
| Dérive du modèle ML dans le temps | Retraining périodique automatique planifié (V1) ou manuel |

### Philosophie de Gestion des Alarmes

ML_Elec adopte une approche de gestion des alarmes adaptée au contexte PME, inspirée des principes ISA-18.2 / IEC 62682 :

- **Priorisation :** 3 niveaux de sévérité (info / warning / critique) avec fréquence d'affichage contrôlée — les alarmes critiques sont toujours visibles, les infos sont consolidées
- **Prévention du flood :** Regroupement automatique des alarmes corrélées sur un même équipement (une alerte agrégée au lieu de N alertes individuelles)
- **Shelving :** L'ingénieur peut mettre en sourdine une alerte pour une durée configurable (1h, 8h, 24h) avec justification obligatoire tracée dans l'audit log
- **Rationalisation :** Mécanisme de feedback utilisateur (FR10) alimentant un tableau de bord de la qualité des alarmes (ratio utile/fausse alerte par équipement)
- **Objectif V1 :** ≤ 10 alarmes actives simultanées par site en conditions normales — au-delà, le système affiche un indicateur de surcharge et recommande un ajustement des seuils

### Traçabilité des Modifications (MOC Simplifié)

Pour un outil ciblant des PME (pas des sites classés SEVESO), ML_Elec implémente un processus MOC (Management of Change) simplifié :

- Toute modification de seuil d'alerte, de règle de corrélation ou de configuration de pipeline est tracée dans l'audit log avec : auteur, horodatage, valeur avant/après, justification optionnelle (cf. FR31)
- Les modifications de seuils critiques (niveau "critique") nécessitent une confirmation explicite avec avertissement : *"Vous modifiez un seuil critique — cette action sera tracée"*
- **Autorité des seuils critiques :** seuls les rôles Engineer et Admin peuvent modifier un seuil de niveau "critique" (Engineer au niveau équipement, Admin au niveau site)
- Historique des modifications consultable par l'ingénieur et l'administrateur avec filtrage par équipement, type de modification et période
- Rollback de configuration possible sur les 30 dernières modifications (cf. FR33)

## Innovation & Patterns Novateurs

Au-delà des contraintes domaine, ML_Elec introduit plusieurs innovations qui constituent des paris à valider. Cette section cartographie les zones d'innovation, le paysage concurrentiel et l'approche de validation.

### Zones d'Innovation Détectées

**1. Triple combinaison inédite : Edge ML + LLM explicatif + No-code**
Aucun outil du marché ne combine simultanément : inférence ML embarquée offline, explication en langage naturel via LLM, et configuration no-code par drag & drop. ML_Elec est le premier à assembler ces trois couches en open-source. (Cf. Executive Summary — différenciateurs.)

**2. Interface "mémoire machine" — rupture UX**
Le paradigme dominant des outils industriels est le dashboard à seuils : valeur actuelle vs seuil d'alerte. ML_Elec introduit la **courbe de santé temporelle** comme interface principale — l'historique et la trajectoire remplacent la valeur instantanée. Rupture de paradigme UX pour le secteur, validée par Journey 1 (Karim).

**3. Edge AI sans infrastructure**
L'inférence ML tourne sur du hardware à < 200€ (NUC Intel, Raspberry Pi 5) sans GPU, sans cloud, sans data scientist. Réduction de 10x+ du coût d'entrée par rapport aux solutions concurrentes. Validée par Journey 6 (offline total).

**4. LLM appliqué à l'électrotechnique industrielle**
Utilisation du LLM non pas pour générer du contenu, mais pour **traduire des anomalies techniques en langage de terrain** — cas d'usage émergent. Validée par Journey 7 (Q&A conversationnel).

### Contexte Marché & Paysage Concurrentiel

- **Concurrents directs** : Augury, SKF Enlight, Siemens MindSphere — tous cloud-first, tous > 50k€/an, tous nécessitent intégration par prestataire
- **Outils open-source existants** : Apache Kafka + Grafana = pipeline data, pas de ML embarqué, pas d'explication LLM
- **Fenêtre de marché** : 18–24 mois avant que les acteurs établis intègrent LLM + edge dans leurs offres SMB

### Approche de Validation

| Innovation | Hypothèse à valider | Méthode |
|---|---|---|
| Onboarding 15 min | Un technicien non-formé peut déployer seul | Test utilisateur avec 5 techniciens terrain |
| LLM industriel | L'explication LLM est comprise et utile | A/B test : alerte technique vs alerte LLM, mesure de l'action prise |
| Courbe de santé | L'interface temporelle est plus actionnable que les seuils | Comparaison temps de décision avec interface classique |
| Edge ML sans GPU | Isolation Forest / Autoencoder tournent en < 5s sur NUC | Benchmark sur hardware cible |

### Mitigations des Risques d'Innovation

| Risque | Mitigation |
|---|---|
| LLM génère des explications incorrectes | Validation humaine obligatoire avant action, mention "suggestion" claire |
| Modèle ML non adapté au type de machine | Période d'apprentissage explicite + indicateur de confiance du modèle |
| Onboarding 15 min non atteignable | Simplification progressive UI, tests utilisateurs itératifs avant release |

## Exigences Spécifiques IoT/Edge + SaaS B2B

Cette section traduit les contraintes d'innovation et de domaine en exigences concrètes pour un déploiement edge industriel.

### Vue d'ensemble du type de projet

ML_Elec est un runtime edge déployé sur site industriel, gérant un réseau de capteurs/équipements sur réseau local, avec interface web accessible aux différents profils. Chaque instance couvre un site industriel unique avec potentiellement plusieurs dizaines d'appareils surveillés simultanément.

### Hardware & Déploiement Edge

| Cible | Spécifications |
|---|---|
| **Minimum** | Raspberry Pi 5 (8GB RAM, SSD recommandé) |
| **Standard** | Intel NUC Gen 12+ (16GB RAM, 256GB SSD) |
| **Industriel** | Serveur x86 standard (Ubuntu Server 22.04+, Docker) |
| **Système hôte** | Linux uniquement (Ubuntu 22.04 LTS, Debian 12) |
| **Runtime** | Docker + Docker Compose — zéro dépendance native |

**Profil énergétique V1 :** runtime prévu pour des équipements alimentés secteur (pas de contrainte batterie), avec consommation edge recommandée ≤ 25W en fonctionnement nominal.

### Mises à Jour

- **Mécanisme V1** : Mise à jour en un clic depuis le dashboard (pull Docker image + restart)
- **Offline** : Mise à jour depuis registre local (air-gap) ou Docker Hub si réseau disponible
- **Rollback** : Retour version précédente en cas d'échec
- **Notifications** : Alerte dashboard quand nouvelle version disponible

### Modèles ML — Stockage & Partage

| Phase | Comportement |
|---|---|
| **V1** | **Modèle généraliste par défaut distribué avec le logiciel — gratuit, open-source** |
| **V1** | Modèles appris localement stockés en local (ONNX + Parquet), export manuel possible |
| **V2** | Marketplace de modèles — modèles spécialisés par type d'équipement (potentiellement premium) |

### Architecture Multi-Devices (Single-Site)

- **Instance unique** par site — une installation Docker = un réseau OT
- **Multi-devices** : Surveillance de 1 à 50 appareils par instance (cf. NFR18)
- **Multi-sites** : Chaque site = instance indépendante en V1 — fédération P2P en V2

### Modèle de Permissions (RBAC + UX par rôle)

| Rôle | Permissions | Vue par défaut |
|---|---|---|
| **Viewer** (Karim) | Lecture seule — dashboards, alertes, historiques | Dashboard santé machines + courbes |
| **Manager** (M. Benali) | Lecture seule | Vue synthétique ROI — pannes évitées, économies, disponibilité |
| **Engineer** (Sofiane) | Configuration complète — pipelines, modèles, règles, plugins | Éditeur + dashboard complet |
| **Admin** | Gestion utilisateurs, mise à jour système, export | Console administration |

- **Audit log V1** : Toute modification de configuration tracée avec utilisateur + horodatage + valeur avant/après (cf. FR31, NFR11)
- Authentification locale en V1 (username/password) — SSO/LDAP en V2
- Pas d'exposition internet par défaut

### Modèle Commercial

| Composant | Licence | Modèle |
|---|---|---|
| **Core runtime** | MIT open-source | Gratuit |
| **Modèle ML par défaut** | MIT open-source | Distribué avec le logiciel, gratuit |
| **Plugins standards** | MIT open-source | Gratuit, communautaires |
| **Plugins premium** | Propriétaire | Vente unitaire (connecteurs GMAO certifiés, modèles spécialisés…) |
| **Support entreprise** | Commercial | À définir V2+ |

## Scope & Développement Phasé

### Stratégie MVP

**Approche :** MVP Expérience + Validation — valider que Karim détecte une anomalie réelle en moins de 15 minutes et que Sofiane déploie sans prestataire. Tout ce qui ne sert pas ces deux validations est V2+.

**Équipe minimale V1 :**
- 1 dev full-stack TypeScript (pipeline + dashboard)
- 1 dev ML/data (modèles + inférence edge)
- 1 dev DevOps/infra (Docker, OPC-UA, MQTT)
- *(1 solo dev complet possible mais timeline ×2)*

### Découpage UX V1/V1.x

- **V1 (release blocker UX) :** boucle décisionnelle complète "voir en 3s, décider en 60s" sur J1/J2/J3/J6, incluant tri par risque, action recommandée, feedback 1 tap et fallback offline
- **V1.x (post-release rapide) :** optimisation de lisibilité terrain (densité/hiérarchie), instrumentation UX par persona, et amélioration des boucles de récupération de confiance
- **Critère de passage V1 → V1.x :** objectifs NFR28 et NFR29 atteints sur 2 cycles de mesure consécutifs

### V1 — Must-Have (18 features)

| Feature | Justification |
|---|---|
| Ingestion MQTT + OPC-UA | Sans données, rien ne tourne |
| Import CSV données historiques | Validation ML avant connexion matériel — cf. FR44 |
| Détection anomalie (Isolation Forest + Autoencoder) | Cœur de valeur #1 |
| Modèle RUL + horizon de défaillance | Cœur de valeur #2 |
| Modèle ML généraliste par défaut (distribué, gratuit) | Élimine le cold start — sans ça l'onboarding échoue — cf. FR46 |
| Explication LLM en langage naturel (API, désactivable) | Différenciateur clé — Karim doit comprendre l'alerte |
| Suggestion d'action concrète | Commander la pièce, planifier l'arrêt, surveiller — cf. FR45 |
| Q&A conversationnel | Composante LLM interactive — cf. FR13, Journey 7 |
| Éditeur no-code visuel (drag & drop pipelines) | Sofiane doit déployer sans code |
| Dashboard santé + courbe temporelle | Interface "mémoire machine" — moment waouh |
| Onboarding ≤ 15 min (cf. NFR01) | Release blocker absolu |
| RBAC 4 rôles + UX par rôle | Sécurité + adoption décideur |
| Audit log | Traçabilité OT non négociable |
| Mécanisme de mise à jour 1-clic + rollback | Maintenance en production |
| 6 types de plugins avec contrats documentés | Extensibilité de base — cf. FR47 |
| Déploiement Docker single-command | Prérequis onboarding 15 min |
| Gestion des alarmes (priorisation, shelving, feedback) | Qualité des alertes = confiance utilisateur |
| Fonctionnement offline total | Différenciateur #1 — validé par Journey 6 |

### V2 — Post-MVP

LLM local (Ollama), marketplace modèles + plugins, intégration GMAO native (SAP PM, Maximo), synchronisation multi-sites P2P, SSO/LDAP, app mobile, rapport PDF automatique (synthèse mensuelle pour M. Benali), BYOM (Bring Your Own Model), export rapport PDF avancé, modèles ML spécialisés par type d'équipement (premium)

### V3 — Vision

Fédération pair-à-pair de PMEs, modèles anonymisés partagés, standard ouvert européen, plugins premium avancés, Digital Twin (jumeau numérique simulation scénarios maintenance), AR overlay mobile

### Analyse des Risques

| Type | Risque | Mitigation |
|---|---|---|
| **Technique** | LLM API indisponible en contexte OT isolé | LLM désactivable dès V1 — produit fonctionne sans |
| **Technique** | Modèle généraliste insuffisant pour équipements atypiques | Indicateur de confiance + documentation des limites |
| **Marché** | Onboarding 15 min non validé avant release | Tests utilisateurs obligatoires (5 techniciens terrain) avant v1.0 |
| **Ressources** | Projet solo → timeline allongée | Scope V1 figé — aucune feature V2 avant que les 18 must-have soient livrés |

## Exigences Fonctionnelles

Les 52 exigences fonctionnelles ci-dessous constituent le **contrat de capacités** de ML_Elec V1. Chaque FR est traçable aux user journeys et critères de succès définis en amont. Numérotation séquentielle par capacité.

### 1. Acquisition & Ingestion de Données

- **FR01** : L'ingénieur peut connecter une source de données MQTT à un pipeline de traitement
- **FR02** : L'ingénieur peut connecter une source de données OPC-UA à un pipeline de traitement
- **FR03** : Le système peut détecter et signaler dans l'interface les signaux de données manquants, hors plage ou corrompus
- **FR04** : L'ingénieur peut configurer les paramètres de connexion (broker, topic, fréquence) sans modifier de code
- **FR44** : L'ingénieur peut importer un fichier CSV de données historiques pour valider le modèle ML avant connexion au matériel réel
- **FR05** : Le système peut ingérer des données en continu (24h/7j) avec reprise automatique après interruption, sans perte de données — la métrique de disponibilité est définie dans NFR12

### 2. Détection d'Anomalies & Prédiction

- **FR06** : Le système peut détecter des anomalies sur les signaux d'un équipement via le modèle généraliste par défaut
- **FR46** : Le système est livré avec un modèle ML généraliste pré-entraîné (bundled dans l'image Docker), fonctionnel dès l'installation sans configuration ni entraînement préalable, couvrant les anomalies courantes sur moteurs électriques
  - **Critères d'acceptation :** Exécution sur 3 datasets publics de référence : (1) CWRU Bearing Dataset, (2) NASA Turbofan (C-MAPSS), (3) FEMTO Bearing Dataset. Le modèle généraliste détecte ≥ 70% des anomalies connues avec un taux de faux positifs < 20% sur chaque dataset, sans aucune configuration ni entraînement spécifique.
- **FR07** : Le système peut calculer une durée de vie résiduelle (RUL) et un horizon de défaillance estimé
- **FR08** : Le système peut générer une alerte avec niveau de sévérité (info / warning / critique)
- **FR09** : L'ingénieur peut configurer des seuils d'alerte et des règles de corrélation pour réduire les faux positifs
- **FR10** : Le technicien peut soumettre un feedback sur une alerte en 1 tap (utile / fausse alerte), avec note de contexte optionnelle et confirmation immédiate
- **FR11** : Le système peut afficher un indicateur de confiance du modèle ML sur chaque alerte
- **FR50** : Le système peut signaler un contexte ambigu avant action (température ambiante élevée, capteur récemment ajouté, modèle en phase d'apprentissage)
- **FR45** : Le système peut générer une suggestion d'action concrète pour chaque alerte (commander la pièce, planifier l'arrêt, surveiller), affichée dans l'interface aux côtés de l'explication LLM

### 3. Explication & Intelligence (LLM)

- **FR12** : Le système peut générer une explication en langage naturel pour chaque anomalie détectée, compréhensible par un technicien BTS sans formation ML (taux ≥ 80% mesuré par test utilisateur), incluant cause probable, sévérité et horizon temporel
  - **Critères d'acceptation :** Test avec 5 techniciens BTS terrain (sans formation ML). Chaque technicien lit 5 explications générées par le LLM et reformule l'action à entreprendre. Taux de reformulation correcte ≥ 80% (≥ 20/25 réponses correctes). Temps de lecture moyen par explication < 30 secondes.
- **FR13** : Le technicien peut poser une question en langage naturel sur l'état d'un équipement et recevoir une réponse en < 10s, basée sur les données de l'équipement (historique, alertes, RUL). Le système indique clairement lorsqu'une question dépasse sa capacité
- **FR14** : L'administrateur peut activer ou désactiver le module LLM
- **FR15** : L'administrateur peut configurer le fournisseur LLM et la clé API

### 4. Éditeur No-Code & Pipelines

- **FR16** : L'ingénieur peut créer un pipeline de traitement par drag & drop sans code
- **FR17** : L'ingénieur peut modifier un pipeline existant sans redéploiement
- **FR18** : L'ingénieur peut installer un plugin depuis la bibliothèque en un clic
- **FR19** : Le développeur peut créer et publier un plugin selon l'architecture documentée
- **FR47** : Le système supporte 6 types de plugins avec contrats d'interface documentés : (1) connecteurs protocole, (2) transformateurs de données, (3) modèles ML, (4) blocs de visualisation, (5) canaux de notification, (6) exportateurs
- **FR20** : L'ingénieur peut tester un pipeline en mode simulation sur des données historiques avant déploiement en production — cf. Journey 3 (Sofiane valide le pipeline avant mise en service)

### 5. Dashboards & Visualisation

- **FR21** : Le technicien peut visualiser la courbe de santé temporelle de chaque équipement surveillé
- **FR22** : Le technicien peut consulter l'historique des alertes et leur statut
- **FR23** : Le manager peut accéder à une vue synthétique ROI (pannes évitées, économies, disponibilité)
- **FR24** : L'ingénieur peut accéder à une vue technique affichant : liste des pipelines et leur statut, modèles ML actifs et indicateur de confiance, 1 000 dernières entrées logs avec filtrage par composant et sévérité
- **FR25** : Tout utilisateur peut exporter des données en CSV/JSON
- **FR26** : L'ingénieur peut partager un lien de vue en lecture seule
- **FR48** : Le système affiche un bandeau d'état discret (online / offline / resynchronisation) avec statut de la file locale des actions utilisateur
- **FR49** : Le technicien peut voir la fraîcheur des données (timestamp de la dernière mesure fiable) sur la vue machine et chaque alerte

### 6. Gestion des Équipements

- **FR27** : L'ingénieur peut créer, nommer et archiver des équipements surveillés dans le système

### 7. Notifications

- **FR28** : Le technicien peut recevoir une notification (push, email ou SMS) lors d'une alerte sur un équipement surveillé, incluant machine, sévérité, horizon et lien direct vers la fiche concernée
- **FR28b** : Le technicien peut configurer par machine les canaux de notification, les niveaux de sévérité déclencheurs et les plages horaires de réception
- **FR51** : Le technicien peut partager un résumé décisionnel (machine, risque, action recommandée, horizon) via WhatsApp/email/copie ; en mode offline, l'envoi est mis en file locale puis synchronisé au retour réseau

### 8. Administration & Sécurité

- **FR29** : L'administrateur peut créer, modifier et supprimer des comptes utilisateurs avec attribution de rôle
- **FR30** : Chaque utilisateur accède à une interface adaptée à son rôle dès la connexion
- **FR31** : Le système enregistre un audit log de toute modification de configuration avec auteur, horodatage et valeur avant/après
- **FR32** : L'administrateur peut déclencher une mise à jour du logiciel depuis le dashboard
- **FR33** : L'administrateur peut effectuer un rollback vers la version précédente
- **FR34** : Le système notifie l'administrateur lorsqu'une nouvelle version est disponible
- **FR35** : Le système distingue les plugins open-source des plugins sous licence propriétaire et en informe l'utilisateur lors de l'installation
- **FR36** : L'administrateur peut consulter l'état de santé du système (services actifs, espace disque, connexions protocoles, uptime)
- **FR37** : Le système peut générer une alerte système lorsqu'une ressource critique atteint un seuil (disque, mémoire, déconnexion broker)

### 9. Onboarding & Déploiement

- **FR38** : Un nouvel utilisateur peut déployer ML_Elec via une commande Docker unique
- **FR39** : Un nouvel utilisateur peut connecter un premier équipement et visualiser un premier signal actif sur le dashboard dans un délai ≤ 15 min après déploiement (cf. NFR01), sans aide externe
  - **Critères d'acceptation :** Test avec 5 ingénieurs électrotechniques sans expérience Docker. Chaque participant part de zéro (machine vierge avec Docker pré-installé). Succès = premier signal actif visible sur le dashboard. Cible : 5/5 participants complètent en ≤ 15 min. Échec = plus de 1 participant dépasse 15 min ou abandonne.
- **FR40** : Le système fournit un guide d'onboarding en 5 étapes max, du premier lancement à la visualisation du premier signal, avec indicateur de progression, fallback données simulées si nécessaire et taux de complétion cible ≥ 90%
- **FR41** : L'ingénieur peut déployer ML_Elec dans un environnement totalement isolé (air-gap)
- **FR42** : Le développeur peut accéder à une documentation de contribution et créer un plugin dans un environnement local Docker

### 10. Internationalisation

- **FR43** : L'utilisateur peut sélectionner la langue de l'interface parmi les langues supportées (FR + EN minimum V1)

### Modèle de Données Conceptuel

Les entités ci-dessous constituent le cœur du modèle de données de ML_Elec. Ce modèle est intentionnellement haut niveau — il définit les concepts métier et leurs relations sans prescrire d'implémentation technique.

| Entité | Description | Relations clés |
|---|---|---|
| **Équipement** | Machine physique surveillée (moteur, compresseur…) | possède N Capteurs, génère N Alertes, est assigné à N Pipelines |
| **Capteur** | Source de données liée à un équipement (vibration, température…) | appartient à 1 Équipement, alimente N Pipelines |
| **Pipeline** | Chaîne de traitement configurée (ingestion → ML → alerte) | traite N Capteurs, utilise 1 ModèleML, génère N Alertes |
| **Alerte** | Événement détecté avec sévérité et explication | concerne 1 Équipement, générée par 1 Pipeline, possède N Feedbacks |
| **ModèleML** | Modèle d'inférence (généraliste ou entraîné localement) | utilisé par N Pipelines, produit des scores et RUL |
| **Utilisateur** | Personne authentifiée avec un rôle RBAC | possède 1 Rôle, auteur de N entrées AuditLog |
| **AuditLog** | Trace immuable d'une modification de configuration | référence 1 Utilisateur, horodaté, append-only |
| **Plugin** | Extension installable (protocole, modèle, visualisation…) | appartient à 1 des 6 types documentés (cf. FR47), utilisable par N Pipelines |

## Exigences Non-Fonctionnelles

Les 29 NFR ci-dessous définissent **comment** le système doit performer. Seules les catégories pertinentes pour ML_Elec sont documentées — chaque exigence est spécifique, mesurable et testable.

### Performance

- **NFR01** : L'onboarding complet (déploiement Docker → premier signal actif affiché) est réalisable en **≤ 15 minutes** par un ingénieur électrotechnique sans expérience Docker préalable — **release blocker absolu** — vérifié par test utilisateur avec 5 techniciens BTS terrain
- **NFR02** : Les pages dashboard se chargent en **< 2 secondes** sur hardware cible (mini-PC entrée de gamme, 4 Go RAM) — vérifié par test de charge avec 10 équipements affichés simultanément
- **NFR03** : Latence entre détection d'anomalie et notification utilisateur **< 5 secondes** — vérifié par injection d'anomalie simulée et mesure du délai bout-en-bout
- **NFR04** : Ingestion de données continue 24h/7j à **≥ 1 000 points/seconde** par équipement sans perte ni dégradation — vérifié par test de charge avec 50 équipements simulés générant 1 000 pts/s chacun pendant 24h
- **NFR05** : Inférence ML exécutée en **< 500 ms** par cycle d'analyse sur hardware cible — vérifié par benchmark sur Raspberry Pi 5 (8GB) et Intel NUC Gen 12
- **NFR06** : Taux de faux positifs du modèle généraliste par défaut **< 5%** sur les datasets de référence, mesuré sur fenêtre glissante de 30 jours — vérifié par exécution sur 3 datasets publics de référence (CWRU Bearing, NASA Turbofan, FEMTO Bearing)
- **NFR28** : Sur le dashboard opérationnel, la machine prioritaire est identifiable en **≤ 3 secondes** dans **≥ 90%** des sessions terrain — vérifié par instrumentation UX (timestamp ouverture dashboard → ouverture machine prioritaire)
- **NFR29** : Après ouverture d'une alerte warning/critique, une action utilisateur (planifier / surveiller / partager) est enregistrée en **≤ 60 secondes** dans **≥ 80%** des cas — vérifié par instrumentation UX (timestamp ouverture alerte → timestamp action)

### Sécurité

- **NFR07** : Données stockées chiffrées au repos via **AES-256** — vérifié par audit de la couche de stockage et tentative de lecture directe des fichiers sans clé
- **NFR08** : Toutes les communications chiffrées en transit : **TLS 1.2+** pour API REST, MQTT over TLS, OPC-UA Security Mode Sign & Encrypt — vérifié par scan TLS (sslyze ou équivalent) sur tous les ports exposés
- **NFR09** : RBAC appliqué sur **100% des endpoints API** — aucune escalade de privilèges possible entre les 4 rôles (Viewer, Manager, Engineer, Admin) — vérifié par test de pénétration automatisé sur chaque endpoint avec chaque rôle
- **NFR10** : Conformité **RGPD** : aucune donnée personnelle ni donnée industrielle transmise hors de l'instance locale sans consentement explicite de l'administrateur ; clés API et secrets stockés chiffrés, jamais en clair — vérifié par audit réseau (capture trafic sortant pendant 48h) et revue du code de stockage des secrets
- **NFR11** : Audit log **inaltérable** (append-only, non supprimable par un utilisateur quel que soit son rôle), conservé **minimum 12 mois** avec horodatage, auteur et valeur avant/après — vérifié par tentative de suppression/modification avec rôle Admin et vérification d'intégrité après 12 mois

### Fiabilité

- **NFR12** : Disponibilité système **≥ 99,5%** (hors maintenance planifiée annoncée), mesurée par monitoring interne (cf. FR36) — vérifié par calcul d'uptime sur 30 jours consécutifs via health check endpoint
- **NFR13** : **Zéro perte de données** en cas de coupure électrique brutale — mécanisme de journalisation préalable à l'écriture (write-ahead) sur toutes les écritures critiques — vérifié par test de coupure alimentation simulée (kill -9) et vérification d'intégrité des données après redémarrage
- **NFR14** : **Dégradation gracieuse** si le LLM est indisponible ou désactivé — toutes les fonctions ML, dashboards et alertes restent opérationnelles avec un temps de réponse dégradé de **< 20%** par rapport au mode nominal ; seules les explications en langage naturel et le Q&A sont désactivés — vérifié par désactivation du module LLM et exécution de la suite de tests fonctionnels complète
- **NFR15** : Reprise automatique après redémarrage matériel sans intervention manuelle en **< 60 secondes** — vérifié par redémarrage forcé du hardware et mesure du temps jusqu'au premier health check positif
- **NFR16** : Backup et restore de la configuration complète (pipelines, utilisateurs, seuils, modèles) en **une seule commande** CLI — vérifié par exécution backup, suppression de l'instance, restore et validation fonctionnelle complète
- **NFR17** : Fonctionnement **offline continu illimité** — aucune fonctionnalité critique (ingestion, détection, alertes, dashboards) ne dépend d'une connexion internet — vérifié par déconnexion réseau totale et exécution de la suite de tests fonctionnels pendant 72h

### Scalabilité

- **NFR18** : Support de **1 à 50 équipements** surveillés par instance sans dégradation de performance **> 10%** — vérifié par test de charge progressif (1, 10, 25, 50 équipements simulés) et mesure des temps de réponse dashboard et latence détection
- **NFR19** : Support de **1 à 20 utilisateurs concurrents** par instance avec temps de réponse conforme aux NFR de performance — vérifié par test de charge avec 20 sessions simultanées exécutant des scénarios d'usage typiques
- **NFR20** : Rotation automatique des données historiques configurable par l'administrateur — rétention par défaut **12 mois**, extensible selon espace disque disponible — vérifié par configuration de rétention à 1 jour et vérification de la suppression automatique des données dépassées
- **NFR21** : Architecture plugin extensible — ajout de nouveaux types de capteurs, protocoles ou modèles ML via les 6 types de plugins documentés (cf. FR47), avec un temps d'intégration d'un nouveau plugin **< 4 heures** pour un développeur familier avec l'API — vérifié par développement d'un plugin de test depuis la documentation et mesure du temps de bout-en-bout

### Accessibilité

- **NFR22** : Interface conforme **WCAG 2.1 niveau AA** minimum — contraste, navigation clavier, lecteur d'écran, et parité clavier sur les parcours critiques J1/J2/J3/J6 — vérifié par audit automatisé (axe-core ou Lighthouse) avec score ≥ 90% et validation manuelle sur ces parcours
- **NFR23** : Interface **responsive** et touch-first — utilisable sur tablette (≥ 768px) pour usage terrain par les techniciens de maintenance, avec cibles tactiles minimales **44×44 px** sur actions primaires — vérifié par test fonctionnel sur tablette physique (iPad / Samsung Tab) couvrant les journeys Karim et Sofiane
- **NFR24** : Architecture **i18n** dès V1 — français et anglais livrés, ajout d'une langue supplémentaire possible en **< 2 heures** par un traducteur non-technique (fichier de traduction unique, sans recompilation ni redéploiement) — vérifié par ajout d'une langue de test (espagnol) et mesure du temps nécessaire

### Intégration

- **NFR25** : Protocoles industriels **MQTT 3.1.1+** et **OPC-UA** supportés nativement avec configuration guidée et validation de connexion — vérifié par test de connexion à un broker MQTT (Mosquitto) et un serveur OPC-UA (open62541) avec données simulées
- **NFR26** : Export de données en **CSV et JSON** depuis tout dashboard en **< 10 secondes** pour un jeu de 100 000 lignes — vérifié par export d'un dataset de 100 000 lignes et mesure du temps de génération
- **NFR27** : **API REST documentée** (OpenAPI 3.0) exposant les endpoints de lecture (équipements, alertes, métriques) pour intégration avec systèmes tiers (GMAO, ERP) — endpoints d'écriture réservés V2+ — vérifié par génération de la documentation OpenAPI et exécution de requêtes de test sur chaque endpoint documenté

## Glossaire

| Terme | Définition |
|---|---|
| **RUL** | Remaining Useful Life — durée de vie résiduelle estimée d'un équipement avant défaillance |
| **MQTT** | Message Queuing Telemetry Transport — protocole de messagerie léger pour l'IoT et les capteurs |
| **OPC-UA** | Open Platform Communications Unified Architecture — standard de communication industrielle pour l'automatisation |
| **ONNX** | Open Neural Network Exchange — format portable de modèles ML interopérable entre frameworks |
| **Parquet** | Format de fichier colonnaire optimisé pour le stockage et l'analyse de grandes quantités de données |
| **DuckDB** | Base de données analytique embarquée, conçue pour l'analyse de données locales sans serveur |
| **LLM** | Large Language Model — modèle de langage capable de générer et comprendre du texte en langage naturel |
| **RBAC** | Role-Based Access Control — contrôle d'accès basé sur les rôles utilisateur |
| **BYOM** | Bring Your Own Model — possibilité de remplacer le modèle ML par défaut par un modèle personnalisé |
| **MOC** | Management of Change — processus de traçabilité des modifications de configuration |
| **GMAO** | Gestion de Maintenance Assistée par Ordinateur (CMMS en anglais) |
| **Edge** | Calcul en périphérie — traitement des données sur site, au plus près de la source, sans dépendance cloud |
| **Isolation Forest** | Algorithme ML de détection d'anomalies basé sur l'isolation des observations atypiques |
| **Autoencoder** | Réseau de neurones entraîné à reconstruire ses entrées, utilisé pour la détection d'anomalies par erreur de reconstruction |
| **IEC 62443** | Norme internationale de cybersécurité pour les systèmes d'automatisation et de contrôle industriels |
| **VLAN** | Virtual Local Area Network — réseau local virtuel pour la segmentation réseau |

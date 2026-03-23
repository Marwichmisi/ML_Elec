---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments:
  - "_bmad-output/brainstorming/brainstorming-session-2026-03-03-00-13-24.md"
  - "_bmad-output/planning-artifacts/research/market-supervision-industrielle-hybride-research-2026-03-03.md"
  - "_bmad-output/planning-artifacts/research/domain-genie-electrique-electrotechnique-research-2026-03-03.md"
  - "_bmad-output/planning-artifacts/research/technical-faisabilite-ml-elec-research-2026-03-03.md"
date: 2026-03-04
author: Marwane
---

# Product Brief: ML_Elec

<!-- Content will be appended sequentially through collaborative workflow steps -->

---

## Executive Summary

ML_Elec est une plateforme open source de monitoring électrique et de maintenance prédictive conçue pour les techniciens terrain — sans data scientist, sans infrastructure cloud imposée, sans complexité de configuration. Elle détecte les anomalies sur les moteurs et systèmes électrotechniques avant la panne, explique les résultats en langage naturel via un LLM intégré, et s'installe en moins de 15 minutes. ML_Elec comble le fossé entre les outils industriels puissants mais inaccessibles (Ignition, AVEVA) et les stacks DIY instables, en offrant la première plateforme de maintenance prédictive véritablement accessible aux PME industrielles — open source, offline-first, et extensible via un écosystème de plugins visuels.

---

## Core Vision

### Problem Statement

Les moteurs électriques et systèmes électrotechniques tombent en panne de façon imprévue. Chaque arrêt non planifié coûte entre 1 000 € et 1 000 000 €/heure selon le secteur industriel — sans compter les coûts de remplacement prématuré des équipements. Ce n'est pas un problème de technologie : la maintenance prédictive existe depuis des années. C'est un problème d'accessibilité fondamental : les solutions actuelles sont trop chères, trop complexes à configurer, ou réservées aux équipes disposant de data scientists dédiés. Le technicien terrain — celui qui connaît ses machines mieux que quiconque — est exclu de cette révolution.

### Problem Impact

- **Arrêts machine non planifiés** : pertes financières immédiates allant de 1 000 € à 1 000 000 €/heure selon le secteur (IBM Data Breach Report, 2025)
- **Remplacements prématurés** : CAPEX gaspillé sur des équipements encore fonctionnels
- **Demande massive non satisfaite** : 73 % des responsables industriels identifient la maintenance prédictive comme priorité n°1, mais la majorité des PME industrielles n'ont pas les outils accessibles pour la mettre en œuvre (Grand View Research IIoT, 2025)
- **Marché PME sous-servi** : le segment PME industrielle (50–500 employés) représente un TAM estimé à 2–3 milliards € en Europe, sans solution adaptée à ce jour

### Pourquoi les Solutions Existantes Échouent

| Solution | Problème |
|---|---|
| **Ignition** (Inductive Automation) | $7 500–15 000+ de licence, courbe d'apprentissage élevée, ML non natif |
| **Grafana** | Excellent en visualisation, zéro connectivité PLC/OPC-UA native, zéro ML |
| **AVEVA / Schneider** | Enterprise-grade, tarification opaque, roadmap anxiogène post-acquisition |
| **Stacks DIY** (Node-RED + InfluxDB) | Gratuites mais instables, non maintenables en production critique |
| **Azure ML / AWS SageMaker** | Requièrent un data scientist pour configurer et maintenir les modèles |

Aucune solution ne combine simultanément : connectivité industrielle native + ML intégré + accessibilité technicien + prix PME + offline-first.

### Proposed Solution

ML_Elec est une plateforme edge open source qui :
- Se connecte aux équipements industriels via **MQTT et OPC-UA** (protocoles natifs, zéro middleware)
- **Détecte les anomalies automatiquement** avec un modèle ML pré-installé (Isolation Forest / Autoencoder) — fonctionnel dès l'installation
- **Explique en langage naturel** ce qui se passe grâce à un LLM intégré (le technicien comprend sans être data scientist)
- **Fonctionne offline-first** : données souveraines, résilience edge, synchronisation intelligente à la reconnexion
- **S'installe en 15 minutes** : de l'installation à la première courbe affichée avec détection d'anomalie active

### Key Differentiators

1. **🧠 Accessible sans data scientist** — Modèle ML pré-configuré livré à l'installation + LLM qui traduit les prédictions en langage humain. Un technicien comprend le résultat sans formation spécialisée.
2. **⚡ Onboarding 15 minutes** — De l'installation à la première détection d'anomalie active. Contrainte de design non-négociable, pas une ambition vague.
3. **🔓 Open source & offline-first** — Souveraineté totale des données, zéro lock-in, fonctionne en edge industriel sans connexion permanente. Stockage Parquet/DuckDB léger et interrogeable.
4. **🎯 Architecture IA bicouche unique** — Couche 1 : ML pour prédire (modèle remplaçable via BYOM). Couche 2 : LLM pour expliquer et dialoguer (API providers configurables). Aucun outil open source industriel ne propose ça nativement.
5. **💰 Pricing accessible PME** — Open source, zéro licence, zéro per-tag. Le seul coût est l'infrastructure edge (à partir de ~80 € sur Raspberry Pi).

---

## Target Users

### Utilisateurs Primaires

---

**👤 Persona 1 — Karim, Technicien de Maintenance (35 ans)**

> *"Je veux savoir que le moteur va tomber en panne AVANT qu'il le fasse. Pas après."*

**Contexte terrain :** Karim travaille dans une PME industrielle de 80 personnes (agroalimentaire, papeterie ou chimie). BTS Électrotechnique, tablette en main sur le terrain, mains souvent occupées, environnement bruyant et chaud. Zéro formation data science. C'est **l'utilisateur quotidien** de ML_Elec — il ne l'installe pas, il l'utilise.

**Sa douleur aujourd'hui :**
- Aucun outil centralisé de monitoring — surveillance visuelle uniquement
- Chaque panne = arrêt production non planifié de 24–48h → 15 000 € à 100 000 € de pertes
- Les solutions existantes (Ignition, Grafana) sont trop chères ou trop complexes pour lui

**Ce que ML_Elec lui apporte :**
- Interface mobile-first, alertes couleur vert/orange/rouge, texte lisible d'un coup d'œil
- Langage naturel : *"Le moteur M3 présente un défaut de roulement. Intervention recommandée d'ici 10 jours."*
- Onboarding **Quick Start 15 minutes** — plug-and-play, modèle ML par défaut, zéro configuration

**Son moment "aha" :** Première intervention préventive réussie. Son responsable lui demande comment il a su. Il montre ML_Elec.

---

**👤 Persona 2 — Sofiane, Ingénieur Électricien (42 ans)**

> *"Je veux configurer une fois, analyser en profondeur, et ne plus jamais appeler un consultant externe."*

**Contexte :** Ingénieur en bureau d'études ou responsable maintenance dans une PME de 200 personnes. Maîtrise les normes IEC, les variateurs ABB/Siemens, les schémas de protection. C'est **le champion d'adoption** — il évalue, déploie, et forme Karim. C'est aussi lui qui convainc M. Benali de signer.

**Sa douleur aujourd'hui :**
- Configure manuellement les seuils d'alarme sur chaque équipement — fastidieux
- Les solutions ML nécessitent un data scientist qu'il n'a pas
- Les données quittent le site avec les solutions cloud — problème de souveraineté

**Ce que ML_Elec lui apporte :**
- Configure les connexions OPC-UA via l'éditeur visuel React Flow sans code
- Crée des règles personnalisées, remplace le modèle ML par le sien (BYOM)
- Analyse les tendances longue durée avec SQL sur les fichiers Parquet locaux
- Onboarding **Power Setup 1–2 heures** — configuration avancée, workflows custom

**Son rôle de décideur :** Il est le **champion interne**. C'est lui qui évalue ML_Elec, pilote le déploiement sur 3 moteurs, démontre le ROI à M. Benali, puis déploie sur toute l'installation.

---

### Utilisateurs Secondaires

**🏭 M. Benali, Responsable de Production / Directeur Technique (50 ans)**

N'utilise pas ML_Elec directement. Il consulte les synthèses que Sofiane lui présente. Son seul intérêt : **ROI visible et chiffré**. Il a besoin de voir : *"On a évité 3 pannes ce mois-ci, soit 45 000 € de pertes évitées."* C'est lui qui débloque le budget. ML_Elec doit prévoir un rapport d'impact automatique pour ce moment de vente interne.

---

**🎓 Yasmine, Étudiante en Génie Électrique (22 ans) — Persona de Croissance Communautaire**

Utilisatrice secondaire aujourd'hui, ingénieure cliente dans 2 ans. Elle utilise ML_Elec pour ses projets de fin d'études sur des bancs de test moteur, contribue à la communauté open source, et devient un vecteur d'adoption organique dans l'industrie. Pas un segment commercial immédiat — un investissement dans la croissance long terme de la communauté.

---

### User Journeys

**🔧 Karim — Du problème à l'ambassadeur :**

| Étape | Ce qui se passe | Émotion |
|---|---|---|
| **Découverte** | Sofiane installe ML_Elec et lui montre l'interface | Curieux, sceptique |
| **Onboarding** | Quick Start 15 min — voit les premières courbes de ses moteurs | *"C'est vraiment simple ?"* |
| **Usage quotidien** | Consulte le score de santé chaque matin sur sa tablette | Rassuré, en contrôle |
| **Moment "aha"** | Première alerte préventive → intervient avant la panne → évite 3 jours d'arrêt | Fier, converti |
| **Long terme** | Ambassadeur interne — forme ses collègues, propose l'extension à d'autres lignes | Champion |

**👷 Sofiane — De l'évaluation au déploiement :**

| Étape | Ce qui se passe | Émotion |
|---|---|---|
| **Découverte** | Cherche "maintenance prédictive open source" → trouve ML_Elec sur GitHub | Intrigué |
| **Onboarding** | Power Setup 1-2h — connecte OPC-UA, configure workflows, valide modèle ML | En confiance |
| **Pilote** | Déploie sur 3 moteurs critiques, collecte 30 jours de données | Analytique |
| **Moment "aha"** | Présente le rapport d'impact à M. Benali — 2 pannes évitées, ROI prouvé | Convaincu |
| **Long terme** | Déploiement sur toute l'installation, partenaire de l'écosystème ML_Elec | Ambassadeur |

---

## Success Metrics

ML_Elec est un succès quand les techniciens réduisent les pannes imprévues — et quand les responsables peuvent le prouver en euros.

### North Star Metric

**Nombre de moteurs surveillés en production active**
_(moteurs avec au moins une alerte générée dans les 30 derniers jours)_

Ce chiffre capture la vraie valeur créée : des machines réellement surveillées, des alertes réelles générées. C'est la métrique qui prouve que ML_Elec fonctionne — pas juste installé.

---

### User Success Metrics

| Métrique | Cible | Mesure |
|---|---|---|
| **Time-to-First-Alert** | < 24h post-installation (idéal : < 1h) | Timestamp first alert − timestamp install |
| **Ratio Préventif/Réactif** | > 50% alertes → intervention avant panne | (Pannes évitées / Pannes totales) × 100 |
| **ROI Automatique** | Calculé et affiché à chaque panne évitée | € économisés = (durée arrêt estimée) × (coût/h site) |
| **Rétention 30 jours** | > 60% des déploiements actifs après 30j | Déploiements actifs J30 / Déploiements totaux |
| **Fréquence d'ouverture** | ≥ 3×/semaine par utilisateur actif | Sessions hebdomadaires par utilisateur |

**Signal d'alarme :** Ouverture < 3×/semaine = problème de valeur perçue.
Taux de faux positifs > 30% = churn silencieux à venir.

---

### Business Objectives

**3 mois (validation produit) :**
- 5 déploiements pilotes actifs en production réelle
- 1 panne évitée documentée avec ROI chiffré
- Onboarding 15 min validé par 3 utilisateurs indépendants (sans aide)

**12 mois (traction marché) :**
- 50 déploiements actifs en production
- 500 GitHub Stars (signal de communauté)
- 1 retour client publié (case study avec ROI chiffré)
- Taux de rétention 30 jours > 60%

---

### Key Performance Indicators

| KPI | Cible 3 mois | Cible 12 mois | Source |
|---|---|---|---|
| Déploiements actifs | 5 | 50 | Telemetry (opt-in) |
| GitHub Stars | 50 | 500 | GitHub |
| Time-to-First-Alert médian | < 2h | < 1h | Logs install |
| Ratio préventif/réactif | Mesuré | > 50% | Dashboard ML_Elec |
| Taux rétention 30j | Mesuré | > 60% | Telemetry (opt-in) |
| Faux positifs | Mesuré | < 20% | Feedback utilisateur |
| NPS (opt-in) | Premier feedback | > 40 | In-app survey |

---

## MVP Scope

### Core Features (V1.0 — Launch)

Le MVP de ML_Elec est défini par une contrainte non-négociable :
**de l'installation à la première prédiction de panne en moins de 15 minutes.**

#### 🔌 Acquisition de données
- **OPC-UA** — connectivité PLC/SCADA native (standard industriel PME)
- **MQTT** — connectivité IoT/capteurs légers
- **Import CSV** — validation du ML sur données historiques avant connexion matériel réel

#### 🧠 Détection & Prédiction
- Modèle ML pré-installé (Isolation Forest / Autoencoder) actif dès l'installation
- Score de santé moteur en temps réel (vert / orange / rouge)
- Détection automatique d'anomalies sans configuration manuelle

#### 🔔 Alertes & Explications
- Alerte structurée : *"Moteur M3 : anomalie détectée, intervenir sous 10 jours"*
- Explication LLM en langage naturel : *"Signature vibratoire anormale, probable défaut de roulement"*
- Historique de la courbe de dégradation (visualisation de la tendance)
- **Suggestion d'action concrète** : commander la pièce, planifier l'arrêt, surveiller
- Alertes **SMS / Email** — notification hors interface

#### 📊 Reporting
- **Rapport PDF automatique** — synthèse mensuelle pour M. Benali
  (pannes évitées, ROI chiffré en €, état des moteurs)

---

### Extensions V1.x (Post-Launch, premières itérations)

| Feature | Valeur | Complexité | Priorité |
|---|---|---|---|
| **BYOM** — Bring Your Own Model | Sofiane remplace le modèle ML par le sien | Moyenne | V1.1 |
| **API publique — lecture seule** | Expose scores de santé + alertes (pas d'OAuth, pas de webhooks complexes) | Faible | V1.1 |
| **Export CSV + Webhook sortant** | Sofiane exporte les données vers SAP/Maximo manuellement | Faible | V1.2 |

---

### Out of Scope V1

| Feature | Raison | Horizon |
|---|---|---|
| **Multi-site** | Complexité architecture, pas de besoin V1 validé | V2 |
| **App mobile native** iOS/Android | PWA responsive couvre Karim en V1 | V2 |
| **Marketplace de plugins** | Écosystème tiers prématuré sans base communautaire | V2 |
| **Intégration GMAO native** (SAP, Maximo) | Projet dans le projet — remplacé par webhook sortant | V2 |

---

### MVP Success Criteria

Le MVP V1.0 est validé quand :

1. ✅ **Onboarding** : 3 utilisateurs indépendants voient leur première alerte en < 15 min sans aide
2. ✅ **Prédiction** : au moins 1 panne réelle évitée documentée avec ROI chiffré
3. ✅ **Rétention** : > 60% des déploiements pilotes actifs après 30 jours
4. ✅ **Confiance** : taux de faux positifs < 20% sur les 30 premiers jours

**Signal go/no-go V1.x :** 4 critères atteints → déployer les extensions V1.x.

---

### Future Vision (V2+)

- **Plateforme multi-site** : dashboard centralisé pour flottes industrielles (50+ moteurs)
- **App mobile native** : alertes push + AR overlay sur le moteur pour Karim
- **Marketplace de plugins** : connecteurs, modèles ML, dashboards communautaires
- **Intégration GMAO native** : création automatique d'ordres de travail SAP/Maximo
- **ML fédéré** : apprentissage partagé entre déploiements (opt-in, données souveraines)
- **Digital Twin** : jumeau numérique pour simulation de scénarios de maintenance

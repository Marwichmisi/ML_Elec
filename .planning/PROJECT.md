# ML_Elec

## What This Is

ML_Elec est une plateforme logicielle modulaire open source (Apache 2.0) pour la maintenance prédictive et l'intelligence industrielle. Basée sur une architecture micro-noyau avec plugins, elle transforme des signaux techniques de capteurs en décisions fiables, prioritaires et actionnables avant la panne. Conçue pour tourner sur des équipements légers (Raspberry Pi, ESP32), elle cible les responsables maintenance, chefs d'atelier et ingénieurs qui ont besoin d'un système modulaire et extensible.

## Core Value

La collecte de données capteurs en temps réel et la détection d'anomalies fiables — c'est le fondamental. Si tout le reste échoue, le système doit acquérir les signaux et détecter les dérives.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Core micro-noyau Go : bus NATS, gestion plugins, API REST, SQLite
- [ ] Plugin acquisition MQTT : collecte données capteurs ESP32 → Core
- [ ] Plugin détection d'anomalie : seuil + détection basique (scikit-learn)
- [ ] Dashboard React + TypeScript : visualisation temps réel des métriques
- [ ] Communication Core ↔ Plugins via Child process + JSON-RPC stdin/stdout
- [ ] Configuration centralisée (core + plugins activables/désactivables)
- [ ] Démo live soutenance : moteur électrique + capteurs vibration/température/courant
- [ ] README professionnel avec getting started en 5 minutes
- [ ] Documentation technique pour développeurs de plugins

### Out of Scope

- Digital twin 3D — trop complexe pour v1
- Intégration GMAO/ERP — plugins premium futurs
- Mémoire opérationnelle (CBR + Knowledge Graph) — v2+
- Monétisation plugins premium — v2+
- Support cloud/SaaS — v1 est on-premise uniquement
- Mobile app — web-first

## Context

- **Projet de soutenance 3ème année** avec démo réelle capteurs
- **Premier projet open source** de Marwane
- **Recherche existante** : analyse des plateformes (ThingsBoard, EdgeX, FIWARE), briques OS réutilisables, évaluation from scratch vs réutilisation
- **Brainstorming complet** : 4 piliers fondamentaux identifiés (confiance, actionnabilité, mémoire, priorisation métier)
- **Technologies validées** : Core Go, plugins Python, dashboard React/TS, MQTT, SQLite
- **Architecture** : micro-noyau + plugins isolés (child process + JSON-RPC)
- **Bus interne** : NATS (performant, multi-langage, léger)
- **Budget matériel** : ~200€ (ESP32 + Raspberry Pi + capteurs)
- **Durée** : 12 mois

## Constraints

- **Licence**: Apache 2.0 — compatible usage commercial, pas de copyleft
- **Hardware**: ESP32 + Raspberry Pi + capteurs vibration/température/courant (~200€)
- **Edge-first**: doit tourner sur Raspberry Pi, pas besoin de serveur puissant
- **Isolation plugins**: crash plugin ne doit pas tuer le core
- **Multi-langage plugins**: Go core, Python plugins (et d'autres langages possibles)
- **Délai soutenance**: ~12 mois

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Core en Go | Performance, compilation statique, binary unique, concurrence native | — Pending |
| Plugins en Python | Écosystème ML/IoT riche, facilité développement, communauté | — Pending |
| Dashboard React + TypeScript | Type safety, écosystème riche, maintainabilité | — Pending |
| Bus interne NATS | Performant, multi-langage (Go+Python), léger, embedded possible | — Pending |
| Communication JSON-RPC stdin/stdout | Isolation processus, modèle éprouvé (HashiCorp go-plugin) | — Pending |
| Stockage SQLite | Simplicité, pas de dépendance serveur, suffisant pour v1 | — Pending |
| MQTT comme protocole IoT | Standard industriel, léger, adapté ESP32 | — Pending |
| Licence Apache 2.0 | Permissive, protection brevets, professionnelle | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-06-30 after initialization*

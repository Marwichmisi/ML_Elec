# Requirements: ML_Elec

**Defined:** 2026-06-30
**Core Value:** La collecte de données capteurs en temps réel et la détection d'anomalies fiables

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Core

- [ ] **CORE-01**: Micro-noyau Go avec bus NATS embedded pour communication interne
- [ ] **CORE-02**: Plugin manager avec isolation child process (HashiCorp go-plugin pattern)
- [ ] **CORE-03**: API REST pour interaction externe (dashboard, configuration)
- [ ] **CORE-04**: Stockage SQLite avec mode WAL pour séries temporelles
- [ ] **CORE-05**: Configuration centralisée (core + plugins activables/désactivables)
- [ ] **CORE-06**: Plugin SDK avec contrats versionnés
- [ ] **CORE-07**: Contrainte taille core < 5000 LOC (prévention core bloat)

### Acquisition

- [ ] **ACQ-01**: Plugin acquisition MQTT : collecte données capteurs ESP32 → Core via NATS
- [ ] **ACQ-02**: Support MQTT QoS 0/1/2 pour fiabilité variable
- [ ] **ACQ-03**: Stockage données capteurs en SQLite avec timestamps
- [ ] **ACQ-04**: Gestion des assets/machines (enregistrement, hiérarchie)

### Détection

- [ ] **DET-01**: Détection par seuil configurable (température, vibration, courant)
- [ ] **DET-02**: Détection ML basique avec IsolationForest (scikit-learn)
- [ ] **DET-03**: Scoring de confiance des alertes (0-100%)
- [ ] **DET-04**: Boucle de feedback opérateur (confirmer/rejeter alerte)

### Visualisation

- [ ] **VIS-01**: Dashboard temps réel avec métriques capteurs en direct
- [ ] **VIS-02**: Graphiques historiques (séries temporelles)
- [ ] **VIS-03**: Contexte décisionnel (cause probable, action recommandée, priorité)
- [ ] **VIS-04**: Widgets configurables par l'utilisateur
- [ ] **VIS-05**: Alertes avec notifications (email pour v1)

### Démo Soutenance

- [ ] **DEMO-01**: Scénario démo live : moteur électrique + capteurs vibration/température/courant
- [ ] **DEMO-02**: Pipeline end-to-end fonctionnel : capteur → MQTT → Core → Détection → Dashboard
- [ ] **DEMO-03**: README professionnel avec getting started en 5 minutes
- [ ] **DEMO-04**: Documentation technique pour développeurs de plugins

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Acquisition Avancée

- **ACQ-05**: Support Modbus TCP
- **ACQ-06**: Support OPC UA
- **ACQ-07**: Fusion multi-capteurs (vibration + température + courant combinés)

### Mémoire Opérationnelle

- **MEM-01**: Fiche de cas (Case-Based Reasoning) après incident
- **MEM-02**: Règles métier (DMN/Drools) pour scénarios connus
- **MEM-03**: Knowledge Graph pour relations actif-capteur-symptôme-intervention
- **MEM-04**: Validation mémoire (tracé si hypothèse retenue/rejetée)

### Intégration

- **INT-01**: Intégration GMAO/ERP
- **INT-02**: Support cloud/SaaS
- **INT-03**: Mobile app

### Avancé

- **ADV-01**: Digital twin 2D/3D
- **ADV-02**: Apprentissage supervisé après 50+ événements confirmés
- **ADV-03**: Migration TimescaleDB pour haute cardinalité

## Out of Scope

| Feature | Reason |
|---------|--------|
| Digital twin 3D | Trop complexe pour v1, coût rendu élevé |
| Intégration GMAO/ERP | Plugins premium futurs, pas core value |
| Mémoire opérationnelle (CBR + KG) | v2+, nécessite historique de données |
| Cloud/SaaS | v1 on-premise uniquement, edge-first |
| Mobile app | Web-first, mobile après |
| OAuth/2FA | Pas critique pour MVP industriel |
| Multi-tenant | v1 = single user/site |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| CORE-01 | Phase 1 | Pending |
| CORE-02 | Phase 1 | Pending |
| CORE-03 | Phase 1 | Pending |
| CORE-04 | Phase 1 | Pending |
| CORE-05 | Phase 1 | Pending |
| CORE-06 | Phase 1 | Pending |
| CORE-07 | Phase 1 | Pending |
| ACQ-01 | Phase 2 | Pending |
| ACQ-02 | Phase 2 | Pending |
| ACQ-03 | Phase 2 | Pending |
| ACQ-04 | Phase 2 | Pending |
| DET-01 | Phase 3 | Pending |
| DET-02 | Phase 3 | Pending |
| DET-03 | Phase 3 | Pending |
| DET-04 | Phase 3 | Pending |
| VIS-01 | Phase 4 | Pending |
| VIS-02 | Phase 4 | Pending |
| VIS-03 | Phase 4 | Pending |
| VIS-04 | Phase 4 | Pending |
| VIS-05 | Phase 4 | Pending |
| DEMO-01 | Phase 5 | Pending |
| DEMO-02 | Phase 5 | Pending |
| DEMO-03 | Phase 5 | Pending |
| DEMO-04 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 24 total
- Mapped to phases: 24
- Unmapped: 0 ✓

## Agent Development Notes

**Compétences obligatoires pour l'agent développeur:**
- **Go skills** : Les skills Go (golang-*) sont obligatoires pour le développement du core. L'agent doit suivre les bonnes pratiques Go (error handling, naming, structs/interfaces, testing, lint) pour produire du code de qualité.
- **Flutter skills** : Les skills Flutter sont obligatoires pour le développement du dashboard/UI. L'agent doit les utiliser pour la couche présentation.

---
*Requirements defined: 2026-06-30*
*Last updated: 2026-06-30 after initial definition*

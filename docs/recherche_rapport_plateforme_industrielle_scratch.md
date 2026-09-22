
# Évaluation de la Viabilité d'une Plateforme Industrielle "From Scratch"

**Auteur :** Manus AI
**Date :** 17 Juin 2026

## 1. Diagnostic de Viabilité

La construction d'une plateforme logicielle industrielle modulaire de maintenance prédictive "from scratch" (sans réutiliser de briques open source majeures comme base technique) est **viable mais très coûteuse et risquée**. 

Le développement d'un micro-kernel et d'un système de plugins robustes exige une expertise pointue en architecture système et en communication inter-processus (IPC). Si cette approche offre un contrôle total et une légèreté maximale, elle implique de recréer des mécanismes complexes (gestion du cycle de vie, isolation, sécurité, résilience) qui sont généralement fournis par des frameworks existants. Le risque principal réside dans la création d'une dette technique importante sur des composants non différenciants (comme le bus d'événements ou le stockage de séries temporelles) au détriment des fonctionnalités métier (IA, maintenance prédictive).

Cependant, pour un produit ciblant l'industrie, les ingénieurs et l'IoT simple (ESP32, Raspberry Pi), une architecture légère et agnostique peut se justifier si les choix technologiques sont pragmatiques et si l'on accepte d'intégrer certaines briques fondamentales (comme la base de données) plutôt que de les recréer.

## 2. Comparaison des Options Technologiques

### 2.1. Le Core (Micro-Kernel)

Le choix du langage pour le cœur du système est crucial pour garantir la performance, la sécurité et la légèreté.

| Langage | Avantages | Inconvénients | Complexité | Recommandation (choix figé Go) |
| :--- | :--- | :--- | :--- | :--- |
| **Go** | Concurrence facile (goroutines), compilation statique, excellent pour les microservices et l'IPC, cross-compile ARM trivial, libs indus matures (Paho, gopcua, modbus), HashiCorp go-plugin éprouvé. | Garbage collector (pauses ms, acceptables pour démo 1-50Hz). | Moyenne | **Recommandé et figé pour V1/MVP soutenance**. |
| **Rust** | Performance native, sécurité mémoire garantie, très léger. | Courbe d'apprentissage très raide, écosystème moins mature pour protocoles industriels, retarderait le MVP 12 mois. | Élevée | **Écarté pour tout le projet**. |
| **Python** | Développement rapide, écosystème IA/Data riche. | Lenteur (GIL), consommation mémoire, difficile à packager comme binaire autonome. | Faible | Non recommandé pour le Core. |
| **Node.js** | Asynchrone par nature, écosystème riche. | Single-threaded, consommation mémoire, moins adapté aux contraintes industrielles strictes. | Moyenne | Non recommandé. |

![Arbitrage Performance vs Complexité](core_tradeoff.png)

### 2.2. Les Plugins

Le développement des plugins prioritairement en Python est un choix **très pertinent**. Python est le standard de facto pour l'analyse de données, l'IA et l'ingénierie. Il permet aux utilisateurs finaux (ingénieurs, data scientists) de créer facilement leurs propres modules d'acquisition ou d'analyse.

### 2.3. Communication Core ↔ Plugins (IPC)

La communication entre le Core et les plugins isolés (processus enfants) doit être rapide et fiable.

| Technologie | Avantages | Inconvénients | Latence estimée | Recommandation |
| :--- | :--- | :--- | :--- | :--- |
| **JSON-RPC sur stdin/stdout** | Extrêmement simple, pas de port réseau à gérer, isolation naturelle. | Parsing JSON coûteux, difficile à déboguer, flux unidirectionnel complexe à gérer. | ~150 µs | **Recommandé pour le MVP** (simplicité). |
| **gRPC sur Unix Domain Sockets (UDS)** | Typage fort (Protobuf), très performant, streaming bidirectionnel natif. | Complexité de mise en œuvre, dépendance lourde (gRPC). | ~116 µs [1] | **Recommandé pour la V1/V2** (robustesse). |
| **Sockets locaux (TCP/UDS) purs** | Très rapide, contrôle total. | Tout le protocole de sérialisation/routage est à recréer. | ~11 µs | Trop complexe à maintenir. |

![Comparaison Latence IPC](ipc_latency_comparison.png)

### 2.4. Stockage des Séries Temporelles

Recréer un moteur de séries temporelles "from scratch" est un **piège de complexité majeur**. La gestion de la compression, de la rétention et des requêtes analytiques sur des millions de points est extrêmement difficile.

| Option | Avantages | Inconvénients | Recommandation |
| :--- | :--- | :--- | :--- |
| **Moteur "Fait Maison"** | Contrôle total, intégration parfaite. | Effort de R&D colossal, risque de perte de données, performances incertaines. | **Non recommandé**. |
| **PostgreSQL + TimescaleDB** | Robuste, SQL standard, gère les métadonnées et les séries temporelles, très performant [2]. | Empreinte mémoire plus importante, trop lourd pour MVP RPi soutenance. | Option V2 post-soutenance uniquement. |
| **SQLite (avec extensions)** | Ultra-léger, parfait pour l'edge/embarqué RPi, zéro dépendance serveur, idéal soutenance. | Limité en concurrence et en volume de données. | **Recommandé et figé pour MVP**. |

## 3. Architecture Cible V1

L'architecture proposée repose sur un modèle **Edge-first**, où la plateforme peut tourner localement sur un équipement industriel (IPC, Raspberry Pi) tout en pouvant se synchroniser avec un cloud.

![Architecture V1](architecture_v1.png)

**Composants clés (stack figée Go) :**
*   **Core en Go :** Gère le cycle de vie des processus enfants (plugins Python), le routage des messages (Bus NATS/MQTT), l'authentification et l'accès au stockage. Choix figé, Rust écarté.
*   **Plugins (Python) :** Exécutés comme des processus isolés. Ils communiquent avec le Core Go via JSON-RPC sur stdin/stdout (pour le MVP) ou gRPC sur UDS (pour la V1 finale).
*   **Stockage :** SQLite pour le MVP edge (RPi) - choix figé soutenance. PostgreSQL + TimescaleDB uniquement en option V2 post-soutenance (trop lourd pour RPi MVP).
*   **Frontend :** React + TypeScript servie par le Core Go, communiquant via API REST ou WebSocket.

## 4. Évaluation des Risques et Coûts Cachés

1.  **Gestion du cycle de vie des processus (Watchdog) :** Maintenir des processus enfants Python en vie, gérer les crashs, les fuites mémoire et les redémarrages sans perdre de données est complexe.
2.  **Sérialisation/Désérialisation :** Le passage constant de données (surtout à haute fréquence) entre le Core et les plugins via JSON peut devenir un goulot d'étranglement CPU.
3.  **Packaging et Distribution :** Distribuer une application contenant un binaire Core Go et un environnement Python isolé avec ses dépendances (pour les plugins) sur différentes architectures (x86, ARM via `GOARCH=arm64`) est un défi DevOps majeur, mais simplifié par le binaire statique Go.
4.  **Sécurité :** L'isolation des plugins (sandboxing) pour éviter qu'un plugin malveillant ou buggé ne compromette le système hôte nécessite des mécanismes OS spécifiques (cgroups, namespaces sous Linux).

## 5. Recommandation Finale et MVP

**Que construire d'abord (MVP - stack figée Go) ?**
1.  Un **Core minimal en Go (choix définitif, Rust écarté)** capable de spawner un processus enfant.
2.  Un système de communication **JSON-RPC sur stdin/stdout** entre Core Go et plugins Python.
3.  Un **Plugin d'acquisition MQTT puis Modbus TCP** (en Python).
4.  Un **Plugin de stockage** écrivant dans **SQLite (MVP edge)**.
5.  Un **Dashboard basique React+TS** affichant les données en temps réel.
6.  **Optionnel P5 visual-flow (Node-RED, hors MVP)** : plugin Node.js pour programmation visuelle des flux, managé par Core Go, via MQTT + REST uniquement.

**Pourquoi ?**
Ce MVP valide le risque technique principal : la stabilité et la performance de l'architecture "Core + Plugins isolés via IPC". Il permet de démontrer la modularité de la solution sans s'éparpiller sur des fonctionnalités avancées (IA, alertes complexes) qui pourront être ajoutées sous forme de nouveaux plugins dans les versions ultérieures.

---
**Références :**
[1] F. Werner, "Using gRPC for (local) inter-process communication", MPI-HD, 2021. URL: https://www.mpi-hd.mpg.de/personalhomes/fwerner/research/2021/09/grpc-for-ipc/
[2] A. Kulkarni, "When Boring is Awesome: Building a scalable time-series database on PostgreSQL", Timescale Blog, 2017. URL: https://medium.com/timescale/when-boring-is-awesome-building-a-scalable-time-series-database-on-postgresql-2900ea453ee2
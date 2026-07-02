
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

| Langage | Avantages | Inconvénients | Complexité | Recommandation |
| :--- | :--- | :--- | :--- | :--- |
| **Rust** | Performance native, sécurité mémoire garantie, très léger, idéal pour l'embarqué/edge. | Courbe d'apprentissage très raide, écosystème moins mature pour certains protocoles industriels. | Élevée | **Recommandé** pour une V1 robuste et performante. |
| **Go** | Concurrence facile (goroutines), compilation statique, excellent pour les microservices et l'IPC. | Garbage collector (latence imprévisible), moins bas niveau que Rust. | Moyenne | **Alternative solide** si l'équipe maîtrise moins Rust. |
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
| **PostgreSQL + TimescaleDB** | Robuste, SQL standard, gère les métadonnées et les séries temporelles, très performant [2]. | Empreinte mémoire plus importante qu'une solution purement embarquée. | **Recommandé**. |
| **SQLite (avec extensions)** | Ultra-léger, parfait pour l'edge/embarqué. | Limité en concurrence et en volume de données. | Alternative pour l'IoT très contraint. |

## 3. Architecture Cible V1

L'architecture proposée repose sur un modèle **Edge-first**, où la plateforme peut tourner localement sur un équipement industriel (IPC, Raspberry Pi) tout en pouvant se synchroniser avec un cloud.

![Architecture V1](architecture_v1.png)

**Composants clés :**
*   **Core (Rust ou Go) :** Gère le cycle de vie des processus enfants (plugins), le routage des messages (Bus), l'authentification et l'accès au stockage.
*   **Plugins (Python) :** Exécutés comme des processus isolés. Ils communiquent avec le Core via JSON-RPC sur stdin/stdout (pour le MVP) ou gRPC sur UDS (pour la V1 finale).
*   **Stockage :** PostgreSQL avec l'extension TimescaleDB pour unifier le stockage relationnel (utilisateurs, configuration) et les séries temporelles (données capteurs).
*   **Frontend :** Une application web légère (React ou Vue) servie par le Core, communiquant via une API REST ou WebSocket.

## 4. Évaluation des Risques et Coûts Cachés

1.  **Gestion du cycle de vie des processus (Watchdog) :** Maintenir des processus enfants Python en vie, gérer les crashs, les fuites mémoire et les redémarrages sans perdre de données est complexe.
2.  **Sérialisation/Désérialisation :** Le passage constant de données (surtout à haute fréquence) entre le Core et les plugins via JSON peut devenir un goulot d'étranglement CPU.
3.  **Packaging et Distribution :** Distribuer une application contenant un binaire Core (Rust/Go) et un environnement Python isolé avec ses dépendances (pour les plugins) sur différentes architectures (x86, ARM) est un défi DevOps majeur.
4.  **Sécurité :** L'isolation des plugins (sandboxing) pour éviter qu'un plugin malveillant ou buggé ne compromette le système hôte nécessite des mécanismes OS spécifiques (cgroups, namespaces sous Linux).

## 5. Recommandation Finale et MVP

**Que construire d'abord (MVP) ?**
1.  Un **Core minimal en Go** (plus rapide à développer que Rust pour une première itération) capable de spawner un processus enfant.
2.  Un système de communication **JSON-RPC sur stdin/stdout**.
3.  Un **Plugin d'acquisition Modbus TCP** (en Python).
4.  Un **Plugin de stockage** écrivant dans une base PostgreSQL/TimescaleDB.
5.  Un **Dashboard basique** affichant les données en temps réel.

**Pourquoi ?**
Ce MVP valide le risque technique principal : la stabilité et la performance de l'architecture "Core + Plugins isolés via IPC". Il permet de démontrer la modularité de la solution sans s'éparpiller sur des fonctionnalités avancées (IA, alertes complexes) qui pourront être ajoutées sous forme de nouveaux plugins dans les versions ultérieures.

---
**Références :**
[1] F. Werner, "Using gRPC for (local) inter-process communication", MPI-HD, 2021. URL: https://www.mpi-hd.mpg.de/personalhomes/fwerner/research/2021/09/grpc-for-ipc/
[2] A. Kulkarni, "When Boring is Awesome: Building a scalable time-series database on PostgreSQL", Timescale Blog, 2017. URL: https://medium.com/timescale/when-boring-is-awesome-building-a-scalable-time-series-database-on-postgresql-2900ea453ee2
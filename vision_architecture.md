# Vision et Architecture : Plateforme d'Intelligence Industrielle (Industrie 4.0)

## 1. Concept Principal
Une plateforme logicielle légère, modulaire et agnostique pour le génie électrique et l'électrotechnique, spécialisée dans la **maintenance prédictive** et le suivi d'équipements industriels (moteurs, compresseurs, etc.). 

Le logiciel ne se veut pas être un monolithe lourd, mais un **"Hub" central** fonctionnant avec un système d'extensions (architecture Micro-Kernel, similaire à VS Code).

## 2. Public Cible
- **Entreprises et Industrie lourde** : Pour le monitoring en temps réel, la prédiction de pannes et l'optimisation des coûts.
- **Ingénieurs professionnels** : Pour le développement de modèles d'IA sur mesure et l'intégration de protocoles spécifiques.
- **Étudiants et Curieux** : Grâce à une version de base accessible et facile à déployer sur des projets IoT simples (ESP32, Raspberry Pi).

## 3. Architecture Technique Validée (Core + Plugins)

### Le Cœur (Core System)
Le cœur du logiciel est extrêmement léger. Il a pour unique responsabilité de :
- Gérer le bus d'événements (Event Bus).
- Gérer les utilisateurs et les permissions.
- Héberger le moteur de base de données (Séries Temporelles, ex: TimescaleDB).
- Gérer le cycle de vie des plugins (chargement, activation, communication).
- Fournir une interface utilisateur de base (Dashboard unifié).

### L'Écosystème de Plugins (Extensions)
Toute fonctionnalité supplémentaire est un plugin (hot-pluggable). L'utilisateur n'installe que ce dont il a besoin.
1. **Plugins d'Acquisition (In/Out)** : Traduisent les protocoles spécifiques en format standard pour le Core.
   - *Exemples* : Plugin MQTT (ESP32, IoT), Plugin Modbus TCP, Plugin OPC UA, Plugin HTTP REST.
2. **Plugins de Traitement & IA** :
   - Plugin de détection d'anomalies de base (Non supervisé).
   - Module d'intégration "Bring Your Own Model" (BYOM) : Permet aux Data Scientists d'importer leurs propres modèles (ex: ONNX entraînés sur Kaggle/PyTorch).
3. **Plugins d'Action et d'Export** :
   - Plugin Alerte SMS/Email.
   - Plugin d'export Excel/CSV.

## 4. Fonctionnalités de Base (À définir pour la V1)
*(Le système Core a besoin de fonctionnalités "out of the box" pour ne pas être une coquille vide. Les plugins viennent s'y greffer.)*

- **Acquisition de base** : Support natif d'une API simple (HTTP/JSON ou MQTT basique) pour valider l'ingestion de données (vibrations, température, courant).
- **Stockage optimisé** : Sauvegarde des séries temporelles à haute fréquence.
- **Visualisation dynamique** : Un dashboard natif pour visualiser les courbes de données en temps réel et l'historique (sans dépendre d'outils externes complexes comme Grafana pour les cas simples).
- **Moteur d'IA (V1)** : Un algorithme statistique simple de détection de déviation (anomalie non supervisée) fourni par défaut.
- **Gestionnaire d'extensions** : L'interface permettant d'activer/désactiver les plugins.

## 5. Le Système de Plugins (Les 4 Piliers)

Pour garantir une flexibilité maximale, l'architecture repose sur quatre familles de plugins, écrites prioritairement en **Python** pour faciliter leur développement par la communauté et les Data Scientists.

### 5.1 Les Plugins "Data Sources" (Connecteurs)
- **Rôle** : Traduire les protocoles industriels et IoT en format standard pour le Core.
- **Exemples** : Modbus TCP (Automates), OPC UA, MQTT Avancé, LoRaWAN, API Météo.
- **Fonctionnement** : Un script Python qui tourne en tâche de fond, se connecte au matériel et pousse les données vers le Core via une API interne ou un broker (gRPC / MQTT local).

### 5.2 Les Plugins "Analytics & AI" (Cerveaux)
- **Rôle** : Analyser les données en temps réel pour prédire des pannes (Maintenance Prédictive) ou effectuer des calculs électriques.
- **Exemples** : Algorithme Deep Learning personnalisé (BYOM), Calcul de déphasage (Cos φ), Analyse FFT de vibrations.
- **Fonctionnement** : Le script Python souscrit aux flux de données du Core, effectue ses inférences (ex: via TensorFlow/PyTorch) et publie des "Alertes" ou de nouvelles métriques vers le Core.

### 5.3 Les Plugins "Visualisation" (Dashboards & 3D)
- **Rôle** : Étendre l'interface utilisateur avec de nouveaux widgets.
- **Exemples** : Jumeau Numérique (Digital Twin) intégrant des modèles 3D (Blender/.glb) réagissant aux données, Cartographie d'usine en 2D, Graphes spectraux complexes.
- **Fonctionnement** : Ces plugins fourniront probablement du code Frontend (Javascript/Vue/React) que le Core injectera dynamiquement dans son interface Web.

### 5.4 Les Plugins "Actions & Notifications" (Bras Actuateurs)
- **Rôle** : Réagir aux alertes générées par le système ou par les plugins d'IA.
- **Exemples** : Envoi de SMS/Email (Twilio/SendGrid), Création de ticket de maintenance (Jira/SAP), Envoi d'une consigne d'arrêt d'urgence à un automate.
- **Fonctionnement** : Un script Python qui écoute le bus d'événements (Event Bus) du Core et déclenche une action externe lorsqu'une condition est remplie.

## 6. L'Intégration et le Langage de Développement (Python)
Le choix stratégique est de baser le développement des plugins (en particulier l'acquisition, l'IA et les actions) sur **Python**.
- **Pourquoi ?** Python est le langage standard de l'IA (Machine Learning) et du traitement de données. Il est extrêmement accessible pour les ingénieurs (électriques, data scientists) et les étudiants, qui n'ont pas forcément vocation à maîtriser des langages de backend complexes (Go, Rust, etc.).
- **Comment ? (Techniquement)** :
  1. Le plugin est un dossier contenant a minima un `plugin_info.json` (nom, version, permissions requises) et un script principal `main.py`.
  2. Architecture **"Child Process"** (façon VS Code) : Le Core lance le plugin Python en tant que processus enfant indépendant en tâche de fond. Ainsi, si un plugin IA effectue un calcul trop lourd ou plante, cela n'affecte ni les performances du Core ni l'interface utilisateur.
  3. **Communication Standard** : Le Core et le script Python s'échangent les données (demandes d'actions, alertes, données de capteurs) via l'entrée/sortie standard (stdin/stdout) ou un Socket TCP local, en utilisant un format lisible comme **JSON-RPC**. C'est une méthode extrêmement légère, robuste, et standardisée (similaire au fonctionnement du LSP).

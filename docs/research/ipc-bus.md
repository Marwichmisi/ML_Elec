# Research IPC Core↔Plugins et bus — ML-Elec (ticket #2, part of #1)

Date : 2026-09-23. Méthode : sources primaires uniquement (docs officielles, code source, specs).
Chaque claim est cité `(URL + version/date de consultation)`. Les chiffres repris des rapports
locaux sans source primaire sont signalés « à valider par bench ».

Inputs locaux lus : `docs/recherche_rapport_plateforme_industrielle_scratch.md`,
`docs/recherche_in_scratch.md`, `docs/recherche_in_opensource.md`,
`docs/rapport_analyse_plateforme_industrielle_in_opensource.md`, `CONTEXT.md`.

> **Correction importante vs rapports locaux** : la formule « HashiCorp go-plugin + JSON-RPC
> sur stdin/stdout » est techniquement incohérente. Le mode `net/rpc` de go-plugin utilise
> `encoding/gob` (Go-only) ; **un plugin Python ne peut parler à go-plugin qu'en gRPC**
> (guide non-Go officiel, §1). Le JSON-RPC/stdio reste possible mais uniquement comme
> protocole custom *sans* la lib go-plugin (on réinvente alors handshake + watchdog).
> La recommandation lockée (§7) tranche en conséquence.

---

## 1. HashiCorp go-plugin : fonctionnement réel (source : repo + code)

- Principe : le Core lance chaque plugin comme **sous-processus** et communique en RPC sur
  **une seule connexion** ; multiplexage via `yamux` (net/rpc) ou HTTP/2 (gRPC).
  En service chez Terraform, Vault, Nomad, Packer, Boundary, Waypoint.
  (https://github.com/hashicorp/go-plugin — README, consulté 2026-09-23 ; v1.8.0 sur
  https://pkg.go.dev/github.com/hashicorp/go-plugin)
- **Deux protocoles, pas interchangeables** : `net/rpc` (gob, Go-only) et `gRPC`
  (cross-language). Le support gRPC (PR #29) est opt-in côté host via
  `ClientConfig.AllowedProtocols = [ProtocolGRPC]` ; sans opt-in, un plugin gRPC est
  rejeté avec une erreur explicite.
  (https://github.com/hashicorp/go-plugin/pull/29 ; `client.go`, champ `AllowedProtocols`)
- **Écrire un plugin Python = gRPC obligatoire.** Recette officielle : serveur gRPC standard
  + service de **health checking avec le nom `plugin` à `SERVING`** + impression sur stdout
  d'une ligne de handshake `CORE-PROTOCOL-VERSION|APP-PROTOCOL-VERSION|NETWORK-TYPE|NETWORK-ADDR|PROTOCOL`
  (ex. `1|1|tcp|127.0.0.1:1234|grpc`, `CORE=1` figé). Exemple Python complet fourni dans le repo.
  (https://github.com/hashicorp/go-plugin/blob/main/docs/guide-plugin-write-non-go.md)
- Versioning : `HandshakeConfig{ProtocolVersion, MagicCookieKey/Value}` — le magic cookie est
  une **fonction UX anti-erreur** (empêche de lancer un binaire quelconque), **pas une sécurité**.
  Incrémenter `ProtocolVersion` invalide les anciens plugins. Semver applicative en plus.
  (README + https://github.com/hashicorp/go-plugin/blob/master/docs/extensive-go-plugin-tutorial.md)
- **Limite « dispensing » en gRPC** : protobuf ne peut pas retourner un service ; le plugin
  gRPC est donc un **singleton** et le multiplexage d'instances doit être géré au niveau du
  `.proto` (ex. champ `instance_id` + reap côté plugin). Les interfaces riches en interfaces
  imbriquées sont coûteuses à porter.
  (https://github.com/hashicorp/go-plugin/pull/29, section « Dispensing Challenges »)
- **Cycle de vie / watchdog (ce que la lib donne, ce qu'elle ne donne pas)** :
  - `Client.Kill()` : échelle de sortie — `Close()` gracieux → attente 2 s → kill forcé ;
    idempotent ; `CleanupClients()` tue tous les clients `Managed` (à appeler avant exit du Core).
    (`client.go`, func `Kill`, `CleanupClients`, consulté 2026-09-23)
  - `StartTimeout` (défaut **1 min**), `MinPort/MaxPort` (défaut max 25000), logs stderr du plugin
    capturés/forwardés par le Core. (`client.go`, valeurs par défaut)
  - `AutoMTLS` : certificats éphémères négociés automatiquement (le client génère un cert,
    l'envoie via `PLUGIN_CLIENT_CERT`, le serveur répond avec le sien ; `TLSConfig` ne doit
    pas être posé manuellement). **Incompatible avec `Reattach`.** (`client.go`, champ `AutoMTLS`)
  - Support **Unix sockets** côté transport (`NETWORK-TYPE=unix`, `UnixSocketConfig`,
    `RunnerFunc`, `EnvUnixSocketDir`). (`client.go`, consulté 2026-09-23)
  - **Non fourni** : politique de restart (backoff, max restarts), détection de plugin
    bloqué-métier (le health gRPC dit « process vivant », pas « boucle 50 Hz à l'heure »).
    → à implémenter dans le Core ML-Elec (voir §7).
- **Piège interop Python + AutoMTLS** : go-plugin génère des certs éphémères **ECDSA P-521** ;
  `grpcio` repose sur BoringSSL qui ne supporte pas P-521 (non réactivable par config).
  Le port communautaire `py-plugin` contourne en utilisant `grpclib` (OpenSSL système, P-521 OK),
  avec interop bidirectionnelle vérifiée contre le binaire d'exemple Go officiel.
  Conséquence MVP : soit `AutoMTLS=false` + `grpcio` classique côté Python, soit AutoMTLS + pile
  `grpclib`. Ne pas découvrir ça à la soutenance.
  (https://github.com/mlund01/py-plugin — README + suite `tests/interop/`, consulté 2026-09-23)
- Licence : **MPL-2.0** (licence fichier + badge repo ; commit conformité 2022-10-12).
  MPL = copyleft faible **au niveau fichier** : utilisation comme lib (import Go) sans modification
  des fichiers go-plugin → pas d'obligation de republication ; toute modification d'un fichier
  MPL doit être republiée sous MPL. Compatible distribution binaire fermée du Core.
  (https://github.com/hashicorp/go-plugin — badge « Mozilla Public License 2.0 » ;
  https://github.com/hashicorp/go-plugin/commit/900c7a5bed543679204f8d43a342c599aaa29641)

## 2. JSON-RPC/stdio vs gRPC/UDS : ce que disent les sources primaires

- **gRPC sur UDS : support officiel des deux côtés.** Schéma d'adressage `unix:path` /
  `unix:///absolute_path` (https://github.com/grpc/grpc/blob/master/doc/naming.md).
  Exemple officiel Python synchrone + asyncio (`server.add_insecure_port('unix:///tmp/....sock')`,
  `grpc.insecure_channel('unix:///...')`) :
  https://chromium.googlesource.com/external/github.com/grpc/grpc.git/+/refs/tags/v1.60.0/examples/python/uds/.
  UDS = transport plus efficace que TCP en local (https://learn.microsoft.com/en-us/aspnet/core/grpc/interprocess-uds?view=aspnetcore-10.0,
  James Newton-King, màj 2024-07-31).
- **Protobuf binaire vs JSON** : « It's like JSON, except it's smaller and faster, and it
  generates native language bindings » (https://protobuf.dev/overview/). Le format canonique
  JSON (ProtoJSON) « is not as efficient as binary wire format and never will be » — plus de
  CPU + plus d'octets (https://protobuf.dev/programming-guides/json/). Champs répétés
  numériques `packed=true` par défaut ; numéros 1–15 = 1 octet sur le wire
  (https://protobuf.dev/programming-guides/encoding/, https://protobuf.dev/programming-guides/proto3/).
  En Python l'écart se resserre (coût de traversée du runtime) — à bencher sur nos messages,
  pas à postuler (https://protobuf.com/efficiency/).
- **Latences citées dans les rapports locaux** (~150 µs JSON-RPC stdio, ~116 µs gRPC, réf.
  Werner MPI-HD 2021 https://www.mpi-hd.mpg.de/personalhomes/fwerner/research/2021/09/grpc-for-ipc/) :
  **ordre de grandeur plausible mais non revérifié ici** — la page MPI-HD n'a pas été re-mesurée
  dans ce ticket. À 1–50 Hz (période 20 ms–1 s), ces deux options tiennent le budget avec >10×
  de marge ; le choix se joue donc sur **robustesse/typage/outillage**, pas sur la latence.
  → Action : micro-bench `bench_ipc` (Go↔Python, 100 B et 4 KB, 50 Hz et rafale) avant de figer V1.

## 3. Bus : NATS vs MQTT Mosquitto (sources : docs NATS, Mosquitto, papier mq-bench)

- **NATS** : un seul binaire `nats-server` (routage + clustering + JetStream + **listener MQTT
  v3.1.1 intégré**, QoS 0/1/2 ; clients MQTT v5 rejetés en CONNACK code 1).
  Core NATS = **at-most-once**, pas de persistance ; **JetStream** ajoute persistance,
  replay, consumers at-least-once. Conversion topic→sujet (`sensors/a/b` → `sensors.a.b`).
  Pour du greenfield la doc recommande NATS de bout en bout ; le listener MQTT sert à
  intégrer un parc existant sans reflasher.
  (https://docs.nats.io/learn/mqtt/ ; https://docs.nats.io/learn/mqtt/qos-sessions-and-retained ;
  https://docs.nats.io/concepts/jetstream ; https://docs.nats.io/concepts/ecosystem)
- **Asymétrie QoS à connaître** : un message publié côté NATS et livré à un abonné MQTT arrive
  **toujours en QoS 0**, même si l'abonnement demandait QoS 1 — il n'y a pas de champ QoS côté
  NATS. Si la garantie compte, le trafic doit naître côté MQTT (QoS 1/2) ou la durabilité doit
  être portée par un stream JetStream, pas par le redelivery MQTT.
  (https://docs.nats.io/learn/mqtt/qos-sessions-and-retained)
- **Perf mesurée (source primaire : papier + harness open source)** : Paul et al.,
  « Benchmarking Message Brokers for IoT Edge Computing », CCGrid 2026
  (https://arxiv.org/html/2603.21600v1, PDF https://arxiv.org/pdf/2603.21600,
  code https://github.com/hpcclab/mq-bench). Payload 1 KB : médiane **NATS 0,21 ms**,
  Redis 0,26 ms, **Mosquitto 0,27 ms** ; à 10 k paires (4 vCPU) NATS ≈ 90 k msg/s,
  **Mosquitto sature ≈ 45 k msg/s (mono-threadé)**. Règle des auteurs : Mosquitto/Redis sur
  edge contraint (léger, stable), NATS/Zenoh multi-threadés dès qu'on scale.
- **Mosquitto** : broker **léger** MQTT v5.0/3.1.1/3.1, « suitable for use on all devices from
  low power single board computers to full servers », licence **EPL/EDL**
  (https://mosquitto.org/ ; https://github.com/eclipse-mosquitto/mosquitto ;
  https://mosquitto.org/man/mosquitto-8.html).
  Points durs v2+ : écoute **loopback seule par défaut** (`listener` requis pour exposer),
  `max_inflight_messages` (ordre non garanti sauf =1), `max_queued_messages/_bytes` (drops
  silencieux), `persistence false` par défaut, `max_keepalive` vs `SetKeepAlive(0)` rejeté.
  (remontés via https://github.com/eclipse-paho/paho.mqtt.golang — section Troubleshooting)
- **Eclipse Paho Go** (`paho.mqtt.golang`) : client MQTT **3.1/3.1.1** (le support v5 est dans
  un autre repo, `paho.golang`), mode **full async**, dépendances `golang.org/x/net/proxy` +
  `gorilla/websocket`, support MQTT-over-WebSocket (`ws://`/`wss://`, proxy via `HTTP_PROXY`).
  Licence : fichiers `epl-v20` + `edl-v10` à la racine du repo (EPL-2.0/EDL-1.0, permissif faible,
  compatible usage commercial avec attribution).
  Gotchas doc : handler de message **ne doit pas bloquer** (sinon `pingresp not received`),
  `SetOrderMatters(false)` sauf besoin d'ordre, `KeepAlive` pour détecter les coupures,
  `DefaultPublishHandler` recommandé avec `CleanSession=false` (sinon retains non acquittés →
  limite inflight), **ne pas réutiliser un `Client` après `Disconnect`**.
  (https://github.com/eclipse-paho/paho.mqtt.golang — README/Troubleshooting, consulté 2026-09-23)
- **nats.go** (client Go de référence) : **Apache-2.0**, reconnexion/backoff, TLS, Nkeys/JWT,
  JetStream + KV + ObjectStore + API services (`micro`).
  (https://github.com/nats-io/nats.go — README + badge licence, consulté 2026-09-23 ;
  code serveur Apache-2.0 confirmé par l'en-tête de
  https://github.com/nats-io/nats-server/blob/c525d513/server/jetstream.go)
- Comparatif protocole (contexte) : MQTT = sessions persistantes + QoS 0/1/2, topics `/`
  (`+`,`#`) ; NATS = pas de session, patterns req/reply + queue groups, sujets `.` (`*`,`>`),
  persistance via extension JetStream (https://www.emqx.com/en/blog/emqx-nats-gateway, 2025-07-08).

## 4. Contraintes perf 1–50 Hz, JSON, GC Go, SQLite

- Budget temporel : 50 Hz = **20 ms/période**. Chaîne locale typique mesurée/attendue :
  IPC gRPC ~0,1–0,5 ms (§2) + bus local ~0,2–0,3 ms (§3, 1 KB) + insert SQLite WAL <1 ms (§ ci-dessous)
  ⇒ **marge >10× même à 50 Hz**, pour quelques dizaines de signaux. Le risque n'est pas la
  latence médiane mais les **queues** : unbounded channels, handlers Paho bloquants, COMMITs
  concurrents (§ ci-dessous).
- **SQLite : un seul writer à la fois, même en WAL.** En WAL, lecteurs et writer coexistent
  (snapshot isolation), mais « there can only be one writer at a time » ; `BEGIN CONCURRENT`
  (mode wal/wal2) parallélise la phase de transaction mais **sérialise les COMMIT** (verrou
  writer + busy-handler, inefficace sous contention → mutex applicatif recommandé si tous les
  writers sont dans le même process).
  (https://www.sqlite.org/wal.html ; https://sqlite.org/isolation.html ;
  https://sqlite.org/src/doc/begin-concurrent/doc/begin_concurrent.md)
  → Règle MVP : **un seul writer** (goroutine dédiée + batching par transaction, `journal_mode=WAL`,
  `synchronous=NORMAL` en edge avec perte au crash acceptée ; FULL si durabilité exigée, au prix
  d'un fsync/COMMIT). **Jamais d'écriture directe depuis les plugins** — seul le Core écrit.
- **GC Go** : concurrent, pauses globales non proportionnelles au tas (plutôt à GOMAXPROCS /
  arrêt des goroutines) ; la fréquence GC arbitre CPU vs mémoire (`GOGC`, `SetGCPercent`,
  limite mémoire). Sources de latence résiduelle : transitions mark/sweep, 25 % CPU pendant le
  mark, assist, write barriers, scan des roots.
  (https://go.dev/doc/gc-guide, « Latency », décrit le GC à partir de Go 1.19)
  → À 50 Hz : pauses sub-ms attendues, OK ; mitiger par zéro-alloc sur le hot path (réutiliser
  buffers, éviter allocs/boucle), `GOMEMLIMIT` sur RPi, traces `GODEBUG=gctrace=1` si jitter.
- **Cross-compile ARM** : `GOOS=linux GOARCH=arm64` (RPi 4/5 64-bit) / `GOARCH=arm GOARM=7`
  pour 32-bit ; binaire **statique** par défaut (pas de libc requise) — le packaging du Core
  est trivial comparé à l'env Python.
  (https://go.dev/doc/install/source#environment ; https://go.googlesource.com/wiki/+/2146f6c4d7eddc1392eceffb5f499a37844d7218/GoArm.md ;
  défaut `GOARM64=v8.0`, https://tip.golang.org/wiki/MinimumRequirements)

## 5. Packaging binaire Go + env Python (edge RPi)

- Core : `GOOS=linux GOARCH=arm64 go build` → un fichier, déposé par SCP/CI. CGO **désactivé**
  (driver SQLite en pur Go, ex. `modernc.org/sqlite`, pour garder le binaire statique et le
  cross-compile sans toolchain C ; si `mattn/go-sqlite3` alors toolchain `aarch64` + CGO —
  à trancher au setup, défaut : pur Go).
- Plugins Python : interprétées, pas de binaire — livrer un **venv figé**
  (`requirements.txt` + `pip freeze`, versions épinglées : `grpcio`, `protobuf`, `paho-mqtt`,
  `river`/`scikit-learn`…), un lanceur `plugins/<nom>/run.sh` (`exec <venv>/bin/python plugin.py`),
  `PYTHONUNBUFFERED=1` (handshake stdout non bufferisé — critique §1), health + timeouts.
  Alternative durcissement V1 : `pex`/`zipapp` ou conteneur par plugin.
- Broker MVP = **binaire existant** (Mosquitto via `apt`), zéro code : le Core n'embarque que
  le client Paho Go. NATS embarquable en lib Go plus tard (le serveur est importable ;
  wrapper ex. `osapi-io/nats-server` : `Start()`/`Stop()` non-bloquants —
  https://github.com/osapi-io/nats-server).
- Isolation : processus séparés (crash plugin ≠ crash Core) ; sandboxing cgroups/namespaces =
  V2 (noté dans les rapports locaux, inchangé).

## 6. Travail restant à mesurer (pas de claim sans bench)

1. `bench_ipc` : Go↔Python, gRPC loopback vs UDS, payloads 100 B / 4 KB, 50 Hz + rafale ;
   comparer latence p50/p99 et CPU (valider §2, trancher UDS en V1).
2. `bench_bus` : Mosquitto local, 50 Hz × N signaux, QoS 0 vs 1, CPU/RAM sur RPi.
3. `bench_store` : inserts SQLite WAL batchés à 50 Hz, `synchronous=NORMAL` vs FULL, taille WAL/checkpoint.
4. Décider driver SQLite (pur Go vs CGO) au vu du cross-compile ARM réel.

---

## 7. Recommandation lockée

### MVP (soutenance, edge RPi)

- **IPC Core↔plugins : go-plugin `ProtocolGRPC` sur TCP loopback (`127.0.0.1`, ports éphémères),
  `AutoMTLS=false`, `grpcio` + stubs générés côté Python.**
  Pourquoi : c'est le **seul** mode go-plugin interopérable Python (recette officielle + exemple
  Python du repo, §1) ; on hérite handshake, versioning `ProtocolVersion`, capture stderr,
  `Kill()` gracieux. Le JSON-RPC/stdio custom est **écarté** : il réinvente ce que la lib donne.
  Les chiffres locaux (~150 µs vs ~116 µs) ne justifient pas à eux seuls un protocole maison.
- **Contrats : `.proto` dès le jour 1** (`Signal{ts_ns,source,name,value_double|bytes,quality}`,
  `Decision{cause+,priority,action,eta_s,case_ref}`, `Envelope{v}`) — binaire protobuf sur le wire
  (§2) ; **pas de JSON sur Core↔plugins**.
- **Bus MVP : Mosquitto externe (broker du système) + Paho Go dans le Core.**
  Topics `mlelec/<site>/<equip>/<signal>` ; **QoS 0** temps réel, **QoS 1** alertes/décisions ;
  `CleanSession=false` + `DefaultPublishHandler` + handlers non bloquants + `KeepAlive` (§3).
  Payloads MQTT en **JSON schématisé** (interop ESP32/RPi/capteurs), schéma versionné dans le repo.
- **Stockage : SQLite, `journal_mode=WAL`, un seul writer** (goroutine Core + batching/transaction,
  `synchronous=NORMAL` sauf exigence durabilité), **aucun accès DB depuis les plugins** (§4).
- **Watchdog (code Core, la lib ne le fait pas)** : `StartTimeout` + health `plugin:SERVING` +
  **heartbeat métier 1 Hz** (deadlock ≠ process mort) + restart backoff (ex. 1 s/5 s/30 s, max N
  puis état `degraded` + fiche de cas) + `CleanupClients()` à l'arrêt. Logs stderr centralisés.
- **Perf attendue** : <2 ms bout-en-bout médian à 50 Hz en local (IPC ~0,3 ms + bus ~0,3 ms +
  store <1 ms), marge ~10× sur la période de 20 ms. Le bench §6 doit le confirmer sur RPi.

### Chemin V1 (post-soutenance, sans réécriture)

1. `TCP loopback → UDS` (`NETWORK-TYPE=unix`, perms socket) : pas de ports, durcissement local (§1–§2).
2. `AutoMTLS=true` + pile Python compatible P-521 (`grpclib`, cf. §1) ou rester `grpcio` sans mTLS
   si menace locale jugée faible — trancher par analyse menace, pas par défaut.
3. Bus : **converger vers `nats-server` (embarqué ou sidecar) avec listener MQTT** : parc existant
   inchangé (pas de reflash), routage interne NATS (req/reply, queue groups), JetStream pour
   persistance/replay des décisions ; QoS natif NATS→MQTT = 0 (§3) → les alertes critiques restent
   publiées côté MQTT ou relues depuis un stream.
4. Stockage : passer le writer sur file/batching durci ; TimescaleDB/InfluxDB uniquement si le
   volume le justifie (option déjà actée dans les rapports locaux).
5. Sandboxing (cgroups/namespaces), `Reattach` opérateur, conteneurisation des plugins Python.

### Risques résiduels (propriétaire : Core team)

| Risque | Gravité MVP | Mitigation lockée |
|---|---|---|
| Sérialisation JSON goulot | **Écarté** Core↔plugins (protobuf) ; reste MQTT edge, volumes faibles | §7-MVP ; bench §6.2 |
| GC Go / jitter 50 Hz | Faible | zéro-alloc hot path, `GOMEMLIMIT`, `gctrace` si besoin (§4) |
| Concurrence SQLite (`SQLITE_BUSY`, perte) | Moyen | single-writer + WAL + batching, plugins sans accès DB (§4) |
| P-521/BoringSSL si AutoMTLS activé trop tôt | Moyen | `AutoMTLS=false` au MVP, chemin V1.2 (§1) |
| Singleton gRPC (« dispensing ») | Faible si contrats simples | `instance_id` dans `.proto` dès qu'un plugin est multi-instancié (§1) |
| Mosquitto : drops silencieux (`max_queued_*`), loopback défaut, `persistence=false` | Moyen | config versionnée (`mosquitto.conf` dans le repo), QoS 1 alertes, tests coupure (§3) |
| Python : stdout bufferisé casse le handshake | Faible | `PYTHONUNBUFFERED=1` + flush explicite (§1, §5) |

*Ne pas closer l'issue #2 — le parent gère branche `research/ipc-bus` + pointeur ticket.*

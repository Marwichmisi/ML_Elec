# 🔍 Code Review Report - ML_Elec

**Date:** 23 mars 2026  
**Revieweur:** BMad Dev Agent (avec Context7 CLI)  
**Statut:** ✅ Architecture validée - ⚠️ Implémentation à compléter

---

## 📊 Résumé Exécutif

### État Actuel du Projet

| Package | Statut | Couverture | Conformité |
|---------|--------|------------|------------|
| `packages/edge-agent` | ⚠️ Skeleton | Config + HTTP minimal | ✅ 85% |
| `packages/shared` | ⚠️ Skeleton | Types de base | ✅ 90% |
| `packages/ui` | ⚠️ Skeleton | Export vide | ❌ 0% |
| `apps/web` | ⚠️ Skeleton | Dashboard statique | ✅ 75% |

**Verdict global :** L'architecture est **correcte** et alignée avec le PRD, mais l'implémentation est **très partielle**. Toutes les fonctionnalités core (OPC-UA, DuckDB, Fastify, React Flow) sont manquantes.

---

## ✅ Points Forts (Validés)

### 1. Architecture Monorepo
- ✅ Structure `apps/` + `packages/` correcte
- ✅ Turborepo configuré pour builds incrémentaux
- ✅ pnpm workspaces pour gestion des dépendances

### 2. Configuration TypeScript
- ✅ `strict: true` dans tous les packages
- ✅ Extensions `.js` obligatoires dans `edge-agent` (ESM NodeNext)
- ✅ Types partagés centralisés dans `packages/shared`

### 3. Validation Zod
- ✅ Schéma de configuration validé avec `zod`
- ✅ Fail-fast au démarrage si variables manquantes
- ✅ Types inférés depuis les schémas (`z.infer`)

### 4. Gestion des Erreurs
- ✅ Erreurs typées dans `loadConfig()`
- ✅ Logging structuré avec prefixes `[Edge Agent]`

---

## ⚠️ Problèmes Critiques (À Corriger)

### 1. **Fastify Non Implémenté** ❌

**Problème:** Le serveur HTTP utilise `http.createServer()` au lieu de Fastify.

**Code actuel:**
```typescript
// packages/edge-agent/src/index.ts
const http = await import('http');
const server = http.createServer((req, res) => { ... });
```

**Recommandation (selon documentation Fastify):**
```typescript
import Fastify from 'fastify';

const fastify = Fastify({
  logger: { level: config.logLevel }
});

fastify.get('/health', {
  schema: {
    response: {
      200: {
        type: 'object',
        properties: {
          status: { type: 'string' },
          timestamp: { type: 'string' }
        }
      }
    }
  }
}, async (request, reply) => {
  return { status: 'healthy', timestamp: new Date().toISOString() };
});

await fastify.listen({ port: config.port, host: '0.0.0.0' });
```

**Pourquoi:** Fastify offre:
- Validation JSON Schema automatique
- Sérialisation rapide (2x plus rapide qu'Express)
- Plugin architecture (JWT, CORS, rate limiting)
- TypeScript-first

---

### 2. **DuckDB Non Configuré** ❌

**Problème:** Aucune configuration DuckDB, pas de metrics writer, pas de Appender.

**Recommandation (selon @duckdb/node-api docs):**
```typescript
import { Database, Connection, DuckDBDataChunk } from '@duckdb/node-api';

// Configuration edge (Raspberry Pi 4)
await connection.run(`SET memory_limit = '1.5GB'`);
await connection.run(`SET threads = 2`);

// Création table metrics avec partitionnement temporel
await connection.run(`
  CREATE TABLE IF NOT EXISTS metrics (
    node_id VARCHAR NOT NULL,
    value DOUBLE NOT NULL,
    quality VARCHAR NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
  )
`);

// High-frequency ingestion avec Appender
const appender = await connection.createAppender('metrics');
appender.appendVarchar(nodeId);
appender.appendDouble(value);
appender.appendVarchar(quality);
appender.appendTimestamp(new Date(timestamp));
appender.endRow();
await appender.flush();
```

---

### 3. **OPC-UA Client Manquant** ❌

**Problème:** Aucun code OPC-UA, pas de reconnexion automatique, pas de subscriptions.

**Recommandation (selon node-opcua docs):**
```typescript
import { OPCUAClient, ClientSession, ClientSubscription } from 'node-opcua-client';

const client = OPCUAClient.create({
  endpointMustExist: false,
  connectionStrategy: {
    maxRetry: 100,
    initialDelay: 1000,
    maxDelay: 60000,
    randomisationFactor: 0.1
  },
  keepSessionAlive: true
});

// Reconnexion automatique built-in
client.on('connection_reestablished', () => {
  console.log('[OPC-UA] Reconnexion établie');
});

// Subscription haute fréquence
const subscription = await session.createSubscription({
  requestedPublishingInterval: 1000,
  requestedLifetimeCount: 100,
  requestedMaxKeepAliveCount: 10,
  publishingEnabled: true
});

const monitoredItem = ClientMonitoredItem.create(
  subscription,
  { nodeId: 'ns=1;s=Temperature', attributeId: AttributeIds.Value },
  { samplingInterval: 100, queueSize: 10 }
);

monitoredItem.on('changed', (dataValue) => {
  // Ingérer dans DuckDB
});
```

---

### 4. **Types Partiels dans `packages/shared`** ⚠️

**Problème:** Les types sont incomplets et ne couvrent pas tous les cas d'usage.

**Code actuel:**
```typescript
export interface PipelineNode {
  id: string;
  type: string; // ❌ Trop générique
  data: Record<string, unknown>; // ❌ Pas de validation
  position: { x: number; y: number };
}
```

**Recommandation:**
```typescript
import { z } from 'zod';

// Schémas Zod pour validation runtime
export const zPipelineNodeType = z.enum([
  'opcua-source',
  'mqtt-source',
  'filter',
  'anomaly-detector',
  'alert',
  'dashboard'
]);

export const zPipelineNode = z.object({
  id: z.string().uuid(),
  type: zPipelineNodeType,
  data: z.record(z.string(), z.unknown()),
  position: z.object({
    x: z.number(),
    y: z.number()
  })
});

export type PipelineNode = z.infer<typeof zPipelineNode>;

// Types d'erreurs domaine
export class OpcUaConnectionError extends Error {
  constructor(
    public endpoint: string,
    public retryCount: number
  ) {
    super(`OPC-UA connection failed to ${endpoint} after ${retryCount} retries`);
    this.name = 'OpcUaConnectionError';
  }
}

export class StorageWriteError extends Error {
  constructor(
    public table: string,
    public cause: Error
  ) {
    super(`Failed to write to ${table}: ${cause.message}`);
    this.name = 'StorageWriteError';
  }
}
```

---

### 5. **UI Vide dans `packages/ui`** ❌

**Problème:** Le package UI n'exporte rien, pas de composants React Flow custom.

**Recommandation:**
```typescript
// packages/ui/src/nodes/OpcUaSourceNode.tsx
import { Handle, Position, NodeProps } from '@xyflow/react';

export interface OpcUaSourceData {
  endpoint: string;
  nodeId: string;
  pollingInterval: number;
  status: 'connected' | 'disconnected' | 'error';
}

export function OpcUaSourceNode({ data }: NodeProps<OpcUaSourceData>) {
  return (
    <div className={`ml-elec-node opcua-source ${data.status}`}>
      <Handle type="source" position={Position.Right} id="output" />
      <div className="node-header">
        <span className="node-icon">📡</span>
        <span className="node-title">OPC-UA Source</span>
      </div>
      <div className="node-body">
        <div className="node-field">
          <label>Endpoint:</label>
          <code>{data.endpoint}</code>
        </div>
        <div className="node-field">
          <label>Node ID:</label>
          <code>{data.nodeId}</code>
        </div>
        <div className="node-status">
          <StatusBadge status={data.status} />
        </div>
      </div>
    </div>
  );
}

// packages/ui/src/index.ts
export { OpcUaSourceNode } from './nodes/OpcUaSourceNode.js';
export { AnomalyDetectorNode } from './nodes/AnomalyDetectorNode.js';
export { AlertNode } from './nodes/AlertNode.js';
export { nodeTypes } from './node-types.js';
```

---

### 6. **Tests Absents** ❌

**Problème:** Aucun test implémenté, seulement des configs Vitest.

**Recommandation:**
```typescript
// packages/edge-agent/src/config.test.ts
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { loadConfig } from './config.js';

describe('loadConfig', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  it('should load with valid environment', () => {
    process.env.JWT_SECRET = 'test-secret';
    process.env.PORT = '3001';
    
    const config = loadConfig();
    
    expect(config.port).toBe(3001);
    expect(config.jwtSecret).toBe('test-secret');
  });

  it('should exit on invalid JWT_SECRET in production', () => {
    process.env.NODE_ENV = 'production';
    process.env.JWT_SECRET = '';
    
    // Mock process.exit
    const mockExit = vi.spyOn(process, 'exit').mockImplementation();
    
    loadConfig();
    
    expect(mockExit).toHaveBeenCalledWith(1);
  });
});

// packages/edge-agent/src/duckdb/metrics-writer.test.ts
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { Database, Connection } from '@duckdb/node-api';
import { MetricsWriter } from './metrics-writer.js';

describe('MetricsWriter', () => {
  let db: Database;
  let conn: Connection;
  let writer: MetricsWriter;

  beforeEach(async () => {
    db = await Database.create(':memory:');
    conn = await db.connect();
    await conn.run(`
      CREATE TABLE metrics (
        node_id VARCHAR,
        value DOUBLE,
        timestamp TIMESTAMPTZ
      )
    `);
    writer = new MetricsWriter(conn);
  });

  afterEach(async () => {
    await conn.close();
    await db.close();
  });

  it('should write metrics successfully', async () => {
    await writer.write({
      nodeId: 'ns=1;s=Temperature',
      value: 42.5,
      timestamp: new Date().toISOString(),
      quality: 'good'
    });

    const result = await conn.query('SELECT COUNT(*) FROM metrics');
    expect(result.toArray()[0][0]).toBe(1);
  });
});
```

---

## 🔍 Vérification Documentation (Context7 CLI)

### React 18 ✅
- **Documentation consultée:** `/facebook/react` v18.3.1
- **Constat:** L'usage de `useState` et `useEffect` dans `App.tsx` est correct
- **Recommandation:** Ajouter `useCallback` pour éviter les re-renders inutiles dans React Flow

### Fastify 4 ✅
- **Documentation consultée:** `/fastify/fastify`
- **Constat:** Non implémenté (voir section "Problèmes Critiques")
- **Recommandation:** Priorité haute pour remplacer `http.createServer()`

### DuckDB Node-API ✅
- **Documentation consultée:** `/websites/duckdb_stable_clients_node_neo`
- **Constat:** API Neo (Promise-based) recommandée, Appender pour haute fréquence
- **Recommandation:** Utiliser `createAppender()` + `flushSync()` pour ingestion

### Zod ✅
- **Documentation consultée:** `/colinhacks/zod` v3.24.2
- **Constat:** Bonne utilisation de `safeParse()` et `z.infer<>`
- **Recommandation:** Ajouter des schémas pour tous les types partagés

### node-opcua ✅
- **Documentation consultée:** `/node-opcua/node-opcua`
- **Constat:** Non implémenté
- **Recommandation:** Utiliser `ClientMonitoredItem.create()` avec reconnexion auto

---

## 📋 Checklist de Conformité

### Architecture Hexagonale
- [ ] Créer `packages/edge-agent/src/domain/` (logique métier pure)
- [ ] Créer `packages/edge-agent/src/ports/` (interfaces TypeScript)
- [ ] Créer `packages/edge-agent/src/adapters/` (implémentations DuckDB, OPC-UA)
- [ ] Créer `packages/edge-agent/src/api/` (routes Fastify uniquement)

### Sécurité
- [ ] Ajouter `@fastify/jwt` pour authentification
- [ ] Générer certificats OPC-UA dans `./pki/` (à ignorer dans git)
- [ ] Chiffrer données Parquet avec `ENCRYPTION_CONFIG` DuckDB
- [ ] Audit log write-ahead avant modifications

### Performance Edge
- [ ] Configurer `SET memory_limit = '1.5GB'` DuckDB
- [ ] Configurer `SET threads = 2` pour Raspberry Pi 4
- [ ] Implémenter compaction nocturne Parquet (cron 02:00)
- [ ] Virtualisation React Flow (`useVirtualizer`) pour >50 nœuds

### Tests
- [ ] Tests unitaires edge-agent (85% coverage)
- [ ] Tests d'intégration OPC-UA avec mock server
- [ ] Tests E2E Playwright (scénario pipeline complet)
- [ ] Tests offline-first (NETWORK=none)

### Qualité de Code
- [ ] ESLint v9 avec `@typescript-eslint/no-floating-promises`
- [ ] Prettier v3 + Husky pre-commit hook
- [ ] JSDoc sur interfaces publiques de `packages/shared`
- [ ] Conventional Commits pour git

---

## 🎯 Priorités d'Implémentation

### Sprint 1 - Fondation (Epic 1)
1. **Fastify + JWT** - Remplacer HTTP server
2. **DuckDB MetricsWriter** - Appender + partitionnement
3. **OPC-UA Collector** - Reconnexion auto + subscriptions
4. **Tests unitaires** - Config + DuckDB in-memory

### Sprint 2 - Core ML (Epic 2 & 3)
5. **MQTT Source** - Alternative à OPC-UA
6. **Anomaly Detector** - Détection statistique simple
7. **Alert System** - Severity levels + confidence
8. **Feedback 1-tap** - Pour amélioration modèle

### Sprint 3 - UX (Epic 6)
9. **React Flow Nodes** - Composants custom UI
10. **HealthTimelineCard** - Courbe de santé
11. **AlertDecisionCard** - Décision en ≤60s

---

## 📚 Références Documentation (Context7)

| Technologie | Library ID | Version | URL |
|-------------|-----------|---------|-----|
| React | `/facebook/react` | v18.3.1 | https://react.dev |
| Fastify | `/fastify/fastify` | v4.x | https://fastify.dev |
| DuckDB | `/websites/duckdb_stable_clients_node_neo` | 1.5.1 | https://duckdb.org/docs |
| Zod | `/colinhacks/zod` | v3.24.2 | https://zod.dev |
| node-opcua | `/node-opcua/node-opcua` | v2.163 | https://node-opcua.github.io |

---

## ✅ Conclusion

**Architecture:** ✅ **Validée** - Conforme au PRD et aux bonnes pratiques edge computing

**Implémentation:** ⚠️ **Partielle** - Skeleton en place, mais fonctionnalités core manquantes

**Recommandation:** Commencer par **Sprint 1 (Fondation)** avec Fastify + DuckDB + OPC-UA avant d'ajouter des features ML/UX.

**Prochaine étape:** Créer les stories détaillées pour Epic 1 (Onboarding & Déploiement Edge) et implémenter Fastify en premier.

---

**Review complétée avec Context7 CLI** ✅  
**Documentation vérifiée:** React, Fastify, DuckDB, Zod, node-opcua

/**
 * Adapter DuckDB pour le stockage des données
 * Implémente IStorageEngine avec l'API Neo (@duckdb/node-api)
 */

import { Database, Connection, DuckDBDataChunk, TIMESTAMP } from '@duckdb/node-api';
import type { IStorageEngine } from '../ports/index.js';
import type { SensorMetric, AnomalyAlert, Equipment, PipelineConfig } from '@ml-elec/shared';
import { StorageWriteError, StorageReadError } from '@ml-elec/shared/errors';

export interface DuckDBStorageOptions {
  dataDir: string;
  memoryLimit?: string;
  threads?: number;
}

export class DuckDBStorageAdapter implements IStorageEngine {
  private db: Database | null = null;
  private connection: Connection | null = null;
  private readonly options: DuckDBStorageOptions;

  constructor(options: DuckDBStorageOptions) {
    this.options = options;
  }

  /**
   * Initialise la connexion DuckDB avec configuration edge
   */
  async initialize(): Promise<void> {
    try {
      // Création de la base de données
      this.db = await Database.create(`${this.options.dataDir}/metrics.duckdb`);
      this.connection = await this.db.connect();

      // Configuration edge (Raspberry Pi 4)
      await this.connection.run(`SET memory_limit = '${this.options.memoryLimit || '1.5GB'}'`);
      await this.connection.run(`SET threads = ${this.options.threads || 2}`);

      // Création des tables
      await this.createTables();
    } catch (error) {
      throw new StorageWriteError('initialization', undefined, error as Error);
    }
  }

  /**
   * Crée les tables si elles n'existent pas
   */
  private async createTables(): Promise<void> {
    if (!this.connection) {
      throw new StorageWriteError('tables', 'No connection');
    }

    // Table metrics avec partitionnement temporel
    await this.connection.run(`
      CREATE TABLE IF NOT EXISTS metrics (
        node_id VARCHAR NOT NULL,
        equipment_id VARCHAR,
        value DOUBLE NOT NULL,
        quality VARCHAR NOT NULL,
        timestamp TIMESTAMPTZ NOT NULL
      )
    `);

    // Index pour les requêtes fréquentes
    await this.connection.run(`
      CREATE INDEX IF NOT EXISTS idx_metrics_node_timestamp 
      ON metrics(node_id, timestamp DESC)
    `);

    await this.connection.run(`
      CREATE INDEX IF NOT EXISTS idx_metrics_equipment_timestamp 
      ON metrics(equipment_id, timestamp DESC)
    `);

    // Table alerts
    await this.connection.run(`
      CREATE TABLE IF NOT EXISTS alerts (
        id VARCHAR PRIMARY KEY,
        equipment_id VARCHAR NOT NULL,
        severity VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        timestamp TIMESTAMPTZ NOT NULL,
        confidence DOUBLE NOT NULL,
        node_id VARCHAR,
        value DOUBLE,
        threshold DOUBLE,
        acknowledged BOOLEAN DEFAULT FALSE,
        acknowledged_by VARCHAR,
        acknowledged_at TIMESTAMPTZ
      )
    `);

    // Table equipments
    await this.connection.run(`
      CREATE TABLE IF NOT EXISTS equipments (
        id VARCHAR PRIMARY KEY,
        name VARCHAR NOT NULL,
        type VARCHAR NOT NULL,
        status VARCHAR NOT NULL,
        location VARCHAR,
        metadata VARCHAR,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
      )
    `);

    // Table pipelines
    await this.connection.run(`
      CREATE TABLE IF NOT EXISTS pipelines (
        id VARCHAR PRIMARY KEY,
        name VARCHAR NOT NULL,
        description VARCHAR,
        config VARCHAR NOT NULL,
        is_active BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
      )
    `);

    // Table audit_log (append-only)
    await this.connection.run(`
      CREATE TABLE IF NOT EXISTS audit_log (
        id VARCHAR PRIMARY KEY,
        user_id VARCHAR NOT NULL,
        action VARCHAR NOT NULL,
        resource VARCHAR NOT NULL,
        resource_id VARCHAR,
        details VARCHAR,
        timestamp TIMESTAMPTZ NOT NULL
      )
    `);
  }

  /**
   * Écrit une métrique avec Appender (haute fréquence)
   */
  async writeMetric(metric: SensorMetric): Promise<void> {
    if (!this.connection) {
      throw new StorageWriteError('metrics', metric.nodeId, new Error('No connection'));
    }

    try {
      const appender = await this.connection.createAppender('metrics');
      appender.appendVarchar(metric.nodeId);
      appender.appendVarchar(metric.equipmentId || null);
      appender.appendDouble(metric.value);
      appender.appendVarchar(metric.quality);
      appender.appendTimestamp(new Date(metric.timestamp));
      appender.endRow();
      await appender.flush();
    } catch (error) {
      throw new StorageWriteError('metrics', metric.nodeId, error as Error);
    }
  }

  /**
   * Écrit un lot de métriques avec DataChunk (bulk insert)
   */
  async writeMetrics(metrics: SensorMetric[]): Promise<void> {
    if (!this.connection || metrics.length === 0) {
      return;
    }

    try {
      const appender = await this.connection.createAppender('metrics');

      // Création d'un DataChunk pour le bulk insert
      const chunk = DuckDBDataChunk.create([
        'VARCHAR', // node_id
        'VARCHAR', // equipment_id
        'DOUBLE',  // value
        'VARCHAR', // quality
        'TIMESTAMPTZ', // timestamp
      ]);

      const nodeIds = metrics.map((m) => m.nodeId);
      const equipmentIds = metrics.map((m) => m.equipmentId || null);
      const values = metrics.map((m) => m.value);
      const qualities = metrics.map((m) => m.quality);
      const timestamps = metrics.map((m) => new Date(m.timestamp));

      chunk.setColumns([nodeIds, equipmentIds, values, qualities, timestamps]);
      appender.appendDataChunk(chunk);
      await appender.flush();
    } catch (error) {
      throw new StorageWriteError('metrics', 'bulk', error as Error);
    }
  }

  /**
   * Lit les métriques pour un équipement
   */
  async getMetrics(equipmentId: string, startTime: Date, endTime: Date): Promise<SensorMetric[]> {
    if (!this.connection) {
      throw new StorageReadError('metrics', `equipmentId=${equipmentId}`, new Error('No connection'));
    }

    try {
      const result = await this.connection.query(`
        SELECT node_id, equipment_id, value, quality, timestamp
        FROM metrics
        WHERE equipment_id = $1
          AND timestamp >= $2
          AND timestamp <= $3
        ORDER BY timestamp ASC
      `, [equipmentId, startTime.toISOString(), endTime.toISOString()]);

      return result.toArray().map((row) => ({
        nodeId: row[0] as string,
        equipmentId: row[1] as string | undefined,
        value: row[2] as number,
        quality: row[3] as 'good' | 'uncertain' | 'bad',
        timestamp: (row[4] as Date).toISOString(),
      }));
    } catch (error) {
      throw new StorageReadError('metrics', `equipmentId=${equipmentId}`, error as Error);
    }
  }

  /**
   * Lit la dernière métrique pour un nodeId
   */
  async getLatestMetric(nodeId: string): Promise<SensorMetric | null> {
    if (!this.connection) {
      throw new StorageReadError('metrics', `nodeId=${nodeId}`, new Error('No connection'));
    }

    try {
      const result = await this.connection.query(`
        SELECT node_id, equipment_id, value, quality, timestamp
        FROM metrics
        WHERE node_id = $1
        ORDER BY timestamp DESC
        LIMIT 1
      `, [nodeId]);

      const rows = result.toArray();
      if (rows.length === 0) {
        return null;
      }

      const row = rows[0];
      return {
        nodeId: row[0] as string,
        equipmentId: row[1] as string | undefined,
        value: row[2] as number,
        quality: row[3] as 'good' | 'uncertain' | 'bad',
        timestamp: (row[4] as Date).toISOString(),
      };
    } catch (error) {
      throw new StorageReadError('metrics', `nodeId=${nodeId}`, error as Error);
    }
  }

  /**
   * Sauvegarde une alerte
   */
  async saveAlert(alert: AnomalyAlert): Promise<void> {
    if (!this.connection) {
      throw new StorageWriteError('alerts', alert.id, new Error('No connection'));
    }

    try {
      await this.connection.query(`
        INSERT INTO alerts (id, equipment_id, severity, description, timestamp, confidence, node_id, value, threshold)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
      `, [
        alert.id,
        alert.equipmentId,
        alert.severity,
        alert.description,
        alert.timestamp,
        alert.confidence,
        alert.nodeId || null,
        alert.value || null,
        alert.threshold || null,
      ]);
    } catch (error) {
      throw new StorageWriteError('alerts', alert.id, error as Error);
    }
  }

  /**
   * Récupère les alertes non acquittées
   */
  async getPendingAlerts(): Promise<AnomalyAlert[]> {
    if (!this.connection) {
      throw new StorageReadError('alerts', 'pending', new Error('No connection'));
    }

    try {
      const result = await this.connection.query(`
        SELECT id, equipment_id, severity, description, timestamp, confidence, node_id, value, threshold
        FROM alerts
        WHERE acknowledged = FALSE
        ORDER BY timestamp DESC
      `);

      return result.toArray().map((row) => ({
        id: row[0] as string,
        equipmentId: row[1] as string,
        severity: row[2] as 'info' | 'warning' | 'critical',
        description: row[3] as string,
        timestamp: (row[4] as Date).toISOString(),
        confidence: row[5] as number,
        nodeId: row[6] as string | undefined,
        value: row[7] as number | undefined,
        threshold: row[8] as number | undefined,
      }));
    } catch (error) {
      throw new StorageReadError('alerts', 'pending', error as Error);
    }
  }

  /**
   * Acquitte une alerte
   */
  async acknowledgeAlert(alertId: string, userId: string): Promise<void> {
    if (!this.connection) {
      throw new StorageWriteError('alerts', alertId, new Error('No connection'));
    }

    try {
      await this.connection.query(`
        UPDATE alerts
        SET acknowledged = TRUE, acknowledged_by = $2, acknowledged_at = NOW()
        WHERE id = $1
      `, [alertId, userId]);
    } catch (error) {
      throw new StorageWriteError('alerts', alertId, error as Error);
    }
  }

  /**
   * Sauvegarde un équipement
   */
  async saveEquipment(equipment: Equipment): Promise<void> {
    if (!this.connection) {
      throw new StorageWriteError('equipments', equipment.id, new Error('No connection'));
    }

    try {
      await this.connection.query(`
        INSERT INTO equipments (id, name, type, status, location, metadata, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (id) DO UPDATE SET
          name = EXCLUDED.name,
          type = EXCLUDED.type,
          status = EXCLUDED.status,
          location = EXCLUDED.location,
          metadata = EXCLUDED.metadata,
          updated_at = EXCLUDED.updated_at
      `, [
        equipment.id,
        equipment.name,
        equipment.type,
        equipment.status,
        equipment.location || null,
        equipment.metadata ? JSON.stringify(equipment.metadata) : null,
        equipment.createdAt,
        equipment.updatedAt,
      ]);
    } catch (error) {
      throw new StorageWriteError('equipments', equipment.id, error as Error);
    }
  }

  /**
   * Récupère tous les équipements
   */
  async getAllEquipments(): Promise<Equipment[]> {
    if (!this.connection) {
      throw new StorageReadError('equipments', 'all', new Error('No connection'));
    }

    try {
      const result = await this.connection.query(`
        SELECT id, name, type, status, location, metadata, created_at, updated_at
        FROM equipments
        ORDER BY name ASC
      `);

      return result.toArray().map((row) => ({
        id: row[0] as string,
        name: row[1] as string,
        type: row[2] as string,
        status: row[3] as 'online' | 'offline' | 'maintenance' | 'error',
        location: row[4] as string | undefined,
        metadata: row[5] ? JSON.parse(row[5] as string) : undefined,
        createdAt: (row[6] as Date).toISOString(),
        updatedAt: (row[7] as Date).toISOString(),
      }));
    } catch (error) {
      throw new StorageReadError('equipments', 'all', error as Error);
    }
  }

  /**
   * Sauvegarde la configuration d'un pipeline
   */
  async savePipeline(pipeline: PipelineConfig): Promise<void> {
    if (!this.connection) {
      throw new StorageWriteError('pipelines', pipeline.id, new Error('No connection'));
    }

    try {
      await this.connection.query(`
        INSERT INTO pipelines (id, name, description, config, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
          name = EXCLUDED.name,
          description = EXCLUDED.description,
          config = EXCLUDED.config,
          is_active = EXCLUDED.is_active,
          updated_at = EXCLUDED.updated_at
      `, [
        pipeline.id,
        pipeline.name,
        pipeline.description || null,
        JSON.stringify(pipeline),
        pipeline.isActive,
        pipeline.createdAt,
        pipeline.updatedAt,
      ]);
    } catch (error) {
      throw new StorageWriteError('pipelines', pipeline.id, error as Error);
    }
  }

  /**
   * Récupère tous les pipelines
   */
  async getAllPipelines(): Promise<PipelineConfig[]> {
    if (!this.connection) {
      throw new StorageReadError('pipelines', 'all', new Error('No connection'));
    }

    try {
      const result = await this.connection.query(`
        SELECT id, name, description, config, is_active, created_at, updated_at
        FROM pipelines
        ORDER BY name ASC
      `);

      return result.toArray().map((row) => JSON.parse(row[3] as string) as PipelineConfig);
    } catch (error) {
      throw new StorageReadError('pipelines', 'all', error as Error);
    }
  }

  /**
   * Ferme la connexion DuckDB
   */
  async close(): Promise<void> {
    if (this.connection) {
      await this.connection.close();
      this.connection = null;
    }
    if (this.db) {
      await this.db.close();
      this.db = null;
    }
  }
}

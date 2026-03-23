/**
 * Tests unitaires pour DuckDB Storage Adapter
 * Utilise une base DuckDB in-memory pour les tests
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { DuckDBStorageAdapter } from './duckdb-storage.adapter.js';

describe('DuckDBStorageAdapter', () => {
  let storage: DuckDBStorageAdapter;

  beforeEach(async () => {
    storage = new DuckDBStorageAdapter({
      dataDir: ':memory:',
      memoryLimit: '512MB',
      threads: 1,
    });
    await storage.initialize();
  });

  afterEach(async () => {
    await storage.close();
  });

  describe('writeMetric', () => {
    it('should write a single metric successfully', async () => {
      const metric = {
        nodeId: 'ns=1;s=Temperature',
        value: 42.5,
        timestamp: new Date().toISOString(),
        quality: 'good' as const,
        equipmentId: '550e8400-e29b-41d4-a716-446655440000',
      };

      await expect(storage.writeMetric(metric)).resolves.not.toThrow();
    });

    it('should write multiple metrics in bulk', async () => {
      const metrics = [
        {
          nodeId: 'ns=1;s=Temperature',
          value: 42.5,
          timestamp: new Date().toISOString(),
          quality: 'good' as const,
        },
        {
          nodeId: 'ns=1;s=Pressure',
          value: 1013.25,
          timestamp: new Date().toISOString(),
          quality: 'good' as const,
        },
      ];

      await expect(storage.writeMetrics(metrics)).resolves.not.toThrow();
    });
  });

  describe('getMetrics', () => {
    it('should return empty array when no metrics exist', async () => {
      const equipmentId = '550e8400-e29b-41d4-a716-446655440000';
      const startTime = new Date('2024-01-01');
      const endTime = new Date();

      const metrics = await storage.getMetrics(equipmentId, startTime, endTime);

      expect(metrics).toEqual([]);
    });

    it('should return metrics for equipment', async () => {
      const equipmentId = '550e8400-e29b-41d4-a716-446655440000';
      const now = new Date().toISOString();

      // Write metrics
      await storage.writeMetric({
        nodeId: 'ns=1;s=Temperature',
        value: 42.5,
        timestamp: now,
        quality: 'good',
        equipmentId,
      });

      await storage.writeMetric({
        nodeId: 'ns=1;s=Pressure',
        value: 1013.25,
        timestamp: now,
        quality: 'good',
        equipmentId,
      });

      // Read metrics
      const startTime = new Date('2024-01-01');
      const endTime = new Date('2030-12-31');
      const metrics = await storage.getMetrics(equipmentId, startTime, endTime);

      expect(metrics).toHaveLength(2);
      expect(metrics.map((m) => m.nodeId)).toContain('ns=1;s=Temperature');
      expect(metrics.map((m) => m.nodeId)).toContain('ns=1;s=Pressure');
    });
  });

  describe('getLatestMetric', () => {
    it('should return null when no metric exists', async () => {
      const result = await storage.getLatestMetric('ns=1;s=NonExistent');
      expect(result).toBeNull();
    });

    it('should return the latest metric for a node', async () => {
      const nodeId = 'ns=1;s=Temperature';
      const baseTime = new Date('2024-01-01T12:00:00Z');

      // Write multiple metrics with different timestamps
      for (let i = 0; i < 5; i++) {
        const timestamp = new Date(baseTime.getTime() + i * 1000).toISOString();
        await storage.writeMetric({
          nodeId,
          value: 40 + i,
          timestamp,
          quality: 'good',
        });
      }

      const latest = await storage.getLatestMetric(nodeId);

      expect(latest).not.toBeNull();
      expect(latest!.value).toBe(44); // Last value
    });
  });

  describe('saveAlert', () => {
    it('should save an alert successfully', async () => {
      const alert = {
        id: 'alert-001',
        equipmentId: '550e8400-e29b-41d4-a716-446655440000',
        severity: 'warning' as const,
        description: 'Temperature too high',
        timestamp: new Date().toISOString(),
        confidence: 0.95,
        nodeId: 'ns=1;s=Temperature',
        value: 85.5,
        threshold: 80.0,
      };

      await expect(storage.saveAlert(alert)).resolves.not.toThrow();
    });
  });

  describe('getPendingAlerts', () => {
    it('should return only unacknowledged alerts', async () => {
      const alert1 = {
        id: 'alert-001',
        equipmentId: '550e8400-e29b-41d4-a716-446655440000',
        severity: 'warning' as const,
        description: 'Alert 1',
        timestamp: new Date().toISOString(),
        confidence: 0.9,
      };

      const alert2 = {
        id: 'alert-002',
        equipmentId: '550e8400-e29b-41d4-a716-446655440000',
        severity: 'critical' as const,
        description: 'Alert 2',
        timestamp: new Date().toISOString(),
        confidence: 0.99,
      };

      await storage.saveAlert(alert1);
      await storage.saveAlert(alert2);

      const pending = await storage.getPendingAlerts();

      expect(pending).toHaveLength(2);
    });
  });

  describe('saveEquipment', () => {
    it('should save an equipment successfully', async () => {
      const equipment = {
        id: '550e8400-e29b-41d4-a716-446655440000',
        name: 'Motor A1',
        type: 'motor',
        status: 'online' as const,
        location: 'Building A',
        metadata: { serial: '12345' },
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      await expect(storage.saveEquipment(equipment)).resolves.not.toThrow();
    });

    it('should update an existing equipment', async () => {
      const equipment = {
        id: '550e8400-e29b-41d4-a716-446655440000',
        name: 'Motor A1',
        type: 'motor',
        status: 'online' as const,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      await storage.saveEquipment(equipment);

      // Update
      const updated = {
        ...equipment,
        status: 'maintenance' as const,
        updatedAt: new Date().toISOString(),
      };

      await storage.saveEquipment(updated);

      const equipments = await storage.getAllEquipments();
      const found = equipments.find((e) => e.id === equipment.id);

      expect(found?.status).toBe('maintenance');
    });
  });

  describe('getAllEquipments', () => {
    it('should return empty array when no equipments', async () => {
      const equipments = await storage.getAllEquipments();
      expect(equipments).toEqual([]);
    });

    it('should return all equipments sorted by name', async () => {
      await storage.saveEquipment({
        id: '550e8400-e29b-41d4-a716-446655440002',
        name: 'Motor B',
        type: 'motor',
        status: 'online' as const,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      });

      await storage.saveEquipment({
        id: '550e8400-e29b-41d4-a716-446655440001',
        name: 'Motor A',
        type: 'motor',
        status: 'online' as const,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      });

      const equipments = await storage.getAllEquipments();

      expect(equipments).toHaveLength(2);
      expect(equipments[0].name).toBe('Motor A'); // Sorted
      expect(equipments[1].name).toBe('Motor B');
    });
  });
});

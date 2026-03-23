/**
 * Tests unitaires pour la configuration
 * Basé sur Vitest documentation officielle
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';

describe('loadConfig', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    vi.resetModules();
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  it('should load with valid environment', async () => {
    // Setup environment variables
    process.env.NODE_ENV = 'test';
    process.env.PORT = '3001';
    process.env.JWT_SECRET = 'test-secret-key-minimum-32-characters-long';
    process.env.DUCKDB_DATA_DIR = './test-data/duckdb';
    process.env.MQTT_BROKER_URL = 'mqtt://localhost:1883';
    process.env.OPCUA_ENDPOINT = 'opc.tcp://localhost:4840';
    process.env.LOG_LEVEL = 'debug';

    // Import after setting env
    const { loadConfig } = await import('../config.js');
    const config = loadConfig();

    expect(config.port).toBe(3001);
    expect(config.jwtSecret).toBe('test-secret-key-minimum-32-characters-long');
    expect(config.nodeEnv).toBe('test');
    expect(config.logLevel).toBe('debug');
  });

  it('should use default values when not provided', async () => {
    process.env.JWT_SECRET = 'test-secret-key-minimum-32-characters-long';

    const { loadConfig } = await import('../config.js');
    const config = loadConfig();

    expect(config.port).toBe(3001); // default
    expect(config.nodeEnv).toBe('development'); // default
    expect(config.logLevel).toBe('info'); // default
    expect(config.duckDbDataDir).toBe('./data/duckdb'); // default
  });

  it('should exit on invalid JWT_SECRET in production', async () => {
    process.env.NODE_ENV = 'production';
    process.env.JWT_SECRET = 'short';
    process.env.PORT = '3001';

    const mockExit = vi.spyOn(process, 'exit').mockImplementation((() => {}) as any);

    // Import and call - should exit
    await import('../config.js').catch(() => {
      // Expected to fail
    });

    expect(mockExit).toHaveBeenCalledWith(1);
  });

  it('should validate PORT as number', async () => {
    process.env.JWT_SECRET = 'test-secret-key-minimum-32-characters-long';
    process.env.PORT = '8080';

    const { loadConfig } = await import('../config.js');
    const config = loadConfig();

    expect(config.port).toBe(8080);
  });

  it('should validate MQTT_BROKER_URL as URL', async () => {
    process.env.JWT_SECRET = 'test-secret-key-minimum-32-characters-long';
    process.env.MQTT_BROKER_URL = 'invalid-url';

    const mockExit = vi.spyOn(process, 'exit').mockImplementation((() => {}) as any);

    await import('../config.js').catch(() => {
      // Expected to fail
    });

    expect(mockExit).toHaveBeenCalledWith(1);
  });

  it('should accept valid MQTT URL', async () => {
    process.env.JWT_SECRET = 'test-secret-key-minimum-32-characters-long';
    process.env.MQTT_BROKER_URL = 'mqtt://broker.example.com:1883';

    const { loadConfig } = await import('../config.js');
    const config = loadConfig();

    expect(config.mqttBrokerUrl).toBe('mqtt://broker.example.com:1883');
  });
});

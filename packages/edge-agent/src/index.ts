/**
 * ML_Elec Edge Agent - API REST pour maintenance prédictive
 * Basé sur Fastify avec plugins: JWT, CORS, Helmet, Validation
 * 
 * Architecture inspirée de: /fastify/fastify documentation officielle
 */

import Fastify, { FastifyInstance, FastifyServerOptions } from 'fastify';
import cors from '@fastify/cors';
import helmet from '@fastify/helmet';
import jwt from '@fastify/jwt';
import multipart from '@fastify/multipart';
import { config, loadConfig } from './config.js';
import { registerHealthRoutes } from './routes/health.routes.js';
import { registerAuthRoutes } from './routes/auth.routes.js';
import { registerMetricsRoutes } from './routes/metrics.routes.js';
import { registerEquipmentsRoutes } from './routes/equipments.routes.js';
import { registerPipelinesRoutes } from './routes/pipelines.routes.js';
import { DuckDBStorageAdapter } from './adapters/duckdb-storage.adapter.js';
import { OpcUaCollectorAdapter } from './adapters/opcua-collector.adapter.js';
import { MqttCollectorAdapter } from './adapters/mqtt-collector.adapter.js';

export interface AppBuildOptions {
  logger?: FastifyServerOptions['logger'];
}

/**
 * Crée l'instance Fastify avec tous les plugins
 */
export async function buildApp(options: AppBuildOptions = {}): Promise<FastifyInstance> {
  const fastify = Fastify({
    logger: options.logger ?? {
      level: config.logLevel,
      transport: config.nodeEnv === 'development' ? {
        target: 'pino-pretty',
        options: {
          translateTime: 'HH:MM:ss Z',
          ignore: 'pid,hostname',
        },
      } : undefined,
    },
  });

  // ============================================================================
  // PLUGINS GLOBAUX
  // ============================================================================

  // CORS - Configuration pour accès local + config utilisateur
  await fastify.register(cors, {
    origin: config.nodeEnv === 'production' ? false : true, // Production: false = même origine
    credentials: true,
    allowedHeaders: ['Content-Type', 'Authorization'],
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'],
  });

  // Helmet - Sécurité headers HTTP
  await fastify.register(helmet, {
    contentSecurityPolicy: config.nodeEnv === 'production',
    crossOriginEmbedderPolicy: false, // Nécessaire pour React Flow
  });

  // JWT - Authentification
  await fastify.register(jwt, {
    secret: config.jwtSecret,
    sign: {
      expiresIn: config.jwtExpiresIn,
    },
  });

  // Multipart - Upload de fichiers (CSV, configs)
  await fastify.register(multipart, {
    limits: {
      fileSize: 10 * 1024 * 1024, // 10MB max
    },
  });

  // ============================================================================
  // INITIALISATION STORAGE
  // ============================================================================

  const storage = new DuckDBStorageAdapter({
    dataDir: config.duckDbDataDir,
    memoryLimit: config.memoryLimit,
    threads: config.threads,
  });

  await storage.initialize();
  fastify.decorate('storage', storage);

  // ============================================================================
  // INITIALISATION COLLECTORS
  // ============================================================================

  const opcUaCollector = new OpcUaCollectorAdapter({
    endpoint: config.opcUaEndpoint,
    onMetric: async (metric) => {
      // Callback pour ingestion automatique
      await storage.writeMetric(metric);
    },
  });

  const mqttCollector = new MqttCollectorAdapter({
    brokerUrl: config.mqttBrokerUrl,
    onMetric: async (metric) => {
      await storage.writeMetric(metric);
    },
  });

  fastify.decorate('opcUaCollector', opcUaCollector);
  fastify.decorate('mqttCollector', mqttCollector);

  // ============================================================================
  // ROUTES
  // ============================================================================

  // Health check (public)
  await registerHealthRoutes(fastify);

  // Authentication (public)
  await registerAuthRoutes(fastify);

  // Routes protégées par JWT
  fastify.register(async function protectedRoutes(instance) {
    instance.addHook('preValidation', async (request, reply) => {
      await request.jwtVerify();
    });

    await registerMetricsRoutes(instance);
    await registerEquipmentsRoutes(instance);
    await registerPipelinesRoutes(instance);
  });

  // ============================================================================
  // HOOKS DE FERMETURE
  // ============================================================================

  fastify.addHook('onClose', async (instance) => {
    instance.log.info('Fermeture de l\'application...');
    
    await opcUaCollector.stop();
    await mqttCollector.stop();
    await storage.close();
    
    instance.log.info('Application fermée');
  });

  return fastify;
}

/**
 * Démarre le serveur
 */
export async function start() {
  // Recharger la config pour validation fail-fast
  loadConfig();

  const fastify = await buildApp();

  try {
    await fastify.listen({
      port: config.port,
      host: '0.0.0.0',
    });

    console.log(`
╔═══════════════════════════════════════════════════════════╗
║             🚀 ML_Elec Edge Agent Démarré                ║
╠═══════════════════════════════════════════════════════════╣
║  URL: http://localhost:${config.port}                      ║
║  Health: http://localhost:${config.port}/health            ║
║  Mode: ${config.nodeEnv.padEnd(13)}                          ║
║  Port: ${config.port.toString().padEnd(13)}                         ║
║  OPC-UA: ${config.opcUaEndpoint.padEnd(40)} ║
║  MQTT: ${config.mqttBrokerUrl.padEnd(43)} ║
╚═══════════════════════════════════════════════════════════╝
    `);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
}

// Démarrage si exécuté directement
if (process.argv[1]?.endsWith('index.js') || process.argv[1]?.endsWith('index.ts')) {
  start();
}

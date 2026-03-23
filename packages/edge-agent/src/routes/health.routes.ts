/**
 * Routes de health check (publiques, sans auth)
 */

import { FastifyInstance } from 'fastify';

export async function registerHealthRoutes(fastify: FastifyInstance) {
  // Health check basique
  fastify.get('/health', async (request, reply) => {
    return {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      version: process.env.npm_package_version || 'unknown',
    };
  });

  // Ready check (vérifie les dépendances)
  fastify.get('/ready', async (request, reply) => {
    const checks = {
      storage: (fastify as any).storage !== undefined,
      opcUa: (fastify as any).opcUaCollector !== undefined,
      mqtt: (fastify as any).mqttCollector !== undefined,
    };

    const isReady = Object.values(checks).every(Boolean);

    return {
      status: isReady ? 'ready' : 'not_ready',
      timestamp: new Date().toISOString(),
      checks,
    };
  });

  // Live check (juste pour vérifier que le processus tourne)
  fastify.get('/live', async (request, reply) => {
    return {
      status: 'alive',
      timestamp: Date.now(),
    };
  });
}

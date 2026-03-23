/**
 * Routes de metrics (protégées par JWT)
 */

import { FastifyInstance } from 'fastify';
import { z } from 'zod';

const metricsQuerySchema = z.object({
  equipmentId: z.string().uuid(),
  startTime: z.string().datetime(),
  endTime: z.string().datetime(),
});

export async function registerMetricsRoutes(fastify: FastifyInstance) {
  // Récupérer les metrics d'un équipement
  fastify.get('/metrics', async (request, reply) => {
    try {
      const query = metricsQuerySchema.parse(request.query);
      const storage = (fastify as any).storage;

      const metrics = await storage.getMetrics(
        query.equipmentId,
        new Date(query.startTime),
        new Date(query.endTime)
      );

      return reply.send({
        data: metrics,
        count: metrics.length,
      });
    } catch (error) {
      if (error instanceof z.ZodError) {
        return reply.code(400).send({
          error: 'Validation error',
          details: error.errors,
        });
      }
      throw error;
    }
  });

  // Dernière métrique d'un node
  fastify.get('/metrics/latest/:nodeId', async (request, reply) => {
    const { nodeId } = request.params as { nodeId: string };
    const storage = (fastify as any).storage;

    const metric = await storage.getLatestMetric(nodeId);

    if (!metric) {
      return reply.code(404).send({
        error: 'Not found',
        message: `No metric found for node ${nodeId}`,
      });
    }

    return reply.send(metric);
  });

  // Écrire une métrique (pour tests ou ingestion manuelle)
  fastify.post('/metrics', async (request, reply) => {
    const storage = (fastify as any).storage;
    const body = request.body as any;

    await storage.writeMetric({
      nodeId: body.nodeId,
      value: body.value,
      timestamp: body.timestamp || new Date().toISOString(),
      quality: body.quality || 'good',
      equipmentId: body.equipmentId,
    });

    return reply.code(201).send({
      status: 'created',
    });
  });
}

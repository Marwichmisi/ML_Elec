/**
 * Routes des pipelines (protégées par JWT)
 */

import { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { zPipelineConfig } from '@ml-elec/shared';

const createPipelineSchema = zPipelineConfig.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

const updatePipelineSchema = zPipelineConfig.partial().omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

export async function registerPipelinesRoutes(fastify: FastifyInstance) {
  // Lister tous les pipelines
  fastify.get('/pipelines', async (request, reply) => {
    const storage = (fastify as any).storage;
    const pipelines = await storage.getAllPipelines();

    return reply.send({
      data: pipelines,
      total: pipelines.length,
    });
  });

  // Récupérer un pipeline
  fastify.get('/pipelines/:id', async (request, reply) => {
    const { id } = request.params as { id: string };
    const storage = (fastify as any).storage;
    const pipelines = await storage.getAllPipelines();

    const pipeline = pipelines.find((p: any) => p.id === id);

    if (!pipeline) {
      return reply.code(404).send({
        error: 'Not found',
        message: `Pipeline ${id} not found`,
      });
    }

    return reply.send(pipeline);
  });

  // Créer un pipeline
  fastify.post('/pipelines', async (request, reply) => {
    try {
      const body = createPipelineSchema.parse(request.body);
      const storage = (fastify as any).storage;

      const pipeline = {
        ...body,
        id: crypto.randomUUID(),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      await storage.savePipeline(pipeline);

      return reply.code(201).send(pipeline);
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

  // Mettre à jour un pipeline
  fastify.patch('/pipelines/:id', async (request, reply) => {
    try {
      const { id } = request.params as { id: string };
      const body = updatePipelineSchema.parse(request.body);
      const storage = (fastify as any).storage;

      const pipelines = await storage.getAllPipelines();
      const existing = pipelines.find((p: any) => p.id === id);

      if (!existing) {
        return reply.code(404).send({
          error: 'Not found',
          message: `Pipeline ${id} not found`,
        });
      }

      const updated = {
        ...existing,
        ...body,
        id,
        updatedAt: new Date().toISOString(),
      };

      await storage.savePipeline(updated);

      return reply.send(updated);
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

  // Activer/Désactiver un pipeline
  fastify.post('/pipelines/:id/activate', async (request, reply) => {
    const { id } = request.params as { id: string };
    const { activate } = request.body as { activate: boolean };
    const storage = (fastify as any).storage;
    const pipelines = await storage.getAllPipelines();

    const pipeline = pipelines.find((p: any) => p.id === id);

    if (!pipeline) {
      return reply.code(404).send({
        error: 'Not found',
        message: `Pipeline ${id} not found`,
      });
    }

    pipeline.isActive = activate;
    pipeline.updatedAt = new Date().toISOString();

    await storage.savePipeline(pipeline);

    return reply.send(pipeline);
  });
}

/**
 * Routes des équipements (protégées par JWT)
 */

import { FastifyInstance } from 'fastify';
import { z } from 'zod';
import { zEquipment } from '@ml-elec/shared';

const createEquipmentSchema = zEquipment.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

const updateEquipmentSchema = zEquipment.partial().omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

export async function registerEquipmentsRoutes(fastify: FastifyInstance) {
  // Lister tous les équipements
  fastify.get('/equipments', async (request, reply) => {
    const storage = (fastify as any).storage;
    const equipments = await storage.getAllEquipments();

    return reply.send({
      data: equipments,
      total: equipments.length,
    });
  });

  // Récupérer un équipement
  fastify.get('/equipments/:id', async (request, reply) => {
    const { id } = request.params as { id: string };
    const storage = (fastify as any).storage;
    const equipments = await storage.getAllEquipments();

    const equipment = equipments.find((e: any) => e.id === id);

    if (!equipment) {
      return reply.code(404).send({
        error: 'Not found',
        message: `Equipment ${id} not found`,
      });
    }

    return reply.send(equipment);
  });

  // Créer un équipement
  fastify.post('/equipments', async (request, reply) => {
    try {
      const body = createEquipmentSchema.parse(request.body);
      const storage = (fastify as any).storage;

      const equipment = {
        ...body,
        id: crypto.randomUUID(),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      await storage.saveEquipment(equipment);

      return reply.code(201).send(equipment);
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

  // Mettre à jour un équipement
  fastify.patch('/equipments/:id', async (request, reply) => {
    try {
      const { id } = request.params as { id: string };
      const body = updateEquipmentSchema.parse(request.body);
      const storage = (fastify as any).storage;

      const equipments = await storage.getAllEquipments();
      const existing = equipments.find((e: any) => e.id === id);

      if (!existing) {
        return reply.code(404).send({
          error: 'Not found',
          message: `Equipment ${id} not found`,
        });
      }

      const updated = {
        ...existing,
        ...body,
        id,
        updatedAt: new Date().toISOString(),
      };

      await storage.saveEquipment(updated);

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

  // Supprimer un équipement
  fastify.delete('/equipments/:id', async (request, reply) => {
    const { id } = request.params as { id: string };
    const storage = (fastify as any).storage;
    const equipments = await storage.getAllEquipments();

    const existing = equipments.find((e: any) => e.id === id);

    if (!existing) {
      return reply.code(404).send({
        error: 'Not found',
        message: `Equipment ${id} not found`,
      });
    }

    // TODO: Implémenter delete dans storage
    return reply.code(501).send({
      error: 'Not implemented',
      message: 'Delete equipment not yet implemented',
    });
  });
}

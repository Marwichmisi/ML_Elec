/**
 * Routes d'authentification (publiques)
 * Basé sur @fastify/jwt documentation officielle
 */

import { FastifyInstance } from 'fastify';
import { z } from 'zod';

const loginSchema = z.object({
  username: z.string().min(1),
  password: z.string().min(1),
});

export async function registerAuthRoutes(fastify: FastifyInstance) {
  // Login - génère un token JWT
  fastify.post('/auth/login', async (request, reply) => {
    try {
      const body = loginSchema.parse(request.body);

      // TODO: Implémenter la vérification réelle des credentials
      // Pour l'instant, auth simplifiée pour démo
      if (body.username === 'admin' && body.password === 'admin') {
        const token = fastify.jwt.sign({
          username: body.username,
          role: 'admin',
        });

        return reply.send({
          token,
          expiresIn: '24h',
          user: {
            username: body.username,
            role: 'admin',
          },
        });
      }

      return reply.code(401).send({
        error: 'Invalid credentials',
        message: 'Username or password is incorrect',
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

  // Refresh token
  fastify.post('/auth/refresh', async (request, reply) => {
    try {
      await request.jwtVerify();
      
      const user = request.user as { username: string; role: string };
      const newToken = fastify.jwt.sign(user);

      return reply.send({
        token: newToken,
        expiresIn: '24h',
      });
    } catch (error) {
      return reply.code(401).send({
        error: 'Invalid token',
        message: 'Token is expired or invalid',
      });
    }
  });

  // Get current user
  fastify.get('/auth/me', async (request, reply) => {
    await request.jwtVerify();
    
    return reply.send({
      user: request.user,
    });
  });
}

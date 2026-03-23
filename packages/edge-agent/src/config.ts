/**
 * Configuration de l'application avec validation Zod
 * Fail-fast si une variable requise est manquante
 */

import { z } from 'zod';
import dotenv from 'dotenv';

// Charger les variables d'environnement
dotenv.config();

/**
 * Schéma de validation des variables d'environnement
 */
const envSchema = z.object({
  NODE_ENV: z.enum(['development', 'production', 'test']).default('development'),
  PORT: z.string().transform((val) => parseInt(val, 10)).default('3001'),
  
  // JWT
  JWT_SECRET: z.string().min(32, 'JWT_SECRET doit faire au moins 32 caractères'),
  JWT_EXPIRES_IN: z.string().default('24h'),
  
  // DuckDB
  DUCKDB_DATA_DIR: z.string().min(1).default('./data/duckdb'),
  DUCKDB_MEMORY_LIMIT: z.string().default('1.5GB'),
  DUCKDB_THREADS: z.string().transform((val) => parseInt(val, 10)).default('2'),
  
  // MQTT
  MQTT_BROKER_URL: z.string().url().default('mqtt://localhost:1883'),
  
  // OPC-UA
  OPCUA_ENDPOINT: z.string().default('opc.tcp://localhost:4840'),
  
  // Logging
  LOG_LEVEL: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
});

export type Config = z.infer<typeof envSchema>;

/**
 * Charge et valide la configuration depuis les variables d'environnement
 * Fail-fast si une variable requise est manquante
 */
export function loadConfig(): Config {
  const result = envSchema.safeParse(process.env);

  if (!result.success) {
    console.error('❌ Configuration invalide:');
    result.error.errors.forEach((err) => {
      console.error(`  - ${err.path.join('.')}: ${err.message}`);
    });
    console.error('\nCopiez .env.example vers .env et remplissez les valeurs requises.');
    process.exit(1);
  }

  return result.data;
}

// Export de la configuration chargée
export const config = loadConfig();

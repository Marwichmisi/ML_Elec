/**
 * Configuration Vitest pour edge-agent
 * Basé sur la documentation officielle Vitest
 */

import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    include: ['src/**/*.test.ts'],
    exclude: ['node_modules', 'dist'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      threshold: {
        statements: 85,
        branches: 80,
        functions: 80,
        lines: 85,
      },
      include: ['src/**/*.ts'],
      exclude: ['src/**/*.test.ts', 'node_modules', 'dist'],
    },
    setupFiles: [],
    server: {
      deps: {
        inline: ['@duckdb/node-api'],
      },
    },
  },
});

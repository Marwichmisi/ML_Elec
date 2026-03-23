/**
 * Types et interfaces partagés pour ML_Elec
 * Tous les types doivent être validés par Zod pour la validation runtime
 */

import { z } from 'zod';

// ============================================================================
// ENUMS ET TYPES DE BASE
// ============================================================================

/**
 * Niveaux de sévérité des alertes
 */
export const zAlertSeverity = z.enum(['info', 'warning', 'critical']);
export type AlertSeverity = z.infer<typeof zAlertSeverity>;

/**
 * Statuts d'un équipement
 */
export const zEquipmentStatus = z.enum(['online', 'offline', 'maintenance', 'error']);
export type EquipmentStatus = z.infer<typeof zEquipmentStatus>;

/**
 * Types de sources de données
 */
export const zDataSourceType = z.enum(['mqtt', 'opcua', 'csv']);
export type DataSourceType = z.infer<typeof zDataSourceType>;

/**
 * Qualité d'une donnée OPC-UA/MQTT
 */
export const zDataQuality = z.enum(['good', 'uncertain', 'bad']);
export type DataQuality = z.infer<typeof zDataQuality>;

/**
 * Rôles utilisateurs pour RBAC
 */
export const zUserRole = z.enum(['admin', 'engineer', 'technician', 'viewer']);
export type UserRole = z.infer<typeof zUserRole>;

// ============================================================================
// SCHÉMAS DE DONNÉES MÉTIER
// ============================================================================

/**
 * Métrique brute d'un capteur
 */
export const zSensorMetric = z.object({
  nodeId: z.string().min(1),
  value: z.number(),
  timestamp: z.string().datetime(),
  quality: zDataQuality,
  equipmentId: z.string().uuid().optional(),
});
export type SensorMetric = z.infer<typeof zSensorMetric>;

/**
 * Alerte de détection d'anomalie
 */
export const zAnomalyAlert = z.object({
  id: z.string().uuid(),
  equipmentId: z.string().uuid(),
  severity: zAlertSeverity,
  description: z.string().min(1),
  timestamp: z.string().datetime(),
  confidence: z.number().min(0).max(1),
  nodeId: z.string().optional(),
  value: z.number().optional(),
  threshold: z.number().optional(),
});
export type AnomalyAlert = z.infer<typeof zAnomalyAlert>;

/**
 * Équipement industriel (moteur, compresseur, etc.)
 */
export const zEquipment = z.object({
  id: z.string().uuid(),
  name: z.string().min(1),
  type: z.string(),
  status: zEquipmentStatus,
  location: z.string().optional(),
  metadata: z.record(z.string(), z.unknown()).optional(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});
export type Equipment = z.infer<typeof zEquipment>;

/**
 * Capteur lié à un équipement
 */
export const zSensor = z.object({
  id: z.string().uuid(),
  equipmentId: z.string().uuid(),
  name: z.string().min(1),
  nodeId: z.string().min(1),
  unit: z.string().optional(),
  dataType: z.enum(['temperature', 'vibration', 'current', 'voltage', 'pressure', 'other']),
  createdAt: z.string().datetime(),
});
export type Sensor = z.infer<typeof zSensor>;

// ============================================================================
// PIPELINES DE TRAITEMENT (React Flow)
// ============================================================================

/**
 * Types de nœuds dans un pipeline
 */
export const zPipelineNodeType = z.enum([
  'opcua-source',
  'mqtt-source',
  'csv-source',
  'filter',
  'transform',
  'anomaly-detector',
  'alert',
  'dashboard',
  'export',
]);
export type PipelineNodeType = z.infer<typeof zPipelineNodeType>;

/**
 * Données spécifiques selon le type de nœud
 */
export const zPipelineNodeData = z.object({
  // Configuration commune
  label: z.string().min(1),
  description: z.string().optional(),
  
  // Pour les sources
  endpoint: z.string().optional(),
  nodeId: z.string().optional(),
  pollingInterval: z.number().positive().optional(),
  
  // Pour MQTT
  topic: z.string().optional(),
  qos: z.number().min(0).max(2).optional(),
  
  // Pour les filtres
  filterExpression: z.string().optional(),
  
  // Pour les détecteurs d'anomalies
  threshold: z.number().optional(),
  sensitivity: z.number().min(0).max(1).optional(),
  
  // Pour les alertes
  severity: zAlertSeverity.optional(),
  emailRecipients: z.array(z.string().email()).optional(),
  
  // État runtime
  status: z.enum(['idle', 'running', 'error', 'stopped']).default('idle'),
  lastError: z.string().optional(),
  processedCount: z.number().default(0),
});
export type PipelineNodeData = z.infer<typeof zPipelineNodeData>;

/**
 * Nœud dans un pipeline React Flow
 */
export const zPipelineNode = z.object({
  id: z.string().uuid(),
  type: zPipelineNodeType,
  data: zPipelineNodeData,
  position: z.object({
    x: z.number(),
    y: z.number(),
  }),
});
export type PipelineNode = z.infer<typeof zPipelineNode>;

/**
 * Connexion entre deux nœuds
 */
export const zPipelineEdge = z.object({
  id: z.string().uuid(),
  source: z.string(),
  target: z.string(),
  sourceHandle: z.string().optional(),
  targetHandle: z.string().optional(),
});
export type PipelineEdge = z.infer<typeof zPipelineEdge>;

/**
 * Configuration complète d'un pipeline
 */
export const zPipelineConfig = z.object({
  id: z.string().uuid(),
  name: z.string().min(1),
  description: z.string().optional(),
  nodes: z.array(zPipelineNode),
  edges: z.array(zPipelineEdge),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
  isActive: z.boolean().default(false),
});
export type PipelineConfig = z.infer<typeof zPipelineConfig>;

// ============================================================================
// CONFIGURATION ET UTILITAIRES
// ============================================================================

/**
 * Options de configuration pour l'agent edge
 */
export const zEdgeAgentConfig = z.object({
  port: z.number().positive().default(3001),
  jwtSecret: z.string().min(32),
  jwtExpiresIn: z.string().default('24h'),
  duckDbDataDir: z.string().min(1),
  mqttBrokerUrl: z.string().url(),
  opcUaEndpoint: z.string(),
  logLevel: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
  memoryLimit: z.string().default('1.5GB'),
  threads: z.number().positive().max(4).default(2),
});
export type EdgeAgentConfig = z.infer<typeof zEdgeAgentConfig>;

/**
 * Résultat d'une opération (pattern Result<T, E>)
 */
export type Result<T, E = Error> =
  | { success: true; data: T }
  | { success: false; error: E };

/**
 * Paginated response
 */
export const zPaginatedResponse = <T extends z.ZodType>(schema: T) =>
  z.object({
    data: z.array(schema),
    total: z.number(),
    page: z.number(),
    pageSize: z.number(),
    totalPages: z.number(),
  });

// ============================================================================
// EXPORTS PUBLICS
// ============================================================================

export {
  zAlertSeverity,
  zEquipmentStatus,
  zDataSourceType,
  zDataQuality,
  zUserRole,
  zSensorMetric,
  zAnomalyAlert,
  zEquipment,
  zSensor,
  zPipelineNodeType,
  zPipelineNodeData,
  zPipelineNode,
  zPipelineEdge,
  zPipelineConfig,
  zEdgeAgentConfig,
};

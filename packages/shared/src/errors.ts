/**
 * Erreurs domaine pour ML_Elec
 * Toutes les erreurs métier étendent Error avec des propriétés typées
 */

/**
 * Erreur de connexion OPC-UA
 */
export class OpcUaConnectionError extends Error {
  public readonly endpoint: string;
  public readonly retryCount: number;
  public readonly isRetryable = true;

  constructor(endpoint: string, retryCount: number, cause?: Error) {
    super(`OPC-UA connection failed to ${endpoint} after ${retryCount} retries${cause ? `: ${cause.message}` : ''}`);
    this.name = 'OpcUaConnectionError';
    this.endpoint = endpoint;
    this.retryCount = retryCount;
    if (cause) this.cause = cause;
  }
}

/**
 * Erreur d'écriture dans le stockage
 */
export class StorageWriteError extends Error {
  public readonly table: string;
  public readonly nodeId?: string;

  constructor(table: string, nodeId?: string, cause?: Error) {
    super(`Failed to write to ${table}${nodeId ? ` for node ${nodeId}` : ''}${cause ? `: ${cause.message}` : ''}`);
    this.name = 'StorageWriteError';
    this.table = table;
    this.nodeId = nodeId;
    if (cause) this.cause = cause;
  }
}

/**
 * Erreur de lecture dans le stockage
 */
export class StorageReadError extends Error {
  public readonly table: string;
  public readonly query?: string;

  constructor(table: string, query?: string, cause?: Error) {
    super(`Failed to read from ${table}${query ? ` with query: ${query}` : ''}${cause ? `: ${cause.message}` : ''}`);
    this.name = 'StorageReadError';
    this.table = table;
    this.query = query;
    if (cause) this.cause = cause;
  }
}

/**
 * Erreur de validation de données
 */
export class ValidationError extends Error {
  public readonly field?: string;
  public readonly value?: unknown;

  constructor(message: string, field?: string, value?: unknown) {
    super(message);
    this.name = 'ValidationError';
    this.field = field;
    this.value = value;
  }
}

/**
 * Erreur d'authentification
 */
export class AuthenticationError extends Error {
  public readonly userId?: string;
  public readonly reason: string;

  constructor(reason: string, userId?: string) {
    super(`Authentication failed: ${reason}`);
    this.name = 'AuthenticationError';
    this.userId = userId;
    this.reason = reason;
  }
}

/**
 * Erreur d'autorisation
 */
export class AuthorizationError extends Error {
  public readonly userId: string;
  public readonly requiredPermission: string;
  public readonly resource: string;

  constructor(userId: string, requiredPermission: string, resource: string) {
    super(`User ${userId} is not authorized to access ${resource} (requires ${requiredPermission})`);
    this.name = 'AuthorizationError';
    this.userId = userId;
    this.requiredPermission = requiredPermission;
    this.resource = resource;
  }
}

/**
 * Erreur de configuration
 */
export class ConfigurationError extends Error {
  public readonly key: string;
  public readonly expectedType?: string;

  constructor(key: string, message: string, expectedType?: string) {
    super(`Configuration error for ${key}: ${message}`);
    this.name = 'ConfigurationError';
    this.key = key;
    this.expectedType = expectedType;
  }
}

/**
 * Erreur de pipeline
 */
export class PipelineError extends Error {
  public readonly pipelineId: string;
  public readonly nodeId?: string;

  constructor(pipelineId: string, message: string, nodeId?: string) {
    super(`Pipeline ${pipelineId} error${nodeId ? ` at node ${nodeId}` : ''}: ${message}`);
    this.name = 'PipelineError';
    this.pipelineId = pipelineId;
    this.nodeId = nodeId;
  }
}

/**
 * Erreur de détection d'anomalie
 */
export class AnomalyDetectionError extends Error {
  public readonly equipmentId: string;
  public readonly algorithm?: string;

  constructor(equipmentId: string, message: string, algorithm?: string) {
    super(`Anomaly detection failed for equipment ${equipmentId}${algorithm ? ` using ${algorithm}` : ''}: ${message}`);
    this.name = 'AnomalyDetectionError';
    this.equipmentId = equipmentId;
    this.algorithm = algorithm;
  }
}

/**
 * Erreur MQTT
 */
export class MqttConnectionError extends Error {
  public readonly brokerUrl: string;
  public readonly retryCount: number;
  public readonly isRetryable = true;

  constructor(brokerUrl: string, retryCount: number, cause?: Error) {
    super(`MQTT connection failed to ${brokerUrl} after ${retryCount} retries${cause ? `: ${cause.message}` : ''}`);
    this.name = 'MqttConnectionError';
    this.brokerUrl = brokerUrl;
    this.retryCount = retryCount;
    if (cause) this.cause = cause;
  }
}

/**
 * Erreur de plugin
 */
export class PluginError extends Error {
  public readonly pluginId: string;
  public readonly action: string;

  constructor(pluginId: string, action: string, message: string) {
    super(`Plugin ${pluginId} failed during ${action}: ${message}`);
    this.name = 'PluginError';
    this.pluginId = pluginId;
    this.action = action;
  }
}

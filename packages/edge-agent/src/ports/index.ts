/**
 * Ports (interfaces) pour l'architecture hexagonale
 * Ces interfaces définissent les contrats que les adapters doivent implémenter
 */

import type { SensorMetric, AnomalyAlert, Equipment, PipelineConfig } from '@ml-elec/shared';

/**
 * Port pour le stockage des données
 */
export interface IStorageEngine {
  /**
   * Écrit une métrique dans la base de données
   */
  writeMetric(metric: SensorMetric): Promise<void>;

  /**
   * Écrit un lot de métriques (bulk insert)
   */
  writeMetrics(metrics: SensorMetric[]): Promise<void>;

  /**
   * Lit les métriques pour un équipement donné
   */
  getMetrics(equipmentId: string, startTime: Date, endTime: Date): Promise<SensorMetric[]>;

  /**
   * Lit la dernière métrique pour un nodeId
   */
  getLatestMetric(nodeId: string): Promise<SensorMetric | null>;

  /**
   * Sauvegarde une alerte
   */
  saveAlert(alert: AnomalyAlert): Promise<void>;

  /**
   * Récupère les alertes non acquittées
   */
  getPendingAlerts(): Promise<AnomalyAlert[]>;

  /**
   * Acquitte une alerte
   */
  acknowledgeAlert(alertId: string, userId: string): Promise<void>;

  /**
   * Sauvegarde un équipement
   */
  saveEquipment(equipment: Equipment): Promise<void>;

  /**
   * Récupère tous les équipements
   */
  getAllEquipments(): Promise<Equipment[]>;

  /**
   * Sauvegarde la configuration d'un pipeline
   */
  savePipeline(pipeline: PipelineConfig): Promise<void>;

  /**
   * Récupère tous les pipelines
   */
  getAllPipelines(): Promise<PipelineConfig[]>;

  /**
   * Ferme la connexion à la base de données
   */
  close(): Promise<void>;
}

/**
 * Port pour le collecteur de données OPC-UA
 */
export interface IOpcUaCollector {
  /**
   * Démarre la collecte
   */
  start(): Promise<void>;

  /**
   * Arrête la collecte
   */
  stop(): Promise<void>;

  /**
   * Ajoute un nœud à surveiller
   */
  addNodeToMonitor(nodeId: string, equipmentId?: string): Promise<void>;

  /**
   * Retire un nœud de la surveillance
   */
  removeNodeFromMonitor(nodeId: string): Promise<void>;

  /**
   * État de la connexion
   */
  isConnected(): boolean;

  /**
   * État du collecteur
   */
  isRunning(): boolean;
}

/**
 * Port pour le collecteur de données MQTT
 */
export interface IMqttCollector {
  /**
   * Démarre la collecte
   */
  start(): Promise<void>;

  /**
   * Arrête la collecte
   */
  stop(): Promise<void>;

  /**
   * S'abonne à un topic
   */
  subscribe(topic: string): Promise<void>;

  /**
   * Se désabonne d'un topic
   */
  unsubscribe(topic: string): Promise<void>;

  /**
   * Publie un message
   */
  publish(topic: string, payload: string, qos?: 0 | 1 | 2): Promise<void>;

  /**
   * État de la connexion
   */
  isConnected(): boolean;

  /**
   * État du collecteur
   */
  isRunning(): boolean;
}

/**
 * Port pour le détecteur d'anomalies
 */
export interface IAnomalyDetector {
  /**
   * Analyse une métrique et retourne true si anomalie
   */
  analyze(metric: SensorMetric): Promise<boolean>;

  /**
   * Calcule le score d'anomalie (0-1)
   */
  getAnomalyScore(metric: SensorMetric): Promise<number>;

  /**
   * Entraîne le modèle avec des données historiques
   */
  train(equipmentId: string, metrics: SensorMetric[]): Promise<void>;

  /**
   * Réinitialise le modèle pour un équipement
   */
  reset(equipmentId: string): Promise<void>;
}

/**
 * Port pour le service d'alertes
 */
export interface IAlertService {
  /**
   * Crée une nouvelle alerte
   */
  createAlert(alert: Omit<AnomalyAlert, 'id' | 'timestamp'>): Promise<AnomalyAlert>;

  /**
   * Notifie les recipients d'une alerte
   */
  notify(alert: AnomalyAlert): Promise<void>;

  /**
   * Acquitte une alerte
   */
  acknowledge(alertId: string, userId: string): Promise<void>;
}

/**
 * Callback pour les nouvelles métriques
 */
export type MetricCallback = (metric: SensorMetric) => void | Promise<void>;

/**
 * Callback pour les nouvelles alertes
 */
export type AlertCallback = (alert: AnomalyAlert) => void | Promise<void>;

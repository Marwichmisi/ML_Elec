/**
 * Adapter OPC-UA pour la collecte de données industrielles
 * Implémente IOpcUaCollector avec node-opcua
 * Reconnexion automatique avec backoff exponentiel
 */

import {
  OPCUAClient,
  ClientSession,
  ClientSubscription,
  ClientMonitoredItem,
  AttributeIds,
  TimestampsToReturn,
  ReadValueIdOptions,
  MonitoringParametersOptions,
  OPCUAClientOptions,
  ConnectionStrategyOptions,
  SubscriptionOptions,
} from 'node-opcua-client';
import { DataValue } from 'node-opcua-data-value';
import type { IOpcUaCollector, MetricCallback } from '../ports/index.js';
import type { SensorMetric } from '@ml-elec/shared';
import { OpcUaConnectionError } from '@ml-elec/shared/errors';

export interface OpcUaCollectorOptions {
  endpoint: string;
  pollingInterval?: number;
  samplingInterval?: number;
  queueSize?: number;
  publishingInterval?: number;
  onMetric?: MetricCallback;
}

export class OpcUaCollectorAdapter implements IOpcUaCollector {
  private client: OPCUAClient | null = null;
  private session: ClientSession | null = null;
  private subscription: ClientSubscription | null = null;
  private readonly options: OpcUaCollectorOptions;
  private readonly monitoredNodes = new Map<string, ClientMonitoredItem>();
  private isRunningFlag = false;
  private reconnectAttempts = 0;
  private readonly maxReconnectAttempts = 100;

  constructor(options: OpcUaCollectorOptions) {
    this.options = {
      pollingInterval: options.pollingInterval || 1000,
      samplingInterval: options.samplingInterval || 100,
      queueSize: options.queueSize || 10,
      publishingInterval: options.publishingInterval || 1000,
      ...options,
    };
  }

  /**
   * Démarre la collecte avec reconnexion automatique
   */
  async start(): Promise<void> {
    this.isRunningFlag = true;
    await this.connect();
  }

  /**
   * Se connecte au serveur OPC-UA avec stratégie de reconnexion
   */
  private async connect(): Promise<void> {
    try {
      // Configuration du client avec reconnexion automatique
      const clientOptions: OPCUAClientOptions = {
        connectionStrategy: {
          maxRetry: this.maxReconnectAttempts,
          initialDelay: 1000,
          maxDelay: 60000,
          randomisationFactor: 0.1,
        } as ConnectionStrategyOptions,
        keepSessionAlive: true,
        endpointMustExist: false,
        requestedSessionTimeout: 60000,
      };

      this.client = OPCUAClient.create(clientOptions);

      // Gestion des événements de reconnexion
      this.client.on('connection_reestablished', () => {
        console.log('[OPC-UA] Connexion rétablie');
        this.reconnectAttempts = 0;
        this.onSessionReady();
      });

      this.client.on('connection_failed', (err: Error) => {
        console.error('[OPC-UA] Échec de connexion:', err.message);
        this.reconnectAttempts++;
      });

      this.client.on('close', () => {
        console.log('[OPC-UA] Connexion fermée');
        if (this.isRunningFlag && this.reconnectAttempts < this.maxReconnectAttempts) {
          console.log('[OPC-UA] Tentative de reconnexion...');
        }
      });

      // Connexion au serveur
      console.log(`[OPC-UA] Connexion à ${this.options.endpoint}...`);
      await this.client.connect(this.options.endpoint);
      console.log('[OPC-UA] Connecté');

      // Création de la session
      await this.createSession();
    } catch (error) {
      throw new OpcUaConnectionError(this.options.endpoint, this.reconnectAttempts, error as Error);
    }
  }

  /**
   * Crée une session OPC-UA
   */
  private async createSession(): Promise<void> {
    if (!this.client) {
      throw new OpcUaConnectionError(this.options.endpoint, this.reconnectAttempts);
    }

    try {
      this.session = await this.client.createSession({
        requestedSessionTimeout: 60000,
        actualSessionTimeout: 60000,
      });

      console.log('[OPC-UA] Session créée');
      await this.onSessionReady();
    } catch (error) {
      throw new OpcUaConnectionError(this.options.endpoint, this.reconnectAttempts, error as Error);
    }
  }

  /**
   * Initialisée quand la session est prête
   */
  private async onSessionReady(): Promise<void> {
    if (!this.session) {
      return;
    }

    // Création de la subscription
    const subscriptionOptions: SubscriptionOptions = {
      requestedPublishingInterval: this.options.publishingInterval,
      requestedLifetimeCount: 100,
      requestedMaxKeepAliveCount: 10,
      maxNotificationsPerPublish: 100,
      publishingEnabled: true,
      priority: 10,
    };

    this.subscription = await this.session.createSubscription(subscriptionOptions);
    console.log('[OPC-UA] Subscription créée (ID:', this.subscription.subscriptionId, ')');

    // Gestion des événements de subscription
    this.subscription
      .on('started', () => {
        console.log('[OPC-UA] Subscription démarrée');
      })
      .on('keepalive', () => {
        console.log('[OPC-UA] Keepalive');
      })
      .on('terminated', () => {
        console.log('[OPC-UA] Subscription terminée');
      })
      .on('internal_error', (err: Error) => {
        console.error('[OPC-UA] Erreur interne:', err);
      });

    // Recréer les monitored items après reconnexion
    for (const [nodeId] of this.monitoredNodes) {
      await this.addNodeToMonitorInternal(nodeId);
    }
  }

  /**
   * Ajoute un nœud à surveiller (interne, sans vérification session)
   */
  private async addNodeToMonitorInternal(nodeId: string): Promise<void> {
    if (!this.subscription) {
      return;
    }

    try {
      const itemToMonitor: ReadValueIdOptions = {
        nodeId,
        attributeId: AttributeIds.Value,
      };

      const parameters: MonitoringParametersOptions = {
        samplingInterval: this.options.samplingInterval,
        discardOldest: true,
        queueSize: this.options.queueSize,
      };

      const monitoredItem = ClientMonitoredItem.create(
        this.subscription,
        itemToMonitor,
        parameters,
        TimestampsToReturn.Both
      );

      // Gestion des changements de valeur
      monitoredItem.on('changed', (dataValue: DataValue) => {
        if (dataValue.value && this.options.onMetric) {
          const metric: SensorMetric = {
            nodeId,
            value: dataValue.value.value as number,
            timestamp: dataValue.sourceTimestamp?.toISOString() || new Date().toISOString(),
            quality: this.mapOpcUaQuality(dataValue.statusCode?.value),
          };

          this.options.onMetric(metric);
        }
      });

      monitoredItem.on('err', (err: Error) => {
        console.error(`[OPC-UA] Erreur sur ${nodeId}:`, err);
      });

      this.monitoredNodes.set(nodeId, monitoredItem);
      console.log('[OPC-UA] Nœud surveillé:', nodeId);
    } catch (error) {
      console.error(`[OPC-UA] Échec surveillance ${nodeId}:`, error);
    }
  }

  /**
   * Ajoute un nœud à surveiller
   */
  async addNodeToMonitor(nodeId: string, _equipmentId?: string): Promise<void> {
    if (!this.subscription) {
      console.log('[OPC-UA] Pas de subscription, ajout en attente');
      // La session n'est pas encore prête, le nœud sera ajouté lors de onSessionReady
      this.monitoredNodes.set(nodeId, null as unknown as ClientMonitoredItem);
      return;
    }

    await this.addNodeToMonitorInternal(nodeId);
  }

  /**
   * Retire un nœud de la surveillance
   */
  async removeNodeFromMonitor(nodeId: string): Promise<void> {
    const monitoredItem = this.monitoredNodes.get(nodeId);
    if (monitoredItem) {
      try {
        await monitoredItem.terminate();
        this.monitoredNodes.delete(nodeId);
        console.log('[OPC-UA] Nœud retiré:', nodeId);
      } catch (error) {
        console.error(`[OPC-UA] Échec suppression ${nodeId}:`, error);
      }
    }
  }

  /**
   * Map le statut OPC-UA vers la qualité ML_Elec
   */
  private mapOpcUaQuality(statusCode?: number): 'good' | 'uncertain' | 'bad' {
    if (statusCode === undefined || statusCode === 0) {
      return 'good';
    }
    // Bit 0-1: 00 = Good, 01 = Uncertain, 10 = Bad
    const qualityBits = statusCode & 0x03;
    switch (qualityBits) {
      case 0:
        return 'good';
      case 1:
        return 'uncertain';
      case 2:
      case 3:
        return 'bad';
      default:
        return 'uncertain';
    }
  }

  /**
   * Arrête la collecte
   */
  async stop(): Promise<void> {
    this.isRunningFlag = false;

    // Arrêter tous les monitored items
    for (const [nodeId, item] of this.monitoredNodes.entries()) {
      try {
        await item.terminate();
      } catch (error) {
        console.error(`[OPC-UA] Échec arrêt ${nodeId}:`, error);
      }
    }
    this.monitoredNodes.clear();

    // Terminer la subscription
    if (this.subscription) {
      try {
        await this.subscription.terminate();
      } catch (error) {
        console.error('[OPC-UA] Échec termination subscription:', error);
      }
      this.subscription = null;
    }

    // Fermer la session
    if (this.session) {
      try {
        await this.session.close();
      } catch (error) {
        console.error('[OPC-UA] Échec fermeture session:', error);
      }
      this.session = null;
    }

    // Déconnecter le client
    if (this.client) {
      try {
        await this.client.disconnect();
      } catch (error) {
        console.error('[OPC-UA] Échec déconnexion:', error);
      }
      this.client = null;
    }

    console.log('[OPC-UA] Arrêté');
  }

  /**
   * État de la connexion
   */
  isConnected(): boolean {
    return this.client !== null && this.session !== null;
  }

  /**
   * État du collecteur
   */
  isRunning(): boolean {
    return this.isRunningFlag;
  }

  /**
   * Récupère la session (pour usage avancé)
   */
  getSession(): ClientSession | null {
    return this.session;
  }
}

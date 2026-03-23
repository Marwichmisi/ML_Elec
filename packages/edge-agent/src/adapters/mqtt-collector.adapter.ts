/**
 * Adapter MQTT pour la collecte de données
 * Basé sur mqtt.js (https://github.com/mqttjs/MQTT.js)
 * Reconnexion automatique avec backoff exponentiel
 */

import mqtt, { MqttClient, IClientOptions } from 'mqtt';
import type { IMqttCollector, MetricCallback } from '../ports/index.js';
import type { SensorMetric } from '@ml-elec/shared';
import { MqttConnectionError } from '@ml-elec/shared/errors';

export interface MqttCollectorOptions {
  brokerUrl: string;
  clientId?: string;
  username?: string;
  password?: string;
  qos?: 0 | 1 | 2;
  onMetric?: MetricCallback;
}

export class MqttCollectorAdapter implements IMqttCollector {
  private client: MqttClient | null = null;
  private readonly options: MqttCollectorOptions;
  private readonly subscribedTopics = new Set<string>();
  private isRunningFlag = false;
  private reconnectAttempts = 0;

  constructor(options: MqttCollectorOptions) {
    this.options = {
      qos: 1,
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
   * Se connecte au broker MQTT avec stratégie de reconnexion
   */
  private async connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        const clientId = this.options.clientId || `ml-elec-${Date.now()}`;

        const options: IClientOptions = {
          clientId,
          username: this.options.username,
          password: this.options.password,
          reconnectPeriod: 5000, // 5 secondes
          connectTimeout: 30000, // 30 secondes
          clean: true,
          resubscribe: true, // Resubscribe automatiquement après reconnexion
        };

        console.log(`[MQTT] Connexion à ${this.options.brokerUrl}...`);
        this.client = mqtt.connect(this.options.brokerUrl, options);

        // Gestion de la connexion
        this.client.on('connect', () => {
          console.log('[MQTT] Connecté');
          this.reconnectAttempts = 0;

          // Resubscribe aux topics
          for (const topic of this.subscribedTopics) {
            this.subscribeInternal(topic);
          }

          resolve();
        });

        // Gestion des erreurs
        this.client.on('error', (err: Error) => {
          console.error('[MQTT] Erreur:', err.message);
          if (!this.client?.reconnecting) {
            reject(new MqttConnectionError(this.options.brokerUrl, this.reconnectAttempts, err));
          }
        });

        // Gestion de la reconnexion
        this.client.on('reconnect', () => {
          this.reconnectAttempts++;
          console.log(`[MQTT] Tentative de reconnexion (${this.reconnectAttempts})...`);
        });

        // Gestion des messages
        this.client.on('message', (topic: string, message: Buffer) => {
          this.handleMessage(topic, message);
        });

        // Gestion de la déconnexion
        this.client.on('close', () => {
          console.log('[MQTT] Déconnecté');
          if (this.isRunningFlag && !this.client?.disconnecting) {
            console.log('[MQTT] Reconnexion automatique...');
          }
        });
      } catch (error) {
        reject(new MqttConnectionError(this.options.brokerUrl, this.reconnectAttempts, error as Error));
      }
    });
  }

  /**
   * Gère les messages MQTT reçus
   */
  private handleMessage(topic: string, message: Buffer): void {
    try {
      const payload = JSON.parse(message.toString());

      // Format attendu: { nodeId, value, timestamp?, quality? }
      const metric: SensorMetric = {
        nodeId: payload.nodeId || topic,
        value: payload.value,
        timestamp: payload.timestamp || new Date().toISOString(),
        quality: payload.quality || 'good',
      };

      if (this.options.onMetric) {
        this.options.onMetric(metric);
      }
    } catch (error) {
      console.error(`[MQTT] Erreur parsing message ${topic}:`, error);
    }
  }

  /**
   * S'abonne à un topic (interne)
   */
  private subscribeInternal(topic: string): void {
    if (!this.client) {
      return;
    }

    this.client.subscribe(topic, { qos: this.options.qos }, (err) => {
      if (err) {
        console.error(`[MQTT] Échec subscription ${topic}:`, err);
      } else {
        console.log(`[MQTT] Abonné à ${topic}`);
      }
    });
  }

  /**
   * S'abonne à un topic
   */
  async subscribe(topic: string): Promise<void> {
    this.subscribedTopics.add(topic);
    
    if (this.client && this.client.connected) {
      this.subscribeInternal(topic);
    }
  }

  /**
   * Se désabonne d'un topic
   */
  async unsubscribe(topic: string): Promise<void> {
    this.subscribedTopics.delete(topic);

    if (this.client) {
      return new Promise((resolve, reject) => {
        this.client.unsubscribe(topic, (err) => {
          if (err) {
            reject(err);
          } else {
            console.log(`[MQTT] Désabonné de ${topic}`);
            resolve();
          }
        });
      });
    }
  }

  /**
   * Publie un message
   */
  async publish(topic: string, payload: string, qos?: 0 | 1 | 2): Promise<void> {
    if (!this.client) {
      throw new MqttConnectionError(this.options.brokerUrl, 0, new Error('Not connected'));
    }

    return new Promise((resolve, reject) => {
      this.client!.publish(topic, payload, { qos: qos || this.options.qos }, (err) => {
        if (err) {
          reject(err);
        } else {
          resolve();
        }
      });
    });
  }

  /**
   * Arrête la collecte
   */
  async stop(): Promise<void> {
    this.isRunningFlag = false;

    if (this.client) {
      return new Promise((resolve) => {
        this.client!.end(true, () => {
          this.client = null;
          console.log('[MQTT] Arrêté');
          resolve();
        });
      });
    }
  }

  /**
   * État de la connexion
   */
  isConnected(): boolean {
    return this.client?.connected || false;
  }

  /**
   * État du collecteur
   */
  isRunning(): boolean {
    return this.isRunningFlag;
  }

  /**
   * Récupère le client (pour usage avancé)
   */
  getClient(): MqttClient | null {
    return this.client;
  }
}

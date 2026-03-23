/**
 * Registry des types de nœuds React Flow
 * Basé sur la documentation @xyflow/react
 */

import type { NodeTypes } from '@xyflow/react';
import { OpcUaSourceNode } from './nodes/OpcUaSourceNode.js';
import { MqttSourceNode } from './nodes/MqttSourceNode.js';
import { AnomalyDetectorNode } from './nodes/AnomalyDetectorNode.js';
import { AlertNode } from './nodes/AlertNode.js';

/**
 * Map des types de nœuds custom pour React Flow
 * À passer au prop `nodeTypes` du composant ReactFlow
 */
export const nodeTypes: NodeTypes = {
  'opcua-source': OpcUaSourceNode,
  'mqtt-source': MqttSourceNode,
  'anomaly-detector': AnomalyDetectorNode,
  'alert': AlertNode,
};

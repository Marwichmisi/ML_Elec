/**
 * UI Components for ML_Elec
 * React Flow (@xyflow/react) custom nodes and components
 * Basé sur la documentation officielle: /xyflow/xyflow
 */

// Node types
export { OpcUaSourceNode, type OpcUaSourceData } from './nodes/OpcUaSourceNode.js';
export { MqttSourceNode, type MqttSourceData } from './nodes/MqttSourceNode.js';
export { AnomalyDetectorNode, type AnomalyDetectorData } from './nodes/AnomalyDetectorNode.js';
export { AlertNode, type AlertNodeData } from './nodes/AlertNode.js';

// Node type registry
export { nodeTypes } from './node-types.js';

// Components
export { StatusBadge } from './components/StatusBadge.js';
export { NodeHeader } from './components/NodeHeader.js';
export { NodeField } from './components/NodeField.js';

// Hooks
export { useNodeConnections } from './hooks/useNodeConnections.js';
export { useNodesData } from './hooks/useNodesData.js';

// Types
export type { PipelineNodeData } from '@ml-elec/shared';

// Version
export const version = '1.0.0';

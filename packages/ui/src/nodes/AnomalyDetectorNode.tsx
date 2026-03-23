/**
 * Nœud Détecteur d'Anomalies pour React Flow
 */

import { memo } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { NodeHeader } from '../components/NodeHeader.js';
import { NodeField } from '../components/NodeField.js';
import { StatusBadge } from '../components/StatusBadge.js';

export interface AnomalyDetectorData {
  label: string;
  description?: string;
  threshold?: number;
  sensitivity?: number;
  algorithm?: string;
  status: 'idle' | 'running' | 'error' | 'stopped';
  lastError?: string;
  anomaliesDetected?: number;
}

function AnomalyDetectorNodeComponent({ data }: NodeProps<AnomalyDetectorData>) {
  const statusColor = {
    idle: '#9ca3af',
    running: '#22c55e',
    error: '#ef4444',
    stopped: '#6b7280',
  }[data.status];

  return (
    <div
      className="ml-elec-node anomaly-detector"
      style={{
        minWidth: '280px',
        background: '#fff',
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
      }}
    >
      {/* Input handle */}
      <Handle
        type="target"
        position={Position.Left}
        id="input"
        style={{
          background: '#3b82f6',
          border: '2px solid #fff',
        }}
      />

      {/* Output handle */}
      <Handle
        type="source"
        position={Position.Right}
        id="output"
        style={{
          background: statusColor,
          border: '2px solid #fff',
        }}
      />

      <NodeHeader
        icon="🔍"
        title={data.label || 'Détecteur d\'Anomalies'}
        statusColor={statusColor}
      />

      <div
        style={{
          padding: '12px',
          fontSize: '13px',
          color: '#374151',
        }}
      >
        {data.algorithm && (
          <NodeField label="Algorithme" value={data.algorithm} />
        )}
        {data.threshold !== undefined && (
          <NodeField label="Seuil" value={data.threshold.toFixed(2)} />
        )}
        {data.sensitivity !== undefined && (
          <NodeField
            label="Sensibilité"
            value={`${Math.round(data.sensitivity * 100)}%`}
          />
        )}

        <div style={{ marginTop: '8px' }}>
          <StatusBadge status={data.status} />
        </div>

        {data.lastError && (
          <div
            style={{
              marginTop: '8px',
              padding: '6px 8px',
              background: '#fef2f2',
              border: '1px solid #fecaca',
              borderRadius: '4px',
              fontSize: '12px',
              color: '#dc2626',
            }}
          >
            {data.lastError}
          </div>
        )}

        {data.anomaliesDetected !== undefined && (
          <div
            style={{
              marginTop: '8px',
              fontSize: '11px',
              color: '#6b7280',
            }}
          >
            {data.anomaliesDetected.toLocaleString()} anomalies détectées
          </div>
        )}
      </div>
    </div>
  );
}

export const AnomalyDetectorNode = memo(AnomalyDetectorNodeComponent);

/**
 * Nœud MQTT Source pour React Flow
 */

import { memo } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { NodeHeader } from '../components/NodeHeader.js';
import { NodeField } from '../components/NodeField.js';
import { StatusBadge } from '../components/StatusBadge.js';

export interface MqttSourceData {
  label: string;
  description?: string;
  brokerUrl?: string;
  topic: string;
  qos?: 0 | 1 | 2;
  status: 'idle' | 'running' | 'error' | 'stopped';
  lastError?: string;
  processedCount?: number;
}

function MqttSourceNodeComponent({ data }: NodeProps<MqttSourceData>) {
  const statusColor = {
    idle: '#9ca3af',
    running: '#22c55e',
    error: '#ef4444',
    stopped: '#6b7280',
  }[data.status];

  return (
    <div
      className="ml-elec-node mqtt-source"
      style={{
        minWidth: '280px',
        background: '#fff',
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
      }}
    >
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
        icon="📨"
        title={data.label || 'MQTT Source'}
        statusColor={statusColor}
      />

      <div
        style={{
          padding: '12px',
          fontSize: '13px',
          color: '#374151',
        }}
      >
        <NodeField label="Topic" value={data.topic} mono />
        <NodeField
          label="QoS"
          value={data.qos !== undefined ? `QoS ${data.qos}` : 'N/A'}
        />

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

        {data.processedCount !== undefined && (
          <div
            style={{
              marginTop: '8px',
              fontSize: '11px',
              color: '#6b7280',
            }}
          >
            {data.processedCount.toLocaleString()} messages traités
          </div>
        )}
      </div>
    </div>
  );
}

export const MqttSourceNode = memo(MqttSourceNodeComponent);

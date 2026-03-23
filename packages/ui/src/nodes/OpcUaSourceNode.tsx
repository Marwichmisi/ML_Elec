/**
 * Nœud OPC-UA Source pour React Flow
 * Permet de configurer une source de données OPC-UA
 */

import { memo } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { NodeHeader } from '../components/NodeHeader.js';
import { NodeField } from '../components/NodeField.js';
import { StatusBadge } from '../components/StatusBadge.js';

export interface OpcUaSourceData {
  label: string;
  description?: string;
  endpoint: string;
  nodeId: string;
  pollingInterval?: number;
  status: 'idle' | 'running' | 'error' | 'stopped';
  lastError?: string;
  processedCount?: number;
}

function OpcUaSourceNodeComponent({ data }: NodeProps<OpcUaSourceData>) {
  const statusColor = {
    idle: '#9ca3af',
    running: '#22c55e',
    error: '#ef4444',
    stopped: '#6b7280',
  }[data.status];

  return (
    <div
      className="ml-elec-node opcua-source"
      style={{
        minWidth: '280px',
        background: '#fff',
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
      }}
    >
      {/* Source handle */}
      <Handle
        type="source"
        position={Position.Right}
        id="output"
        style={{
          background: statusColor,
          border: '2px solid #fff',
        }}
      />

      {/* Header */}
      <NodeHeader
        icon="📡"
        title={data.label || 'OPC-UA Source'}
        statusColor={statusColor}
      />

      {/* Body */}
      <div
        style={{
          padding: '12px',
          fontSize: '13px',
          color: '#374151',
        }}
      >
        <NodeField label="Endpoint" value={data.endpoint} />
        <NodeField label="Node ID" value={data.nodeId} mono />
        <NodeField
          label="Polling"
          value={data.pollingInterval ? `${data.pollingInterval}ms` : 'N/A'}
        />

        {/* Status */}
        <div style={{ marginTop: '8px' }}>
          <StatusBadge status={data.status} />
        </div>

        {/* Error message */}
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

        {/* Stats */}
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

export const OpcUaSourceNode = memo(OpcUaSourceNodeComponent);

/**
 * Nœud Alert pour React Flow
 */

import { memo } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { NodeHeader } from '../components/NodeHeader.js';
import { NodeField } from '../components/NodeField.js';

export interface AlertNodeData {
  label: string;
  description?: string;
  severity?: 'info' | 'warning' | 'critical';
  emailRecipients?: string[];
  status: 'idle' | 'running' | 'error' | 'stopped';
  alertsSent?: number;
}

const severityConfig = {
  info: { color: '#3b82f6', label: 'Info', bg: '#dbeafe' },
  warning: { color: '#f59e0b', label: 'Warning', bg: '#fef3c7' },
  critical: { color: '#dc2626', label: 'Critical', bg: '#fee2e2' },
};

function AlertNodeComponent({ data }: NodeProps<AlertNodeData>) {
  const severity = data.severity || 'info';
  const config = severityConfig[severity];

  return (
    <div
      className="ml-elec-node alert"
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
          background: config.color,
          border: '2px solid #fff',
        }}
      />

      <NodeHeader
        icon="🚨"
        title={data.label || 'Alerte'}
        statusColor={config.color}
      />

      <div
        style={{
          padding: '12px',
          fontSize: '13px',
          color: '#374151',
        }}
      >
        {/* Severity badge */}
        <div
          style={{
            display: 'inline-block',
            padding: '4px 10px',
            borderRadius: '9999px',
            fontSize: '12px',
            fontWeight: '500',
            backgroundColor: config.bg,
            color: config.color,
            marginBottom: '8px',
          }}
        >
          {config.label}
        </div>

        {data.emailRecipients && data.emailRecipients.length > 0 && (
          <NodeField
            label="Recipients"
            value={`${data.emailRecipients.length} email(s)`}
          />
        )}

        {data.alertsSent !== undefined && (
          <div
            style={{
              marginTop: '8px',
              fontSize: '11px',
              color: '#6b7280',
            }}
          >
            {data.alertsSent.toLocaleString()} alertes envoyées
          </div>
        )}
      </div>
    </div>
  );
}

export const AlertNode = memo(AlertNodeComponent);

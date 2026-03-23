/**
 * Badge de statut pour les nœuds React Flow
 */

interface StatusBadgeProps {
  status: 'idle' | 'running' | 'error' | 'stopped';
  size?: 'sm' | 'md';
}

const statusConfig = {
  idle: { label: 'En attente', color: '#9ca3af', bg: '#f3f4f6' },
  running: { label: 'En cours', color: '#16a34a', bg: '#dcfce7' },
  error: { label: 'Erreur', color: '#dc2626', bg: '#fef2f2' },
  stopped: { label: 'Arrêté', color: '#4b5563', bg: '#f3f4f6' },
};

export function StatusBadge({ status, size = 'sm' }: StatusBadgeProps) {
  const config = statusConfig[status];
  const fontSize = size === 'sm' ? '11px' : '12px';
  const padding = size === 'sm' ? '3px 8px' : '4px 10px';

  return (
    <span
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: '6px',
        padding,
        borderRadius: '9999px',
        fontSize,
        fontWeight: '500',
        backgroundColor: config.bg,
        color: config.color,
      }}
    >
      <span
        style={{
          width: size === 'sm' ? '6px' : '8px',
          height: size === 'sm' ? '6px' : '8px',
          borderRadius: '50%',
          backgroundColor: config.color,
        }}
      />
      {config.label}
    </span>
  );
}

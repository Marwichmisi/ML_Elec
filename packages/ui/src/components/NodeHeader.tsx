/**
 * En-tête de nœud React Flow
 */

interface NodeHeaderProps {
  icon: string;
  title: string;
  statusColor?: string;
}

export function NodeHeader({ icon, title, statusColor }: NodeHeaderProps) {
  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        padding: '10px 12px',
        borderBottom: '1px solid #e5e7eb',
        background: statusColor
          ? `linear-gradient(to right, ${statusColor}15, transparent)`
          : '#f9fafb',
      }}
    >
      <span style={{ fontSize: '18px' }}>{icon}</span>
      <span
        style={{
          fontSize: '14px',
          fontWeight: '600',
          color: '#111827',
        }}
      >
        {title}
      </span>
    </div>
  );
}

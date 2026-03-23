/**
 * Champ de donnée pour les nœuds React Flow
 */

interface NodeFieldProps {
  label: string;
  value?: string | number;
  mono?: boolean;
}

export function NodeField({ label, value, mono }: NodeFieldProps) {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '6px',
      }}
    >
      <span
        style={{
          fontSize: '12px',
          color: '#6b7280',
          fontWeight: '500',
        }}
      >
        {label}
      </span>
      <span
        style={{
          fontSize: '12px',
          color: '#111827',
          fontFamily: mono ? 'ui-monospace, monospace' : undefined,
          backgroundColor: mono ? '#f3f4f6' : 'transparent',
          padding: mono ? '2px 6px' : '0',
          borderRadius: mono ? '4px' : '0',
        }}
      >
        {value || '—'}
      </span>
    </div>
  );
}

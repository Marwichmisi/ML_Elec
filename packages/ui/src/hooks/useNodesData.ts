/**
 * Hook pour récupérer les données de plusieurs nœuds
 * Basé sur la documentation @xyflow/react
 */

import { useMemo } from 'react';
import { useNodes } from '@xyflow/react';

/**
 * Récupère les données de plusieurs nœuds par leurs IDs
 */
export function useNodesData<T = any>(nodeIds: string[]): Map<string, T> {
  const nodes = useNodes();

  return useMemo(() => {
    const dataMap = new Map<string, T>();

    for (const nodeId of nodeIds) {
      const node = nodes.find((n) => n.id === nodeId);
      if (node) {
        dataMap.set(nodeId, node.data as T);
      }
    }

    return dataMap;
  }, [nodeIds, nodes]);
}

/**
 * Récupère les données d'un seul nœud par son ID
 */
export function useNodeData<T = any>(nodeId: string): T | null {
  const nodes = useNodes();

  return useMemo(() => {
    const node = nodes.find((n) => n.id === nodeId);
    return node ? (node.data as T) : null;
  }, [nodeId, nodes]);
}

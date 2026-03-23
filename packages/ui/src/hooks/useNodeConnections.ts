/**
 * Hook pour récupérer les données de connexion d'un nœud
 * Basé sur la documentation @xyflow/react
 */

import { useCallback, useMemo } from 'react';
import { useNodeId, useReactFlow } from '@xyflow/react';

export interface UseNodeConnectionsResult {
  /** IDs des nœuds sources connectés */
  sourceNodeIds: string[];
  /** IDs des nœuds cibles connectés */
  targetNodeIds: string[];
  /** Données des nœuds sources */
  sourceData: any[];
  /** Données des nœuds cibles */
  targetData: any[];
}

/**
 * Récupère les connexions d'un nœud
 * Doit être utilisé dans un composant enfant d'un nœud React Flow
 */
export function useNodeConnections(): UseNodeConnectionsResult {
  const nodeId = useNodeId();
  const { getEdges, getNodes } = useReactFlow();

  return useMemo(() => {
    if (!nodeId) {
      return {
        sourceNodeIds: [],
        targetNodeIds: [],
        sourceData: [],
        targetData: [],
      };
    }

    const edges = getEdges();
    const nodes = getNodes();

    // Trouver les connexions sortantes (ce nœud est source)
    const outgoingEdges = edges.filter((edge) => edge.source === nodeId);
    const targetNodeIds = outgoingEdges.map((edge) => edge.target);

    // Trouver les connexions entrantes (ce nœud est cible)
    const incomingEdges = edges.filter((edge) => edge.target === nodeId);
    const sourceNodeIds = incomingEdges.map((edge) => edge.source);

    // Récupérer les données des nœuds connectés
    const sourceData = nodes
      .filter((node) => sourceNodeIds.includes(node.id))
      .map((node) => node.data);

    const targetData = nodes
      .filter((node) => targetNodeIds.includes(node.id))
      .map((node) => node.data);

    return {
      sourceNodeIds,
      targetNodeIds,
      sourceData,
      targetData,
    };
  }, [nodeId, getEdges, getNodes]);
}

export function detectCycle(edges: Map<string, string[]>): Error | null {
  const WHITE = 0, GRAY = 1, BLACK = 2;
  const color = new Map<string, number>();
  function dfs(node: string): Error | null {
    color.set(node, GRAY);
    for (const neighbor of edges.get(node) ?? []) {
      const c = color.get(neighbor) ?? WHITE;
      if (c === GRAY) return new Error(`cycle detected at ${neighbor}`);
      if (c === WHITE) { const err = dfs(neighbor); if (err) return err; }
    }
    color.set(node, BLACK);
    return null;
  }
  for (const node of edges.keys()) {
    if ((color.get(node) ?? WHITE) === WHITE) {
      const err = dfs(node);
      if (err) return err;
    }
  }
  return null;
}

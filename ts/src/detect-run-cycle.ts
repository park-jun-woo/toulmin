import type { Graph } from "./graph.js";   // 타입 전용(순환 안전)

// detectRunCycle walks the graph-of-graphs reachable from root via each node's
// runGraph edge, keyed by Graph object identity, using a 3-color DFS
// (WHITE=unvisited, GRAY=visiting, BLACK=done). Re-entering a GRAY graph is an
// execution cycle; a shared sub-graph reached by two paths (diamond DAG) is legal
// thanks to BLACK. Execution composition must be a DAG.
export function detectRunCycle(root: Graph): Error | null {
  const WHITE = 0, GRAY = 1, BLACK = 2;
  const color = new Map<Graph, number>();   // ★ Graph 객체 정체성으로 키잉
  function dfs(g: Graph): Error | null {
    color.set(g, GRAY);
    for (const r of g.rules) {
      const sub = r.runGraph;
      if (!sub) continue;
      const c = color.get(sub) ?? WHITE;
      if (c === GRAY) return new Error(`toulmin: run cycle detected at graph "${sub.name}"`);
      if (c === WHITE) { const err = dfs(sub); if (err) return err; }
    }
    color.set(g, BLACK);
    return null;
  }
  return dfs(root);
}

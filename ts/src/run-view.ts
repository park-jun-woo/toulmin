import type { NodeEvent, RunView } from "./types.js";
import type { DefeatEdge } from "./defeat-edge.js";
import { shortName } from "./short-name.js";

// createRunView builds an immutable, read-only RunView backed by the snapshot
// maps. all() and attackers() return copies so callers cannot mutate the snapshot.
export function createRunView(
  order: NodeEvent[],
  byName: Map<string, NodeEvent>,
  attackers: Map<string, NodeEvent[]>,
): RunView {
  return {
    all: () => order.slice(),
    get: (name) => byName.get(name),
    attackers: (name) => (attackers.get(name) ?? []).slice(),
  };
}

// buildAttackerEvents groups defeats by `to` (edges are already `to ← from`, so we
// do NOT invert) and collects each `from`'s NodeEvent from byName, keyed by the
// `to` node's shortName.
export function buildAttackerEvents(
  defeats: DefeatEdge[],
  byName: Map<string, NodeEvent>,
): Map<string, NodeEvent[]> {
  const attackers = new Map<string, NodeEvent[]>();
  for (const d of defeats) {
    const to = shortName(d.to);
    const ev = byName.get(shortName(d.from));
    if (!ev) continue;
    const list = attackers.get(to) ?? [];
    list.push(ev);
    attackers.set(to, list);
  }
  return attackers;
}

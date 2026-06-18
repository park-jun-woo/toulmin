import type { Context, Trace, TraceEntry } from "./types.js";

export function createTrace(nodes: TraceEntry[], ctx: Context): Trace {
  return {
    all: () => nodes,
    get: (name) => nodes.find((e) => e.name === name),
    ctx: () => ctx,
  };
}

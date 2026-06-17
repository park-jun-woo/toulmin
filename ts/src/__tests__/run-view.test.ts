import { describe, it, expect } from "vitest";
import { createRunView, buildAttackerEvents } from "../run-view.js";
import { NodeEventType, type NodeEvent } from "../types.js";
import type { DefeatEdge } from "../defeat-edge.js";

function ev(name: string): NodeEvent {
  return { name, role: "rule", type: NodeEventType.Active, verdict: 1, evidence: null };
}

describe("createRunView", () => {
  it("all() returns a defensive copy in order", () => {
    const a = ev("a");
    const b = ev("b");
    const view = createRunView([a, b], new Map([["a", a], ["b", b]]), new Map());

    const all = view.all();
    expect(all.map(e => e.name)).toEqual(["a", "b"]);
    all.pop();
    expect(view.all()).toHaveLength(2); // copy — mutation does not leak
  });

  it("get() returns the event when present and undefined when missing", () => {
    const a = ev("a");
    const view = createRunView([a], new Map([["a", a]]), new Map());
    expect(view.get("a")).toBe(a);
    expect(view.get("nope")).toBeUndefined();
  });

  it("attackers() returns a copy when present and [] when absent", () => {
    const a = ev("a");
    const b = ev("b");
    const attackers = new Map<string, NodeEvent[]>([["a", [b]]]);
    const view = createRunView([a, b], new Map([["a", a], ["b", b]]), attackers);

    const atk = view.attackers("a"); // present → ?? left branch
    expect(atk.map(e => e.name)).toEqual(["b"]);
    atk.pop();
    expect(view.attackers("a")).toHaveLength(1); // copy

    expect(view.attackers("b")).toEqual([]); // absent → ?? right branch ([])
  });
});

describe("buildAttackerEvents", () => {
  function edge(from: string, to: string): DefeatEdge {
    return { from, to } as DefeatEdge;
  }

  it("groups attackers by `to`, accumulating multiple onto one target", () => {
    const r1 = ev("r1");
    const r2 = ev("r2");
    const byName = new Map([["r1", r1], ["r2", r2]]);
    // two edges to the same target "w" → second hit takes the `?? []` left branch
    const defeats = [edge("r1", "w"), edge("r2", "w")];

    const attackers = buildAttackerEvents(defeats, byName);
    expect(attackers.get("w")?.map(e => e.name)).toEqual(["r1", "r2"]);
  });

  it("skips edges whose `from` has no NodeEvent (continue branch)", () => {
    const r1 = ev("r1");
    const byName = new Map([["r1", r1]]); // "ghost" intentionally absent
    const defeats = [edge("ghost", "w"), edge("r1", "w")];

    const attackers = buildAttackerEvents(defeats, byName);
    // ghost skipped; only r1 recorded
    expect(attackers.get("w")?.map(e => e.name)).toEqual(["r1"]);
  });

  it("returns an empty map when there are no defeats", () => {
    const attackers = buildAttackerEvents([], new Map());
    expect(attackers.size).toBe(0);
  });
});

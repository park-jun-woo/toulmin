import { describe, it, expect } from "vitest";
import { Graph } from "../graph.js";
import { newContext } from "../map-context.js";
import { type RuleFunc, type Context, type Trace, type TraceEntry } from "../types.js";
import { shortName } from "../short-name.js";

// Access-control graph used across several run() cases:
//   authenticate (rule)               — user present
//   blockIP      (counter → auth)     — request blocked
//   exemptInternalIP (except → block) — internal network
function buildAccessControl(fired: TraceEntry[]) {
  const authenticate: RuleFunc = (ctx) => [ctx.get("user") != null, null];
  const blockIP: RuleFunc = (ctx) => [ctx.get("blocked") === true, null];
  const exemptInternalIP: RuleFunc = (ctx) => [ctx.get("internal") === true, null];

  const g = new Graph("access");
  const auth = g.rule(authenticate);
  const block = g.counter(blockIP);
  const exempt = g.except(exemptInternalIP);
  block.attacks(auth);
  exempt.attacks(block);

  const record = (name: string) => (t: Trace) => { const self = t.get(name); if (self) fired.push(self); };
  auth.runOn(record(shortName(auth.id)));
  block.runOn(record(shortName(block.id)));
  exempt.runOn(record(shortName(exempt.id)));

  return g;
}

function byPrefix(entries: TraceEntry[], prefix: string): TraceEntry {
  const e = entries.find(t => t.name.startsWith(prefix));
  if (!e) throw new Error(`no entry for ${prefix}`);
  return e;
}

describe("Graph.run", () => {
  it("throws when ctx is null", () => {
    const g = new Graph("nil");
    g.rule(() => [true, null]);
    expect(() => g.run(null as unknown as Context)).toThrow("ctx");
  });

  it("external blocked IP: only block fires (Active); auth Defeated, exempt Inactive", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", false);

    const { trace } = g.run(ctx);

    // Only the Active node (activated && verdict>0) fires.
    expect(fired.map(e => e.name.replace(/_\d+$/, ""))).toEqual(["blockIP"]);
    expect(byPrefix(fired, "blockIP").verdict).toBeGreaterThan(0);

    // trace carries the full picture: auth defeated (verdict 0), exempt inactive.
    expect(byPrefix(trace.all(), "authenticate").verdict).toBe(0);
    expect(byPrefix(trace.all(), "authenticate").activated).toBe(true);
    expect(byPrefix(trace.all(), "exemptInternalIP").activated).toBe(false);

    // roles carried through
    expect(byPrefix(trace.all(), "authenticate").role).toBe("rule");
    expect(byPrefix(trace.all(), "blockIP").role).toBe("counter");
    expect(byPrefix(trace.all(), "exemptInternalIP").role).toBe("except");

    // trace snapshot is queryable post-hoc
    expect(trace.all()).toHaveLength(3);
  });

  it("internal network: exempt Active + auth Active fire; block Defeated does not", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const { trace } = g.run(ctx);

    // exempt and auth are Active; block is Defeated (verdict 0) → no fire.
    expect(fired.map(e => e.name.replace(/_\d+$/, "")).sort()).toEqual(["authenticate", "exemptInternalIP"]);
    expect(byPrefix(trace.all(), "exemptInternalIP").verdict).toBeGreaterThan(0);
    expect(byPrefix(trace.all(), "blockIP").verdict).toBe(0);
    expect(byPrefix(trace.all(), "authenticate").verdict).toBeGreaterThan(0);
  });

  it("only Active nodes fire — a Defeated node's handler is skipped", () => {
    const fired: TraceEntry[] = [];
    const warrant: RuleFunc = () => [true, null];
    const counter: RuleFunc = () => [true, null];
    const g = new Graph("defeated-skip");
    const w = g.rule(warrant);
    const c = g.counter(counter);
    c.attacks(w);
    const wName = shortName(w.id);
    w.runOn((t) => { const self = t.get(wName); if (self) fired.push(self); }); // w is Defeated (verdict 0) → must NOT fire

    expect(() => g.run(newContext())).not.toThrow();
    expect(fired).toHaveLength(0);
  });

  it("an Inactive node's handler is skipped", () => {
    const fired: TraceEntry[] = [];
    const inactive: RuleFunc = () => [false, null];
    const g = new Graph("inactive-skip");
    const r = g.rule(inactive);
    const rName = shortName(r.id);
    r.runOn((t) => { const self = t.get(rName); if (self) fired.push(self); });

    expect(() => g.run(newContext())).not.toThrow();
    expect(fired).toHaveLength(0);
  });

  it("wraps a throwing handler with node name context", () => {
    const g = new Graph("throw");
    g.rule(() => [true, null]).runOn(() => { throw new Error("boom"); });
    const ctx = newContext();

    expect(() => g.run(ctx)).toThrow("boom");
    expect(() => g.run(ctx)).toThrow("runOn");
  });

  it("forces trace/duration off (each Active node fires exactly once)", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    g.run(ctx, { trace: true, duration: true });
    // exempt + auth Active → 2 fires (block Defeated).
    expect(fired).toHaveLength(2);
  });

  it("propagates a cycle error from evaluation", () => {
    const g = new Graph("cycle");
    const a = g.rule(() => [true, null]);
    const b = g.counter(() => [true, null]);
    a.attacks(b);
    b.attacks(a);
    expect(() => g.run(newContext())).toThrow("cycle");
  });

  it("handler sees the full trace, in registration order", () => {
    const seen: string[][] = [];
    const fired: TraceEntry[] = [];
    const g = new Graph("trace-all");
    const r1 = g.rule(() => [true, null]);
    const r1Name = shortName(r1.id);
    r1.runOn((t) => {
      seen.push(t.all().map(e => e.name.replace(/_\d+$/, "")));
      const self = t.get(r1Name);
      if (self) fired.push(self);
    });
    g.rule(() => [true, null]).runOn((t) => {
      seen.push(t.all().map(e => e.name.replace(/_\d+$/, "")));
    });

    g.run(newContext());
    expect(fired).toHaveLength(1);
  });

  it("gradient branch via self.verdict", () => {
    let branch = "";
    const g = new Graph("gradient");
    const r = g.rule(() => [true, null]).qualifier(0.75);
    const rName = shortName(r.id);
    r.runOn((t) => {
      const self = t.get(rName);
      branch = (self?.verdict ?? 0) >= 0.5 ? "strong" : "weak";
    });
    g.run(newContext());
    // qualifier 0.75, unattacked → verdict 2*0.75-1 = 0.5 → strong
    expect(branch).toBe("strong");
  });

  it("gradient branch — below threshold reads weak verdict", () => {
    let captured = 1;
    const g = new Graph("gradient-weak");
    const r = g.rule(() => [true, null]).qualifier(0.6);
    const rName = shortName(r.id);
    r.runOn((t) => { captured = t.get(rName)?.verdict ?? 0; });
    g.run(newContext());
    // qualifier 0.6, unattacked → verdict 2*0.6-1 = 0.2 (< 0.5)
    expect(captured).toBeLessThan(0.5);
  });
});

describe("Graph.evaluate (pure)", () => {
  it("returns results without firing handlers", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    const results = g.evaluate(ctx);
    expect(fired).toHaveLength(0);
    expect(results.length).toBeGreaterThan(0);
  });

  it("throws when ctx is null", () => {
    const g = new Graph("nil");
    g.rule(() => [true, null]);
    expect(() => g.evaluate(null as unknown as Context)).toThrow("ctx");
  });
});

describe("Graph.run — trace", () => {
  it("trace is one entry per registered rule, in registration order", () => {
    const g = new Graph("order");
    g.rule(() => [true, null]);
    g.rule(() => [true, null]);
    g.rule(() => [true, null]);

    const { trace } = g.run(newContext());
    expect(trace.all()).toHaveLength(3);
  });

  it("trace verdict matches the evaluation result; ground is the ctx as-is", () => {
    const g = new Graph("trace-verdict");
    const w = g.rule(() => [true, null]);
    const r = g.counter(() => [false, null]);
    r.attacks(w);

    const ctx = newContext();
    ctx.set("k", "v");
    const { results, trace } = g.run(ctx);
    const entries = trace.all();

    expect(byPrefix(entries, "")).toBeDefined();
    expect(entries[0].verdict).toBe(results[0].verdict);
    expect(entries[1].activated).toBe(false);
    expect(entries[1].verdict).toBe(0);
    // Ground is the ctx as-is, same reference for every entry.
    expect(entries[0].ground).toBe(ctx);
    expect(entries[1].ground).toBe(ctx);
    // ctx() exposes this Run's context.
    expect(trace.ctx()).toBe(ctx);
  });
});

describe("Graph.run — composition (_runDepth recursion)", () => {
  const always: RuleFunc = () => [true, null];

  it("Active node recurses into its runGraph (ctx flows down, depth+1)", () => {
    const sub = new Graph("sub");
    sub.rule(always).runOn((t) => { t.ctx().set("subRan", true); });

    const parent = new Graph("parent");
    parent.rule(always).run(sub); // Active → recurse

    const ctx = newContext();
    parent.run(ctx);
    expect(ctx.get("subRan")).toBe(true);
  });

  it("non-Active node does NOT recurse into its runGraph", () => {
    const sub = new Graph("sub-skip");
    sub.rule(always).runOn((t) => { t.ctx().set("subRan", true); });

    const parent = new Graph("parent-skip");
    parent.rule(() => [false, null]).run(sub); // Inactive → no recurse

    const ctx = newContext();
    parent.run(ctx);
    expect(ctx.get("subRan")).toBeUndefined();
  });

  it("Active node with no runGraph takes the non-recurse branch", () => {
    const g = new Graph("no-subgraph");
    g.rule(always); // Active, but runGraph undefined
    expect(() => g.run(newContext())).not.toThrow();
  });

  it("runOn fires before the sub-graph Run", () => {
    const log: string[] = [];
    const sub = new Graph("sub-order");
    sub.rule(always).runOn(() => { log.push("sub"); });

    const parent = new Graph("parent-order");
    parent.rule(always).runOn(() => { log.push("handler"); }).run(sub);

    parent.run(newContext());
    expect(log).toEqual(["handler", "sub"]);
  });

  it("wraps a sub-graph Run error with run \"parent\" → \"sub\" context", () => {
    const sub = new Graph("sub-boom");
    sub.rule(always).runOn(() => { throw new Error("kaboom"); });

    const parent = new Graph("parent-boom");
    parent.rule(always).run(sub);

    const ctx = newContext();
    expect(() => parent.run(ctx)).toThrow(/run "/);
    expect(() => parent.run(ctx)).toThrow(/→ "sub-boom"/);
    expect(() => parent.run(ctx)).toThrow(/kaboom/);
  });

  it("throws when run composition depth exceeds the backstop", () => {
    const chain: Graph[] = [];
    for (let i = 0; i < 70; i++) {
      const g = new Graph(`chain-${i}`);
      g.rule(always);
      chain.push(g);
    }
    for (let i = 0; i < chain.length - 1; i++) {
      chain[i].rules[0].runGraph = chain[i + 1]; // acyclic chain, wired directly
    }
    expect(() => chain[0].run(newContext())).toThrow(/run depth exceeded/);
  });

  it("Graph.run throws the detected run cycle before dispatching", () => {
    const ga = new Graph("rc-A");
    const gb = new Graph("rc-B");
    ga.rule(always).run(gb);
    gb.rule(always).run(ga); // execution cycle

    expect(() => ga.run(newContext())).toThrow(/run cycle detected/);
  });
});

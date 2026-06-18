import { describe, it, expect } from "vitest";
import { Graph } from "../graph.js";
import { newContext } from "../map-context.js";
import { NodeEventType, type RuleFunc, type Context, type NodeEvent, type RunView } from "../types.js";

// Access-control graph used across several run() cases:
//   authenticate (rule)               — user present
//   blockIP      (counter → auth)     — request blocked
//   exemptInternalIP (except → block) — internal network
function buildAccessControl(fired: NodeEvent[]) {
  const authenticate: RuleFunc = (ctx) => [ctx.get("user") != null, null];
  const blockIP: RuleFunc = (ctx) => [ctx.get("blocked") === true, null];
  const exemptInternalIP: RuleFunc = (ctx) => [ctx.get("internal") === true, null];

  const g = new Graph("access");
  const auth = g.rule(authenticate);
  const block = g.counter(blockIP);
  const exempt = g.except(exemptInternalIP);
  block.attacks(auth);
  exempt.attacks(block);

  const record = (_ctx: Context, ev: NodeEvent, _view: RunView) => { fired.push(ev); };
  auth.onActive(record).onDefeated(record).onInactive(record);
  block.onActive(record).onDefeated(record).onInactive(record);
  exempt.onActive(record).onDefeated(record).onInactive(record);

  return g;
}

function byPrefix(fired: NodeEvent[], prefix: string): NodeEvent {
  const ev = fired.find(e => e.name.startsWith(prefix));
  if (!ev) throw new Error(`no event for ${prefix}`);
  return ev;
}

describe("Graph.run", () => {
  it("throws when ctx is null", () => {
    const g = new Graph("nil");
    g.rule(() => [true, null]);
    expect(() => g.run(null as unknown as Context)).toThrow("ctx");
  });

  it("external blocked IP: block Active, auth Defeated (verdict 0), exempt Inactive", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", false);

    const { view } = g.run(ctx);

    // classifyEvent: Active (verdict > 0), Defeated (verdict <= 0), Inactive
    expect(byPrefix(fired, "blockIP").type).toBe(NodeEventType.Active);
    expect(byPrefix(fired, "blockIP").verdict).toBeGreaterThan(0);
    expect(byPrefix(fired, "authenticate").type).toBe(NodeEventType.Defeated);
    expect(byPrefix(fired, "authenticate").verdict).toBe(0);
    expect(byPrefix(fired, "exemptInternalIP").type).toBe(NodeEventType.Inactive);

    // roles carried through
    expect(byPrefix(fired, "authenticate").role).toBe("rule");
    expect(byPrefix(fired, "blockIP").role).toBe("counter");
    expect(byPrefix(fired, "exemptInternalIP").role).toBe("except");

    // view snapshot is queryable post-hoc
    expect(view.all()).toHaveLength(3);
  });

  it("internal network: exempt Active, block Defeated, auth Active", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    g.run(ctx);

    expect(byPrefix(fired, "exemptInternalIP").type).toBe(NodeEventType.Active);
    expect(byPrefix(fired, "blockIP").type).toBe(NodeEventType.Defeated);
    expect(byPrefix(fired, "authenticate").type).toBe(NodeEventType.Active);
  });

  it("skips nodes with no handler registered for the fired event type", () => {
    const fired: NodeEvent[] = [];
    const onlyInactive: RuleFunc = () => [true, null]; // Active, but only onInactive registered
    const g = new Graph("no-handler");
    const r = g.rule(onlyInactive);
    r.onInactive((_c, ev) => fired.push(ev)); // no onActive → selectHandler returns undefined

    expect(() => g.run(newContext())).not.toThrow();
    expect(fired).toHaveLength(0); // handler skipped (!h continue)
  });

  it("wraps a throwing handler with node name and event-type context", () => {
    const g = new Graph("throw");
    g.rule(() => [true, null]).onActive(() => { throw new Error("boom"); });
    const ctx = newContext();

    expect(() => g.run(ctx)).toThrow("boom");
    expect(() => g.run(ctx)).toThrow("Active");
  });

  it("forces trace/duration off (each node fires exactly once)", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    g.run(ctx, { trace: true, duration: true });
    expect(fired).toHaveLength(3);
  });

  it("propagates a cycle error from evaluation", () => {
    const g = new Graph("cycle");
    const a = g.rule(() => [true, null]);
    const b = g.counter(() => [true, null]);
    a.attacks(b);
    b.attacks(a);
    expect(() => g.run(newContext())).toThrow("cycle");
  });

  it("diamond: two counters attack one warrant (accumulating defeat edges per target)", () => {
    const fired: NodeEvent[] = [];
    const w = (() => [true, null]) as RuleFunc;
    const r1 = (() => [true, null]) as RuleFunc;
    const r2 = (() => [true, null]) as RuleFunc;
    const g = new Graph("diamond");
    const ww = g.rule(w);
    const cc1 = g.counter(r1);
    const cc2 = g.counter(r2);
    cc1.attacks(ww);
    cc2.attacks(ww); // second edge to the same target → edges.get(to) ?? [] left branch
    ww.onDefeated((_c, ev, view) => {
      fired.push(ev);
      expect(view.attackers(ev.name)).toHaveLength(2);
    });

    g.run(newContext());
    expect(fired).toHaveLength(1);
  });
});

describe("Graph.evaluate (pure)", () => {
  it("returns results without firing handlers", () => {
    const fired: NodeEvent[] = [];
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

describe("Graph.run — composition (_runDepth recursion)", () => {
  const always: RuleFunc = () => [true, null];

  it("Active node recurses into its runGraph (ctx flows down, depth+1)", () => {
    const sub = new Graph("sub");
    sub.rule(always).onActive((ctx) => { ctx.set("subRan", true); });

    const parent = new Graph("parent");
    parent.rule(always).run(sub); // Active → recurse

    const ctx = newContext();
    parent.run(ctx);
    expect(ctx.get("subRan")).toBe(true);
  });

  it("non-Active node does NOT recurse into its runGraph", () => {
    const sub = new Graph("sub-skip");
    sub.rule(always).onActive((ctx) => { ctx.set("subRan", true); });

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

  it("wraps a sub-graph Run error with run \"parent\" → \"sub\" context", () => {
    const sub = new Graph("sub-boom");
    sub.rule(always).onActive(() => { throw new Error("kaboom"); });

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

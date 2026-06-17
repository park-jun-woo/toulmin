import { describe, it, expect } from "vitest";
import {
  Graph,
  newContext,
  NodeEventType,
  type RuleFunc,
  type Context,
  type NodeEvent,
  type RunView,
} from "../src/index.js";

// Access-control graph:
//   authenticate (rule)  — user present
//   blockIP      (counter, attacks authenticate) — request blocked
//   exemptInternalIP (except, attacks blockIP)    — internal network
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

function viewByPrefix(view: RunView, prefix: string): NodeEvent {
  const ev = view.all().find(e => e.name.startsWith(prefix));
  if (!ev) throw new Error(`no event for ${prefix}`);
  return ev;
}

describe("Run — node event handlers", () => {
  it("3-events: external blocked IP → block Active, auth Defeated", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", false);

    g.run(ctx);

    expect(byPrefix(fired, "blockIP").type).toBe(NodeEventType.Active);
    expect(byPrefix(fired, "authenticate").type).toBe(NodeEventType.Defeated);
    expect(byPrefix(fired, "exemptInternalIP").type).toBe(NodeEventType.Inactive);
  });

  it("3-events: internal network → exempt Active, block Defeated, auth Active", () => {
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
    expect(byPrefix(fired, "authenticate").verdict).toBeCloseTo(1 / 3, 5);
  });

  it("roles are carried on events", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    g.run(ctx);

    expect(byPrefix(fired, "authenticate").role).toBe("rule");
    expect(byPrefix(fired, "blockIP").role).toBe("counter");
    expect(byPrefix(fired, "exemptInternalIP").role).toBe("except");
  });

  it("full pass — unreached nodes still fire Inactive", () => {
    const fired: NodeEvent[] = [];
    const inactiveW: RuleFunc = () => [false, null];
    const unreachedC: RuleFunc = () => [false, null];

    const g = new Graph("full");
    const w = g.rule(inactiveW);
    const c = g.counter(unreachedC);
    c.attacks(w);
    w.onInactive((_ctx, ev) => fired.push(ev));
    c.onInactive((_ctx, ev) => fired.push(ev));

    // lazy evaluate never reaches c (w is inactive → attackers not calc'd)
    g.run(newContext());

    expect(fired.map(e => e.type)).toEqual([NodeEventType.Inactive, NodeEventType.Inactive]);
    expect(byPrefix(fired, "unreachedC").type).toBe(NodeEventType.Inactive);
  });

  it("evaluate stays pure — handlers never fire", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const r1 = g.evaluate(ctx);
    const r2 = g.evaluate(ctx);

    expect(fired).toHaveLength(0);
    expect(r1[0].verdict).toBe(r2[0].verdict);
  });

  it("run results match evaluate results", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const evalResults = g.evaluate(ctx);
    const runResult = g.run(ctx);

    expect(runResult.results.map(r => r.verdict)).toEqual(evalResults.map(r => r.verdict));
  });

  it("handler throw stops & propagates with node + event context", () => {
    const authenticate: RuleFunc = (ctx) => [ctx.get("user") != null, null];
    const g = new Graph("throw");
    g.rule(authenticate).onActive(() => { throw new Error("boom"); });
    const ctx = newContext();
    ctx.set("user", "alice");

    expect(() => g.run(ctx)).toThrow("authenticate");
    expect(() => g.run(ctx)).toThrow("Active");
    expect(() => g.run(ctx)).toThrow("boom");
  });

  it("events order = registration order", () => {
    const fired: NodeEvent[] = [];
    const ruleA: RuleFunc = () => [true, null];
    const ruleB: RuleFunc = () => [true, null];
    const ruleC: RuleFunc = () => [true, null];

    const g = new Graph("order");
    const rec = (_ctx: Context, ev: NodeEvent) => fired.push(ev);
    g.rule(ruleA).onActive(rec);
    g.rule(ruleB).onActive(rec);
    g.rule(ruleC).onActive(rec);

    const result = g.run(newContext());

    const order = result.view.all().map(e => e.name.replace(/_\d+$/, ""));
    expect(order).toEqual(["ruleA", "ruleB", "ruleC"]);
    expect(fired.map(e => e.name.replace(/_\d+$/, ""))).toEqual(["ruleA", "ruleB", "ruleC"]);
  });

  it("verdict === 0 → Defeated", () => {
    const fired: NodeEvent[] = [];
    const warrant: RuleFunc = () => [true, null];
    const counter: RuleFunc = () => [true, null];

    const g = new Graph("zero");
    const w = g.rule(warrant);
    const c = g.counter(counter);
    c.attacks(w);
    w.onActive((_c, ev) => fired.push(ev)).onDefeated((_c, ev) => fired.push(ev));

    g.run(newContext());

    const ev = byPrefix(fired, "warrant");
    expect(ev.verdict).toBe(0);
    expect(ev.type).toBe(NodeEventType.Defeated);
  });

  it("run forces trace/duration off (handlers fire once per node)", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    // even when caller asks for trace, run ignores it; each node fires exactly once
    g.run(ctx, { trace: true, duration: true });

    expect(fired).toHaveLength(3);
    const names = fired.map(e => e.name.replace(/_\d+$/, "")).sort();
    expect(names).toEqual(["authenticate", "blockIP", "exemptInternalIP"]);
  });
});

describe("Run — RunView", () => {
  it("handler sees all 3 nodes via view.all()", () => {
    const seen: string[][] = [];
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    // build a graph whose handlers inspect the full view
    const authenticate: RuleFunc = (c) => [c.get("user") != null, null];
    const blockIP: RuleFunc = (c) => [c.get("blocked") === true, null];
    const exemptInternalIP: RuleFunc = (c) => [c.get("internal") === true, null];
    const g2 = new Graph("view-all");
    const auth = g2.rule(authenticate);
    const block = g2.counter(blockIP);
    const exempt = g2.except(exemptInternalIP);
    block.attacks(auth);
    exempt.attacks(block);
    const inspect = (_c: Context, _ev: NodeEvent, view: RunView) => {
      seen.push(view.all().map(e => e.name.replace(/_\d+$/, "")));
    };
    auth.onActive(inspect);
    block.onDefeated(inspect);
    exempt.onActive(inspect);

    g2.run(ctx);

    // every handler invocation saw all three registered nodes, in registration order
    expect(seen.length).toBe(3);
    for (const names of seen) {
      expect(names).toEqual(["authenticate", "blockIP", "exemptInternalIP"]);
    }
  });

  it("view is immutable — one handler mutating ctx does not change another handler's view", () => {
    const ruleA: RuleFunc = () => [true, null];
    const ruleB: RuleFunc = () => [true, null];

    const g = new Graph("view-immutable");
    const a = g.rule(ruleA);
    const b = g.rule(ruleB);

    const verdictsSeenForA: (number | undefined)[] = [];
    // first handler mutates ctx
    a.onActive((c) => { c.set("mutated", true); });
    // second handler reads the view's record for ruleA — must reflect the full-pass snapshot
    b.onActive((_c, _ev, view) => {
      verdictsSeenForA.push(view.all().find(e => e.name.startsWith("ruleA"))?.verdict);
    });

    const result = g.run(newContext());

    const aVerdict = result.view.all().find(e => e.name.startsWith("ruleA"))?.verdict;
    expect(verdictsSeenForA).toHaveLength(1);
    expect(verdictsSeenForA[0]).toBe(aVerdict);

    // all() returns a copy — mutating it does not affect later reads
    const copy = result.view.all();
    copy.pop();
    expect(result.view.all()).toHaveLength(2);
  });

  it("view.attackers('authenticate') returns the blockIP event", () => {
    let attackers: NodeEvent[] = [];
    const authenticate: RuleFunc = (c) => [c.get("user") != null, null];
    const blockIP: RuleFunc = (c) => [c.get("blocked") === true, null];
    const exemptInternalIP: RuleFunc = (c) => [c.get("internal") === true, null];

    const g = new Graph("view-attackers");
    const auth = g.rule(authenticate);
    const block = g.counter(blockIP);
    const exempt = g.except(exemptInternalIP);
    block.attacks(auth);
    exempt.attacks(block);

    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", false);

    auth.onDefeated((_c, ev, view) => { attackers = view.attackers(ev.name); });

    const result = g.run(ctx);

    expect(attackers).toHaveLength(1);
    expect(attackers[0].name.startsWith("blockIP")).toBe(true);

    // also queryable from the returned view (using the snapshot's actual name)
    const authEv = viewByPrefix(result.view, "authenticate");
    const post = result.view.attackers(authEv.name);
    expect(post.map(e => e.name.startsWith("blockIP"))).toEqual([true]);
  });

  it("view.get('nope') → undefined", () => {
    const ruleA: RuleFunc = () => [true, null];
    const g = new Graph("view-get-miss");
    g.rule(ruleA);
    const result = g.run(newContext());
    expect(result.view.get("nope")).toBeUndefined();
  });

  it("gradient branch via view.get(...)?.verdict", () => {
    const branch: string[] = [];
    const authenticate: RuleFunc = (c) => [c.get("user") != null, null];
    const blockIP: RuleFunc = (c) => [c.get("blocked") === true, null];

    const g = new Graph("view-gradient");
    const auth = g.rule(authenticate);
    const block = g.counter(blockIP);
    block.attacks(auth);

    block.onActive((_c, ev, view) => {
      // consult the attacked node's gradient verdict via the view
      const target = view.all().find(e => e.name.startsWith("authenticate"));
      const v = target?.verdict ?? 0;
      if (ev.verdict >= 0.5 && v <= 0) branch.push("hardBlock");
      else branch.push("softFlag");
    });

    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);

    g.run(ctx);

    expect(branch).toEqual(["hardBlock"]);
  });

  it("run().view is queryable post-hoc", () => {
    const fired: NodeEvent[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const { view } = g.run(ctx);

    expect(view.all()).toHaveLength(3);
    expect(viewByPrefix(view, "authenticate").type).toBe(NodeEventType.Active);
    expect(viewByPrefix(view, "blockIP").type).toBe(NodeEventType.Defeated);
    expect(viewByPrefix(view, "exemptInternalIP").type).toBe(NodeEventType.Active);
  });
});

// touch

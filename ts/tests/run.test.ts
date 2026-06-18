import { describe, it, expect } from "vitest";
import {
  Graph,
  newContext,
  type RuleFunc,
  type Trace,
  type TraceEntry,
} from "../src/index.js";
import { shortName } from "../src/short-name.js";

// Access-control graph:
//   authenticate (rule)  — user present
//   blockIP      (counter, attacks authenticate) — request blocked
//   exemptInternalIP (except, attacks blockIP)    — internal network
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

describe("Run — node handlers (Active only)", () => {
  it("external blocked IP → only block fires (Active); auth Defeated, exempt Inactive in trace", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", false);

    const { trace } = g.run(ctx);

    expect(fired.map(e => e.name.replace(/_\d+$/, ""))).toEqual(["blockIP"]);
    expect(byPrefix(trace.all(), "blockIP").verdict).toBeGreaterThan(0);
    expect(byPrefix(trace.all(), "authenticate").verdict).toBe(0);
    expect(byPrefix(trace.all(), "exemptInternalIP").activated).toBe(false);
  });

  it("internal network → exempt + auth fire (Active); block Defeated does not", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const { trace } = g.run(ctx);

    expect(fired.map(e => e.name.replace(/_\d+$/, "")).sort()).toEqual(["authenticate", "exemptInternalIP"]);
    expect(byPrefix(trace.all(), "exemptInternalIP").verdict).toBeGreaterThan(0);
    expect(byPrefix(trace.all(), "blockIP").verdict).toBe(0);
    expect(byPrefix(trace.all(), "authenticate").verdict).toBeCloseTo(1 / 3, 5);
  });

  it("roles are carried on trace entries", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const { trace } = g.run(ctx);

    expect(byPrefix(trace.all(), "authenticate").role).toBe("rule");
    expect(byPrefix(trace.all(), "blockIP").role).toBe("counter");
    expect(byPrefix(trace.all(), "exemptInternalIP").role).toBe("except");
  });

  it("full pass — unreached nodes still appear in the trace as inactive", () => {
    const inactiveW: RuleFunc = () => [false, null];
    const unreachedC: RuleFunc = () => [false, null];

    const g = new Graph("full");
    const w = g.rule(inactiveW);
    const c = g.counter(unreachedC);
    c.attacks(w);

    const { trace } = g.run(newContext());

    expect(trace.all().map(e => e.activated)).toEqual([false, false]);
    expect(byPrefix(trace.all(), "unreachedC").activated).toBe(false);
  });

  it("evaluate stays pure — handlers never fire", () => {
    const fired: TraceEntry[] = [];
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
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const evalResults = g.evaluate(ctx);
    const runResult = g.run(ctx);

    expect(runResult.results.map(r => r.verdict)).toEqual(evalResults.map(r => r.verdict));
  });

  it("handler throw stops & propagates with node + runOn context", () => {
    const authenticate: RuleFunc = (ctx) => [ctx.get("user") != null, null];
    const g = new Graph("throw");
    g.rule(authenticate).runOn(() => { throw new Error("boom"); });
    const ctx = newContext();
    ctx.set("user", "alice");

    expect(() => g.run(ctx)).toThrow("authenticate");
    expect(() => g.run(ctx)).toThrow("runOn");
    expect(() => g.run(ctx)).toThrow("boom");
  });

  it("trace order = registration order", () => {
    const ruleA: RuleFunc = () => [true, null];
    const ruleB: RuleFunc = () => [true, null];
    const ruleC: RuleFunc = () => [true, null];

    const g = new Graph("order");
    g.rule(ruleA);
    g.rule(ruleB);
    g.rule(ruleC);

    const { trace } = g.run(newContext());
    const order = trace.all().map(e => e.name.replace(/_\d+$/, ""));
    expect(order).toEqual(["ruleA", "ruleB", "ruleC"]);
  });

  it("verdict === 0 → Defeated → handler skipped", () => {
    const fired: TraceEntry[] = [];
    const warrant: RuleFunc = () => [true, null];
    const counter: RuleFunc = () => [true, null];

    const g = new Graph("zero");
    const w = g.rule(warrant);
    const c = g.counter(counter);
    c.attacks(w);
    const wName = shortName(w.id);
    w.runOn((t) => { const self = t.get(wName); if (self) fired.push(self); });

    const { trace } = g.run(newContext());

    expect(fired).toHaveLength(0);
    expect(byPrefix(trace.all(), "warrant").verdict).toBe(0);
  });

  it("run forces trace/duration off (Active nodes fire once)", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    g.run(ctx, { trace: true, duration: true });

    const names = fired.map(e => e.name.replace(/_\d+$/, "")).sort();
    expect(names).toEqual(["authenticate", "exemptInternalIP"]);
  });
});

describe("Run — trace argument (node inspection)", () => {
  it("handler sees all nodes via the trace argument, in registration order", () => {
    const seen: string[][] = [];
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const authenticate: RuleFunc = (c) => [c.get("user") != null, null];
    const blockIP: RuleFunc = (c) => [c.get("blocked") === true, null];
    const exemptInternalIP: RuleFunc = (c) => [c.get("internal") === true, null];
    const g2 = new Graph("view-all");
    const auth = g2.rule(authenticate);
    const block = g2.counter(blockIP);
    const exempt = g2.except(exemptInternalIP);
    block.attacks(auth);
    exempt.attacks(block);
    const inspect = (t: Trace) => {
      seen.push(t.all().map(e => e.name.replace(/_\d+$/, "")));
    };
    auth.runOn(inspect);
    block.runOn(inspect);
    exempt.runOn(inspect);

    g2.run(ctx);

    // exempt + auth Active fire (block Defeated). Each invocation saw all three nodes.
    expect(seen.length).toBe(2);
    for (const names of seen) {
      expect(names).toEqual(["authenticate", "blockIP", "exemptInternalIP"]);
    }
  });

  it("trace verdict is the full-pass snapshot — visible to every handler", () => {
    const ruleA: RuleFunc = () => [true, null];
    const ruleB: RuleFunc = () => [true, null];

    const g = new Graph("trace-snapshot");
    const a = g.rule(ruleA);
    const b = g.rule(ruleB);

    const verdictsSeenForA: (number | undefined)[] = [];
    a.runOn((t) => { t.ctx().set("mutated", true); });
    b.runOn((t) => {
      verdictsSeenForA.push(t.all().find(e => e.name.startsWith("ruleA"))?.verdict);
    });

    const { trace } = g.run(newContext());

    const aVerdict = trace.all().find(e => e.name.startsWith("ruleA"))?.verdict;
    expect(verdictsSeenForA).toHaveLength(1);
    expect(verdictsSeenForA[0]).toBe(aVerdict);
  });

  it("self carries the firing node's continuous verdict (gradient branch)", () => {
    const branch: string[] = [];
    const authenticate: RuleFunc = (c) => [c.get("user") != null, null];
    const blockIP: RuleFunc = (c) => [c.get("blocked") === true, null];

    const g = new Graph("view-gradient");
    const auth = g.rule(authenticate);
    const block = g.counter(blockIP);
    block.attacks(auth);

    const blockName = shortName(block.id);
    block.runOn((t) => {
      const target = t.all().find(e => e.name.startsWith("authenticate"));
      const v = target?.verdict ?? 0;
      const self = t.get(blockName);
      if ((self?.verdict ?? 0) >= 0.5 && v <= 0) branch.push("hardBlock");
      else branch.push("softFlag");
    });

    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);

    g.run(ctx);

    expect(branch).toEqual(["hardBlock"]);
  });

  it("run().trace is queryable post-hoc", () => {
    const fired: TraceEntry[] = [];
    const g = buildAccessControl(fired);
    const ctx = newContext();
    ctx.set("user", "alice");
    ctx.set("blocked", true);
    ctx.set("internal", true);

    const { trace } = g.run(ctx);

    expect(trace.all()).toHaveLength(3);
    expect(byPrefix(trace.all(), "authenticate").verdict).toBeGreaterThan(0);
    expect(byPrefix(trace.all(), "blockIP").verdict).toBe(0);
    expect(byPrefix(trace.all(), "exemptInternalIP").verdict).toBeGreaterThan(0);
  });
});

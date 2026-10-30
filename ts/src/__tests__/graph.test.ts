import { describe, it, expect } from "vitest";
import { Graph } from "../graph.js";
import { newContext } from "../map-context.js";
import { Strength, type RuleFunc, type Context, type TraceEntry, type Spec } from "../types.js";

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

  const record = (self: TraceEntry) => { fired.push(self); };
  auth.runOn(record);
  block.runOn(record);
  exempt.runOn(record);

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

    // exempt and auth are Active; block is gated out by except (verdict -1, was 0) → no fire.
    expect(fired.map(e => e.name.replace(/_\d+$/, "")).sort()).toEqual(["authenticate", "exemptInternalIP"]);
    expect(byPrefix(trace.all(), "exemptInternalIP").verdict).toBeGreaterThan(0);
    expect(byPrefix(trace.all(), "blockIP").verdict).toBe(-1);
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
    w.runOn((self) => { fired.push(self); }); // w is Defeated (verdict 0) → must NOT fire

    expect(() => g.run(newContext())).not.toThrow();
    expect(fired).toHaveLength(0);
  });

  it("an Inactive node's handler is skipped", () => {
    const fired: TraceEntry[] = [];
    const inactive: RuleFunc = () => [false, null];
    const g = new Graph("inactive-skip");
    const r = g.rule(inactive);
    r.runOn((self) => { fired.push(self); });

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
    g.rule(() => [true, null]).runOn((self, t) => {
      seen.push(t.all().map(e => e.name.replace(/_\d+$/, "")));
      fired.push(self);
    });
    g.rule(() => [true, null]).runOn((_self, t) => {
      seen.push(t.all().map(e => e.name.replace(/_\d+$/, "")));
    });

    g.run(newContext());
    expect(fired).toHaveLength(1);
  });

  it("gradient branch via self.verdict", () => {
    let branch = "";
    const g = new Graph("gradient");
    const r = g.rule(() => [true, null]).qualifier(0.75);
    r.runOn((self) => {
      branch = self.verdict >= 0.5 ? "strong" : "weak";
    });
    g.run(newContext());
    // qualifier 0.75, unattacked → verdict 2*0.75-1 = 0.5 → strong
    expect(branch).toBe("strong");
  });

  it("gradient branch — below threshold reads weak verdict", () => {
    let captured = 1;
    const g = new Graph("gradient-weak");
    const r = g.rule(() => [true, null]).qualifier(0.6);
    r.runOn((self) => { captured = self.verdict; });
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
    sub.rule(always).runOn((_self, t) => { t.ctx().set("subRan", true); });

    const parent = new Graph("parent");
    parent.rule(always).run(sub); // Active → recurse

    const ctx = newContext();
    parent.run(ctx);
    expect(ctx.get("subRan")).toBe(true);
  });

  it("non-Active node does NOT recurse into its runGraph", () => {
    const sub = new Graph("sub-skip");
    sub.rule(always).runOn((_self, t) => { t.ctx().set("subRan", true); });

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

// ── verdict semantics (verdict rule functions, except gate, attack force) ──

class TagSpec implements Spec {
  constructor(public tag: string) {}
  specName() { return "TagSpec"; }
  validate() {}
}

const warrant: RuleFunc = () => [true, null];
const rebuttal: RuleFunc = () => [true, null];
const defeater: RuleFunc = () => [true, null];
const close = (a: number, b: number) => expect(Math.abs(a - b)).toBeLessThan(1e-12);

describe("verdict rule functions — boolean compatibility", () => {
  it("boolean rules give exactly the same verdicts as ±1 rules and the legacy formula", () => {
    const boolOf = (b: boolean): RuleFunc => () => [b, null];
    const numOf = (b: boolean): RuleFunc => () => [b ? 1 : -1, null];
    // w ← c1 (q 0.7), c2 (q 0.4, inactive), c3 (q 0.3); c1 ← c4 (q 0.9)
    const build = (useBool: boolean) => {
      const g = new Graph("compat");
      const mk = useBool ? boolOf : numOf;
      const w = g.rule(mk(true)).with(new TagSpec("w")).qualifier(0.8);
      const c1 = g.counter(mk(true)).with(new TagSpec("c1")).qualifier(0.7);
      const c2 = g.counter(mk(false)).with(new TagSpec("c2")).qualifier(0.4);
      const c3 = g.counter(mk(true)).with(new TagSpec("c3")).qualifier(0.3);
      const c4 = g.counter(mk(true)).with(new TagSpec("c4")).qualifier(0.9);
      c1.attacks(w);
      c2.attacks(w);
      c3.attacks(w);
      c4.attacks(c1);
      return g;
    };
    // legacy formula: raw = q/(1+Σ(v_j+1)/2), v = 2raw−1
    const vc4 = 2 * 0.9 - 1;
    const vc1 = 2 * (0.7 / (1.0 + (vc4 + 1.0) / 2.0)) - 1;
    const vc3 = 2 * 0.3 - 1;
    const legacy = 2 * (0.8 / (1.0 + ((vc1 + 1.0) / 2.0 + 0.0 + (vc3 + 1.0) / 2.0))) - 1;
    for (const trace of [false, true]) {
      const rb = build(true).evaluate(newContext(), { trace });
      const rn = build(false).evaluate(newContext(), { trace });
      expect(rb[0].verdict).toBe(legacy);
      expect(rn[0].verdict).toBe(legacy);
      expect(rb[0].applicability).toBe(1.0);
    }
  });

  it("false maps to -1 (not applicable): no result, ruleVerdict -1", () => {
    const g = new Graph("false");
    g.rule(() => [false, "ev"]);
    expect(g.evaluate(newContext())).toHaveLength(0);
    const { trace } = g.run(newContext());
    expect(trace.all()[0].ruleVerdict).toBe(-1);
    expect(trace.all()[0].activated).toBe(false);
    expect(trace.all()[0].evidence).toBe("ev");
  });
});

describe("verdict rule functions — range and NaN", () => {
  for (const bad of [1.5, -2.0, NaN, Infinity]) {
    it(`verdict ${bad} is an evaluation error naming the rule (evaluate, trace, run)`, () => {
      const g = new Graph("range");
      const w = g.rule(warrant);
      g.counter(() => [bad, null]).attacks(w);
      const calls = [
        () => g.evaluate(newContext()),
        () => g.evaluate(newContext(), { trace: true }),
        () => g.run(newContext()),
      ];
      for (const call of calls) {
        expect(call).toThrow(/^rule ".+": rule verdict .+ out of range \[-1, 1\]$/);
      }
    });
  }

  it("a non-number, non-boolean verdict is an evaluation error", () => {
    const g = new Graph("type");
    g.rule((() => ["yes", null]) as unknown as RuleFunc);
    expect(() => g.evaluate(newContext())).toThrow(/rule verdict must be a number or boolean, got string/);
  });

  it("boundary values -1 and +1 are accepted", () => {
    const g = new Graph("bounds");
    g.rule(() => [1, null]);
    g.rule(() => [-1, null]);
    const res = g.evaluate(newContext());
    expect(res).toHaveLength(1);
    expect(res[0].verdict).toBe(1);
  });
});

describe("verdict rule functions — partial verdicts", () => {
  it("v = 0 means half applies (NOT not-applicable): verdict 0, applicability 0.5", () => {
    const g = new Graph("half");
    g.rule(() => [0, "half"]);
    const res = g.evaluate(newContext());
    expect(res).toHaveLength(1);
    expect(res[0]).toMatchObject({ verdict: 0, applicability: 0.5, evidence: "half" });
  });

  it("qualifier scales applicability: v 0.5, q 0.8 → verdict 0.2", () => {
    const g = new Graph("q");
    g.rule(() => [0.5, null]).qualifier(0.8);
    const res = g.evaluate(newContext());
    close(res[0].verdict, 0.2);
    expect(res[0].applicability).toBe(0.75);
  });

  it("v = -1 is not applicable and attackers are not evaluated (lazy)", () => {
    let attacked = false;
    const g = new Graph("lazy");
    const w = g.rule(() => [-1, null]);
    g.counter(() => { attacked = true; return [1, null]; }).attacks(w);
    expect(g.evaluate(newContext())).toHaveLength(0);
    expect(attacked).toBe(false);
  });

  it("v = -0.5 is active and evaluates its attackers", () => {
    let attacked = false;
    const g = new Graph("barely");
    const w = g.rule(() => [-0.5, null]);
    g.counter(() => { attacked = true; return [-1, null]; }).attacks(w);
    const res = g.evaluate(newContext());
    expect(res).toHaveLength(1);
    expect(res[0].verdict).toBe(-0.5);
    expect(attacked).toBe(true);
  });

  it("trace carries each rule's own verdict; a partial rule is Active and fires runOn", () => {
    const fired: number[] = [];
    const g = new Graph("rv");
    const w = g.rule(() => [0.8, null]).runOn((self: TraceEntry) => { fired.push(self.ruleVerdict); });
    g.counter(() => [-1, null]).attacks(w);
    const want: Record<string, number> = { rule: 0.8, counter: -1 };
    const { trace } = g.run(newContext());
    expect(fired).toEqual([0.8]);
    for (const e of trace.all()) {
      expect(e.ruleVerdict).toBe(want[e.role]);
      expect(e.activated).toBe(e.ruleVerdict > -1);
    }
    const res = g.evaluate(newContext(), { trace: true });
    for (const e of res[0].trace!) expect(e.ruleVerdict).toBe(want[e.role]);
  });

  it("run trace reports ruleVerdict -1 for a node that never ran", () => {
    const g = new Graph("never");
    const w = g.rule(warrant);
    g.counter(rebuttal).attacks(w);
    g.rules[1].fn = undefined as unknown as RuleFunc;   // no function → never runs
    const { trace } = g.run(newContext());
    expect(trace.all()[1].ruleVerdict).toBe(-1);
  });
});

describe("except gate and applicability", () => {
  const build = (v: number, asExcept: boolean, n: number) => {
    const g = new Graph("gate");
    const w = g.rule(warrant);
    for (let i = 0; i < n; i++) {
      const fn: RuleFunc = () => [v, null];
      const a = (asExcept ? g.except(fn) : g.counter(fn)).with(new TagSpec(String.fromCharCode(97 + i)));
      a.attacks(w);
    }
    return g;
  };
  const cases = [
    { name: "three counters stay in the denominator", v: 1, except: false, verdict: -0.5, applic: 1 },
    { name: "three full excepts gate to 0", v: 1, except: true, verdict: -1, applic: 0 },
    { name: "three partial excepts", v: 0.5, except: true, verdict: 2 * 0.25 ** 3 - 1, applic: 0.25 ** 3 },
    { name: "inactive excepts leave the gate at 1", v: -1, except: true, verdict: 1, applic: 1 },
  ];
  for (const c of cases) {
    it(c.name, () => {
      for (const trace of [false, true]) {
        const res = build(c.v, c.except, 3).evaluate(newContext(), { trace });
        expect(res).toHaveLength(1);
        close(res[0].verdict, c.verdict);
        close(res[0].applicability, c.applic);
      }
    });
  }

  it("a node registered with role except gates even when its strength is not Defeater", () => {
    const g = new Graph("role-except");
    const w = g.rule(warrant);
    g.except(defeater).attacks(w);
    g.rules[1].strength = Strength.Defeasible;   // role map still says "except"
    expect(g.evaluate(newContext())[0].verdict).toBe(-1);
  });

  it("strict rules ignore both counters and excepts", () => {
    const g = new Graph("strict");
    const w = g.rule(warrant).qualifier(0.9);
    g.counter(rebuttal).attacks(w);
    g.except(defeater).attacks(w);
    g.rules[0].strength = Strength.Strict;
    const res = g.evaluate(newContext());
    close(res[0].verdict, 0.8);
    expect(res[0].applicability).toBe(1);
  });
});

describe("attack force", () => {
  it("defaults to 1 when omitted", () => {
    const g = new Graph("default");
    const w = g.rule(warrant);
    g.counter(rebuttal).attacks(w);
    expect(g.defeats[0].force).toBe(1);
    expect(g.evaluate(newContext())[0].verdict).toBe(0);
  });

  it("an edge without a stored force counts as 1", () => {
    const g = new Graph("raw-edge");
    const w = g.rule(warrant);
    const c = g.counter(rebuttal);
    g.defeats.push({ from: c.id, to: w.id } as unknown as { from: string; to: string; force: number });
    expect(g.evaluate(newContext())[0].verdict).toBe(0);
  });

  const cases = [
    { name: "counter force 0", except: false, force: 0, verdict: 1 },
    { name: "counter force 0.5", except: false, force: 0.5, verdict: 1 / 3 },   // S = 0.5 → raw 2/3
    { name: "counter force 1", except: false, force: 1, verdict: 0 },
    { name: "except force 0", except: true, force: 0, verdict: 1 },
    { name: "except force 0.5", except: true, force: 0.5, verdict: 0 },          // a_eff = 0.5
    { name: "except force 0.25", except: true, force: 0.25, verdict: 0.5 },
  ];
  for (const c of cases) {
    it(`scales the attack: ${c.name}`, () => {
      for (const trace of [false, true]) {
        const g = new Graph("force");
        const w = g.rule(warrant);
        const a = c.except ? g.except(defeater) : g.counter(rebuttal);
        a.attacks(w, c.force);
        close(g.evaluate(newContext(), { trace })[0].verdict, c.verdict);
      }
    });
  }

  it("half-force except leaves a partial defeat (counter verdict 0 → warrant 1/3)", () => {
    const g = new Graph("half-except");
    const w = g.rule(warrant);
    const r = g.counter(rebuttal);
    r.attacks(w);
    g.except(defeater).attacks(r, 0.5);
    close(g.evaluate(newContext())[0].verdict, 1 / 3);
  });
});

describe("evaluation paths — caching, errors, missing nodes", () => {
  it("a shared attacker is computed once and reused (diamond), with and without trace", () => {
    for (const trace of [false, true]) {
      let calls = 0;
      const g = new Graph("diamond");
      const w = g.rule(warrant);
      const r1 = g.counter(rebuttal).with(new TagSpec("r1"));
      const r2 = g.counter(rebuttal).with(new TagSpec("r2"));
      const d = g.counter(() => { calls++; return [0, null]; });   // half-applying shared rebuttal
      r1.attacks(w);
      r2.attacks(w);
      d.attacks(r1);
      d.attacks(r2);
      const res = g.evaluate(newContext(), { trace });
      expect(calls).toBe(1);
      // d: verdict 0 → each rᵢ: S = 0.5 → 1/3 → w: S = 2·(2/3) → raw 3/7
      close(res[0].verdict, 2 * (3 / 7) - 1);
      if (trace) expect(res[0].trace).toHaveLength(4);
    }
  });

  it("after a rule error, remaining attackers are skipped and the error propagates", () => {
    for (const trace of [false, true]) {
      let reached = false;
      const g = new Graph("err-skip");
      const w = g.rule(warrant);
      g.counter(() => { throw new Error("first"); }).attacks(w);
      g.counter(() => { reached = true; return [1, null]; }).attacks(w);
      expect(() => g.evaluate(newContext(), { trace })).toThrow(/rule ".+": Error: first/);
      expect(reached).toBe(false);
    }
  });

  it("an attacker without a function counts as not applicable (trace path)", () => {
    const g = new Graph("no-fn-trace");
    const w = g.rule(warrant);
    g.counter(rebuttal).attacks(w);
    g.rules[1].fn = undefined as unknown as RuleFunc;
    const res = g.evaluate(newContext(), { trace: true });
    expect(res[0].verdict).toBe(1);
    expect(res[0].trace).toHaveLength(1);
  });

  it("duration records a per-rule elapsed time alongside ruleVerdict", () => {
    const g = new Graph("duration");
    g.rule(() => [0.5, null]);
    const res = g.evaluate(newContext(), { duration: true });
    const e = res[0].trace![0];
    expect(e.duration).toBeGreaterThanOrEqual(0);
    expect(e.ruleVerdict).toBe(0.5);
    expect(e.qualifier).toBe(1);
  });

  it("run's full pass surfaces an error from a node unreached during the warrant pass", () => {
    const g = new Graph("full-pass-err");
    const w = g.rule(() => [-1, null]);                 // not applicable → attackers skipped
    g.counter(() => [2, null]).attacks(w);              // evaluated only by the full pass
    expect(g.evaluate(newContext())).toHaveLength(0);  // evaluate never reaches it
    expect(() => g.run(newContext())).toThrow(/out of range \[-1, 1\]/);
  });

  it("run trace falls back to role \"rule\" when a node has no registered role", () => {
    const g = new Graph("no-role");
    g.rule(warrant);
    g.roles.clear();
    const { trace } = g.run(newContext());
    expect(trace.all()[0].role).toBe("rule");
  });
});

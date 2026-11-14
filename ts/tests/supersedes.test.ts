import { describe, it, expect } from "vitest";
import { Graph, Strength, newContext, type RuleFunc, type Rule, type TraceEntry } from "../src/index.js";

// Cross-port VIP scenario: vipDiscount supersedes generalDiscount. The superseder
// excludes only while its effective applicability a_eff > 0 — except gates count,
// counter rebuttals do not.
const generalDiscount: RuleFunc = () => [true, null];
const vipDiscount: RuleFunc = () => [true, null];
const suspendedVIP: RuleFunc = () => [true, null];
const vipDisputed: RuleFunc = () => [true, null];

function entry(entries: TraceEntry[], prefix: string): TraceEntry {
  const e = entries.find(t => t.name.startsWith(prefix));
  if (!e) throw new Error(`no entry for ${prefix}`);
  return e;
}

function vipGraph(wire: (g: Graph, vip: Rule) => void): Graph {
  const g = new Graph("vip");
  const general = g.rule(generalDiscount);
  const vip = g.rule(vipDiscount);
  vip.supersedes(general);
  wire(g, vip);
  return g;
}

describe("Supersedes — superseder excludes only while a_eff > 0", () => {
  const cases: Array<[string, (g: Graph, vip: Rule) => void, number, number, boolean]> = [
    ["no exception", () => {}, 1, 1, true],
    ["exception fully gates", (g, vip) => g.except(suspendedVIP).attacks(vip), -1, 0, false],
    ["exception force 0.5", (g, vip) => g.except(suspendedVIP).attacks(vip, 0.5), 0, 0.5, true],
    ["counter keeps priority", (g, vip) => g.counter(vipDisputed).attacks(vip), 0, 1, true],
  ];

  for (const [name, wire, vipVerdict, vipAEff, wantExcluded] of cases) {
    it(name, () => {
      const g = vipGraph(wire);
      const { results, trace } = g.run(newContext());
      const ge = entry(trace.all(), "generalDiscount");
      const ve = entry(trace.all(), "vipDiscount");
      expect(ge.excluded).toBe(wantExcluded);
      expect(ge.verdict).toBe(1);                    // reporting-only
      expect(ve.excluded).toBe(false);
      expect(ve.activated).toBe(true);
      expect(ve.verdict).toBe(vipVerdict);
      expect(results.find(r => r.name.startsWith("vipDiscount"))!.applicability).toBe(vipAEff);

      const traced = g.evaluate(newContext(), { trace: true });
      const gw = traced.find(r => r.name.startsWith("generalDiscount"))!;
      expect(entry(gw.trace!, "generalDiscount").excluded).toBe(wantExcluded);
    });
  }

  it("strict superseder ignores the exception gate (a_eff = a)", () => {
    const g = vipGraph((g, vip) => g.except(suspendedVIP).attacks(vip));
    g.rules[1].strength = Strength.Strict;
    const { trace } = g.run(newContext());
    expect(entry(trace.all(), "generalDiscount").excluded).toBe(true);
  });
});

describe("Supersedes — chain and cycle", () => {
  const x: RuleFunc = () => [true, null];
  const y: RuleFunc = () => [true, null];
  const z: RuleFunc = () => [true, null];

  function chain(gateX: boolean): TraceEntry[] {
    const g = new Graph("chain");
    const rz = g.rule(z);
    const ry = g.rule(y);
    const rx = g.rule(x);
    rx.supersedes(ry);
    ry.supersedes(rz);
    if (gateX) g.except(suspendedVIP).attacks(rx);
    return g.run(newContext()).trace.all();
  }

  it("an excluded superseder still excludes its own target", () => {
    const t = chain(false);
    expect([entry(t, "x").excluded, entry(t, "y").excluded, entry(t, "z").excluded]).toEqual([false, true, true]);
  });

  it("a fully gated superseder loses only its own exclusion", () => {
    const t = chain(true);
    expect([entry(t, "x").excluded, entry(t, "y").excluded, entry(t, "z").excluded]).toEqual([false, false, true]);
  });

  it("a supersede cycle is rejected", () => {
    const g = new Graph("cycle");
    const ra = g.rule(x);
    const rb = g.rule(y);
    ra.supersedes(rb);
    rb.supersedes(ra);
    expect(() => g.evaluate(newContext())).toThrow(/supersede cycle detected/);
  });

  it("with() rewrites supersede edges to the new id", () => {
    const spec = { specName: () => "s", validate: () => {} };
    const g = new Graph("with");
    const general = g.rule(generalDiscount);
    const vip = g.rule(vipDiscount);
    vip.supersedes(general);
    vip.with(spec);
    general.with(spec);
    expect(g.supersedes).toEqual([{ from: vip.id, to: general.id }]);
    expect(entry(g.run(newContext()).trace.all(), "generalDiscount").excluded).toBe(true);
  });
});

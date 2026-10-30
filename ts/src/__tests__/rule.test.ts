import { describe, it, expect } from "vitest";
import { Graph } from "../graph.js";
import type { RuleFunc } from "../types.js";

const always: RuleFunc = () => [true, null];

describe("Rule.run", () => {
  it("wires a non-null sub-graph and returns the same rule (fluent)", () => {
    const sub = new Graph("sub");
    sub.rule(always);

    const parent = new Graph("parent");
    const r = parent.rule(always);
    const returned = r.run(sub);

    expect(returned).toBe(r);                       // returns this for chaining
    expect(parent.rules[0].runGraph).toBe(sub);     // runGraph recorded on the meta
  });

  it("throws when the sub-graph is null", () => {
    const g = new Graph("null-sub");
    const r = g.rule(always);
    // @ts-expect-error — intentional null to exercise the guard
    expect(() => r.run(null)).toThrow(/non-null sub-graph/);
  });
});

describe("Rule.attacks — force", () => {
  const warrant: RuleFunc = () => [true, null];
  const rebuttal: RuleFunc = () => [true, null];

  it("records force 1 when omitted", () => {
    const g = new Graph("default");
    const w = g.rule(warrant);
    g.counter(rebuttal).attacks(w);
    expect(g.defeats).toEqual([{ from: g.rules[1].name, to: w.id, force: 1 }]);
  });

  it("rejects more than one value, NaN, and values outside [0, 1]", () => {
    const g = new Graph("validate");
    const w = g.rule(warrant);
    const c = g.counter(rebuttal);
    expect(() => c.attacks(w, 0.5, 0.5)).toThrow("at most one force value");
    for (const bad of [NaN, -0.1, 1.1]) {
      expect(() => c.attacks(w, bad)).toThrow("force must be between 0.0 and 1.0");
    }
    expect(g.defeats).toHaveLength(0);
    c.attacks(w, 0);
    c.attacks(w, 1);
    expect(g.defeats.map(d => d.force)).toEqual([0, 1]);
  });
});

describe("qualifier validation", () => {
  it("accepts values in [0, 1] and returns the rule for chaining", () => {
    const g = new Graph("q-ok");
    const r = g.rule(always);
    expect(r.qualifier(0)).toBe(r);
    expect(g.rules[0].qualifier).toBe(0);
    r.qualifier(1);
    expect(g.rules[0].qualifier).toBe(1);
  });

  it("rejects values below 0 and above 1", () => {
    const g = new Graph("q-range");
    const r = g.rule(always).qualifier(0.6);
    expect(() => r.qualifier(-0.1)).toThrow("qualifier must be between 0.0 and 1.0");
    expect(() => r.qualifier(1.1)).toThrow("qualifier must be between 0.0 and 1.0");
    expect(g.rules[0].qualifier).toBe(0.6);
  });

  it("rejects NaN", () => {
    const g = new Graph("q-nan");
    expect(() => g.rule(always).qualifier(NaN)).toThrow("qualifier must be between 0.0 and 1.0");
    expect(g.rules[0].qualifier).toBe(1.0);
  });
});

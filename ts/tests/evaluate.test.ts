import { describe, it, expect } from "vitest";
import { Graph, newContext, type RuleFunc } from "../src/index.js";

describe("Evaluate", () => {
  it("single warrant — verdict 1.0", () => {
    const g = new Graph("test");
    const fn: RuleFunc = (ctx, specs) => [true, null];
    g.rule(fn);
    const results = g.evaluate(newContext());
    expect(results).toHaveLength(1);
    expect(results[0].verdict).toBe(1.0);
  });

  it("warrant + counter — verdict 0.0", () => {
    const g = new Graph("test");
    const wFn: RuleFunc = (ctx, specs) => [true, null];
    const rFn: RuleFunc = (ctx, specs) => [true, null];
    const w = g.rule(wFn);
    const r = g.counter(rFn);
    r.attacks(w);
    const results = g.evaluate(newContext());
    expect(results[0].verdict).toBe(0.0);
  });

  it("compensation chain — verdict 1/3", () => {
    const g = new Graph("test");
    const wFn: RuleFunc = (ctx, specs) => [true, null];
    const rFn: RuleFunc = (ctx, specs) => [true, null];
    const dFn: RuleFunc = (ctx, specs) => [true, null];
    const w = g.rule(wFn);
    const r = g.counter(rFn);
    const d = g.except(dFn);
    r.attacks(w);
    d.attacks(r);
    const results = g.evaluate(newContext());
    expect(results[0].verdict).toBeCloseTo(1 / 3, 5);
  });

  it("diamond DAG — shared attacker", () => {
    const g = new Graph("test");
    const wFn: RuleFunc = (ctx, specs) => [true, null];
    const r1Fn: RuleFunc = (ctx, specs) => [true, null];
    const r2Fn: RuleFunc = (ctx, specs) => [true, null];
    const dFn: RuleFunc = (ctx, specs) => [true, null];
    const w = g.rule(wFn);
    const r1 = g.counter(r1Fn);
    const r2 = g.counter(r2Fn);
    const d = g.except(dFn);
    r1.attacks(w);
    r2.attacks(w);
    d.attacks(r1);
    d.attacks(r2);
    const results = g.evaluate(newContext());
    expect(results[0].verdict).toBe(0.0);
  });

  it("inactive rule — no result", () => {
    const g = new Graph("test");
    const fn: RuleFunc = (ctx, specs) => [false, null];
    g.rule(fn);
    expect(g.evaluate(newContext())).toHaveLength(0);
  });

  it("cycle detection", () => {
    const g = new Graph("test");
    const aFn: RuleFunc = (ctx, specs) => [true, null];
    const bFn: RuleFunc = (ctx, specs) => [true, null];
    const a = g.rule(aFn);
    const b = g.counter(bFn);
    a.attacks(b);
    b.attacks(a);
    expect(() => g.evaluate(newContext())).toThrow("cycle");
  });

  it("null ctx throws", () => {
    const g = new Graph("test");
    const fn: RuleFunc = (ctx, specs) => [true, null];
    g.rule(fn);
    expect(() => g.evaluate(null as any)).toThrow("ctx");
  });

  it("rule throw → error", () => {
    const g = new Graph("test");
    const fn: RuleFunc = () => { throw new Error("boom"); };
    g.rule(fn);
    expect(() => g.evaluate(newContext())).toThrow("boom");
  });

  it("qualifier 0.0 — verdict -1.0", () => {
    const g = new Graph("test");
    const fn: RuleFunc = (ctx, specs) => [true, null];
    g.rule(fn).qualifier(0.0);
    const results = g.evaluate(newContext());
    expect(results[0].verdict).toBe(-1.0);
  });

  it("evaluate twice — idempotent", () => {
    const g = new Graph("test");
    const fn: RuleFunc = (ctx, specs) => [true, null];
    g.rule(fn);
    const r1 = g.evaluate(newContext());
    const r2 = g.evaluate(newContext());
    expect(r1[0].verdict).toBe(r2[0].verdict);
  });
});

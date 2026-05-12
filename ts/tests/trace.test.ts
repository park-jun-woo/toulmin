import { describe, it, expect } from "vitest";
import { Graph, newContext, type RuleFunc } from "../src/index.js";

describe("Trace", () => {
  it("collects trace entries", () => {
    const g = new Graph("test");
    const wFn: RuleFunc = (ctx, specs) => [true, "warrant-evidence"];
    const rFn: RuleFunc = (ctx, specs) => [true, "counter-evidence"];
    const w = g.rule(wFn);
    const r = g.counter(rFn);
    r.attacks(w);
    const results = g.evaluate(newContext(), { trace: true });
    expect(results[0].trace).toHaveLength(2);
    expect(results[0].trace![0].role).toBe("rule");
    expect(results[0].trace![1].role).toBe("counter");
    expect(results[0].trace![0].activated).toBe(true);
  });

  it("duration measures time", () => {
    const g = new Graph("test");
    const fn: RuleFunc = (ctx, specs) => [true, null];
    g.rule(fn);
    const results = g.evaluate(newContext(), { duration: true });
    expect(results[0].trace).toBeDefined();
    expect(results[0].trace![0].duration).toBeGreaterThanOrEqual(0);
  });

  it("inactive rule appears in trace", () => {
    const g = new Graph("test");
    const wFn: RuleFunc = (ctx, specs) => [true, null];
    const rFn: RuleFunc = (ctx, specs) => [false, null];
    const w = g.rule(wFn);
    const r = g.counter(rFn);
    r.attacks(w);
    const results = g.evaluate(newContext(), { trace: true });
    expect(results[0].trace).toHaveLength(2);
    expect(results[0].trace![1].activated).toBe(false);
  });
});

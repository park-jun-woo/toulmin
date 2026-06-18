import { describe, it, expect } from "vitest";
import {
  Graph,
  newContext,
  detectRunCycle,
  type RuleFunc,
  type Context,
  type NodeEvent,
} from "../src/index.js";

const always: RuleFunc = () => [true, null];

describe("Run composition — rule.run(subGraph)", () => {
  it("basic: Active node → sub-graph Run (ctx side effect)", () => {
    const parentRule: RuleFunc = () => [true, null];
    const subRule: RuleFunc = () => [true, null];

    const sub = new Graph("sub");
    sub.rule(subRule).onActive((ctx) => { ctx.set("subRan", true); });

    const parent = new Graph("parent");
    parent.rule(parentRule).run(sub);

    const ctx = newContext();
    parent.run(ctx);

    expect(ctx.get("subRan")).toBe(true);
  });

  it("Active-only: Defeated node does not run its runGraph", () => {
    const warrant: RuleFunc = () => [true, null];
    const counter: RuleFunc = () => [true, null];
    const subRule: RuleFunc = () => [true, null];

    const sub = new Graph("sub-defeated");
    sub.rule(subRule).onActive((ctx) => { ctx.set("subRan", true); });

    const parent = new Graph("parent-defeated");
    const w = parent.rule(warrant);
    const c = parent.counter(counter);
    c.attacks(w);
    w.run(sub);   // w is Defeated (verdict 0) → must NOT run sub

    const ctx = newContext();
    parent.run(ctx);

    expect(ctx.get("subRan")).toBeUndefined();
  });

  it("Active-only: Inactive node does not run its runGraph", () => {
    const inactive: RuleFunc = () => [false, null];
    const subRule: RuleFunc = () => [true, null];

    const sub = new Graph("sub-inactive");
    sub.rule(subRule).onActive((ctx) => { ctx.set("subRan", true); });

    const parent = new Graph("parent-inactive");
    parent.rule(inactive).run(sub);

    const ctx = newContext();
    parent.run(ctx);

    expect(ctx.get("subRan")).toBeUndefined();
  });

  it("handler-then-sub ordering: onActive fires before sub Run", () => {
    const log: string[] = [];
    const parentRule: RuleFunc = () => [true, null];
    const subRule: RuleFunc = () => [true, null];

    const sub = new Graph("sub-order");
    sub.rule(subRule).onActive(() => { log.push("sub"); });

    const parent = new Graph("parent-order");
    parent.rule(parentRule).onActive(() => { log.push("handler"); }).run(sub);

    parent.run(newContext());

    expect(log).toEqual(["handler", "sub"]);
  });

  it("ctx flows down: sub rule reads value set by parent handler", () => {
    let seen: unknown;
    const parentRule: RuleFunc = () => [true, null];
    const subRule: RuleFunc = (ctx) => [true, ctx.get("token")];

    const sub = new Graph("sub-ctx");
    sub.rule(subRule).onActive((_ctx, ev) => { seen = ev.evidence; });

    const parent = new Graph("parent-ctx");
    parent.rule(parentRule).onActive((ctx) => { ctx.set("token", "from-parent"); }).run(sub);

    parent.run(newContext());

    expect(seen).toBe("from-parent");
  });

  it("cycle A→B→A → run() throws", () => {
    const ruleA: RuleFunc = () => [true, null];
    const ruleB: RuleFunc = () => [true, null];

    const ga = new Graph("cycle-A");
    const gb = new Graph("cycle-B");
    ga.rule(ruleA).run(gb);
    gb.rule(ruleB).run(ga);

    const ctx = newContext();
    expect(() => ga.run(ctx)).toThrow(/run cycle detected/);
    // detectRunCycle agrees directly
    expect(detectRunCycle(ga)).not.toBeNull();
  });

  it("depth guard: long acyclic chain exceeds runMaxDepth → throws", () => {
    const chain: Graph[] = [];
    for (let i = 0; i < 70; i++) {
      // registration is per-graph, so reusing one fn across graphs is fine
      const g = new Graph(`chain-${i}`);
      g.rule(always);
      chain.push(g);
    }
    for (let i = 0; i < chain.length - 1; i++) {
      chain[i].rules[0].runGraph = chain[i + 1];   // acyclic chain via direct wiring
    }
    // sanity: acyclic, so detectRunCycle passes — only the depth backstop fires
    expect(detectRunCycle(chain[0])).toBeNull();
    expect(() => chain[0].run(newContext())).toThrow(/run depth exceeded/);
  });

  it("verdict isolation: sub verdicts do not affect parent results", () => {
    const parentRule: RuleFunc = () => [true, null];
    const subWarrant: RuleFunc = () => [true, null];
    const subCounter: RuleFunc = () => [true, null];

    const sub = new Graph("sub-iso");
    const sw = sub.rule(subWarrant);
    const sc = sub.counter(subCounter);
    sc.attacks(sw);   // sub warrant is Defeated (verdict 0)

    const parent = new Graph("parent-iso");
    parent.rule(parentRule).run(sub);

    const ctx = newContext();
    const evalResults = parent.evaluate(ctx);
    const runResult = parent.run(ctx);

    expect(runResult.results.map(r => r.verdict)).toEqual(evalResults.map(r => r.verdict));
    expect(runResult.results).toHaveLength(1);
    expect(runResult.results[0].verdict).toBeCloseTo(1, 5);
  });

  it("diamond DAG is legal: shared sub-graph Run via two paths", () => {
    let bottomRuns = 0;
    const topLeft: RuleFunc = () => [true, null];
    const topRight: RuleFunc = () => [true, null];
    const leftRule: RuleFunc = () => [true, null];
    const rightRule: RuleFunc = () => [true, null];
    const bottomRule: RuleFunc = () => [true, null];

    const bottom = new Graph("diamond-bottom");
    bottom.rule(bottomRule).onActive(() => { bottomRuns++; });

    const left = new Graph("diamond-left");
    left.rule(leftRule).run(bottom);

    const right = new Graph("diamond-right");
    right.rule(rightRule).run(bottom);

    const top = new Graph("diamond-top");
    top.rule(topLeft).run(left);
    top.rule(topRight).run(right);

    expect(detectRunCycle(top)).toBeNull();
    top.run(newContext());

    expect(bottomRuns).toBe(2);   // run twice (once per path), still a legal DAG
  });

  it("error propagation: sub handler throw wrapped as run \"...\" → \"...\"", () => {
    const parentRule: RuleFunc = () => [true, null];
    const subRule: RuleFunc = () => [true, null];

    const sub = new Graph("sub-boom");
    sub.rule(subRule).onActive(() => { throw new Error("boom"); });

    const parent = new Graph("parent-boom");
    parent.rule(parentRule).run(sub);

    const ctx = newContext();
    expect(() => parent.run(ctx)).toThrow(/run "/);
    expect(() => parent.run(ctx)).toThrow(/→ "sub-boom"/);
    expect(() => parent.run(ctx)).toThrow(/boom/);
  });

  it("run(null) is a registration error", () => {
    const g = new Graph("null-sub");
    const r = g.rule(always);
    // @ts-expect-error — intentional null for guard
    expect(() => r.run(null)).toThrow(/non-null sub-graph/);
  });
});

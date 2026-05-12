import { describe, it, expect } from "vitest";
import { Graph, newContext, type RuleFunc, type Spec } from "../src/index.js";

class TestSpec implements Spec {
  constructor(public value: string) {}
  specName() { return "TestSpec"; }
  validate() {}
}

const isAuth: RuleFunc = (ctx, specs) => [ctx.get("user") != null, null];
const isBlocked: RuleFunc = (ctx, specs) => [true, null];

describe("Graph construction", () => {
  it("registers rule, counter, except", () => {
    const g = new Graph("test");
    const auth = g.rule(isAuth);
    const blocked = g.counter(isBlocked);
    blocked.attacks(auth);
    expect(g.rules).toHaveLength(2);
    expect(g.defeats).toHaveLength(1);
  });

  it("throws on duplicate registration", () => {
    const g = new Graph("test");
    g.rule(isAuth);
    expect(() => g.rule(isAuth)).toThrow("duplicate");
  });

  it("with() updates defeat edges", () => {
    const fn: RuleFunc = (ctx, specs) => [true, null];
    const g = new Graph("test");
    const r1 = g.rule(fn);
    const r2 = g.counter(isBlocked);
    r2.attacks(r1);
    r1.with(new TestSpec("admin"));
    expect(g.defeats[0].to).toContain("#");
  });

  it("same func different spec — separate rules", () => {
    const fn: RuleFunc = (ctx, specs) => [true, null];
    const g = new Graph("test");
    g.rule(fn).with(new TestSpec("admin"));
    g.rule(fn).with(new TestSpec("editor"));
    expect(g.rules).toHaveLength(2);
    expect(g.rules[0].name).not.toBe(g.rules[1].name);
  });

  it("qualifier validates range", () => {
    const fn1: RuleFunc = (ctx, specs) => [true, null];
    const fn2: RuleFunc = (ctx, specs) => [true, null];
    const g = new Graph("test");
    expect(() => g.rule(fn1).qualifier(1.5)).toThrow("qualifier");
    expect(() => g.rule(fn2).qualifier(-0.1)).toThrow("qualifier");
  });
});

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

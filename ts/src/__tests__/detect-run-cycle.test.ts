import { describe, it, expect } from "vitest";
import { Graph } from "../graph.js";
import { detectRunCycle } from "../detect-run-cycle.js";
import type { RuleFunc } from "../types.js";

const always: RuleFunc = () => [true, null];

describe("detectRunCycle", () => {
  it("acyclic chain returns null (covers !sub continue + WHITE recurse + BLACK done)", () => {
    // root has two rules: one without a runGraph (→ !sub continue),
    // one wired to a sub-graph (→ WHITE recurse → BLACK).
    const sub = new Graph("sub");
    sub.rule(always); // leaf, no runGraph

    const root = new Graph("root");
    root.rule(always);            // no runGraph → !sub continue
    root.rule(() => [false, null]).run(sub);

    expect(detectRunCycle(root)).toBeNull();
  });

  it("direct self-cycle A→A returns an Error (GRAY re-entry)", () => {
    const g = new Graph("self");
    const r = g.rule(always);
    r.run(g); // points back to itself

    const err = detectRunCycle(g);
    expect(err).not.toBeNull();
    expect(err!.message).toMatch(/run cycle detected at graph "self"/);
  });

  it("indirect cycle A→B→A propagates the Error up the recursion", () => {
    const ga = new Graph("cycle-A");
    const gb = new Graph("cycle-B");
    ga.rule(always).run(gb);
    gb.rule(always).run(ga); // back-edge to GRAY A

    const err = detectRunCycle(ga);
    expect(err).not.toBeNull();
    expect(err!.message).toMatch(/run cycle detected/);
  });

  it("diamond DAG is legal — shared BLACK sub-graph reached twice returns null", () => {
    const bottom = new Graph("d-bottom");
    bottom.rule(always);

    const left = new Graph("d-left");
    left.rule(always).run(bottom);

    const right = new Graph("d-right");
    right.rule(always).run(bottom);

    const top = new Graph("d-top");
    top.rule(always).run(left);
    top.rule(() => [false, null]).run(right); // second path reaches BLACK bottom → skip

    expect(detectRunCycle(top)).toBeNull();
  });
});

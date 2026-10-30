import type { Spec, NodeHandler } from "./types.js";
import type { Graph } from "./graph.js";
import type { RuleFunc } from "./types.js";
import { ruleID } from "./rule-id.js";
import { attackForce } from "./attack-force.js";

export class Rule {
  id: string;
  private graph: Graph;
  private idx: number;
  private fn: RuleFunc;

  constructor(id: string, graph: Graph, idx: number, fn: RuleFunc) {
    this.id = id; this.graph = graph; this.idx = idx; this.fn = fn;
  }

  /**
   * Declares a defeat edge: this rule attacks target.
   * The optional force F in [0, 1] (default 1) scales the attack: a rebuttal
   * (counter or rule attacker) adds F·(v+1)/2 to the target's denominator, and an
   * exception (except) multiplies the target's applicability by 1 − F·(v+1)/2.
   * Throws if more than one force is given, or force is NaN or outside [0, 1].
   */
  attacks(target: Rule, ...force: number[]): void {
    const f = attackForce(force);
    this.graph.defeats.push({ from: this.id, to: target.id, force: f });
  }

  /**
   * Declares a supersede edge: while this rule's effective applicability
   * a_eff > 0 — its own (v+1)/2 after every except gate; counter rebuttals do
   * not count, strict rules keep a_eff = (v+1)/2 — target is crisply excluded
   * from the tally (TraceEntry.excluded). An exception that fully gates this
   * rule lifts the exclusion. The edge never enters the h-Categoriser traversal.
   */
  supersedes(target: Rule): void {
    this.graph.supersedes.push({ from: this.id, to: target.id });
  }

  with(spec: Spec): Rule {
    spec.validate();
    const meta = this.graph.rules[this.idx];
    meta.specs.push(spec);
    const oldID = meta.name;
    const newID = ruleID(this.fn, meta.specs);
    meta.name = newID;
    this.id = newID;
    const role = this.graph.roles.get(oldID);
    if (role !== undefined) {
      this.graph.roles.delete(oldID);
      this.graph.roles.set(newID, role);
    }
    for (const edge of [...this.graph.defeats, ...this.graph.supersedes]) {
      if (edge.from === oldID) edge.from = newID;
      if (edge.to === oldID) edge.to = newID;
    }
    return this;
  }

  qualifier(q: number): Rule {
    if (Number.isNaN(q) || q < 0.0 || q > 1.0) throw new Error("qualifier must be between 0.0 and 1.0");
    this.graph.rules[this.idx].qualifier = q;
    return this;
  }

  runOn(h: NodeHandler): Rule {
    this.graph.rules[this.idx].runOn = h;
    return this;
  }

  run(g: Graph): Rule {
    if (g == null) throw new Error("toulmin: run requires a non-null sub-graph");
    this.graph.rules[this.idx].runGraph = g;
    return this;
  }
}

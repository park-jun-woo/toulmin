import type { Spec, NodeHandler } from "./types.js";
import type { Graph } from "./graph.js";
import type { RuleFunc } from "./types.js";
import { ruleID } from "./rule-id.js";

export class Rule {
  id: string;
  private graph: Graph;
  private idx: number;
  private fn: RuleFunc;

  constructor(id: string, graph: Graph, idx: number, fn: RuleFunc) {
    this.id = id; this.graph = graph; this.idx = idx; this.fn = fn;
  }

  attacks(target: Rule): void {
    this.graph.defeats.push({ from: this.id, to: target.id });
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
    for (const edge of this.graph.defeats) {
      if (edge.from === oldID) edge.from = newID;
      if (edge.to === oldID) edge.to = newID;
    }
    return this;
  }

  qualifier(q: number): Rule {
    if (q < 0.0 || q > 1.0) throw new Error("qualifier must be between 0.0 and 1.0");
    this.graph.rules[this.idx].qualifier = q;
    return this;
  }

  onActive(h: NodeHandler): Rule {
    this.graph.rules[this.idx].onActive = h;
    return this;
  }

  onDefeated(h: NodeHandler): Rule {
    this.graph.rules[this.idx].onDefeated = h;
    return this;
  }

  onInactive(h: NodeHandler): Rule {
    this.graph.rules[this.idx].onInactive = h;
    return this;
  }
}

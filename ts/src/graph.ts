import type { Context, RuleFunc, EvalOption, EvalResult, Specs, TraceEntry } from "./types.js";
import { Strength, EvalMethod } from "./types.js";
import type { DefeatEdge } from "./defeat-edge.js";
import type { RuleMeta } from "./rule-meta.js";
import { Rule } from "./rule.js";
import { ruleID } from "./rule-id.js";
import { shortName } from "./short-name.js";
import { detectCycle } from "./detect-cycle.js";

export class Graph {
  readonly name: string;
  rules: RuleMeta[] = [];
  roles = new Map<string, string>();
  defeats: DefeatEdge[] = [];

  constructor(name: string) { this.name = name; }

  rule(fn: RuleFunc): Rule { return this._register(fn, Strength.Defeasible, "rule"); }
  counter(fn: RuleFunc): Rule { return this._register(fn, Strength.Defeasible, "counter"); }
  except(fn: RuleFunc): Rule { return this._register(fn, Strength.Defeater, "except"); }

  private _register(fn: RuleFunc, strength: Strength, role: string): Rule {
    const id = ruleID(fn, []);
    if (this.roles.has(id)) throw new Error(`duplicate rule registration: ${id}`);
    const idx = this.rules.length;
    this.rules.push({ name: id, qualifier: 1.0, strength, defeats: [], specs: [], fn });
    this.roles.set(id, role);
    return new Rule(id, this, idx, fn);
  }

  evaluate(ctx: Context, option?: EvalOption): EvalResult[] {
    if (ctx == null) throw new Error("ctx must not be null");
    const opt = resolveOption(option);

    // Build eval context
    const fnMap = new Map<string, RuleFunc>();
    const qualMap = new Map<string, number>();
    const strMap = new Map<string, Strength>();
    const specsMap = new Map<string, Specs>();
    const edges = new Map<string, string[]>();
    const roleMap = new Map(this.roles);

    for (const r of this.rules) {
      fnMap.set(r.name, r.fn);
      qualMap.set(r.name, r.qualifier);
      strMap.set(r.name, r.strength);
      specsMap.set(r.name, r.specs);
    }
    for (const d of this.defeats) {
      const list = edges.get(d.to) ?? [];
      list.push(d.from);
      edges.set(d.to, list);
    }

    // Cycle detection
    const cycleErr = detectCycle(edges);
    if (cycleErr) throw cycleErr;

    // Build attacker set
    const attackerSet = new Set<string>();
    for (const attackers of edges.values()) {
      for (const a of attackers) attackerSet.add(a);
    }

    // Evaluation state
    const ran = new Set<string>();
    const active = new Map<string, boolean>();
    const evidence = new Map<string, unknown>();
    const verdictCache = new Map<string, number>();
    let trace: TraceEntry[] = [];
    let err: Error | null = null;

    function isWarrant(name: string): boolean {
      const str = strMap.get(name) ?? Strength.Defeasible;
      if (str === Strength.Defeater) return false;
      return !attackerSet.has(name);
    }

    function inferRole(id: string): string {
      if (strMap.get(id) === Strength.Defeater) return "except";
      if (attackerSet.has(id)) return "counter";
      return "rule";
    }

    function calc(id: string): number {
      if (err) return -1.0;
      const cached = verdictCache.get(id);
      if (cached !== undefined) return cached;
      const fn = fnMap.get(id);
      if (!fn) return -1.0;
      if (!ran.has(id)) {
        ran.add(id);
        try {
          const [act, ev] = fn(ctx, specsMap.get(id) ?? []);
          active.set(id, act);
          evidence.set(id, ev);
        } catch (e) {
          err = new Error(`rule "${id}": ${e}`);
          return -1.0;
        }
      }
      if (!active.get(id)) return -1.0;
      let sum = 0;
      if (strMap.get(id) !== Strength.Strict) {
        for (const aid of edges.get(id) ?? []) {
          sum += (calc(aid) + 1.0) / 2.0;
        }
      }
      const q = qualMap.get(id) ?? 1.0;
      const v = 2 * (q / (1 + sum)) - 1;
      verdictCache.set(id, v);
      return v;
    }

    function calcTrace(id: string, withDuration: boolean): number {
      if (err) return -1.0;
      const cached = verdictCache.get(id);
      if (cached !== undefined) return cached;
      const fn = fnMap.get(id);
      if (!fn) return -1.0;
      if (!ran.has(id)) {
        ran.add(id);
        const start = withDuration ? performance.now() : 0;
        try {
          const [act, ev] = fn(ctx, specsMap.get(id) ?? []);
          active.set(id, act);
          evidence.set(id, ev);
        } catch (e) {
          err = new Error(`rule "${id}": ${e}`);
          return -1.0;
        }
        const elapsed = withDuration ? performance.now() - start : undefined;
        trace.push({
          name: id,
          role: inferRole(id),
          activated: active.get(id) ?? false,
          qualifier: qualMap.get(id) ?? 1.0,
          evidence: evidence.get(id),
          specs: specsMap.get(id),
          duration: elapsed,
        });
      }
      if (!active.get(id)) return -1.0;
      let sum = 0;
      if (strMap.get(id) !== Strength.Strict) {
        for (const aid of edges.get(id) ?? []) {
          sum += (calcTrace(aid, withDuration) + 1.0) / 2.0;
        }
      }
      const q = qualMap.get(id) ?? 1.0;
      const v = 2 * (q / (1 + sum)) - 1;
      verdictCache.set(id, v);
      return v;
    }

    function reset(): void {
      ran.clear();
      active.clear();
      evidence.clear();
      verdictCache.clear();
      trace = [];
      err = null;
    }

    // Main evaluation loop
    const results: EvalResult[] = [];
    for (const r of this.rules) {
      if (!isWarrant(r.name)) continue;
      if (opt.trace) reset();

      const verdict = opt.trace
        ? calcTrace(r.name, opt.duration)
        : calc(r.name);

      if (err) throw err;
      if (!active.get(r.name)) continue;

      const result: EvalResult = {
        name: shortName(r.name),
        verdict,
        evidence: evidence.get(r.name),
      };
      if (opt.trace) {
        result.trace = trace.map(t => ({ ...t, name: shortName(t.name) }));
      }
      results.push(result);
    }
    return results;
  }
}

function resolveOption(opt?: EvalOption): { method: EvalMethod; trace: boolean; duration: boolean } {
  const resolved = {
    method: opt?.method ?? EvalMethod.Matrix,
    trace: opt?.trace ?? false,
    duration: opt?.duration ?? false,
  };
  if (resolved.duration) resolved.trace = true;
  return resolved;
}

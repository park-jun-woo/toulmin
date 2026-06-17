import type { Context, RuleFunc, EvalOption, EvalResult, Specs, TraceEntry, NodeEvent, NodeHandler, RunResult } from "./types.js";
import { Strength, EvalMethod, NodeEventType } from "./types.js";
import type { DefeatEdge } from "./defeat-edge.js";
import type { RuleMeta } from "./rule-meta.js";
import { Rule } from "./rule.js";
import { ruleID } from "./rule-id.js";
import { shortName } from "./short-name.js";
import { detectCycle } from "./detect-cycle.js";
import { createRunView, buildAttackerEvents } from "./run-view.js";

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
    return this._evaluate(ctx, resolveOption(option), false).results;
  }

  run(ctx: Context, option?: EvalOption): RunResult {
    if (ctx == null) throw new Error("ctx must not be null");
    const opt = resolveOption(option);
    opt.trace = false;      // 강제 off — full pass는 공유·비초기화 상태에서만 정확
    opt.duration = false;
    const st = this._evaluate(ctx, opt, true);   // full pass

    // 디스패치 전 1회 — 불변 스냅샷
    const order: NodeEvent[] = [];
    const byName = new Map<string, NodeEvent>();
    for (const r of this.rules) {
      const type = classifyEvent(st.active.get(r.name) ?? false, st.verdictCache.get(r.name) ?? 0);
      const ne: NodeEvent = {
        name: shortName(r.name),
        role: this.roles.get(r.name) ?? "rule",
        type,
        verdict: st.verdictCache.get(r.name) ?? 0,
        evidence: st.evidence.get(r.name),
      };
      order.push(ne);
      byName.set(ne.name, ne);
    }
    const attackers = buildAttackerEvents(this.defeats, byName);
    const view = createRunView(order, byName, attackers);

    for (let i = 0; i < order.length; i++) {
      const ne = order[i];
      const meta = this.rules[i];          // 인덱스 대응 — shortName 매핑 불필요
      const h = selectHandler(meta, ne.type);
      if (!h) continue;
      try {
        h(ctx, ne, view);
      } catch (e) {
        throw new Error(`handler "${shortName(meta.name)}" (${NodeEventType[ne.type]}): ${e}`);
      }
    }
    return { results: st.results, view };
  }

  private _evaluate(ctx: Context, opt: { method: EvalMethod; trace: boolean; duration: boolean }, full: boolean): {
    results: EvalResult[];
    active: Map<string, boolean>;
    verdictCache: Map<string, number>;
    evidence: Map<string, unknown>;
  } {
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

    if (full) {                                // run은 trace를 강제 off하므로 calc(non-trace)만 탐
      for (const r of this.rules) {
        if (!ran.has(r.name)) calc(r.name);    // 미도달 노드까지 평가
      }
      if (err) throw err;
    }

    return { results, active, verdictCache, evidence };
  }
}

function classifyEvent(active: boolean, verdict: number): NodeEventType {
  if (!active) return NodeEventType.Inactive;
  return verdict > 0 ? NodeEventType.Active : NodeEventType.Defeated;
}

function selectHandler(r: RuleMeta, t: NodeEventType): NodeHandler | undefined {
  switch (t) {
    case NodeEventType.Active: return r.onActive;
    case NodeEventType.Defeated: return r.onDefeated;
    default: return r.onInactive;
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

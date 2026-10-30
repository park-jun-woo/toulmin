import type { Context, RuleFunc, EvalOption, EvalResult, Specs, TraceEntry, RunResult } from "./types.js";
import { Strength, EvalMethod } from "./types.js";
import type { DefeatEdge, SupersedeEdge } from "./defeat-edge.js";
import type { RuleMeta } from "./rule-meta.js";
import { Rule } from "./rule.js";
import { ruleID } from "./rule-id.js";
import { shortName } from "./short-name.js";
import { detectCycle } from "./detect-cycle.js";
import { detectRunCycle } from "./detect-run-cycle.js";
import { createTrace } from "./trace.js";
import { safeCall } from "./safe-call.js";

const runMaxDepth = 64;   // 깊이 백스톱 (모듈 스코프)

export class Graph {
  readonly name: string;
  rules: RuleMeta[] = [];
  roles = new Map<string, string>();
  defeats: DefeatEdge[] = [];
  supersedes: SupersedeEdge[] = [];

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
    const err = detectRunCycle(this);       // ★ 최상위 1회
    if (err) throw err;
    const opt = resolveOption(option);
    opt.trace = false;      // 강제 off — full pass는 공유·비초기화 상태에서만 정확
    opt.duration = false;
    return this._runDepth(ctx, opt, 0);
  }

  private _runDepth(ctx: Context, opt: { method: EvalMethod; trace: boolean; duration: boolean }, depth: number): RunResult {
    if (depth > runMaxDepth) {
      throw new Error(`toulmin: run depth exceeded ${runMaxDepth} (possible runaway composition)`);
    }
    const st = this._evaluate(ctx, opt, true);   // full pass

    // full pass 상태(active/verdictCache/evidence)와 메타(specs/qualifier/role)로
    // 전 노드 TraceEntry[]를 등록 순서로 조립 (인덱스 i ↔ this.rules[i] 구조적 보장).
    const entries: TraceEntry[] = [];
    for (const r of this.rules) {
      entries.push({
        name: shortName(r.name),
        role: this.roles.get(r.name) ?? "rule",
        activated: st.active.get(r.name) ?? false,
        ruleVerdict: st.ruleVerdict.get(r.name) ?? -1.0,
        qualifier: r.qualifier,
        verdict: st.verdictCache.get(r.name) ?? 0,
        evidence: st.evidence.get(r.name),
        ground: ctx,
        specs: r.specs,
        excluded: st.excluded.has(r.name),
      });
    }

    const tr = createTrace(entries, ctx);

    for (let i = 0; i < entries.length; i++) {
      const meta = this.rules[i];          // 인덱스 대응
      // Active = activated && verdict > 0 인 노드만 발화/실행
      if (!(entries[i].activated && entries[i].verdict > 0)) continue;
      // (a) runOn 핸들러
      if (meta.runOn) {
        try {
          meta.runOn(entries[i], tr);
        } catch (e) {
          throw new Error(`runOn "${entries[i].name}": ${e}`);
        }
      }
      // (b) 실행 간선 — Active면 하위 그래프 Run (ctx 아래로, depth+1)
      if (meta.runGraph) {
        try {
          meta.runGraph._runDepth(ctx, opt, depth + 1);   // 같은 클래스 → private 접근 가능
        } catch (e) {
          throw new Error(`run "${entries[i].name}" → "${meta.runGraph.name}": ${e}`);
        }
      }
    }
    return { results: st.results, trace: tr };
  }

  private _evaluate(ctx: Context, opt: { method: EvalMethod; trace: boolean; duration: boolean }, full: boolean): {
    results: EvalResult[];
    active: Map<string, boolean>;
    ruleVerdict: Map<string, number>;
    verdictCache: Map<string, number>;
    evidence: Map<string, unknown>;
    excluded: Set<string>;
  } {
    // Build eval context
    const fnMap = new Map<string, RuleFunc>();
    const qualMap = new Map<string, number>();
    const strMap = new Map<string, Strength>();
    const specsMap = new Map<string, Specs>();
    const edges = new Map<string, string[]>();
    const forces = new Map<string, number[]>();   // forces[to][i] = force of edges[to][i]
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
      const fl = forces.get(d.to) ?? [];
      fl.push(d.force ?? 1.0);
      forces.set(d.to, fl);
    }

    // Cycle detection
    const cycleErr = detectCycle(edges);
    if (cycleErr) throw cycleErr;

    // Supersede edges: target → superseders. Never traversed; a cycle is an authoring error.
    const supersedeBy = new Map<string, string[]>();
    const supersedeAdj = new Map<string, string[]>();
    for (const s of this.supersedes) {
      supersedeBy.set(s.to, [...(supersedeBy.get(s.to) ?? []), s.from]);
      supersedeAdj.set(s.from, [...(supersedeAdj.get(s.from) ?? []), s.to]);
    }
    const supCycleErr = detectCycle(supersedeAdj);
    if (supCycleErr) throw new Error(`toulmin: supersede ${supCycleErr.message}`);

    // Build attacker set
    const attackerSet = new Set<string>();
    for (const attackers of edges.values()) {
      for (const a of attackers) attackerSet.add(a);
    }

    // Evaluation state
    const ran = new Set<string>();
    const active = new Map<string, boolean>();          // derived: ruleVerdict > -1
    const ruleVerdict = new Map<string, number>();      // rule function's own v in [-1, 1]
    const applicability = new Map<string, number>();    // a_eff
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

    // Exceptions (except / Defeater) gate the target's applicability multiplicatively;
    // every other attacker is a rebuttal in the denominator.
    function isExceptAttacker(id: string): boolean {
      return strMap.get(id) === Strength.Defeater || roleMap.get(id) === "except";
    }

    // Runs the rule function once; on error sets err and returns false.
    function runRule(id: string, fn: RuleFunc): boolean {
      const out = safeCall(fn, ctx, specsMap.get(id)!);   // every registered id has specs
      if (!out.ok) {
        err = new Error(`rule "${id}": ${out.error}`);
        return false;
      }
      ruleVerdict.set(id, out.verdict);
      active.set(id, out.verdict > -1.0);
      evidence.set(id, out.evidence);
      return true;
    }

    // a = (v+1)/2; S = Σ_rebuttal F·(v_j+1)/2; a_eff = a · Π_except (1 − F·(v_j+1)/2)
    // raw = q · a_eff / (1 + S); verdict = 2·raw − 1. Strict ignores every attacker.
    function combine(id: string, evalFn: (aid: string) => number): number {
      const a = (ruleVerdict.get(id)! + 1.0) / 2.0;   // only active (ran) nodes reach combine
      let sum = 0;
      let gate = 1.0;
      if (strMap.get(id) !== Strength.Strict) {
        const fl = forces.get(id) ?? [];
        (edges.get(id) ?? []).forEach((aid, i) => {
          const x = fl[i] * (evalFn(aid) + 1.0) / 2.0;   // forces is parallel to edges
          if (isExceptAttacker(aid)) gate *= 1.0 - x;
          else sum += x;
        });
      }
      const aEff = a * gate;
      applicability.set(id, aEff);
      const q = qualMap.get(id)!;
      const v = 2 * (q * aEff / (1 + sum)) - 1;
      verdictCache.set(id, v);
      return v;
    }

    // A node whose rule verdict is -1 (not applicable) returns -1 without
    // evaluating its attackers (lazy evaluation).
    function calc(id: string): number {
      if (err) return -1.0;
      const cached = verdictCache.get(id);
      if (cached !== undefined) return cached;
      const fn = fnMap.get(id);
      if (!fn) return -1.0;
      if (!ran.has(id)) {
        ran.add(id);
        if (!runRule(id, fn)) return -1.0;
      }
      if (!active.get(id)) return -1.0;
      return combine(id, calc);
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
        if (!runRule(id, fn)) return -1.0;
        const elapsed = withDuration ? performance.now() - start : undefined;
        trace.push({
          name: id,
          role: inferRole(id),
          activated: active.get(id)!,
          ruleVerdict: ruleVerdict.get(id)!,
          qualifier: qualMap.get(id)!,
          verdict: 0,                       // evaluate-trace path: verdict not assembled here (Go parity: zero-value)
          evidence: evidence.get(id),
          specs: specsMap.get(id),
          duration: elapsed,
          excluded: false,                  // patched per warrant by resolveSupersedes
        });
      }
      if (!active.get(id)) return -1.0;
      return combine(id, (aid) => calcTrace(aid, withDuration));
    }

    // Reporting-only pass: a target is excluded when any superseder has a_eff > 0
    // (except gates count, rebuttals do not; a non-active superseder keeps a_eff 0).
    // Every superseder is calc'd (idempotent) — no short-circuit. An excluded
    // superseder keeps its power to exclude.
    function resolveSupersedes(): Set<string> {
      const excluded = new Set<string>();
      for (const [target, superseders] of supersedeBy) {
        let applicable = false;
        for (const from of superseders) {
          calc(from);
          if ((applicability.get(from) ?? 0) > 0) applicable = true;
        }
        if (applicable) excluded.add(target);
      }
      return excluded;
    }

    function reset(): void {
      ran.clear();
      active.clear();
      ruleVerdict.clear();
      applicability.clear();
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
        applicability: applicability.get(r.name)!,   // active warrant → combine ran
      };
      if (opt.trace) {
        const excluded = resolveSupersedes();
        if (err) throw err;
        result.trace = trace.map(t => ({ ...t, name: shortName(t.name), excluded: excluded.has(t.name) }));
      }
      results.push(result);
    }

    if (full) {                                // run은 trace를 강제 off하므로 calc(non-trace)만 탐
      for (const r of this.rules) {
        if (!ran.has(r.name)) calc(r.name);    // 미도달 노드까지 평가
      }
      if (err) throw err;
    }

    const excluded = full ? resolveSupersedes() : new Set<string>();
    if (err) throw err;
    return { results, active, ruleVerdict, verdictCache, evidence, excluded };
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

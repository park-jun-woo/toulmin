from __future__ import annotations
import time
from typing import Any

from rulecat.types import (
    Context, EvalMethod, EvalOption, EvalResult, Strength, TraceEntry,
    RunResult,
)
from rulecat.trace import Trace
from rulecat.defeat_edge import DefeatEdge
from rulecat.rule_meta import RuleMeta
from rulecat.rule import Rule
from rulecat.rule_id import rule_id
from rulecat.short_name import short_name
from rulecat.to_verdict import to_verdict
from rulecat.detect_cycle import detect_cycle
from rulecat.detect_run_cycle import detect_run_cycle


run_max_depth = 64  # 깊이 백스톱


class Graph:
    def __init__(self, name: str) -> None:
        self.name = name
        self.rules: list[RuleMeta] = []
        self.roles: dict[str, str] = {}
        self.defeats: list[DefeatEdge] = []

    def rule(self, fn: Any) -> Rule:
        return self._register(fn, Strength.DEFEASIBLE, "rule")

    def counter(self, fn: Any) -> Rule:
        return self._register(fn, Strength.DEFEASIBLE, "counter")

    def except_(self, fn: Any) -> Rule:
        return self._register(fn, Strength.DEFEATER, "except")

    def _register(self, fn: Any, strength: Strength, role: str) -> Rule:
        rid = rule_id(fn, [])
        if rid in self.roles:
            raise ValueError(f"duplicate rule registration: {rid}")
        idx = len(self.rules)
        self.rules.append(RuleMeta(
            name=rid, qualifier=1.0, strength=int(strength),
            defeats=[], specs=[], fn=fn,
        ))
        self.roles[rid] = role
        return Rule(rid, self, idx, fn)

    def evaluate(self, ctx: Context, option: EvalOption | None = None) -> list[EvalResult]:
        if ctx is None:
            raise ValueError("ctx must not be None")
        results, *_ = self._evaluate(ctx, _resolve_option(option), full=False)
        return results

    def run(self, ctx: Context, option: EvalOption | None = None) -> RunResult:
        if ctx is None:
            raise ValueError("ctx must not be None")
        err = detect_run_cycle(self)  # 최상위 1회
        if err is not None:
            raise RuntimeError(f"toulmin: {err}")
        opt = _resolve_option(option)
        opt = EvalOption(method=opt.method, trace=False, duration=False)  # 강제 off
        return self._run_depth(ctx, opt, 0)

    def _run_depth(self, ctx: Context, opt: EvalOption, depth: int) -> RunResult:
        if depth > run_max_depth:
            raise RuntimeError(
                f"toulmin: run depth exceeded {run_max_depth} "
                "(possible runaway composition)"
            )
        results, active, verdict_cache, evidence, rule_verdict = self._evaluate(
            ctx, opt, full=True
        )

        # 디스패치 전 1회 — 전 노드 TraceEntry를 등록 순서로 조립
        # (인덱스 i ↔ self.rules[i] 구조적 보장)
        entries: list[TraceEntry] = []
        for r in self.rules:
            entries.append(TraceEntry(
                name=short_name(r.name),
                role=self.roles.get(r.name, "rule"),
                activated=active.get(r.name, False),
                rule_verdict=rule_verdict.get(r.name, -1.0),
                qualifier=r.qualifier,
                verdict=verdict_cache.get(r.name, 0.0),
                evidence=evidence.get(r.name),
                ground=ctx,
                specs=r.specs,
            ))

        tr = Trace(entries, ctx)

        for i, e in enumerate(entries):
            meta = self.rules[i]            # 인덱스 대응
            active_flag = e.activated and e.verdict > 0  # Active 하나만 발화
            if not active_flag:
                continue
            # (a) run_on 핸들러 — 먼저
            if meta.run_on is not None:
                try:
                    meta.run_on(e, tr)
                except Exception as exc:
                    raise RuntimeError(
                        f'run_on "{e.name}": {exc}'
                    ) from exc
            # (b) 실행 간선 — Active면 하위 그래프 Run (ctx 아래로, depth+1)
            if meta.run_graph is not None:
                try:
                    meta.run_graph._run_depth(ctx, opt, depth + 1)
                except Exception as exc:
                    raise RuntimeError(
                        f'run "{e.name}" → "{meta.run_graph.name}": {exc}'
                    ) from exc
        return RunResult(results=results, trace=tr)

    def _evaluate(
        self, ctx: Context, opt: EvalOption, full: bool
    ) -> tuple[
        list[EvalResult], dict[str, bool], dict[str, float], dict[str, Any],
        dict[str, float],
    ]:
        # Build eval maps
        fn_map: dict[str, Any] = {}
        qual_map: dict[str, float] = {}
        str_map: dict[str, int] = {}
        specs_map: dict[str, list[Any]] = {}
        edges: dict[str, list[str]] = {}
        forces: dict[str, list[float]] = {}  # forces[to][i] = edges[to][i]의 F

        for r in self.rules:
            fn_map[r.name] = r.fn
            qual_map[r.name] = r.qualifier
            str_map[r.name] = r.strength
            specs_map[r.name] = r.specs

        for d in self.defeats:
            edges.setdefault(d.to, []).append(d.from_)
            forces.setdefault(d.to, []).append(d.force)

        # Cycle detection
        cycle_err = detect_cycle(edges)
        if cycle_err:
            raise RuntimeError(cycle_err)

        # Attacker set
        attacker_set: set[str] = set()
        for attackers in edges.values():
            for a in attackers:
                attacker_set.add(a)

        # Eval state
        ran: set[str] = set()
        active: dict[str, bool] = {}          # 파생값: rule_verdict > -1
        rule_verdict: dict[str, float] = {}   # 규칙 함수 자신의 v ∈ [-1, 1]
        applicability: dict[str, float] = {}  # a_eff
        evidence: dict[str, Any] = {}
        verdict_cache: dict[str, float] = {}
        trace_entries: list[TraceEntry] = []
        err: list[Exception | None] = [None]  # mutable container for closure

        def is_warrant(name: str) -> bool:
            s = str_map.get(name, Strength.DEFEASIBLE)
            if s == Strength.DEFEATER:
                return False
            return name not in attacker_set

        def infer_role(id: str) -> str:
            if str_map.get(id) == Strength.DEFEATER:
                return "except"
            if id in attacker_set:
                return "counter"
            return "rule"

        def is_except_attacker(id: str) -> bool:
            # 예외(Except/Defeater)는 곱셈 게이트, 그 밖의 공격자는 반박(분모)
            return (
                str_map.get(id) == Strength.DEFEATER
                or self.roles.get(id) == "except"
            )

        def call_rule(id: str) -> bool:
            """규칙 함수를 실행해 결과를 기록한다. 오류면 err를 세우고 False."""
            try:
                raw_v, ev = fn_map[id](ctx, specs_map.get(id, []))
                v = to_verdict(raw_v)
            except Exception as e:
                err[0] = RuntimeError(f'rule "{id}": {e}')
                return False
            rule_verdict[id] = v
            active[id] = v > -1.0
            evidence[id] = ev
            return True

        def combine(id: str, eval_fn: Any) -> float:
            # a = (v+1)/2, S = Σ_반박 F·(v_j+1)/2,
            # a_eff = a · Π_예외 (1 − F·(v_j+1)/2), raw = q·a_eff/(1+S)
            a = (rule_verdict[id] + 1.0) / 2.0
            total = 0.0
            gate = 1.0
            if str_map.get(id) != Strength.STRICT:  # strict는 공격을 모두 무시
                fs = forces.get(id, [])
                for i, aid in enumerate(edges.get(id, [])):
                    f = fs[i] if i < len(fs) else 1.0
                    x = f * (eval_fn(aid) + 1.0) / 2.0
                    if is_except_attacker(aid):
                        gate *= 1.0 - x
                    else:
                        total += x
            a_eff = a * gate
            applicability[id] = a_eff
            q = qual_map.get(id, 1.0)
            v = 2.0 * (q * a_eff / (1.0 + total)) - 1.0
            verdict_cache[id] = v
            return v

        def calc(id: str) -> float:
            if err[0]:
                return -1.0
            cached = verdict_cache.get(id)
            if cached is not None:
                return cached
            if fn_map.get(id) is None:
                return -1.0
            if id not in ran:
                ran.add(id)
                if not call_rule(id):
                    return -1.0
            if not active.get(id):
                return -1.0  # 해당 없음: 공격자 평가를 건너뜀(지연 평가)
            return combine(id, calc)

        def calc_trace(id: str, with_duration: bool) -> float:
            if err[0]:
                return -1.0
            cached = verdict_cache.get(id)
            if cached is not None:
                return cached
            if fn_map.get(id) is None:
                return -1.0
            if id not in ran:
                ran.add(id)
                start = time.perf_counter() if with_duration else 0.0
                if not call_rule(id):
                    return -1.0
                elapsed = time.perf_counter() - start if with_duration else None
                trace_entries.append(TraceEntry(
                    name=id,
                    role=infer_role(id),
                    activated=active.get(id, False),
                    rule_verdict=rule_verdict[id],
                    qualifier=qual_map.get(id, 1.0),
                    evidence=evidence.get(id),
                    specs=specs_map.get(id, []),
                    duration=elapsed,
                ))
            if not active.get(id):
                return -1.0  # 해당 없음: 공격자 평가를 건너뜀(지연 평가)
            return combine(id, lambda aid: calc_trace(aid, with_duration))

        def reset() -> None:
            ran.clear()
            active.clear()
            rule_verdict.clear()
            applicability.clear()
            evidence.clear()
            verdict_cache.clear()
            trace_entries.clear()
            err[0] = None

        # Main evaluation loop
        results: list[EvalResult] = []
        for r in self.rules:
            if not is_warrant(r.name):
                continue
            if opt.trace:
                reset()
            verdict = (
                calc_trace(r.name, opt.duration) if opt.trace
                else calc(r.name)
            )
            if err[0]:
                raise err[0]
            if not active.get(r.name):
                continue
            result = EvalResult(
                name=short_name(r.name),
                verdict=verdict,
                evidence=evidence.get(r.name),
                applicability=applicability.get(r.name, 0.0),
            )
            if opt.trace:
                result.trace = [
                    TraceEntry(
                        name=short_name(t.name),
                        role=t.role,
                        activated=t.activated,
                        rule_verdict=t.rule_verdict,
                        qualifier=t.qualifier,
                        evidence=t.evidence,
                        specs=t.specs,
                        duration=t.duration,
                    )
                    for t in trace_entries
                ]
            results.append(result)

        if full:  # run이 trace를 강제 off하므로 calc(non-trace)만 탐
            for r in self.rules:
                if r.name not in ran:
                    calc(r.name)  # 미도달 노드까지 평가
            if err[0]:
                raise err[0]

        return results, active, verdict_cache, evidence, rule_verdict


def _resolve_option(opt: EvalOption | None) -> EvalOption:
    if opt is None:
        return EvalOption()
    if opt.duration:
        opt = EvalOption(method=opt.method, trace=True, duration=True)
    return opt

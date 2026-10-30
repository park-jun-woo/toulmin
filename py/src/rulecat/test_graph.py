import math
import unittest

import pytest

from rulecat import Graph, MapContext, TraceEntry
from rulecat.graph import run_max_depth
from rulecat.types import EvalOption


class RunTest(unittest.TestCase):
    def test_none_ctx_raises(self):
        g = Graph("t")
        g.rule(lambda ctx, specs: (True, None))
        with self.assertRaises(ValueError):
            g.run(None)

    def test_dispatch_calls_handler_and_skips_handlerless(self):
        g = Graph("t")
        seen = []
        a = g.rule(lambda ctx, specs: (True, None))
        a.run_on(lambda node, t: seen.append(node.verdict))
        # rule without any handler -> run_on is None -> skipped
        g.rule(lambda ctx, specs: (True, None))
        res = g.run(MapContext())
        self.assertEqual(len(seen), 1)
        self.assertGreater(seen[0], 0)
        self.assertEqual(len(res.trace.all()), 2)

    def test_handler_exception_wrapped_with_cause(self):
        g = Graph("t")
        cause = RuntimeError("boom")

        def bad(node, t):
            raise cause

        g.rule(lambda ctx, specs: (True, None)).run_on(bad)
        with self.assertRaises(RuntimeError) as cm:
            g.run(MapContext())
        self.assertIs(cm.exception.__cause__, cause)
        self.assertIn("boom", str(cm.exception))
        self.assertIn("run_on", str(cm.exception))

    def test_only_active_node_fires(self):
        # warrant attacked by an active counter -> warrant verdict 0.0 (defeated),
        # so only the counter (active, verdict>0) fires.
        g = Graph("t")
        fired = []
        w = g.rule(lambda ctx, specs: (True, None))
        c = g.counter(lambda ctx, specs: (True, None))
        c.attacks(w)
        w.run_on(lambda node, t: fired.append(("w", node.verdict)))
        c.run_on(lambda node, t: fired.append(("c", node.verdict)))
        g.run(MapContext())
        self.assertEqual([name for name, _ in fired], ["c"])

    def test_trace_exposes_verdict_for_other_nodes(self):
        # gradient/verdict query from a handler via the trace list.
        g = Graph("t")
        from rulecat.short_name import short_name
        branch = []
        w = g.rule(lambda ctx, specs: (True, None))
        c = g.counter(lambda ctx, specs: (True, None))
        c.attacks(w)
        w_name = short_name(w.id)

        def h(node, t):
            other = next((e for e in t.all() if e.name == w_name), None)
            branch.append("hard" if (other and other.verdict >= 0.5) else "soft")

        c.run_on(h)
        g.run(MapContext())
        # warrant attacked -> verdict 0.0 -> soft
        self.assertEqual(branch, ["soft"])


class RunCycleTest(unittest.TestCase):
    def test_run_raises_on_run_cycle(self):
        def a_rule(ctx, specs):
            return (True, None)

        g = Graph("cyc")
        g.rule(a_rule).run(g)  # self run-cycle
        with self.assertRaises(RuntimeError) as cm:
            g.run(MapContext())
        self.assertIn("run cycle", str(cm.exception))


class RunDepthTest(unittest.TestCase):
    def test_depth_guard_raises(self):
        g = Graph("deep")
        g.rule(lambda ctx, specs: (True, None))
        opt = EvalOption(method=0, trace=False, duration=False)
        with self.assertRaises(RuntimeError) as cm:
            g._run_depth(MapContext(), opt, run_max_depth + 1)
        self.assertIn("depth exceeded", str(cm.exception))

    def test_active_node_runs_sub_graph(self):
        def parent_rule(ctx, specs):
            return (True, None)

        def sub_rule(ctx, specs):
            return (True, None)

        sub = Graph("sub-run")
        sub.rule(sub_rule).run_on(lambda node, t: t.ctx().set("ran", True))

        parent = Graph("parent-run")
        parent.rule(parent_rule).run(sub)

        ctx = MapContext()
        parent.run(ctx)
        self.assertTrue(ctx.get("ran"))

    def test_inactive_node_with_run_graph_skips_sub(self):
        def parent_rule(ctx, specs):
            return (False, None)

        def sub_rule(ctx, specs):
            return (True, None)

        sub = Graph("sub-skip")
        sub.rule(sub_rule).run_on(lambda node, t: t.ctx().set("ran", True))

        parent = Graph("parent-skip")
        parent.rule(parent_rule).run(sub)

        ctx = MapContext()
        parent.run(ctx)
        self.assertIsNone(ctx.get("ran"))

    def test_sub_run_error_is_wrapped(self):
        cause = RuntimeError("boom")

        def parent_rule(ctx, specs):
            return (True, None)

        def sub_rule(ctx, specs):
            return (True, None)

        def bad(node, t):
            raise cause

        sub = Graph("sub-err")
        sub.rule(sub_rule).run_on(bad)

        parent = Graph("parent-err")
        parent.rule(parent_rule).run(sub)

        with self.assertRaises(RuntimeError) as cm:
            parent.run(MapContext())
        msg = str(cm.exception)
        self.assertTrue(msg.startswith('run "'))
        self.assertIn("sub-err", msg)


class EvaluateInternalTest(unittest.TestCase):
    def test_none_ctx_raises(self):
        g = Graph("t")
        g.rule(lambda ctx, specs: (True, None))
        with self.assertRaises(ValueError):
            g.evaluate(None)

    def test_duplicate_target_and_attackers(self):
        g = Graph("t")
        w = g.rule(lambda ctx, specs: (True, None))
        c1 = g.counter(lambda ctx, specs: (True, None))
        c2 = g.counter(lambda ctx, specs: (True, None))
        c1.attacks(w)
        c2.attacks(w)
        results = g.evaluate(MapContext())
        self.assertTrue(results)

    def test_cycle_raises(self):
        g = Graph("t")
        a = g.rule(lambda ctx, specs: (True, None))
        b = g.counter(lambda ctx, specs: (True, None))
        a.attacks(b)
        b.attacks(a)
        with self.assertRaises(RuntimeError):
            g.evaluate(MapContext())


# --- 규칙 함수 verdict 반환·예외 게이트·force 의미론(Go 판 c084dc2·fad0ea9과 같은 값) ---


class Tag:
    def __init__(self, value: str):
        self.value = value

    def spec_name(self) -> str:
        return "Tag"

    def validate(self) -> None:
        pass


def _const(v, ev=None):
    return lambda ctx, specs: (v, ev)


# --- boolean 호환 ---


def _compat_graph(use_bool: bool) -> Graph:
    # w ← c1 (q 0.7), c2 (q 0.4, 비활성), c3 (q 0.3); c1 ← c4 (q 0.9)
    def of(b):
        if use_bool:
            return lambda ctx, specs: (b, None)
        v = 1.0 if b else -1.0
        return lambda ctx, specs: (v, None)

    g = Graph("compat")
    w = g.rule(of(True)).with_spec(Tag("w")).qualifier(0.8)
    c1 = g.counter(of(True)).with_spec(Tag("c1")).qualifier(0.7)
    c2 = g.counter(of(False)).with_spec(Tag("c2")).qualifier(0.4)
    c3 = g.counter(of(True)).with_spec(Tag("c3")).qualifier(0.3)
    c4 = g.counter(of(True)).with_spec(Tag("c4")).qualifier(0.9)
    c1.attacks(w)
    c2.attacks(w)
    c3.attacks(w)
    c4.attacks(c1)
    return g


@pytest.mark.parametrize("trace", [False, True])
def test_bool_compat_equals_legacy_formula(trace):
    # 옛 공식: raw = q/(1+Σ(v_j+1)/2), v = 2raw−1 — 비트 단위로 같아야 한다
    vc4 = 2 * 0.9 - 1
    vc1 = 2 * (0.7 / (1.0 + (vc4 + 1.0) / 2.0)) - 1
    vc3 = 2 * 0.3 - 1
    legacy = 2 * (0.8 / (1.0 + ((vc1 + 1.0) / 2.0 + 0.0 + (vc3 + 1.0) / 2.0))) - 1
    opt = EvalOption(trace=trace)
    rb = _compat_graph(True).evaluate(MapContext(), opt)
    rf = _compat_graph(False).evaluate(MapContext(), opt)
    assert rb[0].verdict == legacy
    assert rf[0].verdict == legacy
    assert rb[0].applicability == 1.0


def test_bool_false_is_not_applicable():
    g = Graph("t")
    g.rule(_const(False))
    assert g.evaluate(MapContext()) == []


def test_int_verdict_is_numeric_not_bool():
    # int 1은 +1, int 0은 "반쯤 해당"(bool False와 다르다)
    g = Graph("t")
    g.rule(_const(1)).with_spec(Tag("one"))
    g.rule(_const(0)).with_spec(Tag("zero"))
    res = g.evaluate(MapContext())
    assert [r.verdict for r in res] == [1.0, 0.0]
    assert [r.applicability for r in res] == [1.0, 0.5]


# --- 부분 verdict ---


def test_zero_means_half_applies():
    g = Graph("t")
    g.rule(_const(0.0, "half"))
    res = g.evaluate(MapContext())
    assert len(res) == 1
    assert res[0].verdict == 0.0
    assert res[0].applicability == 0.5
    assert res[0].evidence == "half"


def test_qualifier_scales_partial_verdict():
    # v = 0.5 → a = 0.75; q = 0.8 → raw = 0.6 → verdict 0.2
    g = Graph("t")
    g.rule(_const(0.5)).qualifier(0.8)
    res = g.evaluate(MapContext())
    assert math.isclose(res[0].verdict, 0.2, abs_tol=1e-12)
    assert res[0].applicability == 0.75


def test_minus_one_is_inactive_and_lazy():
    attacked = []
    g = Graph("t")
    w = g.rule(_const(-1.0))

    def counter(ctx, specs):
        attacked.append(True)
        return 1.0, None

    g.counter(counter).attacks(w)
    assert g.evaluate(MapContext()) == []
    assert attacked == []


def test_barely_applies_evaluates_attackers():
    attacked = []
    g = Graph("t")
    w = g.rule(_const(-0.5))

    def counter(ctx, specs):
        attacked.append(True)
        return -1.0, None

    g.counter(counter).attacks(w)
    res = g.evaluate(MapContext())
    assert len(res) == 1
    assert res[0].verdict == -0.5
    assert attacked == [True]


# --- verdict 범위·NaN·타입 오류 ---


@pytest.mark.parametrize("bad,shown", [
    (1.5, "1.5"), (-2.0, "-2"), (float("nan"), "NaN"), (float("inf"), "+Inf"),
])
@pytest.mark.parametrize("mode", ["plain", "trace", "run"])
def test_verdict_out_of_range_is_error(bad, shown, mode):
    g = Graph("range")
    w = g.rule(_const(True))

    def bad_counter(ctx, specs):
        return bad, None

    c = g.counter(bad_counter)
    c.attacks(w)
    with pytest.raises(RuntimeError) as ei:
        if mode == "plain":
            g.evaluate(MapContext())
        elif mode == "trace":
            g.evaluate(MapContext(), EvalOption(trace=True))
        else:
            g.run(MapContext())
    assert str(ei.value) == (
        f'rule "{c.id}": rule verdict {shown} out of range [-1, 1]'
    )


@pytest.mark.parametrize("bad", [None, "yes", (1.0,)])
def test_verdict_non_number_is_error(bad):
    g = Graph("t")
    g.rule(_const(bad))
    with pytest.raises(RuntimeError, match="rule verdict must be bool or float"):
        g.evaluate(MapContext())


def test_verdict_bounds_are_inclusive():
    g = Graph("t")
    g.rule(_const(1.0)).with_spec(Tag("hi"))
    g.rule(_const(-1.0)).with_spec(Tag("lo"))
    res = g.evaluate(MapContext())
    assert [r.verdict for r in res] == [1.0]


# --- 예외 게이트와 applicability ---


def _gate_graph(v: float, as_except: bool, n: int) -> Graph:
    g = Graph("gate")
    w = g.rule(_const(True))
    fn = _const(v)
    for i in range(n):
        add = g.except_ if as_except else g.counter
        add(fn).with_spec(Tag(chr(ord("a") + i))).attacks(w)
    return g


@pytest.mark.parametrize("v,as_except,verdict,applic", [
    (1.0, False, -0.5, 1.0),                    # S = 3 → raw 1/4 (변화 없음)
    (1.0, True, -1.0, 0.0),                     # gate = 0³
    (0.5, True, 2 * 0.25 ** 3 - 1, 0.25 ** 3),  # gate = (1−0.75)³
    (-1.0, True, 1.0, 1.0),                     # 비활성 예외는 게이트 1
])
@pytest.mark.parametrize("trace", [False, True])
def test_except_gate(v, as_except, verdict, applic, trace):
    res = _gate_graph(v, as_except, 3).evaluate(MapContext(), EvalOption(trace=trace))
    assert len(res) == 1
    assert math.isclose(res[0].verdict, verdict, abs_tol=1e-12)
    assert math.isclose(res[0].applicability, applic, abs_tol=1e-12)


def test_except_gate_and_counter_combine():
    # w(v=1) ← counter(1.0), except(0.0): a_eff = 1·(1−0.5) = 0.5, S = 1
    # raw = 0.5/2 = 0.25 → verdict −0.5
    g = Graph("t")
    w = g.rule(_const(True))
    g.counter(_const(1.0)).attacks(w)
    g.except_(_const(0.0)).attacks(w)
    res = g.evaluate(MapContext())
    assert res[0].applicability == 0.5
    assert res[0].verdict == -0.5


def test_strict_ignores_counter_and_except():
    from rulecat.types import Strength
    g = Graph("t")
    w = g.rule(_const(0.5))
    g.rules[0].strength = int(Strength.STRICT)
    g.counter(_const(True)).attacks(w)
    g.except_(_const(True)).attacks(w)
    res = g.evaluate(MapContext())
    assert res[0].applicability == 0.75
    assert res[0].verdict == 0.5


# --- force ---


@pytest.mark.parametrize("as_except,force,verdict", [
    (False, 0.0, 1.0),
    (False, 0.5, 1.0 / 3.0),   # S = 0.5 → raw 2/3
    (False, 1.0, 0.0),
    (True, 0.0, 1.0),
    (True, 0.5, 0.0),          # a_eff = 0.5
    (True, 0.25, 0.5),
])
@pytest.mark.parametrize("trace", [False, True])
def test_attack_force(as_except, force, verdict, trace):
    g = Graph("force")
    w = g.rule(_const(True))
    a = (g.except_ if as_except else g.counter)(_const(True))
    a.attacks(w, force)
    res = g.evaluate(MapContext(), EvalOption(trace=trace))
    assert math.isclose(res[0].verdict, verdict, abs_tol=1e-12)


# --- Qualifier NaN ---


def test_qualifier_rejects_nan():
    g = Graph("t")
    r = g.rule(_const(True))
    with pytest.raises(ValueError, match="qualifier must be between 0.0 and 1.0"):
        r.qualifier(float("nan"))
    assert g.rules[0].qualifier == 1.0


# --- trace의 rule_verdict ---


def test_rule_verdict_in_run_and_evaluate_trace():
    fired = {}
    g = Graph("rv")
    w = g.rule(_const(0.8))
    w.run_on(lambda node, t: fired.setdefault("w", node.rule_verdict))
    g.counter(_const(-1.0)).attacks(w)

    rr = g.run(MapContext())
    assert fired == {"w": 0.8}
    want = {"rule": 0.8, "counter": -1.0}
    for e in rr.trace.all():
        assert e.rule_verdict == want[e.role]
        assert e.activated == (e.rule_verdict > -1)

    res = g.evaluate(MapContext(), EvalOption(trace=True))
    assert {e.role: e.rule_verdict for e in res[0].trace} == want


def test_trace_entry_rule_verdict_default():
    # 실행되지 않은 노드의 rule_verdict는 -1(activated=False와 일치)
    e = TraceEntry()
    assert e.rule_verdict == -1.0 and e.activated is False



# --- 평가 경로의 가지(캐시·공유 비활성 공격자·외부 공격자·오류 전파) ---


@pytest.mark.parametrize("trace", [False, True])
def test_shared_attacker_is_cached(trace):
    # w ← r1, r2; d가 r1·r2를 모두 공격 → 두 번째 방문은 verdict 캐시를 쓴다
    calls = []

    def d_fn(ctx, specs):
        calls.append(1)
        return 0.0, None

    g = Graph("t")
    w = g.rule(_const(True))
    r1 = g.counter(_const(True)).with_spec(Tag("r1"))
    r2 = g.counter(_const(True)).with_spec(Tag("r2"))
    d = g.counter(d_fn)
    r1.attacks(w)
    r2.attacks(w)
    d.attacks(r1)
    d.attacks(r2)
    res = g.evaluate(MapContext(), EvalOption(trace=trace))
    assert calls == [1]
    # d verdict 0 → r1·r2 raw = 1/1.5 → verdict 1/3 → w raw = 1/(1+2·(2/3)) = 3/7
    assert math.isclose(res[0].verdict, 2 * 3 / 7 - 1)
    if trace:
        assert len(res[0].trace) == 4


@pytest.mark.parametrize("trace", [False, True])
def test_shared_inactive_attacker_runs_once(trace):
    # 비활성 노드는 캐시되지 않지만 ran에 남아 다시 실행되지 않는다
    calls = []

    def d_fn(ctx, specs):
        calls.append(1)
        return -1.0, None

    g = Graph("t")
    w = g.rule(_const(True))
    r1 = g.counter(_const(True)).with_spec(Tag("r1"))
    r2 = g.counter(_const(True)).with_spec(Tag("r2"))
    d = g.except_(d_fn)
    r1.attacks(w)
    r2.attacks(w)
    d.attacks(r1)
    d.attacks(r2)
    res = g.evaluate(MapContext(), EvalOption(trace=trace))
    assert calls == [1]
    assert math.isclose(res[0].verdict, 2 * (1 / 3) - 1)


@pytest.mark.parametrize("trace", [False, True])
def test_foreign_attacker_is_not_applicable(trace):
    # 다른 그래프의 규칙이 건 간선: 이 그래프에 함수가 없으므로 -1(무시)
    g = Graph("t")
    other = Graph("other")
    w = g.rule(_const(True))
    other.counter(_const(True)).attacks(w)
    g.defeats.extend(other.defeats)
    res = g.evaluate(MapContext(), EvalOption(trace=trace))
    assert res[0].verdict == 1.0
    assert res[0].applicability == 1.0


@pytest.mark.parametrize("trace", [False, True])
def test_error_stops_remaining_attackers(trace):
    seen = []

    def boom(ctx, specs):
        raise RuntimeError("boom")

    def later(ctx, specs):
        seen.append(1)
        return True, None

    g = Graph("t")
    w = g.rule(_const(True))
    g.counter(boom).attacks(w)
    g.counter(later).attacks(w)
    with pytest.raises(RuntimeError, match='^rule ".+": boom$'):
        g.evaluate(MapContext(), EvalOption(trace=trace))
    assert seen == []


def test_run_full_pass_surfaces_error_in_unreached_node():
    # evaluate는 비활성 warrant 뒤의 공격자에 닿지 않지만 run의 전체 패스는 닿는다
    g = Graph("t")
    w = g.rule(_const(False))
    g.counter(_const(2.0)).attacks(w)
    assert g.evaluate(MapContext()) == []
    with pytest.raises(RuntimeError, match=r"rule verdict 2 out of range \[-1, 1\]"):
        g.run(MapContext())


if __name__ == "__main__":
    unittest.main()

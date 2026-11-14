"""배제(Supersedes): 배제자의 적용도 a_eff > 0일 때만 대상을 크리스프하게 배제한다.

네 포트 공통 VIP 시나리오: vip_discount가 general_discount를 배제한다. 예외 게이트는
반영하고 반박(counter)은 반영하지 않는다.
"""
import pytest
from rulecat import EvalOption, Graph, MapContext, Strength
from rulecat.short_name import short_name


def general_discount(ctx, specs):
    return (True, None)


def vip_discount(ctx, specs):
    return (True, None)


def suspended_vip(ctx, specs):
    return (True, None)


def vip_disputed(ctx, specs):
    return (True, None)


def _x(ctx, specs):
    return (True, None)


def _y(ctx, specs):
    return (True, None)


def _z(ctx, specs):
    return (True, None)


def _vip_graph(wire):
    g = Graph("vip")
    general = g.rule(general_discount)
    vip = g.rule(vip_discount)
    vip.supersedes(general)
    wire(g, vip)
    return g, short_name(general.id), short_name(vip.id)


@pytest.mark.parametrize("wire, vip_verdict, vip_a_eff, want_excluded", [
    pytest.param(lambda g, vip: None, 1.0, 1.0, True, id="no exception"),
    pytest.param(lambda g, vip: g.except_(suspended_vip).attacks(vip), -1.0, 0.0, False,
                 id="exception fully gates"),
    pytest.param(lambda g, vip: g.except_(suspended_vip).attacks(vip, force=0.5), 0.0, 0.5, True,
                 id="exception force 0.5"),
    pytest.param(lambda g, vip: g.counter(vip_disputed).attacks(vip), 0.0, 1.0, True,
                 id="counter keeps priority"),
])
def test_vip_supersede(wire, vip_verdict, vip_a_eff, want_excluded):
    g, gn, vn = _vip_graph(wire)
    rr = g.run(MapContext())
    ge = rr.trace.get(gn)
    ve = rr.trace.get(vn)
    assert ge.excluded is want_excluded
    assert ge.verdict == 1.0  # 보고 전용
    assert ve.excluded is False
    assert ve.activated is True
    assert ve.verdict == vip_verdict
    vr = next(r for r in rr.results if r.name == vn)
    assert vr.applicability == vip_a_eff

    traced = g.evaluate(MapContext(), EvalOption(trace=True))
    gw = next(r for r in traced if r.name == gn)
    assert next(t for t in gw.trace if t.name == gn).excluded is want_excluded


def test_strict_superseder_ignores_exception_gate():
    g, gn, _ = _vip_graph(lambda g, vip: g.except_(suspended_vip).attacks(vip))
    g.rules[1].strength = int(Strength.STRICT)
    assert g.run(MapContext()).trace.get(gn).excluded is True


def _chain(gate_x):
    g = Graph("chain")
    z = g.rule(_z)
    y = g.rule(_y)
    x = g.rule(_x)
    x.supersedes(y)
    y.supersedes(z)
    if gate_x:
        g.except_(suspended_vip).attacks(x)
    tr = g.run(MapContext()).trace
    return [tr.get(short_name(r.id)).excluded for r in (x, y, z)]


def test_excluded_superseder_still_excludes():
    assert _chain(False) == [False, True, True]


def test_gated_superseder_loses_only_its_own_exclusion():
    assert _chain(True) == [False, False, True]


def test_supersede_cycle_rejected():
    g = Graph("cycle")
    a = g.rule(_x)
    b = g.rule(_y)
    a.supersedes(b)
    b.supersedes(a)
    with pytest.raises(RuntimeError, match="supersede cycle detected"):
        g.evaluate(MapContext())


class _Spec:
    def spec_name(self):
        return "s"

    def validate(self):
        pass


def test_with_spec_rewrites_supersede_edges():
    g = Graph("with")
    general = g.rule(general_discount)
    vip = g.rule(vip_discount)
    vip.supersedes(general)
    vip.with_spec(_Spec())
    general.with_spec(_Spec())
    assert [(e.from_, e.to) for e in g.supersedes] == [(vip.id, general.id)]
    assert g.run(MapContext()).trace.get(short_name(general.id)).excluded is True

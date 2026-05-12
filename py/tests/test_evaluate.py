import pytest
from rulecat import Graph, MapContext


def test_single_warrant():
    g = Graph("test")
    g.rule(lambda ctx, specs: (True, None))
    results = g.evaluate(MapContext())
    assert len(results) == 1
    assert results[0].verdict == 1.0


def test_warrant_counter():
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (True, None))
    r = g.counter(lambda ctx, specs: (True, None))
    r.attacks(w)
    results = g.evaluate(MapContext())
    assert results[0].verdict == 0.0


def test_compensation_chain():
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (True, None))
    r = g.counter(lambda ctx, specs: (True, None))
    d = g.except_(lambda ctx, specs: (True, None))
    r.attacks(w)
    d.attacks(r)
    results = g.evaluate(MapContext())
    assert abs(results[0].verdict - 1/3) < 1e-5


def test_diamond_dag():
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (True, None))
    r1 = g.counter(lambda ctx, specs: (True, None))
    r2 = g.counter(lambda ctx, specs: (True, None))
    d = g.except_(lambda ctx, specs: (True, None))
    r1.attacks(w)
    r2.attacks(w)
    d.attacks(r1)
    d.attacks(r2)
    results = g.evaluate(MapContext())
    assert results[0].verdict == 0.0


def test_inactive_rule():
    g = Graph("test")
    g.rule(lambda ctx, specs: (False, None))
    results = g.evaluate(MapContext())
    assert len(results) == 0


def test_cycle_detection():
    g = Graph("test")
    a = g.rule(lambda ctx, specs: (True, None))
    b = g.counter(lambda ctx, specs: (True, None))
    a.attacks(b)
    b.attacks(a)
    with pytest.raises(RuntimeError, match="cycle"):
        g.evaluate(MapContext())


def test_none_ctx():
    g = Graph("test")
    g.rule(lambda ctx, specs: (True, None))
    with pytest.raises(ValueError, match="ctx"):
        g.evaluate(None)


def test_rule_exception():
    def bad_fn(ctx, specs):
        raise RuntimeError("boom")
    g = Graph("test")
    g.rule(bad_fn)
    with pytest.raises(RuntimeError, match="boom"):
        g.evaluate(MapContext())


def test_qualifier_zero():
    g = Graph("test")
    g.rule(lambda ctx, specs: (True, None)).qualifier(0.0)
    results = g.evaluate(MapContext())
    assert results[0].verdict == -1.0


def test_evaluate_idempotent():
    g = Graph("test")
    g.rule(lambda ctx, specs: (True, None))
    r1 = g.evaluate(MapContext())
    r2 = g.evaluate(MapContext())
    assert r1[0].verdict == r2[0].verdict

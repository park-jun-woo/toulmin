from rulecat import Graph, MapContext, EvalOption


def test_trace_entries():
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (True, "warrant-ev"))
    r = g.counter(lambda ctx, specs: (True, "counter-ev"))
    r.attacks(w)
    results = g.evaluate(MapContext(), EvalOption(trace=True))
    assert len(results[0].trace) == 2
    assert results[0].trace[0].role == "rule"
    assert results[0].trace[1].role == "counter"
    assert results[0].trace[0].activated is True


def test_duration():
    g = Graph("test")
    g.rule(lambda ctx, specs: (True, None))
    results = g.evaluate(MapContext(), EvalOption(duration=True))
    assert results[0].trace is not None
    assert results[0].trace[0].duration is not None
    assert results[0].trace[0].duration >= 0


def test_inactive_in_trace():
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (True, None))
    r = g.counter(lambda ctx, specs: (False, None))
    r.attacks(w)
    results = g.evaluate(MapContext(), EvalOption(trace=True))
    assert len(results[0].trace) == 2
    assert results[0].trace[1].activated is False

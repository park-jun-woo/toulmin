import pytest
from rulecat import Graph, MapContext


def _auth(ctx, specs):
    return (True, "auth")


def _block(ctx, specs):
    return (bool(ctx.get("blocked")), "block")


def _exempt(ctx, specs):
    return (bool(ctx.get("internal")), "exempt")


def _build_access_control():
    """authenticate (rule) <- blockIP (counter) <- exemptInternalIP (except)."""
    g = Graph("access-control")
    auth = g.rule(_auth)
    block = g.counter(_block)
    exempt = g.except_(_exempt)
    block.attacks(auth)
    exempt.attacks(block)
    return g, auth, block, exempt


def short_of(rule):
    from rulecat.short_name import short_name
    return short_name(rule.id)


def _entry(trace, name):
    return next((e for e in trace if e.name == name), None)


def _is_active(trace, name):
    e = _entry(trace, name)
    return e is not None and e.activated and e.verdict > 0


def test_active_only_external_blocked_ip():
    g, auth, block, exempt = _build_access_control()
    fired = []

    def on(name):
        def h(node, t):
            fired.append(name)
        return h

    for r, n in ((auth, "auth"), (block, "block"), (exempt, "exempt")):
        r.run_on(on(n))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    result = g.run(ctx)

    # block active (verdict>0); auth defeated (verdict 0.0); exempt inactive
    assert fired == ["block"]
    assert _is_active(result.trace.all(), short_of(block))
    assert not _is_active(result.trace.all(), short_of(auth))
    assert not _is_active(result.trace.all(), short_of(exempt))


def test_active_only_internal_network():
    g, auth, block, exempt = _build_access_control()
    fired = []

    def on(name):
        def h(node, t):
            fired.append(name)
        return h

    for r, n in ((auth, "auth"), (block, "block"), (exempt, "exempt")):
        r.run_on(on(n))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", True)
    result = g.run(ctx)

    # exempt active -> block defeated -> auth active again
    assert _is_active(result.trace.all(), short_of(exempt))
    assert _is_active(result.trace.all(), short_of(auth))
    assert not _is_active(result.trace.all(), short_of(block))
    assert set(fired) == {"exempt", "auth"}


def test_full_pass_includes_unreached_node():
    # Lazy evaluate: warrant W is inactive, so its attacker C is never reached
    # during evaluate, but run's full pass still produces a trace entry for it.
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (False, None))
    c = g.counter(lambda ctx, specs: (False, None))
    c.attacks(w)

    fired = []
    c.run_on(lambda node, t: fired.append("c"))

    # evaluate does not reach C and never fires handlers.
    g.evaluate(MapContext())
    assert fired == []

    result = g.run(MapContext())
    # inactive node fires nothing, but appears in the trace.
    assert fired == []
    names = [e.name for e in result.trace.all()]
    assert short_of(c) in names


def test_evaluate_pure_and_idempotent():
    g, auth, block, exempt = _build_access_control()
    fired = []
    for r in (auth, block, exempt):
        r.run_on(lambda node, t: fired.append(t))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    r1 = g.evaluate(ctx)
    r2 = g.evaluate(ctx)

    assert fired == []  # evaluate never fires handlers
    assert len(r1) == len(r2)
    assert [x.verdict for x in r1] == [x.verdict for x in r2]


def test_handler_raise_stops_and_propagates_with_context():
    g = Graph("test")
    w = g.rule(_auth)

    cause = RuntimeError("boom")

    def bad(node, t):
        raise cause

    w.run_on(bad)

    with pytest.raises(RuntimeError) as exc_info:
        g.run(MapContext())

    msg = str(exc_info.value)
    assert "_auth" in msg or "auth" in msg
    assert "run_on" in msg
    assert "boom" in msg
    assert exc_info.value.__cause__ is cause


def test_handlers_fire_in_registration_order():
    g = Graph("test")
    order = []

    def make(name):
        def h(node, t):
            order.append(name)
        return h

    a = g.rule(lambda ctx, specs: (True, None))
    b = g.rule(lambda ctx, specs: (True, None))
    c = g.rule(lambda ctx, specs: (True, None))
    a.run_on(make("a"))
    b.run_on(make("b"))
    c.run_on(make("c"))

    result = g.run(MapContext())
    assert order == ["a", "b", "c"]
    assert len(result.trace.all()) == 3


def test_verdict_zero_node_does_not_fire():
    # warrant attacked by one active counter -> verdict 0.0 -> not Active
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (True, None))
    c = g.counter(lambda ctx, specs: (True, None))
    c.attacks(w)

    seen = []
    w.run_on(lambda node, t: seen.append(node.verdict))

    result = g.run(MapContext())
    assert seen == []  # defeated (verdict 0.0) -> no fire
    w_entry = _entry(result.trace.all(), short_of(w))
    assert w_entry.verdict == 0.0
    assert w_entry.activated is True


# --- trace-based node queries from handlers ---


def test_trace_sees_all_three_nodes():
    g, auth, block, exempt = _build_access_control()
    seen = []

    def h(node, t):
        seen.append(sorted(e.name for e in t.all()))

    # block is the active node here.
    block.run_on(h)

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    expected = sorted([short_of(auth), short_of(block), short_of(exempt)])
    assert seen and seen[0] == expected


def test_trace_self_is_firing_node():
    g, auth, block, exempt = _build_access_control()
    captured = []

    block.run_on(lambda node, t: captured.append(node.name))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    assert captured == [short_of(block)]


def test_trace_ground_is_ctx():
    g, auth, block, exempt = _build_access_control()
    grounds = []

    block.run_on(lambda node, t: grounds.append(node.ground))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    assert grounds and all(gr is ctx for gr in grounds)


def test_returned_trace_is_queryable_posthoc():
    g, auth, block, exempt = _build_access_control()

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    result = g.run(ctx)

    trace = result.trace
    assert len(trace.all()) == 3
    assert trace.get(short_of(auth)).verdict == 0.0   # defeated
    assert trace.get(short_of(block)).verdict > 0     # active
    assert trace.get("nope") is None


def test_gradient_branch_via_verdict():
    # blocked=False: block inactive -> auth unattacked -> verdict 1.0 -> hard
    g, auth, block, exempt = _build_access_control()
    branch = []
    auth_name = short_of(auth)

    def h(node, t):
        other = _entry(t.all(), auth_name)
        branch.append("hard" if (other and other.verdict >= 0.5) else "soft")

    auth.run_on(h)

    ctx = MapContext()
    ctx.set("blocked", False)
    ctx.set("internal", False)
    g.run(ctx)
    assert branch == ["hard"]


def test_gradient_branch_soft_when_defeated():
    # internal=True: exempt active -> block defeated (0.0) -> auth verdict ~0.33 -> soft
    g, auth, block, exempt = _build_access_control()
    branch = []
    auth_name = short_of(auth)

    def h(node, t):
        other = _entry(t.all(), auth_name)
        branch.append("hard" if (other and other.verdict >= 0.5) else "soft")

    exempt.run_on(h)

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", True)
    g.run(ctx)

    assert branch == ["soft"]

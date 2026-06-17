import pytest
from rulecat import Graph, MapContext, NodeEventType


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


def test_three_events_external_blocked_ip():
    g, auth, block, exempt = _build_access_control()
    events: dict[str, NodeEventType] = {}

    def on(name):
        def h(ctx, ev, view):
            events[name] = ev.type
        return h

    for r, n in ((auth, "auth"), (block, "block"), (exempt, "exempt")):
        r.on_active(on(n)).on_defeated(on(n)).on_inactive(on(n))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    assert events["block"] == NodeEventType.ACTIVE
    assert events["auth"] == NodeEventType.DEFEATED
    assert events["exempt"] == NodeEventType.INACTIVE


def test_three_events_internal_network():
    g, auth, block, exempt = _build_access_control()
    events: dict[str, NodeEventType] = {}

    def on(name):
        def h(ctx, ev, view):
            events[name] = ev.type
        return h

    for r, n in ((auth, "auth"), (block, "block"), (exempt, "exempt")):
        r.on_active(on(n)).on_defeated(on(n)).on_inactive(on(n))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", True)
    g.run(ctx)

    assert events["exempt"] == NodeEventType.ACTIVE
    assert events["block"] == NodeEventType.DEFEATED
    assert events["auth"] == NodeEventType.ACTIVE


def test_full_pass_fires_inactive_for_unreached_node():
    # Lazy evaluate: warrant W is inactive, so its attacker C is never reached.
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (False, None))
    c = g.counter(lambda ctx, specs: (False, None))
    c.attacks(w)

    fired = []
    c.on_inactive(lambda ctx, ev, view: fired.append(ev.type))

    # evaluate does not reach C and never fires handlers.
    g.evaluate(MapContext())
    assert fired == []

    result = g.run(MapContext())
    assert fired == [NodeEventType.INACTIVE]
    # Unreached node still produced an event in run's view.
    names = [e.name for e in result.view.all()]
    assert short_of(c) in names


def test_evaluate_pure_and_idempotent():
    g, auth, block, exempt = _build_access_control()
    fired = []
    for r in (auth, block, exempt):
        r.on_active(lambda ctx, ev, view: fired.append(ev))
        r.on_defeated(lambda ctx, ev, view: fired.append(ev))
        r.on_inactive(lambda ctx, ev, view: fired.append(ev))

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

    def bad(ctx, ev, view):
        raise cause

    w.on_active(bad)

    with pytest.raises(RuntimeError) as exc_info:
        g.run(MapContext())

    msg = str(exc_info.value)
    assert "_auth" in msg or "auth" in msg
    assert "ACTIVE" in msg
    assert "boom" in msg
    assert exc_info.value.__cause__ is cause


def test_events_order_matches_registration_order():
    g = Graph("test")
    order = []

    def make(name):
        def h(ctx, ev, view):
            order.append(name)
        return h

    a = g.rule(lambda ctx, specs: (True, None))
    b = g.rule(lambda ctx, specs: (True, None))
    c = g.rule(lambda ctx, specs: (True, None))
    a.on_active(make("a"))
    b.on_active(make("b"))
    c.on_active(make("c"))

    result = g.run(MapContext())
    assert order == ["a", "b", "c"]
    assert len(result.view.all()) == 3


def test_verdict_zero_is_defeated():
    # warrant attacked by one active counter -> verdict 0.0 -> DEFEATED
    g = Graph("test")
    w = g.rule(lambda ctx, specs: (True, None))
    c = g.counter(lambda ctx, specs: (True, None))
    c.attacks(w)

    seen = []
    w.on_active(lambda ctx, ev, view: seen.append(("active", ev.verdict)))
    w.on_defeated(lambda ctx, ev, view: seen.append(("defeated", ev.verdict)))

    g.run(MapContext())
    assert seen == [("defeated", 0.0)]


# --- RunView tests ---


def test_view_sees_all_three_nodes():
    g, auth, block, exempt = _build_access_control()
    seen = []

    def h(ctx, ev, view):
        seen.append(sorted(e.name for e in view.all()))

    auth.on_active(h).on_defeated(h).on_inactive(h)

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    expected = sorted([short_of(auth), short_of(block), short_of(exempt)])
    assert seen and seen[0] == expected


def test_view_is_immutable_across_handlers():
    g, auth, block, exempt = _build_access_control()
    captured = {}

    def h(name):
        def handler(ctx, ev, view):
            # mutate the returned copy; must not affect later handlers
            view.all().append(ev)
            captured[name] = len(view.all())
        return handler

    auth.on_active(h("auth")).on_defeated(h("auth")).on_inactive(h("auth"))
    block.on_active(h("block")).on_defeated(h("block")).on_inactive(h("block"))
    exempt.on_active(h("exempt")).on_defeated(h("exempt")).on_inactive(h("exempt"))

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    # every handler sees the same 3-node snapshot regardless of prior mutation
    assert set(captured.values()) == {3}


def test_view_attackers_returns_block_ip():
    g, auth, block, exempt = _build_access_control()
    found = {}

    def h(ctx, ev, view):
        found["auth"] = [a.name for a in view.attackers(short_of(auth))]
        found["block"] = [a.name for a in view.attackers(short_of(block))]

    auth.on_active(h).on_defeated(h).on_inactive(h)

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    assert found["auth"] == [short_of(block)]
    assert found["block"] == [short_of(exempt)]


def test_view_get_missing_returns_none():
    g, auth, block, exempt = _build_access_control()
    seen = []

    def h(ctx, ev, view):
        seen.append(view.get("nope"))

    auth.on_active(h).on_defeated(h).on_inactive(h)

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    g.run(ctx)

    assert seen and all(x is None for x in seen)


def _gradient_branch(internal):
    g, auth, block, exempt = _build_access_control()
    branch = []

    def h(ctx, ev, view):
        node = view.get(short_of(auth))
        if node is not None and node.verdict >= 0.5:
            branch.append("hard")
        else:
            branch.append("soft")

    block.on_active(h).on_defeated(h).on_inactive(h)

    ctx = MapContext()
    ctx.set("blocked", False)
    ctx.set("internal", internal)
    g.run(ctx)
    return branch


def test_view_gradient_branch_via_verdict():
    # blocked=False: block inactive -> auth unattacked -> verdict 1.0 -> hard
    assert _gradient_branch(internal=False) == ["hard"]


def test_view_gradient_branch_soft_when_defeated():
    # internal=True: exempt active -> block defeated (0.0) -> auth verdict ~0.33 -> soft
    g, auth, block, exempt = _build_access_control()
    branch = []

    def h(ctx, ev, view):
        node = view.get(short_of(auth))
        branch.append("hard" if (node and node.verdict >= 0.5) else "soft")

    exempt.on_active(h).on_defeated(h).on_inactive(h)

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", True)
    g.run(ctx)

    assert branch == ["soft"]


def test_returned_view_is_queryable_posthoc():
    g, auth, block, exempt = _build_access_control()

    ctx = MapContext()
    ctx.set("blocked", True)
    ctx.set("internal", False)
    result = g.run(ctx)

    view = result.view
    assert view is not None
    assert len(view.all()) == 3
    assert view.get(short_of(auth)).type == NodeEventType.DEFEATED
    assert view.get(short_of(block)).type == NodeEventType.ACTIVE
    assert [a.name for a in view.attackers(short_of(auth))] == [short_of(block)]
    assert view.get("nope") is None
